package http

import (
	"net/http"

	"github.com/Tecsisa/foulkon/api"
	"github.com/julienschmidt/httprouter"
)

// REQUESTS

type CreatePolicyRequest struct {
	Name       string          `json:"name,omitempty"`
	Path       string          `json:"path,omitempty"`
	Statements []api.Statement `json:"statements,omitempty"`
}

type UpdatePolicyRequest struct {
	Name       string          `json:"name,omitempty"`
	Path       string          `json:"path,omitempty"`
	Statements []api.Statement `json:"statements,omitempty"`
}

// RESPONSES

type ListPoliciesResponse struct {
	Policies []string `json:"policies,omitempty"`
	Limit    int      `json:"limit"`
	Offset   int      `json:"offset"`
	Total    int      `json:"total"`
}

type ListAllPoliciesResponse struct {
	Policies []api.PolicyIdentity `json:"policies,omitempty"`
	Limit    int                  `json:"limit"`
	Offset   int                  `json:"offset"`
	Total    int                  `json:"total"`
}

type ListAttachedGroupsResponse struct {
	Groups []api.PolicyGroups `json:"groups,omitempty"`
	Limit  int                `json:"limit"`
	Offset int                `json:"offset"`
	Total  int                `json:"total"`
}

// HANDLERS

func (wh *WorkerHandler) HandleAddPolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	request := &CreatePolicyRequest{}
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, request)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}

	// Call policy API to create policy
	response, err := wh.worker.PolicyApi.AddPolicy(requestInfo, request.Name, request.Path, filterData.Org, request.Statements)
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusCreated)
}

func (wh *WorkerHandler) HandleGetPolicyByName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}

	// Call policy API to retrieve policy
	response, err := wh.worker.PolicyApi.GetPolicyByName(requestInfo, filterData.Org, filterData.PolicyName)
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (wh *WorkerHandler) HandleListPolicies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}
	// Call policy API to list policies
	result, total, err := wh.worker.PolicyApi.ListPolicies(requestInfo, filterData)
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
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (wh *WorkerHandler) HandleListAllPolicies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}
	// Call policy API to list all policies
	result, total, err := wh.worker.PolicyApi.ListPolicies(requestInfo, filterData)
	// Create response
	response := &ListAllPoliciesResponse{
		Policies: result,
		Offset:   filterData.Offset,
		Limit:    filterData.Limit,
		Total:    total,
	}
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (wh *WorkerHandler) HandleUpdatePolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	request := &UpdatePolicyRequest{}
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, request)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}
	// Call policy API to update policy
	response, err := wh.worker.PolicyApi.UpdatePolicy(requestInfo, filterData.Org, filterData.PolicyName, request.Name, request.Path, request.Statements)
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (wh *WorkerHandler) HandleRemovePolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}
	// Call policy API to remove policy
	err := wh.worker.PolicyApi.RemovePolicy(requestInfo, filterData.Org, filterData.PolicyName)
	wh.processHttpResponse(r, w, requestInfo, nil, err, http.StatusNoContent)
}

func (wh *WorkerHandler) HandleListAttachedGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}
	// Call policy API to list attached groups
	result, total, err := wh.worker.PolicyApi.ListAttachedGroups(requestInfo, filterData)
	// Create response
	response := &ListAttachedGroupsResponse{
		Groups: result,
		Offset: filterData.Offset,
		Limit:  filterData.Limit,
		Total:  total,
	}
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}
