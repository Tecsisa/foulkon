package http

import (
	"github.com/julienschmidt/httprouter"
	"github.com/tecsisa/authorizr/authorizr"
	"net/http"
)

const (
	// Constants for values in url
	USER_ID   = "/:userid"
	GROUP_ID  = "/:groupid"
	POLICY_ID = "/:policyid"

	// API root reference
	API_ROOT      = "/api"
	API_VERSION_1 = API_ROOT + "/v1"

	// User API urls
	USER_ROOT_URL      = API_VERSION_1 + "/users"
	USER_ID_URL        = USER_ROOT_URL + USER_ID
	USER_ID_GROUPS_URL = USER_ID_URL + "/groups"

	// Group API urls
	GROUP_ROOT_URL           = API_VERSION_1 + "/groups"
	GROUP_ID_URL             = GROUP_ROOT_URL + GROUP_ID
	GROUP_ID_USERS_URL       = GROUP_ID_URL + "/users"
	GROUP_ID_USERS_ID_URL    = GROUP_ID_USERS_URL + USER_ID
	GROUP_ID_POLICIES_URL    = GROUP_ID_URL + "/policies"
	GROUP_ID_POLICIES_ID_URL = GROUP_ID_POLICIES_URL + POLICY_ID

	// Policy API urls
	POLICY_ROOT_URL = API_VERSION_1 + "/policies"
	POLICY_ID_URL   = POLICY_ROOT_URL + POLICY_ID
)

// Handler returns an http.Handler for the APIs.
func Handler(core *authorizr.Core) http.Handler {
	// Create the muxer to handle the actual endpoints
	router := httprouter.New()

	// User api
	userHandler := UserHandler{core: core}
	router.GET(USER_ROOT_URL, userHandler.handleGetUsers)
	router.POST(USER_ROOT_URL, userHandler.handlePostUsers)

	router.GET(USER_ID_URL, userHandler.handleGetUserId)
	router.DELETE(USER_ID_URL, userHandler.handleDeleteUserId)

	router.GET(USER_ID_GROUPS_URL, userHandler.handleUserIdGroups)

	// Group api
	groupHandler := GroupHandler{core: core}
	router.POST(GROUP_ROOT_URL, groupHandler.handleCreateGroup)
	router.GET(GROUP_ROOT_URL, groupHandler.handleListGroups)

	router.DELETE(GROUP_ID_URL, groupHandler.handleDeleteGroup)
	router.GET(GROUP_ID_URL, groupHandler.handleGetGroup)
	router.PUT(GROUP_ID_URL, groupHandler.handleUpdateGroup)

	router.GET(GROUP_ID_USERS_URL, groupHandler.handleListMembers)
	router.POST(GROUP_ID_USERS_URL, groupHandler.handleAddMember)

	router.DELETE(GROUP_ID_USERS_ID_URL, groupHandler.handleRemoveMember)

	router.POST(GROUP_ID_POLICIES_URL, groupHandler.handleAttachGroupPolicy)
	router.GET(GROUP_ID_POLICIES_URL, groupHandler.handleListAtachhedGroupPolicies)

	router.DELETE(GROUP_ID_POLICIES_ID_URL, groupHandler.handleDetachGroupPolicy)

	// Policy api
	policyHandler := PolicyHandler{core: core}
	router.GET(POLICY_ROOT_URL, policyHandler.handleListPolicies)
	router.POST(POLICY_ROOT_URL, policyHandler.handleCreatePolicy)

	router.DELETE(POLICY_ID_URL, policyHandler.handleDeletePolicy)
	router.GET(POLICY_ID_URL, policyHandler.handleGetPolicy)

	// Return handler
	return router
}
