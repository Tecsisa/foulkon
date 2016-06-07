package auth

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/emanoelxavier/openid2go/openid"
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
		if verr, ok := e.(*openid.ValidationError); ok {
			http.Error(rw, verr.Message, verr.HTTPStatus)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}
		logger.Errorf("There was an error in authentication. %v", e.Error())
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
