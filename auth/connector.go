package auth

import (
	"net/http"
)

// Authenticator system, with connector and basic admin authentication
type Authenticator struct {
	Connector     AuthConnector
	adminUser     string
	adminPassword string
}

// NewAuthenticator returns a configured Authenticator with associated connector
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

func (a *Authenticator) Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var handler http.Handler
		if isAdmin(r, a.adminUser, a.adminPassword) {
			// Admin check
			handler = h

		} else {
			// Connector
			handler = a.Connector.Authenticate(h)
		}

		handler.ServeHTTP(w, r)
	})
}

// GetAuthenticatedUser retrieves user from request
func (a *Authenticator) GetAuthenticatedUser(r *http.Request) (string, bool) {
	if isAdmin(r, a.adminUser, a.adminPassword) {
		return a.adminUser, true
	}
	return a.Connector.RetrieveUserID(*r), false
}

func isAdmin(r *http.Request, adminUser string, adminPassword string) bool {
	username, password, ok := r.BasicAuth()
	// Password is never stored in DB
	return ok && username == adminUser && password == adminPassword
}
