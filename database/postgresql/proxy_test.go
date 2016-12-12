package postgresql

import (
	"testing"

	"time"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/database"
	"github.com/stretchr/testify/assert"
)

func TestPostgresRepo_GetProxyResources(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousResources []ProxyResource
		// Postgres Repo Args
		filter *api.Filter
		// Expected result
		expectedResponse []api.ProxyResource
		expectedError    *database.Error
	}{
		"OkCase1": {
			previousResources: []ProxyResource{
				{
					ID:          "ID",
					Name:        "name",
					Path:        "path",
					Org:         "org",
					Host:        "host",
					Url:         "/url",
					Method:      "Method",
					UrnResource: "urnr",
					Action:      "action",
					Urn:         "urn",
					CreateAt:    now.UnixNano(),
					UpdateAt:    now.UnixNano(),
				},
				{
					ID:          "ID2",
					Name:        "name2",
					Path:        "path2",
					Org:         "org2",
					Host:        "host2",
					Url:         "/url2",
					Method:      "Method2",
					UrnResource: "urnr2",
					Action:      "action2",
					Urn:         "urn2",
					CreateAt:    now.UnixNano(),
					UpdateAt:    now.UnixNano(),
				},
			},
			filter: &api.Filter{
				PathPrefix: "path",
				Offset:     0,
				Limit:      20,
				OrderBy:    "urn desc",
			},
			expectedResponse: []api.ProxyResource{
				{
					ID:   "ID2",
					Name: "name2",
					Path: "path2",
					Org:  "org2",
					Resource: api.ResourceEntity{
						Host:   "host2",
						Url:    "/url2",
						Method: "Method2",
						Urn:    "urnr2",
						Action: "action2",
					},
					Urn:      "urn2",
					CreateAt: now,
					UpdateAt: now,
				},
				{
					ID:   "ID",
					Name: "name",
					Path: "path",
					Org:  "org",
					Resource: api.ResourceEntity{
						Host:   "host",
						Url:    "/url",
						Method: "Method",
						Urn:    "urnr",
						Action: "action",
					},
					Urn:      "urn",
					CreateAt: now,
					UpdateAt: now,
				},
			},
		},
		"OkCase2": {
			previousResources: []ProxyResource{
				{
					ID:          "ID",
					Name:        "name",
					Path:        "path",
					Org:         "org",
					Host:        "host",
					Url:         "/url",
					Method:      "Method",
					UrnResource: "urnr",
					Action:      "action",
					Urn:         "urn",
					CreateAt:    now.UnixNano(),
					UpdateAt:    now.UnixNano(),
				},
				{
					ID:          "ID2",
					Name:        "name2",
					Path:        "path2",
					Org:         "org",
					Host:        "host2",
					Url:         "/url2",
					Method:      "Method2",
					UrnResource: "urnr2",
					Action:      "action2",
					Urn:         "urn2",
					CreateAt:    now.UnixNano(),
					UpdateAt:    now.UnixNano(),
				},
			},
			filter: &api.Filter{
				PathPrefix: "path",
				Org:        "org",
				Offset:     0,
				Limit:      20,
			},
			expectedResponse: []api.ProxyResource{
				{
					ID:   "ID",
					Name: "name",
					Path: "path",
					Org:  "org",
					Resource: api.ResourceEntity{
						Host:   "host",
						Url:    "/url",
						Method: "Method",
						Urn:    "urnr",
						Action: "action",
					},
					Urn:      "urn",
					CreateAt: now,
					UpdateAt: now,
				},
				{
					ID:   "ID2",
					Name: "name2",
					Path: "path2",
					Org:  "org",
					Resource: api.ResourceEntity{
						Host:   "host2",
						Url:    "/url2",
						Method: "Method2",
						Urn:    "urnr2",
						Action: "action2",
					},
					Urn:      "urn2",
					CreateAt: now,
					UpdateAt: now,
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

		// Call to repository to get proxy resources
		res, total, err := repoDB.GetProxyResources(test.filter)
		assert.Nil(t, err, "Error in test case %v", n)

		// Check total
		assert.Equal(t, total, len(test.expectedResponse), "Error in test case %v", n)

		// Check response
		assert.Equal(t, test.expectedResponse, res, "Error in test case %v", n)
	}
}

func TestPostgresRepo_GetProxyResourceByName(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousResource *ProxyResource
		// Postgres Repo Args
		org  string
		name string
		// Expected result
		expectedResponse *api.ProxyResource
		expectedError    *database.Error
	}{
		"OkCase": {
			previousResource: &ProxyResource{
				ID:          "ID",
				Name:        "name",
				Path:        "path",
				Org:         "org",
				Host:        "host",
				Url:         "/url",
				Method:      "Method",
				UrnResource: "urn2",
				Action:      "action",
				Urn:         "urn",
				CreateAt:    now.UnixNano(),
				UpdateAt:    now.UnixNano(),
			},
			org:  "org",
			name: "name",
			expectedResponse: &api.ProxyResource{
				ID:   "ID",
				Name: "name",
				Path: "path",
				Org:  "org",
				Resource: api.ResourceEntity{
					Host:   "host",
					Url:    "/url",
					Method: "Method",
					Urn:    "urn2",
					Action: "action",
				},
				Urn:      "urn",
				CreateAt: now,
				UpdateAt: now,
			},
		},
		"ErrorCaseProxyResourceNotExist": {
			previousResource: &ProxyResource{
				ID:          "ID",
				Name:        "name",
				Path:        "path",
				Org:         "org",
				Host:        "host",
				Url:         "/url",
				Method:      "Method",
				UrnResource: "urn2",
				Action:      "action",
				Urn:         "urn",
				CreateAt:    now.UnixNano(),
				UpdateAt:    now.UnixNano(),
			},
			org:  "org1",
			name: "name2",
			expectedError: &database.Error{
				Code:    database.PROXY_RESOURCE_NOT_FOUND,
				Message: "Proxy resource with organization org1 and name name2 not found",
			},
		},
	}

	for n, test := range testcases {
		// Clean proxy_resource database
		cleanProxyResourcesTable(t, n)

		// Insert previous data
		if test.previousResource != nil {
			insertProxyResource(t, n, *test.previousResource)
		}

		// Call to repository to get proxy resource
		res, err := repoDB.GetProxyResourceByName(test.org, test.name)
		if test.expectedError != nil {
			dbError, _ := err.(*database.Error)
			assert.Equal(t, test.expectedError, dbError, "Error in test case %v", n)
		} else {
			assert.Nil(t, err, "Error in test case %v", n)

			// Check response
			assert.Equal(t, test.expectedResponse, res, "Error in test case %v", n)
		}
	}
}

func TestPostgresRepo_AddProxyResource(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousResource *ProxyResource
		// Postgres Repo Args
		proxyResource *api.ProxyResource
		// Expected result
		expectedResponse *api.ProxyResource
		expectedError    *database.Error
	}{
		"OkCase": {
			proxyResource: &api.ProxyResource{
				ID:   "ID",
				Name: "name",
				Path: "path",
				Org:  "org",
				Resource: api.ResourceEntity{
					Host:   "host",
					Url:    "/url",
					Method: "Method",
					Urn:    "urn2",
					Action: "action",
				},
				Urn:      "urn",
				CreateAt: now,
				UpdateAt: now,
			},
			expectedResponse: &api.ProxyResource{
				ID:   "ID",
				Name: "name",
				Path: "path",
				Org:  "org",
				Resource: api.ResourceEntity{
					Host:   "host",
					Url:    "/url",
					Method: "Method",
					Urn:    "urn2",
					Action: "action",
				},
				Urn:      "urn",
				CreateAt: now,
				UpdateAt: now,
			},
		},
		"ErrorCaseUserAlreadyExist": {
			previousResource: &ProxyResource{
				ID:          "ID",
				Name:        "name",
				Path:        "path",
				Org:         "org",
				Host:        "host",
				Url:         "/url",
				Method:      "Method",
				UrnResource: "urn2",
				Action:      "action",
				Urn:         "urn",
				CreateAt:    now.UnixNano(),
				UpdateAt:    now.UnixNano(),
			},
			proxyResource: &api.ProxyResource{
				ID:   "ID",
				Name: "name",
				Path: "path",
				Org:  "org",
				Resource: api.ResourceEntity{
					Host:   "host",
					Url:    "/url",
					Method: "Method",
					Urn:    "urn2",
					Action: "action",
				},
				Urn:      "urn",
				CreateAt: now,
				UpdateAt: now,
			},
			expectedError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "pq: duplicate key value violates unique constraint \"proxy_resources_pkey\"",
			},
		},
	}

	for n, test := range testcases {
		// Clean proxy_resource database
		cleanProxyResourcesTable(t, n)

		// Insert previous data
		if test.previousResource != nil {
			insertProxyResource(t, n, *test.previousResource)
		}
		// Call to repository to store a proxy resource
		proxyResource, err := repoDB.AddProxyResource(*test.proxyResource)
		if test.expectedError != nil {
			dbError, _ := err.(*database.Error)
			assert.Equal(t, test.expectedError, dbError, "Error in test case %v", n)
		} else {
			assert.Nil(t, err, "Error in test case %v", n)

			// Check response
			assert.Equal(t, test.expectedResponse, proxyResource, "Error in test case %v", n)
			// Check database
			count := getProxyResourcesCountFiltered(t, n, test.expectedResponse.ID, test.expectedResponse.Name, test.expectedResponse.Org,
				test.expectedResponse.Path, test.expectedResponse.Urn, test.expectedResponse.CreateAt.UnixNano(), test.expectedResponse.UpdateAt.UnixNano())
			assert.Equal(t, 1, count, "Error in test case %v", n)
		}
	}
}
