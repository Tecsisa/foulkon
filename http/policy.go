package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tecsisa/authorizr/authorizr"
)

type PolicyHandler struct {
	core *authorizr.Core
}

func (p *PolicyHandler) handleListPolicies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (p *PolicyHandler) handleCreatePolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (p *PolicyHandler) handleDeletePolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (p *PolicyHandler) handleGetPolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
