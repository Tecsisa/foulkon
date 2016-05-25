package http

import (
	"encoding/json"
	"net/http"

	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/tecsisa/authorizr/api"
)

// Requests

type CreateUserRequest struct {
	ExternalID string `json:"ExternalID, omitempty"`
	Path       string `json:"Path, omitempty"`
}

type UpdateUserRequest struct {
	Path string `json:"Path, omitempty"`
}

// Responses

type CreateUserResponse struct {
	User *api.User
}

type UpdateUserResponse struct {
	User *api.User
}

type GetUsersResponse struct {
	Users []api.User
}

type GetUserByIdResponse struct {
	User *api.User
}

type GetGroupsByUserIdResponse struct {
	Groups []api.Group
}

// This method returns a list of users that belongs to Org param and have PathPrefix
func (a *AuthHandler) handleGetUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Retrieve PathPrefix
	pathPrefix := r.URL.Query().Get("PathPrefix")

	// Call user API
	result, err := a.core.AuthApi.GetListUsers(a.core.Authenticator.RetrieveUserID(*r), pathPrefix)
	if err != nil {
		a.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		switch apiError.Code {
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			RespondForbidden(w)
		default: // Unexpected API error
			RespondInternalServerError(w)
		}
		return
	}

	// Create response
	response := &GetUsersResponse{
		Users: result,
	}

	// Return data
	RespondOk(w, response)
}

// This method creates the user passed by form request and return the user created
func (a *AuthHandler) handlePostUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Decode request
	request := CreateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		a.core.Logger.Errorln(err)
		RespondBadRequest(w)
		return
	}

	// Call user API to create an user
	result, err := a.core.AuthApi.AddUser(a.core.Authenticator.RetrieveUserID(*r), request.ExternalID, request.Path)

	// Error handling
	if err != nil {
		a.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		switch apiError.Code {
		case api.USER_ALREADY_EXIST:
			RespondConflict(w)
		case api.INVALID_PARAMETER_ERROR:
			RespondBadRequest(w)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			RespondForbidden(w)
		default: // Unexpected API error
			RespondInternalServerError(w)
		}
		return
	}

	response := &CreateUserResponse{
		User: result,
	}

	// Write user to response
	RespondOk(w, response)
}

func (a *AuthHandler) handlePutUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Decode request
	request := UpdateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		a.core.Logger.Errorln(err)
		RespondBadRequest(w)
		return
	}

	// Check parameters
	if len(strings.TrimSpace(request.Path)) == 0 {
		a.core.Logger.Errorf("There are mising parameters: Path %v", request.Path)
		RespondBadRequest(w)
		return
	}

	// Retrieve user id from path
	id := ps.ByName(USER_ID)

	// Call user API to update user
	result, err := a.core.AuthApi.UpdateUser(a.core.Authenticator.RetrieveUserID(*r), id, request.Path)

	// Error handling
	if err != nil {
		a.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		switch apiError.Code {
		case api.USER_BY_EXTERNAL_ID_NOT_FOUND:
			RespondNotFound(w)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			RespondForbidden(w)
		case api.INVALID_PARAMETER_ERROR:
			RespondBadRequest(w)
		default: // Unexpected API error
			RespondInternalServerError(w)
		}
		return
	}

	// Create response
	response := &UpdateUserResponse{
		User: result,
	}

	// Write user to response
	RespondOk(w, response)
}

func (a *AuthHandler) handleGetUserId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Retrieve user id from path
	id := ps.ByName(USER_ID)

	// Call user API to retrieve user
	result, err := a.core.AuthApi.GetUserByExternalId(a.core.Authenticator.RetrieveUserID(*r), id)

	// Error handling
	if err != nil {
		a.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		switch apiError.Code {
		case api.USER_BY_EXTERNAL_ID_NOT_FOUND:
			RespondNotFound(w)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			RespondForbidden(w)
		default: // Unexpected API error
			RespondInternalServerError(w)
		}
		return
	}

	response := GetUserByIdResponse{
		User: result,
	}

	// Write user to response
	RespondOk(w, response)
}

func (a *AuthHandler) handleDeleteUserId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Retrieve user id from path
	id := ps.ByName(USER_ID)

	// Call user API to delete user
	err := a.core.AuthApi.RemoveUserById(a.core.Authenticator.RetrieveUserID(*r), id)

	if err != nil {
		a.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		switch apiError.Code {
		case api.USER_BY_EXTERNAL_ID_NOT_FOUND:
			RespondNotFound(w)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			RespondForbidden(w)
		default: // Unexpected API error
			RespondInternalServerError(w)
		}
		return
	}

	RespondNoContent(w)
}

func (a *AuthHandler) handleUserIdGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Retrieve users using path
	id := ps.ByName(USER_ID)

	result, err := a.core.AuthApi.GetGroupsByUserId(a.core.Authenticator.RetrieveUserID(*r), id)

	if err != nil {
		a.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		switch apiError.Code {
		case api.USER_BY_EXTERNAL_ID_NOT_FOUND:
			RespondNotFound(w)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			RespondForbidden(w)
		default: // Unexpected API error
			RespondInternalServerError(w)
		}
		return
	}

	response := GetGroupsByUserIdResponse{
		Groups: result,
	}

	// Write user to response
	RespondOk(w, response)
}

func (a *AuthHandler) handleOrgListUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//TODO: Unimplemented
}
