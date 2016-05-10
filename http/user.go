package http

import (
	"encoding/json"
	"net/http"

	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/authorizr"
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
	// Retrieve PathPrefix
	pathPrefix := r.URL.Query().Get("PathPrefix")

	// Call user API
	result, err := u.core.UserApi.GetListUsers(pathPrefix)
	if err != nil {
		u.core.Logger.Errorln(err)
		RespondInternalServerError(w)
		return
	}

	// Create response
	response := &GetUsersResponse{
		Users: result,
	}

	// Return data
	RespondOk(w, response)
}

// This method create the user passed by form request and return the user created
func (u *UserHandler) handlePostUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Decode request
	request := CreateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		u.core.Logger.Errorln(err)
		RespondBadRequest(w)
		return
	}

	// Check parameters
	if len(strings.TrimSpace(request.ExternalID)) == 0 ||
		len(strings.TrimSpace(request.Path)) == 0 {
		u.core.Logger.Errorf("There are mising parameters: ExternalID %v, Path %v", request.ExternalID, request.Path)
		RespondBadRequest(w)
		return
	}

	// Call user API to create an user
	result, err := u.core.UserApi.AddUser(createUserFromRequest(request))

	// Error handling
	if err != nil {
		u.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		if apiError.Code == api.USER_ALREADY_EXIST {
			RespondConflict(w)
			return
		} else { // Unexpected API error
			RespondInternalServerError(w)
			return
		}
	}

	response := &CreateUserResponse{
		User: result,
	}

	// Write user to response
	RespondOk(w, response)
}

func (u *UserHandler) handlePutUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//TODO: Unimplemented
}

func (u *UserHandler) handleGetUserId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Retrieve user id from path
	id := ps.ByName(USER_ID)

	// Call user API to retrieve user
	result, err := u.core.UserApi.GetUserByExternalId(id)

	// Error handling
	if err != nil {
		u.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		if apiError.Code == api.USER_BY_EXTERNAL_ID_NOT_FOUND {
			RespondNotFound(w)
			return
		} else { // Unexpected API error
			RespondInternalServerError(w)
			return
		}
	}

	response := GetUserByIdResponse{
		User: result,
	}

	// Write user to response
	RespondOk(w, response)
}

func (u *UserHandler) handleDeleteUserId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Retrieve user id from path
	id := ps.ByName(USER_ID)

	// Call user API to delete user
	err := u.core.UserApi.RemoveUserById(id)

	// Check if there were errors
	if err != nil {
		u.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		// If user doesn't exist
		if apiError.Code == api.USER_BY_EXTERNAL_ID_NOT_FOUND {
			RespondNotFound(w)
		} else { // Unexpected error
			RespondInternalServerError(w)
		}
	} else { // Respond without content
		RespondNoContent(w)
	}
}

func (u *UserHandler) handleUserIdGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Retrieve users using path
	id := ps.ByName(USER_ID)

	result, err := u.core.UserApi.GetGroupsByUserId(id)
	if err != nil {
		RespondInternalServerError(w)
	}
	b, err := json.Marshal(result)
	if err != nil {
		RespondInternalServerError(w)
	}

	w.Write(b)
}

func (u *UserHandler) handleOrgListUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//TODO: Unimplemented
}

func createUserFromRequest(request CreateUserRequest) api.User {
	path := request.Path + "/" + request.ExternalID
	urn := api.CreateUrn("", api.RESOURCE_USER, path)
	user := api.User{
		ID:         uuid.NewV4().String(),
		ExternalID: request.ExternalID,
		Path:       path,
		Urn:        urn,
	}

	return user
}
