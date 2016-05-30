package http

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/authorizr"
)

const (
	// Constants for values in url
	USER_ID     = "userid"
	GROUP_NAME  = "groupname"
	POLICY_NAME = "policyname"
	ORG_NAME    = "orgname"

	// URI Path param prefix
	URI_PATH_PREFIX = "/:"

	// API root reference
	API_ROOT      = "/api"
	API_VERSION_1 = API_ROOT + "/v1"

	// Organization API ROOT
	ORG_ROOT = "/organizations/:" + ORG_NAME

	// User API urls
	USER_ROOT_URL      = API_VERSION_1 + "/users"
	USER_ID_URL        = USER_ROOT_URL + URI_PATH_PREFIX + USER_ID
	USER_ID_GROUPS_URL = USER_ID_URL + "/groups"

	// Group organization API urls
	GROUP_ORG_ROOT_URL       = API_VERSION_1 + ORG_ROOT + "/groups"
	GROUP_ID_URL             = GROUP_ORG_ROOT_URL + URI_PATH_PREFIX + GROUP_NAME
	GROUP_ID_USERS_URL       = GROUP_ID_URL + "/users"
	GROUP_ID_USERS_ID_URL    = GROUP_ID_USERS_URL + URI_PATH_PREFIX + USER_ID
	GROUP_ID_POLICIES_URL    = GROUP_ID_URL + "/policies"
	GROUP_ID_POLICIES_ID_URL = GROUP_ID_POLICIES_URL + URI_PATH_PREFIX + POLICY_NAME

	// Policy API urls
	POLICY_ROOT_URL      = API_VERSION_1 + ORG_ROOT + "/policies"
	POLICY_ID_URL        = POLICY_ROOT_URL + URI_PATH_PREFIX + POLICY_NAME
	POLICY_ID_GROUPS_URL = POLICY_ROOT_URL + URI_PATH_PREFIX + POLICY_NAME + "/groups"

	// Authorization URLs
	EFFECT_URL = API_VERSION_1 + "/resources"
)

type AuthHandler struct {
	core *authorizr.Core
}

func (a *AuthHandler) TransactionLog(r *http.Request, authenticatedUser *api.AuthenticatedUser, status int, msg string) {

	// TODO: X-Forwarded headers
	//for header, _ := range r.Header {
	//	println(header, ": ", r.Header.Get(header))
	//}

	a.core.Logger.WithFields(logrus.Fields{
		"RequestID": uuid.NewV4().String(),
		"Method":    r.Method,
		"URI":       r.RequestURI,
		"Address":   r.RemoteAddr,
		"User":      authenticatedUser.Identifier,
		"Status":    status,
	}).Info(msg)
}

// Handler returns an http.Handler for the APIs.
func Handler(core *authorizr.Core) http.Handler {
	// Create the muxer to handle the actual endpoints
	router := httprouter.New()

	authHandler := AuthHandler{core: core}

	// User api
	router.GET(USER_ROOT_URL, authHandler.handleGetUsers)
	router.POST(USER_ROOT_URL, authHandler.handlePostUsers)

	router.GET(USER_ID_URL, authHandler.handleGetUserId)
	router.PUT(USER_ID_URL, authHandler.handlePutUser)
	router.DELETE(USER_ID_URL, authHandler.handleDeleteUserId)

	router.GET(USER_ID_GROUPS_URL, authHandler.handleUserIdGroups)

	// Special endpoint with organization URI for users
	router.GET(API_VERSION_1+ORG_ROOT+"/users", authHandler.handleOrgListUsers)

	// Group api
	router.POST(GROUP_ORG_ROOT_URL, authHandler.handleCreateGroup)
	router.GET(GROUP_ORG_ROOT_URL, authHandler.handleListGroups)

	router.DELETE(GROUP_ID_URL, authHandler.handleDeleteGroup)
	router.GET(GROUP_ID_URL, authHandler.handleGetGroup)
	router.PUT(GROUP_ID_URL, authHandler.handleUpdateGroup)

	router.GET(GROUP_ID_USERS_URL, authHandler.handleListMembers)

	router.POST(GROUP_ID_USERS_ID_URL, authHandler.handleAddMember)
	router.DELETE(GROUP_ID_USERS_ID_URL, authHandler.handleRemoveMember)

	router.GET(GROUP_ID_POLICIES_URL, authHandler.handleListAttachedGroupPolicies)

	router.POST(GROUP_ID_POLICIES_ID_URL, authHandler.handleAttachGroupPolicy)
	router.DELETE(GROUP_ID_POLICIES_ID_URL, authHandler.handleDetachGroupPolicy)

	// Special endpoint without organization URI for groups
	router.GET(API_VERSION_1+"/groups", authHandler.handleListAllGroups)

	// Policy api
	router.GET(POLICY_ROOT_URL, authHandler.handleListPolicies)
	router.POST(POLICY_ROOT_URL, authHandler.handleCreatePolicy)

	router.DELETE(POLICY_ID_URL, authHandler.handleDeletePolicy)
	router.GET(POLICY_ID_URL, authHandler.handleGetPolicy)
	router.PUT(POLICY_ID_URL, authHandler.handleUpdatePolicy)

	router.GET(POLICY_ID_GROUPS_URL, authHandler.handleGetPolicyAttachedGroups)

	// Special endpoint without organization URI for policies
	router.GET(API_VERSION_1+"/policies", authHandler.handleListAllPolicies)

	// Get effect endpoint
	router.GET(EFFECT_URL, authHandler.handleGetEffectByUserActionResource)

	// Return handler
	return core.Authenticator.Authenticate(router)
}

// HTTP responses

// 2xx RESPONSES

func (a *AuthHandler) RespondOk(r *http.Request, authenticatedUser *api.AuthenticatedUser, w http.ResponseWriter, value interface{}) {
	b, err := json.Marshal(value)
	if err != nil {
		a.RespondInternalServerError(r, authenticatedUser, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
	a.TransactionLog(r, authenticatedUser, http.StatusOK, "Request processed")

}

func (a *AuthHandler) RespondNoContent(r *http.Request, authenticatedUser *api.AuthenticatedUser, w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
	a.TransactionLog(r, authenticatedUser, http.StatusNoContent, "Request processed")
}

// 4xx RESPONSES
func (a *AuthHandler) RespondNotFound(r *http.Request, authenticatedUser *api.AuthenticatedUser, w http.ResponseWriter, apiError *api.Error) {
	w, err := writeErrorWithStatus(w, apiError, http.StatusNotFound)
	if err != nil {
		a.RespondInternalServerError(r, authenticatedUser, w)
		return
	}
	a.TransactionLog(r, authenticatedUser, http.StatusNotFound, "Request processed")
}

func (a *AuthHandler) RespondBadRequest(r *http.Request, authenticatedUser *api.AuthenticatedUser, w http.ResponseWriter, apiError *api.Error) {
	w, err := writeErrorWithStatus(w, apiError, http.StatusBadRequest)
	if err != nil {
		a.RespondInternalServerError(r, authenticatedUser, w)
		return
	}
	a.TransactionLog(r, authenticatedUser, http.StatusBadRequest, "Bad Request")
}

func (a *AuthHandler) RespondConflict(r *http.Request, authenticatedUser *api.AuthenticatedUser, w http.ResponseWriter, apiError *api.Error) {
	w, err := writeErrorWithStatus(w, apiError, http.StatusConflict)
	if err != nil {
		a.RespondInternalServerError(r, authenticatedUser, w)
		return
	}
	a.TransactionLog(r, authenticatedUser, http.StatusConflict, "Resource conflict")
}

func (a *AuthHandler) RespondForbidden(r *http.Request, authenticatedUser *api.AuthenticatedUser, w http.ResponseWriter, apiError *api.Error) {
	w, err := writeErrorWithStatus(w, apiError, http.StatusForbidden)
	if err != nil {
		a.RespondInternalServerError(r, authenticatedUser, w)
		return
	}
	a.TransactionLog(r, authenticatedUser, http.StatusForbidden, "Forbidden")
}

// 5xx RESPONSES

func (a *AuthHandler) RespondInternalServerError(r *http.Request, authenticatedUser *api.AuthenticatedUser, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	a.TransactionLog(r, authenticatedUser, http.StatusInternalServerError, "Server error")
}

// Private Helper Methods
func writeErrorWithStatus(w http.ResponseWriter, apiError *api.Error, statusCode int) (http.ResponseWriter, error) {
	b, err := json.Marshal(apiError)
	if err != nil {
		return nil, err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(b)
	return w, nil
}
