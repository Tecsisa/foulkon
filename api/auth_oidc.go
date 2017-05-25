package api

import (
	"fmt"
	"net/url"
	"time"

	"github.com/Tecsisa/foulkon/database"
	"github.com/satori/go.uuid"
)

// TYPE DEFINITIONS

// Authenticator OIDC domain
type OidcProvider struct {
	ID          string       `json:"id, omitempty"`
	Name        string       `json:"name, omitempty"`
	Path        string       `json:"path, omitempty"`
	Urn         string       `json:"urn, omitempty"`
	CreateAt    time.Time    `json:"createAt, omitempty"`
	UpdateAt    time.Time    `json:"updateAt, omitempty"`
	IssuerURL   string       `json:"issuerURL, omitempty"`
	OidcClients []OidcClient `json:"oidcClients, omitempty"`
}

type OidcClient struct {
	Name string `json:"name, omitempty"`
}

func (op OidcProvider) String() string {
	return fmt.Sprintf("[id: %v, name: %v, path: %v, urn: %v, createAt: %v, updateAt: %v, issuerURL: %v, clients: %v]",
		op.ID, op.Name, op.Path, op.Urn, op.CreateAt.Format("2006-01-02 15:04:05 MST"),
		op.UpdateAt.Format("2006-01-02 15:04:05 MST"), op.IssuerURL, op.OidcClients)
}

func (op OidcClient) String() string {
	return fmt.Sprintf("name: %v", op.Name)
}

func (op OidcProvider) GetUrn() string {
	return op.Urn
}

// AUTHENTICATOR OIDC API IMPLEMENTATION

func (api WorkerAPI) AddOidcProvider(requestInfo RequestInfo, name string, path string, issuerURL string, oidcClients []string) (*OidcProvider, error) {
	// Validate fields
	if !IsValidName(name) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: name %v", name),
		}
	}
	if !IsValidPath(path) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: path %v", path),
		}

	}
	if _, err := url.ParseRequestURI(issuerURL); err != nil {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: issuerUrl %v", issuerURL),
		}
	}
	err := AreValidOidcClientNames(oidcClients)
	if err != nil {
		apiError := err.(*Error)
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: apiError.Message,
		}

	}

	oidcProvider := createOidcProvider(name, path, issuerURL, oidcClients)

	// Check restrictions
	oidcProvidersFiltered, err := api.GetAuthorizedOidcProviders(requestInfo, oidcProvider.Urn, AUTH_OIDC_ACTION_CREATE_PROVIDER, []OidcProvider{oidcProvider})
	if err != nil {
		return nil, err
	}
	if len(oidcProvidersFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, oidcProvider.Urn),
		}
	}

	// Check if OIDC provider already exists
	_, err = api.AuthOidcRepo.GetOidcProviderByName(name)

	// Check if OIDC provider could be retrieved
	if err != nil {
		// Transform to DB error
		dbError := err.(*database.Error)
		switch dbError.Code {
		// OIDC provider doesn't exist in DB
		case database.AUTH_OIDC_PROVIDER_NOT_FOUND:
			// Create OIDC provider
			createdOidcProvider, err := api.AuthOidcRepo.AddOidcProvider(oidcProvider)

			// Check if there is an unexpected error in DB
			if err != nil {
				//Transform to DB error
				dbError := err.(*database.Error)
				return nil, &Error{
					Code:    UNKNOWN_API_ERROR,
					Message: dbError.Message,
				}
			}

			LogOperation(requestInfo.RequestID, requestInfo.Identifier, fmt.Sprintf("OIDC provider created %+v", createdOidcProvider))
			return createdOidcProvider, nil
		default: // Unexpected error
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	} else { // Fail if OIDC provider exists
		return nil, &Error{
			Code:    AUTH_OIDC_PROVIDER_ALREADY_EXIST,
			Message: fmt.Sprintf("Unable to create OIDC provider, OIDC provider with name %v already exist", name),
		}
	}
}

//TODO @rsoleto: implement it!!
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

// PRIVATE HELPER METHODS

func createOidcProvider(name string, path string, issuerURL string, oidcClients []string) OidcProvider {
	urn := CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, path, name)
	oidcClientsApi := []OidcClient{}
	for _, oc := range oidcClients {
		oidcClientsApi = append(oidcClientsApi, OidcClient{Name: oc})
	}
	oidcProvider := OidcProvider{
		ID:          uuid.NewV4().String(),
		Name:        name,
		Path:        path,
		CreateAt:    time.Now().UTC(),
		UpdateAt:    time.Now().UTC(),
		Urn:         urn,
		IssuerURL:   issuerURL,
		OidcClients: oidcClientsApi,
	}

	return oidcProvider
}
