package http

import (
	"encoding/json"
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
}

type ListAllPoliciesResponse struct {
	Policies []api.PolicyIdentity `json:"policies, omitempty"`
}

type ListAttachedGroupsResponse struct {
	Groups []string `json:"groups, omitempty"`
}

// HANDLERS

func (h *WorkerHandler) HandleAddPolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Retrieve Organization
	org := ps.ByName(ORG_NAME)

	// Decode request
	request := CreatePolicyRequest{}
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

	// Store this policy
	response, err := h.worker.PolicyApi.AddPolicy(requestInfo, request.Name, request.Path, org, request.Statements)

	// Error handling
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.POLICY_ALREADY_EXIST:
			h.RespondConflict(r, requestInfo, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			h.RespondBadRequest(r, requestInfo, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, requestInfo, w, apiError)
		default:
			h.RespondInternalServerError(r, requestInfo, w)
		}
		return
	}

	// Write policy to response
	h.RespondCreated(r, requestInfo, w, response)
}

func (h *WorkerHandler) HandleGetPolicyByName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Retrieve org and policy name from request path
	orgId := ps.ByName(ORG_NAME)
	policyName := ps.ByName(POLICY_NAME)

	// Call policies API to retrieve policy
	response, err := h.worker.PolicyApi.GetPolicyByName(requestInfo, orgId, policyName)

	// Check errors
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.POLICY_BY_ORG_AND_NAME_NOT_FOUND:
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

	// Return policy
	h.RespondOk(r, requestInfo, w, response)
}

func (h *WorkerHandler) HandleListPolicies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Retrieve org from path
	org := ps.ByName(ORG_NAME)

	// Retrieve query param if exist
	pathPrefix := r.URL.Query().Get("PathPrefix")

	// Call policy API to retrieve policies
	result, err := h.worker.PolicyApi.ListPolicies(requestInfo, org, pathPrefix)
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.INVALID_PARAMETER_ERROR:
			h.RespondBadRequest(r, requestInfo, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, requestInfo, w, apiError)
		default: // Unexpected API error
			h.RespondInternalServerError(r, requestInfo, w)
		}
		return
	}

	// Create response
	policies := []string{}
	for _, policy := range result {
		policies = append(policies, policy.Name)
	}
	response := &ListPoliciesResponse{
		Policies: policies,
	}

	// Return policies
	h.RespondOk(r, requestInfo, w, response)
}

func (h *WorkerHandler) HandleListAllPolicies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// get Org and PathPrefix from request, so the query can be filtered
	pathPrefix := r.URL.Query().Get("PathPrefix")

	// Call policies API to retrieve policies
	result, err := h.worker.PolicyApi.ListPolicies(requestInfo, "", pathPrefix)
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
	response := &ListAllPoliciesResponse{
		Policies: result,
	}

	// Return policies
	h.RespondOk(r, requestInfo, w, response)
}

func (h *WorkerHandler) HandleUpdatePolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Decode request
	request := UpdatePolicyRequest{}
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

	// Retrieve policy, org from path
	org := ps.ByName(ORG_NAME)
	policyName := ps.ByName(POLICY_NAME)

	// Call policy API to update policy
	response, err := h.worker.PolicyApi.UpdatePolicy(requestInfo, org, policyName, request.Name, request.Path, request.Statements)

	// Check errors
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.POLICY_BY_ORG_AND_NAME_NOT_FOUND:
			h.RespondNotFound(r, requestInfo, w, apiError)
		case api.POLICY_ALREADY_EXIST:
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

	// Write policy to response
	h.RespondOk(r, requestInfo, w, response)
}

func (h *WorkerHandler) HandleRemovePolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Retrieve org and policy name from request path
	orgId := ps.ByName(ORG_NAME)
	policyName := ps.ByName(POLICY_NAME)

	// Call API to delete policy
	err := h.worker.PolicyApi.RemovePolicy(requestInfo, orgId, policyName)

	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.POLICY_BY_ORG_AND_NAME_NOT_FOUND:
			h.RespondNotFound(r, requestInfo, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			h.RespondBadRequest(r, requestInfo, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, requestInfo, w, apiError)
		default: // Unexpected API error
			h.RespondInternalServerError(r, requestInfo, w)
		}
		return
	}

	h.RespondNoContent(r, requestInfo, w)
}

func (h *WorkerHandler) HandleListAttachedGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Retrieve org and policy name from request path
	orgId := ps.ByName(ORG_NAME)
	policyName := ps.ByName(POLICY_NAME)

	// Call policies API to retrieve attached groups
	result, err := h.worker.PolicyApi.ListAttachedGroups(requestInfo, orgId, policyName)
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.POLICY_BY_ORG_AND_NAME_NOT_FOUND:
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

	// Create response
	response := &ListAttachedGroupsResponse{
		Groups: result,
	}

	// Return groups
	h.RespondOk(r, requestInfo, w, response)
}
