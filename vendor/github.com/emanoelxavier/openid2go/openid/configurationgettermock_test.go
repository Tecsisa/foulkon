package openid

import "testing"

type configurationGetterMock struct {
	t     *testing.T
	Calls chan Call
}

func newConfigurationGetterMock(t *testing.T) *configurationGetterMock {
	return &configurationGetterMock{t, make(chan Call)}
}

type getConfigurationCall struct {
	iss string
}

type getConfigurationResponse struct {
	config configuration
	err    error
}

func (c *configurationGetterMock) getConfiguration(iss string) (configuration, error) {
	c.Calls <- &getConfigurationCall{iss}
	gr := (<-c.Calls).(*getConfigurationResponse)
	return gr.config, gr.err
}

func (c *configurationGetterMock) assertGetConfiguration(iss string, config configuration, err error) {
	call := (<-c.Calls).(*getConfigurationCall)
	if iss != anything && call.iss != iss {
		c.t.Error("Expected getConfiguration with", iss, "but was", call.iss)
	}
	c.Calls <- &getConfigurationResponse{config, err}
}

func (c *configurationGetterMock) close() {
	close(c.Calls)
}

func (c *configurationGetterMock) assertDone() {
	if _, more := <-c.Calls; more {
		c.t.Fatal("Did not expect more calls.")
	}
}
