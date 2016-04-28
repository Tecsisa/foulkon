package http

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tecsisa/authorizr/authorizr"
)

const (
	// Constants for values in url
	USER_ID   = "userid"
	GROUP_ID  = "groupid"
	POLICY_ID = "policyid"
	ORG_ID    = "orgid"

	// URI Path param prefix
	URI_PATH_PREFIX = "/:"

	// API root reference
	API_ROOT      = "/api"
	API_VERSION_1 = API_ROOT + "/v1"

	// Organization API ROOT
	ORG_ROOT = "/organization/:" + ORG_ID

	// User API urls
	USER_ROOT_URL      = API_VERSION_1 + "/users"
	USER_ID_URL        = USER_ROOT_URL + URI_PATH_PREFIX + USER_ID
	USER_ID_GROUPS_URL = USER_ID_URL + "/groups"

	// Group organization API urls
	GROUP_ORG_ROOT_URL       = API_VERSION_1 + ORG_ROOT + "/groups"
	GROUP_ID_URL             = GROUP_ORG_ROOT_URL + URI_PATH_PREFIX + GROUP_ID
	GROUP_ID_USERS_URL       = GROUP_ID_URL + "/users"
	GROUP_ID_USERS_ID_URL    = GROUP_ID_USERS_URL + URI_PATH_PREFIX + USER_ID
	GROUP_ID_POLICIES_URL    = GROUP_ID_URL + "/policies"
	GROUP_ID_POLICIES_ID_URL = GROUP_ID_POLICIES_URL + URI_PATH_PREFIX + POLICY_ID

	// Policy API urls
	POLICY_ROOT_URL      = API_VERSION_1 + ORG_ROOT + "/policies"
	POLICY_ID_URL        = POLICY_ROOT_URL + URI_PATH_PREFIX + POLICY_ID
	POLICY_ID_GROUPS_URL = POLICY_ROOT_URL + URI_PATH_PREFIX + POLICY_ID + "/groups"
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
	router.PUT(USER_ID_URL, userHandler.handlePutUser)
	router.DELETE(USER_ID_URL, userHandler.handleDeleteUserId)

	router.GET(USER_ID_GROUPS_URL, userHandler.handleUserIdGroups)

	// Special endpoint with organization URI for users
	router.GET(API_VERSION_1+ORG_ROOT+"/users", userHandler.handleOrgListUsers)

	// Group api
	groupHandler := GroupHandler{core: core}
	router.POST(GROUP_ORG_ROOT_URL, groupHandler.handleCreateGroup)
	router.GET(GROUP_ORG_ROOT_URL, groupHandler.handleListGroups)

	router.DELETE(GROUP_ID_URL, groupHandler.handleDeleteGroup)
	router.GET(GROUP_ID_URL, groupHandler.handleGetGroup)
	router.PUT(GROUP_ID_URL, groupHandler.handleUpdateGroup)

	router.GET(GROUP_ID_USERS_URL, groupHandler.handleListMembers)
	router.POST(GROUP_ID_USERS_URL, groupHandler.handleAddMember)

	router.DELETE(GROUP_ID_USERS_ID_URL, groupHandler.handleRemoveMember)

	router.POST(GROUP_ID_POLICIES_URL, groupHandler.handleAttachGroupPolicy)
	router.GET(GROUP_ID_POLICIES_URL, groupHandler.handleListAtachhedGroupPolicies)

	router.DELETE(GROUP_ID_POLICIES_ID_URL, groupHandler.handleDetachGroupPolicy)

	// Special endpoint without organization URI for groups
	router.GET(API_VERSION_1+"/groups", groupHandler.handleListAllGroups)

	// Policy api
	policyHandler := PolicyHandler{core: core}
	router.GET(POLICY_ROOT_URL, policyHandler.handleListPolicies)
	router.POST(POLICY_ROOT_URL, policyHandler.handleCreatePolicy)

	router.DELETE(POLICY_ID_URL, policyHandler.handleDeletePolicy)
	router.GET(POLICY_ID_URL, policyHandler.handleGetPolicy)

	router.GET(POLICY_ID_GROUPS_URL, policyHandler.handleGetPolicyAttachedGroups)

	// Special endpoint without organization URI for policies
	router.GET(API_VERSION_1+"/policies", policyHandler.handleListAllPolicies)

	// Return handler
	return router
}

// HTTP responses

// 2xx RESPONSES

func RespondOk(w http.ResponseWriter, value interface{}) {
	b, err := json.Marshal(value)
	if err != nil {
		RespondInternalServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)

}

func RespondNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// 4xx RESPONSES
func RespondNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func RespondBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func RespondConflict(w http.ResponseWriter) {
	w.WriteHeader(http.StatusConflict)
}

// 5xx RESPONSES

func RespondInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}
