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

func TestWorkerHandler_HandleAddUser(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// API method args
		request *CreateUserRequest
		// Expected result
		expectedStatusCode int
		expectedResponse   *api.User
		expectedError      api.Error
		// Manager Results
		addUserResult *api.User
		// Manager Errors
		addUserErr error
	}{
		"OkCase": {
			request: &CreateUserRequest{
				ExternalID: "UserID",
				Path:       "Path",
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
				UpdateAt:   now,
			},
			addUserResult: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
				UpdateAt:   now,
			},
		},
		"ErrorCaseMalformedRequest": {
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "EOF",
			},
		},
		"ErrorCaseUserAlreadyExist": {
			request: &CreateUserRequest{
				ExternalID: "UserID",
				Path:       "Path",
			},
			expectedStatusCode: http.StatusConflict,
			expectedError: api.Error{
				Code:    api.USER_ALREADY_EXIST,
				Message: "User already exist",
			},
			addUserErr: &api.Error{
				Code:    api.USER_ALREADY_EXIST,
				Message: "User already exist",
			},
		},
		"ErrorCaseInvalidParameterError": {
			request: &CreateUserRequest{
				ExternalID: "UserID",
				Path:       "Path",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
			addUserErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
		},
		"ErrorCaseUnauthorizedResourcesError": {
			request: &CreateUserRequest{
				ExternalID: "UserID",
				Path:       "Path",
			},
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			addUserErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseUnknownApiError": {
			request: &CreateUserRequest{
				ExternalID: "UserID",
				Path:       "Path",
			},
			expectedStatusCode: http.StatusInternalServerError,
			addUserErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[AddUserMethod][0] = test.addUserResult
		testApi.ArgsOut[AddUserMethod][1] = test.addUserErr

		var body *bytes.Buffer
		if test.request != nil {
			jsonObject, err := json.Marshal(test.request)
			assert.Nil(t, err, "Error in test case %v", n)
			body = bytes.NewBuffer(jsonObject)
		}
		if body == nil {
			body = bytes.NewBuffer([]byte{})
		}

		url := fmt.Sprintf(server.URL + USER_ROOT_URL)
		req, err := http.NewRequest(http.MethodPost, url, body)
		assert.Nil(t, err, "Error in test case %v", n)

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if test.request != nil {
			// Check received parameters
			assert.Equal(t, test.request.ExternalID, testApi.ArgsIn[AddUserMethod][1], "Error in test case %v", n)
			assert.Equal(t, test.request.Path, testApi.ArgsIn[AddUserMethod][2], "Error in test case %v", n)
		}

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusCreated:
			response := &api.User{}
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
			// Check result
			assert.Equal(t, test.expectedError, apiError, "Error in test case %v", n)
		}
	}
}

func TestWorkerHandler_HandleGetUserByExternalID(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// API method args
		externalID   string
		offset       string
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedResponse   *api.User
		expectedError      api.Error
		// Manager Results
		getUserByExternalIdResult *api.User
		// Manager Errors
		getUserByExternalIdErr error
	}{
		"OkCase": {
			externalID:         "UserID",
			expectedStatusCode: http.StatusOK,
			expectedResponse: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
				UpdateAt:   now,
			},
			getUserByExternalIdResult: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
				UpdateAt:   now,
			},
		},
		"ErrorCaseInvalidRequest": {
			externalID:         "UserID",
			offset:             "-1",
			ignoreArgsIn:       true,
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
		},
		"ErrorCaseUserNotExist": {
			externalID:         "UserID",
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.USER_BY_EXTERNAL_ID_NOT_FOUND,
				Message: "User not exist",
			},
			getUserByExternalIdErr: &api.Error{
				Code:    api.USER_BY_EXTERNAL_ID_NOT_FOUND,
				Message: "User not exist",
			},
		},
		"ErrorCaseInvalidParameterError": {
			externalID:         "InvalidID",
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
			getUserByExternalIdErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
		},
		"ErrorCaseUnauthorizedResourcesError": {
			externalID:         "UnauthorizedID",
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			getUserByExternalIdErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseUnknownApiError": {
			externalID:         "ExceptionID",
			expectedStatusCode: http.StatusInternalServerError,
			getUserByExternalIdErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[GetUserByExternalIdMethod][0] = test.getUserByExternalIdResult
		testApi.ArgsOut[GetUserByExternalIdMethod][1] = test.getUserByExternalIdErr

		url := fmt.Sprintf(server.URL+USER_ROOT_URL+"/%v", test.externalID)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		q := req.URL.Query()
		q.Add("Offset", test.offset)
		req.URL.RawQuery = q.Encode()

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameters
			assert.Equal(t, test.externalID, testApi.ArgsIn[GetUserByExternalIdMethod][1], "Error in test case %v", n)
		}

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusOK:
			response := &api.User{}
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
			// Check result
			assert.Equal(t, test.expectedError, apiError, "Error in test case %v", n)
		}
	}
}

func TestWorkerHandler_HandleListUsers(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		filter       *api.Filter
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedResponse   GetUserExternalIDsResponse
		expectedError      api.Error
		// Manager Results
		getUserListResult []string
		totalResult       int
		// Manager Errors
		getUserListErr error
	}{
		"OkCase": {
			filter: &api.Filter{
				PathPrefix: "myPath",
				Offset:     0,
				Limit:      0,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: GetUserExternalIDsResponse{
				ExternalIDs: []string{"userId1", "userId2"},
				Offset:      0,
				Limit:       0,
				Total:       2,
			},
			getUserListResult: []string{"userId1", "userId2"},
			totalResult:       2,
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
		"ErrorCaseUnauthorizedError": {
			filter: &api.Filter{
				PathPrefix: "myPath",
			},
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Error",
			},
			getUserListErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Error",
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
			getUserListErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid Path",
			},
		},
		"ErrorCaseUnknownApiError": {
			filter:             testFilter,
			expectedStatusCode: http.StatusInternalServerError,
			getUserListErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[ListUsersMethod][0] = test.getUserListResult
		testApi.ArgsOut[ListUsersMethod][1] = test.totalResult
		testApi.ArgsOut[ListUsersMethod][2] = test.getUserListErr

		url := fmt.Sprintf(server.URL + USER_ROOT_URL)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		addQueryParams(test.filter, req)

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameter
			filterData, ok := testApi.ArgsIn[ListUsersMethod][1].(*api.Filter)
			if ok {
				// Check result
				assert.Equal(t, test.filter, filterData, "Error in test case %v", n)
			}
		}

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusOK:
			getUserExternalIDsResponse := GetUserExternalIDsResponse{}
			err = json.NewDecoder(res.Body).Decode(&getUserExternalIDsResponse)
			assert.Nil(t, err, "Error in test case %v", n)
			// Check result
			assert.Equal(t, test.expectedResponse, getUserExternalIDsResponse, "Error in test case %v", n)
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

func TestWorkerHandler_HandleUpdateUser(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// API method args
		request *UpdateUserRequest
		// Expected result
		expectedStatusCode int
		expectedResponse   *api.User
		expectedError      api.Error
		// Manager Results
		updateUserResult *api.User
		// Manager Errors
		updateUserErr error
	}{
		"OkCase": {
			request: &UpdateUserRequest{
				Path: "NewPath",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
				UpdateAt:   now,
			},
			updateUserResult: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
				UpdateAt:   now,
			},
		},
		"ErrorCaseMalformedRequest": {
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "EOF",
			},
		},
		"ErrorCaseUserNotExist": {
			request: &UpdateUserRequest{
				Path: "NewPath",
			},
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.USER_BY_EXTERNAL_ID_NOT_FOUND,
				Message: "User not exist",
			},
			updateUserErr: &api.Error{
				Code:    api.USER_BY_EXTERNAL_ID_NOT_FOUND,
				Message: "User not exist",
			},
		},
		"ErrorCaseInvalidParameterError": {
			request: &UpdateUserRequest{
				Path: "InvalidPath",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
			updateUserErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
		},
		"ErrorCaseUnauthorizedResourcesError": {
			request: &UpdateUserRequest{
				Path: "NewPath",
			},
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			updateUserErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseUnknownApiError": {
			request: &UpdateUserRequest{
				Path: "NewPath",
			},
			expectedStatusCode: http.StatusInternalServerError,
			updateUserErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[UpdateUserMethod][0] = test.updateUserResult
		testApi.ArgsOut[UpdateUserMethod][1] = test.updateUserErr

		var body *bytes.Buffer
		if test.request != nil {
			jsonObject, err := json.Marshal(test.request)
			assert.Nil(t, err, "Error in test case %v", n)
			body = bytes.NewBuffer(jsonObject)
		}
		if body == nil {
			body = bytes.NewBuffer([]byte{})
		}

		url := fmt.Sprintf(server.URL + USER_ROOT_URL + "/userid")
		req, err := http.NewRequest(http.MethodPut, url, body)
		assert.Nil(t, err, "Error in test case %v", n)

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if test.request != nil {
			// Check received parameters
			assert.Equal(t, "userid", testApi.ArgsIn[UpdateUserMethod][1], "Error in test case %v", n)
			assert.Equal(t, test.request.Path, testApi.ArgsIn[UpdateUserMethod][2], "Error in test case %v", n)
		}

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusOK:
			response := &api.User{}
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
			// Check result
			assert.Equal(t, test.expectedError, apiError, "Error in test case %v", n)
		}
	}
}

func TestWorkerHandler_HandleRemoveUser(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		externalID   string
		offset       string
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedError      api.Error
		// Manager Errors
		removeUserByIdErr error
	}{
		"OkCase": {
			externalID:         "UserID",
			expectedStatusCode: http.StatusNoContent,
		},
		"ErrorCaseInvalidRequest": {
			externalID:         "UserID",
			offset:             "-1",
			ignoreArgsIn:       true,
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
			removeUserByIdErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset -1",
			},
		},
		"ErrorCaseUserNotExist": {
			externalID:         "UserID",
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.USER_BY_EXTERNAL_ID_NOT_FOUND,
				Message: "User not exist",
			},
			removeUserByIdErr: &api.Error{
				Code:    api.USER_BY_EXTERNAL_ID_NOT_FOUND,
				Message: "User not exist",
			},
		},
		"ErrorCaseInvalidParameterError": {
			externalID:         "InvalidID",
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
			removeUserByIdErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
		},
		"ErrorCaseUnauthorizedResourcesError": {
			externalID:         "UnauthorizedID",
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			removeUserByIdErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseUnknownApiError": {
			externalID:         "ExceptionID",
			expectedStatusCode: http.StatusInternalServerError,
			removeUserByIdErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[RemoveUserMethod][0] = test.removeUserByIdErr

		url := fmt.Sprintf(server.URL+USER_ROOT_URL+"/%v", test.externalID)
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		q := req.URL.Query()
		q.Add("Offset", test.offset)
		req.URL.RawQuery = q.Encode()

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameters
			assert.Equal(t, test.externalID, testApi.ArgsIn[RemoveUserMethod][1], "Error in test case %v", n)
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
			// Check result
			assert.Equal(t, test.expectedError, apiError, "Error in test case %v", n)
		}
	}
}

func TestWorkerHandler_HandleListGroupsByUser(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// API method args
		filter       *api.Filter
		ignoreArgsIn bool
		// Expected result
		expectedStatusCode int
		expectedResponse   GetGroupsByUserIdResponse
		expectedError      api.Error
		// Manager Results
		getGroupsByUserIdResult []api.UserGroups
		totalGroupsResult       int
		// Manager Errors
		getGroupsByUserIdErr error
	}{
		"OkCase": {
			filter: &api.Filter{
				ExternalID: "UserID",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: GetGroupsByUserIdResponse{
				Groups: []api.UserGroups{
					{
						Org:      "org1",
						Name:     "group1",
						CreateAt: now,
					},
					{
						Org:      "org2",
						Name:     "group2",
						CreateAt: now,
					},
				},
				Offset: 0,
				Limit:  0,
				Total:  2,
			},
			getGroupsByUserIdResult: []api.UserGroups{
				{
					Org:      "org1",
					Name:     "group1",
					CreateAt: now,
				},
				{
					Org:      "org2",
					Name:     "group2",
					CreateAt: now,
				},
			},
			totalGroupsResult: 2,
		},
		"ErrorCaseInvalidFilterParams": {
			filter: &api.Filter{
				PathPrefix: "",
				ExternalID: "UserID",
				Limit:      -1,
			},
			ignoreArgsIn:       true,
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Limit -1",
			},
		},
		"ErrorCaseUserNotExist": {
			filter: &api.Filter{
				ExternalID: "UserID",
			},
			expectedStatusCode: http.StatusNotFound,
			expectedError: api.Error{
				Code:    api.USER_BY_EXTERNAL_ID_NOT_FOUND,
				Message: "User not exist",
			},
			getGroupsByUserIdErr: &api.Error{
				Code:    api.USER_BY_EXTERNAL_ID_NOT_FOUND,
				Message: "User not exist",
			},
		},
		"ErrorCaseInvalidParameterError": {
			filter: &api.Filter{
				ExternalID: "InvalidID",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
			getGroupsByUserIdErr: &api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter",
			},
		},
		"ErrorCaseUnauthorizedResourcesError": {
			filter: &api.Filter{
				ExternalID: "UnauthorizedID",
			},
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
			getGroupsByUserIdErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized",
			},
		},
		"ErrorCaseUnknownApiError": {
			filter: &api.Filter{
				ExternalID: "ExceptionID",
			},
			expectedStatusCode: http.StatusInternalServerError,
			getGroupsByUserIdErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[ListGroupsByUserMethod][0] = test.getGroupsByUserIdResult
		testApi.ArgsOut[ListGroupsByUserMethod][1] = test.totalGroupsResult
		testApi.ArgsOut[ListGroupsByUserMethod][2] = test.getGroupsByUserIdErr

		url := fmt.Sprintf(server.URL+USER_ROOT_URL+"/%v/groups", test.filter.ExternalID)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, err, "Error in test case %v", n)

		addQueryParams(test.filter, req)

		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		if !test.ignoreArgsIn {
			// Check received parameters
			filterData, ok := testApi.ArgsIn[ListGroupsByUserMethod][1].(*api.Filter)
			if ok {
				// Check result
				assert.Equal(t, test.filter, filterData, "Error in test case %v", n)
			}
		}

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusOK:
			getGroupsByUserIdResponse := GetGroupsByUserIdResponse{}
			err = json.NewDecoder(res.Body).Decode(&getGroupsByUserIdResponse)
			assert.Nil(t, err, "Error in test case %v", n)
			// Check result
			assert.Equal(t, test.expectedResponse, getGroupsByUserIdResponse, "Error in test case %v", n)
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
