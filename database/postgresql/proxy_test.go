package postgresql

import (
	"testing"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/database"
	"github.com/kylelemons/godebug/pretty"
)

func TestPostgresRepo_GetProxyResources(t *testing.T) {
	testcases := map[string]struct {
		// Previous data
		previousResources []ProxyResource
		// Postgres Repo Args
		proxyResources []api.ProxyResource
		// Expected result
		expectedResponse []api.ProxyResource
		expectedError    *database.Error
	}{

		"OkCase": {
			previousResources: []ProxyResource{
				{
					ID:     "ID",
					Host:   "host",
					Url:    "/url",
					Method: "Method",
					Urn:    "urn",
					Action: "action",
				},
			},
			proxyResources: []api.ProxyResource{
				{
					ID:     "ID",
					Host:   "host",
					Url:    "/url",
					Method: "Method",
					Urn:    "urn",
					Action: "action",
				},
			},
			expectedResponse: []api.ProxyResource{
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
	}

	for n, test := range testcases {
		// Clean proxy_resource database
		cleanProxyResourcesTable()

		// Insert previous data
		if test.previousResources != nil {
			for _, previousResource := range test.previousResources {
				if err := insertProxyResource(previousResource); err != nil {
					t.Errorf("Test %v failed. Unexpected error inserting previous proxy resources: %v", n, err)
					continue
				}
			}
		}

		// Call to repository to get resources
		res, err := repoDB.GetProxyResources()
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error: %v", n, err)
			continue
		}
		// Check response
		if diff := pretty.Compare(res, test.expectedResponse); diff != "" {
			t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
			continue
		}
	}
}
