package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"fmt"

	"github.com/kylelemons/godebug/pretty"
	"github.com/tecsisa/authorizr/api"
)

func TestWorkerHandler_HandlePolicyList(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		org        string
		pathPrefix string
		// Expected result
		expectedStatusCode int
		expectedResponse   ListPoliciesResponse
		expectedError      api.Error
		// API Results
		getPolicyListResult []api.PolicyIdentity
		// API Errors
		getPolicyListErr error
	}{
		"OkCase": {
			org:                "org1",
			pathPrefix:         "path",
			expectedStatusCode: http.StatusOK,
			expectedResponse: ListPoliciesResponse{
				[]string{"policy1"},
			},
			getPolicyListResult: []api.PolicyIdentity{
				{
					Org:  "org1",
					Name: "policy1",
				},
			},
		},
		"OkCaseNoOrg": {
			pathPrefix:         "path",
			expectedStatusCode: http.StatusOK,
			expectedResponse: ListPoliciesResponse{
				[]string{"policy1", "policy2"},
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
		},
		"ErrorCaseInvalidParameterError": {
			org:                "org1",
			pathPrefix:         "path",
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
			org:                "org1",
			pathPrefix:         "path",
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
			expectedStatusCode: http.StatusInternalServerError,
			getPolicyListErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[GetPolicyListMethod][0] = test.getPolicyListResult
		testApi.ArgsOut[GetPolicyListMethod][1] = test.getPolicyListErr

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/policies?PathPrefix=%v", test.org, test.pathPrefix)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error creating http request %v", n, err)
			continue
		}

		if test.pathPrefix != "" {
			q := req.URL.Query()
			q.Add("PathPrefix", test.pathPrefix)
			req.URL.RawQuery = q.Encode()
		}

		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error calling server %v", n, err)
			continue
		}
		if testApi.ArgsIn[GetPolicyListMethod][1] != test.org {
			t.Errorf("Test case %v. Received different Org (wanted:%v / received:%v)", n, test.org, testApi.ArgsIn[GetPolicyListMethod][1])
			continue
		}
		if testApi.ArgsIn[GetPolicyListMethod][2] != test.pathPrefix {
			t.Errorf("Test case %v. Received different PathPrefix (wanted:%v / received:%v)", n, test.pathPrefix, testApi.ArgsIn[GetPolicyListMethod][2])
			continue
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

func TestWorkerHandler_HandleCreatePolicy(t *testing.T) {
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
						Action: []string{
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
				Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]api.Statement{
					{
						Effect: "allow",
						Action: []string{
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
				Path:     "/path/",
				Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]api.Statement{
					{
						Effect: "allow",
						Action: []string{
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
						Action: []string{
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
						Action: []string{
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
						Action: []string{
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
						Action: []string{
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

func TestWorkerHandler_HandleDeletePolicy(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		org        string
		policyName string
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

		testApi.ArgsOut[DeletePolicyMethod][0] = test.deletePolicyErr

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/policies/%v", test.org, test.policyName)
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error creating http request %v", n, err)
			continue
		}

		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error calling server %v", n, err)
			continue
		}

		// Check received parameters
		if testApi.ArgsIn[DeletePolicyMethod][1] != test.org {
			t.Errorf("Test case %v. Received different Org (wanted:%v / received:%v)", n, test.org, testApi.ArgsIn[DeletePolicyMethod][1])
			continue
		}
		if testApi.ArgsIn[DeletePolicyMethod][2] != test.policyName {
			t.Errorf("Test case %v. Received different Name (wanted:%v / received:%v)", n, test.policyName, testApi.ArgsIn[DeletePolicyMethod][2])
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
						Action: []string{
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
				Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]api.Statement{
					{
						Effect: "allow",
						Action: []string{
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
				Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]api.Statement{
					{
						Effect: "allow",
						Action: []string{
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
						Action: []string{
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
						Action: []string{
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
						Action: []string{
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
						Action: []string{
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

func TestWorkerHandler_HandleGetPolicy(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		// API method args
		org        string
		policyName string
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
				Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]api.Statement{
					{
						Effect: "allow",
						Action: []string{
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
				Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]api.Statement{
					{
						Effect: "allow",
						Action: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
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

func TestWorkerHandler_HandleGetAttachedGroups(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		org        string
		policyName string
		// Expected result
		expectedStatusCode int
		expectedResponse   GetPolicyGroupsResponse
		expectedError      *api.Error
		// API Results
		getPolicyGroupsResult []string
		// API Errors
		getPolicyGroupsErr error
	}{
		"OkCase": {
			org:                "org1",
			policyName:         "p1",
			expectedStatusCode: http.StatusOK,
			expectedResponse: GetPolicyGroupsResponse{
				Groups: []string{"group1", "group2"},
			},
			getPolicyGroupsResult: []string{"group1", "group2"},
		},
		"ErrorCaseNotFound": {
			org:                "org1",
			policyName:         "p1",
			expectedStatusCode: http.StatusNotFound,
			expectedError: &api.Error{
				Code: api.POLICY_BY_ORG_AND_NAME_NOT_FOUND,
			},
			getPolicyGroupsErr: &api.Error{
				Code: api.POLICY_BY_ORG_AND_NAME_NOT_FOUND,
			},
		},
		"ErrorCaseUnauthorized": {
			org:                "org1",
			policyName:         "p1",
			expectedStatusCode: http.StatusForbidden,
			expectedError: &api.Error{
				Code: api.UNAUTHORIZED_RESOURCES_ERROR,
			},
			getPolicyGroupsErr: &api.Error{
				Code: api.UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"ErrorCaseInvalidParam": {
			org:                "org1",
			policyName:         "p1",
			expectedStatusCode: http.StatusBadRequest,
			expectedError: &api.Error{
				Code: api.INVALID_PARAMETER_ERROR,
			},
			getPolicyGroupsErr: &api.Error{
				Code: api.INVALID_PARAMETER_ERROR,
			},
		},
		"ErrorCaseInternalServerError": {
			org:                "org1",
			policyName:         "p1",
			expectedStatusCode: http.StatusInternalServerError,
			expectedError: &api.Error{
				Code: api.UNKNOWN_API_ERROR,
			},
			getPolicyGroupsErr: &api.Error{
				Code: api.UNKNOWN_API_ERROR,
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[GetAttachedGroupsMethod][0] = test.getPolicyGroupsResult
		testApi.ArgsOut[GetAttachedGroupsMethod][1] = test.getPolicyGroupsErr

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/policies/%v/groups", test.org, test.policyName)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error creating http request %v", n, err)
			continue
		}

		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error calling server %v", n, err)
			continue
		}

		// Check received parameters
		if testApi.ArgsIn[GetAttachedGroupsMethod][1] != test.org {
			t.Errorf("Test case %v. Received different org (wanted:%v / received:%v)", n, test.org, testApi.ArgsIn[GetAttachedGroupsMethod][1])
			continue
		}
		if testApi.ArgsIn[GetAttachedGroupsMethod][2] != test.policyName {
			t.Errorf("Test case %v. Received different Name (wanted:%v / received:%v)", n, test.policyName, testApi.ArgsIn[GetAttachedGroupsMethod][2])
			continue
		}

		// check status code
		if test.expectedStatusCode != res.StatusCode {
			t.Errorf("Test case %v. Received different http status code (wanted:%v / received:%v)", n, test.expectedStatusCode, res.StatusCode)
			continue
		}

		switch res.StatusCode {
		case http.StatusOK:
			GetAttachedGroupsMethodResponse := GetPolicyGroupsResponse{}
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

func TestWorkerHandler_HandleListAllPolicies(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		pathPrefix string
		// Expected result
		expectedStatusCode int
		expectedResponse   ListAllPoliciesResponse
		expectedError      api.Error
		// Manager Results
		getPolicyListResult []api.PolicyIdentity
		// Manager Errors
		getPolicyListErr error
	}{
		"OkCase": {
			pathPrefix:         "path",
			expectedStatusCode: http.StatusOK,
			expectedResponse: ListAllPoliciesResponse{
				[]api.PolicyIdentity{
					{
						Org:  "org1",
						Name: "policy1",
					},
					{
						Org:  "org1",
						Name: "policy2",
					},
				},
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
		},
		"ErrorCaseInvalidParameterError": {
			pathPrefix:         "path",
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
			pathPrefix:         "path",
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
			pathPrefix:         "path",
			expectedStatusCode: http.StatusInternalServerError,
			getPolicyListErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[GetPolicyListMethod][0] = test.getPolicyListResult
		testApi.ArgsOut[GetPolicyListMethod][1] = test.getPolicyListErr

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/policies?PathPrefix=%v", test.pathPrefix)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error creating http request %v", n, err)
			continue
		}

		if test.pathPrefix != "" {
			q := req.URL.Query()
			q.Add("PathPrefix", test.pathPrefix)
			req.URL.RawQuery = q.Encode()
		}

		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error calling server %v", n, err)
			continue
		}
		if testApi.ArgsIn[GetPolicyListMethod][1] != "" {
			t.Errorf("Test case %v. Received different Org (wanted:%v / received:%v)", n, "", testApi.ArgsIn[GetPolicyListMethod][1])
			continue
		}
		if testApi.ArgsIn[GetPolicyListMethod][2] != test.pathPrefix {
			t.Errorf("Test case %v. Received different PathPrefix (wanted:%v / received:%v)", n, test.pathPrefix, testApi.ArgsIn[GetPolicyListMethod][2])
			continue
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
