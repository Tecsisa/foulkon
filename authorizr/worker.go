package authorizr

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
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/auth"
	"github.com/tecsisa/authorizr/database/postgresql"
)

// aux var for ${OS_ENV_VAR} regex
var rEnvVar, _ = regexp.Compile(`^\$\{(\w+)\}$`)
var db *sql.DB
var worker_logfile *os.File
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

// Create a Worker using configuration values
func NewWorker(config *toml.TomlTree) (*Worker, error) {

	// Create logger
	var logOut io.Writer
	var err error
	logOut = os.Stdout
	loggerType := getDefaultValue(config, "logger.type", "Stdout")
	if loggerType == "file" {
		logFileDir := getDefaultValue(config, "logger.file.dir", "/tmp/authorizr.log")
		worker_logfile, err = os.OpenFile(logFileDir, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return nil, err
		}
		logOut = worker_logfile
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

	switch getMandatoryValue(config, "database.type") {
	case "postgres": // PostgreSQL DB
		logger.Info("Connecting to postgres database")
		gormDB, err := postgresql.InitDb(getMandatoryValue(config, "database.postgres.datasourcename"),
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

	authApi.Logger = *logger

	// Instantiate Auth Connector
	var authConnector auth.AuthConnector
	switch getMandatoryValue(config, "authenticator.type") {
	case "oidc":
		issuer := getMandatoryValue(config, "authenticator.oidc.issuer")
		clientsids := getMandatoryValue(config, "authenticator.oidc.clientids")
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

	adminUser := getMandatoryValue(config, "admin.username")
	adminPassword := getMandatoryValue(config, "admin.password")
	if len(strings.TrimSpace(adminUser)) < 1 || len(strings.TrimSpace(adminPassword)) < 1 {
		err := errors.New(fmt.Sprintf("Admin user config unexpected adminUser:%v, adminpassword:%v", adminUser, adminPassword))
		logger.Error(err)
		return nil, err
	}

	authenticator := auth.NewAuthenticator(authConnector, adminUser, adminPassword)
	logger.Infof("Created authenticator with admin username %v", adminUser)

	return &Worker{
		Host:          getMandatoryValue(config, "server.host"),
		Port:          getMandatoryValue(config, "server.port"),
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

func CloseWorker() {
	status := 0
	if err := db.Close(); err != nil {
		logger.Errorf("Couldn't close DB connection: %v", err)
		status = 1
	}
	if err := worker_logfile.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't close logfile: %v", err)
		status = 1
	}
	os.Exit(status)
}

// This aux method returns mandatory config value or finishes program execution
func getMandatoryValue(config *toml.TomlTree, key string) string {
	if !config.Has(key) {
		fmt.Fprintf(os.Stderr, "Cannot retrieve configuration value %v", key)
		os.Exit(1)
		return ""
	} else {
		value := getVar(config, key)
		if value == "" {
			fmt.Fprintf(os.Stderr, "Cannot retrieve configuration value %v", key)
			os.Exit(1)
		}
		return value
	}

}

// This aux method returns a value if defined in config file. Else, returns default value
func getDefaultValue(config *toml.TomlTree, key string, def string) string {
	value := def
	if config.Has(key) {
		value = getVar(config, key)
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
