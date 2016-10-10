package logger

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/Tecsisa/foulkon/middleware"
	"github.com/kylelemons/godebug/pretty"
)

func TestRequestLoggerMiddleware_Action(t *testing.T) {
	testMessage := "TestMessage"
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testMessage))
		w.WriteHeader(http.StatusOK)
	})

	// Create logger
	logBuffer := bytes.NewBufferString("")
	log := &logrus.Logger{
		Out:       logBuffer,
		Formatter: &logrus.TextFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}

	mw := NewRequestLoggerMiddleware(log)
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
	testLogBuffer := bytes.NewBufferString("")
	testLogger := &logrus.Logger{
		Out:       testLogBuffer,
		Formatter: &logrus.TextFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}
	testLogger.WithFields(logrus.Fields{
		"requestID": req.Header.Get(middleware.REQUEST_ID_HEADER),
		"method":    req.Method,
		"URI":       req.RequestURI,
		"address":   req.RemoteAddr,
		"user":      req.Header.Get(middleware.USER_ID_HEADER),
	}).Info("")
	if diff := pretty.Compare(string(logBuffer.Bytes()), string(testLogBuffer.Bytes())); diff != "" {
		t.Errorf("Test failed. Received different messages (received/wanted) %v", diff)
		return
	}

}
