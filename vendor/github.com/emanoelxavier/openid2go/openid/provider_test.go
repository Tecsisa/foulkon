package openid

import "testing"

func Test_validateProviders_EmptyProviderList(t *testing.T) {
	var ps providers
	se := ps.validate()
	expectSetupError(t, se, SetupErrorEmptyProviderCollection)

	ps = make([]Provider, 0)
	se = ps.validate()
	expectSetupError(t, se, SetupErrorEmptyProviderCollection)

}

func Test_validateProvider_EmptyIssuer(t *testing.T) {
	p := Provider{}
	se := p.validate()
	expectSetupError(t, se, SetupErrorInvalidIssuer)
}

func Test_validateProvider_EmptyClientIDs(t *testing.T) {
	p := Provider{Issuer: "https://test"}
	se := p.validate()
	expectSetupError(t, se, SetupErrorInvalidClientIDs)
}

func Test_validateProvider_ValidProvider(t *testing.T) {
	p := Provider{Issuer: "https://test", ClientIDs: []string{"clientID"}}
	se := p.validate()

	if se != nil {
		t.Error("An error was returned but not expected", se)
	}
}

func Test_validateProviders_OneInvalidProvider(t *testing.T) {
	p := Provider{Issuer: "https://test", ClientIDs: []string{"clientID"}}
	ps := []Provider{p, Provider{}}

	se := providers(ps).validate()
	expectSetupError(t, se, SetupErrorInvalidIssuer)
}

func Test_validateProviders_AllValidProviders(t *testing.T) {
	p := Provider{Issuer: "https://test", ClientIDs: []string{"clientID"}}
	ps := []Provider{p, p}

	se := providers(ps).validate()

	if se != nil {
		t.Error("An error was returned but not expected", se)
	}
}

func expectSetupError(t *testing.T, e error, sec SetupErrorCode) {
	if e == nil {
		t.Error("An error was expected but not returned")
	}

	if se, ok := e.(*SetupError); ok {
		if se.Code != sec {
			t.Error("Expected error code", sec, "but was", se.Code)
		}
	} else {
		t.Errorf("Expected error type '*SetupError' but was %T", e)
	}
}
