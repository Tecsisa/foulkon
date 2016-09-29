package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"fmt"

	"github.com/Tecsisa/foulkon/api"
	"github.com/kylelemons/godebug/pretty"
)

func TestWorkerHandler_HandleAddPolicy(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		// API method args
		org     string
		request *CreatePolicyRequest
		// Expected result
		expectedStatusCode int
		expectedResponse   *api.Policy
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
			expectedResponse: &api.Policy{
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
			if err != nil {
				t.Errorf("Test case %v. Unexpected marshalling api request %v", n, err)
				continue
			}
			body = bytes.NewBuffer(jsonObject)
		}
		if body == nil {
			body = bytes.NewBuffer([]byte{})
		}

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/policies", test.org)
		req, err := http.NewRequest(http.MethodPost, url, body)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error creating http request %v", n, err)
			continue
		}

		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error calling server %v", n, err)
			continue
		}

		if test.request != nil {
			// Check received parameters
			if testApi.ArgsIn[AddPolicyMethod][1] != test.request.Name {
				t.Errorf("Test case %v. Received different name (wanted:%v / received:%v)", n, test.request.Name, testApi.ArgsIn[AddPolicyMethod][1])
				continue
			}
			if testApi.ArgsIn[AddPolicyMethod][2] != test.request.Path {
				t.Errorf("Test case %v. Received different path (wanted:%v / received:%v)", n, test.request.Path, testApi.ArgsIn[AddPolicyMethod][2])
				continue
			}
			if testApi.ArgsIn[AddPolicyMethod][3] != test.org {
				t.Errorf("Test case %v. Received different org (wanted:%v / received:%v)", n, test.org, testApi.ArgsIn[AddPolicyMethod][3])
				continue
			}
			if diff := pretty.Compare(testApi.ArgsIn[AddPolicyMethod][4], test.request.Statements); diff != "" {
				t.Errorf("Test %v failed. Received different statements (received/wanted) %v",
					n, diff)
				continue
			}
		}

		// check status code
		if test.expectedStatusCode != res.StatusCode {
			t.Errorf("Test case %v. Received different http status code (wanted:%v / received:%v)", n, test.expectedStatusCode, res.StatusCode)
			continue
		}

		switch res.StatusCode {
		case http.StatusCreated:
			response := api.Policy{}
			err = json.NewDecoder(res.Body).Decode(&response)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(response, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
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

func TestWorkerHandler_HandleGetPolicyByName(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		// API method args
		org        string
		policyName string
		offset     string
		// Expected result
		expectedStatusCode int
		expectedResponse   *api.Policy
		expectedError      *api.Error
		// Manager Results
		getPolicyByNameResult *api.Policy
		// Manager Errors
		getPolicyByNameErr error
	}{
		"OkCase": {
			org:                "org1",
			policyName:         "p1",
			expectedStatusCode: http.StatusOK,
			expectedResponse: &api.Policy{
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
			expectedStatusCode: http.StatusBadRequest,
			getPolicyByNameErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
			expectedError: &api.Error{
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
			expectedError: &api.Error{
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
			expectedError: &api.Error{
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
			expectedError: &api.Error{
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
			expectedError: &api.Error{
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
		if err != nil {
			t.Errorf("Test case %v. Unexpected error creating http request %v", n, err)
			continue
		}

		q := req.URL.Query()
		q.Add("Offset", test.offset)
		req.URL.RawQuery = q.Encode()

		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error calling server %v", n, err)
			continue
		}

		// Check received parameters
		if testApi.ArgsIn[GetPolicyByNameMethod][1] != test.org {
			t.Errorf("Test case %v. Received different org (wanted:%v / received:%v)", n, test.org, testApi.ArgsIn[GetPolicyByNameMethod][1])
			continue
		}
		if testApi.ArgsIn[GetPolicyByNameMethod][2] != test.policyName {
			t.Errorf("Test case %v. Received different Name (wanted:%v / received:%v)", n, test.policyName, testApi.ArgsIn[GetPolicyByNameMethod][2])
			continue
		}

		// check status code
		if test.expectedStatusCode != res.StatusCode {
			t.Errorf("Test case %v. Received different http status code (wanted:%v / received:%v)", n, test.expectedStatusCode, res.StatusCode)
			continue
		}

		switch res.StatusCode {
		case http.StatusOK:
			response := api.Policy{}
			err = json.NewDecoder(res.Body).Decode(&response)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(response, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
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
		totalGroupsResult   int
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
			totalGroupsResult: 1,
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
			totalGroupsResult: 2,
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
		testApi.ArgsOut[ListPoliciesMethod][1] = test.totalGroupsResult
		testApi.ArgsOut[ListPoliciesMethod][2] = test.getPolicyListErr

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/policies", test.filter.Org)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error creating http request %v", n, err)
			continue
		}

		addQueryParams(test.filter, req)

		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error calling server %v", n, err)
			continue
		}

		if !test.ignoreArgsIn {
			// Check received parameters
			filterData, ok := testApi.ArgsIn[ListPoliciesMethod][1].(*api.Filter)
			if ok {
				// Check result
				if diff := pretty.Compare(filterData, test.filter); diff != "" {
					t.Errorf("Test %v failed. Received different filters (received/wanted) %v", n, diff)
					continue
				}
			}
		}

		if test.expectedStatusCode != res.StatusCode {
			t.Errorf("Test case %v. Received different http status code (wanted:%v / received:%v)", n, test.expectedStatusCode, res.StatusCode)
			continue
		}

		switch res.StatusCode {
		case http.StatusOK:
			listPoliciesResponse := ListPoliciesResponse{}
			err = json.NewDecoder(res.Body).Decode(&listPoliciesResponse)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(listPoliciesResponse, test.expectedResponse); diff != "" {
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
				t.Errorf("Test %v failed. Received different error response (received/wanted) %v",
					n, diff)
				continue
			}
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
		if err != nil {
			t.Errorf("Test case %v. Unexpected error creating http request %v", n, err)
			continue
		}

		addQueryParams(test.filter, req)

		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error calling server %v", n, err)
			continue
		}
		if !test.ignoreArgsIn {
			// Check received parameters
			filterData, ok := testApi.ArgsIn[ListPoliciesMethod][1].(*api.Filter)
			if ok {
				// Check result
				if diff := pretty.Compare(filterData, test.filter); diff != "" {
					t.Errorf("Test %v failed. Received different filters (received/wanted) %v", n, diff)
					continue
				}
			}
		}

		if test.expectedStatusCode != res.StatusCode {
			t.Errorf("Test case %v. Received different http status code (wanted:%v / received:%v)", n, test.expectedStatusCode, res.StatusCode)
			continue
		}

		switch res.StatusCode {
		case http.StatusOK:
			listAllPoliciesResponse := ListAllPoliciesResponse{}
			err = json.NewDecoder(res.Body).Decode(&listAllPoliciesResponse)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(listAllPoliciesResponse, test.expectedResponse); diff != "" {
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
				t.Errorf("Test %v failed. Received different error response (received/wanted) %v",
					n, diff)
				continue
			}
		}
	}
}

func TestWorkerHandler_HandleUpdatePolicy(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		// API method args
		org     string
		request *UpdatePolicyRequest
		// Expected result
		expectedStatusCode int
		expectedResponse   *api.Policy
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
			expectedResponse: &api.Policy{
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
			if err != nil {
				t.Errorf("Test case %v. Unexpected marshalling api request %v", n, err)
				continue
			}
			body = bytes.NewBuffer(jsonObject)
		}
		if body == nil {
			body = bytes.NewBuffer([]byte{})
		}

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/policies/policy1", test.org)
		req, err := http.NewRequest(http.MethodPut, url, body)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error creating http request %v", n, err)
			continue
		}

		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error calling server %v", n, err)
			continue
		}

		if test.request != nil {
			// Check received parameters
			if testApi.ArgsIn[UpdatePolicyMethod][1] != test.org {
				t.Errorf("Test case %v. Received different Org (wanted:%v / received:%v)", n, test.org, testApi.ArgsIn[UpdatePolicyMethod][1])
				continue
			}
			if testApi.ArgsIn[UpdatePolicyMethod][2] != "policy1" {
				t.Errorf("Test case %v. Received different Name (wanted:%v / received:%v)", n, "policy1", testApi.ArgsIn[UpdatePolicyMethod][2])
				continue
			}
			if testApi.ArgsIn[UpdatePolicyMethod][3] != test.request.Name {
				t.Errorf("Test case %v. Received different newName (wanted:%v / received:%v)", n, test.request.Name, testApi.ArgsIn[UpdatePolicyMethod][3])
				continue
			}
			if testApi.ArgsIn[UpdatePolicyMethod][4] != test.request.Path {
				t.Errorf("Test case %v. Received different Path (wanted:%v / received:%v)", n, test.request.Path, testApi.ArgsIn[UpdatePolicyMethod][4])
				continue
			}
			if diff := pretty.Compare(testApi.ArgsIn[UpdatePolicyMethod][5], test.request.Statements); diff != "" {
				t.Errorf("Test %v failed. Received different statements (received/wanted) %v",
					n, diff)
				continue
			}
		}

		// check status code
		if test.expectedStatusCode != res.StatusCode {
			t.Errorf("Test case %v. Received different http status code (wanted:%v / received:%v)", n, test.expectedStatusCode, res.StatusCode)
			continue
		}

		switch res.StatusCode {
		case http.StatusOK:
			response := api.Policy{}
			err = json.NewDecoder(res.Body).Decode(&response)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(response, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
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

func TestWorkerHandler_HandleRemovePolicy(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		org        string
		policyName string
		offset     string
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
			org:        "org1",
			policyName: "p1",
			offset:     "-1",
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
		if err != nil {
			t.Errorf("Test case %v. Unexpected error creating http request %v", n, err)
			continue
		}

		q := req.URL.Query()
		q.Add("Offset", test.offset)
		req.URL.RawQuery = q.Encode()

		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error calling server %v", n, err)
			continue
		}

		// Check received parameters
		if testApi.ArgsIn[RemovePolicyMethod][1] != test.org {
			t.Errorf("Test case %v. Received different Org (wanted:%v / received:%v)", n, test.org, testApi.ArgsIn[RemovePolicyMethod][1])
			continue
		}
		if testApi.ArgsIn[RemovePolicyMethod][2] != test.policyName {
			t.Errorf("Test case %v. Received different Name (wanted:%v / received:%v)", n, test.policyName, testApi.ArgsIn[RemovePolicyMethod][2])
			continue
		}

		// check status code
		if test.expectedStatusCode != res.StatusCode {
			t.Errorf("Test case %v. Received different http status code (wanted:%v / received:%v)", n, test.expectedStatusCode, res.StatusCode)
			continue
		}

		switch res.StatusCode {
		case http.StatusNoContent:
			// No message expected
			continue
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

func TestWorkerHandler_HandleListAttachedGroups(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		filter       *api.Filter
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedResponse   ListAttachedGroupsResponse
		expectedError      api.Error
		// Manager Results
		getPolicyGroupsResult []string
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
				Groups: []string{"group1", "group2"},
				Offset: 0,
				Limit:  0,
				Total:  2,
			},
			getPolicyGroupsResult: []string{"group1", "group2"},
			totalGroupsResult:     2,
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
		if err != nil {
			t.Errorf("Test case %v. Unexpected error creating http request %v", n, err)
			continue
		}

		addQueryParams(test.filter, req)

		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error calling server %v", n, err)
			continue
		}

		if !test.ignoreArgsIn {
			// Check received parameters
			filterData, ok := testApi.ArgsIn[ListAttachedGroupsMethod][1].(*api.Filter)
			if ok {
				// Check result
				if diff := pretty.Compare(filterData, test.filter); diff != "" {
					t.Errorf("Test %v failed. Received different filters (received/wanted) %v", n, diff)
					continue
				}
			}
		}

		// check status code
		if test.expectedStatusCode != res.StatusCode {
			t.Errorf("Test case %v. Received different http status code (wanted:%v / received:%v)", n, test.expectedStatusCode, res.StatusCode)
			continue
		}

		switch res.StatusCode {
		case http.StatusOK:
			GetAttachedGroupsMethodResponse := ListAttachedGroupsResponse{}
			err = json.NewDecoder(res.Body).Decode(&GetAttachedGroupsMethodResponse)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(GetAttachedGroupsMethodResponse, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
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
