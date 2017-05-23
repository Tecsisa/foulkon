package postgresql

import "github.com/Tecsisa/foulkon/api"

// AUTH OIDC PROVIDER REPOSITORY IMPLEMENTATION

func (pr PostgresRepo) AddOidcProvider(oidcProvider api.OidcProvider) (*api.OidcProvider, error) {
	return nil, nil
}

func (pr PostgresRepo) GetOidcProviderByName(name string) (*api.OidcProvider, error) {
	return nil, nil
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
