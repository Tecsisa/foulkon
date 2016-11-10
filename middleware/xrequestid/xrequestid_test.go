package xrequestid

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tecsisa/foulkon/middleware"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestXRequestIdMiddleware_Action(t *testing.T) {
	testMessage := "TestMessage"
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testMessage))
		w.WriteHeader(http.StatusOK)
	})

	mw := NewXRequestIdMiddleware()
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

	// Check Header
	header := req.Header.Get(middleware.REQUEST_ID_HEADER)

	_, err = uuid.FromString(header)
	assert.Nil(t, err, "Error in test")
}

func TestXRequestIdMiddleware_GetInfo(t *testing.T) {
	mw := NewXRequestIdMiddleware()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	testuuid := uuid.NewV4().String()
	req.Header.Set(middleware.REQUEST_ID_HEADER, testuuid)
	mc := new(middleware.MiddlewareContext)
	mw.GetInfo(req, mc)
	// Check request id value from context
	assert.Equal(t, mc.XRequestId, testuuid, "Error in test")
}
