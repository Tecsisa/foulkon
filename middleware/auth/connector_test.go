package auth

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sirupsen/logrus/hooks/test"
	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/middleware"
	"github.com/stretchr/testify/assert"
)

// Aux connector
type TestConnector struct {
	userID          string
	unauthenticated bool
}

func (tc *TestConnector) Authenticate(h http.Handler) http.Handler {
	if tc.unauthenticated {
		// Reset value
		tc.unauthenticated = false
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		})
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set(middleware.USER_ID_HEADER, tc.userID)
		h.ServeHTTP(w, r)
	})
}

func (tc TestConnector) RetrieveUserID(r http.Request) string {
	return tc.userID
}

func TestAuthenticatorMiddleware_Action(t *testing.T) {
	testMessage := "TestMessage"
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testMessage))
		w.WriteHeader(http.StatusOK)
	})
	// Create logger
	testLogger, hook := test.NewNullLogger()
	api.Log = testLogger
	testcases := map[string]struct {
		// Middleware args
		userID             string
		password           string
		unauthenticated    bool
		admin              bool
		expectedLog        string
		expectedStatusCode int
		testConnectorNull  bool
	}{
		"OkCase": {
			userID:             "UserId",
			unauthenticated:    false,
			expectedStatusCode: http.StatusOK,
		},
		"OkCaseAdmin": {
			userID:             "admin",
			password:           "admin",
			unauthenticated:    false,
			expectedStatusCode: http.StatusOK,
			admin:              true,
		},
		"OkCaseInvalidAdmin": {
			userID:             "admin",
			password:           "fail",
			unauthenticated:    false,
			expectedStatusCode: http.StatusOK,
			admin:              true,
			expectedLog:        "Trying to connect as admin, admin user/password invalid, delegating to connector...",
		},
		"OkCaseUnautenticated": {
			userID:             "UserId",
			unauthenticated:    true,
			expectedStatusCode: http.StatusUnauthorized,
		},
		"OkNoAuthAuthenticationMethod": {
			userID:             "UserId",
			unauthenticated:    true,
			expectedStatusCode: http.StatusUnauthorized,
			testConnectorNull:  true,
		},
	}

	for n, testcase := range testcases {
		var mw *AuthenticatorMiddleware
		if testcase.testConnectorNull {
			mw = NewAuthenticatorMiddleware(nil, "admin", "admin")
		} else {
			mw = NewAuthenticatorMiddleware(&TestConnector{userID: testcase.userID, unauthenticated: testcase.unauthenticated}, "admin", "admin")
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		if testcase.admin {
			req.SetBasicAuth(testcase.userID, testcase.password)
		}
		w := httptest.NewRecorder()
		mw.Action(testHandler).ServeHTTP(w, req)
		res := w.Result()
		// Check status code
		assert.Equal(t, testcase.expectedStatusCode, res.StatusCode, "Error in test case %v", n)

		// Check body
		if res.StatusCode == http.StatusOK {
			// Check logger
			if testcase.expectedLog != "" {
				assert.Equal(t, testcase.expectedLog, hook.LastEntry().Message, "Error in test case %v", n)
			}

			buffer := new(bytes.Buffer)
			_, err := buffer.ReadFrom(res.Body)
			assert.Nil(t, err, "Error in test case %v", n)

			assert.Equal(t, string(buffer.Bytes()), testMessage, "Error in test case %v", n)
			// Check Header
			userID := req.Header.Get(middleware.USER_ID_HEADER)
			assert.Equal(t, testcase.userID, userID, "Error in test case %v", n)
		}
	}
}

func TestAuthenticatorMiddleware_GetInfo(t *testing.T) {
	testMessage := "TestMessage"
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testMessage))
		w.WriteHeader(http.StatusOK)
	})
	testcases := map[string]struct {
		// Middleware args
		userID             string
		password           string
		unauthenticated    bool
		admin              bool
		expectedStatusCode int
	}{
		"OkCase": {
			userID:             "UserId",
			unauthenticated:    false,
			expectedStatusCode: http.StatusOK,
		},
		"OkCaseAdmin": {
			userID:             "admin",
			password:           "admin",
			unauthenticated:    false,
			expectedStatusCode: http.StatusOK,
			admin:              true,
		},
	}

	for n, testcase := range testcases {
		mw := NewAuthenticatorMiddleware(&TestConnector{userID: testcase.userID, unauthenticated: testcase.unauthenticated}, "admin", "admin")
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		if testcase.admin {
			req.SetBasicAuth(testcase.userID, testcase.password)
		}
		w := httptest.NewRecorder()
		mw.Action(testHandler).ServeHTTP(w, req)
		mc := new(middleware.MiddlewareContext)
		mw.GetInfo(req, mc)

		// Check user id
		assert.Equal(t, testcase.userID, mc.UserId, "Error in test case %v", n)
		// Check admin privilege
		assert.Equal(t, testcase.admin, mc.Admin, "Error in test case %v", n)
	}
}
