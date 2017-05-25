package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Tecsisa/foulkon/api"
	"github.com/stretchr/testify/assert"
)

func TestWorkerHandler_HandleAddOidcProvider(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		// API method args
		request *CreateOidcProviderRequest
		// Expected result
		expectedStatusCode int
		expectedResponse   api.OidcProvider
		expectedError      api.Error
		// Manager Results
		addOidcProviderResult *api.OidcProvider
		// Manager Errors
		addOidcProviderErr error
	}{
		"OkCase": {
			request: &CreateOidcProviderRequest{
				Name: "test",
				Path: "/path/",
				OidcClients: []string{
					"client1",
				},
				IssuerURL: "https://test.com",
			},
			addOidcProviderResult: &api.OidcProvider{
				ID:        "test1",
				Name:      "test",
				Path:      "/path/",
				CreateAt:  now,
				UpdateAt:  now,
				Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "test"),
				IssuerURL: "https://test.com",
				OidcClients: []api.OidcClient{
					{
						Name: "client1",
					},
				},
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse: api.OidcProvider{
				ID:        "test1",
				Name:      "test",
				Path:      "/path/",
				CreateAt:  now,
				UpdateAt:  now,
				Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "test"),
				IssuerURL: "https://test.com",
				OidcClients: []api.OidcClient{
					{
						Name: "client1",
					},
				},
			},
		},
		"ErrorCaseMalformedRequest": {
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "EOF",
			},
		},
		"ErrorCaseOidcProviderAlreadyExists": {
			request: &CreateOidcProviderRequest{
				Name: "test",
				Path: "/path/",
				OidcClients: []string{
					"client1",
				},
				IssuerURL: "https://test.com",
			},
			addOidcProviderErr: &api.Error{
				Code: api.AUTH_OIDC_PROVIDER_ALREADY_EXIST,
			},
			expectedStatusCode: http.StatusConflict,
			expectedError: api.Error{
				Code: api.AUTH_OIDC_PROVIDER_ALREADY_EXIST,
			},
		},
		"ErrorCaseInvalidParameter": {
			request: &CreateOidcProviderRequest{
				Name: "test",
				Path: "/path/**",
				OidcClients: []string{
					"client1",
				},
				IssuerURL: "https://test.com",
			},
			addOidcProviderErr: &api.Error{
				Code: api.INVALID_PARAMETER_ERROR,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code: api.INVALID_PARAMETER_ERROR,
			},
		},
		"ErrorCaseUnauthorized": {
			request: &CreateOidcProviderRequest{
				Name: "test",
				Path: "/path/",
				OidcClients: []string{
					"client1",
				},
				IssuerURL: "https://test.com",
			},
			addOidcProviderErr: &api.Error{
				Code: api.UNAUTHORIZED_RESOURCES_ERROR,
			},
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code: api.UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"ErrorCaseInternalServerError": {
			request: &CreateOidcProviderRequest{
				Name: "test",
				Path: "/path/",
				OidcClients: []string{
					"client1",
				},
				IssuerURL: "https://test.com",
			},
			addOidcProviderErr: &api.Error{
				Code: api.UNKNOWN_API_ERROR,
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedError: api.Error{
				Code: api.UNKNOWN_API_ERROR,
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[AddOidcProviderMethod][0] = test.addOidcProviderResult
		testApi.ArgsOut[AddOidcProviderMethod][1] = test.addOidcProviderErr

		var body *bytes.Buffer
		if test.request != nil {
			jsonObject, err := json.Marshal(test.request)
			assert.Nil(t, err, "Error in test case %v", n)
			body = bytes.NewBuffer(jsonObject)
		}
		if body == nil {
			body = bytes.NewBuffer([]byte{})
		}

		url := fmt.Sprintf(server.URL + OIDC_AUTH_ROOT_URL)
		req, err := http.NewRequest(http.MethodPost, url, body)
		assert.Nil(t, err, "Error in test case %v", n)

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if test.request != nil {
			// Check received parameters
			assert.Equal(t, test.request.Name, testApi.ArgsIn[AddOidcProviderMethod][1], "Error in test case %v", n)
			assert.Equal(t, test.request.Path, testApi.ArgsIn[AddOidcProviderMethod][2], "Error in test case %v", n)
			assert.Equal(t, test.request.IssuerURL, testApi.ArgsIn[AddOidcProviderMethod][3], "Error in test case %v", n)
			assert.Equal(t, test.request.OidcClients, testApi.ArgsIn[AddOidcProviderMethod][4], "Error in test case %v", n)
		}

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusCreated:
			response := api.OidcProvider{}
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

func TestWorkerHandler_HandleGetOidcProviderByName(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		// API method args
		oidcProviderName string
		offset           string
		ignoreArgsIn     bool
		// Expected result
		expectedStatusCode int
		expectedResponse   api.OidcProvider
		expectedError      api.Error
		// Manager Results
		getOidcProviderByNameResult *api.OidcProvider
		// Manager Errors
		getOidcProviderByNameErr error
	}{
		"OkCase": {
			oidcProviderName:   "op1",
			expectedStatusCode: http.StatusOK,
			expectedResponse: api.OidcProvider{
				ID:        "test1",
				Name:      "test",
				Path:      "/path/",
				CreateAt:  now,
				UpdateAt:  now,
				Urn:       api.CreateUrn("org1", api.RESOURCE_POLICY, "/path/", "test"),
				IssuerURL: "https://test.com",
				OidcClients: []api.OidcClient{
					{
						Name: "client1",
					},
				},
			},
			getOidcProviderByNameResult: &api.OidcProvider{
				ID:        "test1",
				Name:      "test",
				Path:      "/path/",
				CreateAt:  now,
				UpdateAt:  now,
				Urn:       api.CreateUrn("org1", api.RESOURCE_POLICY, "/path/", "test"),
				IssuerURL: "https://test.com",
				OidcClients: []api.OidcClient{
					{
						Name: "client1",
					},
				},
			},
		},
		"ErrorCaseInvalidRequest": {
			oidcProviderName:   "op1",
			offset:             "-1",
			ignoreArgsIn:       true,
			expectedStatusCode: http.StatusBadRequest,
			getOidcProviderByNameErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
		},
		"ErrorCaseOidcProviderNotFound": {
			oidcProviderName:   "op1",
			expectedStatusCode: http.StatusNotFound,
			getOidcProviderByNameErr: &api.Error{
				Code: api.AUTH_OIDC_PROVIDER_BY_NAME_NOT_FOUND,
			},
			expectedError: api.Error{
				Code: api.AUTH_OIDC_PROVIDER_BY_NAME_NOT_FOUND,
			},
		},
		"ErrorCaseUnauthorized": {
			oidcProviderName:   "op1",
			expectedStatusCode: http.StatusForbidden,
			getOidcProviderByNameErr: &api.Error{
				Code: api.UNAUTHORIZED_RESOURCES_ERROR,
			},
			expectedError: api.Error{
				Code: api.UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"ErrorCaseInvalidParam": {
			oidcProviderName:   "op1",
			expectedStatusCode: http.StatusBadRequest,
			getOidcProviderByNameErr: &api.Error{
				Code: api.INVALID_PARAMETER_ERROR,
			},
			expectedError: api.Error{
				Code: api.INVALID_PARAMETER_ERROR,
			},
		},
		"ErrorCaseInternalServerError": {
			oidcProviderName:   "op1",
			expectedStatusCode: http.StatusInternalServerError,
			getOidcProviderByNameErr: &api.Error{
				Code: api.UNKNOWN_API_ERROR,
			},
			expectedError: api.Error{
				Code: api.UNKNOWN_API_ERROR,
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[GetOidcProviderByNameMethod][0] = test.getOidcProviderByNameResult
		testApi.ArgsOut[GetOidcProviderByNameMethod][1] = test.getOidcProviderByNameErr

		url := fmt.Sprintf(server.URL+OIDC_AUTH_ROOT_URL+"/%v", test.oidcProviderName)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		q := req.URL.Query()
		q.Add("Offset", test.offset)
		req.URL.RawQuery = q.Encode()

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameters
			assert.Equal(t, test.oidcProviderName, testApi.ArgsIn[GetOidcProviderByNameMethod][1], "Error in test case %v", n)
		}

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusOK:
			response := api.OidcProvider{}
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
