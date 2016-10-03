package http

import (
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
	Limit  int      `json:"limit, omitempty"`
	Offset int      `json:"offset, omitempty"`
	Total  int      `json:"total, omitempty"`
}

type ListAllGroupsResponse struct {
	Groups []api.GroupIdentity `json:"groups, omitempty"`
	Limit  int                 `json:"limit, omitempty"`
	Offset int                 `json:"offset, omitempty"`
	Total  int                 `json:"total, omitempty"`
}

type ListMembersResponse struct {
	Members []string `json:"members, omitempty"`
	Limit   int      `json:"limit, omitempty"`
	Offset  int      `json:"offset, omitempty"`
	Total   int      `json:"total, omitempty"`
}

type ListAttachedGroupPoliciesResponse struct {
	AttachedPolicies []string `json:"policies, omitempty"`
	Limit            int      `json:"limit, omitempty"`
	Offset           int      `json:"offset, omitempty"`
	Total            int      `json:"total, omitempty"`
}

// HANDLERS

func (h *WorkerHandler) HandleAddGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	request := &CreateGroupRequest{}
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, request)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}
	// Call group API to create group
	response, err := h.worker.GroupApi.AddGroup(requestInfo, filterData.Org, request.Name, request.Path)
	h.processHttpResponse(r, w, requestInfo, response, err, http.StatusCreated)
}

func (h *WorkerHandler) HandleGetGroupByName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}
	// Call group API to retrieve group
	response, err := h.worker.GroupApi.GetGroupByName(requestInfo, filterData.Org, filterData.GroupName)
	h.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (h *WorkerHandler) HandleListGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}
	// Call group API to retrieve group list
	result, total, err := h.worker.GroupApi.ListGroups(requestInfo, filterData)
	groups := []string{}
	for _, group := range result {
		groups = append(groups, group.Name)
	}
	// Create response
	response := &ListGroupsResponse{
		Groups: groups,
		Offset: filterData.Offset,
		Limit:  filterData.Limit,
		Total:  total,
	}
	h.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (h *WorkerHandler) HandleListAllGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}
	// Call group API to get all groups
	result, total, err := h.worker.GroupApi.ListGroups(requestInfo, filterData)
	// Create response
	response := &ListAllGroupsResponse{
		Groups: result,
		Offset: filterData.Offset,
		Limit:  filterData.Limit,
		Total:  total,
	}
	h.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (h *WorkerHandler) HandleUpdateGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	request := &UpdateGroupRequest{}
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, request)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}
	// Call group API to update group
	response, err := h.worker.GroupApi.UpdateGroup(requestInfo, filterData.Org, filterData.GroupName, request.Name, request.Path)
	h.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (h *WorkerHandler) HandleRemoveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}
	// Call group API to remove group
	err := h.worker.GroupApi.RemoveGroup(requestInfo, filterData.Org, filterData.GroupName)
	h.processHttpResponse(r, w, requestInfo, nil, err, http.StatusNoContent)
}

func (h *WorkerHandler) HandleAddMember(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}
	// Call group API to add member to group
	err := h.worker.GroupApi.AddMember(requestInfo, filterData.ExternalID, filterData.GroupName, filterData.Org)
	h.processHttpResponse(r, w, requestInfo, nil, err, http.StatusNoContent)
}

func (h *WorkerHandler) HandleRemoveMember(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}
	// Call group API to delete member from group
	err := h.worker.GroupApi.RemoveMember(requestInfo, filterData.ExternalID, filterData.GroupName, filterData.Org)
	h.processHttpResponse(r, w, requestInfo, nil, err, http.StatusNoContent)
}

func (h *WorkerHandler) HandleListMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}
	// Call group API to list members of group
	result, total, err := h.worker.GroupApi.ListMembers(requestInfo, filterData)
	response := &ListMembersResponse{
		Members: result,
		Offset:  filterData.Offset,
		Limit:   filterData.Limit,
		Total:   total,
	}
	h.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (h *WorkerHandler) HandleAttachPolicyToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}
	// Call group API to attach policy to group
	err := h.worker.GroupApi.AttachPolicyToGroup(requestInfo, filterData.Org, filterData.GroupName, filterData.PolicyName)
	h.processHttpResponse(r, w, requestInfo, nil, err, http.StatusNoContent)
}

func (h *WorkerHandler) HandleDetachPolicyToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}
	// Call group API to detach policy from group
	err := h.worker.GroupApi.DetachPolicyToGroup(requestInfo, filterData.Org, filterData.GroupName, filterData.PolicyName)
	h.processHttpResponse(r, w, requestInfo, nil, err, http.StatusNoContent)
}

func (h *WorkerHandler) HandleListAttachedGroupPolicies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := h.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}
	// Call group API to list group policies
	result, total, err := h.worker.GroupApi.ListAttachedGroupPolicies(requestInfo, filterData)
	// Create response
	response := &ListAttachedGroupPoliciesResponse{
		AttachedPolicies: result,
		Offset:           filterData.Offset,
		Limit:            filterData.Limit,
		Total:            total,
	}
	h.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}
