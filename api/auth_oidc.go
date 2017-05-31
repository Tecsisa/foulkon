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
	ID          string       `json:"id,omitempty"`
	Name        string       `json:"name,omitempty"`
	Path        string       `json:"path,omitempty"`
	Urn         string       `json:"urn,omitempty"`
	CreateAt    time.Time    `json:"createAt,omitempty"`
	UpdateAt    time.Time    `json:"updateAt,omitempty"`
	IssuerURL   string       `json:"issuerUrl,omitempty"`
	OidcClients []OidcClient `json:"clients,omitempty"`
}

type OidcClient struct {
	Name string `json:"name,omitempty"`
}

func (op OidcProvider) String() string {
	return fmt.Sprintf("[id: %v, name: %v, path: %v, urn: %v, createAt: %v, updateAt: %v, issuerUrl: %v, clients: %v]",
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

func (api WorkerAPI) GetOidcProviderByName(requestInfo RequestInfo, name string) (*OidcProvider, error) {
	// Validate fields
	if !IsValidName(name) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: name %v", name),
		}
	}

	// Call repo to retrieve the OIDC Provider
	oidcProvider, err := api.AuthOidcRepo.GetOidcProviderByName(name)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		// OIDC Provider doesn't exist in DB
		if dbError.Code == database.AUTH_OIDC_PROVIDER_NOT_FOUND {
			return nil, &Error{
				Code:    AUTH_OIDC_PROVIDER_BY_NAME_NOT_FOUND,
				Message: dbError.Message,
			}
		}
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	// Check restrictions
	oidcProvidersFiltered, err := api.GetAuthorizedOidcProviders(requestInfo, oidcProvider.Urn, AUTH_OIDC_ACTION_GET_PROVIDER, []OidcProvider{*oidcProvider})
	if err != nil {
		return nil, err
	}

	if len(oidcProvidersFiltered) > 0 {
		oidcProviderFiltered := oidcProvidersFiltered[0]
		return &oidcProviderFiltered, nil
	}
	return nil, &Error{
		Code: UNAUTHORIZED_RESOURCES_ERROR,
		Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
			requestInfo.Identifier, oidcProvider.Urn),
	}
}

func (api WorkerAPI) ListOidcProviders(requestInfo RequestInfo, filter *Filter) ([]string, int, error) {
	// Validate fields
	var total int
	orderByValidColumns := api.AuthOidcRepo.OrderByValidColumns(AUTH_OIDC_ACTION_LIST_PROVIDERS)
	err := validateFilter(filter, orderByValidColumns)
	if err != nil {
		return nil, total, err
	}

	// Call repo to retrieve the OIDC Providers
	oidcProviders, total, err := api.AuthOidcRepo.GetOidcProvidersFiltered(filter)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, total, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	// Check restrictions to list
	var urnPrefix string
	if len(filter.Org) == 0 {
		urnPrefix = "*"
	} else {
		urnPrefix = GetUrnPrefix("", RESOURCE_AUTH_OIDC_PROVIDER, filter.PathPrefix)
	}
	oidcProvidersFiltered, err := api.GetAuthorizedOidcProviders(requestInfo, urnPrefix, AUTH_OIDC_ACTION_LIST_PROVIDERS, oidcProviders)
	if err != nil {
		return nil, total, err
	}

	oidcProviderNames := []string{}
	for _, op := range oidcProvidersFiltered {
		oidcProviderNames = append(oidcProviderNames, op.Name)
	}

	return oidcProviderNames, total, nil
}

func (api WorkerAPI) UpdateOidcProvider(requestInfo RequestInfo, oidcProviderName string, newName string, newPath string, newIssuerUrl string,
	newClients []string) (*OidcProvider, error) {
	// Validate fields
	if !IsValidName(newName) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: new name %v", newName),
		}
	}
	if !IsValidPath(newPath) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: new path %v", newPath),
		}

	}
	if _, err := url.ParseRequestURI(newIssuerUrl); err != nil {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: issuerUrl %v", newIssuerUrl),
		}
	}
	err := AreValidOidcClientNames(newClients)
	if err != nil {
		apiError := err.(*Error)
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: apiError.Message,
		}

	}

	// Call repo to retrieve the old OIDC Provider
	oldOidcProvider, err := api.GetOidcProviderByName(requestInfo, oidcProviderName)
	if err != nil {
		return nil, err
	}

	// Check restrictions
	oidcProvidersFiltered, err := api.GetAuthorizedOidcProviders(requestInfo, oldOidcProvider.Urn,
		AUTH_OIDC_ACTION_UPDATE_PROVIDER, []OidcProvider{*oldOidcProvider})
	if err != nil {
		return nil, err
	}
	if len(oidcProvidersFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, oldOidcProvider.Urn),
		}
	}

	// Check if OIDC Provider with "newName" exists
	targetOidcProvider, err := api.GetOidcProviderByName(requestInfo, newName)

	if err == nil && targetOidcProvider.ID != oldOidcProvider.ID {
		// OIDC Provider already exists
		return nil, &Error{
			Code:    AUTH_OIDC_PROVIDER_ALREADY_EXIST,
			Message: fmt.Sprintf("OIDC Provider name: %v already exists", newName),
		}
	}

	if err != nil {
		if apiError := err.(*Error); apiError.Code != AUTH_OIDC_PROVIDER_BY_NAME_NOT_FOUND {
			return nil, err
		}
	}

	auxOidcProvider := OidcProvider{
		Urn: CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, newPath, newName),
	}

	// Check restrictions
	oidcProvidersFiltered, err = api.GetAuthorizedOidcProviders(requestInfo, auxOidcProvider.Urn,
		AUTH_OIDC_ACTION_UPDATE_PROVIDER, []OidcProvider{auxOidcProvider})
	if err != nil {
		return nil, err
	}
	if len(oidcProvidersFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, auxOidcProvider.Urn),
		}
	}

	// Create OIDC Clients
	oidcClients := []OidcClient{}
	for _, oc := range newClients {
		oidcClients = append(oidcClients, OidcClient{Name: oc})
	}

	oidcProvider := OidcProvider{
		ID:          oldOidcProvider.ID,
		Name:        newName,
		Path:        newPath,
		Urn:         auxOidcProvider.Urn,
		CreateAt:    oldOidcProvider.CreateAt,
		UpdateAt:    time.Now().UTC(),
		IssuerURL:   newIssuerUrl,
		OidcClients: oidcClients,
	}

	// Update OIDC Provider
	updatedOidcProvider, err := api.AuthOidcRepo.UpdateOidcProvider(oidcProvider)

	// Check unexpected DB error
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	LogOperation(requestInfo.RequestID, requestInfo.Identifier, fmt.Sprintf("OIDC Provider updated from %+v to %+v",
		oldOidcProvider, updatedOidcProvider))
	return updatedOidcProvider, nil
}

func (api WorkerAPI) RemoveOidcProvider(requestInfo RequestInfo, name string) error {
	// Call repo to retrieve the OIDC provider
	oidcProvider, err := api.GetOidcProviderByName(requestInfo, name)
	if err != nil {
		return err
	}

	// Check restrictions
	oidcProvidersFiltered, err := api.GetAuthorizedOidcProviders(requestInfo, oidcProvider.Urn, AUTH_OIDC_ACTION_DELETE_PROVIDER, []OidcProvider{*oidcProvider})
	if err != nil {
		return err
	}
	if len(oidcProvidersFiltered) < 1 {
		return &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, oidcProvider.Urn),
		}
	}

	err = api.AuthOidcRepo.RemoveOidcProvider(oidcProvider.ID)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	LogOperation(requestInfo.RequestID, requestInfo.Identifier, fmt.Sprintf("OIDC Provider deleted %v", oidcProvider))
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
