package api

import "github.com/Tecsisa/foulkon/database"

// TYPE DEFINITIONS

// ProxyResource domain
type ProxyResource struct {
	ID     string `json:"id, omitempty"`
	Host   string `json:"host, omitempty"`
	Url    string `json:"url, omitempty"`
	Method string `json:"method, omitempty"`
	Urn    string `json:"urn, omitempty"`
	Action string `json:"action, omitempty"`
}

// GetProxyResources return proxy resources
func (api ProxyAPI) GetProxyResources() ([]ProxyResource, error) {
	resources, err := api.ProxyRepo.GetProxyResources()

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
