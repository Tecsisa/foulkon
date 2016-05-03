package auth

import "net/http"
import "github.com/emanoelxavier/openid2go/openid"

// This struct represent a connector for OIDC that implements interface of auth connector
type OIDCAuthConnector struct {
	configuration openid.Configuration
}

func InitOIDCConnector() (AuthConnector, error) {
	configuration, _ := openid.NewConfiguration(openid.ProvidersGetter(getProviders))
	return &OIDCAuthConnector{
		configuration: *configuration,
	}, nil

}

// This method retrieve data from a request an check if user is correctly authenticated
func (c OIDCAuthConnector) Authenticate(h http.Handler) http.Handler {
	return openid.Authenticate(&c.configuration, h)
}

func getProviders() ([]openid.Provider, error) {
	provider, err := openid.NewProvider("http://le001.mad.es.tecsisa.com:5556", []string{"9jCU4aaDHjV-y59SSlGwfrmpdo4mIkGBW4E41QvI-X0=@127.0.0.1"})

	if err != nil {
		return nil, err
	}

	return []openid.Provider{provider}, nil
}
