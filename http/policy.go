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
	Policies []api.PolicyIdentity
}

type ListAllPoliciesResponse struct {
	Policies []api.PolicyIdentity
}

type GetPolicyGroupsResponse struct {
	Groups []api.GroupIdentity
}

func (a *AuthHandler) handleListPolicies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := a.core.Authenticator.RetrieveUserID(*r)
	// Retrieve org from path
	org := ps.ByName(ORG_NAME)

	// Retrieve query param if exist
	pathPrefix := r.URL.Query().Get("PathPrefix")

	// Call policy API to retrieve policies
	result, err := a.core.AuthApi.GetListPolicies(a.core.Authenticator.RetrieveUserID(*r), org, pathPrefix)
	if err != nil {
		a.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		switch apiError.Code {
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &userID, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &userID, w, apiError)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &userID, w)
		}
		return
	}

	// Create response
	response := &ListPoliciesResponse{
		Policies: result,
	}

	// Return data
	a.RespondOk(r, &userID, w, response)
}

func (a *AuthHandler) handleCreatePolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := a.core.Authenticator.RetrieveUserID(*r)
	// Retrieve Organization
	org := ps.ByName(ORG_NAME)

	// Decode request
	request := CreatePolicyRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		a.core.Logger.Errorln(err)
		a.RespondBadRequest(r, &userID, w, &api.Error{Code: api.INVALID_PARAMETER_ERROR, Message: err.Error()})
		return
	}

	// Store this policy
	storedPolicy, err := a.core.AuthApi.AddPolicy(a.core.Authenticator.RetrieveUserID(*r), request.Name, request.Path, org, &request.Statements)

	// Error handling
	if err != nil {
		a.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		switch apiError.Code {
		case api.POLICY_ALREADY_EXIST:
			a.RespondConflict(r, &userID, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &userID, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &userID, w, apiError)
		default:
			a.RespondInternalServerError(r, &userID, w)
		}
		return
	}

	response := &CreatePolicyResponse{
		Policy: storedPolicy,
	}

	// Write group to response
	a.RespondOk(r, &userID, w, response)
}

func (a *AuthHandler) handleDeletePolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := a.core.Authenticator.RetrieveUserID(*r)
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
			a.RespondNotFound(r, &userID, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &userID, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &userID, w, apiError)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &userID, w)
		}
		return
	}

	a.RespondNoContent(r, &userID, w)
}

func (a *AuthHandler) handleUpdatePolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := a.core.Authenticator.RetrieveUserID(*r)
	// Decode request
	request := UpdatePolicyRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		a.core.Logger.Errorln(err)
		a.RespondBadRequest(r, &userID, w, &api.Error{Code: api.INVALID_PARAMETER_ERROR, Message: err.Error()})
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
			a.RespondNotFound(r, &userID, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &userID, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &userID, w, apiError)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &userID, w)
		}
		return
	}

	// Create response
	response := &UpdatePolicyResponse{
		Policy: result,
	}

	// Write policy to response
	a.RespondOk(r, &userID, w, response)
}

func (a *AuthHandler) handleGetPolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := a.core.Authenticator.RetrieveUserID(*r)
	// Retrieve org and policy name from request path
	orgId := ps.ByName(ORG_NAME)
	policyName := ps.ByName(POLICY_NAME)

	// Call policies API to retrieve policy
	result, err := a.core.AuthApi.GetPolicyByName(a.core.Authenticator.RetrieveUserID(*r), orgId, policyName)

	// Check errors
	if err != nil {
		a.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		switch apiError.Code {
		case api.POLICY_BY_ORG_AND_NAME_NOT_FOUND:
			a.RespondNotFound(r, &userID, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &userID, w, apiError)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &userID, w)
		}
		return
	}

	// Create response
	response := &GetPolicyResponse{
		Policy: result,
	}

	// Return data
	a.RespondOk(r, &userID, w, response)
}

func (a *AuthHandler) handleGetPolicyAttachedGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := a.core.Authenticator.RetrieveUserID(*r)
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
			a.RespondNotFound(r, &userID, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &userID, w, apiError)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &userID, w)
		}
		return
	}

	// Create response
	response := &GetPolicyGroupsResponse{
		Groups: result,
	}

	// Return data
	a.RespondOk(r, &userID, w, response)
}

func (a *AuthHandler) handleListAllPolicies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := a.core.Authenticator.RetrieveUserID(*r)
	// get Org and PathPrefix from request, so the query can be filtered
	org := r.URL.Query().Get("Org")
	pathPrefix := r.URL.Query().Get("PathPrefix")

	// Call policies API to retrieve policies
	result, err := a.core.AuthApi.GetListPolicies(a.core.Authenticator.RetrieveUserID(*r), org, pathPrefix)
	if err != nil {
		a.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		switch apiError.Code {
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &userID, w, apiError)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &userID, w)
		}
		return
	}

	// Create response
	response := &ListAllPoliciesResponse{
		Policies: result,
	}

	// Return data
	a.RespondOk(r, &userID, w, response)
}
