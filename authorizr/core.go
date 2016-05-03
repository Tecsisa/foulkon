package authorizr

import (
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/auth"
	"github.com/tecsisa/authorizr/database/postgresql"
)

// Core is the manager of authorizR. This use abstractions of connectors for backends,
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

type CoreConfig struct {
	LogFile        io.Writer
	DatasourceName string
}

func NewCore(coreconfig *CoreConfig) (*Core, error) {

	// Create logger
	logger := &log.Logger{
		Out:       coreconfig.LogFile,
		Formatter: &log.JSONFormatter{},
		Hooks:     make(log.LevelHooks),
		Level:     log.InfoLevel,
	}

	logger.Info("Accesing to db with DSN " + coreconfig.DatasourceName)
	// Start DB
	db, err := postgresql.InitDb(coreconfig.DatasourceName)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	// Instantiate Auth Connector
	authConnector, err := auth.InitOIDCConnector()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	authenticator := auth.NewAuthenticator(authConnector, "authorizr")

	// Instantiate APIs
	userApi := &api.UsersAPI{
		UserRepo: postgresql.PostgresRepo{
			Dbmap: db,
		},
	}

	groupApi := &api.GroupsAPI{
		GroupRepo: postgresql.PostgresRepo{
			Dbmap: db,
		},
	}

	return &Core{
		Logger:        logger,
		Authenticator: authenticator,
		UserApi:       userApi,
		GroupApi:      groupApi,
		Logger:        logger,
	}, nil
}
