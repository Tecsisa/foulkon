package openid

import (
	"crypto/rsa"
	"errors"
	"net/http"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

func Test_getSigningKey_WhenGetProvidersReturnsError(t *testing.T) {
	pm, _, _, _, tv := createIDTokenValidator(t)

	ee := errors.New("Error getting providers.")

	go func() {
		pm.assertGetProviders(nil, ee)
		pm.close()
	}()

	sk, err := tv.getSigningKey(nil)

	if sk != nil {
		t.Error("The returned signing key should be nil.")
	}

	if err == nil {
		t.Fatal("An error was expected but not returned.")
	}

	if err.Error() != ee.Error() {
		t.Error("Expected error", ee, ", but got", err)
	}

	pm.assertDone()
}

func Test_getSigningKey_WhenGetProvidersReturnsEmptyCollection(t *testing.T) {
	pm, _, _, _, tv := createIDTokenValidator(t)

	go func() {
		pm.assertGetProviders(nil, nil)
		pm.assertGetProviders([]Provider{}, nil)
		pm.close()
	}()

	_, err := tv.getSigningKey(nil)
	expectSetupError(t, err, SetupErrorEmptyProviderCollection)

	_, err = tv.getSigningKey(nil)
	expectSetupError(t, err, SetupErrorEmptyProviderCollection)

	pm.assertDone()

}

func Test_getSigningKey_UsingTokenWithInvalidIssuerType(t *testing.T) {
	pm, _, _, _, tv := createIDTokenValidator(t)

	go func() {
		pm.assertGetProviders([]Provider{{Issuer: "https://issuer", ClientIDs: []string{"client"}}}, nil)
		pm.close()
	}()

	jt := jwt.New(jwt.SigningMethodRS256)
	jt.Claims.(jwt.MapClaims)["iss"] = 0 // The expected issuer type is string, not int.
	sk, err := tv.getSigningKey(jt)

	if sk != nil {
		t.Error("The returned signing key should be nil.")
	}

	expectValidationError(t, err, ValidationErrorInvalidIssuerType, http.StatusUnauthorized, nil)
	pm.assertDone()
}

func Test_getSigningKey_UsingTokenWithEmptyIssuer(t *testing.T) {
	pm, _, _, _, tv := createIDTokenValidator(t)

	go func() {
		pm.assertGetProviders([]Provider{{Issuer: "https://issuer", ClientIDs: []string{"client"}}}, nil)
		pm.assertGetProviders([]Provider{{Issuer: "https://issuer", ClientIDs: []string{"client"}}}, nil)

		pm.close()
	}()

	jt := jwt.New(jwt.SigningMethodRS256)

	// The token has no 'iss' claim
	sk, err := tv.getSigningKey(jt)

	if sk != nil {
		t.Error("The returned signing key should be nil.")
	}

	expectValidationError(t, err, ValidationErrorInvalidIssuerType, http.StatusUnauthorized, nil)

	// The token has '' as 'iss' claim
	jt.Claims.(jwt.MapClaims)["iss"] = ""
	sk, err = tv.getSigningKey(jt)

	if sk != nil {
		t.Error("The returned signing key should be nil.")
	}

	expectValidationError(t, err, ValidationErrorInvalidIssuer, http.StatusUnauthorized, nil)

	pm.assertDone()
}

func Test_getSigningKey_UsingTokenWithUnknownIssuer(t *testing.T) {
	pm, _, _, _, tv := createIDTokenValidator(t)

	go func() {
		pm.assertGetProviders([]Provider{{Issuer: "https://issuer", ClientIDs: []string{"client"}}}, nil)
		pm.close()
	}()

	jt := jwt.New(jwt.SigningMethodRS256)
	jt.Claims.(jwt.MapClaims)["iss"] = "http://unknown"

	// The token has no 'iss' claim
	sk, err := tv.getSigningKey(jt)

	if sk != nil {
		t.Error("The returned signing key should be nil.")
	}

	expectValidationError(t, err, ValidationErrorIssuerNotFound, http.StatusUnauthorized, nil)
	pm.assertDone()
}

func Test_getSigningKey_UsingTokenWithInvalidAudienceType(t *testing.T) {
	pm, _, _, _, tv := createIDTokenValidator(t)

	go func() {
		pm.assertGetProviders([]Provider{{Issuer: "https://issuer", ClientIDs: []string{"client"}}}, nil)
		pm.close()
	}()

	jt := jwt.New(jwt.SigningMethodRS256)
	jt.Claims.(jwt.MapClaims)["iss"] = "https://issuer"
	jt.Claims.(jwt.MapClaims)["aud"] = 0 // Expected 'aud' type is string

	sk, err := tv.getSigningKey(jt)

	if sk != nil {
		t.Error("The returned signing key should be nil.")
	}

	expectValidationError(t, err, ValidationErrorInvalidAudienceType, http.StatusUnauthorized, nil)
	pm.assertDone()
}

func Test_getSigningKey_UsingTokenWithInvalidAudience(t *testing.T) {
	pm, _, _, _, tv := createIDTokenValidator(t)

	go func() {
		pm.assertGetProviders([]Provider{{Issuer: "https://issuer", ClientIDs: []string{"client"}}}, nil)
		pm.assertGetProviders([]Provider{{Issuer: "https://issuer", ClientIDs: []string{"client"}}}, nil)
		pm.close()
	}()

	jt := jwt.New(jwt.SigningMethodRS256)
	jt.Claims.(jwt.MapClaims)["iss"] = "https://issuer"

	// No audience claim
	sk, err := tv.getSigningKey(jt)

	if sk != nil {
		t.Error("The returned signing key should be nil.")
	}

	expectValidationError(t, err, ValidationErrorInvalidAudienceType, http.StatusUnauthorized, nil)

	// Empty audience claim.
	jt.Claims.(jwt.MapClaims)["aud"] = ""
	sk, err = tv.getSigningKey(jt)

	if sk != nil {
		t.Error("The returned signing key should be nil.")
	}

	expectValidationError(t, err, ValidationErrorInvalidAudience, http.StatusUnauthorized, nil)
	pm.assertDone()

}

func Test_getSigningKey_UsingTokenWithUnknownAudience(t *testing.T) {
	pm, _, _, _, tv := createIDTokenValidator(t)

	go func() {
		pm.assertGetProviders([]Provider{{Issuer: "https://issuer", ClientIDs: []string{"client1", "client2"}}}, nil)
		pm.close()
	}()

	jt := jwt.New(jwt.SigningMethodRS256)
	jt.Claims.(jwt.MapClaims)["iss"] = "https://issuer"
	jt.Claims.(jwt.MapClaims)["aud"] = "client3" // unknown audience

	sk, err := tv.getSigningKey(jt)

	if sk != nil {
		t.Error("The returned signing key should be nil.")
	}

	expectValidationError(t, err, ValidationErrorAudienceNotFound, http.StatusUnauthorized, nil)
	pm.assertDone()
}

func Test_getSigningKey_UsingTokenWithUnknownMultipleAudiences(t *testing.T) {
	pm, _, _, _, tv := createIDTokenValidator(t)

	go func() {
		pm.assertGetProviders([]Provider{{Issuer: "https://issuer", ClientIDs: []string{"client1", "client2"}}}, nil)
		pm.close()
	}()

	jt := jwt.New(jwt.SigningMethodRS256)
	jt.Claims.(jwt.MapClaims)["iss"] = "https://issuer"
	jt.Claims.(jwt.MapClaims)["aud"] = []interface{}{"client3", "client4"} // unknown audiences

	sk, err := tv.getSigningKey(jt)

	if sk != nil {
		t.Error("The returned signing key should be nil.")
	}

	expectValidationError(t, err, ValidationErrorAudienceNotFound, http.StatusUnauthorized, nil)
	pm.assertDone()
}

func Test_getSigningKey_UsingTokenWithInvalidSubjectType(t *testing.T) {
	pm, _, _, _, tv := createIDTokenValidator(t)

	go func() {
		pm.assertGetProviders([]Provider{{Issuer: "https://issuer", ClientIDs: []string{"client"}}}, nil)
		pm.close()
	}()

	jt := jwt.New(jwt.SigningMethodRS256)
	jt.Claims.(jwt.MapClaims)["iss"] = "https://issuer"
	jt.Claims.(jwt.MapClaims)["aud"] = "client"
	jt.Claims.(jwt.MapClaims)["sub"] = 0 // The expected 'sub' claim type is string
	sk, err := tv.getSigningKey(jt)

	if sk != nil {
		t.Error("The returned signing key should be nil.")
	}

	expectValidationError(t, err, ValidationErrorInvalidSubjectType, http.StatusUnauthorized, nil)
	pm.assertDone()
}

func Test_getSigningKey_UsingValidToken_WhenSigningKeyGetterReturnsError(t *testing.T) {
	pm, _, sm, _, tv := createIDTokenValidator(t)

	iss := "https://issuer"
	keyID := "kid"
	ee := &ValidationError{Code: ValidationErrorIssuerNotFound, HTTPStatus: http.StatusUnauthorized}

	go func() {
		pm.assertGetProviders([]Provider{{Issuer: iss, ClientIDs: []string{"client"}}}, nil)
		sm.assertGetSigningKey(iss, keyID, nil, ee)
		pm.close()
		sm.close()
	}()

	jt := jwt.New(jwt.SigningMethodRS256)
	jt.Claims.(jwt.MapClaims)["iss"] = iss
	jt.Claims.(jwt.MapClaims)["aud"] = "client"
	jt.Claims.(jwt.MapClaims)["sub"] = "subject1"
	jt.Header["kid"] = keyID

	_, err := tv.getSigningKey(jt)

	expectValidationError(t, err, ee.Code, ee.HTTPStatus, nil)
	pm.assertDone()
	sm.assertDone()
}

func Test_getSigningKey_UsingValidToken_WhenSigningKeyGetterSucceeds(t *testing.T) {
	pm, _, sm, kp, tv := createIDTokenValidator(t)

	iss := "https://issuer"
	keyID := "kid"
	esk := "signingKey"
	pk := &rsa.PublicKey{N: nil, E: 345}

	go func() {
		pm.assertGetProviders([]Provider{{Issuer: iss, ClientIDs: []string{"client"}}}, nil)
		sm.assertGetSigningKey(iss, keyID, []byte(esk), nil)
		kp.assertParse([]byte(esk), pk, nil)
		pm.close()
		sm.close()
		kp.close()
	}()

	jt := jwt.New(jwt.SigningMethodRS256)
	jt.Claims.(jwt.MapClaims)["iss"] = iss
	jt.Claims.(jwt.MapClaims)["aud"] = "client"
	jt.Claims.(jwt.MapClaims)["sub"] = "subject1"
	jt.Header["kid"] = keyID

	rsk, err := tv.getSigningKey(jt)

	if err != nil {
		t.Error("An error was returned but not expected.", err)
	}

	expectSigningKey(t, rsk, jt, pk)

	pm.assertDone()
	sm.assertDone()
}

func Test_getSigningKey_UsingValidToken_WithoutKeyIdentifier_WhenSigningKeyGetterSucceeds(t *testing.T) {
	pm, _, sm, kp, tv := createIDTokenValidator(t)

	iss := "https://issuer"
	keyID := ""
	esk := "signingKey"
	pk := &rsa.PublicKey{N: nil, E: 345}

	go func() {
		pm.assertGetProviders([]Provider{{Issuer: iss, ClientIDs: []string{"client"}}}, nil)
		sm.assertGetSigningKey(iss, keyID, []byte(esk), nil)
		kp.assertParse([]byte(esk), pk, nil)
		pm.close()
		sm.close()
		kp.close()
	}()

	jt := jwt.New(jwt.SigningMethodRS256)
	jt.Claims.(jwt.MapClaims)["iss"] = iss
	jt.Claims.(jwt.MapClaims)["aud"] = "client"
	jt.Claims.(jwt.MapClaims)["sub"] = "subject1"

	rsk, err := tv.getSigningKey(jt)

	if err != nil {
		t.Error("An error was returned but not expected.", err)
	}

	expectSigningKey(t, rsk, jt, pk)

	pm.assertDone()
	sm.assertDone()
}

func Test_getSigningKey_UsingValidTokenWithMultipleAudiences(t *testing.T) {
	pm, _, sm, kp, tv := createIDTokenValidator(t)

	iss := "https://issuer"
	keyID := "kid"
	esk := "signingKey"
	pk := &rsa.PublicKey{N: nil, E: 345}

	go func() {
		pm.assertGetProviders([]Provider{{Issuer: iss, ClientIDs: []string{"client"}}}, nil)
		sm.assertGetSigningKey(iss, keyID, []byte(esk), nil)
		kp.assertParse([]byte(esk), pk, nil)
		pm.close()
		sm.close()
		kp.close()
	}()

	jt := jwt.New(jwt.SigningMethodRS256)
	jt.Claims.(jwt.MapClaims)["iss"] = iss
	jt.Claims.(jwt.MapClaims)["aud"] = []interface{}{"unknown", "client"}
	jt.Claims.(jwt.MapClaims)["sub"] = "subject1"
	jt.Header["kid"] = keyID

	rsk, err := tv.getSigningKey(jt)

	if err != nil {
		t.Error("An error was returned but not expected.", err)
	}

	expectSigningKey(t, rsk, jt, pk)

	pm.assertDone()
	sm.assertDone()
}

func Test_renewAndGetSigningKey_UsingValidToken_WhenFlushCachedSigningKeysReturnsError(t *testing.T) {
	_, _, sm, _, tv := createIDTokenValidator(t)

	ee := &ValidationError{Code: ValidationErrorIssuerNotFound, HTTPStatus: http.StatusUnauthorized}
	go func() {
		sm.assertFlushCachedSigningKeys(anything, ee)
		sm.close()
	}()

	jt := jwt.New(jwt.SigningMethodRS256)
	jt.Claims.(jwt.MapClaims)["iss"] = ""

	_, err := tv.renewAndGetSigningKey(jt)

	expectValidationError(t, err, ee.Code, ee.HTTPStatus, nil)

	sm.assertDone()
}

func Test_renewAndGetSigningKey_UsingValidToken_WhenGetSigningKeyReturnsError(t *testing.T) {
	_, _, sm, _, tv := createIDTokenValidator(t)

	ee := &ValidationError{Code: ValidationErrorIssuerNotFound, HTTPStatus: http.StatusUnauthorized}
	go func() {
		sm.assertFlushCachedSigningKeys(anything, nil)
		sm.assertGetSigningKey(anything, anything, nil, ee)
		sm.close()
	}()

	jt := jwt.New(jwt.SigningMethodRS256)
	jt.Claims.(jwt.MapClaims)["iss"] = ""
	jt.Header["kid"] = ""

	_, err := tv.renewAndGetSigningKey(jt)

	expectValidationError(t, err, ee.Code, ee.HTTPStatus, nil)

	sm.assertDone()
}

func Test_renewAndGetSigningKey_UsingValidToken_WhenGetSigningKeySucceeds(t *testing.T) {
	_, _, sm, kp, tv := createIDTokenValidator(t)
	esk := "signingKey"
	pk := &rsa.PublicKey{N: nil, E: 365}

	go func() {
		sm.assertFlushCachedSigningKeys(anything, nil)
		sm.assertGetSigningKey(anything, anything, []byte(esk), nil)
		kp.assertParse([]byte(esk), pk, nil)
		sm.close()
		kp.close()
	}()

	jt := jwt.New(jwt.SigningMethodRS256)
	jt.Claims.(jwt.MapClaims)["iss"] = ""
	jt.Header["kid"] = ""

	rsk, err := tv.renewAndGetSigningKey(jt)

	if err != nil {
		t.Error("An error was returned but not expected.", err)
	}

	expectSigningKey(t, rsk, jt, pk)

	sm.assertDone()
}

func Test_validate_WhenParserReturnsErrorFirstTime(t *testing.T) {
	_, jm, _, _, tv := createIDTokenValidator(t)

	je := &jwt.ValidationError{Errors: jwt.ValidationErrorNotValidYet}
	ee := &ValidationError{Code: ValidationErrorJwtValidationFailure, HTTPStatus: http.StatusUnauthorized}

	go func() {
		jm.assertParse(anything, nil, je)
		jm.close()
	}()

	_, err := tv.validate(anything)

	expectValidationError(t, err, ee.Code, ee.HTTPStatus, ee.Err)

	jm.assertDone()
}

func Test_validate_WhenParserSuceedsFirstTime(t *testing.T) {
	_, jm, _, _, tv := createIDTokenValidator(t)

	jt := &jwt.Token{}

	go func() {
		jm.assertParse(anything, jt, nil)
		jm.close()
	}()

	rjt, err := tv.validate(anything)

	if err != nil {
		t.Error("Unexpected error was returned.", err)
	}

	if rjt != jt {
		t.Errorf("Expected %+v, but got %+v.", jt, rjt)
	}

	jm.assertDone()
}

func Test_validate_WhenParserReturnsErrorSecondTime(t *testing.T) {
	_, jm, _, _, tv := createIDTokenValidator(t)

	jfe := &jwt.ValidationError{Errors: jwt.ValidationErrorSignatureInvalid}
	je := &jwt.ValidationError{Errors: jwt.ValidationErrorMalformed}
	ee := &ValidationError{Code: ValidationErrorJwtValidationFailure, HTTPStatus: http.StatusBadRequest}

	go func() {
		jm.assertParse(anything, nil, jfe)
		jm.assertParse(anything, nil, je)
		jm.close()
	}()

	_, err := tv.validate(anything)

	expectValidationError(t, err, ee.Code, ee.HTTPStatus, ee.Err)

	jm.assertDone()
}

func Test_validate_WhenParserReturnsSignatureInvalidErrorSecondTime(t *testing.T) {
	_, jm, _, _, tv := createIDTokenValidator(t)

	je := &jwt.ValidationError{Errors: jwt.ValidationErrorSignatureInvalid}
	ee := &ValidationError{Code: ValidationErrorJwtValidationFailure, HTTPStatus: http.StatusUnauthorized}

	go func() {
		jm.assertParse(anything, nil, je)
		jm.assertParse(anything, nil, je)
		jm.close()
	}()

	_, err := tv.validate(anything)

	expectValidationError(t, err, ee.Code, ee.HTTPStatus, ee.Err)

	jm.assertDone()
}

func Test_validate_WhenParserSuceedsSecondTime(t *testing.T) {
	_, jm, _, _, tv := createIDTokenValidator(t)

	jfe := &jwt.ValidationError{Errors: jwt.ValidationErrorSignatureInvalid}

	jt := &jwt.Token{}

	go func() {
		jm.assertParse(anything, jt, jfe)
		jm.assertParse(anything, jt, nil)
		jm.close()
	}()

	rjt, err := tv.validate(anything)

	if err != nil {
		t.Error("Unexpected error was returned.", err)
	}

	if rjt != jt {
		t.Errorf("Expected %+v, but got %+v.", jt, rjt)
	}

	jm.assertDone()
}

func expectSigningKey(t *testing.T, rsk interface{}, jt *jwt.Token, esk *rsa.PublicKey) {

	if rsk == nil {
		t.Fatal("The returned signing key was nil.")
	}

	if sk, ok := rsk.(*rsa.PublicKey); ok {
		if sk.E != esk.E {
			t.Error("Expected signing key", esk, "but got", sk)
		}
	} else {
		t.Errorf("Expected signing key type '*rsa.PublicKey', but got %T", rsk)
	}
}

func createIDTokenValidator(t *testing.T) (*providersGetterMock, *jwtParserMock, *signingKeyGetterMock, *rsaParserMock, *idTokenValidator) {
	pm := newProvidersGetterMock(t)
	jm := newJwtParserMock(t)
	sm := newSigningKeyGetterMock(t)
	kp := newRSAParserMock(t)
	return pm, jm, sm, kp, &idTokenValidator{pm.getProviders, jm.parse, sm, kp.parse}
}
