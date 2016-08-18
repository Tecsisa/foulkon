package http

import (
	"encoding/json"
	"net/http"

	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/tecsisa/foulkon/api"
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
	// Decode request
	request := AuthorizeResourcesRequest{}
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

	// Retrieve allowed resources
	result, err := h.worker.AuthzApi.GetAuthorizedExternalResources(requestInfo, request.Action, request.Resources)
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(h.worker.Logger, requestInfo, apiError)
		switch apiError.Code {
		case api.INVALID_PARAMETER_ERROR:
			h.RespondBadRequest(r, requestInfo, w, apiError)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			h.RespondForbidden(r, requestInfo, w, apiError)
		default: // Unexpected API error
			h.RespondInternalServerError(r, requestInfo, w)
		}
		return
	}

	if result == nil || len(result) < 1 {
		h.RespondForbidden(r, requestInfo, w, &api.Error{
			Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to any resource", requestInfo.Identifier),
		})
		return
	}

	response := AuthorizeResourcesResponse{
		ResourcesAllowed: result,
	}

	h.RespondOk(r, requestInfo, w, response)
}
