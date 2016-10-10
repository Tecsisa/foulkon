package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

const (
	TEST_HEADER_NAME = "TestHeaderName"
)

// TestMiddleware that implements middleware interface
type TestMiddleware struct {
	HeaderValue string
}

func (tm *TestMiddleware) Action(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerValue := r.Header.Get(TEST_HEADER_NAME) + tm.HeaderValue
		r.Header.Set(TEST_HEADER_NAME, headerValue)
		next.ServeHTTP(w, r)
	})
}

func (tm *TestMiddleware) GetInfo(r *http.Request, mc *MiddlewareContext) {
	mc.XRequestId = r.Header.Get(TEST_HEADER_NAME)
}

func TestMiddlewareHandler_Handle(t *testing.T) {
	testMessage := "TestMessage"
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testMessage))
		w.WriteHeader(http.StatusOK)
	})
	testcases := map[string]struct {
		middlewares map[string]Middleware
	}{
		"OkTestCase": {
			middlewares: map[string]Middleware{
				REQUEST_LOGGER_MIDDLEWARE: &TestMiddleware{
					HeaderValue: REQUEST_LOGGER_MIDDLEWARE,
				},
				AUTHENTICATOR_MIDDLEWARE: &TestMiddleware{
					HeaderValue: AUTHENTICATOR_MIDDLEWARE,
				},
				XREQUESTID_MIDDLEWARE: &TestMiddleware{
					HeaderValue: XREQUESTID_MIDDLEWARE,
				},
			},
		},
	}

	for x, testcase := range testcases {
		mwh := getMiddlewareHandler(testcase.middlewares)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		mwh.Handle(testHandler).ServeHTTP(w, req)
		res := w.Result()

		// Check status code
		if res.StatusCode != http.StatusOK {
			t.Errorf("Test %v failed. Unexpected status code. (received/wanted) %v / %v", x, w.Code, http.StatusOK)
			continue
		}

		// Check body
		buffer := new(bytes.Buffer)
		if _, err := buffer.ReadFrom(res.Body); err != nil {
			t.Errorf("Test %v failed. Unexpected error reading response: %v.", x, err)
			continue
		}
		if diff := pretty.Compare(string(buffer.Bytes()), testMessage); diff != "" {
			t.Errorf("Test %v failed. Received different errors (received/wanted) %v", x, diff)
			continue
		}

		// Check Header
		expectedHeader := XREQUESTID_MIDDLEWARE + AUTHENTICATOR_MIDDLEWARE + REQUEST_LOGGER_MIDDLEWARE
		if diff := pretty.Compare(req.Header.Get(TEST_HEADER_NAME), expectedHeader); diff != "" {
			t.Errorf("Test %v failed. Received different header value (received/wanted) %v", x, diff)
			continue
		}

	}

}

func TestMiddlewareHandler_GetMiddlewareContext(t *testing.T) {
	testMessage := "TestMessage"
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testMessage))
		w.WriteHeader(http.StatusOK)
	})
	testcases := map[string]struct {
		middlewares     map[string]Middleware
		expectedContext *MiddlewareContext
	}{
		"OkTestCase": {
			middlewares: map[string]Middleware{
				REQUEST_LOGGER_MIDDLEWARE: &TestMiddleware{
					HeaderValue: REQUEST_LOGGER_MIDDLEWARE,
				},
				AUTHENTICATOR_MIDDLEWARE: &TestMiddleware{
					HeaderValue: AUTHENTICATOR_MIDDLEWARE,
				},
				XREQUESTID_MIDDLEWARE: &TestMiddleware{
					HeaderValue: XREQUESTID_MIDDLEWARE,
				},
			},
			expectedContext: &MiddlewareContext{
				XRequestId: XREQUESTID_MIDDLEWARE + AUTHENTICATOR_MIDDLEWARE + REQUEST_LOGGER_MIDDLEWARE,
			},
		},
	}

	for x, testcase := range testcases {
		mwh := getMiddlewareHandler(testcase.middlewares)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		mwh.Handle(testHandler).ServeHTTP(w, req)
		testmc := mwh.GetMiddlewareContext(req)

		// Check context
		if diff := pretty.Compare(testmc, testcase.expectedContext); diff != "" {
			t.Errorf("Test %v failed. Received different context (received/wanted) %v", x, diff)
			continue
		}

	}

}

// Private helper methods
func getMiddlewareHandler(middlewares map[string]Middleware) *MiddlewareHandler {
	return &MiddlewareHandler{Middlewares: middlewares}
}
