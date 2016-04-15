package http

import (
	"net/http"
	"github.com/tecsisa/authorizr/authorizr"
)

// Handler returns an http.Handler for the APIs.
func Handler(core *authorizr.Core) http.Handler {
	// Create the muxer to handle the actual endpoints
	mux := http.NewServeMux()

	// User api
	mux.Handle("/users", handleGetUsers(core))

	// Group api
	mux.Handle("/groups", handleGetGroups(core))

	// Policy api
	mux.Handle("/policies", handleGetPolicy(core))

	// Create request handler
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mux.ServeHTTP(w, req)
		return
	})
}