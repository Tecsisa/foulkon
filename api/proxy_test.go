package api

import (
	"testing"

	"github.com/Tecsisa/foulkon/database"
)

func TestProxyAPI_GetProxyResources(t *testing.T) {
	testcases := map[string]struct {
		wantError error

		getProxyResourcesMethod []ProxyResource
		getProxyResourcesErr    error
	}{
		"OkCase": {
			getProxyResourcesMethod: []ProxyResource{
				{
					ID:     "ID",
					Host:   "host",
					Url:    "/url",
					Method: "Method",
					Urn:    "urn",
					Action: "action",
				},
			},
		},
		"ErrorCaseInternalError": {
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getProxyResourcesErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
	}

	for n, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeProxyTestAPI(testRepo)

		testRepo.ArgsOut[GetProxyResourcesMethod][0] = testcase.getProxyResourcesMethod
		testRepo.ArgsOut[GetProxyResourcesMethod][1] = testcase.getProxyResourcesErr

		resources, err := testAPI.GetProxyResources()
		checkMethodResponse(t, n, testcase.wantError, err, testcase.getProxyResourcesMethod, resources)
	}
}
