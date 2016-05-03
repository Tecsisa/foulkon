package openid

import (
	"fmt"
	"net/http"
	"testing"
)

// Data used for negative tests of GetIdTokenAuthorizationHeader.
var badHeaders = []struct {
	header     string              // The wrong header.
	errorCode  ValidationErrorCode // The expected error code.
	httpStatus int                 // The expected http status code.
}{
	{"", ValidationErrorAuthorizationHeaderNotFound, http.StatusBadRequest},
	{"token", ValidationErrorAuthorizationHeaderWrongFormat, http.StatusBadRequest},
	{"token token token", ValidationErrorAuthorizationHeaderWrongFormat, http.StatusBadRequest},
	{"scheme token", ValidationErrorAuthorizationHeaderWrongSchemeName, http.StatusBadRequest},
	{"bearer token", ValidationErrorAuthorizationHeaderWrongSchemeName, http.StatusBadRequest},
	{"Bearer token token", ValidationErrorAuthorizationHeaderWrongFormat, http.StatusBadRequest},
}

// createRequest creates a request with the given string(headerContent) as the
// http Authorization header and returns that request.
func createRequest(headerContent string) *http.Request {
	r := http.Request{}
	r.Header = http.Header(map[string][]string{})
	r.Header.Set("Authorization", headerContent)
	return &r
}

// expectError validates whether the provided error(e) has
// an error code(c)
func expectError(t *testing.T, e error, headerContent string, errorCode ValidationErrorCode, httpStatus int) {
	if ve, ok := e.(*ValidationError); ok {
		if ve.Code != errorCode {
			t.Errorf("For header %v. Expected error code %v, got %v", headerContent, errorCode, ve.Code)
		}
		if ve.HTTPStatus != httpStatus {
			t.Errorf("For header %v. Expected http status %v, got %v", headerContent, httpStatus, ve.HTTPStatus)
		}
	} else {
		t.Errorf("For header %v. Expected error type 'ValidationError', got %T", headerContent, e)
	}
}

// Tests getIdTokenAuthorizationHeader providing an Authorization header with unexpected content.
func Test_getIDTokenAuthorizationHeader_WrongHeaderContent(t *testing.T) {
	for _, tt := range badHeaders {

		_, err := getIDTokenAuthorizationHeader(createRequest(tt.header))
		expectError(t, err, tt.header, tt.errorCode, tt.httpStatus)
	}
}

// Tests getIdTokenAuthorizationHeader providing a request without Authorization header.
func Test_getIDTokenAuthorizationHeader_NoHeader(t *testing.T) {
	_, err := getIDTokenAuthorizationHeader(&http.Request{})

	expectError(t, err, "No Authorization Header", ValidationErrorAuthorizationHeaderNotFound, http.StatusBadRequest)
}

// Tests getIdTokenAuthorizationHeader providing an Authorization header with expected format.
func Test_getIDTokenAuthorizationHeader_CorrectHeaderContent(t *testing.T) {
	et := "token"
	hc := fmt.Sprintf("Bearer %v", et)
	rt, err := getIDTokenAuthorizationHeader(createRequest(hc))

	if err != nil {
		t.Errorf("The header content %v is valid. Unexpected error", hc)
	}

	if rt != et {
		t.Errorf("Expected result %v, got %v", et, rt)
	}
}
