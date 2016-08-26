package http

import (
	"encoding/json"
	"net/http"

	"github.com/Tecsisa/foulkon/api"
	"github.com/julienschmidt/httprouter"
)

// REQUESTS

type CreateGroupRequest struct {
	Name string `json:"name, omitempty"`
	Path string `json:"path, omitempty"`
}

type UpdateGroupRequest struct {
	Name string `json:"name, omitempty"`
	Path string `json:"path, omitempty"`
}

// RESPONSES

type ListGroupsResponse struct {
	Groups []string `json:"groups, omitempty"`
}

type ListAllGroupsResponse struct {
	Groups []api.GroupIdentity `json:"groups, omitempty"`
}

type ListMembersResponse struct {
	Members []string `json:"members, omitempty"`
}

type ListAttachedGroupPoliciesResponse struct {
	AttachedPolicies []string `json:"policies, omitempty"`
}

// HANDLERS

func (h *WorkerHandler) HandleAddGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Decode request
	request := CreateGroupRequest{}
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

	org := ps.ByName(ORG_NAME)
	// Call group API to create a group
	response, err := h.worker.GroupApi.AddGroup(requestInfo, org, request.Name, request.Path)

	// Error handling
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.GROUP_ALREADY_EXIST:
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

	// Write group to response
	h.RespondCreated(r, requestInfo, w, response)
}

func (h *WorkerHandler) HandleGetGroupByName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Retrieve group org and name from path
	org := ps.ByName(ORG_NAME)
	name := ps.ByName(GROUP_NAME)

	// Call group API to retrieve group
	response, err := h.worker.GroupApi.GetGroupByName(requestInfo, org, name)

	// Error handling
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.GROUP_BY_ORG_AND_NAME_NOT_FOUND:
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

	// Write group to response
	h.RespondOk(r, requestInfo, w, response)
}

func (h *WorkerHandler) HandleListGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Retrieve group org from path
	org := ps.ByName(ORG_NAME)

	// Retrieve query param if exists
	pathPrefix := r.URL.Query().Get("PathPrefix")

	// Call group API to retrieve groups
	result, err := h.worker.GroupApi.ListGroups(requestInfo, org, pathPrefix)
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

	groups := []string{}
	for _, group := range result {
		groups = append(groups, group.Name)
	}

	// Create response
	response := &ListGroupsResponse{
		Groups: groups,
	}

	// Return groups
	h.RespondOk(r, requestInfo, w, response)
}

func (h *WorkerHandler) HandleListAllGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// get PathPrefix from request, so the query can be filtered
	pathPrefix := r.URL.Query().Get("PathPrefix")

	// Call group API to retrieve groups
	result, err := h.worker.GroupApi.ListGroups(requestInfo, "", pathPrefix)
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, requestInfo, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			h.RespondBadRequest(r, requestInfo, w, apiError)
		default:
			h.RespondInternalServerError(r, requestInfo, w)
		}
		return
	}

	// Create response
	response := &ListAllGroupsResponse{
		Groups: result,
	}

	// Return groups
	h.RespondOk(r, requestInfo, w, response)
}

func (h *WorkerHandler) HandleUpdateGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Decode request
	request := UpdateGroupRequest{}
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

	// Retrieve group, org from path
	org := ps.ByName(ORG_NAME)
	groupName := ps.ByName(GROUP_NAME)

	// Call group API to update group
	response, err := h.worker.GroupApi.UpdateGroup(requestInfo, org, groupName, request.Name, request.Path)

	// Check errors
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.GROUP_BY_ORG_AND_NAME_NOT_FOUND:
			h.RespondNotFound(r, requestInfo, w, apiError)
		case api.GROUP_ALREADY_EXIST:
			h.RespondConflict(r, requestInfo, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, requestInfo, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			h.RespondBadRequest(r, requestInfo, w, apiError)
		default:
			h.RespondInternalServerError(r, requestInfo, w)
		}
		return
	}

	// Write group to response
	h.RespondOk(r, requestInfo, w, response)
}

func (h *WorkerHandler) HandleRemoveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Retrieve group org and name from path
	org := ps.ByName(ORG_NAME)
	name := ps.ByName(GROUP_NAME)

	// Call user API to delete group
	err := h.worker.GroupApi.RemoveGroup(requestInfo, org, name)

	// Check if there were errors
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.GROUP_BY_ORG_AND_NAME_NOT_FOUND:
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

	h.RespondNoContent(r, requestInfo, w)
}

func (h *WorkerHandler) HandleAddMember(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Retrieve group, org and user from path
	org := ps.ByName(ORG_NAME)
	user := ps.ByName(USER_ID)
	group := ps.ByName(GROUP_NAME)

	// Call group API to create an group
	err := h.worker.GroupApi.AddMember(requestInfo, user, group, org)
	// Error handling
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.GROUP_BY_ORG_AND_NAME_NOT_FOUND, api.USER_BY_EXTERNAL_ID_NOT_FOUND:
			h.RespondNotFound(r, requestInfo, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, requestInfo, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			h.RespondBadRequest(r, requestInfo, w, apiError)
		case api.USER_IS_ALREADY_A_MEMBER_OF_GROUP:
			h.RespondConflict(r, requestInfo, w, apiError)
		default:
			h.RespondInternalServerError(r, requestInfo, w)
		}
		return
	}

	h.RespondNoContent(r, requestInfo, w)
}

func (h *WorkerHandler) HandleRemoveMember(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Retrieve group, org and user from path
	org := ps.ByName(ORG_NAME)
	user := ps.ByName(USER_ID)
	group := ps.ByName(GROUP_NAME)

	// Call group API to create an group
	err := h.worker.GroupApi.RemoveMember(requestInfo, user, group, org)
	// Error handling
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.GROUP_BY_ORG_AND_NAME_NOT_FOUND, api.USER_BY_EXTERNAL_ID_NOT_FOUND, api.USER_IS_NOT_A_MEMBER_OF_GROUP:
			h.RespondNotFound(r, requestInfo, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, requestInfo, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			h.RespondBadRequest(r, requestInfo, w, apiError)
		default:
			h.RespondInternalServerError(r, requestInfo, w)
		}
		return
	}

	h.RespondNoContent(r, requestInfo, w)
}

func (h *WorkerHandler) HandleListMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Retrieve group, org
	org := ps.ByName(ORG_NAME)
	group := ps.ByName(GROUP_NAME)

	// Call group API to list members
	result, err := h.worker.GroupApi.ListMembers(requestInfo, org, group)

	// Check errors
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.GROUP_BY_ORG_AND_NAME_NOT_FOUND:
			h.RespondNotFound(r, requestInfo, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, requestInfo, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			h.RespondBadRequest(r, requestInfo, w, apiError)
		default:
			h.RespondInternalServerError(r, requestInfo, w)
		}
		return
	}

	// Create response
	response := &ListMembersResponse{
		Members: result,
	}

	// Write GroupMembers to response
	h.RespondOk(r, requestInfo, w, response)
}

func (h *WorkerHandler) HandleAttachPolicyToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Retrieve group, org and policy from path
	org := ps.ByName(ORG_NAME)
	groupName := ps.ByName(GROUP_NAME)
	policyName := ps.ByName(POLICY_NAME)

	// Call group API to attach policy to group
	err := h.worker.GroupApi.AttachPolicyToGroup(requestInfo, org, groupName, policyName)

	// Error handling
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.GROUP_BY_ORG_AND_NAME_NOT_FOUND, api.POLICY_BY_ORG_AND_NAME_NOT_FOUND:
			h.RespondNotFound(r, requestInfo, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, requestInfo, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			h.RespondBadRequest(r, requestInfo, w, apiError)
		case api.POLICY_IS_ALREADY_ATTACHED_TO_GROUP:
			h.RespondConflict(r, requestInfo, w, apiError)
		default: // Unexpected API error
			h.RespondInternalServerError(r, requestInfo, w)
		}
		return

	}

	h.RespondNoContent(r, requestInfo, w)
}

func (h *WorkerHandler) HandleDetachPolicyToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Retrieve group, org and policy from path
	org := ps.ByName(ORG_NAME)
	groupName := ps.ByName(GROUP_NAME)
	policyName := ps.ByName(POLICY_NAME)

	// Call group API to detach policy to group
	err := h.worker.GroupApi.DetachPolicyToGroup(requestInfo, org, groupName, policyName)

	// Error handling
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.GROUP_BY_ORG_AND_NAME_NOT_FOUND, api.POLICY_BY_ORG_AND_NAME_NOT_FOUND, api.POLICY_IS_NOT_ATTACHED_TO_GROUP:
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

	h.RespondNoContent(r, requestInfo, w)
}

func (h *WorkerHandler) HandleListAttachedGroupPolicies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Retrieve group, org from path
	org := ps.ByName(ORG_NAME)
	groupName := ps.ByName(GROUP_NAME)

	// Call group API to retrieve attached policies
	result, err := h.worker.GroupApi.ListAttachedGroupPolicies(requestInfo, org, groupName)
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.GROUP_BY_ORG_AND_NAME_NOT_FOUND:
			h.RespondNotFound(r, requestInfo, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, requestInfo, w, apiError)
		case api.INVALID_PARAMETER_ERROR:
			h.RespondBadRequest(r, requestInfo, w, apiError)
		default:
			h.RespondInternalServerError(r, requestInfo, w)
		}
		return
	}

	// Create response
	response := &ListAttachedGroupPoliciesResponse{
		AttachedPolicies: result,
	}

	// Return group policies
	h.RespondOk(r, requestInfo, w, response)
}
