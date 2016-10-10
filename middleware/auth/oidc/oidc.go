package oidc

import (
	"net/http"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/Tecsisa/foulkon/middleware"
	"github.com/Tecsisa/foulkon/middleware/auth"
	"github.com/emanoelxavier/openid2go/openid"
)

// OIDCAuthConnector represents an OIDC connector that implements interface of auth connector
type OIDCAuthConnector struct {
	configuration openid.Configuration
}

func InitOIDCConnector(logger *log.Logger, provider string, clientids []string) (auth.AuthConnector, error) {
	getProviders := func() ([]openid.Provider, error) {
		provider, err := openid.NewProvider(provider, clientids)

		if err != nil {
			return nil, err
		}

		return []openid.Provider{provider}, nil
	}
	errorHandler := func(e error, rw http.ResponseWriter, r *http.Request) bool {
		requestID := r.Header.Get(middleware.REQUEST_ID_HEADER)
		if validationErr, ok := e.(*openid.ValidationError); ok {
			logger.WithFields(log.Fields{
				"requestID": requestID,
			}).Error(validationErr.Message)
			http.Error(rw, fmt.Sprintf("Error %v", validationErr.Message), validationErr.HTTPStatus)
		} else {
			logger.WithFields(log.Fields{
				"requestID": requestID,
			}).Error("Internal server error")
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
