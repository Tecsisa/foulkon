package http

import (
	"net/http"

	"encoding/json"

	"github.com/julienschmidt/httprouter"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/authorizr"
)

type PolicyHandler struct {
	core *authorizr.Core
}

// Requests
type CreatePolicyRequest struct {
	Name       string          `json:"Name, omitempty"`
	Path       string          `json:"Path, omitempty"`
	Statements []api.Statement `json:"Statements, omitempty"`
}

type UpdatePolicyRequest struct {
	Name       string          `json:"Name, omitempty"`
	Path       string          `json:"Path, omitempty"`
	Statements []api.Statement `json:"Statements, omitempty"`
}

// Responses
type CreatePolicyResponse struct {
	Policy *api.Policy
}

type UpdatePolicyResponse struct {
	Policy *api.Policy
}

type ListPoliciesResponse struct {
	Policies []api.Policy
}

type ListAllPoliciesResponse struct {
	Policies []api.Policy
}

func (p *PolicyHandler) handleListPolicies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Retrieve org from path
	org := ps.ByName(ORG_NAME)

	// Retrieve query param if exist
	pathPrefix := r.URL.Query().Get("PathPrefix")

	// Call group API to retrieve groups
	result, err := p.core.PolicyApi.GetPolicies(org, pathPrefix)
	if err != nil {
		p.core.Logger.Errorln(err)
		RespondInternalServerError(w)
		return
	}

	// Create response
	response := &ListPoliciesResponse{
		Policies: result,
	}

	// Return data
	RespondOk(w, response)
}

func (p *PolicyHandler) handleCreatePolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Retrieve Organization
	org := ps.ByName(ORG_NAME)

	// Decode request
	request := CreatePolicyRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		p.core.Logger.Errorln(err)
		RespondBadRequest(w)
		return
	}

	// Store this policy
	storedPolicy, err := p.core.PolicyApi.AddPolicy(request.Name, request.Path, org, &request.Statements)

	// Error handling
	if err != nil {
		p.core.Logger.Errorln(err)
		switch err.(*api.Error).Code {
		case api.POLICY_ALREADY_EXIST:
			RespondConflict(w)
			return
		case api.INVALID_PARAMETER_ERROR:
			RespondBadRequest(w)
			return
		default:
			RespondInternalServerError(w)
			return
		}
		RespondInternalServerError(w)
		return
	}

	response := &CreatePolicyResponse{
		Policy: storedPolicy,
	}

	// Write group to response
	RespondOk(w, response)
}

func (p *PolicyHandler) handleDeletePolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Retrieve org and policy name from request path
	orgId := ps.ByName(ORG_NAME)
	policyName := ps.ByName(POLICY_NAME)

	// Call API to delete policy
	err := p.core.PolicyApi.DeletePolicy(orgId, policyName)

	if err != nil {
		p.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		switch apiError.Code {
		case api.POLICY_BY_ORG_AND_NAME_NOT_FOUND:
			RespondNotFound(w)
		default: // Unexpected API error
			RespondInternalServerError(w)
		}
		return
	}

	RespondNoContent(w)
}

func (p *PolicyHandler) handleUpdatePolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Decode request
	request := UpdatePolicyRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		p.core.Logger.Errorln(err)
		RespondBadRequest(w)
		return
	}

	// Retrieve policy, org from path
	org := ps.ByName(ORG_NAME)
	policyName := ps.ByName(POLICY_NAME)

	// Check errors
	if err != nil {
		p.core.Logger.Errorln(err)
		RespondBadRequest(w)
		return
	}

	// Call policy API to update policy
	result, err := p.core.PolicyApi.UpdatePolicy(org, policyName, request.Name, request.Path, request.Statements)

	// Check errors
	if err != nil {
		p.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		if apiError.Code == api.POLICY_BY_ORG_AND_NAME_NOT_FOUND {
			RespondNotFound(w)
			return
		} else { // Unexpected API error
			RespondInternalServerError(w)
			return
		}
	}

	// Create response
	response := &UpdatePolicyResponse{
		Policy: result,
	}

	// Write policy to response
	RespondOk(w, response)
}

func (p *PolicyHandler) handleGetPolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (p *PolicyHandler) handleGetPolicyAttachedGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (p *PolicyHandler) handleListAllPolicies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// get Org and PathPrefix from request, so the query can be filtered
	org := r.URL.Query().Get("Org")
	pathPrefix := r.URL.Query().Get("PathPrefix")

	// Call policies API to retrieve policies
	result, err := p.core.PolicyApi.GetPolicies(org, pathPrefix)
	if err != nil {
		p.core.Logger.Errorln(err)
		RespondInternalServerError(w)
		return
	}

	// Create response
	response := &ListAllPoliciesResponse{
		Policies: result,
	}

	// Return data
	RespondOk(w, response)
}
