package authorizr

import (
	"io"

	"errors"
	"os"
	"strings"

	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/auth"
	"github.com/tecsisa/authorizr/database/postgresql"
)

// Core is the manager of authorize. This use abstractions of connectors for backends,
// that you define at startup
type Core struct {
	// APIs
	UserApi   *api.UsersAPI
	GroupApi  *api.GroupsAPI
	PolicyApi *api.PolicyAPI

	// Logger
	Logger *log.Logger

	//  Auth connector
	Authenticator *auth.Authenticator
}

// Core config struct that manage configuration
type CoreConfig struct {
	LoggerConfig        map[string]string `json:"LoggerConfig"`
	DatabaseConfig      map[string]string `json:"DatabaseConfig"`
	AuthConnectorConfig map[string]string `json:"AuthConnectorConfig"`
	AdminUserConfig     map[string]string `json:"AdminUserConfig"`
}

// Create a Core using configuration values
func NewCore(coreconfig *CoreConfig) (*Core, error) {

	// Create logger
	var logOut io.Writer
	logOut = os.Stdout
	if coreconfig.LoggerConfig["log_type"] == "file" {
		logFileDir := coreconfig.LoggerConfig["log_file_dir"]
		file, err := os.OpenFile(logFileDir, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return nil, err
		}
		logOut = file
	}
	// Logger level
	loglevel := log.InfoLevel
	if coreconfig.LoggerConfig["log_level_debug"] == "true" {
		loglevel = log.DebugLevel
	}

	logger := &log.Logger{
		Out:       logOut,
		Formatter: &log.JSONFormatter{},
		Hooks:     make(log.LevelHooks),
		Level:     loglevel,
	}

	// Start DB with APIs
	var userApi *api.UsersAPI
	var groupApi *api.GroupsAPI

	switch coreconfig.DatabaseConfig["db_type"] {
	case "postgres": // Postgres DB
		logger.Info("Connecting to postgres database")
		db, err := postgresql.InitDb(coreconfig.DatabaseConfig["db_postgres_datasourcename"])
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		logger.Info("Connected to postgres database")
		userApi = &api.UsersAPI{
			UserRepo: postgresql.PostgresRepo{
				Dbmap: db,
			},
		}
		groupApi = &api.GroupsAPI{
			GroupRepo: postgresql.PostgresRepo{
				Dbmap: db,
			},
		}
	default:
		err := errors.New("Unexpected db_type value in configuration file (Maybe it is empty)")
		logger.Error(err)
		return nil, err
	}

	// Instantiate Auth Connector
	var authConnector auth.AuthConnector
	switch coreconfig.AuthConnectorConfig["auth_connector_type"] {
	case "oidc":
		issuer := coreconfig.AuthConnectorConfig["auth_connector_oidc_issuer"]
		clientsids := coreconfig.AuthConnectorConfig["auth_connector_oidc_client_ids"]
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

	adminUser := coreconfig.AdminUserConfig["auth_admin_username"]
	adminPassword := coreconfig.AdminUserConfig["auth_admin_password"]
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
		Logger:        logger,
		Authenticator: authenticator,
		UserApi:       userApi,
		GroupApi:      groupApi,
	}, nil
}
