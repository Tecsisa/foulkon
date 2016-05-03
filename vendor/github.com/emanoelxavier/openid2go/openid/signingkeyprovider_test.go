package openid

import (
	"net/http"
	"testing"
)

func Test_getSigningKey_WhenKeyIsCached(t *testing.T) {
	_, keyCache := createSigningKeyProvider(t)

	iss := "issuer"
	kid := "kid1"
	key := "signingKey"
	keyCache.jwksMap[iss] = []signingKey{{keyID: kid, key: []byte(key)}}

	expectKey(t, keyCache, iss, kid, key)
}

func Test_getSigningKey_WhenKeyIsNotCached_WhenProviderReturnsKey(t *testing.T) {
	keyGetter, keyCache := createSigningKeyProvider(t)

	iss := "issuer"
	kid := "kid1"
	key := "signingKey"

	go func() {
		keyGetter.assertGetSigningKeySet(iss, []signingKey{{keyID: kid, key: []byte(key)}}, nil)
		keyGetter.close()
	}()

	expectKey(t, keyCache, iss, kid, key)

	// Validate that the key is cached
	expectCachedKid(t, keyCache, iss, kid, key)

	keyGetter.assertDone()
}

func Test_getSigningKey_WhenProviderReturnsError(t *testing.T) {
	keyGetter, keyCache := createSigningKeyProvider(t)

	iss := "issuer"
	kid := "kid1"
	ee := &ValidationError{Code: ValidationErrorGetJwksFailure, HTTPStatus: http.StatusUnauthorized}

	go func() {
		keyGetter.assertGetSigningKeySet(iss, nil, ee)
		keyGetter.close()
	}()

	rk, re := keyCache.getSigningKey(iss, kid)

	expectValidationError(t, re, ee.Code, ee.HTTPStatus, nil)

	if rk != nil {
		t.Error("A key was returned but not expected")
	}

	cachedKeys := keyCache.jwksMap[iss]
	if len(cachedKeys) != 0 {
		t.Fatal("There shouldnt be cached keys for the targeted issuer.")
	}

	keyGetter.assertDone()
}

func Test_getSigningKey_WhenKeyIsNotFound(t *testing.T) {
	keyGetter, keyCache := createSigningKeyProvider(t)

	iss := "issuer"
	kid := "kid1"
	tkid := "kid2"
	key := "signingKey"

	go func() {
		keyGetter.assertGetSigningKeySet(iss, []signingKey{{keyID: kid, key: []byte(key)}}, nil)
		keyGetter.close()
	}()

	rk, re := keyCache.getSigningKey(iss, tkid)

	expectValidationError(t, re, ValidationErrorKidNotFound, http.StatusUnauthorized, nil)

	if rk != nil {
		t.Error("A key was returned but not expected")
	}

	expectCachedKid(t, keyCache, iss, kid, key)

	keyGetter.assertDone()
}

func Test_flushCachedSigningKeys_FlushedKeysAreDeleted(t *testing.T) {
	_, keyCache := createSigningKeyProvider(t)

	iss := "issuer"
	iss2 := "issuer2"
	kid := "kid1"
	key := "signingKey"
	keyCache.jwksMap[iss] = []signingKey{{keyID: kid, key: []byte(key)}}
	keyCache.jwksMap[iss2] = []signingKey{{keyID: kid, key: []byte(key)}}

	keyCache.flushCachedSigningKeys(iss2)

	dk := keyCache.jwksMap[iss2]

	if dk != nil {
		t.Error("Flushed keys should not be in the cache.")
	}

	expectCachedKid(t, keyCache, iss, kid, key)
}

func Test_flushCachedSigningKey_RetrieveFlushedKey(t *testing.T) {
	keyGetter, keyCache := createSigningKeyProvider(t)

	iss := "issuer"
	kid := "kid1"
	key := "signingKey"

	go func() {
		keyGetter.assertGetSigningKeySet(iss, []signingKey{{keyID: kid, key: []byte(key)}}, nil)
		keyGetter.assertGetSigningKeySet(iss, []signingKey{{keyID: kid, key: []byte(key)}}, nil)

		keyGetter.close()
	}()

	// Get the signing key not yet cached will cache it.
	expectKey(t, keyCache, iss, kid, key)

	// Flush the signing keys for the given provider.
	keyCache.flushCachedSigningKeys(iss)

	// Get the signing key will once again call the provider and cache the keys.

	expectKey(t, keyCache, iss, kid, key)

	// Validate that the key is cached
	expectCachedKid(t, keyCache, iss, kid, key)

	keyGetter.assertDone()
}

func expectCachedKid(t *testing.T, keyProv *signingKeyProvider, iss string, kid string, key string) {

	cachedKeys := keyProv.jwksMap[iss]
	if len(cachedKeys) == 0 {
		t.Fatal("The keys were not cached as expected.")
	}

	foundKid := false
	for _, cachedKey := range cachedKeys {
		if cachedKey.keyID == kid {
			foundKid = true
			if keyStr := string(cachedKey.key); keyStr != key {
				t.Error("Expected key", key, "but got", keyStr)
			}

			continue
		}
	}

	if !foundKid {
		t.Error("A key with key id", kid, "was not found.")
	}
}

func expectKey(t *testing.T, c signingKeyGetter, iss string, kid string, key string) {

	sk, re := c.getSigningKey(iss, kid)

	if re != nil {
		t.Error("An error was returned but not expected.")
	}

	if sk == nil {
		t.Fatal("The returned signing key should not be nil.")
	}

	keyStr := string(sk)

	if keyStr != key {
		t.Error("Expected key", key, "but got", keyStr)
	}
}

func createSigningKeyProvider(t *testing.T) (*signingKeySetGetterMock, *signingKeyProvider) {
	mock := newSigningKeySetGetterMock(t)
	return mock, newSigningKeyProvider(mock)
}
