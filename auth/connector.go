package auth

import (
	"net/http"
)

// Authenticator system, with connector and digest admin authentication
type Authenticator struct {
	Connector AuthConnector
}

func NewAuthenticator(connector AuthConnector, realm string) *Authenticator {
	return &Authenticator{
		Connector: connector,
	}
}

// Interface for authentication that connectors implement
type AuthConnector interface {
	Authenticate(h http.Handler) http.Handler
}

// Authenticate method
func (a *Authenticator) Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var handler http.Handler
		if checkAdmin(r) {
			// Admin check
			handler = h

		} else {
			// Connector
			handler = a.Connector.Authenticate(h)
		}

		handler.ServeHTTP(w, r)
	})
}

// This method check if user is an admin
func checkAdmin(r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	// TODO rsoleto: Hay que cambiarlo para que utilice a través de un fichero de configuración o de BD
	if ok && username == "admin" && password == "admin" {
		return true
	}
	return false
}
