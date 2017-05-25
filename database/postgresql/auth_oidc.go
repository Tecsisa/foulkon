package postgresql

import (
	"time"

	"fmt"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/database"
	"github.com/satori/go.uuid"
)

// AUTH OIDC PROVIDER REPOSITORY IMPLEMENTATION

func (pr PostgresRepo) AddOidcProvider(oidcProvider api.OidcProvider) (*api.OidcProvider, error) {
	// Create OIDC Provider model
	oidcProviderDB := &OidcProvider{
		ID:        oidcProvider.ID,
		Name:      oidcProvider.Name,
		Path:      oidcProvider.Path,
		CreateAt:  oidcProvider.CreateAt.UnixNano(),
		UpdateAt:  oidcProvider.UpdateAt.UnixNano(),
		Urn:       oidcProvider.Urn,
		IssuerURL: oidcProvider.IssuerURL,
	}

	transaction := pr.Dbmap.Begin()

	// Create OIDC Provider
	if err := transaction.Create(oidcProviderDB).Error; err != nil {
		transaction.Rollback()
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Create OIDC Clients
	for _, oidcClientApi := range oidcProvider.OidcClients {
		// Create OIDC Client model
		oidcClientDB := &OidcClient{
			ID:             uuid.NewV4().String(),
			OidcProviderID: oidcProvider.ID,
			Name:           oidcClientApi.Name,
		}
		if err := transaction.Create(oidcClientDB).Error; err != nil {
			transaction.Rollback()
			return nil, &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: err.Error(),
			}
		}
	}

	transaction.Commit()

	// Create API OIDC Provider
	oidcProviderApi := dbOidcProviderToAPIOidcProvider(oidcProviderDB)
	oidcProviderApi.OidcClients = oidcProvider.OidcClients

	return oidcProviderApi, nil
}

func (pr PostgresRepo) GetOidcProviderByName(name string) (*api.OidcProvider, error) {
	oidcProvider := &OidcProvider{}
	query := pr.Dbmap.Where("name like ?", name).First(oidcProvider)

	// Check if OIDC Provider exists
	if query.RecordNotFound() {
		return nil, &database.Error{
			Code:    database.AUTH_OIDC_PROVIDER_NOT_FOUND,
			Message: fmt.Sprintf("OIDC Provider with name %v not found", name),
		}
	}

	// Error Handling
	if err := query.Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Retrieve associated OIDC Clients
	oidcClients := []OidcClient{}
	query = pr.Dbmap.Where("oidc_provider_id like ?", oidcProvider.ID).Find(&oidcClients)
	// Error Handling
	if err := query.Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Create API OidcProvider
	oidcProviderApi := dbOidcProviderToAPIOidcProvider(oidcProvider)
	oidcProviderApi.OidcClients = dbOidcClientsToAPIOidcClients(oidcClients)

	return oidcProviderApi, nil
}

func (pr PostgresRepo) GetOidcProvidersFiltered(filter *api.Filter) ([]api.OidcProvider, int, error) {
	return nil, 0, nil
}

func (pr PostgresRepo) UpdateOidcProvider(oidcProvider api.OidcProvider) (*api.OidcProvider, error) {
	return nil, nil
}

func (pr PostgresRepo) RemoveOidcProvider(name string) error {
	return nil
}

// PRIVATE HELPER METHODS

// Transform a OIDC Provider retrieved from db into a OIDC Provider for API
func dbOidcProviderToAPIOidcProvider(oidcProvider *OidcProvider) *api.OidcProvider {
	return &api.OidcProvider{
		ID:        oidcProvider.ID,
		Name:      oidcProvider.Name,
		Path:      oidcProvider.Path,
		CreateAt:  time.Unix(0, oidcProvider.CreateAt).UTC(),
		UpdateAt:  time.Unix(0, oidcProvider.UpdateAt).UTC(),
		Urn:       oidcProvider.Urn,
		IssuerURL: oidcProvider.IssuerURL,
	}
}

// Transform a list of OIDC clients from db into API OIDC clients
func dbOidcClientsToAPIOidcClients(oidcClients []OidcClient) []api.OidcClient {
	oidcClientsApi := make([]api.OidcClient, len(oidcClients), cap(oidcClients))
	for i, oc := range oidcClients {
		oidcClientsApi[i] = api.OidcClient{
			Name: oc.Name,
		}
	}

	return oidcClientsApi
}
