package http

import (
	"testing"

	"crypto/tls"
	"path/filepath"

	"net/http"

	"time"

	"strings"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/foulkon"
	"github.com/julienschmidt/httprouter"
	"github.com/kylelemons/godebug/pretty"
)

func TestNewWorker(t *testing.T) {
	// Args
	worker := &foulkon.Worker{
		Host:     "host",
		Port:     "port",
		CertFile: "cert",
		KeyFile:  "key",
	}
	handler := httprouter.New()
	// Call func
	srv := NewWorker(worker, handler)
	ws := srv.(*WorkerServer)

	// Check responses
	if diff := pretty.Compare(ws.Addr, worker.Host+":"+worker.Port); diff != "" {
		t.Errorf("Test failed. Received different Addr (received/wanted) %v", diff)
		return
	}

	if diff := pretty.Compare(ws.certFile, worker.CertFile); diff != "" {
		t.Errorf("Test failed. Received different certFile (received/wanted) %v", diff)
		return
	}

	if diff := pretty.Compare(ws.keyFile, worker.KeyFile); diff != "" {
		t.Errorf("Test failed. Received different keyFile (received/wanted) %v", diff)
		return
	}

	if diff := pretty.Compare(ws.Handler, handler); diff != "" {
		t.Errorf("Test failed. Received different keyFile (received/wanted) %v", diff)
		return
	}
}

func TestNewProxy(t *testing.T) {
	testApi := makeTestApi()
	testcases := map[string]struct {
		proxy *foulkon.Proxy

		getProxyResourcesMethod []api.ProxyResource
		getProxyResourcesError  error

		expectedResources []api.ProxyResource
		expectedError     string
	}{
		"OKCase": {
			proxy: &foulkon.Proxy{
				Host:        "host",
				Port:        "port",
				CertFile:    "cert",
				KeyFile:     "key",
				RefreshTime: 10,
				ProxyApi:    testApi,
			},
			getProxyResourcesMethod: []api.ProxyResource{
				{
					ID:     "ID2",
					Host:   "host2",
					Url:    "/url2",
					Method: "Method2",
					Urn:    "urn2",
					Action: "action2",
				},
			},
			expectedResources: []api.ProxyResource{
				{
					ID:     "ID2",
					Host:   "host2",
					Url:    "/url2",
					Method: "Method2",
					Urn:    "urn2",
					Action: "action2",
				},
			},
		},
		"OKCaseEmptyResourceURL": {
			proxy: &foulkon.Proxy{
				Host:        "host",
				Port:        "port",
				CertFile:    "cert",
				KeyFile:     "key",
				RefreshTime: 10,
				ProxyApi:    testApi,
			},
			getProxyResourcesMethod: []api.ProxyResource{
				{
					ID:     "ID2",
					Host:   "host2",
					Url:    "",
					Method: "Method2",
					Urn:    "urn2",
					Action: "action2",
				},
			},
			expectedResources: []api.ProxyResource{
				{
					ID:     "ID2",
					Host:   "host2",
					Url:    "",
					Method: "Method2",
					Urn:    "urn2",
					Action: "action2",
				},
			},
		},
		"OKCaseEmptyResources": {
			proxy: &foulkon.Proxy{
				Host:        "host",
				Port:        "port",
				CertFile:    "cert",
				KeyFile:     "key",
				RefreshTime: 10,
				ProxyApi:    testApi,
			},
			getProxyResourcesMethod: []api.ProxyResource{},
			expectedResources:       []api.ProxyResource{},
		},
		"ErrorCaseGetProxyResources": {
			proxy: &foulkon.Proxy{
				Host:        "host",
				Port:        "port",
				CertFile:    "cert",
				KeyFile:     "key",
				RefreshTime: 10,
				ProxyApi:    testApi,
			},
			getProxyResourcesError: api.Error{
				Code:    INTERNAL_SERVER_ERROR,
				Message: "Unknow error",
			},
			expectedError: "Unexpected error reading proxy resources from database Code: InternalServerError, Message: Unknow error",
		},
	}
	for n, test := range testcases {
		testApi.ArgsOut[GetProxyResourcesMethod][0] = test.getProxyResourcesMethod
		testApi.ArgsOut[GetProxyResourcesMethod][1] = test.getProxyResourcesError

		// Call func
		srv := NewProxy(test.proxy)
		ps := srv.(*ProxyServer)

		if test.expectedError != "" {
			// Check error
			ps := srv.(*ProxyServer)
			ps.resourceLock.Lock()
			if diff := pretty.Compare(test.expectedError, hook.LastEntry().Message); diff != "" {
				t.Errorf("Test %v failed. Received different errors (received/wanted) %v", n, diff)
				ps.resourceLock.Unlock()
				continue
			}
			ps.resourceLock.Unlock()
		} else {
			// Check responses
			if diff := pretty.Compare(ps.Addr, test.proxy.Host+":"+test.proxy.Port); diff != "" {
				t.Errorf("Test %v failed. Received different Addr (received/wanted) %v", n, diff)
				continue
			}

			if diff := pretty.Compare(ps.certFile, test.proxy.CertFile); diff != "" {
				t.Errorf("Test %v failed. Received different certFile (received/wanted) %v", n, diff)
				continue
			}

			if diff := pretty.Compare(ps.keyFile, test.proxy.KeyFile); diff != "" {
				t.Errorf("Test %v failed. Received different keyFile (received/wanted) %v", n, diff)
				continue
			}

			if diff := pretty.Compare(ps.refreshTime, test.proxy.RefreshTime); diff != "" {
				t.Errorf("Test %v failed. Received different refreshTime (received/wanted) %v", n, diff)
				continue
			}

			if diff := pretty.Compare(ps.currentResources, test.expectedResources); diff != "" {
				t.Errorf("Test %v failed. Received different resources (received/wanted) %v", n, diff)
				continue
			}
		}
	}
}

func Test_strSliceContains(t *testing.T) {
	testcases := map[string]struct {
		ss             []string
		s              string
		expectedResult bool
	}{
		"OKcaseTrue": {
			ss:             []string{"a", "b"},
			s:              "a",
			expectedResult: true,
		},
		"OKcaseFalse": {
			ss:             []string{"a", "b"},
			s:              "c",
			expectedResult: false,
		},
	}

	for n, test := range testcases {
		result := strSliceContains(test.ss, test.s)
		if result != test.expectedResult {
			t.Errorf("Test %v failed. Received different responses (received/wanted)", n)
			continue
		}
	}
}

func TestWorkerServer_Configuration(t *testing.T) {
	ws := WorkerServer{}
	err := ws.Configuration()
	if diff := pretty.Compare(err, nil); diff != "" {
		t.Errorf("Test failed. Received different errors (received/wanted) %v", diff)
		return
	}
}

func TestProxyServer_Configuration(t *testing.T) {
	certFile, _ := filepath.Abs("../dist/test/cert.pem")
	keyFile, _ := filepath.Abs("../dist/test/key.pem")
	testcases := map[string]struct {
		ps *ProxyServer

		expectedAddr  string
		expectedError string
	}{
		"OKcase": {
			ps: &ProxyServer{
				certFile: "",
				keyFile:  "",
			},
			expectedAddr: ":http",
		},
		"OKcaseTLS": {
			ps: &ProxyServer{
				certFile: certFile,
				keyFile:  keyFile,
			},
			expectedAddr: ":https",
		},
		"ErrorCaseTLS": {
			ps: &ProxyServer{
				certFile: certFile,
				keyFile:  "",
			},
			expectedError: "open : no such file or directory",
		},
	}

	for n, test := range testcases {
		var err error
		test.ps.TLSConfig = &tls.Config{}
		err = test.ps.Configuration()
		if test.expectedError != "" {
			if err != nil {
				if diff := pretty.Compare(err.Error(), test.expectedError); diff != "" {
					t.Errorf("Test %v failed. Received different errors (received/wanted) %v", n, diff)
					continue
				}
			} else {
				t.Errorf("Test %v failed. No errors received", n)
				continue
			}
		} else {
			if diff := pretty.Compare(test.ps.Addr, test.expectedAddr); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
		}
	}
}

func TestWorkerServer_Run(t *testing.T) {
	certFile, _ := filepath.Abs("../dist/test/cert.pem")
	keyFile, _ := filepath.Abs("../dist/test/key.pem")
	testcases := map[string]struct {
		worker  *foulkon.Worker
		handler http.Handler

		expectedError string
	}{
		"ErrorCaseListen": {
			worker: &foulkon.Worker{
				Host: "fail",
			},
			handler:       httprouter.New(),
			expectedError: "listen tcp: lookup fail",
		},
		"ErrorCaseListenTLS": {
			worker: &foulkon.Worker{
				Host:     "fail",
				CertFile: certFile,
				KeyFile:  keyFile,
			},
			handler:       httprouter.New(),
			expectedError: "listen tcp: lookup fail",
		},
	}

	for n, test := range testcases {
		var err error
		srv := NewWorker(test.worker, test.handler)
		err = srv.Run()
		if !strings.Contains(err.Error(), test.expectedError) {
			t.Errorf("Test %v failed. Received different errors (received: %v / wanted: %v)", n, test.expectedError, err.Error())
			continue
		}
	}
}

func TestProxyServer_Run(t *testing.T) {
	testApi := makeTestApi()
	testcases := map[string]struct {
		proxy *foulkon.Proxy

		expectedResources []api.ProxyResource
		expectedError     string
	}{
		"OKCase": {
			proxy: &foulkon.Proxy{
				RefreshTime: 1 * time.Millisecond,
				ProxyApi:    testApi,
			},
			expectedResources: []api.ProxyResource{
				{
					ID:     "ID2",
					Host:   "host2",
					Url:    "/url2",
					Method: "Method2",
					Urn:    "urn2",
					Action: "action2",
				},
			},
		},
		"ErrorCaseListen": {
			proxy: &foulkon.Proxy{
				Host:        "fail",
				Port:        "53",
				RefreshTime: 1 * time.Millisecond,
				ProxyApi:    testApi,
			},
			expectedError: "listen tcp: lookup fail",
		},
	}
	for n, test := range testcases {
		var err error
		srv := NewProxy(test.proxy)
		srv.Configuration()

		if test.expectedError != "" {
			err = srv.Run()
			if err != nil {
				if !strings.Contains(err.Error(), test.expectedError) {
					t.Errorf("Test %v failed. Received different errors (received: %v / wanted: %v)", n, test.expectedError, err.Error())
					continue
				}
			}
		} else {
			testApi.ArgsOut[GetProxyResourcesMethod][0] = []api.ProxyResource{
				{
					ID:     "ID2",
					Host:   "host2",
					Url:    "/url2",
					Method: "Method2",
					Urn:    "urn2",
					Action: "action2",
				},
			}

			go func() {
				srv.Run()
			}()

			// Wait reloadFunc
			time.Sleep(5 * time.Millisecond)

			ps := srv.(*ProxyServer)

			ps.resourceLock.Lock()
			if diff := pretty.Compare(ps.currentResources, test.expectedResources); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
			ps.resourceLock.Unlock()
		}
	}
}
