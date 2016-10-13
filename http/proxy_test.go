package http

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Tecsisa/foulkon/api"
	"github.com/kylelemons/godebug/pretty"
)

func TestProxyHandler_HandleRequest(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		expectedStatusCode int
		expectedError      api.Error
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
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Bad request",
			},
		},
		"ErrorCaseInvalidResources": {
			expectedStatusCode: http.StatusBadRequest,
			resource:           "/invalidUrn",
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Bad request",
			},
		},
		"ErrorCaseInvalidAction": {
			expectedStatusCode: http.StatusBadRequest,
			resource:           "/invalidAction",
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Bad request",
			},
		},
		"ErrorCaseInvalidHost": {
			expectedStatusCode: http.StatusInternalServerError,
			resource:           "/invalid",
			expectedError: api.Error{
				Code:    INVALID_DEST_HOST_URL,
				Message: "Error calling destination host",
			},
			getAuthorizedExternalResourcesResult: []string{"urn:ews:example:instance1:resource/invalid"},
		},
		"ErrorCaseHostUnreachable": {
			expectedStatusCode: http.StatusInternalServerError,
			resource:           "/fail",
			expectedError: api.Error{
				Code:    HOST_UNREACHABLE,
				Message: "Error calling destination resource",
			},
			getAuthorizedExternalResourcesResult: []string{"urn:ews:example:instance1:resource/fail"},
		},
		"ErrorCaseUnauthenticated": {
			expectedStatusCode: http.StatusForbidden,
			authStatusCode:     http.StatusUnauthorized,
			resource:           USER_ROOT_URL + "/user",
			expectedError: api.Error{
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
			expectedError: api.Error{
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
			expectedError: api.Error{
				Code:    FORBIDDEN_ERROR,
				Message: "Forbidden resource. If you need access, contact the administrator",
			},
		},
		"ErrorCaseAuthBadRequest": {
			expectedStatusCode: http.StatusBadRequest,
			authStatusCode:     http.StatusBadRequest,
			resource:           USER_ROOT_URL + "/user",
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Bad request",
			},
		},
		"ErrorCaseWorkerError": {
			expectedStatusCode: http.StatusInternalServerError,
			resource:           USER_ROOT_URL + "/user",
			expectedError: api.Error{
				Code:    INTERNAL_SERVER_ERROR,
				Message: "There was a problem retrieving authorization, status code 500",
			},
			getAuthorizedExternalResourcesErr: &api.Error{
				Code: api.UNKNOWN_API_ERROR,
			},
		},
		"ErrorCaseUnauthorized": {
			expectedStatusCode: http.StatusForbidden,
			resource:           USER_ROOT_URL + "/user",
			expectedError: api.Error{
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

		if err != nil {
			t.Errorf("Test case %v. Unexpected error creating http request %v", n, err)
			continue
		}

		if test.authStatusCode != 0 {
			authConnector.statusCode = test.authStatusCode
		}

		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error calling server %v", n, err)
			continue
		}

		// check status code
		if test.expectedStatusCode != res.StatusCode {
			t.Errorf("Test case %v. Received different http status code (wanted:%v / received:%v)", n, test.expectedStatusCode, res.StatusCode)
			continue
		}
		switch res.StatusCode {
		case http.StatusOK:
			response := api.User{}
			err = json.NewDecoder(res.Body).Decode(&response)
			if err != nil {
				t.Fatalf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(response, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v",
					n, diff)
				continue
			}
		case http.StatusInternalServerError: // Empty message so continue
			continue
		default:
			apiError := api.Error{}
			err = json.NewDecoder(res.Body).Decode(&apiError)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing error response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(apiError, test.expectedError); diff != "" {
				t.Errorf("Test %v failed. Received different error response (received/wanted) %v", n, diff)
				continue
			}
		}
	}
}
