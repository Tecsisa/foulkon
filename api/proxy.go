package api

import (
	"fmt"
	"time"

	"github.com/Tecsisa/foulkon/database"
	"github.com/satori/go.uuid"
)

// TYPE DEFINITIONS

// ProxyResource domain
type ProxyResource struct {
	ID       string         `json:"id, omitempty"`
	Name     string         `json:"name, omitempty"`
	Org      string         `json:"org, omitempty"`
	Path     string         `json:"path, omitempty"`
	Urn      string         `json:"urn, omitempty"`
	Resource ResourceEntity `json:"resource, omitempty"`
	CreateAt time.Time      `json:"createAt, omitempty"`
	UpdateAt time.Time      `json:"updateAt, omitempty"`
}

// Proxy resource identifier to retrieve them from DB
type ProxyResourceIdentity struct {
	Org  string `json:"org, omitempty"`
	Name string `json:"name, omitempty"`
}

type ResourceEntity struct {
	Host   string `json:"host, omitempty"`
	Path   string `json:"path, omitempty"`
	Method string `json:"method, omitempty"`
	Urn    string `json:"urn, omitempty"`
	Action string `json:"action, omitempty"`
}

func (p ProxyResource) GetUrn() string {
	return p.Urn
}

// GetProxyResources return proxy resources
func (api ProxyAPI) GetProxyResources() ([]ProxyResource, error) {
	resources, _, err := api.ProxyRepo.GetProxyResources(&Filter{})

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	return resources, nil
}

func (api WorkerAPI) AddProxyResource(requestInfo RequestInfo, name string, org string, path string, resource ResourceEntity) (*ProxyResource, error) {
	// Validate fields
	if !IsValidName(name) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: name %v", name),
		}
	}
	if !IsValidOrg(org) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: org %v", org),
		}
	}
	if !IsValidPath(path) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: path %v", path),
		}
	}
	err := IsValidProxyResource(&resource)
	if err != nil {
		return nil, err
	}

	proxyResource := createProxyResource(name, org, path, resource)

	// Check restrictions
	proxyResourcesFiltered, err := api.GetAuthorizedProxyResources(requestInfo, proxyResource.Urn, PROXY_ACTION_CREATE_RESOURCE, []ProxyResource{proxyResource})

	if err != nil {
		return nil, err
	}
	if len(proxyResourcesFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, proxyResource.Urn),
		}
	}

	// Check if proxy resource already exists
	_, err = api.ProxyRepo.GetProxyResourceByName(org, name)

	if err != nil {
		// Transform to DB error
		dbError := err.(*database.Error)
		// Proxy resource doesn't exist in DB
		switch dbError.Code {
		case database.PROXY_RESOURCE_NOT_FOUND:
			// Create proxy resource
			created, err := api.ProxyRepo.AddProxyResource(proxyResource)

			// Check unexpected DB error
			if err != nil {
				//Transform to DB error
				dbError := err.(*database.Error)
				return nil, &Error{
					Code:    UNKNOWN_API_ERROR,
					Message: dbError.Message,
				}
			}
			LogOperation(requestInfo.RequestID, requestInfo.Identifier, fmt.Sprintf("proxy resource created %+v", created))
			return created, nil
		default: // Unexpected error
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	} else {
		return nil, &Error{
			Code:    PROXY_RESOURCE_ALREADY_EXIST,
			Message: fmt.Sprintf("Unable to create proxy resource, proxy resource with org %v and name %v already exist", org, name),
		}
	}
}

func (api WorkerAPI) GetProxyResourceByName(requestInfo RequestInfo, org string, name string) (*ProxyResource, error) {
	// Validate fields
	if !IsValidName(name) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: name %v", name),
		}
	}
	if !IsValidOrg(org) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: org %v", org),
		}
	}

	// Call repo to retrieve the proxy resource
	proxyResource, err := api.ProxyRepo.GetProxyResourceByName(org, name)

	// Error handling
	if err != nil {
		// Transform to DB error
		dbError := err.(*database.Error)
		// Proxy resource doesn't exist in DB
		switch dbError.Code {
		case database.PROXY_RESOURCE_NOT_FOUND:
			return nil, &Error{
				Code:    PROXY_RESOURCE_BY_ORG_AND_NAME_NOT_FOUND,
				Message: dbError.Message,
			}
		default:
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	}

	// Check restrictions
	proxyResourceFiltered, err := api.GetAuthorizedProxyResources(requestInfo, proxyResource.Urn, PROXY_ACTION_GET_PROXY_RESOURCE, []ProxyResource{*proxyResource})
	if err != nil {
		return nil, err
	}

	// Check if we have our user authorized
	if len(proxyResourceFiltered) > 0 {
		proxyResourceFiltered := proxyResourceFiltered[0]
		return &proxyResourceFiltered, nil
	}
	return nil, &Error{
		Code: UNAUTHORIZED_RESOURCES_ERROR,
		Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
			requestInfo.Identifier, proxyResource.Urn),
	}
}

func (api WorkerAPI) UpdateProxyResource(requestInfo RequestInfo, org string, name string, newName string, newPath string, newResource ResourceEntity) (*ProxyResource, error) {
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
	err := IsValidProxyResource(&newResource)
	if err != nil {
		return nil, err
	}

	// Call repo to retrieve the old proxy resource
	oldProxyResource, err := api.GetProxyResourceByName(requestInfo, org, name)
	if err != nil {
		return nil, err
	}

	// Check restrictions
	proxyResourcesFiltered, err := api.GetAuthorizedProxyResources(requestInfo, oldProxyResource.Urn, PROXY_ACTION_UPDATE_RESOURCE, []ProxyResource{*oldProxyResource})
	if err != nil {
		return nil, err
	}
	if len(proxyResourcesFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, oldProxyResource.Urn),
		}
	}

	// Check if a proxy resource with "newName" already exists
	newProxyResource, err := api.GetProxyResourceByName(requestInfo, org, newName)

	if err == nil && oldProxyResource.ID != newProxyResource.ID {
		// Proxy resource already exists
		return nil, &Error{
			Code:    PROXY_RESOURCE_ALREADY_EXIST,
			Message: fmt.Sprintf("Proxy resource name: %v already exists", newName),
		}
	}

	if err != nil {
		if apiError := err.(*Error); apiError.Code != PROXY_RESOURCE_BY_ORG_AND_NAME_NOT_FOUND {
			return nil, err
		}
	}

	auxProxyResource := ProxyResource{
		Urn: CreateUrn(org, RESOURCE_PROXY, newPath, newName),
	}

	// Check restrictions
	proxyResourcesFiltered, err = api.GetAuthorizedProxyResources(requestInfo, auxProxyResource.Urn, PROXY_ACTION_UPDATE_RESOURCE, []ProxyResource{auxProxyResource})
	if err != nil {
		return nil, err
	}
	if len(proxyResourcesFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, auxProxyResource.Urn),
		}
	}

	// Update proxy resource
	proxyResource := ProxyResource{
		ID:       oldProxyResource.ID,
		Name:     newName,
		Path:     newPath,
		Org:      oldProxyResource.Org,
		Urn:      auxProxyResource.Urn,
		Resource: newResource,
		CreateAt: oldProxyResource.CreateAt,
		UpdateAt: time.Now().UTC(),
	}

	updatedProxyResource, err := api.ProxyRepo.UpdateProxyResource(proxyResource)

	// Check unexpected DB error
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	LogOperation(requestInfo.RequestID, requestInfo.Identifier, fmt.Sprintf("Proxy resource updated from %+v to %+v", oldProxyResource, updatedProxyResource))
	return updatedProxyResource, nil
}

func (api WorkerAPI) RemoveProxyResource(requestInfo RequestInfo, org string, name string) error {
	// Call repo to retrieve the proxy resource
	proxyResource, err := api.GetProxyResourceByName(requestInfo, org, name)
	if err != nil {
		return err
	}

	// Check restrictions
	proxyResourcesFiltered, err := api.GetAuthorizedProxyResources(requestInfo, proxyResource.Urn, PROXY_ACTION_DELETE_RESOURCE, []ProxyResource{*proxyResource})
	if err != nil {
		return err
	}
	if len(proxyResourcesFiltered) < 1 {
		return &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, proxyResource.Urn),
		}
	}

	err = api.ProxyRepo.RemoveProxyResource(proxyResource.ID)

	// Error handling
	if err != nil {
		// Transform to DB error
		dbError := err.(*database.Error)
		return &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	LogOperation(requestInfo.RequestID, requestInfo.Identifier, fmt.Sprintf("Proxy resource deleted %+v", proxyResource))
	return nil
}

func (api WorkerAPI) ListProxyResources(requestInfo RequestInfo, filter *Filter) ([]ProxyResourceIdentity, int, error) {
	// Validate fields
	var total int
	orderByValidColumns := api.ProxyRepo.OrderByValidColumns(PROXY_ACTION_LIST_RESOURCES)
	err := validateFilter(filter, orderByValidColumns)
	if err != nil {
		return nil, total, err
	}

	// Call repo to retrieve the proxy resources
	proxyResources, total, err := api.ProxyRepo.GetProxyResources(filter)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, total, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	// Check restrictions
	var urnPrefix string
	if len(filter.Org) == 0 {
		urnPrefix = "*"
	} else {
		urnPrefix = GetUrnPrefix(filter.Org, RESOURCE_PROXY, filter.PathPrefix)
	}
	proxyResourcesFiltered, err := api.GetAuthorizedProxyResources(requestInfo, urnPrefix, PROXY_ACTION_LIST_RESOURCES, proxyResources)
	if err != nil {
		return nil, total, err
	}

	proxyResourcesIDs := []ProxyResourceIdentity{}
	for _, p := range proxyResourcesFiltered {
		proxyResourcesIDs = append(proxyResourcesIDs, ProxyResourceIdentity{
			Org:  p.Org,
			Name: p.Name,
		})
	}

	return proxyResourcesIDs, total, nil
}

// PRIVATE HELPER METHODS

func createProxyResource(name string, org string, path string, resource ResourceEntity) ProxyResource {
	pr := ProxyResource{
		ID:       uuid.NewV4().String(),
		Name:     name,
		Org:      org,
		Path:     path,
		Urn:      CreateUrn(org, RESOURCE_PROXY, path, name),
		Resource: resource,
		CreateAt: time.Now().UTC(),
		UpdateAt: time.Now().UTC(),
	}

	return pr
}
