package http

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tecsisa/authorizr/api"
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

func (a *WorkerHandler) HandleGetAuthorizedExternalResources(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := a.worker.Authenticator.RetrieveUserID(*r)
	requestID := r.Header.Get(REQUEST_ID_HEADER)

	// Decode request
	request := AuthorizeResourcesRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		apiError := &api.Error{
			Code:    api.INVALID_PARAMETER_ERROR,
			Message: err.Error(),
		}
		api.LogErrorMessage(a.worker.Logger, requestID, apiError)
		a.RespondBadRequest(r, &userID, w, apiError)
		return
	}

	// Retrieve allowed resources
	a.worker.Logger.Debugf("Request ID %v. Action %v, Resources %v", requestID, request.Action, request.Resources)
	result, err := a.worker.AuthzApi.GetAuthorizedExternalResources(userID, request.Action, request.Resources)
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogErrorMessage(a.worker.Logger, requestID, apiError)
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

	response := AuthorizeResourcesResponse{
		ResourcesAllowed: result,
	}

	a.RespondOk(r, &userID, w, response)
}
