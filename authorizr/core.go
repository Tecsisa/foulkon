package authorizr

import (
	"io"
	"regexp"

	"errors"
	"os"
	"strings"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/pelletier/go-toml"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/auth"
	"github.com/tecsisa/authorizr/database/postgresql"
)

// Core is the manager of authorize. This use abstractions of connectors for backends,
// that you define at startup
type Core struct {
	// Server config
	Host string
	Port string

	// API
	AuthApi *api.AuthAPI

	// Logger
	Logger *log.Logger

	//  Auth connector
	Authenticator *auth.Authenticator
}

// Create a Core using configuration values
func NewCore(config *toml.TomlTree) (*Core, error) {

	// Create logger
	var logOut io.Writer
	logOut = os.Stdout
	loggerType := getDefaultValue(config, "logger.type", "Stdout")
	if loggerType == "file" {
		logFileDir := getDefaultValue(config, "logger.file.dir", "/tmp/authorizr.log")
		file, err := os.OpenFile(logFileDir, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return nil, err
		}
		logOut = file
	}
	// Logger level
	loglevel, err := log.ParseLevel(getDefaultValue(config, "logger.level", ""))
	if err != nil {
		loglevel = log.InfoLevel
	}

	logger := &log.Logger{
		Out:       logOut,
		Formatter: &log.JSONFormatter{},
		Hooks:     make(log.LevelHooks),
		Level:     loglevel,
	}
	logger.Infof("Logger type: %v, LogLevel: %v", loggerType, log.GetLevel().String())

	// Start DB with API
	var authApi *api.AuthAPI

	switch getMandatoryValue(config, "database.type") {
	case "postgres": // Postgres DB
		logger.Info("Connecting to postgres database")
		db, err := postgresql.InitDb(getMandatoryValue(config, "database.postgres.datasourcename"))
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		logger.Info("Connected to postgres database")

		// Create repository
		repoDB := postgresql.PostgresRepo{
			Dbmap: db,
		}
		authApi = &api.AuthAPI{
			GroupRepo:  repoDB,
			UserRepo:   repoDB,
			PolicyRepo: repoDB,
		}

	default:
		err := errors.New("Unexpected db_type value in configuration file (Maybe it is empty)")
		logger.Error(err)
		return nil, err
	}

	// Instantiate Auth Connector
	var authConnector auth.AuthConnector
	switch getMandatoryValue(config, "authenticator.type") {
	case "oidc":
		issuer := getMandatoryValue(config, "authenticator.oidc.issuer")
		clientsids := getMandatoryValue(config, "authenticator.oidc.clientids")
		authOidcConnector, err := auth.InitOIDCConnector(issuer, strings.Split(clientsids, ";"))
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

	// Create authenticator
	authenticator := auth.NewAuthenticator(authConnector, adminUser, adminPassword)
	logger.Infof("Created authenticator with admin user name %v", adminUser)

	// Return created core
	return &Core{
		Host:          getMandatoryValue(config, "server.host"),
		Port:          getMandatoryValue(config, "server.port"),
		Logger:        logger,
		Authenticator: authenticator,
		AuthApi:       authApi,
	}, nil
}

// This aux method returns mandatory config value or finish program execution
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

// This aux method returns a value if exist or default value
func getDefaultValue(config *toml.TomlTree, key string, def string) string {
	value := def
	if config.Has(key) {
		value = getVar(config, key)
	}
	return value
}

// Check variables in TOML file.
// If the value of a key is '${SOMETHING}', we will search the value in the OS ENV vars
// If the value of a key is something else, return that value
func getVar(config *toml.TomlTree, key string) string {
	r, err := regexp.Compile(`^\$\{(\w+)\}$`)
	if err != nil {
		fmt.Printf("Regexp compilation error.\n")
		os.Exit(1)
		return ""
	}

	value := config.Get(key).(string)
	match := r.FindStringSubmatch(value)
	if match != nil && len(match) > 1 {
		if match[1] != "" {
			return os.Getenv(match[1])
		}
	}
	return value

}
