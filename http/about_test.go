package http

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Tecsisa/foulkon/api"
	"github.com/stretchr/testify/assert"
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
		assert.Nil(t, err, "Error in test case %v", n)

		req.SetBasicAuth(test.adminUser, test.adminPassword)
		res, err := client.Do(req)
		assert.Nil(t, err, "Error in test case %v", n)

		// check status code
		assert.Equal(t, test.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		switch res.StatusCode {
		case http.StatusOK:
			response := Config{}
			err = json.NewDecoder(res.Body).Decode(&response)
			assert.Nil(t, err, "Error in test case %v", n)
			// Check result
			assert.Equal(t, test.expectedConfig, response, "Error in test case %v", n)
		case http.StatusUnauthorized:
			apiError := api.Error{}
			// Check error
			err = json.NewDecoder(res.Body).Decode(&apiError)
			assert.Nil(t, err, "Error in test case %v", n)
			assert.Equal(t, test.expectedError, apiError, "Error in test case %v", n)
		}
	}
}
