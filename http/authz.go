package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// REQUESTS

type AuthorizeResourcesRequest struct {
	Action    string   `json:"action,omitempty"`
	Resources []string `json:"resources,omitempty"`
}

// RESPONSES

type AuthorizeResourcesResponse struct {
	ResourcesAllowed []string `json:"resourcesAllowed,omitempty"`
}

// HANDLERS

func (wh *WorkerHandler) HandleGetAuthorizedExternalResources(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	requestInfo := wh.getRequestInfo(r)
	// Process request
	request := &AuthorizeResourcesRequest{}
	requestInfo, _, apiErr := wh.processHttpRequest(r, w, nil, request)
	if apiErr != nil {
		wh.processHttpResponse(r, w, requestInfo, nil, apiErr, http.StatusBadRequest)
		return
	}

	// Retrieve allowed resources
	result, err := wh.worker.AuthzApi.GetAuthorizedExternalResources(requestInfo, request.Action, request.Resources)
	response := AuthorizeResourcesResponse{
		ResourcesAllowed: result,
	}
	wh.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}
