package postgresql

import (
	"testing"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/database"
	"github.com/stretchr/testify/assert"
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
		cleanProxyResourcesTable(t, n)

		// Insert previous data
		if test.previousResources != nil {
			for _, previousResource := range test.previousResources {
				insertProxyResource(t, n, previousResource)
			}
		}

		// Call to repository to get resources
		res, err := repoDB.GetProxyResources()
		assert.Nil(t, err, "Error in test case %v", n)

		// Check response
		assert.Equal(t, test.expectedResponse, res, "Error in test case %v", n)
	}
}
