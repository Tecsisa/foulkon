package openid

import (
	"bytes"
	"crypto/rsa"
	"testing"
)

type rsaParserMock struct {
	t     *testing.T
	Calls chan Call
}

func newRSAParserMock(t *testing.T) *rsaParserMock {
	return &rsaParserMock{t, make(chan Call)}
}

type rsaParseCall struct {
	key []byte
}

type rsaParseResp struct {
	pk *rsa.PublicKey
	e  error
}

func (p *rsaParserMock) parse(key []byte) (*rsa.PublicKey, error) {
	p.Calls <- &rsaParseCall{key}
	pr := (<-p.Calls).(*rsaParseResp)
	return pr.pk, pr.e
}

func (p *rsaParserMock) assertParse(key []byte, pk *rsa.PublicKey, e error) {
	call := (<-p.Calls).(*rsaParseCall)
	if key != nil && bytes.Compare(key, call.key) != 0 {
		p.t.Error("Expected parse with", key, "but was", call.key)
	}

	p.Calls <- &rsaParseResp{pk, e}
}

func (p *rsaParserMock) close() {
	close(p.Calls)
}

func (p *rsaParserMock) assertDone() {
	if _, more := <-p.Calls; more {
		p.t.Fatal("Did not expect more calls.")
	}
}
