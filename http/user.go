package http

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tecsisa/authorizr/api"
)

// REQUESTS

type CreateUserRequest struct {
	ExternalID string `json:"externalId, omitempty"`
	Path       string `json:"path, omitempty"`
}

type UpdateUserRequest struct {
	Path string `json:"path, omitempty"`
}

// RESPONSES

type GetUserExternalIDsResponse struct {
	ExternalIDs []string `json:"users, omitempty"`
}

type GetGroupsByUserIdResponse struct {
	Groups []api.GroupIdentity `json:"groups, omitempty"`
}

// HANDLERS

func (h *WorkerHandler) HandleAddUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	authenticatedUser := h.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
	// Decode request
	request := CreateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		apiError := &api.Error{
			Code:    api.INVALID_PARAMETER_ERROR,
			Message: err.Error(),
		}
		api.LogErrorMessage(h.worker.Logger, requestID, apiError)
		h.RespondBadRequest(r, &authenticatedUser, w, apiError)
		return
	}

	// Call user API to create an user
	response, err := h.worker.UserApi.AddUser(authenticatedUser, request.ExternalID, request.Path)

	// Error handling
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.USER_ALREADY_EXIST:
			h.RespondConflict(r, &authenticatedUser, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			h.RespondBadRequest(r, &authenticatedUser, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, &authenticatedUser, w, apiError)
		default: // Unexpected API error
			h.RespondInternalServerError(r, &authenticatedUser, w)
		}
		return
	}

	// Write user to response
	h.RespondCreated(r, &authenticatedUser, w, response)
}

func (h *WorkerHandler) HandleGetUserByExternalID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticatedUser := h.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
	// Retrieve user id from path
	id := ps.ByName(USER_ID)

	// Call user API to retrieve user
	response, err := h.worker.UserApi.GetUserByExternalID(authenticatedUser, id)

	// Error handling
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.USER_BY_EXTERNAL_ID_NOT_FOUND:
			h.RespondNotFound(r, &authenticatedUser, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, &authenticatedUser, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			h.RespondBadRequest(r, &authenticatedUser, w, apiError)
		default: // Unexpected API error
			h.RespondInternalServerError(r, &authenticatedUser, w)
		}
		return
	}

	// Write user to response
	h.RespondOk(r, &authenticatedUser, w, response)
}

func (h *WorkerHandler) HandleListUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Retrieve PathPrefix
	pathPrefix := r.URL.Query().Get("PathPrefix")

	// Call user API
	authenticatedUser := h.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
	result, err := h.worker.UserApi.ListUsers(authenticatedUser, pathPrefix)
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, &authenticatedUser, w, apiError)
		default: // Unexpected API error
			h.RespondInternalServerError(r, &authenticatedUser, w)
		}
		return
	}

	// Create response
	response := &GetUserExternalIDsResponse{
		ExternalIDs: result,
	}

	// Return users
	h.RespondOk(r, &authenticatedUser, w, response)
}

func (h *WorkerHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticatedUser := h.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
	// Decode request
	request := UpdateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		apiError := &api.Error{
			Code:    api.INVALID_PARAMETER_ERROR,
			Message: err.Error(),
		}
		api.LogErrorMessage(h.worker.Logger, requestID, apiError)
		h.RespondBadRequest(r, &authenticatedUser, w, apiError)
		return
	}

	// Retrieve user id from path
	id := ps.ByName(USER_ID)

	// Call user API to update user
	response, err := h.worker.UserApi.UpdateUser(authenticatedUser, id, request.Path)

	// Error handling
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.USER_BY_EXTERNAL_ID_NOT_FOUND:
			h.RespondNotFound(r, &authenticatedUser, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, &authenticatedUser, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			h.RespondBadRequest(r, &authenticatedUser, w, apiError)
		default: // Unexpected API error
			h.RespondInternalServerError(r, &authenticatedUser, w)
		}
		return
	}

	// Write user to response
	h.RespondOk(r, &authenticatedUser, w, response)
}

func (h *WorkerHandler) HandleRemoveUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticatedUser := h.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
	// Retrieve user id from path
	id := ps.ByName(USER_ID)

	// Call user API to delete user
	err := h.worker.UserApi.RemoveUser(authenticatedUser, id)

	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.USER_BY_EXTERNAL_ID_NOT_FOUND:
			h.RespondNotFound(r, &authenticatedUser, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, &authenticatedUser, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			h.RespondBadRequest(r, &authenticatedUser, w, apiError)
		default: // Unexpected API error
			h.RespondInternalServerError(r, &authenticatedUser, w)
		}
		return
	}

	h.RespondNoContent(r, &authenticatedUser, w)
}

func (h *WorkerHandler) HandleListGroupsByUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticatedUser := h.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
	// Retrieve users using path
	id := ps.ByName(USER_ID)

	result, err := h.worker.UserApi.ListGroupsByUser(authenticatedUser, id)

	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.USER_BY_EXTERNAL_ID_NOT_FOUND:
			h.RespondNotFound(r, &authenticatedUser, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, &authenticatedUser, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			h.RespondBadRequest(r, &authenticatedUser, w, apiError)
		default: // Unexpected API error
			h.RespondInternalServerError(r, &authenticatedUser, w)
		}
		return
	}

	response := GetGroupsByUserIdResponse{
		Groups: result,
	}

	// Write user to response
	h.RespondOk(r, &authenticatedUser, w, response)
}
