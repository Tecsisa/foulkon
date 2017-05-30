package postgresql

import (
	"testing"
	"time"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/database"
	"github.com/stretchr/testify/assert"
)

func TestPostgresRepo_AddOidcProvider(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousOidcProviders []OidcProvider
		previousOidcClients   []OidcClient
		// Postgres Repo Args
		oidcProviderToCreate *api.OidcProvider
		// Expected result
		expectedResponse *api.OidcProvider
		expectedError    *database.Error
	}{
		"OkCase": {
			oidcProviderToCreate: &api.OidcProvider{
				ID:        "OIDCProviderID",
				Name:      "Name",
				Path:      "Path",
				Urn:       "urn",
				CreateAt:  now,
				UpdateAt:  now,
				IssuerURL: "",
				OidcClients: []api.OidcClient{
					{
						Name: "client1",
					},
					{
						Name: "client2",
					},
					{
						Name: "client3",
					},
				},
			},
			expectedResponse: &api.OidcProvider{
				ID:        "OIDCProviderID",
				Name:      "Name",
				Path:      "Path",
				Urn:       "urn",
				CreateAt:  now,
				UpdateAt:  now,
				IssuerURL: "",
				OidcClients: []api.OidcClient{
					{
						Name: "client1",
					},
					{
						Name: "client2",
					},
					{
						Name: "client3",
					},
				},
			},
		},
		"ErrorCaseAlreadyExists": {
			previousOidcProviders: []OidcProvider{
				{
					ID:        "OIDCProviderID",
					Name:      "Name",
					Path:      "Path",
					Urn:       "urn",
					CreateAt:  now.UnixNano(),
					UpdateAt:  now.UnixNano(),
					IssuerURL: "",
				},
			},
			previousOidcClients: []OidcClient{
				{
					ID:             "ID1",
					OidcProviderID: "OIDCProviderID",
					Name:           "client1",
				},
				{
					ID:             "ID2",
					OidcProviderID: "OIDCProviderID",
					Name:           "client2",
				},
				{
					ID:             "ID3",
					OidcProviderID: "OIDCProviderID",
					Name:           "client3",
				},
			},
			oidcProviderToCreate: &api.OidcProvider{
				ID:        "OIDCProviderID",
				Name:      "Name",
				Path:      "Path",
				Urn:       "urn",
				CreateAt:  now,
				UpdateAt:  now,
				IssuerURL: "",
				OidcClients: []api.OidcClient{
					{
						Name: "client1",
					},
					{
						Name: "client2",
					},
					{
						Name: "client3",
					},
				},
			},
			expectedError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "pq: duplicate key value violates unique constraint \"oidc_providers_pkey\"",
			},
		},
		"ErrorCaseOIDCClientAlreadyExists": {
			previousOidcProviders: []OidcProvider{
				{
					ID:        "OIDCProviderID2",
					Name:      "Name2",
					Path:      "Path2",
					Urn:       "urn2",
					CreateAt:  now.UnixNano(),
					UpdateAt:  now.UnixNano(),
					IssuerURL: "",
				},
			},
			previousOidcClients: []OidcClient{
				{
					ID:             "ID1",
					OidcProviderID: "OIDCProviderID2",
					Name:           "client1",
				},
				{
					ID:             "ID2",
					OidcProviderID: "OIDCProviderID",
					Name:           "client1",
				},
				{
					ID:             "ID3",
					OidcProviderID: "OIDCProviderID2",
					Name:           "client3",
				},
			},
			oidcProviderToCreate: &api.OidcProvider{
				ID:        "OIDCProviderID",
				Name:      "Name",
				Path:      "Path",
				Urn:       "urn",
				CreateAt:  now,
				UpdateAt:  now,
				IssuerURL: "",
				OidcClients: []api.OidcClient{
					{
						Name: "client1",
					},
					{
						Name: "client2",
					},
					{
						Name: "client3",
					},
				},
			},
			expectedError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "pq: duplicate key value violates unique constraint \"idx_oidc_client\"",
			},
		},
	}
	for n, test := range testcases {
		// Clean OIDC Provider databases
		cleanOidcClientsTable(t, n)
		cleanOidcProvidersTable(t, n)

		// Insert previous data
		for _, op := range test.previousOidcProviders {
			insertOidcProvider(t, n, op, nil)
		}
		for _, oc := range test.previousOidcClients {
			insertOidcClient(t, n, oc)
		}
		// Call to repository to store the OIDC Provider
		storedOidcProvider, err := repoDB.AddOidcProvider(*test.oidcProviderToCreate)
		if test.expectedError != nil {
			dbError, _ := err.(*database.Error)
			assert.Equal(t, dbError, test.expectedError, "Error in test case %v", n)
		} else {
			assert.Nil(t, err, "Error in test case %v", n)
			// Check response
			assert.Equal(t, storedOidcProvider, test.expectedResponse, "Error in test case %v", n)
			// Check database
			oidcProviderNumber := getOidcProvidersCountFiltered(t, n, test.oidcProviderToCreate.ID, test.oidcProviderToCreate.Name, test.oidcProviderToCreate.Path,
				test.oidcProviderToCreate.CreateAt.UnixNano(), test.oidcProviderToCreate.UpdateAt.UnixNano(),
				test.oidcProviderToCreate.Urn, test.oidcProviderToCreate.IssuerURL)

			if oidcProviderNumber != 1 {
				t.Errorf("Test %v failed. Received different oidc providers number: %v", n, oidcProviderNumber)
				continue
			}
		}
	}
}

func TestPostgresRepo_GetOidcProviderByName(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		oidcProvider *OidcProvider
		oidcClients  []OidcClient
		// Postgres Repo Args
		name string
		// Expected result
		expectedResponse *api.OidcProvider
		expectedError    *database.Error
	}{
		"OkCase": {
			name: "test",
			oidcProvider: &OidcProvider{
				ID:        "1234",
				Name:      "test",
				Path:      "/path/",
				CreateAt:  now.UnixNano(),
				UpdateAt:  now.UnixNano(),
				Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "test"),
				IssuerURL: "https://test.com",
			},
			oidcClients: []OidcClient{
				{
					ID:             "0123",
					OidcProviderID: "1234",
					Name:           "client1",
				},
			},
			expectedResponse: &api.OidcProvider{
				ID:        "1234",
				Name:      "test",
				Path:      "/path/",
				CreateAt:  now,
				UpdateAt:  now,
				Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "test"),
				IssuerURL: "https://test.com",
				OidcClients: []api.OidcClient{
					{
						Name: "client1",
					},
				},
			},
		},
		"ErrorCaseNotFound": {
			name: "test",
			expectedError: &database.Error{
				Code:    database.AUTH_OIDC_PROVIDER_NOT_FOUND,
				Message: "OIDC Provider with name test not found",
			},
		},
	}

	for n, test := range testcases {
		// Clean OIDC Provider database
		cleanOidcProvidersTable(t, n)
		cleanOidcClientsTable(t, n)

		// Insert previous data
		if test.oidcProvider != nil {
			insertOidcProvider(t, n, *test.oidcProvider, test.oidcClients)
		}
		// Call to repository to get a OIDC Provider
		receivedOidcProvider, err := repoDB.GetOidcProviderByName(test.name)
		if test.expectedError != nil {
			dbError, _ := err.(*database.Error)
			assert.Equal(t, test.expectedError, dbError, "Error in test case %v", n)
		} else {
			assert.Nil(t, err, "Error in test case %v", n)
			// Check response
			assert.Equal(t, test.expectedResponse, receivedOidcProvider, "Error in test case %v", n)
		}
	}
}

func TestPostgresRepo_GetOidcProvidersFiltered(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		oidcProviders []OidcProvider
		oidcClients   []OidcClient
		// Postgres Repo Args
		filter *api.Filter
		// Expected result
		expectedResponse []api.OidcProvider
		expectedError    *database.Error
	}{
		"OkCase": {
			filter: &api.Filter{
				PathPrefix: "/",
				Offset:     0,
				Limit:      20,
				OrderBy:    "name desc",
			},
			oidcProviders: []OidcProvider{
				{
					ID:        "111",
					Name:      "test1",
					Path:      "/path1/",
					CreateAt:  now.UnixNano(),
					UpdateAt:  now.UnixNano(),
					Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/path1/", "test1"),
					IssuerURL: "http://test1.com",
				},
				{
					ID:        "222",
					Name:      "test2",
					Path:      "/path2/",
					CreateAt:  now.UnixNano(),
					UpdateAt:  now.UnixNano(),
					Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/path2/", "test2"),
					IssuerURL: "http://test2.com",
				},
			},
			oidcClients: []OidcClient{
				{
					ID:             "1",
					Name:           "client1",
					OidcProviderID: "111",
				},
				{
					ID:             "2",
					Name:           "client2",
					OidcProviderID: "222",
				},
			},
			expectedResponse: []api.OidcProvider{
				{
					ID:        "222",
					Name:      "test2",
					Path:      "/path2/",
					CreateAt:  now,
					UpdateAt:  now,
					Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/path2/", "test2"),
					IssuerURL: "http://test2.com",
					OidcClients: []api.OidcClient{
						{
							Name: "client2",
						},
					},
				},
				{
					ID:        "111",
					Name:      "test1",
					Path:      "/path1/",
					CreateAt:  now,
					UpdateAt:  now,
					Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/path1/", "test1"),
					IssuerURL: "http://test1.com",
					OidcClients: []api.OidcClient{
						{
							Name: "client1",
						},
					},
				},
			},
		},
		"OKCaseNotFound": {
			filter: &api.Filter{
				PathPrefix: "test",
				Offset:     0,
				Limit:      20,
			},
			expectedResponse: []api.OidcProvider{},
		},
		"ErrorCaseInvalidColumnToOrder": {
			filter: &api.Filter{
				PathPrefix: "/",
				Offset:     0,
				Limit:      20,
				OrderBy:    "nocolumn desc",
			},
			oidcProviders: []OidcProvider{
				{
					ID:        "111",
					Name:      "test1",
					Path:      "/path1/",
					CreateAt:  now.UnixNano(),
					UpdateAt:  now.UnixNano(),
					Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/path1/", "test1"),
					IssuerURL: "http://test1.com",
				},
				{
					ID:        "222",
					Name:      "test2",
					Path:      "/path2/",
					CreateAt:  now.UnixNano(),
					UpdateAt:  now.UnixNano(),
					Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/path2/", "test2"),
					IssuerURL: "http://test2.com",
				},
			},
			oidcClients: []OidcClient{
				{
					ID:             "1",
					Name:           "client1",
					OidcProviderID: "111",
				},
				{
					ID:             "2",
					Name:           "client2",
					OidcProviderID: "222",
				},
			},
			expectedError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "pq: column \"nocolumn\" does not exist",
			},
		},
	}

	for n, test := range testcases {
		// Clean OIDC Provider database
		cleanOidcProvidersTable(t, n)
		cleanOidcClientsTable(t, n)

		// Insert previous data
		for i, oidcProvider := range test.oidcProviders {
			var oidcClient []OidcClient = []OidcClient{test.oidcClients[i]}
			insertOidcProvider(t, n, oidcProvider, oidcClient)
		}
		// Call repository to get OIDC Providers
		receivedOidcProviders, total, err := repoDB.GetOidcProvidersFiltered(test.filter)

		if test.expectedError != nil {
			dbError, _ := err.(*database.Error)
			assert.Equal(t, test.expectedError, dbError, "Error in test case %v", n)
		} else {
			// Check error
			assert.Nil(t, err, "Error in test case %v", n)

			// Check total
			assert.Equal(t, len(test.expectedResponse), total, "Error in test case %v", n)

			// Check response
			assert.Equal(t, test.expectedResponse, receivedOidcProviders, "Error in test case %v", n)
		}

	}
}

func TestPostgresRepo_UpdateOidcProvider(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		previousOidcProviders []OidcProvider
		previousOidcClients   []OidcClient
		oidcProvider          api.OidcProvider
		// Expected result
		expectedResponse *api.OidcProvider
		expectedError    *database.Error
	}{
		"OkCase": {
			previousOidcProviders: []OidcProvider{
				{
					ID:        "111",
					Name:      "test1",
					Path:      "/path1/",
					CreateAt:  now.UnixNano(),
					UpdateAt:  now.UnixNano(),
					Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/path1/", "test1"),
					IssuerURL: "http://test1.com",
				},
			},
			previousOidcClients: []OidcClient{
				{
					ID:             "1",
					Name:           "client1",
					OidcProviderID: "111",
				},
			},
			oidcProvider: api.OidcProvider{
				ID:        "111",
				Name:      "newName",
				Path:      "/newPath/",
				CreateAt:  now,
				UpdateAt:  now,
				Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/newPath/", "newName"),
				IssuerURL: "http://testNew.com",
				OidcClients: []api.OidcClient{
					{
						Name: "clientNew",
					},
					{
						Name: "clientNew2",
					},
				},
			},
			expectedResponse: &api.OidcProvider{
				ID:        "111",
				Name:      "newName",
				Path:      "/newPath/",
				CreateAt:  now,
				UpdateAt:  now,
				Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/newPath/", "newName"),
				IssuerURL: "http://testNew.com",
				OidcClients: []api.OidcClient{
					{
						Name: "clientNew",
					},
					{
						Name: "clientNew2",
					},
				},
			},
		},
		"ErrorCaseDuplicateUrn": {
			previousOidcProviders: []OidcProvider{
				{
					ID:        "222",
					Name:      "test2",
					Path:      "/path2/",
					CreateAt:  now.UnixNano(),
					UpdateAt:  now.UnixNano(),
					Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/path2/", "test2"),
					IssuerURL: "http://test2.com",
				},
				{
					ID:        "111",
					Name:      "test1",
					Path:      "/path1/",
					CreateAt:  now.UnixNano(),
					UpdateAt:  now.UnixNano(),
					Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/path1/", "test1"),
					IssuerURL: "http://test1.com",
				},
			},
			oidcProvider: api.OidcProvider{
				ID:        "222",
				Name:      "test1",
				Path:      "/path1/",
				CreateAt:  now,
				UpdateAt:  now,
				Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/path1/", "test1"),
				IssuerURL: "http://test1.com",
				OidcClients: []api.OidcClient{
					{
						Name: "clientNew",
					},
					{
						Name: "clientNew2",
					},
				},
			},
			expectedError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "pq: duplicate key value violates unique constraint \"oidc_providers_urn_key\"",
			},
		},
		"ErrorCaseClientDuplicated": {
			previousOidcProviders: []OidcProvider{
				{
					ID:        "222",
					Name:      "test2",
					Path:      "/path2/",
					CreateAt:  now.UnixNano(),
					UpdateAt:  now.UnixNano(),
					Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/path2/", "test2"),
					IssuerURL: "http://test2.com",
				},
				{
					ID:        "111",
					Name:      "test1",
					Path:      "/path1/",
					CreateAt:  now.UnixNano(),
					UpdateAt:  now.UnixNano(),
					Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/path1/", "test1"),
					IssuerURL: "http://test1.com",
				},
			},
			oidcProvider: api.OidcProvider{
				ID:        "333",
				Name:      "test3",
				Path:      "/path3/",
				CreateAt:  now,
				UpdateAt:  now,
				Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/path3/", "test3"),
				IssuerURL: "http://test3.com",
				OidcClients: []api.OidcClient{
					{
						Name: "clientNew",
					},
					{
						Name: "clientNew",
					},
				},
			},
			expectedError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "pq: duplicate key value violates unique constraint \"idx_oidc_client\"",
			},
		},
	}

	for n, test := range testcases {
		// Clean OIDC Provider database
		cleanOidcProvidersTable(t, n)
		cleanOidcClientsTable(t, n)

		// Call to repository to add the OIDC Providers
		if test.previousOidcProviders != nil {
			for _, op := range test.previousOidcProviders {
				insertOidcProvider(t, n, op, test.previousOidcClients)
			}
		}
		receivedOidcProvider, err := repoDB.UpdateOidcProvider(test.oidcProvider)
		if test.expectedError != nil {
			dbError, _ := err.(*database.Error)
			assert.Equal(t, test.expectedError, dbError, "Error in test case %v", n)
		} else {
			assert.Nil(t, err, "Error in test case %v", n)
			// Check response
			assert.Equal(t, receivedOidcProvider, test.expectedResponse, "Error in test case %v", n)
			// Check database
			oidcProviderNumber := getOidcProvidersCountFiltered(t, n, test.expectedResponse.ID, test.expectedResponse.Name, test.expectedResponse.Path,
				test.expectedResponse.CreateAt.UnixNano(), test.expectedResponse.UpdateAt.UnixNano(), test.expectedResponse.Urn,
				test.expectedResponse.IssuerURL)
			assert.Equal(t, 1, oidcProviderNumber, "Error in test case %v", n)
			oidcClientNumber := getOidcClientsCountFiltered(t, n, "", test.expectedResponse.ID, "")
			assert.Equal(t, 2, oidcClientNumber, "Error in test case %v", n)
		}
	}
}

func TestPostgresRepo_RemoveOidcProvider(t *testing.T) {

	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousOidcProviders []OidcProvider
		previousOidcClients   []OidcClient
		oidcProviderToDelete  string
		// Expected result
		expectedError *database.Error
	}{
		"OkCase": {
			previousOidcProviders: []OidcProvider{
				{
					ID:        "111",
					Name:      "test1",
					Path:      "/path1/",
					CreateAt:  now.UnixNano(),
					UpdateAt:  now.UnixNano(),
					Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/path1/", "test1"),
					IssuerURL: "http://test1.com",
				},
				{
					ID:        "222",
					Name:      "test2",
					Path:      "/path2/",
					CreateAt:  now.UnixNano(),
					UpdateAt:  now.UnixNano(),
					Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/path2/", "test2"),
					IssuerURL: "http://test2.com",
				},
			},
			previousOidcClients: []OidcClient{
				{
					ID:             "1",
					Name:           "client1",
					OidcProviderID: "111",
				},
				{
					ID:             "2",
					Name:           "client2",
					OidcProviderID: "222",
				},
			},
			oidcProviderToDelete: "111",
		},
	}

	for n, test := range testcases {
		// Clean OIDC Provider database
		cleanOidcProvidersTable(t, n)
		cleanOidcClientsTable(t, n)

		// Insert previous data
		if test.previousOidcProviders != nil {
			for _, op := range test.previousOidcProviders {
				insertOidcProvider(t, n, op, nil)
			}
		}

		if test.previousOidcClients != nil {
			for _, oc := range test.previousOidcClients {
				insertOidcClient(t, n, oc)
			}
		}

		// Call to repository to remove OIDC Provider
		err := repoDB.RemoveOidcProvider(test.oidcProviderToDelete)
		assert.Nil(t, err, "Error in test case %v", n)

		// Check database
		oidcProviderNumber := getOidcProvidersCountFiltered(t, n, test.oidcProviderToDelete, "", "", 0, 0, "", "")
		assert.Equal(t, 0, oidcProviderNumber, "Error in test case %v", n)

		// Check total OIDC Providers
		totalOidcProviderNumber := getOidcProvidersCountFiltered(t, n, "", "", "", 0, 0, "", "")
		assert.Equal(t, 1, totalOidcProviderNumber, "Error in test case %v", n)

		// Check OIDC Clients
		relations := getOidcClientsCountFiltered(t, n, "", test.oidcProviderToDelete, "")
		assert.Equal(t, 0, relations, "Error in test case %v", n)

		// Check total OIDC Clients
		totalRelations := getOidcClientsCountFiltered(t, n, "", "", "")
		assert.Equal(t, 1, totalRelations, "Error in test case %v", n)
	}
}

func Test_dbOidcProviderToAPIOidcProvider(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		dbOidcProvider  *OidcProvider
		apiOidcProvider *api.OidcProvider
	}{
		"OkCase": {
			dbOidcProvider: &OidcProvider{
				ID:        "test1",
				Name:      "test",
				Path:      "/path/",
				CreateAt:  now.UnixNano(),
				UpdateAt:  now.UnixNano(),
				Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "test"),
				IssuerURL: "http://example.com",
			},
			apiOidcProvider: &api.OidcProvider{
				ID:        "test1",
				Name:      "test",
				Path:      "/path/",
				CreateAt:  now,
				UpdateAt:  now,
				Urn:       api.CreateUrn("", api.RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "test"),
				IssuerURL: "http://example.com",
			},
		},
	}

	for n, test := range testcases {
		receivedAPIOidcProvider := dbOidcProviderToAPIOidcProvider(test.dbOidcProvider)
		// Check response
		assert.Equal(t, test.apiOidcProvider, receivedAPIOidcProvider, "Error in test case %v", n)
	}
}

func Test_dbOidcClientsToAPIOidcClients(t *testing.T) {
	testcases := map[string]struct {
		dbOidcClients  []OidcClient
		apiOidcClients []api.OidcClient
	}{
		"OkCase": {
			dbOidcClients: []OidcClient{
				{
					ID:             "0123",
					OidcProviderID: "0123",
					Name:           "client1",
				},
			},
			apiOidcClients: []api.OidcClient{
				{
					Name: "client1",
				},
			},
		},
		"OkCase2": {
			dbOidcClients: []OidcClient{
				{
					ID:             "0123",
					OidcProviderID: "0123",
					Name:           "client1",
				},
				{
					ID:             "01234",
					OidcProviderID: "0123",
					Name:           "client2",
				},
			},
			apiOidcClients: []api.OidcClient{
				{
					Name: "client1",
				},
				{
					Name: "client2",
				},
			},
		},
	}

	for n, test := range testcases {
		receivedAPIOidcClients := dbOidcClientsToAPIOidcClients(test.dbOidcClients)
		// Check response
		assert.Equal(t, test.apiOidcClients, receivedAPIOidcClients, "Error in test case %v", n)
	}
}
