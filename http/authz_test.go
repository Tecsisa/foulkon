package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/tecsisa/authorizr/api"
)

func TestWorkerHandler_HandleGetAuthorizedExternalResources(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		request *AuthorizeResourcesRequest
		// Expected result
		expectedStatusCode int
		expectedResponse   AuthorizeResourcesResponse
		expectedError      api.Error
		// Manager Results
		getAuthorizedExternalResourcesResult []string
		// Manager Errors
		getAuthorizedExternalResourcesErr error
	}{
		"OkCase": {
			request: &AuthorizeResourcesRequest{
				Resources: []string{},
				Action:    api.USER_ACTION_GET_USER,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: AuthorizeResourcesResponse{
				ResourcesAllowed: []string{"resource1", "resource2"},
			},
			getAuthorizedExternalResourcesResult: []string{"resource1", "resource2"},
		},
		"ErrorCaseMalformedRequest": {
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "EOF",
			},
		},
		"ErrorCaseInvalidParameter": {
			request: &AuthorizeResourcesRequest{
				Resources: []string{},
				Action:    api.USER_ACTION_GET_USER,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Error",
			},
			getAuthorizedExternalResourcesErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Error",
			},
		},
		"ErrorCaseUnauthorizedError": {
			request: &AuthorizeResourcesRequest{
				Resources: []string{},
				Action:    api.USER_ACTION_GET_USER,
			},
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Error",
			},
			getAuthorizedExternalResourcesErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Error",
			},
		},
		"ErrorCaseUnknownApiError": {
			request: &AuthorizeResourcesRequest{
				Resources: []string{},
				Action:    api.USER_ACTION_GET_USER,
			},
			expectedStatusCode: http.StatusInternalServerError,
			getAuthorizedExternalResourcesErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[GetAuthorizedExternalResourcesMethod][0] = test.getAuthorizedExternalResourcesResult
		testApi.ArgsOut[GetAuthorizedExternalResourcesMethod][1] = test.getAuthorizedExternalResourcesErr

		var body *bytes.Buffer
		if test.request != nil {
			jsonObject, err := json.Marshal(test.request)
			if err != nil {
				t.Errorf("Test case %v. Unexpected marshalling api request %v", n, err)
				continue
			}
			body = bytes.NewBuffer(jsonObject)
		}
		if body == nil {
			body = bytes.NewBuffer([]byte{})
		}
		req, err := http.NewRequest(http.MethodPost, server.URL+RESOURCE_URL, body)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error creating http request %v", n, err)
			continue
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
			authorizeResourcesResponse := AuthorizeResourcesResponse{}
			err = json.NewDecoder(res.Body).Decode(&authorizeResourcesResponse)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if !reflect.DeepEqual(authorizeResourcesResponse, test.expectedResponse) {
				t.Errorf("Test %v failed. Received different responses (wanted:%v / received:%v)",
					n, test.expectedResponse, authorizeResourcesResponse)
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
