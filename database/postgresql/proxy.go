package postgresql

import (
	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/database"
)

// PROXY REPOSITORY IMPLEMENTATION

func (pr PostgresRepo) GetProxyResources() ([]api.ProxyResource, error) {
	resources := []ProxyResource{}
	query := pr.Dbmap.Find(&resources)

	// Error Handling
	if err := query.Error; err != nil {
		return nil, &database.Error{
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
	return proxyResources, nil
}

// PRIVATE HELPER METHODS

// Transform a proxyResource retrieved from db into a proxyResource for API
func dbResourceToApiResource(pr *ProxyResource) *api.ProxyResource {
	return &api.ProxyResource{
		ID:     pr.ID,
		Host:   pr.Host,
		Url:    pr.Url,
		Method: pr.Method,
		Urn:    pr.Urn,
		Action: pr.Action,
	}
}
