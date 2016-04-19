package authorizr

import (
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/database/postgresql"
	"log"
	"net/http"
)

// Core is the manager of authorizR. This use abstractions of connectors for backends,
// that you define at startup
type Core struct {
	userapi   *api.UsersAPI
	groupapi  *api.GroupAPI
	policyapi *api.PolicyAPI
}

func NewCore() (*Core, error) {
	db := postgresql.InitDb("/tmp/authorizer/sql.bin")

	userapiimp := &api.UsersAPI{
		UserRepo: postgresql.PostgresRepo{
			Dbmap: db,
		},
	}
	return &Core{
		userapi: userapiimp,
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

func RespondError(w http.ResponseWriter, status int, err error) {
	log.Println("Error received ", err)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
}

func RespondOk(w http.ResponseWriter, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
}
