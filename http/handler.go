package http

import (
	"encoding/json"
	"net/http"

	"fmt"
	"strconv"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/foulkon"
	"github.com/julienschmidt/httprouter"
)

const (
	// Constants for values in url
	USER_ID             = "userid"
	GROUP_NAME          = "groupname"
	POLICY_NAME         = "policyname"
	PROXY_RESOURCE_NAME = "proxyresourcename"
	AUTH_PROVIDER_NAME  = "authprovidername"
	ORG_NAME            = "orgname"

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

	// Proxy resource API urls
	PROXY_RESOURCE_ROOT_URL = API_VERSION_1 + ORG_ROOT + "/proxy-resources"
	PROXY_RESOURCE_ID_URL   = PROXY_RESOURCE_ROOT_URL + URI_PATH_PREFIX + PROXY_RESOURCE_NAME

	// Authorization URLs
	RESOURCE_URL = API_VERSION_1 + "/resource"

	// Admin URLs
	ADMIN_ROOT = "/admin"

	// Admin OIDC Authentication API URLs
	OIDC_AUTH_ROOT_URL = API_VERSION_1 + ADMIN_ROOT + "/auth/oidc/providers"
	OIDC_AUTH_ID_URL   = OIDC_AUTH_ROOT_URL + URI_PATH_PREFIX + AUTH_PROVIDER_NAME

	// Foulkon configuration URL
	ABOUT = "/about"
)

// PROXY

type ProxyHandler struct {
	proxy  *foulkon.Proxy
	client *http.Client
}

// WORKER

type WorkerHandler struct {
	worker *foulkon.Worker
}

func (wh *WorkerHandler) processHttpRequest(r *http.Request, w http.ResponseWriter, ps httprouter.Params, request interface{}) (
	requestInfo api.RequestInfo, filterData *api.Filter, apiError *api.Error) {
	// Get Request Info
	requestInfo = wh.getRequestInfo(r)
	// Decode request if passed
	if request != nil {
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			apiError = &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: err.Error(),
			}
			api.LogOperationError(requestInfo.RequestID, requestInfo.Identifier, apiError)
		}
	}
	filterData, err := getFilterData(r, ps)
	if err != nil {
		apiError = err.(*api.Error)
		api.LogOperationError(requestInfo.RequestID, requestInfo.Identifier, apiError)
	}
	return requestInfo, filterData, apiError
}

func (wh *WorkerHandler) processHttpResponse(r *http.Request, w http.ResponseWriter, requestInfo api.RequestInfo, response interface{}, err error, responseCode int) {
	if err != nil {
		// Transform to API errors
		apiError := err.(*api.Error)
		api.LogOperationError(requestInfo.RequestID, requestInfo.Identifier, apiError)
		var statusCode int
		switch apiError.Code {
		case api.USER_ALREADY_EXIST, api.GROUP_ALREADY_EXIST,
			api.USER_IS_ALREADY_A_MEMBER_OF_GROUP,
			api.PROXY_RESOURCE_ALREADY_EXIST,
			api.POLICY_IS_ALREADY_ATTACHED_TO_GROUP, api.POLICY_ALREADY_EXIST,
			api.PROXY_RESOURCES_ROUTES_CONFLICT,
			api.AUTH_OIDC_PROVIDER_ALREADY_EXIST:
			// A conflict occurs
			statusCode = http.StatusConflict
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			// No authorization success
			statusCode = http.StatusForbidden
		case api.USER_BY_EXTERNAL_ID_NOT_FOUND, api.GROUP_BY_ORG_AND_NAME_NOT_FOUND,
			api.USER_IS_NOT_A_MEMBER_OF_GROUP, api.POLICY_IS_NOT_ATTACHED_TO_GROUP,
			api.POLICY_BY_ORG_AND_NAME_NOT_FOUND, api.PROXY_RESOURCE_BY_ORG_AND_NAME_NOT_FOUND,
			api.AUTH_OIDC_PROVIDER_BY_NAME_NOT_FOUND:
			// Resource or relation not found
			statusCode = http.StatusNotFound
		case api.INVALID_PARAMETER_ERROR, api.REGEX_NO_MATCH:
			// Unexpected input in validation parameters
			statusCode = http.StatusBadRequest
		default: // Unexpected API error
			statusCode = http.StatusInternalServerError
		}
		WriteHttpResponse(r, w, requestInfo.RequestID, requestInfo.Identifier, statusCode, apiError)
		return
	}

	// Write response data if everything is ok
	WriteHttpResponse(r, w, requestInfo.RequestID, requestInfo.Identifier, responseCode, response)
}

func (wh *WorkerHandler) getRequestInfo(r *http.Request) api.RequestInfo {
	// Retrieve request information from middleware context
	mc := wh.worker.MiddlewareHandler.GetMiddlewareContext(r)
	return api.RequestInfo{
		Identifier: mc.UserId,
		Admin:      mc.Admin,
		RequestID:  mc.XRequestId,
	}
}

// WorkerHandlerRouter returns http.Handler for the APIs.
func WorkerHandlerRouter(worker *foulkon.Worker) http.Handler {
	// Create the muxer to handle the actual endpoints
	router := httprouter.New()

	workerHandler := WorkerHandler{worker: worker}

	// User api
	router.GET(USER_ROOT_URL, workerHandler.HandleListUsers)
	router.POST(USER_ROOT_URL, workerHandler.HandleAddUser)

	router.GET(USER_ID_URL, workerHandler.HandleGetUserByExternalID)
	router.PUT(USER_ID_URL, workerHandler.HandleUpdateUser)
	router.DELETE(USER_ID_URL, workerHandler.HandleRemoveUser)

	router.GET(USER_ID_GROUPS_URL, workerHandler.HandleListGroupsByUser)

	// Group api
	router.POST(GROUP_ORG_ROOT_URL, workerHandler.HandleAddGroup)
	router.GET(GROUP_ORG_ROOT_URL, workerHandler.HandleListGroups)

	router.DELETE(GROUP_ID_URL, workerHandler.HandleRemoveGroup)
	router.GET(GROUP_ID_URL, workerHandler.HandleGetGroupByName)
	router.PUT(GROUP_ID_URL, workerHandler.HandleUpdateGroup)

	router.GET(GROUP_ID_USERS_URL, workerHandler.HandleListMembers)

	router.POST(GROUP_ID_USERS_ID_URL, workerHandler.HandleAddMember)
	router.DELETE(GROUP_ID_USERS_ID_URL, workerHandler.HandleRemoveMember)

	router.GET(GROUP_ID_POLICIES_URL, workerHandler.HandleListAttachedGroupPolicies)

	router.POST(GROUP_ID_POLICIES_ID_URL, workerHandler.HandleAttachPolicyToGroup)
	router.DELETE(GROUP_ID_POLICIES_ID_URL, workerHandler.HandleDetachPolicyToGroup)

	// Special endpoint without organization URI for groups
	router.GET(API_VERSION_1+"/groups", workerHandler.HandleListAllGroups)

	// Policy api
	router.GET(POLICY_ROOT_URL, workerHandler.HandleListPolicies)
	router.POST(POLICY_ROOT_URL, workerHandler.HandleAddPolicy)

	router.DELETE(POLICY_ID_URL, workerHandler.HandleRemovePolicy)

	router.GET(POLICY_ID_URL, workerHandler.HandleGetPolicyByName)
	router.PUT(POLICY_ID_URL, workerHandler.HandleUpdatePolicy)

	router.GET(POLICY_ID_GROUPS_URL, workerHandler.HandleListAttachedGroups)

	// Special endpoint without organization URI for policies
	router.GET(API_VERSION_1+"/policies", workerHandler.HandleListAllPolicies)

	// Proxy Resources api
	router.GET(PROXY_RESOURCE_ROOT_URL, workerHandler.HandleListProxyResource)
	router.POST(PROXY_RESOURCE_ROOT_URL, workerHandler.HandleAddProxyResource)

	router.DELETE(PROXY_RESOURCE_ID_URL, workerHandler.HandleRemoveProxyResource)

	router.GET(PROXY_RESOURCE_ID_URL, workerHandler.HandleGetProxyResourceByName)
	router.PUT(PROXY_RESOURCE_ID_URL, workerHandler.HandleUpdateProxyResource)

	// Resources authorized endpoint
	router.POST(RESOURCE_URL, workerHandler.HandleGetAuthorizedExternalResources)

	// OIDC authentication api
	router.GET(OIDC_AUTH_ROOT_URL, workerHandler.HandleListOidcProviders)
	router.POST(OIDC_AUTH_ROOT_URL, workerHandler.HandleAddOidcProvider)

	router.DELETE(OIDC_AUTH_ID_URL, workerHandler.HandleRemoveOidcProvider)

	router.GET(OIDC_AUTH_ID_URL, workerHandler.HandleGetOidcProviderByName)
	router.PUT(OIDC_AUTH_ID_URL, workerHandler.HandleUpdateOidcProvider)

	// Current Foulkon configuration
	router.GET(ABOUT, workerHandler.HandleGetCurrentConfig)

	return workerHandler.worker.MiddlewareHandler.Handle(router)
}

// WriteHttpResponse fill a http response with data, controlling marshalling errors
func WriteHttpResponse(r *http.Request, w http.ResponseWriter, requestId string, userId string, statusCode int, value interface{}) {
	if value != nil {
		b, err := json.Marshal(value)
		if err != nil {
			apiErr := &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: err.Error(),
			}
			api.TransactionResponseErrorLog(requestId, userId, r, http.StatusInternalServerError, apiErr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		// Set status code
		w.WriteHeader(statusCode)
		w.Write(b)
		return
	}

	// Set status code by default if interface isn't defined
	w.WriteHeader(statusCode)
}

// Private Helper Methods

func getFilterData(r *http.Request, ps httprouter.Params) (*api.Filter, error) {
	var err error
	// Retrieve Offset
	var offset int
	offs := r.URL.Query().Get("Offset")
	if len(offs) != 0 {
		offset, err = strconv.Atoi(offs)
		if err != nil || offset < 0 {
			return nil, &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: fmt.Sprintf("Invalid parameter: Offset %v", offs),
			}
		}
	}
	// Retrieve Limit
	var limit int
	lmt := r.URL.Query().Get("Limit")
	if len(lmt) != 0 {
		limit, err = strconv.Atoi(lmt)
		if err != nil || limit < 0 {
			return nil, &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: fmt.Sprintf("Invalid parameter: Limit %v", lmt),
			}
		}
	}

	// Retrieve Org
	var org string
	if org = ps.ByName(ORG_NAME); len(org) == 0 {
		org = r.URL.Query().Get("Org")
	}

	return &api.Filter{
		PathPrefix:        r.URL.Query().Get("PathPrefix"),
		Org:               org,
		ExternalID:        ps.ByName(USER_ID),
		PolicyName:        ps.ByName(POLICY_NAME),
		GroupName:         ps.ByName(GROUP_NAME),
		ProxyResourceName: ps.ByName(PROXY_RESOURCE_NAME),
		AuthProviderName:  ps.ByName(AUTH_PROVIDER_NAME),
		Offset:            offset,
		Limit:             limit,
		OrderBy:           r.URL.Query().Get("OrderBy"),
	}, nil
}
