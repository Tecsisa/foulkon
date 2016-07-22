package auth

import (
	"net/http"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/emanoelxavier/openid2go/openid"
)

const (
	USER_ID_HEADER = "USER-ID"
)

// This struct represents an OIDC connector that implements interface of auth connector
type OIDCAuthConnector struct {
	configuration openid.Configuration
}

func InitOIDCConnector(logger *log.Logger, provider string, clientids []string) (AuthConnector, error) {
	getProviders := func() ([]openid.Provider, error) {
		provider, err := openid.NewProvider(provider, clientids)

		if err != nil {
			return nil, err
		}

		return []openid.Provider{provider}, nil
	}
	errorHandler := func(e error, rw http.ResponseWriter, r *http.Request) bool {
		requestID := r.Header.Get("Request-ID")
		if verr, ok := e.(*openid.ValidationError); ok {
			logger.WithFields(log.Fields{
				"requestID": requestID,
			}).Error(verr.Message)
			http.Error(rw, fmt.Sprintf("Error %v", verr.Message), verr.HTTPStatus)
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
func (c OIDCAuthConnector) Authenticate(h http.Handler) http.Handler {
	userHandler := func(u *openid.User, w http.ResponseWriter, r *http.Request) {
		r.Header.Add(USER_ID_HEADER, u.ID)
		h.ServeHTTP(w, r)
	}
	return openid.AuthenticateUser(&c.configuration, openid.UserHandlerFunc(userHandler))
}

// Retrieve user from OIDC token
func (c OIDCAuthConnector) RetrieveUserID(r http.Request) string {
	userID := r.Header.Get(USER_ID_HEADER)
	r.Header.Del(USER_ID_HEADER)
	return userID
}
