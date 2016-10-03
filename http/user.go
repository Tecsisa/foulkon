package http

import (
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
	// Process request
	request := &CreateUserRequest{}
	requestInfo, _, apiErr := h.processHttpRequest(r, w, nil, request)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}

	// Call user API to create user
	response, err := h.worker.UserApi.AddUser(requestInfo, request.ExternalID, request.Path)
	h.processHttpResponse(r, w, requestInfo, response, err, http.StatusCreated)
}

func (h *WorkerHandler) HandleGetUserByExternalID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}

	// Call user API to get user
	response, err := h.worker.UserApi.GetUserByExternalID(requestInfo, filterData.ExternalID)
	h.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (h *WorkerHandler) HandleListUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}

	// Call user API to list users
	result, total, err := h.worker.UserApi.ListUsers(requestInfo, filterData)
	// Create response
	response := &GetUserExternalIDsResponse{
		ExternalIDs: result,
		Offset:      filterData.Offset,
		Limit:       filterData.Limit,
		Total:       total,
	}
	h.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (h *WorkerHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	request := &UpdateUserRequest{}
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, request)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}

	// Call user API to update user
	response, err := h.worker.UserApi.UpdateUser(requestInfo, filterData.ExternalID, request.Path)
	h.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (h *WorkerHandler) HandleRemoveUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}

	// Call user API to delete user
	err := h.worker.UserApi.RemoveUser(requestInfo, filterData.ExternalID)
	h.processHttpResponse(r, w, requestInfo, nil, err, http.StatusNoContent)
}

func (h *WorkerHandler) HandleListGroupsByUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}

	// Call user API to retrieve user's groups
	result, total, err := h.worker.UserApi.ListGroupsByUser(requestInfo, filterData)
	response := GetGroupsByUserIdResponse{
		Groups: result,
		Offset: filterData.Offset,
		Limit:  filterData.Limit,
		Total:  total,
	}
	h.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}
