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
