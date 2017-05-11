package http

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"bytes"
	"fmt"

	"github.com/Tecsisa/foulkon/api"
	"github.com/stretchr/testify/assert"
)

func TestProxyHandler_HandleRequest(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		expectedStatusCode int
		expectedError      *api.Error
		expectedResponse   api.User
		resource           string
		// Manager Results
		getListUsersResult                   []string
		getUserByExternalIdResult            *api.User
		getAuthorizedExternalResourcesResult []string
		// Manager Errors
		getListUsersErr                   error
		getUserByExternalIdErr            error
		getAuthorizedExternalResourcesErr error
		// Authentication
		authStatusCode int
	}{
		"OkCaseAdmin": {
			expectedStatusCode: http.StatusOK,
			resource:           USER_ROOT_URL + "/user",
			expectedResponse: api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
				UpdateAt:   now,
			},
			getUserByExternalIdResult: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
				UpdateAt:   now,
			},
			getAuthorizedExternalResourcesResult: []string{"urn:ews:example:instance1:resource/user"},
		},
		"ErrorCaseInvalidParameter": {
			expectedStatusCode: http.StatusBadRequest,
			resource:           "/urnPrefix",
			expectedError: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Bad request",
			},
		},
		"ErrorCaseInvalidResources": {
			expectedStatusCode: http.StatusBadRequest,
			resource:           "/invalidUrn",
			expectedError: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Bad request",
			},
		},
		"ErrorCaseInvalidAction": {
			expectedStatusCode: http.StatusBadRequest,
			resource:           "/invalidAction",
			expectedError: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Bad request",
			},
		},
		"ErrorCaseInvalidHost": {
			expectedStatusCode: http.StatusInternalServerError,
			resource:           "/invalid",
			expectedError: &api.Error{
				Code:    INVALID_DEST_HOST_URL,
				Message: "Error creating destination host",
			},
			getAuthorizedExternalResourcesResult: []string{"urn:ews:example:instance1:resource/invalid"},
		},
		"ErrorCaseHostUnreachable": {
			expectedStatusCode: http.StatusInternalServerError,
			resource:           "/fail",
			expectedError: &api.Error{
				Code:    HOST_UNREACHABLE,
				Message: "Error calling destination resource",
			},
			getAuthorizedExternalResourcesResult: []string{"urn:ews:example:instance1:resource/fail"},
		},
		"ErrorCaseUnauthenticated": {
			expectedStatusCode: http.StatusForbidden,
			authStatusCode:     http.StatusUnauthorized,
			resource:           USER_ROOT_URL + "/user",
			expectedError: &api.Error{
				Code:    FORBIDDEN_ERROR,
				Message: "Forbidden resource. If you need access, contact the administrator",
			},
			getAuthorizedExternalResourcesResult: []string{"urn:ews:example:instance1:resource/user"},
		},
		"ErrorCaseForbidden": {
			expectedStatusCode: http.StatusForbidden,
			resource:           USER_ROOT_URL + "/user",
			getUserByExternalIdResult: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
				UpdateAt:   now,
			},
			getAuthorizedExternalResourcesErr: &api.Error{
				Code: api.UNAUTHORIZED_RESOURCES_ERROR,
			},
			expectedError: &api.Error{
				Code:    FORBIDDEN_ERROR,
				Message: "Forbidden resource. If you need access, contact the administrator",
			},
		},
		"ErrorCaseForbiddenDifferentResources": {
			expectedStatusCode: http.StatusForbidden,
			resource:           USER_ROOT_URL + "/user",
			getUserByExternalIdResult: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
				UpdateAt:   now,
			},
			getAuthorizedExternalResourcesResult: []string{"urn:ews:example:instance1:resource/forbidden"},
			expectedError: &api.Error{
				Code:    FORBIDDEN_ERROR,
				Message: "Forbidden resource. If you need access, contact the administrator",
			},
		},
		"ErrorCaseAuthBadRequest": {
			expectedStatusCode: http.StatusBadRequest,
			authStatusCode:     http.StatusBadRequest,
			resource:           USER_ROOT_URL + "/user",
			expectedError: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Bad request",
			},
		},
		"ErrorCaseWorkerError": {
			expectedStatusCode: http.StatusInternalServerError,
			resource:           USER_ROOT_URL + "/user",
			expectedError: &api.Error{
				Code:    INTERNAL_SERVER_ERROR,
				Message: "Internal server error. Contact the administrator",
			},
			getAuthorizedExternalResourcesErr: &api.Error{
				Code: api.UNKNOWN_API_ERROR,
			},
		},
		"ErrorCaseUnauthorized": {
			expectedStatusCode: http.StatusForbidden,
			resource:           USER_ROOT_URL + "/user",
			expectedError: &api.Error{
				Code:    FORBIDDEN_ERROR,
				Message: "Forbidden resource. If you need access, contact the administrator",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[GetAuthorizedExternalResourcesMethod][0] = test.getAuthorizedExternalResourcesResult
		testApi.ArgsOut[GetAuthorizedExternalResourcesMethod][1] = test.getAuthorizedExternalResourcesErr
		testApi.ArgsOut[GetUserByExternalIdMethod][0] = test.getUserByExternalIdResult
		testApi.ArgsOut[GetUserByExternalIdMethod][1] = test.getUserByExternalIdErr

		req, err := http.NewRequest(http.MethodGet, proxy.URL+test.resource, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		if test.authStatusCode != 0 {
			authConnector.statusCode = test.authStatusCode
		}

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusOK:
			response := api.User{}
			err = json.NewDecoder(res.Body).Decode(&response)
			if err != nil {
				t.Fatalf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			assert.Equal(t, test.expectedResponse, response, "Error in test case %v", n)
		default:
			if test.expectedError != nil {
				apiError := &api.Error{}
				err = json.NewDecoder(res.Body).Decode(&apiError)
				assert.Nil(t, err, "Error in test case %v", n)
				// Check error
				assert.Equal(t, test.expectedError, apiError, "Error in test case %v", n)
			}
		}
	}
}

func TestWorkerHandler_HandleAddProxyResource(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		// API method args
		org     string
		request *CreateProxyResourceRequest
		// Expected result
		expectedStatusCode int
		expectedResponse   api.ProxyResource
		expectedError      api.Error
		// Manager Results
		createProxyResource *api.ProxyResource
		// Manager Errors
		createProxyResourceErr error
	}{
		"OkCase": {
			org: "org",
			request: &CreateProxyResourceRequest{
				Name: "name1",
				Path: "/path/",
				Resource: api.ResourceEntity{
					Host:   "Host",
					Path:   "/path",
					Method: "GET",
					Urn:    "urn",
					Action: "action",
				},
			},
			createProxyResource: &api.ProxyResource{
				ID:   "ID1",
				Name: "name1",
				Org:  "org",
				Path: "/path/",
				Resource: api.ResourceEntity{
					Host:   "Host",
					Path:   "/path",
					Method: "GET",
					Urn:    "urn",
					Action: "action",
				},
				CreateAt: now,
				UpdateAt: now,
				Urn:      api.CreateUrn("org", api.RESOURCE_PROXY, "/path/", "name1"),
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse: api.ProxyResource{
				ID:   "ID1",
				Name: "name1",
				Org:  "org",
				Path: "/path/",
				Resource: api.ResourceEntity{
					Host:   "Host",
					Path:   "/path",
					Method: "GET",
					Urn:    "urn",
					Action: "action",
				},
				CreateAt: now,
				UpdateAt: now,
				Urn:      api.CreateUrn("org", api.RESOURCE_PROXY, "/path/", "name1"),
			},
		},
		"ErrorCaseMalformedRequest": {
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "EOF",
			},
		},
		"ErrorCaseProxyResourceAlreadyExists": {
			org: "org",
			request: &CreateProxyResourceRequest{
				Name: "name1",
				Path: "/path/",
				Resource: api.ResourceEntity{
					Host:   "Host",
					Path:   "/path",
					Method: "GET",
					Urn:    "urn",
					Action: "action",
				},
			},
			createProxyResourceErr: &api.Error{
				Code: api.PROXY_RESOURCE_ALREADY_EXIST,
			},
			expectedStatusCode: http.StatusConflict,
			expectedError: api.Error{
				Code: api.PROXY_RESOURCE_ALREADY_EXIST,
			},
		},
		"ErrorCaseInvalidParameter": {
			org: "org",
			request: &CreateProxyResourceRequest{
				Name: "name1",
				Path: "/path/",
				Resource: api.ResourceEntity{
					Host:   "Host",
					Path:   "/path",
					Method: "GET",
					Urn:    "urn",
					Action: "action",
				},
			},
			createProxyResourceErr: &api.Error{
				Code: api.INVALID_PARAMETER_ERROR,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code: api.INVALID_PARAMETER_ERROR,
			},
		},
		"ErrorCaseUnauthorized": {
			org: "org",
			request: &CreateProxyResourceRequest{
				Name: "name1",
				Path: "/path/",
				Resource: api.ResourceEntity{
					Host:   "Host",
					Path:   "/path",
					Method: "GET",
					Urn:    "urn",
					Action: "action",
				},
			},
			createProxyResourceErr: &api.Error{
				Code: api.UNAUTHORIZED_RESOURCES_ERROR,
			},
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code: api.UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"ErrorCaseInternalServerError": {
			org: "org1",
			request: &CreateProxyResourceRequest{
				Name: "name1",
				Path: "/path/",
				Resource: api.ResourceEntity{
					Host:   "Host",
					Path:   "/path",
					Method: "GET",
					Urn:    "urn",
					Action: "action",
				},
			},
			createProxyResourceErr: &api.Error{
				Code: api.UNKNOWN_API_ERROR,
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedError: api.Error{
				Code: api.UNKNOWN_API_ERROR,
			},
		},
		"ErrorCaseProxyResourceRouteConflcit": {
			org: "org1",
			request: &CreateProxyResourceRequest{
				Name: "name1",
				Path: "/path/",
				Resource: api.ResourceEntity{
					Host:   "Host",
					Path:   "/path",
					Method: "GET",
					Urn:    "urn",
					Action: "action",
				},
			},
			createProxyResourceErr: &api.Error{
				Code: api.PROXY_RESOURCES_ROUTES_CONFLICT,
			},
			expectedStatusCode: http.StatusConflict,
			expectedError: api.Error{
				Code: api.PROXY_RESOURCES_ROUTES_CONFLICT,
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {
		testApi.ArgsOut[AddProxyResourceMethod][0] = test.createProxyResource
		testApi.ArgsOut[AddProxyResourceMethod][1] = test.createProxyResourceErr

		var body *bytes.Buffer
		if test.request != nil {
			jsonObject, err := json.Marshal(test.request)
			assert.Nil(t, err, "Error in test case %v", n)
			body = bytes.NewBuffer(jsonObject)
		}
		if body == nil {
			body = bytes.NewBuffer([]byte{})
		}

		path := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/proxy-resources", test.org)
		req, err := http.NewRequest(http.MethodPost, path, body)
		assert.Nil(t, err, "Error in test case %v", n)

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if test.request != nil {
			// Check received parameters

			assert.Equal(t, test.request.Name, testApi.ArgsIn[AddProxyResourceMethod][1], "Error in test case %v", n)
			assert.Equal(t, test.request.Path, testApi.ArgsIn[AddProxyResourceMethod][2], "Error in test case %v", n)
			assert.Equal(t, test.org, testApi.ArgsIn[AddProxyResourceMethod][3], "Error in test case %v", n)
			assert.Equal(t, test.request.Resource, testApi.ArgsIn[AddProxyResourceMethod][4], "Error in test case %v", n)
		}

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusCreated:
			response := api.ProxyResource{}
			err = json.NewDecoder(res.Body).Decode(&response)
			assert.Nil(t, err, "Error in test case %v", n)
			// Check result
			assert.Equal(t, test.expectedResponse, response, "Error in test case %v", n)
		case http.StatusInternalServerError: // Empty message so continue
			continue
		default:
			apiError := api.Error{}
			err = json.NewDecoder(res.Body).Decode(&apiError)
			assert.Nil(t, err, "Error in test case %v", n)
			// Check error
			assert.Equal(t, test.expectedError, apiError, "Error in test case %v", n)
		}
	}
}

func TestWorkerHandler_HandleUpdateProxyResource(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		// API method args
		org     string
		request *UpdateProxyResourceRequest
		// Expected result
		expectedStatusCode int
		expectedResponse   api.ProxyResource
		expectedError      api.Error
		// Manager Results
		updateProxyResourceResult *api.ProxyResource
		// Manager Errors
		updateProxyResourceErr error
	}{
		"OKCase": {
			org: "org",
			request: &UpdateProxyResourceRequest{
				Name: "newName",
				Path: "/new/",
				Resource: api.ResourceEntity{
					Host:   "http://new.com",
					Path:   "/new",
					Method: "POST",
					Action: "new:get",
					Urn:    "urn:ews:example:new:resource/get",
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: api.ProxyResource{
				ID:   "ID1",
				Name: "newName",
				Org:  "org",
				Path: "/new/",
				Resource: api.ResourceEntity{
					Host:   "http://new.com",
					Path:   "/new",
					Method: "POST",
					Action: "new:get",
					Urn:    "urn:ews:example:new:resource/get",
				},
				CreateAt: now,
				UpdateAt: now,
				Urn:      api.CreateUrn("org", api.RESOURCE_PROXY, "/new/", "newName"),
			},
			updateProxyResourceResult: &api.ProxyResource{
				ID:   "ID1",
				Name: "newName",
				Org:  "org",
				Path: "/new/",
				Resource: api.ResourceEntity{
					Host:   "http://new.com",
					Path:   "/new",
					Method: "POST",
					Action: "new:get",
					Urn:    "urn:ews:example:new:resource/get",
				},
				CreateAt: now,
				UpdateAt: now,
				Urn:      api.CreateUrn("org", api.RESOURCE_PROXY, "/new/", "newName"),
			},
		},
		"ErrorCaseMalformedRequest": {
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "EOF",
			},
		},
		"ErrorCaseProxyResourceNotFound": {
			request: &UpdateProxyResourceRequest{
				Name: "newName",
				Path: "/NewPath/",
				Resource: api.ResourceEntity{
					Host:   "http://new.com",
					Path:   "/new",
					Method: "POST",
					Action: "new:get",
					Urn:    "urn:ews:example:new:resource/get",
				},
			},
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.PROXY_RESOURCE_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Group not found",
			},
			updateProxyResourceErr: &api.Error{
				Code:    api.PROXY_RESOURCE_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Group not found",
			},
		},
		"ErrorCaseInvalidParameterError": {
			request: &UpdateProxyResourceRequest{
				Name: "newName",
				Path: "InvalidPath",
				Resource: api.ResourceEntity{
					Host:   "http://new.com",
					Path:   "/new",
					Method: "POST",
					Action: "new:get",
					Urn:    "urn:ews:example:new:resource/get",
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
			updateProxyResourceErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
		},
		"ErrorCaseProxyResourceAlreadyExistError": {
			request: &UpdateProxyResourceRequest{
				Name: "newName",
				Path: "newPath",
				Resource: api.ResourceEntity{
					Host:   "http://new.com",
					Path:   "/new",
					Method: "POST",
					Action: "new:get",
					Urn:    "urn:ews:example:new:resource/get",
				},
			},
			expectedStatusCode: http.StatusConflict,
			expectedError: api.Error{
				Code:    api.PROXY_RESOURCE_ALREADY_EXIST,
				Message: "Proxy resource already exist",
			},
			updateProxyResourceErr: &api.Error{
				Code:    api.PROXY_RESOURCE_ALREADY_EXIST,
				Message: "Proxy resource already exist",
			},
		},
		"ErrorCaseUnauthorizedResourcesError": {
			request: &UpdateProxyResourceRequest{
				Name: "newName",
				Path: "/NewPath/",
				Resource: api.ResourceEntity{
					Host:   "http://new.com",
					Path:   "/new",
					Method: "POST",
					Action: "new:get",
					Urn:    "urn:ews:example:new:resource/get",
				},
			},
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			updateProxyResourceErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseUnknownApiError": {
			request: &UpdateProxyResourceRequest{
				Name: "newName",
				Path: "/NewPath/",
				Resource: api.ResourceEntity{
					Host:   "http://new.com",
					Path:   "/new",
					Method: "POST",
					Action: "new:get",
					Urn:    "urn:ews:example:new:resource/get",
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			updateProxyResourceErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
		"ErrorCaseProxyResourceRouteConflict": {
			request: &UpdateProxyResourceRequest{
				Name: "newName",
				Path: "/NewPath/",
				Resource: api.ResourceEntity{
					Host:   "http://new.com",
					Path:   "/new",
					Method: "POST",
					Action: "new:get",
					Urn:    "urn:ews:example:new:resource/get",
				},
			},
			expectedStatusCode: http.StatusConflict,
			updateProxyResourceErr: &api.Error{
				Code:    api.PROXY_RESOURCES_ROUTES_CONFLICT,
				Message: "Error",
			},
			expectedError: api.Error{
				Code:    api.PROXY_RESOURCES_ROUTES_CONFLICT,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[UpdateProxyResourceMethod][0] = test.updateProxyResourceResult
		testApi.ArgsOut[UpdateProxyResourceMethod][1] = test.updateProxyResourceErr

		var body *bytes.Buffer
		if test.request != nil {
			jsonObject, err := json.Marshal(test.request)
			assert.Nil(t, err, "Error in test case %v", n)
			body = bytes.NewBuffer(jsonObject)
		}
		if body == nil {
			body = bytes.NewBuffer([]byte{})
		}

		path := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/proxy-resources/pr", test.org)
		req, err := http.NewRequest(http.MethodPut, path, body)
		assert.Nil(t, err, "Error in test case %v", n)

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if test.request != nil {
			// Check received parameters
			assert.Equal(t, test.org, testApi.ArgsIn[UpdateProxyResourceMethod][1], "Error in test case %v", n)
			assert.Equal(t, "pr", testApi.ArgsIn[UpdateProxyResourceMethod][2], "Error in test case %v", n)
			assert.Equal(t, test.request.Name, testApi.ArgsIn[UpdateProxyResourceMethod][3], "Error in test case %v", n)
			assert.Equal(t, test.request.Path, testApi.ArgsIn[UpdateProxyResourceMethod][4], "Error in test case %v", n)
		}

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusOK:
			response := api.ProxyResource{}
			err = json.NewDecoder(res.Body).Decode(&response)
			assert.Nil(t, err, "Error in test case %v", n)
			// Check result
			assert.Equal(t, test.expectedResponse, response, "Error in test case %v", n)
		case http.StatusInternalServerError: // Empty message so continue
			continue
		default:
			apiError := api.Error{}
			err = json.NewDecoder(res.Body).Decode(&apiError)
			assert.Nil(t, err, "Error in test case %v", n)
			// Check error
			assert.Equal(t, test.expectedError, apiError, "Error in test case %v", n)
		}
	}
}

func TestWorkerHandler_HandleGetProxyResourceByName(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		// API method args
		org          string
		name         string
		offset       string
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedResponse   api.ProxyResource
		expectedError      api.Error
		// Manager Results
		getProxyResourceByNameResult *api.ProxyResource
		// Manager Errors
		getProxyResourceByNameErr error
	}{
		"OKCase": {
			org:                "org",
			name:               "pr",
			expectedStatusCode: http.StatusOK,
			expectedResponse: api.ProxyResource{
				ID:   "prID",
				Name: "pr",
				Path: "/path/",
				Org:  "org",
				Urn:  "urn",
				Resource: api.ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
					Urn:    "urn2",
				},
				CreateAt: now,
				UpdateAt: now,
			},
			getProxyResourceByNameResult: &api.ProxyResource{
				ID:   "prID",
				Name: "pr",
				Path: "/path/",
				Org:  "org",
				Urn:  "urn",
				Resource: api.ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
					Urn:    "urn2",
				},
				CreateAt: now,
				UpdateAt: now,
			},
		},
		"ErrorCaseInvalidRequest": {
			org:                "org",
			name:               "pr",
			offset:             "-1",
			ignoreArgsIn:       true,
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
			getProxyResourceByNameErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
		},
		"ErrorCaseProxyResourceNotFound": {
			org:                "org",
			name:               "pr",
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.PROXY_RESOURCE_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Proxy resource not found",
			},
			getProxyResourceByNameErr: &api.Error{
				Code:    api.PROXY_RESOURCE_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Proxy resource not found",
			},
		},
		"ErrorCaseUnauthorizedResourcesError": {
			name:               "pr",
			org:                "org",
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			getProxyResourceByNameErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseInvalidParameterError": {
			name:               "invalid",
			org:                "org",
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
			getProxyResourceByNameErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
		},
		"ErrorCaseUnknownApiError": {
			name:               "pr",
			org:                "org",
			expectedStatusCode: http.StatusInternalServerError,
			getProxyResourceByNameErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[GetProxyResourceByNameMethod][0] = test.getProxyResourceByNameResult
		testApi.ArgsOut[GetProxyResourceByNameMethod][1] = test.getProxyResourceByNameErr

		path := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/proxy-resources/%v", test.org, test.name)
		req, err := http.NewRequest(http.MethodGet, path, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		q := req.URL.Query()
		q.Add("Offset", test.offset)
		req.URL.RawQuery = q.Encode()

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameters
			assert.Equal(t, test.org, testApi.ArgsIn[GetProxyResourceByNameMethod][1], "Error in test case %v", n)
			assert.Equal(t, test.name, testApi.ArgsIn[GetProxyResourceByNameMethod][2], "Error in test case %v", n)
		}
		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusOK:
			response := api.ProxyResource{}
			err = json.NewDecoder(res.Body).Decode(&response)
			assert.Nil(t, err, "Error in test case %v", n)
			// Check result
			assert.Equal(t, test.expectedResponse, response, "Error in test case %v", n)
		case http.StatusInternalServerError: // Empty message so continue
			continue
		default:
			apiError := api.Error{}
			err = json.NewDecoder(res.Body).Decode(&apiError)
			assert.Nil(t, err, "Error in test case %v", n)
			// Check error
			assert.Equal(t, test.expectedError, apiError, "Error in test case %v", n)
		}
	}
}

func TestWorkerHandler_HandleListProxyResource(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		filter       *api.Filter
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedResponse   ListProxyResourcesResponse
		expectedError      api.Error
		// Manager Results
		getProxyResourceListResult []api.ProxyResourceIdentity
		totalProxyResourcesResult  int
		// Manager Errors
		getProxyResourceListErr error
	}{
		"OkCase": {
			filter: &api.Filter{
				PathPrefix: "/path/",
				Org:        "org1",
				Offset:     0,
				Limit:      0,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: ListProxyResourcesResponse{
				Resources: []string{"name"},
				Offset:    0,
				Limit:     0,
				Total:     1,
			},
			getProxyResourceListResult: []api.ProxyResourceIdentity{
				{
					Name: "name",
					Org:  "org1",
				},
			},
			totalProxyResourcesResult: 1,
		},
		"ErrorCaseInvalidFilterParams": {
			filter: &api.Filter{
				PathPrefix: "",
				Limit:      -1,
			},
			ignoreArgsIn:       true,
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Limit -1",
			},
		},
		"OkCaseNoOrg": {
			filter: &api.Filter{
				PathPrefix: "/path/",
				Offset:     0,
				Limit:      0,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: ListProxyResourcesResponse{
				Resources: []string{"name1", "name2"},
				Offset:    0,
				Limit:     0,
				Total:     2,
			},
			getProxyResourceListResult: []api.ProxyResourceIdentity{
				{
					Org:  "org1",
					Name: "name1",
				},
				{
					Org:  "org2",
					Name: "name2",
				},
			},
			totalProxyResourcesResult: 2,
		},
		"ErrorCaseInvalidParameterError": {
			filter: &api.Filter{
				PathPrefix: "/path/",
				Org:        "org1",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Unauthorized",
			},
			getProxyResourceListErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseUnauthorizedError": {
			filter: &api.Filter{
				PathPrefix: "/path/",
				Org:        "org1",
			},
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			getProxyResourceListErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseUnknownApiError": {
			filter:             testFilter,
			expectedStatusCode: http.StatusInternalServerError,
			getProxyResourceListErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[ListProxyResourcesMethod][0] = test.getProxyResourceListResult
		testApi.ArgsOut[ListProxyResourcesMethod][1] = test.totalProxyResourcesResult
		testApi.ArgsOut[ListProxyResourcesMethod][2] = test.getProxyResourceListErr

		path := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/proxy-resources", test.filter.Org)
		req, err := http.NewRequest(http.MethodGet, path, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		addQueryParams(test.filter, req)

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameters
			filterData, ok := testApi.ArgsIn[ListProxyResourcesMethod][1].(*api.Filter)
			if ok {
				// Check result
				assert.Equal(t, test.filter, filterData, "Error in test case %v", n)
			}
		}

		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusOK:
			listProxyResources := ListProxyResourcesResponse{}
			err = json.NewDecoder(res.Body).Decode(&listProxyResources)
			assert.Nil(t, err, "Error in test case %v", n)
			// Check result
			assert.Equal(t, test.expectedResponse, listProxyResources, "Error in test case %v", n)
		case http.StatusInternalServerError: // Empty message so continue
			continue
		default:
			apiError := api.Error{}
			err = json.NewDecoder(res.Body).Decode(&apiError)
			assert.Nil(t, err, "Error in test case %v", n)
			// Check error
			assert.Equal(t, test.expectedError, apiError, "Error in test case %v", n)
		}
	}
}

func TestWorkerHandler_HandleRemoveProxyResource(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		org          string
		name         string
		offset       string
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedError      api.Error
		// Manager Errors
		removeProxyResource error
	}{
		"OKCase": {
			org:                "org",
			name:               "pr",
			expectedStatusCode: http.StatusNoContent,
		},
		"ErrorCaseInvalidRequest": {
			org:                "org",
			name:               "pr",
			offset:             "-1",
			ignoreArgsIn:       true,
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
			removeProxyResource: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
		},
		"ErrorCaseGroupNotFound": {
			org:                "org",
			name:               "pr",
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.PROXY_RESOURCE_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Group not found",
			},
			removeProxyResource: &api.Error{
				Code:    api.PROXY_RESOURCE_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Group not found",
			},
		},
		"ErrorCaseInvalidParameterError": {
			org:                "org",
			name:               "InvalidID",
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
			removeProxyResource: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
		},
		"ErrorCaseUnauthorizedResourcesError": {
			org:                "org",
			name:               "UnauthorizedID",
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			removeProxyResource: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseUnknownApiError": {
			org:                "org",
			name:               "ExceptionID",
			expectedStatusCode: http.StatusInternalServerError,
			removeProxyResource: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {
		testApi.ArgsOut[RemoveProxyResourceMethod][0] = test.removeProxyResource

		path := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/proxy-resources/%v", test.org, test.name)
		req, err := http.NewRequest(http.MethodDelete, path, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		q := req.URL.Query()
		q.Add("Offset", test.offset)
		req.URL.RawQuery = q.Encode()

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameters
			assert.Equal(t, test.org, testApi.ArgsIn[RemoveProxyResourceMethod][1], "Error in test case %v", n)
			assert.Equal(t, test.name, testApi.ArgsIn[RemoveProxyResourceMethod][2], "Error in test case %v", n)
		}

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusNoContent:
			// No message expected
			continue
		case http.StatusInternalServerError: // Empty message so continue
			continue
		default:
			apiError := api.Error{}
			err = json.NewDecoder(res.Body).Decode(&apiError)
			assert.Nil(t, err, "Error in test case %v", n)
			// Check error
			assert.Equal(t, test.expectedError, apiError, "Error in test case %v", n)
		}
	}
}
