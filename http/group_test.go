package http

import (
	"encoding/json"
	"net/http"
	"testing"

	"time"

	"bytes"
	"fmt"

	"github.com/Tecsisa/foulkon/api"
	"github.com/stretchr/testify/assert"
)

func TestWorkerHandler_HandleAddGroup(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		// API method args
		org     string
		request *CreateGroupRequest
		// Expected result
		expectedStatusCode int
		expectedResponse   api.Group
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
			expectedResponse: api.Group{
				ID:       "GroupID",
				Name:     "group1",
				Path:     "Path",
				Urn:      "Urn",
				Org:      "org1",
				CreateAt: now,
				UpdateAt: now,
			},
			addGroupResult: &api.Group{
				ID:       "GroupID",
				Name:     "group1",
				Path:     "Path",
				Urn:      "Urn",
				Org:      "org1",
				CreateAt: now,
				UpdateAt: now,
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
			assert.Nil(t, err, "Error in test case %v", n)
			body = bytes.NewBuffer(jsonObject)
		}
		if body == nil {
			body = bytes.NewBuffer([]byte{})
		}

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/groups", test.org)
		req, err := http.NewRequest(http.MethodPost, url, body)
		assert.Nil(t, err, "Error in test case %v", n)

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if test.request != nil {
			// Check received parameters
			assert.Equal(t, test.org, testApi.ArgsIn[AddGroupMethod][1], "Error in test case %v", n)
			assert.Equal(t, test.request.Name, testApi.ArgsIn[AddGroupMethod][2], "Error in test case %v", n)
			assert.Equal(t, test.request.Path, testApi.ArgsIn[AddGroupMethod][3], "Error in test case %v", n)
		}

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusCreated:
			response := api.Group{}
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

func TestWorkerHandler_HandleGetGroupByName(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		// API method args
		org          string
		name         string
		offset       string
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedResponse   api.Group
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
			expectedResponse: api.Group{
				ID:       "groupID",
				Name:     "group1",
				Path:     "Path",
				Urn:      "Urn",
				Org:      "Org",
				CreateAt: now,
				UpdateAt: now,
			},
			getGroupByNameResult: &api.Group{
				ID:       "groupID",
				Name:     "group1",
				Path:     "Path",
				Urn:      "Urn",
				Org:      "Org",
				CreateAt: now,
				UpdateAt: now,
			},
		},
		"ErrorCaseInvalidRequest": {
			name:               "group1",
			org:                "org1",
			offset:             "-1",
			ignoreArgsIn:       true,
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
			getGroupByNameErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
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

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/groups/%v", test.org, test.name)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		q := req.URL.Query()
		q.Add("Offset", test.offset)
		req.URL.RawQuery = q.Encode()

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameters
			assert.Equal(t, test.org, testApi.ArgsIn[GetGroupByNameMethod][1], "Error in test case %v", n)
			assert.Equal(t, test.name, testApi.ArgsIn[GetGroupByNameMethod][2], "Error in test case %v", n)
		}
		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusOK:
			response := api.Group{}
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

func TestWorkerHandler_HandleListGroups(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		filter       *api.Filter
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedResponse   ListGroupsResponse
		expectedError      api.Error
		// Manager Results
		getListGroupResult []api.GroupIdentity
		totalGroupsResult  int
		// Manager Errors
		getListGroupsErr error
	}{
		"OkCase": {
			filter: &api.Filter{
				PathPrefix: "path",
				Org:        "org",
				Offset:     0,
				Limit:      0,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: ListGroupsResponse{
				Groups: []string{"group1"},
				Offset: 0,
				Limit:  0,
				Total:  1,
			},
			getListGroupResult: []api.GroupIdentity{
				{
					Org:  "org1",
					Name: "group1",
				},
			},
			totalGroupsResult: 1,
		},
		"ErrorCaseInvalidFilterParams": {
			filter: &api.Filter{
				PathPrefix: "",
				Offset:     -1,
			},
			ignoreArgsIn:       true,
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
		},
		"ErrorCaseUnauthorizedError": {
			filter: &api.Filter{
				PathPrefix: "Path",
				Org:        "org1",
				Offset:     0,
				Limit:      0,
			},
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
		"ErrorCaseInvalidParameterError": {
			filter: &api.Filter{
				PathPrefix: "Invalid",
				Org:        "org1",
				Offset:     0,
				Limit:      0,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid Path",
			},
			getListGroupsErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid Path",
			},
		},
		"ErrorCaseUnknownApiError": {
			filter: &api.Filter{
				PathPrefix: "path",
				Org:        "org1",
				Offset:     0,
				Limit:      0,
			},
			expectedStatusCode: http.StatusInternalServerError,
			getListGroupsErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[ListGroupsMethod][0] = test.getListGroupResult
		testApi.ArgsOut[ListGroupsMethod][1] = test.totalGroupsResult
		testApi.ArgsOut[ListGroupsMethod][2] = test.getListGroupsErr

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/groups", test.filter.Org)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		addQueryParams(test.filter, req)

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameter
			filterData, ok := testApi.ArgsIn[ListGroupsMethod][1].(*api.Filter)
			if ok {
				// Check result
				assert.Equal(t, test.filter, filterData, "Error in test case %v", n)
			}
		}

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusOK:
			listGroupsResponse := ListGroupsResponse{}
			err = json.NewDecoder(res.Body).Decode(&listGroupsResponse)
			assert.Nil(t, err, "Error in test case %v", n)
			// Check result
			assert.Equal(t, test.expectedResponse, listGroupsResponse, "Error in test case %v", n)
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

func TestWorkerHandler_HandleListAllGroups(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		filter       *api.Filter
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedResponse   ListAllGroupsResponse
		expectedError      api.Error
		// Manager Results
		getListAllGroupResult []api.GroupIdentity
		totalGroupsResult     int
		// Manager Errors
		getListAllGroupErr error
	}{
		"OkCase": {
			filter: &api.Filter{
				PathPrefix: "/path/",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: ListAllGroupsResponse{
				Groups: []api.GroupIdentity{
					{
						Org:  "org1",
						Name: "group1",
					},
				},
				Offset: 0,
				Limit:  0,
				Total:  1,
			},
			getListAllGroupResult: []api.GroupIdentity{
				{
					Org:  "org1",
					Name: "group1",
				},
			},
			totalGroupsResult: 1,
		},
		"ErrorCaseInvalidFilterParams": {
			filter: &api.Filter{
				PathPrefix: "",
				Offset:     -1,
			},
			ignoreArgsIn:       true,
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
		},
		"ErrorCaseUnauthorizedError": {
			filter: &api.Filter{
				PathPrefix: "/path/",
				Offset:     0,
				Limit:      0,
			},
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			getListAllGroupErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseInvalidParameterError": {
			filter: &api.Filter{
				PathPrefix: "Invalid",
				Offset:     0,
				Limit:      0,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid Path",
			},
			getListAllGroupErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid Path",
			},
		},
		"ErrorCaseUnknownApiError": {
			filter: &api.Filter{
				PathPrefix: "/path/",
				Offset:     0,
				Limit:      0,
			},
			expectedStatusCode: http.StatusInternalServerError,
			getListAllGroupErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[ListGroupsMethod][0] = test.getListAllGroupResult
		testApi.ArgsOut[ListGroupsMethod][1] = test.totalGroupsResult
		testApi.ArgsOut[ListGroupsMethod][2] = test.getListAllGroupErr

		url := fmt.Sprintf(server.URL + API_VERSION_1 + "/groups")
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		addQueryParams(test.filter, req)

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameter
			filterData, ok := testApi.ArgsIn[ListGroupsMethod][1].(*api.Filter)
			if ok {
				// Check result
				assert.Equal(t, test.filter, filterData, "Error in test case %v", n)
			}
		}

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusOK:
			listAllGroupsResponse := ListAllGroupsResponse{}
			err = json.NewDecoder(res.Body).Decode(&listAllGroupsResponse)
			assert.Nil(t, err, "Error in test case %v", n)
			// Check result
			assert.Equal(t, test.expectedResponse, listAllGroupsResponse, "Error in test case %v", n)
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

func TestWorkerHandler_HandleUpdateGroup(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		// API method args
		org     string
		request *UpdateGroupRequest
		// Expected result
		expectedStatusCode int
		expectedResponse   api.Group
		expectedError      api.Error
		// Manager Results
		updateGroupResult *api.Group
		// Manager Errors
		updateGroupErr error
	}{
		"OkCase": {
			org: "org1",
			request: &UpdateGroupRequest{
				Name: "newName",
				Path: "NewPath",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: api.Group{
				ID:       "GroupID",
				Name:     "newName",
				Path:     "NewPath",
				Urn:      "urn",
				CreateAt: now,
				UpdateAt: now,
			},
			updateGroupResult: &api.Group{
				ID:       "GroupID",
				Name:     "newName",
				Path:     "NewPath",
				Urn:      "urn",
				CreateAt: now,
				UpdateAt: now,
			},
		},
		"ErrorCaseMalformedRequest": {
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "EOF",
			},
		},
		"ErrorCaseGroupNotFound": {
			request: &UpdateGroupRequest{
				Name: "newName",
				Path: "NewPath",
			},
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.GROUP_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Group not found",
			},
			updateGroupErr: &api.Error{
				Code:    api.GROUP_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Group not found",
			},
		},
		"ErrorCaseInvalidParameterError": {
			request: &UpdateGroupRequest{
				Name: "newName",
				Path: "InvalidPath",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
			updateGroupErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
		},
		"ErrorCaseGroupAlreadyExistError": {
			request: &UpdateGroupRequest{
				Name: "newName",
				Path: "newPath",
			},
			expectedStatusCode: http.StatusConflict,
			expectedError: api.Error{
				Code:    api.GROUP_ALREADY_EXIST,
				Message: "Group already exist",
			},
			updateGroupErr: &api.Error{
				Code:    api.GROUP_ALREADY_EXIST,
				Message: "Group already exist",
			},
		},
		"ErrorCaseUnauthorizedResourcesError": {
			request: &UpdateGroupRequest{
				Name: "newName",
				Path: "NewPath",
			},
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			updateGroupErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseUnknownApiError": {
			request: &UpdateGroupRequest{
				Name: "newName",
				Path: "NewPath",
			},
			expectedStatusCode: http.StatusInternalServerError,
			updateGroupErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[UpdateGroupMethod][0] = test.updateGroupResult
		testApi.ArgsOut[UpdateGroupMethod][1] = test.updateGroupErr

		var body *bytes.Buffer
		if test.request != nil {
			jsonObject, err := json.Marshal(test.request)
			assert.Nil(t, err, "Error in test case %v", n)
			body = bytes.NewBuffer(jsonObject)
		}
		if body == nil {
			body = bytes.NewBuffer([]byte{})
		}

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/groups/group1", test.org)
		req, err := http.NewRequest(http.MethodPut, url, body)
		assert.Nil(t, err, "Error in test case %v", n)

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if test.request != nil {
			// Check received parameters
			assert.Equal(t, test.org, testApi.ArgsIn[UpdateGroupMethod][1], "Error in test case %v", n)
			assert.Equal(t, "group1", testApi.ArgsIn[UpdateGroupMethod][2], "Error in test case %v", n)
			assert.Equal(t, test.request.Name, testApi.ArgsIn[UpdateGroupMethod][3], "Error in test case %v", n)
			assert.Equal(t, test.request.Path, testApi.ArgsIn[UpdateGroupMethod][4], "Error in test case %v", n)
		}

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusOK:
			response := api.Group{}
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

func TestWorkerHandler_HandleRemoveGroup(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		org          string
		name         string
		offset       string
		ignoreArgsIn bool
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
		"ErrorCaseInvalidRequest": {
			org:                "org1",
			name:               "group1",
			offset:             "-1",
			ignoreArgsIn:       true,
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
			removeGroupErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
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

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/groups/%v", test.org, test.name)
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		q := req.URL.Query()
		q.Add("Offset", test.offset)
		req.URL.RawQuery = q.Encode()

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameters
			assert.Equal(t, test.org, testApi.ArgsIn[RemoveGroupMethod][1], "Error in test case %v", n)
			assert.Equal(t, test.name, testApi.ArgsIn[RemoveGroupMethod][2], "Error in test case %v", n)
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

func TestWorkerHandler_HandleAddMember(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		org          string
		userID       string
		groupName    string
		offset       string
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedError      api.Error
		// Manager Errors
		addMemberErr error
	}{
		"OkCase": {
			org:                "org1",
			userID:             "user1",
			groupName:          "group1",
			expectedStatusCode: http.StatusNoContent,
		},
		"ErrorCaseInvalidRequest": {
			org:                "org1",
			userID:             "user1",
			groupName:          "group1",
			offset:             "-1",
			ignoreArgsIn:       true,
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
			addMemberErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
		},
		"ErrorCaseGroupNotFoundErr": {
			org:                "org1",
			userID:             "user1",
			groupName:          "Invalid Group",
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.GROUP_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Group Not Found",
			},
			addMemberErr: &api.Error{
				Code:    api.GROUP_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Group Not Found",
			},
		},
		"ErrorCaseUserNotFoundErr": {
			org:                "org1",
			userID:             "Invalid User",
			groupName:          "group1",
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.USER_BY_EXTERNAL_ID_NOT_FOUND,
				Message: "User Not Found",
			},
			addMemberErr: &api.Error{
				Code:    api.USER_BY_EXTERNAL_ID_NOT_FOUND,
				Message: "User Not Found",
			},
		},
		"ErrorCaseUnauthorizedError": {
			org:                "org1",
			userID:             "user1",
			groupName:          "group1",
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			addMemberErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseInvalidParameterErr": {
			org:                "org1",
			userID:             "user1",
			groupName:          "group1",
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid Parameter",
			},
			addMemberErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid Parameter",
			},
		},
		"ErrorCaseUserIsAlreadyMemberErr": {
			org:                "org1",
			userID:             "user1",
			groupName:          "group1",
			expectedStatusCode: http.StatusConflict,
			expectedError: api.Error{
				Code:    api.USER_IS_ALREADY_A_MEMBER_OF_GROUP,
				Message: "User is already a member of group",
			},
			addMemberErr: &api.Error{
				Code:    api.USER_IS_ALREADY_A_MEMBER_OF_GROUP,
				Message: "User is already a member of group",
			},
		},
		"ErrorCaseUnknownApiError": {
			org:                "org1",
			userID:             "user1",
			groupName:          "group1",
			expectedStatusCode: http.StatusInternalServerError,
			addMemberErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[AddMemberMethod][0] = test.addMemberErr

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/groups/%v/users/%v", test.org, test.groupName, test.userID)
		req, err := http.NewRequest(http.MethodPost, url, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		q := req.URL.Query()
		q.Add("Offset", test.offset)
		req.URL.RawQuery = q.Encode()

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameters
			assert.Equal(t, test.userID, testApi.ArgsIn[AddMemberMethod][1], "Error in test case %v", n)
			assert.Equal(t, test.groupName, testApi.ArgsIn[AddMemberMethod][2], "Error in test case %v", n)
			assert.Equal(t, test.org, testApi.ArgsIn[AddMemberMethod][3], "Error in test case %v", n)
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

func TestWorkerHandler_HandleRemoveMember(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		org          string
		userID       string
		groupName    string
		offset       string
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedError      api.Error
		// Manager Errors
		removeMemberErr error
	}{
		"OkCase": {
			org:                "org1",
			userID:             "user1",
			groupName:          "group1",
			expectedStatusCode: http.StatusNoContent,
		},
		"ErrorCaseInvalidRequest": {
			org:                "org1",
			userID:             "user1",
			groupName:          "group1",
			offset:             "-1",
			ignoreArgsIn:       true,
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
			removeMemberErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
		},
		"ErrorCaseGroupNotFoundErr": {
			org:                "org1",
			userID:             "user1",
			groupName:          "Invalid Group",
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.GROUP_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Group Not Found",
			},
			removeMemberErr: &api.Error{
				Code:    api.GROUP_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Group Not Found",
			},
		},
		"ErrorCaseUserNotFoundErr": {
			org:                "org1",
			userID:             "Invalid User",
			groupName:          "group1",
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.USER_BY_EXTERNAL_ID_NOT_FOUND,
				Message: "User Not Found",
			},
			removeMemberErr: &api.Error{
				Code:    api.USER_BY_EXTERNAL_ID_NOT_FOUND,
				Message: "User Not Found",
			},
		},
		"ErrorCaseUserIsNotMemberErr": {
			org:                "org1",
			userID:             "user1",
			groupName:          "group1",
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.USER_IS_NOT_A_MEMBER_OF_GROUP,
				Message: "User is not a member",
			},
			removeMemberErr: &api.Error{
				Code:    api.USER_IS_NOT_A_MEMBER_OF_GROUP,
				Message: "User is not a member",
			},
		},
		"ErrorCaseUnauthorizedError": {
			org:                "org1",
			userID:             "user1",
			groupName:          "group1",
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			removeMemberErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseInvalidParameterErr": {
			org:                "org1",
			userID:             "user1",
			groupName:          "group1",
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid Parameter",
			},
			removeMemberErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid Parameter",
			},
		},
		"ErrorCaseUnknownApiError": {
			org:                "org1",
			userID:             "user1",
			groupName:          "group1",
			expectedStatusCode: http.StatusInternalServerError,
			removeMemberErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[RemoveMemberMethod][0] = test.removeMemberErr

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/groups/%v/users/%v", test.org, test.groupName, test.userID)
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		q := req.URL.Query()
		q.Add("Offset", test.offset)
		req.URL.RawQuery = q.Encode()

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameters
			assert.Equal(t, test.userID, testApi.ArgsIn[RemoveMemberMethod][1], "Error in test case %v", n)
			assert.Equal(t, test.groupName, testApi.ArgsIn[RemoveMemberMethod][2], "Error in test case %v", n)
			assert.Equal(t, test.org, testApi.ArgsIn[RemoveMemberMethod][3], "Error in test case %v", n)
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

func TestWorkerHandler_HandleListMembers(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		// API method args
		filter       *api.Filter
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedResponse   ListMembersResponse
		expectedError      api.Error
		// Manager Results
		getListMembersResult []api.GroupMembers
		totalGroupsResult    int
		// Manager Errors
		getListMembersErr error
	}{
		"OkCase": {
			filter: &api.Filter{
				Org:       "org1",
				GroupName: "group1",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: ListMembersResponse{
				Members: []api.GroupMembers{
					{
						User:     "member1",
						CreateAt: now,
					},
					{
						User:     "member2",
						CreateAt: now,
					},
				},
				Offset: 0,
				Limit:  0,
				Total:  2,
			},
			getListMembersResult: []api.GroupMembers{
				{
					User:     "member1",
					CreateAt: now,
				},
				{
					User:     "member2",
					CreateAt: now,
				},
			},
			totalGroupsResult: 2,
		},
		"ErrorCaseInvalidFilterParams": {
			filter: &api.Filter{
				PathPrefix: "",
				Offset:     -1,
			},
			ignoreArgsIn:       true,
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
		},
		"ErrorCaseGroupNotFoundErr": {
			filter: &api.Filter{
				Org:       "org1",
				GroupName: "group1",
			},
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.GROUP_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Group Not Found",
			},
			getListMembersErr: &api.Error{
				Code:    api.GROUP_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Group Not Found",
			},
		},
		"ErrorCaseInvalidParameterErr": {
			filter: &api.Filter{
				Org:       "org1",
				GroupName: "group1",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid Parameter",
			},
			getListMembersErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid Parameter",
			},
		},
		"ErrorCaseUnauthorizedError": {
			filter: &api.Filter{
				Org:       "org1",
				GroupName: "group1",
			},
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			getListMembersErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseUnknownApiError": {
			filter:             testFilter,
			expectedStatusCode: http.StatusInternalServerError,
			getListMembersErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[ListMembersMethod][0] = test.getListMembersResult
		testApi.ArgsOut[ListMembersMethod][1] = test.totalGroupsResult
		testApi.ArgsOut[ListMembersMethod][2] = test.getListMembersErr

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/groups/%v/users", test.filter.Org, test.filter.GroupName)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		addQueryParams(test.filter, req)

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameter
			filterData, ok := testApi.ArgsIn[ListMembersMethod][1].(*api.Filter)
			if ok {
				// Check result
				assert.Equal(t, test.filter, filterData, "Error in test case %v", n)
			}

		}

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusOK:
			getGroupMembersResponse := ListMembersResponse{}
			err = json.NewDecoder(res.Body).Decode(&getGroupMembersResponse)
			assert.Nil(t, err, "Error in test case %v", n)
			// Check result
			assert.Equal(t, test.expectedResponse, getGroupMembersResponse, "Error in test case %v", n)
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

func TestWorkerHandler_HandleAttachPolicyToGroup(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		org          string
		groupName    string
		policyName   string
		offset       string
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedError      api.Error
		// Manager Errors
		attachGroupPolicyErr error
	}{
		"OkCase": {
			org:                "org1",
			groupName:          "group1",
			policyName:         "policy1",
			expectedStatusCode: http.StatusNoContent,
		},
		"ErrorCaseInvalidRequest": {
			org:                "org1",
			groupName:          "group1",
			policyName:         "policy1",
			offset:             "-1",
			ignoreArgsIn:       true,
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
			attachGroupPolicyErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
		},
		"ErrorCaseGroupNotFoundErr": {
			org:                "org1",
			groupName:          "Invalid Group",
			policyName:         "policy1",
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.GROUP_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Group Not Found",
			},
			attachGroupPolicyErr: &api.Error{
				Code:    api.GROUP_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Group Not Found",
			},
		},
		"ErrorCasePolicyNotFoundErr": {
			org:                "org1",
			groupName:          "group1",
			policyName:         "Invalid Policy",
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.POLICY_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "User Not Found",
			},
			attachGroupPolicyErr: &api.Error{
				Code:    api.POLICY_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "User Not Found",
			},
		},
		"ErrorCaseUnauthorizedError": {
			org:                "org1",
			groupName:          "group1",
			policyName:         "policy1",
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			attachGroupPolicyErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseInvalidParameterErr": {
			org:                "org1",
			groupName:          "group1",
			policyName:         "policy1",
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid Parameter",
			},
			attachGroupPolicyErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid Parameter",
			},
		},
		"ErrorCasePolicyIsAlreadyAttachedErr": {
			org:                "org1",
			groupName:          "group1",
			policyName:         "policy1",
			expectedStatusCode: http.StatusConflict,
			expectedError: api.Error{
				Code:    api.POLICY_IS_ALREADY_ATTACHED_TO_GROUP,
				Message: "Policy is already attached to group",
			},
			attachGroupPolicyErr: &api.Error{
				Code:    api.POLICY_IS_ALREADY_ATTACHED_TO_GROUP,
				Message: "Policy is already attached to group",
			},
		},
		"ErrorCaseUnknownApiError": {
			org:                "org1",
			groupName:          "group1",
			policyName:         "policy1",
			expectedStatusCode: http.StatusInternalServerError,
			attachGroupPolicyErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[AttachPolicyToGroupMethod][0] = test.attachGroupPolicyErr

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/groups/%v/policies/%v", test.org, test.groupName, test.policyName)
		req, err := http.NewRequest(http.MethodPost, url, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		q := req.URL.Query()
		q.Add("Offset", test.offset)
		req.URL.RawQuery = q.Encode()

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameters
			assert.Equal(t, test.org, testApi.ArgsIn[AttachPolicyToGroupMethod][1], "Error in test case %v", n)
			assert.Equal(t, test.groupName, testApi.ArgsIn[AttachPolicyToGroupMethod][2], "Error in test case %v", n)
			assert.Equal(t, test.policyName, testApi.ArgsIn[AttachPolicyToGroupMethod][3], "Error in test case %v", n)
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

func TestWorkerHandler_HandleDetachPolicyToGroup(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		org          string
		groupName    string
		policyName   string
		offset       string
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedError      api.Error
		// Manager Errors
		detachGroupPolicyErr error
	}{
		"OkCase": {
			org:                "org1",
			groupName:          "group1",
			policyName:         "policy1",
			expectedStatusCode: http.StatusNoContent,
		},
		"ErrorCaseInvalidRequest": {
			org:                "org1",
			groupName:          "group1",
			policyName:         "policy1",
			offset:             "-1",
			ignoreArgsIn:       true,
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
			detachGroupPolicyErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
		},
		"ErrorCaseGroupNotFoundErr": {
			org:                "org1",
			groupName:          "Invalid Group",
			policyName:         "policy1",
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.GROUP_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Group Not Found",
			},
			detachGroupPolicyErr: &api.Error{
				Code:    api.GROUP_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Group Not Found",
			},
		},
		"ErrorCasePolicyNotFoundErr": {
			org:                "org1",
			groupName:          "group1",
			policyName:         "Invalid Policy",
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.POLICY_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "User Not Found",
			},
			detachGroupPolicyErr: &api.Error{
				Code:    api.POLICY_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "User Not Found",
			},
		},
		"ErrorCasePolicyIsNotAttachedErr": {
			org:                "org1",
			groupName:          "group1",
			policyName:         "Invalid Policy",
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.POLICY_IS_NOT_ATTACHED_TO_GROUP,
				Message: "Policy is not attached to group",
			},
			detachGroupPolicyErr: &api.Error{
				Code:    api.POLICY_IS_NOT_ATTACHED_TO_GROUP,
				Message: "Policy is not attached to group",
			},
		},
		"ErrorCaseUnauthorizedError": {
			org:                "org1",
			groupName:          "group1",
			policyName:         "policy1",
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			detachGroupPolicyErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseInvalidParameterErr": {
			org:                "org1",
			groupName:          "group1",
			policyName:         "policy1",
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid Parameter",
			},
			detachGroupPolicyErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid Parameter",
			},
		},
		"ErrorCaseUnknownApiError": {
			org:                "org1",
			groupName:          "group1",
			policyName:         "policy1",
			expectedStatusCode: http.StatusInternalServerError,
			detachGroupPolicyErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[DetachPolicyToGroupMethod][0] = test.detachGroupPolicyErr

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/groups/%v/policies/%v", test.org, test.groupName, test.policyName)
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		q := req.URL.Query()
		q.Add("Offset", test.offset)
		req.URL.RawQuery = q.Encode()

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameters
			assert.Equal(t, test.org, testApi.ArgsIn[DetachPolicyToGroupMethod][1], "Error in test case %v", n)
			assert.Equal(t, test.groupName, testApi.ArgsIn[DetachPolicyToGroupMethod][2], "Error in test case %v", n)
			assert.Equal(t, test.policyName, testApi.ArgsIn[DetachPolicyToGroupMethod][3], "Error in test case %v", n)
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

func TestWorkerHandler_HandleListAttachedGroupPolicies(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		// API method args
		filter       *api.Filter
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedResponse   ListAttachedGroupPoliciesResponse
		expectedError      api.Error
		// Manager Results
		getListAttachedGroupPoliciesResult []api.GroupPolicies
		totalGroupsResult                  int
		// Manager Errors
		getListAttachedGroupPoliciesErr error
	}{
		"OkCase": {
			filter: &api.Filter{
				Org:       "org1",
				GroupName: "group1",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: ListAttachedGroupPoliciesResponse{
				AttachedPolicies: []api.GroupPolicies{
					{
						Policy:   "policy1",
						CreateAt: now,
					},
					{
						Policy:   "policy2",
						CreateAt: now,
					},
				},
			},
			getListAttachedGroupPoliciesResult: []api.GroupPolicies{
				{
					Policy:   "policy1",
					CreateAt: now,
				},
				{
					Policy:   "policy2",
					CreateAt: now,
				},
			},
		},
		"ErrorCaseInvalidFilterParams": {
			filter: &api.Filter{
				PathPrefix: "",
				Offset:     -1,
			},
			ignoreArgsIn:       true,
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
		},
		"ErrorCaseGroupNotFoundErr": {
			filter: &api.Filter{
				Org:       "org1",
				GroupName: "group1",
			},
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.GROUP_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Group Not Found",
			},
			getListAttachedGroupPoliciesErr: &api.Error{
				Code:    api.GROUP_BY_ORG_AND_NAME_NOT_FOUND,
				Message: "Group Not Found",
			},
		},
		"ErrorCaseInvalidParameterErr": {
			filter: &api.Filter{
				Org:       "org1",
				GroupName: "group1",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid Parameter",
			},
			getListAttachedGroupPoliciesErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid Parameter",
			},
		},
		"ErrorCaseUnauthorizedError": {
			filter: &api.Filter{
				Org:       "org1",
				GroupName: "group1",
			},
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			getListAttachedGroupPoliciesErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseUnknownApiError": {
			filter:             testFilter,
			expectedStatusCode: http.StatusInternalServerError,
			getListAttachedGroupPoliciesErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[ListAttachedGroupPoliciesMethod][0] = test.getListAttachedGroupPoliciesResult
		testApi.ArgsOut[ListAttachedGroupPoliciesMethod][1] = test.totalGroupsResult
		testApi.ArgsOut[ListAttachedGroupPoliciesMethod][2] = test.getListAttachedGroupPoliciesErr

		url := fmt.Sprintf(server.URL+API_VERSION_1+"/organizations/%v/groups/%v/policies", test.filter.Org, test.filter.GroupName)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		addQueryParams(test.filter, req)

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameter
			filterData, ok := testApi.ArgsIn[ListAttachedGroupPoliciesMethod][1].(*api.Filter)
			if ok {
				// Check result
				assert.Equal(t, test.filter, filterData, "Error in test case %v", n)
			}
		}

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusOK:
			getGroupPoliciesResponse := ListAttachedGroupPoliciesResponse{}
			err = json.NewDecoder(res.Body).Decode(&getGroupPoliciesResponse)
			assert.Nil(t, err, "Error in test case %v", n)
			// Check result
			assert.Equal(t, test.expectedResponse, getGroupPoliciesResponse, "Error in test case %v", n)
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
