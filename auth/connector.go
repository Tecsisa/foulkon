package auth

import (
	"net/http"

	"github.com/tecsisa/authorizr/api"
)

// Authenticator system, with connector and basic admin authentication
type Authenticator struct {
	Connector     AuthConnector
	adminUser     string
	adminPassword string
}

// Returns a configured Authenticator with associated connector
func NewAuthenticator(connector AuthConnector, adminUser string, adminPassword string) *Authenticator {
	return &Authenticator{
		Connector:     connector,
		adminUser:     adminUser,
		adminPassword: adminPassword,
	}
}

// Interface for authentication that connectors implement
type AuthConnector interface {
	Authenticate(h http.Handler) http.Handler
	RetrieveUserID(r http.Request) string
}

// Authenticate method. This checks all problem related to user authentication, so if something goes wrong, it return a error message.
func (a *Authenticator) Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var handler http.Handler
		if checkAdmin(*r, a.adminUser, a.adminPassword) {
			// Admin check
			handler = h

		} else {
			// Connector
			handler = a.Connector.Authenticate(h)
		}

		handler.ServeHTTP(w, r)
	})
}

// Retrieve user from request. This method never fail because Authentication deal with all problems related to this
func (a *Authenticator) RetrieveUserID(r http.Request) api.AuthenticatedUser {
	if checkAdmin(r, a.adminUser, a.adminPassword) {
		return api.AuthenticatedUser{
			Identifier: a.adminUser,
			Admin:      true,
		}
	} else {
		return api.AuthenticatedUser{
			Identifier: a.Connector.RetrieveUserID(r),
			Admin:      false,
		}
	}
}

// This method check if user is an admin
func checkAdmin(r http.Request, adminUser string, adminPassword string) bool {
	username, password, ok := r.BasicAuth()
	// Password is never stored in DB
	if ok && username == adminUser && password == adminPassword {
		return true
	}
	return false
}
