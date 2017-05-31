package oidc

import (
	"net/http"

	"fmt"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/middleware"
	"github.com/Tecsisa/foulkon/middleware/auth"
	"github.com/emanoelxavier/openid2go/openid"
)

// OIDCAuthConnector represents an OIDC connector that implements interface of auth connector
type OIDCAuthConnector struct {
	configuration openid.Configuration
}

// InitOIDCConnector initializes OIDC connector configuration
func InitOIDCConnector(oidcProviders []api.OidcProvider) (auth.AuthConnector, error) {
	getProviders := func() ([]openid.Provider, error) {
		providers := []openid.Provider{}
		for _, oc := range oidcProviders {
			clientIds := []string{}
			for _, clientId := range oc.OidcClients {
				clientIds = append(clientIds, clientId.Name)
			}
			provider, err := openid.NewProvider(oc.IssuerURL, clientIds)
			if err != nil {
				return nil, err
			}
			providers = append(providers, provider)
		}

		return providers, nil
	}
	errorHandler := func(e error, rw http.ResponseWriter, r *http.Request) bool {
		requestID := r.Header.Get(middleware.REQUEST_ID_HEADER)
		if validationErr, ok := e.(*openid.ValidationError); ok {
			apiError := &api.Error{
				Code:    api.AUTHENTICATION_API_ERROR,
				Message: validationErr.Message,
			}
			api.LogOperationError(requestID, "", apiError)
			http.Error(rw, fmt.Sprintf("Error %v", validationErr.Message), validationErr.HTTPStatus)
		} else {
			apiError := &api.Error{
				Code:    api.AUTHENTICATION_API_ERROR,
				Message: validationErr.Message,
			}
			api.LogOperationError(requestID, "", apiError)
			http.Error(rw, "Unexpected error", http.StatusInternalServerError)
		}

		return true
	}
	configuration, _ := openid.NewConfiguration(openid.ProvidersGetter(getProviders), openid.ErrorHandler(errorHandler))
	return &OIDCAuthConnector{
		configuration: *configuration,
	}, nil

}

// This method retrieves data from request an checks if user is correctly authenticated
func (c OIDCAuthConnector) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userHandler := func(u *openid.User, w http.ResponseWriter, r *http.Request) {
			r.Header.Add(middleware.USER_ID_HEADER, u.ID)
			next.ServeHTTP(w, r)
		}
		authenticationHandler := openid.AuthenticateUser(&c.configuration, openid.UserHandlerFunc(userHandler))
		authenticationHandler.ServeHTTP(w, r)
	})

}

// Retrieve user from OIDC token
func (c OIDCAuthConnector) RetrieveUserID(r http.Request) string {
	userID := r.Header.Get(middleware.USER_ID_HEADER)
	return userID
}
