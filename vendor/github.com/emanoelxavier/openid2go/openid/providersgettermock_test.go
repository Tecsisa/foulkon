package openid

import "testing"

type providersGetterMock struct {
	t     *testing.T
	Calls chan Call
}

func newProvidersGetterMock(t *testing.T) *providersGetterMock {
	return &providersGetterMock{t, make(chan Call)}
}

type getProvidersCall struct {
}

type getProvidersResp struct {
	provs []Provider
	e     error
}

func (p *providersGetterMock) getProviders() ([]Provider, error) {
	p.Calls <- &getProvidersCall{}
	gr := (<-p.Calls).(*getProvidersResp)
	return gr.provs, gr.e
}

func (p *providersGetterMock) assertGetProviders(ps []Provider, e error) {
	call := (<-p.Calls).(*getProvidersCall)
	if call == nil {
		p.t.Error("Expected a getProviders call but it was nil.")
	}

	p.Calls <- &getProvidersResp{ps, e}
}

func (p *providersGetterMock) close() {
	close(p.Calls)
}

func (p *providersGetterMock) assertDone() {
	if _, more := <-p.Calls; more {
		p.t.Fatal("Did not expect more calls.")
	}
}
