package openid

import "testing"

type signingKeyGetterMock struct {
	t     *testing.T
	Calls chan Call
}

func newSigningKeyGetterMock(t *testing.T) *signingKeyGetterMock {
	return &signingKeyGetterMock{t, make(chan Call)}
}

type getSigningKeyCall struct {
	iss   string
	keyID string
}

type getSigningKeyResp struct {
	key []byte
	err error
}

type flushCachedSigningKeysCall struct {
	iss string
}

type flushCachedSigningKeysResp struct {
	err error
}

func (s *signingKeyGetterMock) getSigningKey(iss string, keyID string) ([]byte, error) {
	s.Calls <- &getSigningKeyCall{iss, keyID}
	sr := (<-s.Calls).(*getSigningKeyResp)
	return sr.key, sr.err
}

func (s *signingKeyGetterMock) flushCachedSigningKeys(iss string) error {
	s.Calls <- &flushCachedSigningKeysCall{iss}
	sr := (<-s.Calls).(*flushCachedSigningKeysResp)
	return sr.err
}

func (s *signingKeyGetterMock) assertGetSigningKey(iss string, keyID string, key []byte, err error) {
	call := (<-s.Calls).(*getSigningKeyCall)
	if iss != anything && call.iss != iss {
		s.t.Error("Expected getSigningKey with issuer", iss, "but was", call.iss)
	}
	if keyID != anything && call.keyID != keyID {
		s.t.Error("Expected getSigningKey with key ID", keyID, "but was", call.keyID)
	}
	s.Calls <- &getSigningKeyResp{key, err}
}

func (s *signingKeyGetterMock) assertFlushCachedSigningKeys(iss string, err error) {
	call := (<-s.Calls).(*flushCachedSigningKeysCall)
	if iss != anything && call.iss != iss {
		s.t.Error("Expected flushCachedSigningKeys with issuer", iss, "but was", call.iss)
	}

	s.Calls <- &flushCachedSigningKeysResp{err}
}

func (s *signingKeyGetterMock) close() {
	close(s.Calls)
}

func (s *signingKeyGetterMock) assertDone() {
	if _, more := <-s.Calls; more {
		s.t.Fatal("Did not expect more calls.")
	}
}
