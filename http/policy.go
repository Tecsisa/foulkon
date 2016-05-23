package http

import (
	"net/http"

	"encoding/json"

	"github.com/julienschmidt/httprouter"
	"github.com/tecsisa/authorizr/api"
)

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

type GetPolicyResponse struct {
	Policy *api.Policy
}

type ListPoliciesResponse struct {
	Policies []api.Policy
}

type ListAllPoliciesResponse struct {
	Policies []api.Policy
}

type GetPolicyGroupsResponse struct {
	Groups []api.Group
}

func (a *AuthHandler) handleListPolicies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Retrieve org from path
	org := ps.ByName(ORG_NAME)

	// Retrieve query param if exist
	pathPrefix := r.URL.Query().Get("PathPrefix")

	// Call group API to retrieve groups
	result, err := a.core.AuthApi.GetPolicies(a.core.Authenticator.RetrieveUserID(*r), org, pathPrefix)
	if err != nil {
		a.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		switch apiError.Code {
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			RespondForbidden(w)
		default: // Unexpected API error
			RespondInternalServerError(w)
		}
		return
	}

	// Create response
	response := &ListPoliciesResponse{
		Policies: result,
	}

	// Return data
	RespondOk(w, response)
}

func (a *AuthHandler) handleCreatePolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Retrieve Organization
	org := ps.ByName(ORG_NAME)

	// Decode request
	request := CreatePolicyRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		a.core.Logger.Errorln(err)
		RespondBadRequest(w)
		return
	}

	// Store this policy
	storedPolicy, err := a.core.AuthApi.AddPolicy(a.core.Authenticator.RetrieveUserID(*r), request.Name, request.Path, org, &request.Statements)

	// Error handling
	if err != nil {
		a.core.Logger.Errorln(err)
		switch err.(*api.Error).Code {
		case api.POLICY_ALREADY_EXIST:
			RespondConflict(w)
		case api.INVALID_PARAMETER_ERROR:
			RespondBadRequest(w)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			RespondForbidden(w)
		default:
			RespondInternalServerError(w)
		}
		return
	}

	response := &CreatePolicyResponse{
		Policy: storedPolicy,
	}

	// Write group to response
	RespondOk(w, response)
}

func (a *AuthHandler) handleDeletePolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Retrieve org and policy name from request path
	orgId := ps.ByName(ORG_NAME)
	policyName := ps.ByName(POLICY_NAME)

	// Call API to delete policy
	err := a.core.AuthApi.DeletePolicy(a.core.Authenticator.RetrieveUserID(*r), orgId, policyName)

	if err != nil {
		a.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		switch apiError.Code {
		case api.POLICY_BY_ORG_AND_NAME_NOT_FOUND:
			RespondNotFound(w)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			RespondForbidden(w)
		default: // Unexpected API error
			RespondInternalServerError(w)
		}
		return
	}

	RespondNoContent(w)
}

func (a *AuthHandler) handleUpdatePolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Decode request
	request := UpdatePolicyRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		a.core.Logger.Errorln(err)
		RespondBadRequest(w)
		return
	}

	// Retrieve policy, org from path
	org := ps.ByName(ORG_NAME)
	policyName := ps.ByName(POLICY_NAME)

	// Call policy API to update policy
	result, err := a.core.AuthApi.UpdatePolicy(a.core.Authenticator.RetrieveUserID(*r), org, policyName, request.Name, request.Path, request.Statements)

	// Check errors
	if err != nil {
		a.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		switch apiError.Code {
		case api.POLICY_BY_ORG_AND_NAME_NOT_FOUND:
			RespondNotFound(w)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			RespondForbidden(w)
		default: // Unexpected API error
			RespondInternalServerError(w)
		}
		return
	}

	// Create response
	response := &UpdatePolicyResponse{
		Policy: result,
	}

	// Write policy to response
	RespondOk(w, response)
}

func (a *AuthHandler) handleGetPolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Retrieve org and policy name from request path
	orgId := ps.ByName(ORG_NAME)
	policyName := ps.ByName(POLICY_NAME)

	// Call policies API to retrieve policy
	result, err := a.core.AuthApi.GetPolicy(a.core.Authenticator.RetrieveUserID(*r), orgId, policyName)

	// Check errors
	if err != nil {
		a.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		switch apiError.Code {
		case api.POLICY_BY_ORG_AND_NAME_NOT_FOUND:
			RespondNotFound(w)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			RespondForbidden(w)
		default: // Unexpected API error
			RespondInternalServerError(w)
		}
		return
	}

	// Create response
	response := &GetPolicyResponse{
		Policy: result,
	}

	// Return data
	RespondOk(w, response)
}

func (a *AuthHandler) handleGetPolicyAttachedGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Retrieve org and policy name from request path
	orgId := ps.ByName(ORG_NAME)
	policyName := ps.ByName(POLICY_NAME)

	// Call policies API to retrieve attached groups
	result, err := a.core.AuthApi.GetPolicyAttachedGroups(a.core.Authenticator.RetrieveUserID(*r), orgId, policyName)
	if err != nil {
		a.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		switch apiError.Code {
		case api.POLICY_BY_ORG_AND_NAME_NOT_FOUND:
			RespondNotFound(w)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			RespondForbidden(w)
		default: // Unexpected API error
			RespondInternalServerError(w)
		}
		return
	}

	// Create response
	response := &GetPolicyGroupsResponse{
		Groups: result,
	}

	// Return data
	RespondOk(w, response)
}

func (a *AuthHandler) handleListAllPolicies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// get Org and PathPrefix from request, so the query can be filtered
	org := r.URL.Query().Get("Org")
	pathPrefix := r.URL.Query().Get("PathPrefix")

	// Call policies API to retrieve policies
	result, err := a.core.AuthApi.GetPolicies(a.core.Authenticator.RetrieveUserID(*r), org, pathPrefix)
	if err != nil {
		a.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		switch apiError.Code {
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			RespondForbidden(w)
		default: // Unexpected API error
			RespondInternalServerError(w)
		}
		return
	}

	// Create response
	response := &ListAllPoliciesResponse{
		Policies: result,
	}

	// Return data
	RespondOk(w, response)
}
