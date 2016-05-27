package http

import (
	"github.com/julienschmidt/httprouter"
	"github.com/tecsisa/authorizr/api"
	"net/http"
)

// Responses

type GetEffectByUserActionResourceResponse struct {
	EffectRestriction api.EffectRestriction `json:", omitempty"`
}

func (a *AuthHandler) handleGetEffectByUserActionResource(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := a.core.Authenticator.RetrieveUserID(*r)

	// Get parameters from URL
	action := r.URL.Query().Get("Action")
	resourceUrn := r.URL.Query().Get("Urn")

	// Retrieve effect for this user, action and resource urn
	result, err := a.core.AuthApi.GetEffectByUserActionResource(userID, action, resourceUrn)
	if err != nil {
		a.core.Logger.Errorln(err)
		// Transform to API errors
		apiError := err.(*api.Error)
		switch apiError.Code {
		case api.INVALID_PARAMETER_ERROR:
			a.RespondBadRequest(r, &userID, w)
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			a.RespondForbidden(r, &userID, w)
		default: // Unexpected API error
			a.RespondInternalServerError(r, &userID, w)
		}
		return
	}

	a.RespondOk(r, &userID, w, result)

}
