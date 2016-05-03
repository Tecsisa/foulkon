package openid

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/square/go-jose"
)

func Test_getsigningKeySet_WhenGetConfigurationReturnsError(t *testing.T) {
	configGetter, _, _, skProv := createSigningKeySetProvider(t)

	ee := &ValidationError{Code: ValidationErrorGetOpenIdConfigurationFailure, HTTPStatus: http.StatusUnauthorized}

	go func() {
		configGetter.assertGetConfiguration(anything, configuration{}, ee)
		configGetter.close()
	}()

	sk, re := skProv.getSigningKeySet(anything)

	expectValidationError(t, re, ee.Code, ee.HTTPStatus, nil)

	if sk != nil {
		t.Error("The returned signing keys should be nil")
	}

	configGetter.assertDone()
}

func Test_getsigningKeySet_WhenGetJwksReturnsError(t *testing.T) {
	configGetter, jwksGetter, _, skProv := createSigningKeySetProvider(t)

	ee := &ValidationError{Code: ValidationErrorGetJwksFailure, HTTPStatus: http.StatusUnauthorized}

	go func() {
		configGetter.assertGetConfiguration(anything, configuration{}, nil)
		configGetter.close()
		jwksGetter.assertGetJwks(anything, jose.JsonWebKeySet{}, ee)
		jwksGetter.close()

	}()

	sk, re := skProv.getSigningKeySet(anything)

	expectValidationError(t, re, ee.Code, ee.HTTPStatus, nil)

	if sk != nil {
		t.Error("The returned signing keys should be nil")
	}

	configGetter.assertDone()
	jwksGetter.assertDone()
}

func Test_getsigningKeySet_WhenJwkSetIsEmpty(t *testing.T) {
	configGetter, jwksGetter, _, skProv := createSigningKeySetProvider(t)

	ee := &ValidationError{Code: ValidationErrorEmptyJwk, HTTPStatus: http.StatusUnauthorized}

	go func() {
		configGetter.assertGetConfiguration(anything, configuration{}, nil)
		configGetter.close()
		jwksGetter.assertGetJwks(anything, jose.JsonWebKeySet{}, nil)
		jwksGetter.close()

	}()

	sk, re := skProv.getSigningKeySet(anything)

	expectValidationError(t, re, ee.Code, ee.HTTPStatus, nil)

	if sk != nil {
		t.Error("The returned signing keys should be nil")
	}

	configGetter.assertDone()
	jwksGetter.assertDone()
}

func Test_getsigningKeySet_WhenKeyEncodingReturnsError(t *testing.T) {
	configGetter, jwksGetter, pemEncoder, skProv := createSigningKeySetProvider(t)

	ee := &ValidationError{Code: ValidationErrorMarshallingKey, HTTPStatus: http.StatusInternalServerError}
	ejwks := jose.JsonWebKeySet{Keys: []jose.JsonWebKey{{Key: nil}}}

	go func() {
		configGetter.assertGetConfiguration(anything, configuration{}, nil)
		configGetter.close()
		jwksGetter.assertGetJwks(anything, ejwks, nil)
		jwksGetter.close()
		pemEncoder.assertPEMEncodePublicKey(nil, nil, ee)
		pemEncoder.close()
	}()

	sk, re := skProv.getSigningKeySet(anything)

	expectValidationError(t, re, ee.Code, ee.HTTPStatus, nil)

	if sk != nil {
		t.Error("The returned signing keys should be nil")
	}

	configGetter.assertDone()
	jwksGetter.assertDone()
	pemEncoder.assertDone()
}

func Test_getsigningKeySet_WhenKeyEncodingReturnsSuccess(t *testing.T) {
	configGetter, jwksGetter, pemEncoder, skProv := createSigningKeySetProvider(t)

	keys := make([]jose.JsonWebKey, 2)
	encryptedKeys := make([]signingKey, 2)

	for i := 0; i < cap(keys); i = i + 1 {
		keys[i] = jose.JsonWebKey{KeyID: fmt.Sprintf("%v", i), Key: i}
		encryptedKeys[i] = signingKey{keyID: fmt.Sprintf("%v", i), key: []byte(fmt.Sprintf("%v", i))}
	}

	ejwks := jose.JsonWebKeySet{Keys: keys}
	go func() {
		configGetter.assertGetConfiguration(anything, configuration{}, nil)
		jwksGetter.assertGetJwks(anything, ejwks, nil)
		for i, encryptedKey := range encryptedKeys {
			pemEncoder.assertPEMEncodePublicKey(keys[i].Key, encryptedKey.key, nil)
		}
		configGetter.close()
		jwksGetter.close()
		pemEncoder.close()
	}()

	sk, re := skProv.getSigningKeySet(anything)

	if re != nil {
		t.Error("An error was returned but not expected.")
	}

	if sk == nil {
		t.Fatal("The returned signing keys should be not nil")
	}

	if len(sk) != len(encryptedKeys) {
		t.Error("Returned", len(sk), "encrypted keys, but expected", len(encryptedKeys))
	}

	for i, encryptedKey := range encryptedKeys {
		if encryptedKey.keyID != sk[i].keyID {
			t.Error("Key at", i, "should have keyID", encryptedKey.keyID, "but was", sk[i].keyID)
		}
		if string(encryptedKey.key) != string(sk[i].key) {
			t.Error("Key at", i, "should be", encryptedKey.key, "but was", sk[i].key)
		}
	}

	configGetter.assertDone()
	jwksGetter.assertDone()
	pemEncoder.assertDone()

}

func createSigningKeySetProvider(t *testing.T) (*configurationGetterMock, *jwksGetterMock, *pemEncoderMock, signingKeySetProvider) {
	configGetter := newConfigurationGetterMock(t)
	jwksGetter := newJwksGetterMock(t)
	pemEncoder := newPEMEncoderMock(t)

	skProv := signingKeySetProvider{configGetter: configGetter, jwksGetter: jwksGetter, keyEncoder: pemEncoder.pemEncodePublicKey}
	return configGetter, jwksGetter, pemEncoder, skProv
}
