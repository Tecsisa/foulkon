package http

import (
	"encoding/json"
	"github.com/tecsisa/authorizr/api"
	"net/http"
	"reflect"
	"testing"
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
			if !reflect.DeepEqual(getUserExternalIDsResponse, test.expectedResponse) {
				t.Errorf("Test %v failed. Received different responses (wanted:%v / received:%v)",
					n, test.expectedResponse, getUserExternalIDsResponse)
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
			if !reflect.DeepEqual(apiError, test.expectedError) {
				t.Errorf("Test %v failed. Received different error response (wanted:%v / received:%v)",
					n, test.expectedError, apiError)
				continue
			}

		}

	}
}
