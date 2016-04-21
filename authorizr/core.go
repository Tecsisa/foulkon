package authorizr

import (
	"net/http"

	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/database/postgresql"
)

// Core is the manager of authorizR. This use abstractions of connectors for backends,
// that you define at startup
type Core struct {
	// APIs
	userapi   *api.UsersAPI
	groupapi  *api.GroupAPI
	policyapi *api.PolicyAPI

	// Logger
	Logger *log.Logger
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

	// Instantiate APIs
	userApi := &api.UsersAPI{
		UserRepo: postgresql.PostgresRepo{
			Dbmap: db,
		},
	}

	return &Core{
		userapi: userApi,
		Logger:  logger,
	}, nil
}

func (core *Core) GetUserAPI() *api.UsersAPI {
	return core.userapi
}

func (core *Core) GetGroupAPI() *api.GroupAPI {
	return core.groupapi
}

func (core *Core) GetPolicyAPI() *api.PolicyAPI {
	return core.policyapi
}

func (core *Core) RespondError(w http.ResponseWriter, status int, err error) {
	core.Logger.Error(err)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
}

func (core *Core) RespondOk(w http.ResponseWriter, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
}
