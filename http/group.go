package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tecsisa/authorizr/authorizr"
)

type GroupHandler struct {
	core *authorizr.Core
}

func (g *GroupHandler) handleCreateGroup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func (g *GroupHandler) handleDeleteGroup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func (g *GroupHandler) handleGetGroup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func (g *GroupHandler) handleListGroups(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func (g *GroupHandler) handleUpdateGroup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func (g *GroupHandler) handleListMembers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func (g *GroupHandler) handleAddMember(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func (g *GroupHandler) handleRemoveMember(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func (g *GroupHandler) handleAttachGroupPolicy(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func (g *GroupHandler) handleDetachGroupPolicy(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func (g *GroupHandler) handleListAtachhedGroupPolicies(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}
