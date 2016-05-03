package openid

import (
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/square/go-jose"
)

type Call interface{}

const anything = "anything"

type HTTPClientMock struct {
	t     *testing.T
	Calls chan Call
}

func NewHTTPClientMock(t *testing.T) *HTTPClientMock {
	return &HTTPClientMock{t, make(chan Call)}
}

type httpGetCall struct {
	url string
}

type httpGetResp struct {
	resp *http.Response
	err  error
}

type decodeResponseCall struct {
	reader io.Reader
}

type decodeResponseResp struct {
	value interface{}
	err   error
}

func (c *HTTPClientMock) httpGet(url string) (*http.Response, error) {
	c.Calls <- &httpGetCall{url}
	gr := (<-c.Calls).(*httpGetResp)
	return gr.resp, gr.err
}

func (c *HTTPClientMock) assertHttpGet(url string, resp *http.Response, err error) {
	call := (<-c.Calls).(*httpGetCall)
	if url != anything && call.url != url {
		c.t.Error("Expected httpGet with", url, "but was", call.url)
	}
	c.Calls <- &httpGetResp{resp, err}
}

func (c *HTTPClientMock) decodeResponse(reader io.Reader, value interface{}) error {
	c.Calls <- &decodeResponseCall{reader}
	dr := (<-c.Calls).(*decodeResponseResp)
	switch v := value.(type) {
	case *configuration:
		if dr.value != nil {
			v.Issuer = dr.value.(*configuration).Issuer
			v.JwksUri = dr.value.(*configuration).JwksUri
		}
	case *jose.JsonWebKeySet:
		if dr.value != nil {
			v.Keys = dr.value.(*jose.JsonWebKeySet).Keys
		}

	default:
		c.t.Fatalf("Expected value type '*configuration', but was %T", value)

	}

	return dr.err
}

func (c *HTTPClientMock) assertDecodeResponse(response string, value interface{}, err error) {
	call := (<-c.Calls).(*decodeResponseCall)
	if response != anything {
		b, e := ioutil.ReadAll(call.reader)
		if e != nil {
			c.t.Error("Error while reading from the call reader", e)
		}
		s := string(b)

		if s != response {
			c.t.Error("Expected decodeResponse with", response, "but was", s)
		}
	}

	c.Calls <- &decodeResponseResp{value, err}
}

func (c *HTTPClientMock) close() {
	close(c.Calls)
}

func (c *HTTPClientMock) assertDone() {
	if _, more := <-c.Calls; more {
		c.t.Fatal("Did not expect more calls.")
	}
}
