package api

import (
	"encoding/json"
	"strings"

	internalhttp "github.com/Tecsisa/foulkon/http"
)

func (c *ClientAPI) GetAuthorizedResources(action, resources string) (string, error) {
	numResources := strings.Count("resources", ",") + 1
	resourcesJson := make([]string, numResources)
	if err := json.Unmarshal([]byte(resources), &resourcesJson); err != nil {
		return "", err
	}
	body := map[string]interface{}{
		"action":    action,
		"resources": resourcesJson,
	}

	req, err := c.prepareRequest("POST", internalhttp.RESOURCE_URL, body, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}
