package http

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tecsisa/authorizr/api"
)

// Requests
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

// Responses

type ListPoliciesResponse struct {
	Policies []api.PolicyIdentity `json:"policies, omitempty"`
}

type ListAllPoliciesResponse struct {
	Policies []api.PolicyIdentity `json:"policies, omitempty"`
}

type GetPolicyGroupsResponse struct {
	Groups []api.GroupIdentity `json:"groups, omitempty"`
}

func (a *WorkerHandler) HandleListPolicies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticatedUser := a.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
	// Retrieve org from path
	org := ps.ByName(ORG_NAME)

	// Retrieve query param if exist
	pathPrefix := r.URL.Query().Get("PathPrefix")

	// Call policy API to retrieve policies
	result, err := a.worker.PolicyApi.GetPolicyList(authenticatedUser, org, pathPrefix)
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(a.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &authenticatedUser, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &authenticatedUser, w, apiError)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &authenticatedUser, w)
		}
		return
	}

	// Create response
	response := &ListPoliciesResponse{
		Policies: result,
	}

	// Return policies
	a.RespondOk(r, &authenticatedUser, w, response)
}

func (a *WorkerHandler) HandleCreatePolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticatedUser := a.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
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
		api.LogErrorMessage(a.worker.Logger, requestID, apiError)
		a.RespondBadRequest(r, &authenticatedUser, w, apiError)
		return
	}

	// Store this policy
	response, err := a.worker.PolicyApi.AddPolicy(authenticatedUser, request.Name, request.Path, org, request.Statements)

	// Error handling
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(a.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.POLICY_ALREADY_EXIST:
			a.RespondConflict(r, &authenticatedUser, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &authenticatedUser, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &authenticatedUser, w, apiError)
		default:
			a.RespondInternalServerError(r, &authenticatedUser, w)
		}
		return
	}

	// Write policy to response
	a.RespondCreated(r, &authenticatedUser, w, response)
}

func (a *WorkerHandler) HandleDeletePolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticatedUser := a.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)

	// Retrieve org and policy name from request path
	orgId := ps.ByName(ORG_NAME)
	policyName := ps.ByName(POLICY_NAME)

	// Call API to delete policy
	err := a.worker.PolicyApi.DeletePolicy(authenticatedUser, orgId, policyName)

	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(a.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.POLICY_BY_ORG_AND_NAME_NOT_FOUND:
			a.RespondNotFound(r, &authenticatedUser, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &authenticatedUser, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &authenticatedUser, w, apiError)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &authenticatedUser, w)
		}
		return
	}

	a.RespondNoContent(r, &authenticatedUser, w)
}

func (a *WorkerHandler) HandleUpdatePolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticatedUser := a.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
	// Decode request
	request := UpdatePolicyRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		apiError := &api.Error{
			Code:    api.INVALID_PARAMETER_ERROR,
			Message: err.Error(),
		}
		api.LogErrorMessage(a.worker.Logger, requestID, apiError)
		a.RespondBadRequest(r, &authenticatedUser, w, apiError)
		return
	}

	// Retrieve policy, org from path
	org := ps.ByName(ORG_NAME)
	policyName := ps.ByName(POLICY_NAME)

	// Call policy API to update policy
	response, err := a.worker.PolicyApi.UpdatePolicy(authenticatedUser, org, policyName, request.Name, request.Path, request.Statements)

	// Check errors
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(a.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.POLICY_BY_ORG_AND_NAME_NOT_FOUND:
			a.RespondNotFound(r, &authenticatedUser, w, apiError)
		case api.POLICY_ALREADY_EXIST:
			a.RespondConflict(r, &authenticatedUser, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &authenticatedUser, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &authenticatedUser, w, apiError)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &authenticatedUser, w)
		}
		return
	}

	// Write policy to response
	a.RespondOk(r, &authenticatedUser, w, response)
}

func (a *WorkerHandler) HandleGetPolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authethicatedUser := a.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
	// Retrieve org and policy name from request path
	orgId := ps.ByName(ORG_NAME)
	policyName := ps.ByName(POLICY_NAME)

	// Call policies API to retrieve policy
	response, err := a.worker.PolicyApi.GetPolicyByName(authethicatedUser, orgId, policyName)

	// Check errors
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(a.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.POLICY_BY_ORG_AND_NAME_NOT_FOUND:
			a.RespondNotFound(r, &authethicatedUser, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &authethicatedUser, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &authethicatedUser, w, apiError)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &authethicatedUser, w)
		}
		return
	}

	// Return policy
	a.RespondOk(r, &authethicatedUser, w, response)
}

func (a *WorkerHandler) HandleGetPolicyAttachedGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	autheticatedUser := a.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
	// Retrieve org and policy name from request path
	orgId := ps.ByName(ORG_NAME)
	policyName := ps.ByName(POLICY_NAME)

	// Call policies API to retrieve attached groups
	result, err := a.worker.PolicyApi.GetAttachedGroups(autheticatedUser, orgId, policyName)
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(a.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.POLICY_BY_ORG_AND_NAME_NOT_FOUND:
			a.RespondNotFound(r, &autheticatedUser, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &autheticatedUser, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &autheticatedUser, w, apiError)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &autheticatedUser, w)
		}
		return
	}

	// Create response
	response := &GetPolicyGroupsResponse{
		Groups: result,
	}

	// Return groups
	a.RespondOk(r, &autheticatedUser, w, response)
}

func (a *WorkerHandler) HandleListAllPolicies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	autheticatedUser := a.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
	// get Org and PathPrefix from request, so the query can be filtered
	pathPrefix := r.URL.Query().Get("PathPrefix")

	// Call policies API to retrieve policies
	result, err := a.worker.PolicyApi.GetPolicyList(autheticatedUser, "", pathPrefix)
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(a.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &autheticatedUser, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &autheticatedUser, w, apiError)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &autheticatedUser, w)
		}
		return
	}

	// Create response
	response := &ListAllPoliciesResponse{
		Policies: result,
	}

	// Return policies
	a.RespondOk(r, &autheticatedUser, w, response)
}
