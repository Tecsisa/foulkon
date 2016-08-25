package foulkon

import (
	"io"
	"regexp"

	"errors"
	"os"
	"strings"

	"fmt"

	"database/sql"

	log "github.com/Sirupsen/logrus"
	"github.com/pelletier/go-toml"
	"github.com/tecsisa/foulkon/api"
	"github.com/tecsisa/foulkon/auth"
	"github.com/tecsisa/foulkon/database/postgresql"
)

// aux var for ${OS_ENV_VAR} regex
var rEnvVar, _ = regexp.Compile(`^\$\{(\w+)\}$`)
var db *sql.DB
var workerLogfile *os.File
var logger *log.Logger

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

	// Logger
	Logger *log.Logger

	//  Auth connector
	Authenticator *auth.Authenticator
}

// NewWorker creates a Worker using configuration values
func NewWorker(config *toml.TomlTree) (*Worker, error) {

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
		logOut = workerLogfile
	}
	// Logger level. Defaults to INFO
	loglevel, err := log.ParseLevel(getDefaultValue(config, "logger.level", "info"))
	if err != nil {
		loglevel = log.InfoLevel
	}

	logger = &log.Logger{
		Out:       logOut,
		Formatter: &log.JSONFormatter{},
		Hooks:     make(log.LevelHooks),
		Level:     loglevel,
	}
	logger.Infof("Logger type: %v, LogLevel: %v", loggerType, logger.Level.String())

	// Start DB with API
	var authApi api.AuthAPI

	dbType, err := getMandatoryValue(config, "database.type")
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	switch dbType {
	case "postgres": // PostgreSQL DB
		logger.Info("Connecting to postgres database")
		dbdsn, err := getMandatoryValue(config, "database.postgres.datasourcename")
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		gormDB, err := postgresql.InitDb(dbdsn,
			getDefaultValue(config, "database.postgres.idleconns", "5"),
			getDefaultValue(config, "database.postgres.maxopenconns", "20"),
			getDefaultValue(config, "database.postgres.connttl", "300"),
		)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		db = gormDB.DB()
		logger.Info("Connected to postgres database")

		// Create repository
		repoDB := postgresql.PostgresRepo{
			Dbmap: gormDB,
		}
		authApi = api.AuthAPI{
			GroupRepo:  repoDB,
			UserRepo:   repoDB,
			PolicyRepo: repoDB,
		}

	default:
		err := errors.New("Unexpected db_type value in configuration file (Maybe it is empty)")
		logger.Error(err)
		return nil, err
	}

	authApi.Logger = logger

	// Instantiate Auth Connector
	var authConnector auth.AuthConnector
	authType, err := getMandatoryValue(config, "authenticator.type")
	if err != nil {
		return nil, err
	}
	switch authType {
	case "oidc":
		issuer, err := getMandatoryValue(config, "authenticator.oidc.issuer")
		if err != nil {
			return nil, err
		}
		clientsids, err := getMandatoryValue(config, "authenticator.oidc.clientids")
		if err != nil {
			return nil, err
		}
		authOidcConnector, err := auth.InitOIDCConnector(logger, issuer, strings.Split(clientsids, ";"))
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		authConnector = authOidcConnector
		logger.Infof("OIDC connector configured for issuer %v", issuer)
	default:
		err := errors.New("Unexpected auth_connector_type value in configuration file (Maybe it is empty)")
		logger.Error(err)
		return nil, err
	}

	adminUser, err := getMandatoryValue(config, "admin.username")
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	adminPassword, err := getMandatoryValue(config, "admin.password")
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if len(strings.TrimSpace(adminUser)) < 1 || len(strings.TrimSpace(adminPassword)) < 1 {
		err := fmt.Errorf("Admin user config unexpected adminUser:%v, adminpassword:%v", adminUser, adminPassword)
		logger.Error(err)
		return nil, err
	}

	authenticator := auth.NewAuthenticator(authConnector, adminUser, adminPassword)
	logger.Infof("Created authenticator with admin username %v", adminUser)

	host, err := getMandatoryValue(config, "server.host")
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	port, err := getMandatoryValue(config, "server.port")
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &Worker{
		Host:          host,
		Port:          port,
		CertFile:      getDefaultValue(config, "server.certfile", ""),
		KeyFile:       getDefaultValue(config, "server.keyfile", ""),
		Logger:        logger,
		Authenticator: authenticator,
		UserApi:       authApi,
		GroupApi:      authApi,
		PolicyApi:     authApi,
		AuthzApi:      authApi,
	}, nil
}

func CloseWorker() int {
	status := 0
	if err := db.Close(); err != nil {
		logger.Errorf("Couldn't close DB connection: %v", err)
		status = 1
	}
	if err := workerLogfile.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't close logfile: %v", err)
		status = 1
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
