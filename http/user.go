package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/authorizr"
	"strings"
)

type UserHandler struct {
	core *authorizr.Core
}

// Requests
type CreateUserRequest struct {
	ExternalID string `json:"ExternalID, omitempty"`
	Path       string `json:"Path, omitempty"`
}

// Responses
type CreateUserResponse struct {
	User *api.User
}

type GetUsersResponse struct {
	Users []api.User
}

type GetUserByIdResponse struct {
	User *api.User
}

// This method return a list of users that belongs to Org param and have PathPrefix
func (u *UserHandler) handleGetUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u.logRequest(r)

	// Retrieve PathPrefix
	pathPrefix := r.URL.Query().Get("PathPrefix")

	// Call user API
	result, err := u.core.Userapi.GetListUsers(pathPrefix)
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
		Users: result,
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

	// Check parameters
	if len(strings.TrimSpace(request.ExternalID)) == 0 ||
		len(strings.TrimSpace(request.Path)) == 0 {
		u.core.RespondError(w, http.StatusBadRequest, err)
		return
	}

	// Call user API to create an user
	result, err := u.core.Userapi.AddUser(createUserFromRequest(request))

	// Check response
	if err != nil {
		u.core.RespondError(w, http.StatusInternalServerError, err)
		return
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

func (u *UserHandler) handlePutUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//TODO: Unimplemented
}

func (u *UserHandler) handleGetUserId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u.logRequest(r)

	// Retrieve user id from path
	id := ps.ByName(USER_ID)

	// Call user API to retrieve user
	result, err := u.core.Userapi.GetUserById(id)

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

	response := GetUserByIdResponse{
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

func (u *UserHandler) handleDeleteUserId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u.logRequest(r)

	// Retrieve user id from path
	id := ps.ByName(USER_ID)

	// Call user API to delete user
	err := u.core.Userapi.RemoveUserById(id)

	// Check if there were errors
	if err != nil {
		u.core.RespondError(w, http.StatusBadRequest, err)
	} else {
		u.core.RespondOk(w, http.StatusAccepted)
	}
}

func (u *UserHandler) handleUserIdGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u.logRequest(r)

	// Retrieve users using path
	id := ps.ByName(USER_ID)

	result, err := u.core.Userapi.GetGroupsByUserId(id)
	if err != nil {
		u.core.RespondError(w, http.StatusInternalServerError, err)
	}
	b, err := json.Marshal(result)
	if err != nil {
		u.core.RespondError(w, http.StatusBadRequest, err)
	}

	w.Write(b)
}

func (u *UserHandler) handleOrgListUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//TODO: Unimplemented
}

func createUserFromRequest(request CreateUserRequest) api.User {
	path := request.Path + "/" + request.ExternalID
	urn := fmt.Sprintf("urn:iws:iam:user/%v", path)
	user := api.User{
		ExternalID: request.ExternalID,
		Path:       path,
		Urn:        urn,
	}

	return user
}

func (u *UserHandler) logRequest(request *http.Request) {
	u.core.Logger.Debugln(request)
}
