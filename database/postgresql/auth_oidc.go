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
	var total int
	oidcProviders := []OidcProvider{}
	query := pr.Dbmap

	if len(filter.PathPrefix) > 0 {
		query = query.Where("path like ?", filter.PathPrefix+"%")
	}
	if len(filter.OrderBy) > 0 {
		query = query.Order(filter.OrderBy)
	}

	// Error handling
	if err := query.Find(&oidcProviders).Count(&total).Offset(filter.Offset).Limit(filter.Limit).Find(&oidcProviders).Error; err != nil {
		return nil, total, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Transform OIDC Providers to API
	var apiOidcProviders []api.OidcProvider
	if oidcProviders != nil {
		apiOidcProviders = make([]api.OidcProvider, len(oidcProviders), cap(oidcProviders))

		for i, op := range oidcProviders {
			oidcProvider := dbOidcProviderToAPIOidcProvider(&op)

			// Retrieve associated OIDC clients
			oidcClients := []OidcClient{}
			query = pr.Dbmap.Where("oidc_provider_id like ?", oidcProvider.ID).Find(&oidcClients)
			// Error Handling
			if err := query.Error; err != nil {
				return nil, total, &database.Error{
					Code:    database.INTERNAL_ERROR,
					Message: err.Error(),
				}
			}

			oidcProvider.OidcClients = dbOidcClientsToAPIOidcClients(oidcClients)

			// Assign OIDC Provider
			apiOidcProviders[i] = *oidcProvider
		}

	}

	return apiOidcProviders, total, nil
}

func (pr PostgresRepo) UpdateOidcProvider(oidcProvider api.OidcProvider) (*api.OidcProvider, error) {
	oidcProviderDB := OidcProvider{
		ID:        oidcProvider.ID,
		Name:      oidcProvider.Name,
		Path:      oidcProvider.Path,
		CreateAt:  oidcProvider.CreateAt.UTC().UnixNano(),
		UpdateAt:  oidcProvider.UpdateAt.UTC().UnixNano(),
		Urn:       oidcProvider.Urn,
		IssuerURL: oidcProvider.IssuerURL,
	}

	transaction := pr.Dbmap.Begin()

	// Update OIDC Provider
	if err := transaction.Model(&OidcProvider{ID: oidcProvider.ID}).Update(oidcProviderDB).Error; err != nil {
		transaction.Rollback()
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Clean old OIDC Clients
	if err := transaction.Where("oidc_provider_id like ?", oidcProvider.ID).Delete(OidcClient{}).Error; err != nil {
		transaction.Rollback()
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Create new OIDC Clients
	for _, oc := range oidcProvider.OidcClients {
		oidcClientDB := &OidcClient{
			ID:             uuid.NewV4().String(),
			OidcProviderID: oidcProvider.ID,
			Name:           oc.Name,
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

	return &oidcProvider, nil
}

func (pr PostgresRepo) RemoveOidcProvider(id string) error {
	transaction := pr.Dbmap.Begin()

	// Delete OIDC Provider
	transaction.Where("id like ?", id).Delete(&OidcProvider{})
	if err := transaction.Error; err != nil {
		transaction.Rollback()
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Delete all OIDC Clients
	transaction.Where("oidc_provider_id like ?", id).Delete(&OidcClient{})
	if err := transaction.Error; err != nil {
		transaction.Rollback()
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}

	}

	transaction.Commit()
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
