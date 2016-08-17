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
	requestInfo := h.GetRequestInfo(r)
	// Decode request
	request := CreateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		apiError := &api.Error{
			Code:    api.INVALID_PARAMETER_ERROR,
			Message: err.Error(),
		}
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		h.RespondBadRequest(r, requestInfo, w, apiError)
		return
	}

	// Call user API to create an user
	response, err := h.worker.UserApi.AddUser(requestInfo, request.ExternalID, request.Path)

	// Error handling
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.USER_ALREADY_EXIST:
			h.RespondConflict(r, requestInfo, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			h.RespondBadRequest(r, requestInfo, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, requestInfo, w, apiError)
		default: // Unexpected API error
			h.RespondInternalServerError(r, requestInfo, w)
		}
		return
	}

	// Write user to response
	h.RespondCreated(r, requestInfo, w, response)
}

func (h *WorkerHandler) HandleGetUserByExternalID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Retrieve user id from path
	id := ps.ByName(USER_ID)

	// Call user API to retrieve user
	response, err := h.worker.UserApi.GetUserByExternalID(requestInfo, id)

	// Error handling
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.USER_BY_EXTERNAL_ID_NOT_FOUND:
			h.RespondNotFound(r, requestInfo, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, requestInfo, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			h.RespondBadRequest(r, requestInfo, w, apiError)
		default: // Unexpected API error
			h.RespondInternalServerError(r, requestInfo, w)
		}
		return
	}

	// Write user to response
	h.RespondOk(r, requestInfo, w, response)
}

func (h *WorkerHandler) HandleListUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Retrieve PathPrefix
	pathPrefix := r.URL.Query().Get("PathPrefix")
	// Call user API
	result, err := h.worker.UserApi.ListUsers(requestInfo, pathPrefix)
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, requestInfo, w, apiError)
		default: // Unexpected API error
			h.RespondInternalServerError(r, requestInfo, w)
		}
		return
	}

	// Create response
	response := &GetUserExternalIDsResponse{
		ExternalIDs: result,
	}

	// Return users
	h.RespondOk(r, requestInfo, w, response)
}

func (h *WorkerHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Decode request
	request := UpdateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		apiError := &api.Error{
			Code:    api.INVALID_PARAMETER_ERROR,
			Message: err.Error(),
		}
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		h.RespondBadRequest(r, requestInfo, w, apiError)
		return
	}

	// Retrieve user id from path
	id := ps.ByName(USER_ID)

	// Call user API to update user
	response, err := h.worker.UserApi.UpdateUser(requestInfo, id, request.Path)

	// Error handling
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.USER_BY_EXTERNAL_ID_NOT_FOUND:
			h.RespondNotFound(r, requestInfo, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, requestInfo, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			h.RespondBadRequest(r, requestInfo, w, apiError)
		default: // Unexpected API error
			h.RespondInternalServerError(r, requestInfo, w)
		}
		return
	}

	// Write user to response
	h.RespondOk(r, requestInfo, w, response)
}

func (h *WorkerHandler) HandleRemoveUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Retrieve user id from path
	id := ps.ByName(USER_ID)

	// Call user API to delete user
	err := h.worker.UserApi.RemoveUser(requestInfo, id)

	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.USER_BY_EXTERNAL_ID_NOT_FOUND:
			h.RespondNotFound(r, requestInfo, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, requestInfo, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			h.RespondBadRequest(r, requestInfo, w, apiError)
		default: // Unexpected API error
			h.RespondInternalServerError(r, requestInfo, w)
		}
		return
	}

	h.RespondNoContent(r, requestInfo, w)
}

func (h *WorkerHandler) HandleListGroupsByUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Retrieve users using path
	id := ps.ByName(USER_ID)

	result, err := h.worker.UserApi.ListGroupsByUser(requestInfo, id)

	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.USER_BY_EXTERNAL_ID_NOT_FOUND:
			h.RespondNotFound(r, requestInfo, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, requestInfo, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			h.RespondBadRequest(r, requestInfo, w, apiError)
		default: // Unexpected API error
			h.RespondInternalServerError(r, requestInfo, w)
		}
		return
	}

	response := GetGroupsByUserIdResponse{
		Groups: result,
	}

	// Write user to response
	h.RespondOk(r, requestInfo, w, response)
}
