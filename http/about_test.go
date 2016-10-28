package http

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Tecsisa/foulkon/api"
	"github.com/kylelemons/godebug/pretty"
)

func TestWorkerHandler_HandleGetCurrentConfig(t *testing.T) {
	testcases := map[string]struct {
		adminUser     string
		adminPassword string

		badRequest         string
		expectedStatusCode int
		expectedConfig     Config
		expectedError      error
	}{
		"OKCase": {
			adminUser:          "admin",
			adminPassword:      "admin",
			expectedStatusCode: http.StatusOK,
			expectedConfig: Config{
				Logger: LoggerConfig{
					Type:          "test",
					Level:         "test",
					FileDirectory: "test",
				},
				Database: DatabaseConfig{
					Type:         "test",
					IdleConns:    0,
					MaxOpenConns: 0,
					ConnTtl:      0,
				},
				AuthConnector: AuthConnectorConfig{
					Type:   "test",
					Issuer: "test",
				},
				Version: "test",
			},
		},
		"ErrorCaseBadRequest": {
			adminUser:          "admin",
			adminPassword:      "admin",
			badRequest:         "?Offset=L",
			expectedStatusCode: http.StatusBadRequest,
			expectedError: api.Error{
				Code:    api.INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: Offset L",
			},
		},
		"ErrorCaseInvalidAdmin": {
			adminUser:          "admin",
			adminPassword:      "fail",
			expectedStatusCode: http.StatusForbidden,
			expectedError: api.Error{
				Code:    api.UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Unauthorized, user is not admin",
			},
		},
	}

	client := http.DefaultClient

	for n, test := range testcases {
		req, err := http.NewRequest(http.MethodGet, server.URL+"/about"+test.badRequest, nil)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error creating http request %v", n, err)
			continue
		}

		req.SetBasicAuth(test.adminUser, test.adminPassword)
		res, err := client.Do(req)
		if err != nil {
			t.Errorf("Test case %v. Unexpected error calling server %v", n, err)
			continue
		}

		// check status code
		if test.expectedStatusCode != res.StatusCode {
			t.Errorf("Test case %v. Received different http status code (wanted:%v / received:%v)", n, test.expectedStatusCode, res.StatusCode)
			continue
		}

		switch res.StatusCode {
		case http.StatusOK:
			response := Config{}
			err = json.NewDecoder(res.Body).Decode(&response)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing response %v", n, err)
				continue
			}

			// Check result
			if diff := pretty.Compare(response, test.expectedConfig); diff != "" {
				t.Errorf("Test %v failed. Received different response (received/wanted) %v", n, diff)
				continue
			}
		case http.StatusUnauthorized:
			apiError := api.Error{}
			// Check error
			err = json.NewDecoder(res.Body).Decode(&apiError)
			if err != nil {
				t.Errorf("Test case %v. Unexpected error parsing error response %v", n, err)
				continue
			}
			if diff := pretty.Compare(apiError, test.expectedError); diff != "" {
				t.Errorf("Test %v failed. Received different error response (received/wanted) %v",
					n, diff)
				continue
			}
		}
	}
}
