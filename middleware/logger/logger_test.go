package logger

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/middleware"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
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
	assert.Equal(t, http.StatusOK, res.StatusCode, "Error in test")

	// Check body
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(res.Body)
	assert.Nil(t, err, "Error in test")

	assert.Equal(t, string(buffer.Bytes()), testMessage)

	// Check context
	mc := new(middleware.MiddlewareContext)
	mw.GetInfo(req, mc)
	assert.Equal(t, new(middleware.MiddlewareContext), mc)

	// Check logger output
	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)
	assert.Equal(t, "", hook.LastEntry().Message)
	assert.Empty(t, hook.LastEntry().Data["requestID"])
	assert.Equal(t, http.MethodGet, hook.LastEntry().Data["httpMethod"])
	assert.Equal(t, req.RemoteAddr, hook.LastEntry().Data["httpRemoteAddress"])
	assert.Equal(t, req.RequestURI, hook.LastEntry().Data["httpURI"])

}
