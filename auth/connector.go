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
}

// Authenticate method
func (a *Authenticator) Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var handler http.Handler
		if checkAdmin(r, a.adminUser, a.adminPassword) {
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
func checkAdmin(r *http.Request, adminUser string, adminPassword string) bool {
	username, password, ok := r.BasicAuth()
	// Password is never stored in DB
	if ok && username == adminUser && password == adminPassword {
		return true
	}
	return false
}
