package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/authorizr"
	"net/http"
	"strconv"
)

type UserHandler struct {
	core *authorizr.Core
}

// Requests
type CreateUserRequest struct {
	Name string
	Org  string
	Path string
}

// Responses

func (u *UserHandler) handleGetUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Retrieve org
	var org string
	queryorg := r.URL.Query().Get("Org")
	if queryorg == "" {
		authorizr.RespondError(w, http.StatusBadRequest, errors.New("Org missing"))
	} else {
		org = queryorg
	}

	result, err := u.core.GetUserAPI().GetListUsers(org, r.URL.Query().Get("PathPrefix"))
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
	request := CreateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		authorizr.RespondError(w, http.StatusBadRequest, err)
		return
	}

	err = u.core.GetUserAPI().AddUser(createUserFromRequest(request))
	if err != nil {
		authorizr.RespondError(w, http.StatusInternalServerError, err)
	} else {
		authorizr.RespondOk(w, http.StatusCreated)
	}
	return
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

func createUserFromRequest(request CreateUserRequest) api.User {
	urn := fmt.Sprintf("urn:iws:iam:%v:user/%v%v", request.Org, request.Path, request.Name)
	user := api.User{
		Name: request.Name,
		Path: request.Path,
		Urn:  urn,
		Org:  request.Org,
	}

	return user
}
