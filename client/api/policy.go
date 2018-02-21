package api

import (
	"encoding/json"

	"github.com/Tecsisa/foulkon/api"
	internalhttp "github.com/Tecsisa/foulkon/http"
)

func (c *ClientAPI) GetPolicy(organizationId, policyName string) (string, error) {
	req, err := c.prepareRequest("GET", internalhttp.API_VERSION_1+"/organizations/"+organizationId+"/policies/"+policyName, nil, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) GetAllPolicies(pathPrefix, offset, limit, orderBy string) (string, error) {
	urlParams := map[string]string{
		"PathPrefix": pathPrefix,
		"Offset":     offset,
		"Limit":      limit,
		"OrderBy":    orderBy,
	}
	req, err := c.prepareRequest("GET", internalhttp.API_VERSION_1+"/policies", nil, urlParams)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) CreatePolicy(organizationId, policyName, path, statement string) (string, error) {

	statementApi := []api.Statement{}
	if err := json.Unmarshal([]byte(statement), &statementApi); err != nil {
		return "", err
	}
	body := map[string]interface{}{
		"name":       policyName,
		"path":       path,
		"Statements": statementApi,
	}

	req, err := c.prepareRequest("POST", internalhttp.API_VERSION_1+"/organizations/"+organizationId+"/policies", body, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) UpdatePolicy(organizationId, policyName, path, statement string) (string, error) {
	statementApi := []api.Statement{}
	if err := json.Unmarshal([]byte(statement), &statementApi); err != nil {
		return "", err
	}
	body := map[string]interface{}{
		"name":       policyName,
		"path":       path,
		"Statements": statementApi,
	}

	req, err := c.prepareRequest("PUT", internalhttp.API_VERSION_1+"/organizations/"+organizationId+"/policies/"+policyName, body, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) DeletePolicy(organizationId, policyName string) (string, error) {
	req, err := c.prepareRequest("DELETE", internalhttp.API_VERSION_1+"/organizations/"+organizationId+"/policies/"+policyName, nil, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) GetGroupsAttached(organizationId, policyName, offset, limit, orderBy string) (string, error) {
	urlParams := map[string]string{
		"Offset":  offset,
		"Limit":   limit,
		"OrderBy": orderBy,
	}
	req, err := c.prepareRequest("GET", internalhttp.API_VERSION_1+"/organizations/"+organizationId+"/policies/"+policyName+"/groups", nil, urlParams)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) GetPoliciesOrganization(organizationId, pathPrefix, offset, limit, orderBy string) (string, error) {
	urlParams := map[string]string{
		"PathPrefix": pathPrefix,
		"Offset":     offset,
		"Limit":      limit,
		"OrderBy":    orderBy,
	}
	req, err := c.prepareRequest("GET", internalhttp.API_VERSION_1+"/organizations/"+organizationId+"/policies", nil, urlParams)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}
