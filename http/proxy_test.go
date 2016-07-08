package http

import (
	"encoding/json"
	"net/http"
	"testing"

	"reflect"

	"github.com/kylelemons/godebug/pretty"
	"github.com/tecsisa/authorizr/api"
	"time"
)

func TestProxyHandler_HandleRequest(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		expectedStatusCode int
		expectedError      api.Error
		expectedResponse   GetUserByIdResponse
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
		unauthenticated bool
	}{
		"OkCaseAdmin": {
			expectedStatusCode: http.StatusOK,
			resource:           USER_ROOT_URL + "/user",
			expectedResponse: GetUserByIdResponse{
				User: &api.User{
					ID:         "UserID",
					ExternalID: "ExternalID",
					Path:       "Path",
					Urn:        "urn",
					CreateAt:   now,
				},
			},
			getUserByExternalIdResult: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
			},
			getAuthorizedExternalResourcesResult: []string{"urn:ews:example:instance1:resource/user"},
		},
		"ErrorCaseInvalidParameter": {
			expectedStatusCode: http.StatusForbidden,
			resource:           "/urnPrefix",
			expectedError: api.Error{
				Code:    AUTHORIZATION_ERROR,
				Message: "Forbidden resource. If you need access, contact the administrators.",
			},
		},
		"ErrorCaseInvalidResources": {
			expectedStatusCode: http.StatusForbidden,
			resource:           "/invalidUrn",
			expectedError: api.Error{
				Code:    AUTHORIZATION_ERROR,
				Message: "Forbidden resource. If you need access, contact the administrators.",
			},
		},
		"ErrorCaseInvalidAction": {
			expectedStatusCode: http.StatusForbidden,
			resource:           "/invalidAction",
			expectedError: api.Error{
				Code:    AUTHORIZATION_ERROR,
				Message: "Forbidden resource. If you need access, contact the administrators.",
			},
		},
		"ErrorCaseInvalidHost": {
			expectedStatusCode: http.StatusForbidden,
			resource:           "/invalid",
			expectedError: api.Error{
				Code:    INVALID_DEST_HOST_URL,
				Message: "Forbidden resource. If you need access, contact the administrators.",
			},
			getAuthorizedExternalResourcesResult: []string{"urn:ews:example:instance1:resource/invalid"},
		},
		"ErrorCaseInvalidResource": {
			expectedStatusCode: http.StatusForbidden,
			resource:           "/fail",
			expectedError: api.Error{
				Code:    DESTINATION_HOST_RESOURCE_CALL_ERROR,
				Message: "Forbidden resource. If you need access, contact the administrators.",
			},
			getAuthorizedExternalResourcesResult: []string{"urn:ews:example:instance1:resource/fail"},
		},
		"ErrorCaseWorkedsa": {
			expectedStatusCode: http.StatusForbidden,
			resource:           USER_ROOT_URL + "/user",
			expectedError: api.Error{
				Code:    AUTHORIZATION_ERROR,
				Message: "Forbidden resource. If you need access, contact the administrators.",
			},
			getAuthorizedExternalResourcesResult: []string{"urn:ews:example:instance1:resource/user"},
			unauthenticated:                      true,
		},
		"ErrorCaseWorkerError": {
			expectedStatusCode: http.StatusForbidden,
			resource:           USER_ROOT_URL + "/user",
			expectedError: api.Error{
				Code:    AUTHORIZATION_ERROR,
				Message: "Forbidden resource. If you need access, contact the administrators.",
			},
			getAuthorizedExternalResourcesErr: &api.Error{
				Code: api.UNKNOWN_API_ERROR,
			},
		},
		"ErrorCaseNotAllowed": {
			expectedStatusCode: http.StatusForbidden,
			resource:           USER_ROOT_URL + "/user",
			expectedError: api.Error{
				Code:    AUTHORIZATION_ERROR,
				Message: "Forbidden resource. If you need access, contact the administrators.",
			},
		},
		"ErrorCaseUnauthorized": {
			expectedStatusCode: http.StatusForbidden,
			resource:           USER_ROOT_URL + "/user",
			expectedError: api.Error{
				Code:    AUTHORIZATION_ERROR,
				Message: "Forbidden resource. If you need access, contact the administrators.",
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
		if !test.unauthenticated {
			req.SetBasicAuth("admin", "admin")
		} else {
			authConnector.unauthenticated = test.unauthenticated
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
			getUserByIdResponse := GetUserByIdResponse{}
			err = json.NewDecoder(res.Body).Decode(&getUserByIdResponse)
			if err != nil {
				t.Fatalf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(getUserByIdResponse, test.expectedResponse); diff != "" {
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
			if !reflect.DeepEqual(apiError, test.expectedError) {
				t.Errorf("Test %v failed. Received different error response (wanted:%v / received:%v)",
					n, test.expectedError, apiError)
				continue
			}

		}

	}
}
