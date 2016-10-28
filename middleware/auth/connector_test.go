package auth

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sirupsen/logrus/hooks/test"
	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/middleware"
	"github.com/kylelemons/godebug/pretty"
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
	}

	for n, testcase := range testcases {
		mw := NewAuthenticatorMiddleware(&TestConnector{userID: testcase.userID, unauthenticated: testcase.unauthenticated}, "admin", "admin")
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		if testcase.admin {
			req.SetBasicAuth(testcase.userID, testcase.password)
		}
		w := httptest.NewRecorder()
		mw.Action(testHandler).ServeHTTP(w, req)
		res := w.Result()
		// Check status code
		if res.StatusCode != testcase.expectedStatusCode {
			t.Errorf("Test %v failed. Unexpected status code. (received/wanted) %v / %v", n, w.Code, testcase.expectedStatusCode)
			continue
		}
		// Check body
		if res.StatusCode == http.StatusOK {
			// Check logger
			if testcase.expectedLog != "" {
				if diff := pretty.Compare(testcase.expectedLog, hook.LastEntry().Message); diff != "" {
					t.Errorf("Test %v failed. Received different logs (received/wanted) %v", n, diff)
					return
				}
			}

			buffer := new(bytes.Buffer)
			if _, err := buffer.ReadFrom(res.Body); err != nil {
				t.Errorf("Test %v failed. Unexpected error reading response: %v.", n, err)
				return
			}
			if diff := pretty.Compare(string(buffer.Bytes()), testMessage); diff != "" {
				t.Errorf("Test %v failed. Received different errors (received/wanted) %v", n, diff)
				return
			}

			// Check Header
			userID := req.Header.Get(middleware.USER_ID_HEADER)
			if diff := pretty.Compare(userID, testcase.userID); diff != "" {
				t.Errorf("Test %v failed. Received different users (received/wanted) %v", n, diff)
				return
			}
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
		if mc.UserId != testcase.userID {
			t.Errorf("Test %v failed. Different user ids received: (received/wanted) %v/%v.", n, mc.UserId, testcase.userID)
		}

		// Check admin privilege
		if mc.Admin != testcase.admin {
			t.Errorf("Test %v failed. Different user admin privilege: (received/wanted) %v/%v.", n, mc.Admin, testcase.admin)
		}

	}
}
