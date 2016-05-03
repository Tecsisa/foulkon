package openid

import (
	"fmt"
	"net/http"
)

type signingKeyGetter interface {
	flushCachedSigningKeys(issuer string) error
	getSigningKey(issuer string, kid string) ([]byte, error)
}

type signingKeyProvider struct {
	keySetGetter signingKeySetGetter
	jwksMap      map[string][]signingKey
}

func newSigningKeyProvider(kg signingKeySetGetter) *signingKeyProvider {
	keyMap := make(map[string][]signingKey)
	return &signingKeyProvider{kg, keyMap}
}

func (s *signingKeyProvider) flushCachedSigningKeys(issuer string) error {
	delete(s.jwksMap, issuer)
	return nil
}

func (s *signingKeyProvider) refreshSigningKeys(issuer string) error {
	skeys, err := s.keySetGetter.getSigningKeySet(issuer)

	if err != nil {
		return err
	}

	s.jwksMap[issuer] = skeys
	return nil
}

func (s *signingKeyProvider) getSigningKey(issuer string, kid string) ([]byte, error) {
	sk := findKey(s.jwksMap, issuer, kid)

	if sk != nil {
		return sk, nil
	}

	err := s.refreshSigningKeys(issuer)

	if err != nil {
		return nil, err
	}

	sk = findKey(s.jwksMap, issuer, kid)

	if sk == nil {
		return nil, &ValidationError{Code: ValidationErrorKidNotFound, Message: fmt.Sprintf("The jwk set retrieved for the issuer %v does not contain a key identifier %v.", issuer, kid), HTTPStatus: http.StatusUnauthorized}
	}

	return sk, nil
}

func findKey(km map[string][]signingKey, issuer string, kid string) []byte {
	if skSet, ok := km[issuer]; ok {
		if kid == "" {
			return skSet[0].key
		} else {
			for _, sk := range skSet {
				if sk.keyID == kid {
					return sk.key
				}
			}
		}
	}

	return nil
}
