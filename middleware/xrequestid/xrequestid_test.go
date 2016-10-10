package xrequestid

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tecsisa/foulkon/middleware"
	"github.com/kylelemons/godebug/pretty"
	"github.com/satori/go.uuid"
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

	// Check Header
	header := req.Header.Get(middleware.REQUEST_ID_HEADER)

	_, err := uuid.FromString(header)
	if err != nil {
		t.Errorf("Test failed. Invalid uuid received from header: %v", header)
		return
	}

}

func TestXRequestIdMiddleware_GetInfo(t *testing.T) {
	mw := NewXRequestIdMiddleware()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	testuuid := uuid.NewV4().String()
	req.Header.Set(middleware.REQUEST_ID_HEADER, testuuid)
	mc := new(middleware.MiddlewareContext)
	mw.GetInfo(req, mc)
	// Check request id value from context
	if mc.XRequestId != testuuid {
		t.Errorf("Test failed. Received differents ids. (received/wanted) %v / %v", mc.XRequestId, testuuid)
		return
	}
}
