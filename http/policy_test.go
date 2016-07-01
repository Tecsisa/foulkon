package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"
	"github.com/tecsisa/authorizr/api"
)

func TestWorkerHandler_HandleListPolicies(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		org        string
		pathPrefix string
		// Expected result
		expectedStatusCode int
		expectedResponse   ListPoliciesResponse
		expectedError      api.Error
		// Manager Results
		getListPoliciesResult []api.PolicyIdentity
		// Manager Errors
		getListPoliciesErr error
	}{
		"OkCase": {
			org:                "org1",
			pathPrefix:         "path",
			expectedStatusCode: http.StatusOK,
			expectedResponse: ListPoliciesResponse{
				[]api.PolicyIdentity{
					api.PolicyIdentity{
						Org:  "org1",
						Name: "policy1",
					},
				},
			},
			getListPoliciesResult: []api.PolicyIdentity{
				api.PolicyIdentity{
					Org:  "org1",
					Name: "policy1",
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
			getListPoliciesErr: &api.Error{
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
			getListPoliciesErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseUnknownApiError": {
			expectedStatusCode: http.StatusInternalServerError,
			getListPoliciesErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[GetListPoliciesMethod][0] = test.getListPoliciesResult
		testApi.ArgsOut[GetListPoliciesMethod][1] = test.getListPoliciesErr

		req, err := http.NewRequest(http.MethodGet, server.URL+API_VERSION_1+"/organizations/"+test.org+"/policies?PathPrefix="+test.pathPrefix, nil)
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
		if testApi.ArgsIn[GetListPoliciesMethod][1] != test.org {
			t.Errorf("Test case %v. Received different Org (wanted:%v / received:%v)", n, test.org, testApi.ArgsIn[GetListPoliciesMethod][1])
			continue
		}
		if testApi.ArgsIn[GetListPoliciesMethod][2] != test.pathPrefix {
			t.Errorf("Test case %v. Received different PathPrefix (wanted:%v / received:%v)", n, test.pathPrefix, testApi.ArgsIn[GetListPoliciesMethod][2])
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
		expectedResponse   CreatePolicyResponse
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
					api.Statement{
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
					api.Statement{
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
			expectedResponse: CreatePolicyResponse{
				Policy: &api.Policy{
					ID:       "test1",
					Name:     "test",
					Org:      "org1",
					CreateAt: now,
					Path:     "/path/",
					Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path/", "test"),
					Statements: &[]api.Statement{
						api.Statement{
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
					api.Statement{
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
					api.Statement{
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
					api.Statement{
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
					api.Statement{
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

		req, err := http.NewRequest(http.MethodPost, server.URL+API_VERSION_1+"/organizations/"+test.org+"/policies", body)
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
			createPolicyResponse := CreatePolicyResponse{}
			err = json.NewDecoder(res.Body).Decode(&createPolicyResponse)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(createPolicyResponse, test.expectedResponse); diff != "" {
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

		req, err := http.NewRequest(http.MethodDelete, server.URL+API_VERSION_1+"/organizations/"+test.org+"/policies/"+test.policyName, nil)
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

func TestWorkerHandler_HandleGetPolicy(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		// API method args
		org        string
		policyName string
		// Expected result
		expectedStatusCode int
		expectedResponse   GetPolicyResponse
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
			expectedResponse: GetPolicyResponse{
				&api.Policy{
					ID:       "test1",
					Name:     "test",
					Org:      "org1",
					Path:     "/path/",
					CreateAt: now,
					Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path/", "test"),
					Statements: &[]api.Statement{
						api.Statement{
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
			getPolicyByNameResult: &api.Policy{
				ID:       "test1",
				Name:     "test",
				Org:      "org1",
				Path:     "/path/",
				CreateAt: now,
				Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]api.Statement{
					api.Statement{
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

		req, err := http.NewRequest(http.MethodGet, server.URL+API_VERSION_1+"/organizations/"+test.org+"/policies/"+test.policyName, nil)
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
			GetPolicyByNameMethodResponse := GetPolicyResponse{}
			err = json.NewDecoder(res.Body).Decode(&GetPolicyByNameMethodResponse)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(GetPolicyByNameMethodResponse, test.expectedResponse); diff != "" {
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
