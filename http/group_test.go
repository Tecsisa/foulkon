package http

import (
	"encoding/json"
	"net/http"
	"testing"

	"time"

	"bytes"
	"github.com/kylelemons/godebug/pretty"
	"github.com/tecsisa/authorizr/api"
)

func TestWorkerHandler_HandleGetGroup(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		// API method args
		org  string
		name string
		// Expected result
		expectedStatusCode int
		expectedResponse   GetGroupNameResponse
		expectedError      api.Error
		// Manager Results
		getGroupByNameResult *api.Group
		// Manager Errors
		getGroupByNameErr error
	}{
		"OkCase": {
			org:                "org1",
			name:               "group1",
			expectedStatusCode: http.StatusOK,
			expectedResponse: GetGroupNameResponse{
				Group: &api.Group{
					ID:       "groupID",
					Name:     "group1",
					Path:     "Path",
					Urn:      "Urn",
					Org:      "Org",
					CreateAt: now,
				},
			},
			getGroupByNameResult: &api.Group{
				ID:       "groupID",
				Name:     "group1",
				Path:     "Path",
				Urn:      "Urn",
				Org:      "Org",
				CreateAt: now,
			},
		},
		"ErrorCaseGroupNotFound": {
			name:               "group1",
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.GROUP_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Group not found",
			},
			getGroupByNameErr: &api.Error{
				Code:    api.GROUP_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Group not found",
			},
		},
		"ErrorCaseUnauthorizedResourcesError": {
			name:               "group1",
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			getGroupByNameErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseInvalidParameterError": {
			name:               "group1",
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
			getGroupByNameErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
		},
		"ErrorCaseUnknownApiError": {
			name:               "group1",
			expectedStatusCode: http.StatusInternalServerError,
			getGroupByNameErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[GetGroupByNameMethod][0] = test.getGroupByNameResult
		testApi.ArgsOut[GetGroupByNameMethod][1] = test.getGroupByNameErr

		req, err := http.NewRequest(http.MethodGet, server.URL+API_VERSION_1+"/organizations/"+test.org+"/groups/"+test.name, nil)
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
		if testApi.ArgsIn[GetGroupByNameMethod][1] != test.org {
			t.Errorf("Test case %v. Received different org (wanted:%v / received:%v)", n, test.org, testApi.ArgsIn[GetGroupByNameMethod][1])
			continue
		}
		if testApi.ArgsIn[GetGroupByNameMethod][2] != test.name {
			t.Errorf("Test case %v. Received different Name (wanted:%v / received:%v)", n, test.name, testApi.ArgsIn[GetGroupByNameMethod][2])
			continue
		}

		// check status code
		if test.expectedStatusCode != res.StatusCode {
			t.Errorf("Test case %v. Received different http status code (wanted:%v / received:%v)", n, test.expectedStatusCode, res.StatusCode)
			continue
		}

		switch res.StatusCode {
		case http.StatusOK:
			getGroupNameResponse := GetGroupNameResponse{}
			err = json.NewDecoder(res.Body).Decode(&getGroupNameResponse)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(getGroupNameResponse, test.expectedResponse); diff != "" {
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

func TestWorkerHandler_HandleCreateGroup(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		// API method args
		org     string
		request *CreateGroupRequest
		// Expected result
		expectedStatusCode int
		expectedResponse   CreateGroupResponse
		expectedError      api.Error
		// Manager Results
		addGroupResult *api.Group
		// Manager Errors
		addGroupErr error
	}{
		"OkCase": {
			org: "org1",
			request: &CreateGroupRequest{
				Name: "group1",
				Path: "Path",
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse: CreateGroupResponse{
				Group: &api.Group{
					ID:       "GroupID",
					Name:     "group1",
					Path:     "Path",
					Urn:      "Urn",
					Org:      "org1",
					CreateAt: now,
				},
			},
			addGroupResult: &api.Group{
				ID:       "GroupID",
				Name:     "group1",
				Path:     "Path",
				Urn:      "Urn",
				Org:      "org1",
				CreateAt: now,
			},
		},
		"ErrorCaseMalformedRequest": {
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "EOF",
			},
		},
		"ErrorCaseGroupAlreadyExist": {
			org: "org1",
			request: &CreateGroupRequest{
				Name: "group1",
				Path: "Path",
			},
			expectedStatusCode: http.StatusConflict,
			expectedError: api.Error{
				Code:    api.GROUP_ALREADY_EXIST,
				Message: "Group already exist",
			},
			addGroupErr: &api.Error{
				Code:    api.GROUP_ALREADY_EXIST,
				Message: "Group already exist",
			},
		},
		"ErrorCaseInvalidParameterError": {
			org: "org1",
			request: &CreateGroupRequest{
				Name: "group1",
				Path: "Path",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid Parameter",
			},
			addGroupErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid Parameter",
			},
		},
		"ErrorCaseUnauthorizedResourcesError": {
			org: "org1",
			request: &CreateGroupRequest{
				Name: "group1",
				Path: "Path",
			},
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			addGroupErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseUnknownApiError": {
			org: "org1",
			request: &CreateGroupRequest{
				Name: "group1",
				Path: "Path",
			},
			expectedStatusCode: http.StatusInternalServerError,
			addGroupErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[AddGroupMethod][0] = test.addGroupResult
		testApi.ArgsOut[AddGroupMethod][1] = test.addGroupErr

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

		req, err := http.NewRequest(http.MethodPost, server.URL+API_VERSION_1+"/organizations/"+test.org+"/groups", body)
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
			if testApi.ArgsIn[AddGroupMethod][1] != test.org {
				t.Errorf("Test case %v. Received different Org (wanted:%v / received:%v)", n, test.org, testApi.ArgsIn[AddGroupMethod][1])
				continue
			}
			if testApi.ArgsIn[AddGroupMethod][2] != test.request.Name {
				t.Errorf("Test case %v. Received different Name (wanted:%v / received:%v)", n, test.request.Name, testApi.ArgsIn[AddGroupMethod][2])
				continue
			}
			if testApi.ArgsIn[AddGroupMethod][3] != test.request.Path {
				t.Errorf("Test case %v. Received different Path (wanted:%v / received:%v)", n, test.request.Path, testApi.ArgsIn[AddGroupMethod][3])
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
			createGroupResponse := CreateGroupResponse{}
			err = json.NewDecoder(res.Body).Decode(&createGroupResponse)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(createGroupResponse, test.expectedResponse); diff != "" {
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

func TestWorkerHandler_HandleDeleteGroup(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		org  string
		name string
		// Expected result
		expectedStatusCode int
		expectedError      api.Error
		// Manager Errors
		removeGroupErr error
	}{
		"OkCase": {
			org:                "org1",
			name:               "group1",
			expectedStatusCode: http.StatusNoContent,
		},
		"ErrorCaseGroupNotFound": {
			org:                "org1",
			name:               "group1",
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.GROUP_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Group not found",
			},
			removeGroupErr: &api.Error{
				Code:    api.GROUP_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Group not found",
			},
		},
		"ErrorCaseInvalidParameterError": {
			org:                "org1",
			name:               "InvalidID",
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
			removeGroupErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
		},
		"ErrorCaseUnauthorizedResourcesError": {
			org:                "org1",
			name:               "UnauthorizedID",
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			removeGroupErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseUnknownApiError": {
			org:                "org1",
			name:               "ExceptionID",
			expectedStatusCode: http.StatusInternalServerError,
			removeGroupErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[RemoveGroupMethod][0] = test.removeGroupErr

		req, err := http.NewRequest(http.MethodDelete, server.URL+API_VERSION_1+"/organizations/"+test.org+"/groups/"+test.name, nil)
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
		if testApi.ArgsIn[RemoveGroupMethod][1] != test.org {
			t.Errorf("Test case %v. Received different Org (wanted:%v / received:%v)", n, test.org, testApi.ArgsIn[RemoveGroupMethod][1])
			continue
		}
		if testApi.ArgsIn[RemoveGroupMethod][2] != test.name {
			t.Errorf("Test case %v. Received different Name (wanted:%v / received:%v)", n, test.name, testApi.ArgsIn[RemoveGroupMethod][2])
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
				t.Errorf("Test %v failed. Received different error response (received/wanted) %v",
					n, diff)
				continue
			}

		}

	}
}

func TestWorkerHandler_HandleListGroups(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		org        string
		pathPrefix string
		// Expected result
		expectedStatusCode int
		expectedResponse   GetGroupsResponse
		expectedError      api.Error
		// Manager Results
		getListGroupResult []api.GroupIdentity
		// Manager Errors
		getListGroupsErr error
	}{
		"OkCase": {
			org:                "org1",
			pathPrefix:         "path",
			expectedStatusCode: http.StatusOK,
			expectedResponse: GetGroupsResponse{
				[]api.GroupIdentity{
					api.GroupIdentity{
						Org:  "org1",
						Name: "group1",
					},
				},
			},
			getListGroupResult: []api.GroupIdentity{
				api.GroupIdentity{
					Org:  "org1",
					Name: "group1",
				},
			},
		},
		"ErrorCaseUnauthorizedError": {
			pathPrefix:         "Path",
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			getListGroupsErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseUnknownApiError": {
			expectedStatusCode: http.StatusInternalServerError,
			getListGroupsErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[GetListGroupsMethod][0] = test.getListGroupResult
		testApi.ArgsOut[GetListGroupsMethod][1] = test.getListGroupsErr

		req, err := http.NewRequest(http.MethodGet, server.URL+API_VERSION_1+"/organizations/"+test.org+"/groups?PathPrefix="+test.pathPrefix, nil)
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

		// Check received parameter
		if testApi.ArgsIn[GetListGroupsMethod][1] != test.org {
			t.Errorf("Test case %v. Received different PathPrefix (wanted:%v / received:%v)", n, test.org, testApi.ArgsIn[GetListGroupsMethod][1])
			continue
		}
		if testApi.ArgsIn[GetListGroupsMethod][2] != test.pathPrefix {
			t.Errorf("Test case %v. Received different PathPrefix (wanted:%v / received:%v)", n, test.pathPrefix, testApi.ArgsIn[GetListGroupsMethod][2])
			continue
		}

		// check status code
		if test.expectedStatusCode != res.StatusCode {
			t.Errorf("Test case %v. Received different http status code (wanted:%v / received:%v)", n, test.expectedStatusCode, res.StatusCode)
			continue
		}

		switch res.StatusCode {
		case http.StatusOK:
			getGroupsResponse := GetGroupsResponse{}
			err = json.NewDecoder(res.Body).Decode(&getGroupsResponse)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(getGroupsResponse, test.expectedResponse); diff != "" {
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
