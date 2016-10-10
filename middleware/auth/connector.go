package auth

import (
	"net/http"

	"github.com/Tecsisa/foulkon/middleware"
)

// Authenticator middleware system, with connector and basic admin authentication
type AuthenticatorMiddleware struct {
	connector     AuthConnector
	adminUser     string
	adminPassword string
}

// NewAuthenticator returns a configured AuthenticatorMiddleware with associated connector
func NewAuthenticatorMiddleware(connector AuthConnector, adminUser string, adminPassword string) *AuthenticatorMiddleware {
	return &AuthenticatorMiddleware{
		connector:     connector,
		adminUser:     adminUser,
		adminPassword: adminPassword,
	}
}

// Interface for authentication that connectors implement
type AuthConnector interface {
	Authenticate(next http.Handler) http.Handler
	RetrieveUserID(r http.Request) string
}

func (a *AuthenticatorMiddleware) Action(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var handler http.Handler
		if isAdmin(r, a.adminUser, a.adminPassword) {
			// Admin check
			r.Header.Add(middleware.USER_ID_HEADER, a.adminUser)
			handler = next
		} else {
			// Connector
			handler = a.connector.Authenticate(next)
		}

		handler.ServeHTTP(w, r)
	})
}

func (a *AuthenticatorMiddleware) GetInfo(r *http.Request, mc *middleware.MiddlewareContext) {
	mc.UserId, mc.Admin = a.getAuthenticatedUser(r)
}

// GetAuthenticatedUser retrieves user from request
func (a *AuthenticatorMiddleware) getAuthenticatedUser(r *http.Request) (string, bool) {
	if isAdmin(r, a.adminUser, a.adminPassword) {
		return a.adminUser, true
	}
	return a.connector.RetrieveUserID(*r), false
}

func isAdmin(r *http.Request, adminUser string, adminPassword string) bool {
	username, password, ok := r.BasicAuth()
	// Password is never stored in DB
	return ok && username == adminUser && password == adminPassword
}
