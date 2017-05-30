package http

import (
	"net/http"

	"github.com/Tecsisa/foulkon/api"
	"github.com/julienschmidt/httprouter"
)

// REQUESTS

type CreateUserRequest struct {
	ExternalID string `json:"externalId,omitempty"`
	Path       string `json:"path,omitempty"`
}

type UpdateUserRequest struct {
	Path string `json:"path,omitempty"`
}

// RESPONSES

type GetUserExternalIDsResponse struct {
	ExternalIDs []string `json:"users,omitempty"`
	Limit       int      `json:"limit"`
	Offset      int      `json:"offset"`
	Total       int      `json:"total"`
}

type GetGroupsByUserIdResponse struct {
	Groups []api.UserGroups `json:"groups,omitempty"`
	Limit  int              `json:"limit"`
	Offset int              `json:"offset"`
	Total  int              `json:"total"`
}

// HANDLERS

func (wh *WorkerHandler) HandleAddUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Process request
	request := &CreateUserRequest{}
	requestInfo, _, apiErr := wh.processHttpRequest(r, w, nil, request)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}

	// Call user API to create user
	response, err := wh.worker.UserApi.AddUser(requestInfo, request.ExternalID, request.Path)
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusCreated)
}

func (wh *WorkerHandler) HandleGetUserByExternalID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}

	// Call user API to get user
	response, err := wh.worker.UserApi.GetUserByExternalID(requestInfo, filterData.ExternalID)
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (wh *WorkerHandler) HandleListUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}

	// Call user API to list users
	result, total, err := wh.worker.UserApi.ListUsers(requestInfo, filterData)
	// Create response
	response := &GetUserExternalIDsResponse{
		ExternalIDs: result,
		Offset:      filterData.Offset,
		Limit:       filterData.Limit,
		Total:       total,
	}
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (wh *WorkerHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	request := &UpdateUserRequest{}
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, request)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}

	// Call user API to update user
	response, err := wh.worker.UserApi.UpdateUser(requestInfo, filterData.ExternalID, request.Path)
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (wh *WorkerHandler) HandleRemoveUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}

	// Call user API to delete user
	err := wh.worker.UserApi.RemoveUser(requestInfo, filterData.ExternalID)
	wh.processHttpResponse(r, w, requestInfo, nil, err, http.StatusNoContent)
}

func (wh *WorkerHandler) HandleListGroupsByUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}

	// Call user API to retrieve user's groups
	result, total, err := wh.worker.UserApi.ListGroupsByUser(requestInfo, filterData)
	response := GetGroupsByUserIdResponse{
		Groups: result,
		Offset: filterData.Offset,
		Limit:  filterData.Limit,
		Total:  total,
	}
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}
