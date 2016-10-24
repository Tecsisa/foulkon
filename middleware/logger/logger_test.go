package logger

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/test"
	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/middleware"
	"github.com/kylelemons/godebug/pretty"
	"github.com/stretchr/testify/assert"
)

func TestRequestLoggerMiddleware_Action(t *testing.T) {
	testMessage := "TestMessage"
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testMessage))
		w.WriteHeader(http.StatusOK)
	})

	// Create logger
	testLogger, hook := test.NewNullLogger()
	api.Log = testLogger

	mw := NewRequestLoggerMiddleware()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mw.Action(testHandler).ServeHTTP(w, req)
	res := w.Result()
	// Check status code
	if res.StatusCode != http.StatusOK {
		t.Errorf("Test failed. Unexpected status code. (received/wanted) %v / %v", w.Code, http.StatusOK)
		return
	}
	// Check body
	buffer := new(bytes.Buffer)
	if _, err := buffer.ReadFrom(res.Body); err != nil {
		t.Errorf("Test failed. Unexpected error reading response: %v.", err)
		return
	}
	if diff := pretty.Compare(string(buffer.Bytes()), testMessage); diff != "" {
		t.Errorf("Test failed. Received different errors (received/wanted) %v", diff)
		return
	}

	// Check context
	mc := new(middleware.MiddlewareContext)
	mw.GetInfo(req, mc)
	if diff := pretty.Compare(mc, new(middleware.MiddlewareContext)); diff != "" {
		t.Errorf("Test failed. Received different contexts (received/wanted) %v", diff)
		return
	}

	// Check logger output
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)
	assert.Equal(t, "", hook.LastEntry().Message)
	assert.Empty(t, hook.LastEntry().Data["requestID"])
	assert.Equal(t, http.MethodGet, hook.LastEntry().Data["httpMethod"])
	assert.Equal(t, req.RemoteAddr, hook.LastEntry().Data["httpRemoteAddress"])
	assert.Equal(t, req.RequestURI, hook.LastEntry().Data["httpURI"])

}
