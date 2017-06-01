package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"fmt"

	"github.com/Tecsisa/foulkon/api"
	"github.com/stretchr/testify/assert"
)

func TestWorkerHandler_HandleAddPolicy(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// API method args
		org     string
		request *CreatePolicyRequest
		// Expected result
		expectedStatusCode int
		expectedResponse   api.Policy
		expectedError      api.Error
		// Manager Results
		createPolicyResult *api.Policy
		// Manager Errors
		createPolicyErr error
	}{
		"OkCase": {
			org: "org1",
			request: &CreatePolicyRequest{
				Name: "test",
				Path: "/path/",
				Statements: []api.Statement{
					{
						Effect: "allow",
						Actions: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			createPolicyResult: &api.Policy{
				ID:       "test1",
				Name:     "test",
				Org:      "org1",
				Path:     "/path/",
				CreateAt: now,
				UpdateAt: now,
				Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]api.Statement{
					{
						Effect: "allow",
						Actions: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse: api.Policy{
				ID:       "test1",
				Name:     "test",
				Org:      "org1",
				CreateAt: now,
				UpdateAt: now,
				Path:     "/path/",
				Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]api.Statement{
					{
						Effect: "allow",
						Actions: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
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
		"ErrorCasePolicyAlreadyExists": {
			org: "org1",
			request: &CreatePolicyRequest{
				Name: "test",
				Path: "/path/",
				Statements: []api.Statement{
					{
						Effect: "allow",
						Actions: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			createPolicyErr: &api.Error{
				Code: api.POLICY_ALREADY_EXIST,
			},
			expectedStatusCode: http.StatusConflict,
			expectedError: api.Error{
				Code: api.POLICY_ALREADY_EXIST,
			},
		},
		"ErrorCaseInvalidParameter": {
			org: "org1",
			request: &CreatePolicyRequest{
				Name: "test",
				Path: "/path/**",
				Statements: []api.Statement{
					{
						Effect: "allow",
						Actions: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			createPolicyErr: &api.Error{
				Code: api.INVALID_PARAMETER_ERROR,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code: api.INVALID_PARAMETER_ERROR,
			},
		},
		"ErrorCaseUnauthorized": {
			org: "org1",
			request: &CreatePolicyRequest{
				Name: "test",
				Path: "/path/",
				Statements: []api.Statement{
					{
						Effect: "allow",
						Actions: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			createPolicyErr: &api.Error{
				Code: api.UNAUTHORIZED_RESOURCES_ERROR,
			},
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code: api.UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"ErrorCaseInternalServerError": {
			org: "org1",
			request: &CreatePolicyRequest{
				Name: "test",
				Path: "/path/**",
				Statements: []api.Statement{
					{
						Effect: "allow",
						Actions: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			createPolicyErr: &api.Error{
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

		testApi.ArgsOut[AddPolicyMethod][0] = test.createPolicyResult
		testApi.ArgsOut[AddPolicyMethod][1] = test.createPolicyErr

		var body *bytes.Buffer
		if test.request != nil {
			jsonObject, err := json.Marshal(test.request)
			assert.Nil(t, err, "Error in test case %v", n)
			body = bytes.NewBuffer(jsonObject)
		}
		if body == nil {
			body = bytes.NewBuffer([]byte{})
		}

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/policies", test.org)
		req, err := http.NewRequest(http.MethodPost, url, body)
		assert.Nil(t, err, "Error in test case %v", n)

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if test.request != nil {
			// Check received parameters
			assert.Equal(t, test.request.Name, testApi.ArgsIn[AddPolicyMethod][1], "Error in test case %v", n)
			assert.Equal(t, test.request.Path, testApi.ArgsIn[AddPolicyMethod][2], "Error in test case %v", n)
			assert.Equal(t, test.org, testApi.ArgsIn[AddPolicyMethod][3], "Error in test case %v", n)
			assert.Equal(t, test.request.Statements, testApi.ArgsIn[AddPolicyMethod][4], "Error in test case %v", n)
		}

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusCreated:
			response := api.Policy{}
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

func TestWorkerHandler_HandleGetPolicyByName(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// API method args
		org          string
		policyName   string
		offset       string
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedResponse   api.Policy
		expectedError      api.Error
		// Manager Results
		getPolicyByNameResult *api.Policy
		// Manager Errors
		getPolicyByNameErr error
	}{
		"OkCase": {
			org:                "org1",
			policyName:         "p1",
			expectedStatusCode: http.StatusOK,
			expectedResponse: api.Policy{
				ID:       "test1",
				Name:     "test",
				Org:      "org1",
				Path:     "/path/",
				CreateAt: now,
				UpdateAt: now,
				Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]api.Statement{
					{
						Effect: "allow",
						Actions: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			getPolicyByNameResult: &api.Policy{
				ID:       "test1",
				Name:     "test",
				Org:      "org1",
				Path:     "/path/",
				CreateAt: now,
				UpdateAt: now,
				Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]api.Statement{
					{
						Effect: "allow",
						Actions: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
		},
		"ErrorCaseInvalidRequest": {
			org:                "org1",
			policyName:         "p1",
			offset:             "-1",
			ignoreArgsIn:       true,
			expectedStatusCode: http.StatusBadRequest,
			getPolicyByNameErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
		},
		"ErrorCasePolicyNotFound": {
			org:                "org1",
			policyName:         "p1",
			expectedStatusCode: http.StatusNotFound,
			getPolicyByNameErr: &api.Error{
				Code: api.POLICY_BY_ORG_AND_NAME_NOT_FOUND,
			},
			expectedError: api.Error{
				Code: api.POLICY_BY_ORG_AND_NAME_NOT_FOUND,
			},
		},
		"ErrorCaseUnauthorized": {
			org:                "org1",
			policyName:         "p1",
			expectedStatusCode: http.StatusForbidden,
			getPolicyByNameErr: &api.Error{
				Code: api.UNAUTHORIZED_RESOURCES_ERROR,
			},
			expectedError: api.Error{
				Code: api.UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"ErrorCaseInvalidParam": {
			org:                "org1",
			policyName:         "p1",
			expectedStatusCode: http.StatusBadRequest,
			getPolicyByNameErr: &api.Error{
				Code: api.INVALID_PARAMETER_ERROR,
			},
			expectedError: api.Error{
				Code: api.INVALID_PARAMETER_ERROR,
			},
		},
		"ErrorCaseInternalServerError": {
			org:                "org1",
			policyName:         "p1",
			expectedStatusCode: http.StatusInternalServerError,
			getPolicyByNameErr: &api.Error{
				Code: api.UNKNOWN_API_ERROR,
			},
			expectedError: api.Error{
				Code: api.UNKNOWN_API_ERROR,
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[GetPolicyByNameMethod][0] = test.getPolicyByNameResult
		testApi.ArgsOut[GetPolicyByNameMethod][1] = test.getPolicyByNameErr

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/policies/%v", test.org, test.policyName)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		q := req.URL.Query()
		q.Add("Offset", test.offset)
		req.URL.RawQuery = q.Encode()

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameters
			assert.Equal(t, test.org, testApi.ArgsIn[GetPolicyByNameMethod][1], "Error in test case %v", n)
			assert.Equal(t, test.policyName, testApi.ArgsIn[GetPolicyByNameMethod][2], "Error in test case %v", n)
		}

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusOK:
			response := api.Policy{}
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

func TestWorkerHandler_HandleListPolicies(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		filter       *api.Filter
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedResponse   ListPoliciesResponse
		expectedError      api.Error
		// Manager Results
		getPolicyListResult []api.PolicyIdentity
		totalPoliciesResult int
		// Manager Errors
		getPolicyListErr error
	}{
		"OkCase": {
			filter: &api.Filter{
				PathPrefix: "/path/",
				Org:        "org1",
				Offset:     0,
				Limit:      0,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: ListPoliciesResponse{
				Policies: []string{"policy1"},
				Offset:   0,
				Limit:    0,
				Total:    1,
			},
			getPolicyListResult: []api.PolicyIdentity{
				{
					Org:  "org1",
					Name: "policy1",
				},
			},
			totalPoliciesResult: 1,
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
			expectedResponse: ListPoliciesResponse{
				Policies: []string{"policy1", "policy2"},
				Offset:   0,
				Limit:    0,
				Total:    2,
			},
			getPolicyListResult: []api.PolicyIdentity{
				{
					Org:  "org1",
					Name: "policy1",
				},
				{
					Org:  "org2",
					Name: "policy2",
				},
			},
			totalPoliciesResult: 2,
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
			getPolicyListErr: &api.Error{
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
			getPolicyListErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseUnknownApiError": {
			filter:             testFilter,
			expectedStatusCode: http.StatusInternalServerError,
			getPolicyListErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[ListPoliciesMethod][0] = test.getPolicyListResult
		testApi.ArgsOut[ListPoliciesMethod][1] = test.totalPoliciesResult
		testApi.ArgsOut[ListPoliciesMethod][2] = test.getPolicyListErr

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/policies", test.filter.Org)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		addQueryParams(test.filter, req)

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameters
			filterData, ok := testApi.ArgsIn[ListPoliciesMethod][1].(*api.Filter)
			if ok {
				// Check result
				assert.Equal(t, test.filter, filterData, "Error in test case %v", n)
			}
		}

		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusOK:
			listPoliciesResponse := ListPoliciesResponse{}
			err = json.NewDecoder(res.Body).Decode(&listPoliciesResponse)
			assert.Nil(t, err, "Error in test case %v", n)
			// Check result
			assert.Equal(t, test.expectedResponse, listPoliciesResponse, "Error in test case %v", n)
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

func TestWorkerHandler_HandleListAllPolicies(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		filter       *api.Filter
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedResponse   ListAllPoliciesResponse
		expectedError      api.Error
		// Manager Results
		getPolicyListResult []api.PolicyIdentity
		totalGroupsResult   int
		// Manager Errors
		getPolicyListErr error
	}{
		"OkCase": {
			filter: &api.Filter{
				PathPrefix: "path",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: ListAllPoliciesResponse{
				Policies: []api.PolicyIdentity{
					{
						Org:  "org1",
						Name: "policy1",
					},
					{
						Org:  "org1",
						Name: "policy2",
					},
				},
				Offset: 0,
				Limit:  0,
				Total:  2,
			},
			getPolicyListResult: []api.PolicyIdentity{
				{
					Org:  "org1",
					Name: "policy1",
				},
				{
					Org:  "org1",
					Name: "policy2",
				},
			},
			totalGroupsResult: 2,
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
				PathPrefix: "path",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Unauthorized",
			},
			getPolicyListErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseUnauthorizedError": {
			filter: &api.Filter{
				PathPrefix: "path",
			},
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			getPolicyListErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseUnknownApiError": {
			filter: &api.Filter{
				PathPrefix: "path",
			},
			expectedStatusCode: http.StatusInternalServerError,
			getPolicyListErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[ListPoliciesMethod][0] = test.getPolicyListResult
		testApi.ArgsOut[ListPoliciesMethod][1] = test.totalGroupsResult
		testApi.ArgsOut[ListPoliciesMethod][2] = test.getPolicyListErr

		url := fmt.Sprintf(server.URL + API_VERSION_1 + "/policies")
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		addQueryParams(test.filter, req)

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameters
			filterData, ok := testApi.ArgsIn[ListPoliciesMethod][1].(*api.Filter)
			if ok {
				// Check result
				assert.Equal(t, test.filter, filterData, "Error in test case %v", n)
			}
		}

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusOK:
			listAllPoliciesResponse := ListAllPoliciesResponse{}
			err = json.NewDecoder(res.Body).Decode(&listAllPoliciesResponse)
			assert.Nil(t, err, "Error in test case %v", n)
			// Check result
			assert.Equal(t, test.expectedResponse, listAllPoliciesResponse, "Error in test case %v", n)
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

func TestWorkerHandler_HandleUpdatePolicy(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// API method args
		org     string
		request *UpdatePolicyRequest
		// Expected result
		expectedStatusCode int
		expectedResponse   api.Policy
		expectedError      api.Error
		// Manager Results
		updatePolicyResult *api.Policy
		// Manager Errors
		updatePolicyErr error
	}{
		"OkCase": {
			org: "org1",
			request: &UpdatePolicyRequest{
				Name: "policy1",
				Path: "path1",
				Statements: []api.Statement{
					{
						Effect: "allow",
						Actions: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: api.Policy{
				ID:       "test1",
				Name:     "policy1",
				Path:     "/path/",
				Org:      "org1",
				CreateAt: now,
				UpdateAt: now,
				Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]api.Statement{
					{
						Effect: "allow",
						Actions: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			updatePolicyResult: &api.Policy{
				ID:       "test1",
				Name:     "policy1",
				Path:     "/path/",
				Org:      "org1",
				CreateAt: now,
				UpdateAt: now,
				Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]api.Statement{
					{
						Effect: "allow",
						Actions: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
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
		"ErrorCasePolicyNotFound": {
			org: "org1",
			request: &UpdatePolicyRequest{
				Name: "policy1",
				Path: "path1",
				Statements: []api.Statement{
					{
						Effect: "allow",
						Actions: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.POLICY_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Policy not found",
			},
			updatePolicyErr: &api.Error{
				Code:    api.POLICY_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Policy not found",
			},
		},
		"ErrorCasePolicyAlreadyExistError": {
			request: &UpdatePolicyRequest{
				Name: "policy1",
				Path: "path2",
			},
			expectedStatusCode: http.StatusConflict,
			expectedError: api.Error{
				Code:    api.POLICY_ALREADY_EXIST,
				Message: "Policy already exist",
			},
			updatePolicyErr: &api.Error{
				Code:    api.POLICY_ALREADY_EXIST,
				Message: "Policy already exist",
			},
		},
		"ErrorCaseInvalidParameterError": {
			org: "org1",
			request: &UpdatePolicyRequest{
				Name: "policy1",
				Path: "path1",
				Statements: []api.Statement{
					{
						Effect: "allow",
						Actions: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			updatePolicyErr: &api.Error{
				Code: api.INVALID_PARAMETER_ERROR,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code: api.INVALID_PARAMETER_ERROR,
			},
		},
		"ErrorCaseUnauthorizedResourcesError": {
			org: "org1",
			request: &UpdatePolicyRequest{
				Name: "policy1",
				Path: "path1",
				Statements: []api.Statement{
					{
						Effect: "allow",
						Actions: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			updatePolicyErr: &api.Error{
				Code: api.UNAUTHORIZED_RESOURCES_ERROR,
			},
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code: api.UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"ErrorCaseUnknownApiError": {
			org: "org1",
			request: &UpdatePolicyRequest{
				Name: "policy1",
				Path: "path1",
				Statements: []api.Statement{
					{
						Effect: "allow",
						Actions: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			updatePolicyErr: &api.Error{
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

		testApi.ArgsOut[UpdatePolicyMethod][0] = test.updatePolicyResult
		testApi.ArgsOut[UpdatePolicyMethod][1] = test.updatePolicyErr

		var body *bytes.Buffer
		if test.request != nil {
			jsonObject, err := json.Marshal(test.request)
			assert.Nil(t, err, "Error in test case %v", n)
			body = bytes.NewBuffer(jsonObject)
		}
		if body == nil {
			body = bytes.NewBuffer([]byte{})
		}

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/policies/policy1", test.org)
		req, err := http.NewRequest(http.MethodPut, url, body)
		assert.Nil(t, err, "Error in test case %v", n)

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if test.request != nil {
			// Check received parameters
			assert.Equal(t, test.org, testApi.ArgsIn[UpdatePolicyMethod][1], "Error in test case %v", n)
			assert.Equal(t, "policy1", testApi.ArgsIn[UpdatePolicyMethod][2], "Error in test case %v", n)
			assert.Equal(t, test.request.Name, testApi.ArgsIn[UpdatePolicyMethod][3], "Error in test case %v", n)
			assert.Equal(t, test.request.Path, testApi.ArgsIn[UpdatePolicyMethod][4], "Error in test case %v", n)
			assert.Equal(t, test.request.Statements, testApi.ArgsIn[UpdatePolicyMethod][5], "Error in test case %v", n)
		}

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusOK:
			response := api.Policy{}
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

func TestWorkerHandler_HandleRemovePolicy(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		org          string
		policyName   string
		offset       string
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedError      api.Error
		// Manager Errors
		deletePolicyErr error
	}{
		"OkCase": {
			org:                "org1",
			policyName:         "p1",
			expectedStatusCode: http.StatusNoContent,
		},
		"ErrorCaseInvalidRequest": {
			org:          "org1",
			policyName:   "p1",
			offset:       "-1",
			ignoreArgsIn: true,
			deletePolicyErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
		},
		"ErrorCasePolicyNotFound": {
			org:        "org1",
			policyName: "p1",
			deletePolicyErr: &api.Error{
				Code: api.POLICY_BY_ORG_AND_NAME_NOT_FOUND,
			},
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code: api.POLICY_BY_ORG_AND_NAME_NOT_FOUND,
			},
		},
		"ErrorCaseInvalidParam": {
			org:        "org1",
			policyName: "p1",
			deletePolicyErr: &api.Error{
				Code: api.INVALID_PARAMETER_ERROR,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code: api.INVALID_PARAMETER_ERROR,
			},
		},
		"ErrorCaseUnauthorized": {
			org:        "org1",
			policyName: "p1",
			deletePolicyErr: &api.Error{
				Code: api.UNAUTHORIZED_RESOURCES_ERROR,
			},
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code: api.UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"ErrorCaseInternalServerError": {
			org:        "org1",
			policyName: "p1",
			deletePolicyErr: &api.Error{
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

		testApi.ArgsOut[RemovePolicyMethod][0] = test.deletePolicyErr

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/policies/%v", test.org, test.policyName)
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		q := req.URL.Query()
		q.Add("Offset", test.offset)
		req.URL.RawQuery = q.Encode()

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameters
			assert.Equal(t, test.org, testApi.ArgsIn[RemovePolicyMethod][1], "Error in test case %v", n)
			assert.Equal(t, test.policyName, testApi.ArgsIn[RemovePolicyMethod][2], "Error in test case %v", n)
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

func TestWorkerHandler_HandleListAttachedGroups(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// API method args
		filter       *api.Filter
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedResponse   ListAttachedGroupsResponse
		expectedError      api.Error
		// Manager Results
		getPolicyGroupsResult []api.PolicyGroups
		totalGroupsResult     int
		// Manager Errors
		getPolicyGroupsErr error
	}{
		"OkCase": {
			filter: &api.Filter{
				Org:        "org1",
				PolicyName: "p1",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: ListAttachedGroupsResponse{
				Groups: []api.PolicyGroups{
					{
						Group:    "group1",
						CreateAt: now,
					},
					{
						Group:    "group2",
						CreateAt: now,
					},
				},
				Offset: 0,
				Limit:  0,
				Total:  2,
			},
			getPolicyGroupsResult: []api.PolicyGroups{
				{
					Group:    "group1",
					CreateAt: now,
				},
				{
					Group:    "group2",
					CreateAt: now,
				},
			},
			totalGroupsResult: 2,
		},
		"ErrorCaseInvalidFilterParams": {
			filter: &api.Filter{
				Limit: -1,
			},
			ignoreArgsIn:       true,
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Limit -1",
			},
		},
		"ErrorCaseNotFound": {
			filter: &api.Filter{
				Org:        "org1",
				PolicyName: "p1",
			},
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code: api.POLICY_BY_ORG_AND_NAME_NOT_FOUND,
			},
			getPolicyGroupsErr: &api.Error{
				Code: api.POLICY_BY_ORG_AND_NAME_NOT_FOUND,
			},
		},
		"ErrorCaseUnauthorized": {
			filter: &api.Filter{
				Org:        "org1",
				PolicyName: "p1",
			},
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code: api.UNAUTHORIZED_RESOURCES_ERROR,
			},
			getPolicyGroupsErr: &api.Error{
				Code: api.UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"ErrorCaseInvalidParam": {
			filter: &api.Filter{
				Org:        "org1",
				PolicyName: "p1",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code: api.INVALID_PARAMETER_ERROR,
			},
			getPolicyGroupsErr: &api.Error{
				Code: api.INVALID_PARAMETER_ERROR,
			},
		},
		"ErrorCaseInternalServerError": {
			filter: &api.Filter{
				Org:        "org1",
				PolicyName: "p1",
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedError: api.Error{
				Code: api.UNKNOWN_API_ERROR,
			},
			getPolicyGroupsErr: &api.Error{
				Code: api.UNKNOWN_API_ERROR,
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[ListAttachedGroupsMethod][0] = test.getPolicyGroupsResult
		testApi.ArgsOut[ListAttachedGroupsMethod][1] = test.totalGroupsResult
		testApi.ArgsOut[ListAttachedGroupsMethod][2] = test.getPolicyGroupsErr

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/policies/%v/groups", test.filter.Org, test.filter.PolicyName)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		addQueryParams(test.filter, req)

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameters
			filterData, ok := testApi.ArgsIn[ListAttachedGroupsMethod][1].(*api.Filter)
			if ok {
				// Check result
				assert.Equal(t, test.filter, filterData, "Error in test case %v", n)
			}
		}

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusOK:
			GetAttachedGroupsMethodResponse := ListAttachedGroupsResponse{}
			err = json.NewDecoder(res.Body).Decode(&GetAttachedGroupsMethodResponse)
			assert.Nil(t, err, "Error in test case %v", n)
			// Check result
			assert.Equal(t, test.expectedResponse, GetAttachedGroupsMethodResponse, "Error in test case %v", n)
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
