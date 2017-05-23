package api

import "time"

// TYPE DEFINITIONS

// Authenticator OIDC domain
type OidcProvider struct {
	ID          string       `json:"id, omitempty"`
	Name        string       `json:"name, omitempty"`
	Path        string       `json:"path, omitempty"`
	Urn         string       `json:"urn, omitempty"`
	CreateAt    time.Time    `json:"createAt, omitempty"`
	UpdateAt    time.Time    `json:"updateAt, omitempty"`
	IssuerURL   string       `json:"issuerUrl, omitempty"`
	OidcClients []OidcClient `json:"oidcClients, omitempty"`
}

type OidcClient struct {
	ID   string `json:"id, omitempty"`
	Name string `json:"name, omitempty"`
}

// AUTHENTICATOR OIDC API IMPLEMENTATION TODO @rsoleto: implement it!!

func (api WorkerAPI) AddOidcProvider(requestInfo RequestInfo, name string, path string, oidcClients []string) (*OidcProvider, error) {
	return nil, nil
}

func (api WorkerAPI) GetOidcProviderByName(requestInfo RequestInfo, name string) (*OidcProvider, error) {
	return nil, nil
}

func (api WorkerAPI) ListOidcProviders(requestInfo RequestInfo, filter *Filter) ([]OidcProvider, int, error) {
	return nil, 0, nil
}

func (api WorkerAPI) UpdateOidcProvider(requestInfo RequestInfo, oidcProviderName string, newName string, newPath string,
	newClients []string) (*OidcProvider, error) {
	return nil, nil
}

func (api WorkerAPI) RemoveOidcProvider(requestInfo RequestInfo, name string) error {
	return nil
}
