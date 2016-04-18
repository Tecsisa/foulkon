package http

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/authorizr"
	"net/http"
	"strconv"
)

type UserHandler struct {
	core *authorizr.Core
}

func (u *UserHandler) handleGetUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Retrieve users using path
	result, err := u.core.GetUserAPI().GetListUsers("/mipath")
	if err != nil {
		authorizr.RespondError(w, http.StatusInternalServerError, err)
	}
	b, err := json.Marshal(result)
	if err != nil {
		authorizr.RespondError(w, http.StatusBadRequest, err)
	}

	w.Write(b)
}

func (u *UserHandler) handlePostUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Add an user
	err := u.core.GetUserAPI().AddUser(api.User{})
	if err != nil {
		authorizr.RespondError(w, http.StatusBadRequest, err)
	} else {
		authorizr.RespondOk(w, http.StatusCreated)
	}
}

func (u *UserHandler) handleGetUserId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Retrieve users using path
	id, err := strconv.ParseUint(ps.ByName("id"), 10, 64)
	if err != nil {
		authorizr.RespondError(w, http.StatusInternalServerError, err)
	}
	result, err := u.core.GetUserAPI().GetUserById(id)
	if err != nil {
		authorizr.RespondError(w, http.StatusInternalServerError, err)
	}
	b, err := json.Marshal(result)
	if err != nil {
		authorizr.RespondError(w, http.StatusBadRequest, err)
	}

	w.Write(b)
}

func (u *UserHandler) handleDeleteUserId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Add an user
	id, err := strconv.ParseUint(ps.ByName("id"), 10, 64)
	if err != nil {
		authorizr.RespondError(w, http.StatusInternalServerError, err)
	}
	err = u.core.GetUserAPI().RemoveUserById(id)
	if err != nil {
		authorizr.RespondError(w, http.StatusBadRequest, err)
	} else {
		authorizr.RespondOk(w, http.StatusCreated)
	}
}

func (u *UserHandler) handleUserIdGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Retrieve users using path
	id, err := strconv.ParseUint(ps.ByName("id"), 10, 64)
	if err != nil {
		authorizr.RespondError(w, http.StatusInternalServerError, err)
	}
	result, err := u.core.GetUserAPI().GetGroupsByUserId(id)
	if err != nil {
		authorizr.RespondError(w, http.StatusInternalServerError, err)
	}
	b, err := json.Marshal(result)
	if err != nil {
		authorizr.RespondError(w, http.StatusBadRequest, err)
	}

	w.Write(b)
}
