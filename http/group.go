package http

import (
	"net/http"

	"github.com/Tecsisa/foulkon/api"
	"github.com/julienschmidt/httprouter"
)

// REQUESTS

type CreateGroupRequest struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
}

type UpdateGroupRequest struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
}

// RESPONSES

type ListGroupsResponse struct {
	Groups []string `json:"groups,omitempty"`
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
	Total  int      `json:"total"`
}

type ListAllGroupsResponse struct {
	Groups []api.GroupIdentity `json:"groups,omitempty"`
	Limit  int                 `json:"limit"`
	Offset int                 `json:"offset"`
	Total  int                 `json:"total"`
}

type ListMembersResponse struct {
	Members []api.GroupMembers `json:"members,omitempty"`
	Limit   int                `json:"limit"`
	Offset  int                `json:"offset"`
	Total   int                `json:"total"`
}

type ListAttachedGroupPoliciesResponse struct {
	AttachedPolicies []api.GroupPolicies `json:"policies,omitempty"`
	Limit            int                 `json:"limit"`
	Offset           int                 `json:"offset"`
	Total            int                 `json:"total"`
}

// HANDLERS

func (wh *WorkerHandler) HandleAddGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	request := &CreateGroupRequest{}
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, request)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}
	// Call group API to create group
	response, err := wh.worker.GroupApi.AddGroup(requestInfo, filterData.Org, request.Name, request.Path)
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusCreated)
}

func (wh *WorkerHandler) HandleGetGroupByName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}
	// Call group API to retrieve group
	response, err := wh.worker.GroupApi.GetGroupByName(requestInfo, filterData.Org, filterData.GroupName)
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (wh *WorkerHandler) HandleListGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}
	// Call group API to retrieve group list
	result, total, err := wh.worker.GroupApi.ListGroups(requestInfo, filterData)
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
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (wh *WorkerHandler) HandleListAllGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}
	// Call group API to get all groups
	result, total, err := wh.worker.GroupApi.ListGroups(requestInfo, filterData)
	// Create response
	response := &ListAllGroupsResponse{
		Groups: result,
		Offset: filterData.Offset,
		Limit:  filterData.Limit,
		Total:  total,
	}
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (wh *WorkerHandler) HandleUpdateGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	request := &UpdateGroupRequest{}
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, request)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}
	// Call group API to update group
	response, err := wh.worker.GroupApi.UpdateGroup(requestInfo, filterData.Org, filterData.GroupName, request.Name, request.Path)
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (wh *WorkerHandler) HandleRemoveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}
	// Call group API to remove group
	err := wh.worker.GroupApi.RemoveGroup(requestInfo, filterData.Org, filterData.GroupName)
	wh.processHttpResponse(r, w, requestInfo, nil, err, http.StatusNoContent)
}

func (wh *WorkerHandler) HandleAddMember(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}
	// Call group API to add member to group
	err := wh.worker.GroupApi.AddMember(requestInfo, filterData.ExternalID, filterData.GroupName, filterData.Org)
	wh.processHttpResponse(r, w, requestInfo, nil, err, http.StatusNoContent)
}

func (wh *WorkerHandler) HandleRemoveMember(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}
	// Call group API to delete member from group
	err := wh.worker.GroupApi.RemoveMember(requestInfo, filterData.ExternalID, filterData.GroupName, filterData.Org)
	wh.processHttpResponse(r, w, requestInfo, nil, err, http.StatusNoContent)
}

func (wh *WorkerHandler) HandleListMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}
	// Call group API to list members of group
	result, total, err := wh.worker.GroupApi.ListMembers(requestInfo, filterData)
	response := &ListMembersResponse{
		Members: result,
		Offset:  filterData.Offset,
		Limit:   filterData.Limit,
		Total:   total,
	}
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}

func (wh *WorkerHandler) HandleAttachPolicyToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}
	// Call group API to attach policy to group
	err := wh.worker.GroupApi.AttachPolicyToGroup(requestInfo, filterData.Org, filterData.GroupName, filterData.PolicyName)
	wh.processHttpResponse(r, w, requestInfo, nil, err, http.StatusNoContent)
}

func (wh *WorkerHandler) HandleDetachPolicyToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}
	// Call group API to detach policy from group
	err := wh.worker.GroupApi.DetachPolicyToGroup(requestInfo, filterData.Org, filterData.GroupName, filterData.PolicyName)
	wh.processHttpResponse(r, w, requestInfo, nil, err, http.StatusNoContent)
}

func (wh *WorkerHandler) HandleListAttachedGroupPolicies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Process request
	requestInfo, filterData, apiErr := wh.processHttpRequest(r, w, ps, nil)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}
	// Call group API to list group policies
	result, total, err := wh.worker.GroupApi.ListAttachedGroupPolicies(requestInfo, filterData)
	// Create response
	response := &ListAttachedGroupPoliciesResponse{
		AttachedPolicies: result,
		Offset:           filterData.Offset,
		Limit:            filterData.Limit,
		Total:            total,
	}
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}
