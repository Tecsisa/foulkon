package http

import (
	"encoding/json"
	"net/http"

	"github.com/Tecsisa/foulkon/api"
	"github.com/julienschmidt/httprouter"
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
	Limit       int      `json:"limit, omitempty"`
	Offset      int      `json:"offset, omitempty"`
	Total       int      `json:"total, omitempty"`
}

type GetGroupsByUserIdResponse struct {
	Groups []api.GroupIdentity `json:"groups, omitempty"`
	Limit  int                 `json:"limit, omitempty"`
	Offset int                 `json:"offset, omitempty"`
	Total  int                 `json:"total, omitempty"`
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

func (h *WorkerHandler) HandleListUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Retrieve filterData
	filterData, err := getFilterData(r, ps)
	if err != nil {
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		h.RespondBadRequest(r, requestInfo, w, apiError)
		return
	}

	result, total, err := h.worker.UserApi.ListUsers(requestInfo, filterData)
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, requestInfo, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			h.RespondBadRequest(r, requestInfo, w, apiError)
		default: // Unexpected API error
			h.RespondInternalServerError(r, requestInfo, w)
		}
		return
	}

	// Create response
	response := &GetUserExternalIDsResponse{
		ExternalIDs: result,
		Offset:      filterData.Offset,
		Limit:       filterData.Limit,
		Total:       total,
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

	// Retrieve filterData
	filterData, err := getFilterData(r, ps)
	if err != nil {
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		h.RespondBadRequest(r, requestInfo, w, apiError)
		return
	}
	result, total, err := h.worker.UserApi.ListGroupsByUser(requestInfo, filterData)

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
		Offset: filterData.Offset,
		Limit:  filterData.Limit,
		Total:  total,
	}

	// Write user to response
	h.RespondOk(r, requestInfo, w, response)
}
