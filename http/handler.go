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
	AUTHORIZE_URL = API_VERSION_1 + "/authorize"
)

type WorkerHandler struct {
	worker *authorizr.Worker
}

type ProxyHandler struct {
	proxy  *authorizr.Proxy
	client *http.Client
}

func (a *WorkerHandler) TransactionLog(r *http.Request, authenticatedUser *api.AuthenticatedUser, status int, msg string) {

	// TODO: X-Forwarded headers
	//for header, _ := range r.Header {
	//	println(header, ": ", r.Header.Get(header))
	//}

	a.worker.Logger.WithFields(logrus.Fields{
		"RequestID": uuid.NewV4().String(),
		"Method":    r.Method,
		"URI":       r.RequestURI,
		"Address":   r.RemoteAddr,
		"User":      authenticatedUser.Identifier,
		"Status":    status,
	}).Info(msg)
}

// Handler returns an http.Handler for the APIs.
func WorkerHandlerRouter(worker *authorizr.Worker) http.Handler {
	// Create the muxer to handle the actual endpoints
	router := httprouter.New()

	workerHandler := WorkerHandler{worker: worker}

	// User api
	router.GET(USER_ROOT_URL, workerHandler.HandleGetUsers)
	router.POST(USER_ROOT_URL, workerHandler.HandlePostUsers)

	router.GET(USER_ID_URL, workerHandler.handleGetUserId)
	router.PUT(USER_ID_URL, workerHandler.HandlePutUser)
	router.DELETE(USER_ID_URL, workerHandler.handleDeleteUserId)

	router.GET(USER_ID_GROUPS_URL, workerHandler.handleUserIdGroups)

	// Special endpoint with organization URI for users
	router.GET(API_VERSION_1+ORG_ROOT+"/users", workerHandler.handleOrgListUsers)

	// Group api
	router.POST(GROUP_ORG_ROOT_URL, workerHandler.handleCreateGroup)
	router.GET(GROUP_ORG_ROOT_URL, workerHandler.handleListGroups)

	router.DELETE(GROUP_ID_URL, workerHandler.handleDeleteGroup)
	router.GET(GROUP_ID_URL, workerHandler.handleGetGroup)
	router.PUT(GROUP_ID_URL, workerHandler.handleUpdateGroup)

	router.GET(GROUP_ID_USERS_URL, workerHandler.handleListMembers)

	router.POST(GROUP_ID_USERS_ID_URL, workerHandler.handleAddMember)
	router.DELETE(GROUP_ID_USERS_ID_URL, workerHandler.handleRemoveMember)

	router.GET(GROUP_ID_POLICIES_URL, workerHandler.handleListAttachedGroupPolicies)

	router.POST(GROUP_ID_POLICIES_ID_URL, workerHandler.handleAttachGroupPolicy)
	router.DELETE(GROUP_ID_POLICIES_ID_URL, workerHandler.handleDetachGroupPolicy)

	// Special endpoint without organization URI for groups
	router.GET(API_VERSION_1+"/groups", workerHandler.handleListAllGroups)

	// Policy api
	router.GET(POLICY_ROOT_URL, workerHandler.handleListPolicies)
	router.POST(POLICY_ROOT_URL, workerHandler.handleCreatePolicy)

	router.DELETE(POLICY_ID_URL, workerHandler.handleDeletePolicy)
	router.GET(POLICY_ID_URL, workerHandler.handleGetPolicy)
	router.PUT(POLICY_ID_URL, workerHandler.handleUpdatePolicy)

	router.GET(POLICY_ID_GROUPS_URL, workerHandler.handleGetPolicyAttachedGroups)

	// Special endpoint without organization URI for policies
	router.GET(API_VERSION_1+"/policies", workerHandler.handleListAllPolicies)

	// Get effect endpoint
	router.POST(AUTHORIZE_URL, workerHandler.HandleAuthorizeResources)

	// Return handler
	return worker.Authenticator.Authenticate(router)
}

func (h *ProxyHandler) TransactionErrorLog(r *http.Request, transactionID string, msg string) {

	// TODO: X-Forwarded headers
	//for header, _ := range r.Header {
	//	println(header, ": ", r.Header.Get(header))
	//}

	h.proxy.Logger.WithFields(logrus.Fields{
		"RequestID": transactionID,
		"Method":    r.Method,
		"URI":       r.RequestURI,
		"Address":   r.RemoteAddr,
	}).Error(msg)
}

func (h *ProxyHandler) TransactionLog(r *http.Request, transactionID string, msg string) {

	// TODO: X-Forwarded headers
	//for header, _ := range r.Header {
	//	println(header, ": ", r.Header.Get(header))
	//}

	h.proxy.Logger.WithFields(logrus.Fields{
		"RequestID": transactionID,
		"Method":    r.Method,
		"URI":       r.RequestURI,
		"Address":   r.RemoteAddr,
	}).Info(msg)
}

// Handler returns an http.Handler for the Proxy.
func ProxyHandlerRouter(proxy *authorizr.Proxy) http.Handler {
	// Create the muxer to handle the actual endpoints
	router := httprouter.New()

	proxyHandler := ProxyHandler{proxy: proxy, client: http.DefaultClient}

	for _, res := range proxy.APIResources {
		router.Handle(res.Method, res.Url, proxyHandler.handleRequest(res))
	}

	return router
}

// HTTP responses

// 2xx RESPONSES

func (a *WorkerHandler) RespondOk(r *http.Request, authenticatedUser *api.AuthenticatedUser, w http.ResponseWriter, value interface{}) {
	b, err := json.Marshal(value)
	if err != nil {
		a.RespondInternalServerError(r, authenticatedUser, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
	a.TransactionLog(r, authenticatedUser, http.StatusOK, "Request processed")
}

func (a *WorkerHandler) RespondCreated(r *http.Request, authenticatedUser *api.AuthenticatedUser, w http.ResponseWriter, value interface{}) {
	b, err := json.Marshal(value)
	if err != nil {
		a.RespondInternalServerError(r, authenticatedUser, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
	a.TransactionLog(r, authenticatedUser, http.StatusCreated, "Request processed")
}

func (a *WorkerHandler) RespondNoContent(r *http.Request, authenticatedUser *api.AuthenticatedUser, w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
	a.TransactionLog(r, authenticatedUser, http.StatusNoContent, "Request processed")
}

// 4xx RESPONSES
func (a *WorkerHandler) RespondNotFound(r *http.Request, authenticatedUser *api.AuthenticatedUser, w http.ResponseWriter, apiError *api.Error) {
	w, err := writeErrorWithStatus(w, apiError, http.StatusNotFound)
	if err != nil {
		a.RespondInternalServerError(r, authenticatedUser, w)
		return
	}
	a.TransactionLog(r, authenticatedUser, http.StatusNotFound, "Request processed")
}

func (a *WorkerHandler) RespondBadRequest(r *http.Request, authenticatedUser *api.AuthenticatedUser, w http.ResponseWriter, apiError *api.Error) {
	w, err := writeErrorWithStatus(w, apiError, http.StatusBadRequest)
	if err != nil {
		a.RespondInternalServerError(r, authenticatedUser, w)
		return
	}
	a.TransactionLog(r, authenticatedUser, http.StatusBadRequest, "Bad Request")
}

func (a *WorkerHandler) RespondConflict(r *http.Request, authenticatedUser *api.AuthenticatedUser, w http.ResponseWriter, apiError *api.Error) {
	w, err := writeErrorWithStatus(w, apiError, http.StatusConflict)
	if err != nil {
		a.RespondInternalServerError(r, authenticatedUser, w)
		return
	}
	a.TransactionLog(r, authenticatedUser, http.StatusConflict, "Resource conflict")
}

func (a *WorkerHandler) RespondForbidden(r *http.Request, authenticatedUser *api.AuthenticatedUser, w http.ResponseWriter, apiError *api.Error) {
	w, err := writeErrorWithStatus(w, apiError, http.StatusForbidden)
	if err != nil {
		a.RespondInternalServerError(r, authenticatedUser, w)
		return
	}
	a.TransactionLog(r, authenticatedUser, http.StatusForbidden, "Forbidden")
}

// 5xx RESPONSES

func (a *WorkerHandler) RespondInternalServerError(r *http.Request, authenticatedUser *api.AuthenticatedUser, w http.ResponseWriter) {
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
