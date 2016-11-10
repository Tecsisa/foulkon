package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Tecsisa/foulkon/api"
	"github.com/stretchr/testify/assert"
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
			assert.Nil(t, err, "Error in test case %v", n)
			body = bytes.NewBuffer(jsonObject)
		}
		if body == nil {
			body = bytes.NewBuffer([]byte{})
		}
		req, err := http.NewRequest(http.MethodPost, server.URL+RESOURCE_URL, body)
		assert.Nil(t, err, "Error in test case %v", n)

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusOK:
			authorizeResourcesResponse := AuthorizeResourcesResponse{}
			err = json.NewDecoder(res.Body).Decode(&authorizeResourcesResponse)
			assert.Nil(t, err, "Error in test case %v", n)
			// Check result
			assert.Equal(t, test.expectedResponse, authorizeResourcesResponse, "Error in test case %v", n)
		case http.StatusInternalServerError: // Empty message so continue
			continue
		default:
			apiError := api.Error{}
			err = json.NewDecoder(res.Body).Decode(&apiError)
			assert.Nil(t, err, "Error in test case %v", n)
			// Check result
			assert.Equal(t, test.expectedError, apiError, "Error in test case %v", n)
		}
	}
}
