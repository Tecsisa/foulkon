package auth

import (
	"net/http"

	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/emanoelxavier/openid2go/openid"
	"github.com/satori/go.uuid"
)

// This struct represent a connector for OIDC that implements interface of auth connector
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
		transactionID := uuid.NewV4().String()
		if verr, ok := e.(*openid.ValidationError); ok {
			logger.WithFields(log.Fields{
				"RequestID": transactionID,
				"Method":    r.Method,
				"URI":       r.RequestURI,
				"Address":   r.RemoteAddr,
			}).Error(verr.Message)
			http.Error(rw, fmt.Sprintf("TransactionID: %v. Error %v", transactionID, verr.Message), verr.HTTPStatus)
		} else {
			logger.WithFields(log.Fields{
				"RequestID": transactionID,
				"Method":    r.Method,
				"URI":       r.RequestURI,
				"Address":   r.RemoteAddr,
			}).Error("Internal server error")
			http.Error(rw, fmt.Sprintf("TransactionID: %v. Unexpected error", transactionID), http.StatusInternalServerError)
		}
		return true
	}
	configuration, _ := openid.NewConfiguration(openid.ProvidersGetter(getProviders), openid.ErrorHandler(errorHandler))
	return &OIDCAuthConnector{
		configuration: *configuration,
	}, nil

}

// This method retrieve data from a request an check if user is correctly authenticated
func (c OIDCAuthConnector) Authenticate(h http.Handler) http.Handler {
	return openid.Authenticate(&c.configuration, h)
}

// Retrieve user from OIDC token
func (c OIDCAuthConnector) RetrieveUserID(r http.Request) string {
	t, _ := jwt.ParseFromRequest(&r, nil)
	if sub := t.Claims["sub"]; sub != nil {
		return sub.(string)
	}
	return ""
}
