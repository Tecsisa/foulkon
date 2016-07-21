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

func TestWorkerHandler_HandleAddUser(t *testing.T) {
	now := time.Now()
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
			},
			addUserResult: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
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
			if err != nil {
				t.Errorf("Test case %v. Unexpected marshalling api request %v", n, err)
				continue
			}
			body = bytes.NewBuffer(jsonObject)
		}
		if body == nil {
			body = bytes.NewBuffer([]byte{})
		}

		url := fmt.Sprintf(server.URL + USER_ROOT_URL)
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
			if testApi.ArgsIn[AddUserMethod][1] != test.request.ExternalID {
				t.Errorf("Test case %v. Received different ExternalID (wanted:%v / received:%v)", n, test.request.ExternalID, testApi.ArgsIn[AddUserMethod][1])
				continue
			}
			if testApi.ArgsIn[AddUserMethod][2] != test.request.Path {
				t.Errorf("Test case %v. Received different Path (wanted:%v / received:%v)", n, test.request.Path, testApi.ArgsIn[AddUserMethod][2])
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
			response := api.User{}
			err = json.NewDecoder(res.Body).Decode(&response)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(response, test.expectedResponse); diff != "" {
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

func TestWorkerHandler_HandleGetUserByExternalID(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		// API method args
		externalID string
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
			},
			getUserByExternalIdResult: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
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
		if testApi.ArgsIn[GetUserByExternalIdMethod][1] != test.externalID {
			t.Errorf("Test case %v. Received different ExternalID (wanted:%v / received:%v)", n, test.externalID, testApi.ArgsIn[GetUserByExternalIdMethod][1])
			continue
		}

		// check status code
		if test.expectedStatusCode != res.StatusCode {
			t.Errorf("Test case %v. Received different http status code (wanted:%v / received:%v)", n, test.expectedStatusCode, res.StatusCode)
			continue
		}

		switch res.StatusCode {
		case http.StatusOK:
			response := api.User{}
			err = json.NewDecoder(res.Body).Decode(&response)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(response, test.expectedResponse); diff != "" {
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

func TestWorkerHandler_HandleListUsers(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		pathPrefix string
		// Expected result
		expectedStatusCode int
		expectedResponse   GetUserExternalIDsResponse
		expectedError      api.Error
		// Manager Results
		getUserListResult []string
		// Manager Errors
		getUserListErr error
	}{
		"OkCase": {
			pathPrefix:         "myPath",
			expectedStatusCode: http.StatusOK,
			expectedResponse: GetUserExternalIDsResponse{
				ExternalIDs: []string{"userId1", "userId2"},
			},
			getUserListResult: []string{"userId1", "userId2"},
		},
		"ErrorCaseUnauthorizedError": {
			pathPrefix:         "myPath",
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
		"ErrorCaseUnknownApiError": {
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
		testApi.ArgsOut[ListUsersMethod][1] = test.getUserListErr

		url := fmt.Sprintf(server.URL + USER_ROOT_URL)
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

		// Check received parameter
		if testApi.ArgsIn[ListUsersMethod][1] != test.pathPrefix {
			t.Errorf("Test case %v. Received different PathPrefix (wanted:%v / received:%v)", n, test.pathPrefix, testApi.ArgsIn[ListUsersMethod][1])
			continue
		}

		// check status code
		if test.expectedStatusCode != res.StatusCode {
			t.Errorf("Test case %v. Received different http status code (wanted:%v / received:%v)", n, test.expectedStatusCode, res.StatusCode)
			continue
		}

		switch res.StatusCode {
		case http.StatusOK:
			getUserExternalIDsResponse := GetUserExternalIDsResponse{}
			err = json.NewDecoder(res.Body).Decode(&getUserExternalIDsResponse)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(getUserExternalIDsResponse, test.expectedResponse); diff != "" {
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

func TestWorkerHandler_HandleUpdateUser(t *testing.T) {
	now := time.Now()
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
			},
			updateUserResult: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
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
			if err != nil {
				t.Errorf("Test case %v. Unexpected marshalling api request %v", n, err)
				continue
			}
			body = bytes.NewBuffer(jsonObject)
		}
		if body == nil {
			body = bytes.NewBuffer([]byte{})
		}

		url := fmt.Sprintf(server.URL + USER_ROOT_URL + "/userid")
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
			if testApi.ArgsIn[UpdateUserMethod][1] != "userid" {
				t.Errorf("Test case %v. Received different ExternalID (wanted:%v / received:%v)", n, "userid", testApi.ArgsIn[UpdateUserMethod][1])
				continue
			}
			if testApi.ArgsIn[UpdateUserMethod][2] != test.request.Path {
				t.Errorf("Test case %v. Received different Path (wanted:%v / received:%v)", n, test.request.Path, testApi.ArgsIn[UpdateUserMethod][2])
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
			response := api.User{}
			err = json.NewDecoder(res.Body).Decode(&response)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(response, test.expectedResponse); diff != "" {
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

func TestWorkerHandler_HandleRemoveUser(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		externalID string
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
		if testApi.ArgsIn[RemoveUserMethod][1] != test.externalID {
			t.Errorf("Test case %v. Received different ExternalID (wanted:%v / received:%v)", n, test.externalID, testApi.ArgsIn[RemoveUserMethod][1])
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

func TestWorkerHandler_HandleListGroupsByUser(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		externalID string
		// Expected result
		expectedStatusCode int
		expectedResponse   GetGroupsByUserIdResponse
		expectedError      api.Error
		// Manager Results
		getGroupsByUserIdResult []api.GroupIdentity
		// Manager Errors
		getGroupsByUserIdErr error
	}{
		"OkCase": {
			externalID:         "UserID",
			expectedStatusCode: http.StatusOK,
			expectedResponse: GetGroupsByUserIdResponse{
				Groups: []api.GroupIdentity{
					{
						Org:  "org1",
						Name: "group1",
					},
					{
						Org:  "org2",
						Name: "group2",
					},
				},
			},
			getGroupsByUserIdResult: []api.GroupIdentity{
				{
					Org:  "org1",
					Name: "group1",
				},
				{
					Org:  "org2",
					Name: "group2",
				},
			},
		},
		"ErrorCaseUserNotExist": {
			externalID:         "UserID",
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
			externalID:         "InvalidID",
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
			externalID:         "UnauthorizedID",
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
			externalID:         "ExceptionID",
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
		testApi.ArgsOut[ListGroupsByUserMethod][1] = test.getGroupsByUserIdErr

		url := fmt.Sprintf(server.URL+USER_ROOT_URL+"/%v/groups", test.externalID)
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
		if testApi.ArgsIn[ListGroupsByUserMethod][1] != test.externalID {
			t.Errorf("Test case %v. Received different ExternalID (wanted:%v / received:%v)", n, test.externalID, testApi.ArgsIn[ListGroupsByUserMethod][1])
			continue
		}

		// check status code
		if test.expectedStatusCode != res.StatusCode {
			t.Errorf("Test case %v. Received different http status code (wanted:%v / received:%v)", n, test.expectedStatusCode, res.StatusCode)
			continue
		}

		switch res.StatusCode {
		case http.StatusOK:
			getGroupsByUserIdResponse := GetGroupsByUserIdResponse{}
			err = json.NewDecoder(res.Body).Decode(&getGroupsByUserIdResponse)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(getGroupsByUserIdResponse, test.expectedResponse); diff != "" {
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
