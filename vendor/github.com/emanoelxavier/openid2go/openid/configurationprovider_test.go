package openid

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
)

type testBody struct {
	io.Reader
}

func (testBody) Close() error { return nil }

func Test_getConfiguration_UsesCorrectUrl(t *testing.T) {
	c := NewHTTPClientMock(t)
	configurationProvider := httpConfigurationProvider{getConfig: c.httpGet}

	issuer := "https://test"
	configSuffix := "/.well-known/openid-configuration"
	go func() {
		c.assertHttpGet(issuer+configSuffix, nil, errors.New("Read configuration error"))
		c.close()
	}()

	_, e := configurationProvider.getConfiguration(issuer)

	if e == nil {
		t.Error("An error was expected but not returned")
	}

	c.assertDone()
}

func Test_getConfiguration_WhenGetReturnsError(t *testing.T) {
	c := NewHTTPClientMock(t)
	configurationProvider := httpConfigurationProvider{getConfig: c.httpGet}

	readError := errors.New("Read configuration error")
	go func() {
		c.assertHttpGet(anything, nil, readError)
		c.close()
	}()

	_, e := configurationProvider.getConfiguration("issuer")

	expectValidationError(t, e, ValidationErrorGetOpenIdConfigurationFailure, http.StatusUnauthorized, readError)

	c.assertDone()
}

func Test_getConfiguration_WhenGetSucceeds(t *testing.T) {
	c := NewHTTPClientMock(t)
	configurationProvider := httpConfigurationProvider{c.httpGet, c.decodeResponse}

	respBody := "openid configuration"
	resp := &http.Response{Body: testBody{bytes.NewBufferString(respBody)}}

	go func() {
		c.assertHttpGet(anything, resp, nil)
		c.assertDecodeResponse(respBody, nil, nil)
		c.close()
	}()

	_, e := configurationProvider.getConfiguration(anything)

	if e != nil {
		t.Error("An error was returned but not expected", e)
	}

	c.assertDone()
}

func Test_getConfiguration_WhenDecodeResponseReturnsError(t *testing.T) {
	c := NewHTTPClientMock(t)
	configurationProvider := httpConfigurationProvider{c.httpGet, c.decodeResponse}
	decodeError := errors.New("Decode configuration error")
	respBody := "openid configuration"
	resp := &http.Response{Body: testBody{bytes.NewBufferString(respBody)}}

	go func() {
		c.assertHttpGet(anything, resp, nil)
		c.assertDecodeResponse(anything, nil, decodeError)
		c.close()
	}()

	_, e := configurationProvider.getConfiguration(anything)

	expectValidationError(t, e, ValidationErrorDecodeOpenIdConfigurationFailure, http.StatusUnauthorized, decodeError)

	c.assertDone()
}

func Test_getConfiguration_WhenDecodeResponseSucceeds(t *testing.T) {
	c := NewHTTPClientMock(t)
	configurationProvider := httpConfigurationProvider{c.httpGet, c.decodeResponse}
	config := &configuration{"testissuer", "https://testissuer/jwk"}
	respBody := "openid configuration"
	resp := &http.Response{Body: testBody{bytes.NewBufferString(respBody)}}

	go func() {
		c.assertHttpGet(anything, resp, nil)
		c.assertDecodeResponse(anything, config, nil)
		c.close()
	}()

	rc, e := configurationProvider.getConfiguration(anything)

	if e != nil {
		t.Error("An error was returned but not expected", e)
	}

	if rc.Issuer != config.Issuer {
		t.Error("Expected issuer", config.Issuer, "but was", rc.Issuer)
	}

	if rc.JwksUri != config.JwksUri {
		t.Error("Expected jwks uri", config.JwksUri, "but was", rc.JwksUri)
	}

	c.assertDone()
}

func expectValidationError(t *testing.T, e error, vec ValidationErrorCode, status int, inner error) {
	if e == nil {
		t.Error("An error was expected but not returned")
	}

	if ve, ok := e.(*ValidationError); ok {
		if ve.Code != vec {
			t.Error("Expected error code", vec, "but was", ve.Code)
		}
		if ve.HTTPStatus != status {
			t.Error("Expected HTTP status", status, "but was", ve.HTTPStatus)
		}
		if inner != nil && ve.Err.Error() != inner.Error() {
			t.Error("Expected inner error", inner.Error(), ",but was", ve.Err.Error())
		}
	} else {
		t.Errorf("Expected error type '*ValidationError' but was %T", e)
	}
}
