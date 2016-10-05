package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// REQUESTS

type AuthorizeResourcesRequest struct {
	Action    string   `json:"action, omitempty"`
	Resources []string `json:"resources, omitempty"`
}

// RESPONSES

type AuthorizeResourcesResponse struct {
	ResourcesAllowed []string `json:"resourcesAllowed, omitempty"`
}

// HANDLERS

func (h *WorkerHandler) HandleGetAuthorizedExternalResources(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	requestInfo := h.GetRequestInfo(r)
	// Process request
	request := &AuthorizeResourcesRequest{}
	requestInfo, _, apiErr := h.processHttpRequest(r, w, nil, request)
	if apiErr != nil {
		h.RespondBadRequest(r, requestInfo, w, apiErr)
		return
	}

	// Retrieve allowed resources
	result, err := h.worker.AuthzApi.GetAuthorizedExternalResources(requestInfo, request.Action, request.Resources)
	response := AuthorizeResourcesResponse{
		ResourcesAllowed: result,
	}
	h.processHttpResponse(r, w, requestInfo, response, err, http.StatusOK)
}
