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
	now := time.Now().UTC()
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
	now := time.Now().UTC()
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
				Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "test"),
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
				Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "test"),
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

func TestWorkerHandler_HandleListOidcProviders(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		filter       *api.Filter
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedResponse   ListOidcProvidersResponse
		expectedError      api.Error
		// Manager Results
		listOidcProvidersResult []string
		listOidcProvidersTotal  int
		// Manager Errors
		listOidcProvidersErr error
	}{
		"OkCase": {
			filter: &api.Filter{
				PathPrefix: "/path/",
				Offset:     0,
				Limit:      0,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: ListOidcProvidersResponse{
				Providers: []string{"oidcProvider1"},
				Offset:    0,
				Limit:     0,
				Total:     1,
			},
			listOidcProvidersResult: []string{
				"oidcProvider1",
			},
			listOidcProvidersTotal: 1,
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
		"ErrorCaseInvalidParameterError": {
			filter: &api.Filter{
				PathPrefix: "/path/",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
			listOidcProvidersErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
		},
		"ErrorCaseUnauthorizedError": {
			filter: &api.Filter{
				PathPrefix: "/path/",
			},
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			listOidcProvidersErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseUnknownApiError": {
			filter:             testFilter,
			expectedStatusCode: http.StatusInternalServerError,
			listOidcProvidersErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[ListOidcProvidersMethod][0] = test.listOidcProvidersResult
		testApi.ArgsOut[ListOidcProvidersMethod][1] = test.listOidcProvidersTotal
		testApi.ArgsOut[ListOidcProvidersMethod][2] = test.listOidcProvidersErr

		url := fmt.Sprintf(server.URL + OIDC_AUTH_ROOT_URL)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		addQueryParams(test.filter, req)

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameters
			filterData, ok := testApi.ArgsIn[ListOidcProvidersMethod][1].(*api.Filter)
			if ok {
				// Check result
				assert.Equal(t, test.filter, filterData, "Error in test case %v", n)
			}
		}

		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusOK:
			listOidcProvidersResponse := ListOidcProvidersResponse{}
			err = json.NewDecoder(res.Body).Decode(&listOidcProvidersResponse)
			assert.Nil(t, err, "Error in test case %v", n)
			// Check result
			assert.Equal(t, test.expectedResponse, listOidcProvidersResponse, "Error in test case %v", n)
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

func TestWorkerHandler_HandleUpdateOidcProvider(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// API method args
		oidcProviderName string
		request          *UpdateOidcProviderRequest
		// Expected result
		expectedStatusCode int
		expectedResponse   api.OidcProvider
		expectedError      api.Error
		// Manager Results
		updateOidcProviderResult *api.OidcProvider
		// Manager Errors
		updateOidcProviderErr error
	}{
		"OkCase": {
			request: &UpdateOidcProviderRequest{
				Name:        "newName",
				Path:        "NewPath",
				IssuerURL:   "http://test.com",
				OidcClients: []string{"client1", "client2"},
			},
			oidcProviderName:   "oidcProviderName",
			expectedStatusCode: http.StatusOK,
			expectedResponse: api.OidcProvider{
				ID:        "ID",
				Name:      "newName",
				Path:      "NewPath",
				Urn:       "urn",
				CreateAt:  now,
				UpdateAt:  now,
				IssuerURL: "http://test.com",
				OidcClients: []api.OidcClient{
					{
						Name: "client1",
					},
					{
						Name: "client2",
					},
				},
			},
			updateOidcProviderResult: &api.OidcProvider{
				ID:        "ID",
				Name:      "newName",
				Path:      "NewPath",
				Urn:       "urn",
				IssuerURL: "http://test.com",
				CreateAt:  now,
				UpdateAt:  now,
				OidcClients: []api.OidcClient{
					{
						Name: "client1",
					},
					{
						Name: "client2",
					},
				},
			},
		},
		"ErrorCaseMalformedRequest": {
			oidcProviderName:   "oidcProviderName",
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "EOF",
			},
		},
		"ErrorCaseOidcProviderNotFound": {
			oidcProviderName: "oidcProviderName",
			request: &UpdateOidcProviderRequest{
				Name: "newName",
				Path: "NewPath",
			},
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.AUTH_OIDC_PROVIDER_BY_NAME_NOT_FOUND,
				Message: "OidcProvider not found",
			},
			updateOidcProviderErr: &api.Error{
				Code:    api.AUTH_OIDC_PROVIDER_BY_NAME_NOT_FOUND,
				Message: "OidcProvider not found",
			},
		},
		"ErrorCaseInvalidParameterError": {
			oidcProviderName: "oidcProviderName",
			request: &UpdateOidcProviderRequest{
				Name: "newName",
				Path: "InvalidPath",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
			updateOidcProviderErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
		},
		"ErrorCaseOidcProviderAlreadyExistError": {
			oidcProviderName: "oidcProviderName",
			request: &UpdateOidcProviderRequest{
				Name: "newName",
				Path: "newPath",
			},
			expectedStatusCode: http.StatusConflict,
			expectedError: api.Error{
				Code:    api.AUTH_OIDC_PROVIDER_ALREADY_EXIST,
				Message: "OidcProvider already exist",
			},
			updateOidcProviderErr: &api.Error{
				Code:    api.AUTH_OIDC_PROVIDER_ALREADY_EXIST,
				Message: "OidcProvider already exist",
			},
		},
		"ErrorCaseUnauthorizedResourcesError": {
			oidcProviderName: "oidcProviderName",
			request: &UpdateOidcProviderRequest{
				Name: "newName",
				Path: "NewPath",
			},
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			updateOidcProviderErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseUnknownApiError": {
			oidcProviderName: "oidcProviderName",
			request: &UpdateOidcProviderRequest{
				Name: "newName",
				Path: "NewPath",
			},
			expectedStatusCode: http.StatusInternalServerError,
			updateOidcProviderErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[UpdateOidcProviderMethod][0] = test.updateOidcProviderResult
		testApi.ArgsOut[UpdateOidcProviderMethod][1] = test.updateOidcProviderErr

		var body *bytes.Buffer
		if test.request != nil {
			jsonObject, err := json.Marshal(test.request)
			assert.Nil(t, err, "Error in test case %v", n)
			body = bytes.NewBuffer(jsonObject)
		}
		if body == nil {
			body = bytes.NewBuffer([]byte{})
		}

		url := fmt.Sprintf(server.URL+OIDC_AUTH_ROOT_URL+"/%v", test.oidcProviderName)
		req, err := http.NewRequest(http.MethodPut, url, body)
		assert.Nil(t, err, "Error in test case %v", n)

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if test.request != nil {
			// Check received parameters
			assert.Equal(t, test.oidcProviderName, testApi.ArgsIn[UpdateOidcProviderMethod][1], "Error in test case %v", n)
			assert.Equal(t, test.request.Name, testApi.ArgsIn[UpdateOidcProviderMethod][2], "Error in test case %v", n)
			assert.Equal(t, test.request.Path, testApi.ArgsIn[UpdateOidcProviderMethod][3], "Error in test case %v", n)
			assert.Equal(t, test.request.IssuerURL, testApi.ArgsIn[UpdateOidcProviderMethod][4], "Error in test case %v", n)
			assert.Equal(t, test.request.OidcClients, testApi.ArgsIn[UpdateOidcProviderMethod][5], "Error in test case %v", n)
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

func TestWorkerHandler_HandleRemoveOidcProvider(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		name         string
		offset       string
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedError      api.Error
		// Manager Errors
		removeOidcProviderErr error
	}{
		"OkCase": {
			name:               "oidcProvider1",
			expectedStatusCode: http.StatusNoContent,
		},
		"ErrorCaseInvalidRequest": {
			name:               "oidcProvider1",
			offset:             "-1",
			ignoreArgsIn:       true,
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
			removeOidcProviderErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
		},
		"ErrorCaseOidcProviderNotFound": {
			name:               "oidcProvider1",
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.AUTH_OIDC_PROVIDER_BY_NAME_NOT_FOUND,
				Message: "OIDC Provider not found",
			},
			removeOidcProviderErr: &api.Error{
				Code:    api.AUTH_OIDC_PROVIDER_BY_NAME_NOT_FOUND,
				Message: "OIDC Provider not found",
			},
		},
		"ErrorCaseInvalidParameterError": {
			name:               "InvalidID",
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
			removeOidcProviderErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
		},
		"ErrorCaseUnauthorizedResourcesError": {
			name:               "UnauthorizedID",
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			removeOidcProviderErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseUnknownApiError": {
			name:               "ExceptionID",
			expectedStatusCode: http.StatusInternalServerError,
			removeOidcProviderErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[RemoveOidcProviderMethod][0] = test.removeOidcProviderErr

		url := fmt.Sprintf(server.URL+OIDC_AUTH_ROOT_URL+"/%v", test.name)
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		q := req.URL.Query()
		q.Add("Offset", test.offset)
		req.URL.RawQuery = q.Encode()

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameters
			assert.Equal(t, test.name, testApi.ArgsIn[RemoveOidcProviderMethod][1], "Error in test case %v", n)
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
