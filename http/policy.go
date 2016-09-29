package http

import (
	"net/http"

	"github.com/Tecsisa/foulkon/api"
	"github.com/julienschmidt/httprouter"
)

// REQUESTS

type CreatePolicyRequest struct {
	Name       string          `json:"name, omitempty"`
	Path       string          `json:"path, omitempty"`
	Statements []api.Statement `json:"statements, omitempty"`
}

type UpdatePolicyRequest struct {
	Name       string          `json:"name, omitempty"`
	Path       string          `json:"path, omitempty"`
	Statements []api.Statement `json:"statements, omitempty"`
}

// RESPONSES

type ListPoliciesResponse struct {
	Policies []string `json:"policies, omitempty"`
	Limit    int      `json:"limit, omitempty"`
	Offset   int      `json:"offset, omitempty"`
	Total    int      `json:"total, omitempty"`
}

type ListAllPoliciesResponse struct {
	Policies []api.PolicyIdentity `json:"policies, omitempty"`
	Limit    int                  `json:"limit, omitempty"`
	Offset   int                  `json:"offset, omitempty"`
	Total    int                  `json:"total, omitempty"`
}

type ListAttachedGroupsResponse struct {
	Groups []string `json:"groups, omitempty"`
	Limit  int      `json:"limit, omitempty"`
	Offset int      `json:"offset, omitempty"`
	Total  int      `json:"total, omitempty"`
}

// HANDLERS

func (h *WorkerHandler) HandleAddPolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	request := &CreatePolicyRequest{}
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, request)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}

	// Call policy API to create policy
	response, err := h.worker.PolicyApi.AddPolicy(requestInfo, request.Name, request.Path, filterData.Org, request.Statements)
	h.processHttpResponse(r, w, requestInfo, response, err, http.StatusCreated)
}

func (h *WorkerHandler) HandleGetPolicyByName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}

	// Call policy API to retrieve policy
	response, err := h.worker.PolicyApi.GetPolicyByName(requestInfo, filterData.Org, filterData.PolicyName)
	h.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (h *WorkerHandler) HandleListPolicies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}
	// Call policy API to list policies
	result, total, err := h.worker.PolicyApi.ListPolicies(requestInfo, filterData)
	// Create response
	policies := []string{}
	for _, policy := range result {
		policies = append(policies, policy.Name)
	}
	response := &ListPoliciesResponse{
		Policies: policies,
		Offset:   filterData.Offset,
		Limit:    filterData.Limit,
		Total:    total,
	}
	h.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (h *WorkerHandler) HandleListAllPolicies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}
	// Call policy API to list all policies
	result, total, err := h.worker.PolicyApi.ListPolicies(requestInfo, filterData)
	// Create response
	response := &ListAllPoliciesResponse{
		Policies: result,
		Offset:   filterData.Offset,
		Limit:    filterData.Limit,
		Total:    total,
	}
	h.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (h *WorkerHandler) HandleUpdatePolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	request := &UpdatePolicyRequest{}
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, request)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}
	// Call policy API to update policy
	response, err := h.worker.PolicyApi.UpdatePolicy(requestInfo, filterData.Org, filterData.PolicyName, request.Name, request.Path, request.Statements)
	h.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (h *WorkerHandler) HandleRemovePolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}
	// Call policy API to remove policy
	err := h.worker.PolicyApi.RemovePolicy(requestInfo, filterData.Org, filterData.PolicyName)
	h.processHttpResponse(r, w, requestInfo, nil, err, http.StatusNoContent)
}

func (h *WorkerHandler) HandleListAttachedGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}
	// Call policy API to list attached groups
	result, total, err := h.worker.PolicyApi.ListAttachedGroups(requestInfo, filterData)
	// Create response
	response := &ListAttachedGroupsResponse{
		Groups: result,
		Offset: filterData.Offset,
		Limit:  filterData.Limit,
		Total:  total,
	}
	h.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}
