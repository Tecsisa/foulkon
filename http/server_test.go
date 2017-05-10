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
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, worker.Host+":"+worker.Port, ws.Addr, "Error in test")
	assert.Equal(t, worker.CertFile, ws.certFile, "Error in test")
	assert.Equal(t, worker.KeyFile, ws.keyFile, "Error in test")
	assert.Equal(t, handler, ws.Handler, "Error in test")
}

func TestNewProxy(t *testing.T) {
	testApi := makeTestApi()
	testcases := map[string]struct {
		proxy *foulkon.Proxy

		getProxyResourcesMethod []api.ProxyResource
		getProxyResourcesError  error

		expectedResources []api.ProxyResource
		expectedError     string
		panicError        string
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
					ID: "ID2",
					Resource: api.ResourceEntity{
						Host:   "host2",
						Path:   "/path2",
						Method: "Method2",
						Urn:    "urn2",
						Action: "action2",
					},
				},
			},
			expectedResources: []api.ProxyResource{
				{
					ID: "ID2",
					Resource: api.ResourceEntity{
						Host:   "host2",
						Path:   "/path2",
						Method: "Method2",
						Urn:    "urn2",
						Action: "action2",
					},
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
					ID: "ID2",
					Resource: api.ResourceEntity{
						Host:   "host2",
						Path:   "",
						Method: "Method2",
						Urn:    "urn2",
						Action: "action2",
					},
				},
			},
			expectedResources: []api.ProxyResource{
				{
					ID: "ID2",
					Resource: api.ResourceEntity{
						Host:   "host2",
						Path:   "",
						Method: "Method2",
						Urn:    "urn2",
						Action: "action2",
					},
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
		"ErrorCaseDeployingRepeatedResourcePaths": {
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
					ID:   "ID1",
					Name: "name1",
					Org:  "org1",
					Resource: api.ResourceEntity{
						Host:   "host2",
						Path:   "/path/:mypath/",
						Method: "Method2",
						Urn:    "urn2",
						Action: "action2",
					},
				},
				{
					ID:   "ID2",
					Name: "name2",
					Org:  "org1",
					Resource: api.ResourceEntity{
						Host:   "host2",
						Path:   "/path/:mypath/",
						Method: "Method2",
						Urn:    "urn2",
						Action: "action2",
					},
				},
				{
					ID:   "ID3",
					Name: "name3",
					Org:  "org1",
					Resource: api.ResourceEntity{
						Host:   "host2",
						Path:   "/path/*mypath",
						Method: "Method2",
						Urn:    "urn2",
						Action: "action2",
					},
				},
				{
					ID:   "ID3",
					Name: "name3",
					Org:  "org1",
					Resource: api.ResourceEntity{
						Host:   "host2",
						Path:   "/path2/*mypath",
						Method: "Method2",
						Urn:    "urn2",
						Action: "action2",
					},
				},
			},
			expectedResources: []api.ProxyResource{
				{
					ID:   "ID1",
					Name: "name1",
					Org:  "org1",
					Resource: api.ResourceEntity{
						Host:   "host2",
						Path:   "/path/:mypath/",
						Method: "Method2",
						Urn:    "urn2",
						Action: "action2",
					},
				},
				{
					ID:   "ID2",
					Name: "name2",
					Org:  "org1",
					Resource: api.ResourceEntity{
						Host:   "host2",
						Path:   "/path/:mypath/",
						Method: "Method2",
						Urn:    "urn2",
						Action: "action2",
					},
				},
				{
					ID:   "ID3",
					Name: "name3",
					Org:  "org1",
					Resource: api.ResourceEntity{
						Host:   "host2",
						Path:   "/path/*mypath",
						Method: "Method2",
						Urn:    "urn2",
						Action: "action2",
					},
				},
				{
					ID:   "ID3",
					Name: "name3",
					Org:  "org1",
					Resource: api.ResourceEntity{
						Host:   "host2",
						Path:   "/path2/*mypath",
						Method: "Method2",
						Urn:    "urn2",
						Action: "action2",
					},
				},
			},
			panicError: "There was a problem adding proxy resource with name name3 and org org1: " +
				"path segment '*mypath' conflicts with existing wildcard ':mypath' in path '/path/*mypath'",
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
			assert.Equal(t, test.expectedError, hook.LastEntry().Message, "Error in test case %v", n)
			ps.resourceLock.Unlock()
		} else {
			// Check responses
			assert.Equal(t, test.proxy.Host+":"+test.proxy.Port, ps.Addr, "Error in test case %v", n)
			assert.Equal(t, test.proxy.CertFile, ps.certFile, "Error in test case %v", n)
			assert.Equal(t, test.proxy.KeyFile, ps.keyFile, "Error in test case %v", n)
			assert.Equal(t, test.proxy.RefreshTime, ps.refreshTime, "Error in test case %v", n)
			assert.Equal(t, test.expectedResources, ps.currentResources, "Error in test case %v", n)
			// Check if panic errors where caught
			if test.panicError != "" {
				assert.Equal(t, test.panicError, hook.LastEntry().Message, "Error in test case %v", n)
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
		assert.Equal(t, test.expectedResult, result, "Error in test case %v", n)
	}
}

func TestWorkerServer_Configuration(t *testing.T) {
	ws := WorkerServer{}
	err := ws.Configuration()
	assert.Nil(t, err, "Error in test")
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
			ok := assert.NotNil(t, err, "Error in test case %v", n)
			if ok {
				assert.Equal(t, test.expectedError, err.Error(), "Error in test case %v", n)
			}
		} else {
			assert.Equal(t, test.expectedAddr, test.ps.Addr, "Error in test case %v", n)
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
		assert.True(t, strings.Contains(err.Error(), test.expectedError), "Error in test case %v", n)
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
					ID: "ID2",
					Resource: api.ResourceEntity{
						Host:   "host2",
						Path:   "/path2",
						Method: "Method2",
						Urn:    "urn2",
						Action: "action2",
					},
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
				assert.True(t, strings.Contains(err.Error(), test.expectedError), "Error in test case %v", n)
			}
		} else {
			testApi.ArgsOut[GetProxyResourcesMethod][0] = []api.ProxyResource{
				{
					ID: "ID2",
					Resource: api.ResourceEntity{
						Host:   "host2",
						Path:   "/path2",
						Method: "Method2",
						Urn:    "urn2",
						Action: "action2",
					},
				},
			}

			go func() {
				srv.Run()
			}()

			// Wait reloadFunc
			time.Sleep(5 * time.Millisecond)

			ps := srv.(*ProxyServer)

			ps.resourceLock.Lock()
			assert.Equal(t, test.expectedResources, ps.currentResources, "Error in test case %v", n)
			ps.resourceLock.Unlock()
		}
	}
}
