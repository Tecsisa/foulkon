package postgresql

import (
	"fmt"

	"time"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/database"
)

// PROXY REPOSITORY IMPLEMENTATION

func (pr PostgresRepo) GetProxyResourceByName(org string, name string) (*api.ProxyResource, error) {
	proxyResource := &ProxyResource{}
	query := pr.Dbmap.Where("org like ? AND name like ?", org, name).First(proxyResource)

	// Check if proxyResource exists
	if query.RecordNotFound() {
		return nil, &database.Error{
			Code:    database.PROXY_RESOURCE_NOT_FOUND,
			Message: fmt.Sprintf("Proxy resource with organization %v and name %v not found", org, name),
		}
	}
	// Error Handling
	if err := query.Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	return dbResourceToApiResource(proxyResource), nil
}

func (pr PostgresRepo) GetProxyResources(filter *api.Filter) ([]api.ProxyResource, int, error) {
	var total int
	resources := []ProxyResource{}
	query := pr.Dbmap

	if len(filter.Org) > 0 {
		query = query.Where("org like ? ", filter.Org)
	}
	if len(filter.PathPrefix) > 0 {
		query = query.Where("path like ? ", filter.PathPrefix+"%")
	}
	if len(filter.OrderBy) > 0 {
		query = query.Order(filter.OrderBy)
	}

	// Error handling
	if err := query.Find(&resources).Count(&total).Offset(filter.Offset).Limit(filter.Limit).Find(&resources).Error; err != nil {
		return nil, total, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Transform proxyResources to API domain
	var proxyResources []api.ProxyResource
	if resources != nil {
		proxyResources = make([]api.ProxyResource, len(resources), cap(resources))
		for i, pr := range resources {
			proxyResources[i] = *dbResourceToApiResource(&pr)
		}
	}

	return proxyResources, total, nil
}

func (pr PostgresRepo) AddProxyResource(proxyResource api.ProxyResource) (*api.ProxyResource, error) {
	// Create proxyResource model
	proxyResourceDB := &ProxyResource{
		ID:           proxyResource.ID,
		Name:         proxyResource.Name,
		Org:          proxyResource.Org,
		Path:         proxyResource.Path,
		Host:         proxyResource.Resource.Host,
		PathResource: proxyResource.Resource.Path,
		Method:       proxyResource.Resource.Method,
		UrnResource:  proxyResource.Resource.Urn,
		Action:       proxyResource.Resource.Action,
		Urn:          proxyResource.Urn,
		CreateAt:     proxyResource.CreateAt.UnixNano(),
		UpdateAt:     proxyResource.UpdateAt.UnixNano(),
	}

	// Store proxyResource
	err := pr.Dbmap.Create(proxyResourceDB).Error

	// Error handling
	if err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	return dbResourceToApiResource(proxyResourceDB), nil
}

func (pr PostgresRepo) UpdateProxyResource(proxyResource api.ProxyResource) (*api.ProxyResource, error) {
	proxyResourceDB := &ProxyResource{
		ID:           proxyResource.ID,
		Name:         proxyResource.Name,
		Org:          proxyResource.Org,
		Path:         proxyResource.Path,
		Host:         proxyResource.Resource.Host,
		PathResource: proxyResource.Resource.Path,
		Method:       proxyResource.Resource.Method,
		UrnResource:  proxyResource.Resource.Urn,
		Action:       proxyResource.Resource.Action,
		Urn:          proxyResource.Urn,
		CreateAt:     proxyResource.CreateAt.UnixNano(),
		UpdateAt:     proxyResource.UpdateAt.UnixNano(),
	}

	// Store proxyResource
	query := pr.Dbmap.Model(&ProxyResource{ID: proxyResource.ID}).Updates(proxyResourceDB)

	// Error Handling
	if err := query.Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	return &proxyResource, nil
}

func (pr PostgresRepo) RemoveProxyResource(id string) error {
	// Remove proxy resource
	query := pr.Dbmap.Where("id like ?", id).Delete(&ProxyResource{})

	// Error handling
	if err := query.Error; err != nil {
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	return nil
}

// PRIVATE HELPER METHODS

// Transform a proxyResource retrieved from db into a proxyResource for API
func dbResourceToApiResource(pr *ProxyResource) *api.ProxyResource {
	return &api.ProxyResource{
		ID:   pr.ID,
		Name: pr.Name,
		Path: pr.Path,
		Org:  pr.Org,
		Resource: api.ResourceEntity{
			Host:   pr.Host,
			Path:   pr.PathResource,
			Method: pr.Method,
			Urn:    pr.UrnResource,
			Action: pr.Action,
		},
		Urn:      pr.Urn,
		CreateAt: time.Unix(0, pr.CreateAt).UTC(),
		UpdateAt: time.Unix(0, pr.UpdateAt).UTC(),
	}
}
