package openid

import "testing"

type signingKeySetGetterMock struct {
	t     *testing.T
	Calls chan Call
}

func newSigningKeySetGetterMock(t *testing.T) *signingKeySetGetterMock {
	return &signingKeySetGetterMock{t, make(chan Call)}
}

type getSigningKeySetCall struct {
	iss string
}

type getSigningKeySetResponse struct {
	keys []signingKey
	err  error
}

func (c *signingKeySetGetterMock) getSigningKeySet(iss string) ([]signingKey, error) {
	c.Calls <- &getSigningKeySetCall{iss}
	sr := (<-c.Calls).(*getSigningKeySetResponse)
	return sr.keys, sr.err
}

func (c *signingKeySetGetterMock) assertGetSigningKeySet(iss string, keys []signingKey, err error) {
	call := (<-c.Calls).(*getSigningKeySetCall)
	if iss != anything && call.iss != iss {
		c.t.Error("Expected getSigningKeySet with issuer", iss, "but was", call.iss)
	}
	c.Calls <- &getSigningKeySetResponse{keys, err}
}

func (c *signingKeySetGetterMock) close() {
	close(c.Calls)
}

func (c *signingKeySetGetterMock) assertDone() {
	if _, more := <-c.Calls; more {
		c.t.Fatal("Did not expect more calls.")
	}
}
