package foulkon

import (
	"io"
	"regexp"

	"errors"
	"os"
	"strings"

	"fmt"

	"database/sql"

	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/database/postgresql"
	"github.com/Tecsisa/foulkon/middleware"
	"github.com/Tecsisa/foulkon/middleware/auth"
	"github.com/Tecsisa/foulkon/middleware/auth/oidc"
	"github.com/Tecsisa/foulkon/middleware/logger"
	"github.com/Tecsisa/foulkon/middleware/xrequestid"
	"github.com/pelletier/go-toml"
)

const (
	FOULKON_VERSION = "v0.3.0-SNAPSHOT"
)

// aux var for ${OS_ENV_VAR} regex
var rEnvVar, _ = regexp.Compile(`^\$\{(\w+)\}$`)
var db *sql.DB
var workerLogfile *os.File

// Worker is the Authorization server.
type Worker struct {
	// Server config
	Host string
	Port string

	// TLS configuration
	CertFile string
	KeyFile  string

	// APIs
	UserApi   api.UserAPI
	GroupApi  api.GroupAPI
	PolicyApi api.PolicyAPI
	AuthzApi  api.AuthzAPI
	ProxyApi  api.ProxyResourcesAPI

	//  Middleware handler
	MiddlewareHandler *middleware.MiddlewareHandler

	// Current Foulkon configuration
	Config WorkerConfig
}

// WorkerConfig
type WorkerConfig struct {
	// Logger Config
	LoggerType    string
	LoggerLevel   string
	FileDirectory string

	// Database Config
	DBType       string
	IdleConns    int
	MaxOpenConns int
	ConnTtl      int

	// Authenticator Config
	AuthType string
	Issuer   string

	Version string
}

// NewWorker creates a Worker using configuration values
func NewWorker(config *toml.TomlTree) (*Worker, error) {
	var wc WorkerConfig

	// Create logger
	var logOut io.Writer
	var err error
	logOut = os.Stdout
	loggerType := getDefaultValue(config, "logger.type", "Stdout")
	if loggerType == "file" {
		logFileDir := getDefaultValue(config, "logger.file.dir", "/tmp/foulkon.log")
		workerLogfile, err = os.OpenFile(logFileDir, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return nil, err
		}
		wc.FileDirectory = logFileDir
		logOut = workerLogfile
	}
	wc.LoggerType = loggerType

	// Logger level. Defaults to INFO
	loglevel, err := logrus.ParseLevel(getDefaultValue(config, "logger.level", "info"))
	if err != nil {
		loglevel = logrus.InfoLevel
	}
	wc.LoggerLevel = loglevel.String()

	api.Log = &logrus.Logger{
		Out:       logOut,
		Formatter: &logrus.JSONFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     loglevel,
	}
	api.Log.Infof("Logger type: %v, LogLevel: %v", loggerType, api.Log.Level.String())

	// Start DB with API
	var authApi api.WorkerAPI

	dbType, err := getMandatoryValue(config, "database.type")
	if err != nil {
		api.Log.Error(err)
		return nil, err
	}
	wc.DBType = dbType

	switch dbType {
	case "postgres": // PostgreSQL DB
		api.Log.Info("Connecting to postgres database")
		dbdsn, err := getMandatoryValue(config, "database.postgres.datasourcename")
		if err != nil {
			api.Log.Error(err)
			return nil, err
		}
		dbIdleconns := getDefaultValue(config, "database.postgres.idleconns", "5")
		dbMaxopenconns := getDefaultValue(config, "database.postgres.maxopenconns", "20")
		dbConttl := getDefaultValue(config, "database.postgres.connttl", "300")
		gormDB, err := postgresql.InitDb(dbdsn, dbIdleconns, dbMaxopenconns, dbConttl)
		if err != nil {
			api.Log.Error(err)
			return nil, err
		}
		db = gormDB.DB()
		api.Log.Info("Connected to postgres database")

		// Create repository
		repoDB := postgresql.PostgresRepo{
			Dbmap: gormDB,
		}
		authApi = api.WorkerAPI{
			GroupRepo:  repoDB,
			UserRepo:   repoDB,
			PolicyRepo: repoDB,
		}
		wc.IdleConns, _ = strconv.Atoi(dbIdleconns)
		wc.MaxOpenConns, _ = strconv.Atoi(dbMaxopenconns)
		wc.ConnTtl, _ = strconv.Atoi(dbConttl)

	default:
		err := errors.New("Unexpected db_type value in configuration file (Maybe it is empty)")
		api.Log.Error(err)
		return nil, err
	}

	// Instantiate Auth Connector
	var authConnector auth.AuthConnector
	authType, err := getMandatoryValue(config, "authenticator.type")
	if err != nil {
		return nil, err
	}
	wc.AuthType = authType

	switch authType {
	case "oidc":
		issuer, err := getMandatoryValue(config, "authenticator.oidc.issuer")
		if err != nil {
			return nil, err
		}
		wc.Issuer = issuer

		clientsids, err := getMandatoryValue(config, "authenticator.oidc.clientids")
		if err != nil {
			return nil, err
		}
		authOidcConnector, err := oidc.InitOIDCConnector(issuer, strings.Split(clientsids, ";"))
		if err != nil {
			api.Log.Error(err)
			return nil, err
		}
		authConnector = authOidcConnector
		api.Log.Infof("OIDC connector configured for issuer %v", issuer)
	default:
		err := errors.New("Unexpected auth_connector_type value in configuration file (Maybe it is empty)")
		api.Log.Error(err)
		return nil, err
	}

	adminUser, err := getMandatoryValue(config, "admin.username")
	if err != nil {
		api.Log.Error(err)
		return nil, err
	}
	adminPassword, err := getMandatoryValue(config, "admin.password")
	if err != nil {
		api.Log.Error(err)
		return nil, err
	}
	if len(strings.TrimSpace(adminUser)) < 1 || len(strings.TrimSpace(adminPassword)) < 1 {
		err := fmt.Errorf("Admin user config unexpected adminUser:%v, adminpassword:%v", adminUser, adminPassword)
		api.Log.Error(err)
		return nil, err
	}

	// Middlewares
	middlewares := make(map[string]middleware.Middleware)

	// Authenticator middleware
	authenticatorMiddleware := auth.NewAuthenticatorMiddleware(authConnector, adminUser, adminPassword)
	middlewares[middleware.AUTHENTICATOR_MIDDLEWARE] = authenticatorMiddleware
	api.Log.Infof("Created authenticator with admin username %v", adminUser)

	// X-Request-Id middleware
	xrequestidMiddleware := xrequestid.NewXRequestIdMiddleware()
	middlewares[middleware.XREQUESTID_MIDDLEWARE] = xrequestidMiddleware

	// Request Logger middleware
	requestLoggerMiddleware := logger.NewRequestLoggerMiddleware()
	middlewares[middleware.REQUEST_LOGGER_MIDDLEWARE] = requestLoggerMiddleware

	host, err := getMandatoryValue(config, "server.host")
	if err != nil {
		api.Log.Error(err)
		return nil, err
	}
	port, err := getMandatoryValue(config, "server.port")
	if err != nil {
		api.Log.Error(err)
		return nil, err
	}

	wc.Version = FOULKON_VERSION

	return &Worker{
		Host:              host,
		Port:              port,
		CertFile:          getDefaultValue(config, "server.certfile", ""),
		KeyFile:           getDefaultValue(config, "server.keyfile", ""),
		MiddlewareHandler: &middleware.MiddlewareHandler{Middlewares: middlewares},
		UserApi:           authApi,
		GroupApi:          authApi,
		PolicyApi:         authApi,
		AuthzApi:          authApi,
		Config:            wc,
	}, nil
}

func CloseWorker() int {
	status := 0
	if err := db.Close(); err != nil {
		api.Log.Errorf("Couldn't close DB connection: %v", err)
		status = 1
	}
	if workerLogfile != nil {
		if err := workerLogfile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't close logfile: %v", err)
			status = 1
		}
	}
	return status
}

// This aux method returns mandatory config value or any error occurred
func getMandatoryValue(config *toml.TomlTree, key string) (string, error) {
	if !config.Has(key) {
		return "", fmt.Errorf("Cannot retrieve configuration value %v", key)
	}

	value := getVar(config, key)
	if value == "" {
		return "", fmt.Errorf("Cannot retrieve configuration value %v", key)
	}
	return value, nil
}

// This aux method returns a value if defined in config file. Else, returns default value
func getDefaultValue(config *toml.TomlTree, key string, def string) string {
	value := def
	if config.Has(key) {
		value = getVar(config, key)
	} else {
		fmt.Fprintf(os.Stdout, "WARN: using default value in key %v - default value: %v\n", key, def)
	}
	return value
}

// Check variables in TOML file.
// If the value of a key is '${SOME_KEY}', we will search the value in the OS ENV vars
// If the value of a key is 'something_else', returns that as the value
func getVar(config *toml.TomlTree, key string) string {
	value := config.Get(key).(string)
	match := rEnvVar.FindStringSubmatch(value)
	if match != nil && len(match) > 1 {
		if match[1] != "" {
			return os.Getenv(match[1])
		}
	}
	return value
}
