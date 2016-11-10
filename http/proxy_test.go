package http

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

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
