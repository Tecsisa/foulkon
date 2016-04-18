package http

import (
	"github.com/julienschmidt/httprouter"
	"github.com/tecsisa/authorizr/authorizr"
	"net/http"
)

// Handler returns an http.Handler for the APIs.
func Handler(core *authorizr.Core) http.Handler {
	// Create the muxer to handle the actual endpoints
	router := httprouter.New()

	// User api
	userHandler := UserHandler{core: core}
	router.GET("/users", userHandler.handleGetUsers)
	router.POST("/users", userHandler.handlePostUsers)

	router.GET("/users/:id", userHandler.handleGetUserId)
	router.DELETE("/users/:id", userHandler.handleDeleteUserId)

	router.GET("/users/:id/groups", userHandler.handleUserIdGroups)

	// Group api
	//router.GET("/groups", handleGroups(core))

	// Policy api
	//router.GET("/policies", handlePolicies(core))

	// Return handler
	return router
}
