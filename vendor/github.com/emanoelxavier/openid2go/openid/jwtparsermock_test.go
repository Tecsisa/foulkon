package openid

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
)

type jwtParserMock struct {
	t     *testing.T
	Calls chan Call
}

func newJwtParserMock(t *testing.T) *jwtParserMock {
	return &jwtParserMock{t, make(chan Call)}
}

type parseCall struct {
	t  string
	kf jwt.Keyfunc
}

type parseResp struct {
	jt *jwt.Token
	e  error
}

func (p *jwtParserMock) parse(t string, kf jwt.Keyfunc) (*jwt.Token, error) {
	p.Calls <- &parseCall{t, kf}
	pr := (<-p.Calls).(*parseResp)
	return pr.jt, pr.e
}

func (p *jwtParserMock) assertParse(t string, jt *jwt.Token, e error) {
	call := (<-p.Calls).(*parseCall)
	if call.t != anything && t != call.t {
		p.t.Error("Expected parse with", t, "but was", call.t)
	}

	p.Calls <- &parseResp{jt, e}
}

func (p *jwtParserMock) close() {
	close(p.Calls)
}

func (p *jwtParserMock) assertDone() {
	if _, more := <-p.Calls; more {
		p.t.Fatal("Did not expect more calls.")
	}
}
