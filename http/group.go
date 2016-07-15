package http

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tecsisa/authorizr/api"
)

// Requests

type CreateGroupRequest struct {
	Name string `json:"name, omitempty"`
	Path string `json:"path, omitempty"`
}

type UpdateGroupRequest struct {
	Name string `json:"name, omitempty"`
	Path string `json:"path, omitempty"`
}

// Responses

type GetGroupsResponse struct {
	Groups []api.GroupIdentity `json:"groups, omitempty"`
}

type GetGroupMembersResponse struct {
	Members []string `json:"members, omitempty"`
}

type GetGroupPoliciesResponse struct {
	AttachedPolicies []api.PolicyIdentity `json:"policies, omitempty"`
}

func (a *WorkerHandler) HandleCreateGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticatedUser := a.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
	// Decode request
	request := CreateGroupRequest{}
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

	org := ps.ByName(ORG_NAME)
	// Call group API to create a group
	response, err := a.worker.GroupApi.AddGroup(authenticatedUser, org, request.Name, request.Path)

	// Error handling
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(a.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.GROUP_ALREADY_EXIST:
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

	// Write group to response
	a.RespondCreated(r, &authenticatedUser, w, response)
}

func (a *WorkerHandler) HandleDeleteGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticatedUser := a.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
	// Retrieve group org and name from path
	org := ps.ByName(ORG_NAME)
	name := ps.ByName(GROUP_NAME)

	// Call user API to delete group
	err := a.worker.GroupApi.RemoveGroup(authenticatedUser, org, name)

	// Check if there were errors
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(a.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.GROUP_BY_ORG_AND_NAME_NOT_FOUND:
			a.RespondNotFound(r, &authenticatedUser, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &authenticatedUser, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &authenticatedUser, w, apiError)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &authenticatedUser, w)
		}
		return
	} else {
		a.RespondNoContent(r, &authenticatedUser, w)
	}
}

func (a *WorkerHandler) HandleGetGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticatedUser := a.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
	// Retrieve group org and name from path
	org := ps.ByName(ORG_NAME)
	name := ps.ByName(GROUP_NAME)

	// Call group API to retrieve group
	response, err := a.worker.GroupApi.GetGroupByName(authenticatedUser, org, name)

	// Error handling
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(a.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.GROUP_BY_ORG_AND_NAME_NOT_FOUND:
			a.RespondNotFound(r, &authenticatedUser, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &authenticatedUser, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &authenticatedUser, w, apiError)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &authenticatedUser, w)
		}
		return
	}

	// Write group to response
	a.RespondOk(r, &authenticatedUser, w, response)
}

func (a *WorkerHandler) HandleListGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticatedUser := a.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
	// Retrieve group org from path
	org := ps.ByName(ORG_NAME)

	// Retrieve query param if exists
	pathPrefix := r.URL.Query().Get("PathPrefix")

	// Call group API to retrieve groups
	result, err := a.worker.GroupApi.GetGroupList(authenticatedUser, org, pathPrefix)
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(a.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &authenticatedUser, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &authenticatedUser, w, apiError)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &authenticatedUser, w)
		}
		return
	}

	// Create response
	response := &GetGroupsResponse{
		Groups: result,
	}

	// Return groups
	a.RespondOk(r, &authenticatedUser, w, response)

}

func (a *WorkerHandler) HandleUpdateGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticatedUser := a.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
	// Decode request
	request := UpdateGroupRequest{}
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

	// Retrieve group, org from path
	org := ps.ByName(ORG_NAME)
	groupName := ps.ByName(GROUP_NAME)

	// Call group API to update group
	response, err := a.worker.GroupApi.UpdateGroup(authenticatedUser, org, groupName, request.Name, request.Path)

	// Check errors
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(a.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.GROUP_BY_ORG_AND_NAME_NOT_FOUND:
			a.RespondNotFound(r, &authenticatedUser, w, apiError)
		case api.GROUP_ALREADY_EXIST:
			a.RespondConflict(r, &authenticatedUser, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &authenticatedUser, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &authenticatedUser, w, apiError)
		default:
			a.RespondInternalServerError(r, &authenticatedUser, w)
		}
		return
	}

	// Write group to response
	a.RespondOk(r, &authenticatedUser, w, response)
}

func (a *WorkerHandler) HandleListMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticatedUser := a.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
	// Retrieve group, org
	org := ps.ByName(ORG_NAME)
	group := ps.ByName(GROUP_NAME)

	// Call group API to list members
	result, err := a.worker.GroupApi.ListMembers(authenticatedUser, org, group)

	// Check errors
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(a.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.GROUP_BY_ORG_AND_NAME_NOT_FOUND:
			a.RespondNotFound(r, &authenticatedUser, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &authenticatedUser, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &authenticatedUser, w, apiError)
		default:
			a.RespondInternalServerError(r, &authenticatedUser, w)
		}
		return
	}

	// Create response
	response := &GetGroupMembersResponse{
		Members: result,
	}

	// Write GroupMembers to response
	a.RespondOk(r, &authenticatedUser, w, response)

}

func (a *WorkerHandler) HandleAddMember(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticatedUser := a.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
	// Retrieve group, org and user from path
	org := ps.ByName(ORG_NAME)
	user := ps.ByName(USER_ID)
	group := ps.ByName(GROUP_NAME)

	// Call group API to create an group
	err := a.worker.GroupApi.AddMember(authenticatedUser, user, group, org)
	// Error handling
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(a.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.GROUP_BY_ORG_AND_NAME_NOT_FOUND, api.USER_BY_EXTERNAL_ID_NOT_FOUND:
			a.RespondNotFound(r, &authenticatedUser, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &authenticatedUser, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &authenticatedUser, w, apiError)
		case api.USER_IS_ALREADY_A_MEMBER_OF_GROUP:
			a.RespondConflict(r, &authenticatedUser, w, apiError)
		default:
			a.RespondInternalServerError(r, &authenticatedUser, w)
		}
		return
	} else {
		a.RespondNoContent(r, &authenticatedUser, w)
	}
}

func (a *WorkerHandler) HandleRemoveMember(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticatedUser := a.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
	// Retrieve group, org and user from path
	org := ps.ByName(ORG_NAME)
	user := ps.ByName(USER_ID)
	group := ps.ByName(GROUP_NAME)

	// Call group API to create an group
	err := a.worker.GroupApi.RemoveMember(authenticatedUser, user, group, org)
	// Error handling
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(a.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.GROUP_BY_ORG_AND_NAME_NOT_FOUND, api.USER_BY_EXTERNAL_ID_NOT_FOUND, api.USER_IS_NOT_A_MEMBER_OF_GROUP:
			a.RespondNotFound(r, &authenticatedUser, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &authenticatedUser, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &authenticatedUser, w, apiError)
		default:
			a.RespondInternalServerError(r, &authenticatedUser, w)
		}
		return
	} else {
		a.RespondNoContent(r, &authenticatedUser, w)
	}

}

func (a *WorkerHandler) HandleAttachGroupPolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticatedUser := a.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
	// Retrieve group, org and policy from path
	org := ps.ByName(ORG_NAME)
	groupName := ps.ByName(GROUP_NAME)
	policyName := ps.ByName(POLICY_NAME)

	// Call group API to attach policy to group
	err := a.worker.GroupApi.AttachPolicyToGroup(authenticatedUser, org, groupName, policyName)

	// Error handling
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(a.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.GROUP_BY_ORG_AND_NAME_NOT_FOUND, api.POLICY_BY_ORG_AND_NAME_NOT_FOUND:
			a.RespondNotFound(r, &authenticatedUser, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &authenticatedUser, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &authenticatedUser, w, apiError)
		case api.POLICY_IS_ALREADY_ATTACHED_TO_GROUP:
			a.RespondConflict(r, &authenticatedUser, w, apiError)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &authenticatedUser, w)
		}
		return

	} else {
		a.RespondNoContent(r, &authenticatedUser, w)
	}

}

func (a *WorkerHandler) HandleDetachGroupPolicy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticatedUser := a.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
	// Retrieve group, org and policy from path
	org := ps.ByName(ORG_NAME)
	groupName := ps.ByName(GROUP_NAME)
	policyName := ps.ByName(POLICY_NAME)

	// Call group API to detach policy to group
	err := a.worker.GroupApi.DetachPolicyToGroup(authenticatedUser, org, groupName, policyName)

	// Error handling
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(a.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.GROUP_BY_ORG_AND_NAME_NOT_FOUND, api.POLICY_BY_ORG_AND_NAME_NOT_FOUND, api.POLICY_IS_NOT_ATTACHED_TO_GROUP:
			a.RespondNotFound(r, &authenticatedUser, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &authenticatedUser, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &authenticatedUser, w, apiError)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &authenticatedUser, w)
		}
		return

	} else {
		a.RespondNoContent(r, &authenticatedUser, w)
	}
}

func (a *WorkerHandler) HandleListAttachedGroupPolicies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticatedUser := a.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
	// Retrieve group, org from path
	org := ps.ByName(ORG_NAME)
	groupName := ps.ByName(GROUP_NAME)

	// Call group API to retrieve attached policies
	result, err := a.worker.GroupApi.ListAttachedGroupPolicies(authenticatedUser, org, groupName)
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(a.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.GROUP_BY_ORG_AND_NAME_NOT_FOUND:
			a.RespondNotFound(r, &authenticatedUser, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &authenticatedUser, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &authenticatedUser, w, apiError)
		default:
			a.RespondInternalServerError(r, &authenticatedUser, w)
		}
		return
	}

	// Create response
	response := &GetGroupPoliciesResponse{
		AttachedPolicies: result,
	}

	// Return group policies
	a.RespondOk(r, &authenticatedUser, w, response)

}

func (a *WorkerHandler) HandleListAllGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	authenticatedUser := a.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)
	// get PathPrefix from request, so the query can be filtered
	pathPrefix := r.URL.Query().Get("PathPrefix")

	// Call group API to retrieve groups
	result, err := a.worker.GroupApi.GetGroupList(authenticatedUser, "", pathPrefix)
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(a.worker.Logger, requestID, apiError)
		switch apiError.Code {
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &authenticatedUser, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &authenticatedUser, w, apiError)
		default:
			a.RespondInternalServerError(r, &authenticatedUser, w)
		}
		return
	}

	// Create response
	response := &GetGroupsResponse{
		Groups: result,
	}

	// Return groups
	a.RespondOk(r, &authenticatedUser, w, response)

}
