package openid

import "testing"

type pemEncoderMock struct {
	t     *testing.T
	Calls chan Call
}

func newPEMEncoderMock(t *testing.T) *pemEncoderMock {
	return &pemEncoderMock{t, make(chan Call)}
}

type pemEncodeCall struct {
	key interface{}
}

type pemEncodeResponse struct {
	key []byte
	err error
}

func (p *pemEncoderMock) pemEncodePublicKey(key interface{}) ([]byte, error) {
	p.Calls <- &pemEncodeCall{key}
	gr := (<-p.Calls).(*pemEncodeResponse)
	return gr.key, gr.err
}

func (p *pemEncoderMock) assertPEMEncodePublicKey(key interface{}, enkey []byte, err error) {
	call := (<-p.Calls).(*pemEncodeCall)
	if call.key != key {
		p.t.Error("Expected pemEncode key  with", key, "but was", call.key)
	}
	p.Calls <- &pemEncodeResponse{enkey, err}
}

func (p *pemEncoderMock) close() {
	close(p.Calls)
}

func (p *pemEncoderMock) assertDone() {
	if _, more := <-p.Calls; more {
		p.t.Fatal("Did not expect more calls.")
	}
}
