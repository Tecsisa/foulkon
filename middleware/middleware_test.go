package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, http.StatusOK, res.StatusCode, "Error in test case %v", x)

		// Check body
		buffer := new(bytes.Buffer)
		_, err := buffer.ReadFrom(res.Body)
		assert.Nil(t, err, "Error in test case %v", x)

		assert.Equal(t, string(buffer.Bytes()), testMessage)

		// Check Header
		expectedHeader := XREQUESTID_MIDDLEWARE + AUTHENTICATOR_MIDDLEWARE + REQUEST_LOGGER_MIDDLEWARE
		assert.Equal(t, expectedHeader, req.Header.Get(TEST_HEADER_NAME))
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
		assert.Equal(t, testcase.expectedContext, testmc, "Error in test case %v", x)
	}

}

// Private helper methods
func getMiddlewareHandler(middlewares map[string]Middleware) *MiddlewareHandler {
	return &MiddlewareHandler{Middlewares: middlewares}
}
