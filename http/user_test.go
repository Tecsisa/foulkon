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

func TestWorkerHandler_HandleGetUsers(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		pathPrefix string
		// Expected result
		expectedStatusCode int
		expectedResponse   GetUserExternalIDsResponse
		expectedError      api.Error
		// Manager Results
		getListUsersResult []string
		// Manager Errors
		getListUsersErr error
	}{
		"OkCase": {
			pathPrefix:         "myPath",
			expectedStatusCode: http.StatusOK,
			expectedResponse: GetUserExternalIDsResponse{
				ExternalIDs: []string{"userId1", "userId2"},
			},
			getListUsersResult: []string{"userId1", "userId2"},
		},
		"ErrorCaseUnauthorizedError": {
			pathPrefix:         "myPath",
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Error",
			},
			getListUsersErr: &api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Error",
			},
		},
		"ErrorCaseUnknownApiError": {
			expectedStatusCode: http.StatusInternalServerError,
			getListUsersErr: &api.Error{
				Code:    api.UNKNOWN_API_ERROR,
				Message: "Error",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {

		testApi.ArgsOut[GetListUsersMethod][0] = test.getListUsersResult
		testApi.ArgsOut[GetListUsersMethod][1] = test.getListUsersErr

		req, err := http.NewRequest(http.MethodGet, server.URL+USER_ROOT_URL, nil)
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
		if testApi.ArgsIn[GetListUsersMethod][1] != test.pathPrefix {
			t.Errorf("Test case %v. Received different PathPrefix (wanted:%v / received:%v)", n, test.pathPrefix, testApi.ArgsIn[GetListUsersMethod][1])
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

func TestWorkerHandler_HandlePostUsers(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		// API method args
		request *CreateUserRequest
		// Expected result
		expectedStatusCode int
		expectedResponse   CreateUserResponse
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
			expectedResponse: CreateUserResponse{
				User: &api.User{
					ID:         "UserID",
					ExternalID: "ExternalID",
					Path:       "Path",
					Urn:        "urn",
					CreateAt:   now,
				},
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

		req, err := http.NewRequest(http.MethodPost, server.URL+USER_ROOT_URL, body)
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
			createUserResponse := CreateUserResponse{}
			err = json.NewDecoder(res.Body).Decode(&createUserResponse)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(createUserResponse, test.expectedResponse); diff != "" {
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

func TestWorkerHandler_HandlePutUser(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		// API method args
		request *UpdateUserRequest
		// Expected result
		expectedStatusCode int
		expectedResponse   UpdateUserResponse
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
			expectedResponse: UpdateUserResponse{
				User: &api.User{
					ID:         "UserID",
					ExternalID: "ExternalID",
					Path:       "Path",
					Urn:        "urn",
					CreateAt:   now,
				},
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

		req, err := http.NewRequest(http.MethodPut, server.URL+USER_ROOT_URL+"/userid", body)
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
			updateUserResponse := UpdateUserResponse{}
			err = json.NewDecoder(res.Body).Decode(&updateUserResponse)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(updateUserResponse, test.expectedResponse); diff != "" {
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

func TestWorkerHandler_HandleGetUserId(t *testing.T) {
	now := time.Now()
	testcases := map[string]struct {
		// API method args
		externalID string
		// Expected result
		expectedStatusCode int
		expectedResponse   GetUserByIdResponse
		expectedError      api.Error
		// Manager Results
		getUserByExternalIdResult *api.User
		// Manager Errors
		getUserByExternalIdErr error
	}{
		"OkCase": {
			externalID:         "UserID",
			expectedStatusCode: http.StatusOK,
			expectedResponse: GetUserByIdResponse{
				User: &api.User{
					ID:         "UserID",
					ExternalID: "ExternalID",
					Path:       "Path",
					Urn:        "urn",
					CreateAt:   now,
				},
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
			externalID:         "UnauthorizedID",
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

		req, err := http.NewRequest(http.MethodGet, server.URL+USER_ROOT_URL+"/"+test.externalID, nil)
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
			getUserByIdResponse := GetUserByIdResponse{}
			err = json.NewDecoder(res.Body).Decode(&getUserByIdResponse)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}
			// Check result
			if diff := pretty.Compare(getUserByIdResponse, test.expectedResponse); diff != "" {
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
