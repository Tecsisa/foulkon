package postgresql

import (
	"testing"
	"time"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/database"
	"github.com/stretchr/testify/assert"
)

func TestPostgresRepo_AddOidcProvider(t *testing.T) {
	type OidcPoviderDbWithClients struct {
		provider *OidcProvider
		clients  []OidcClient
	}
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousOidcProvider *OidcPoviderDbWithClients
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
			previousOidcProvider: &OidcPoviderDbWithClients{
				provider: &OidcProvider{
					ID:        "OIDCProviderID",
					Name:      "Name",
					Path:      "Path",
					Urn:       "urn",
					CreateAt:  now.UnixNano(),
					UpdateAt:  now.UnixNano(),
					IssuerURL: "",
				},
				clients: []OidcClient{
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
	}
	for n, test := range testcases {
		// Clean OIDC Provider databases
		cleanOidcClientsTable(t, n)
		cleanOidcProvidersTable(t, n)

		// Insert previous data
		if test.previousOidcProvider != nil {
			insertOidcProvider(t, n, *test.previousOidcProvider.provider, test.previousOidcProvider.clients)
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
		// Clean policy database
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
				Urn:       api.CreateUrn("123", api.RESOURCE_POLICY, "/path/", "test"),
				IssuerURL: "http://example.com",
			},
			apiOidcProvider: &api.OidcProvider{
				ID:        "test1",
				Name:      "test",
				Path:      "/path/",
				CreateAt:  now,
				UpdateAt:  now,
				Urn:       api.CreateUrn("123", api.RESOURCE_POLICY, "/path/", "test"),
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
