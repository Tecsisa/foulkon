package http

import (
	"encoding/json"
	"net/http"

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

type GetUserExternalIDsResponse struct {
	ExternalIDs []string
}

type GetUserByIdResponse struct {
	User *api.User
}

type GetGroupsByUserIdResponse struct {
	Groups []api.GroupIdentity
}

// This method returns a list of users that belongs to Org param and have PathPrefix
func (a *AuthHandler) handleGetUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Retrieve PathPrefix
	pathPrefix := r.URL.Query().Get("PathPrefix")

	// Call user API
	userID := a.core.Authenticator.RetrieveUserID(*r)
	result, err := a.core.AuthApi.GetListUsers(userID, pathPrefix)
	if err != nil {
		a.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		switch apiError.Code {
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &userID, w, apiError)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &userID, w)
		}
		return
	}

	// Create response
	response := &GetUserExternalIDsResponse{
		ExternalIDs: result,
	}

	// Return data
	a.RespondOk(r, &userID, w, response)
}

// This method creates the user passed by form request and return the user created
func (a *AuthHandler) handlePostUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := a.core.Authenticator.RetrieveUserID(*r)
	// Decode request
	request := CreateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		a.core.Logger.Errorln(err)
		a.RespondBadRequest(r, &userID, w, &api.Error{Code: api.INVALID_PARAMETER_ERROR, Message: err.Error()})
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
			a.RespondConflict(r, &userID, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &userID, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &userID, w, apiError)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &userID, w)
		}
		return
	}

	response := &CreateUserResponse{
		User: result,
	}

	// Write user to response
	a.RespondOk(r, &userID, w, response)
}

func (a *AuthHandler) handlePutUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := a.core.Authenticator.RetrieveUserID(*r)
	// Decode request
	request := UpdateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		a.core.Logger.Errorln(err)
		a.RespondBadRequest(r, &userID, w, &api.Error{Code: api.INVALID_PARAMETER_ERROR, Message: err.Error()})
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
			a.RespondNotFound(r, &userID, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &userID, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &userID, w, apiError)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &userID, w)
		}
		return
	}

	// Create response
	response := &UpdateUserResponse{
		User: result,
	}

	// Write user to response
	a.RespondOk(r, &userID, w, response)
}

func (a *AuthHandler) handleGetUserId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := a.core.Authenticator.RetrieveUserID(*r)
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
			a.RespondNotFound(r, &userID, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &userID, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &userID, w, apiError)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &userID, w)
		}
		return
	}

	response := GetUserByIdResponse{
		User: result,
	}

	// Write user to response
	a.RespondOk(r, &userID, w, response)
}

func (a *AuthHandler) handleDeleteUserId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := a.core.Authenticator.RetrieveUserID(*r)
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
			a.RespondNotFound(r, &userID, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &userID, w, apiError)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &userID, w)
		}
		return
	}

	a.RespondNoContent(r, &userID, w)
}

func (a *AuthHandler) handleUserIdGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := a.core.Authenticator.RetrieveUserID(*r)
	// Retrieve users using path
	id := ps.ByName(USER_ID)

	result, err := a.core.AuthApi.GetGroupsByUserId(a.core.Authenticator.RetrieveUserID(*r), id)

	if err != nil {
		a.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		switch apiError.Code {
		case api.USER_BY_EXTERNAL_ID_NOT_FOUND:
			a.RespondNotFound(r, &userID, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &userID, w, apiError)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &userID, w)
		}
		return
	}

	response := GetGroupsByUserIdResponse{
		Groups: result,
	}

	// Write user to response
	a.RespondOk(r, &userID, w, response)
}

func (a *AuthHandler) handleOrgListUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//TODO: Unimplemented
}
