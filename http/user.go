package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"errors"

	"github.com/julienschmidt/httprouter"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/authorizr"
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
type CreateUserResponse struct {
	User *api.User
}

type GetUsersResponse struct {
	users []api.User
}

// This method return a list of users that belongs to Org param and have PathPrefix
func (u *UserHandler) handleGetUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u.logRequest(r)

	// Retrieve org
	var org string
	queryorg := r.URL.Query().Get("Org")
	if queryorg == "" {
		u.core.RespondError(w, http.StatusBadRequest, errors.New("Org missing"))
		return
	} else {
		org = queryorg
	}

	// Call user API
	result, err := u.core.GetUserAPI().GetListUsers(org, r.URL.Query().Get("PathPrefix"))
	if err != nil {
		u.core.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	// Check if there are results
	if result == nil {
		u.core.RespondOk(w, http.StatusNotFound)
		return
	}

	// Create response
	response := &GetUsersResponse{
		users: result,
	}

	b, err := json.Marshal(response)
	if err != nil {
		u.core.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	// Return data
	u.core.RespondOk(w, http.StatusOK)
	w.Write(b)
}

// This method create the user passed by form request and return the user created
func (u *UserHandler) handlePostUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u.logRequest(r)

	// Decode request
	request := CreateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		u.core.RespondError(w, http.StatusBadRequest, err)
		return
	}

	// Call user API to create an user
	result, err := u.core.GetUserAPI().AddUser(createUserFromRequest(request))

	// Check response
	if err != nil {
		u.core.RespondError(w, http.StatusInternalServerError, err)
	}
	response := &CreateUserResponse{
		User: result,
	}

	// Write user to response
	b, err := json.Marshal(response)
	if err != nil {
		u.core.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	// Return data
	u.core.RespondOk(w, http.StatusOK)
	w.Write(b)
}

func (u *UserHandler) handleGetUserId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u.logRequest(r)

	// Retrieve user id from path
	id, err := strconv.ParseUint(ps.ByName("id"), 10, 64)
	if err != nil {
		u.core.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	// Call user API to retrieve user
	result, err := u.core.GetUserAPI().GetUserById(id)

	// Check if there were errors
	if err != nil {
		u.core.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	// Check if there are results
	if result == nil {
		u.core.RespondOk(w, http.StatusNotFound)
		return
	}

	// Write user to response
	b, err := json.Marshal(result)
	if err != nil {
		u.core.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	// Return data
	u.core.RespondOk(w, http.StatusOK)
	w.Write(b)
}

func (u *UserHandler) handleDeleteUserId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u.logRequest(r)

	// Retrieve user id from path
	id, err := strconv.ParseUint(ps.ByName("id"), 10, 64)
	if err != nil {
		u.core.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	// Call user API to delete user
	err = u.core.GetUserAPI().RemoveUserById(id)

	// Check if there were errors
	if err != nil {
		u.core.RespondError(w, http.StatusBadRequest, err)
	} else {
		u.core.RespondOk(w, http.StatusCreated)
	}
}

func (u *UserHandler) handleUserIdGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u.logRequest(r)

	// Retrieve users using path
	id, err := strconv.ParseUint(ps.ByName("id"), 10, 64)
	if err != nil {
		u.core.RespondError(w, http.StatusInternalServerError, err)
	}
	result, err := u.core.GetUserAPI().GetGroupsByUserId(id)
	if err != nil {
		u.core.RespondError(w, http.StatusInternalServerError, err)
	}
	b, err := json.Marshal(result)
	if err != nil {
		u.core.RespondError(w, http.StatusBadRequest, err)
	}

	w.Write(b)
}

func createUserFromRequest(request CreateUserRequest) api.User {
	path := request.Path + "/" + request.Name
	urn := fmt.Sprintf("urn:iws:iam:%v:user/%v", request.Org, path)
	user := api.User{
		Name: request.Name,
		Path: path,
		Urn:  urn,
		Org:  request.Org,
	}

	return user
}

func (u *UserHandler) logRequest(request *http.Request) {
	u.core.Logger.Infoln(request)
}
