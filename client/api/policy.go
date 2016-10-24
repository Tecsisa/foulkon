package api

import internalhttp "github.com/Tecsisa/foulkon/http"

func (c *ClientAPI) GetPolicy(organizationId, policyName string) (string, error) {
	req, err := c.prepareRequest("GET", internalhttp.API_VERSION_1+"/organizations/"+organizationId+"/policies/"+policyName, nil, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) GetAllPolicy(pathPrefix, offset, limit, orderBy string) (string, error) {
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

func (c *ClientAPI) CreatePolicy(organizationId, policyName, path, effect, actions, resources string) (string, error) {
	body := map[string]string{
		"name":       policyName,
		"path":       path,
		"Statements": "\"Statements\" : [       {       \"Effect\" : \"allow\",       \"Actions\" : [\"iam:*\"],       \"Resources\" : [\"urn:everything:*\"]       }   ]",
		//"effect":    effect,
		//"actions":   actions,
		//"resources": resources,
	}
	req, err := c.prepareRequest("POST", internalhttp.API_VERSION_1+"/organizations/"+organizationId+"/policies", body, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

//func (c *ClientAPI) UpdatePolicy(organizationId, policyName string) (string, error) {
//	req, err := c.prepareRequest("PUT", internalhttp.API_VERSION_1+"/organizations/"+organizationId+"/policies/"+policyName, nil, nil)
//	if err != nil {
//		return "", err
//	}
//	return c.makeRequest(req)
//}

func (c *ClientAPI) DeletePolicy(organizationId, policyName string) (string, error) {
	req, err := c.prepareRequest("DELETE", internalhttp.API_VERSION_1+"/organizations/"+organizationId+"/policies/"+policyName, nil, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) GetGroupsPolicy(organizationId, policyName, offset, limit, orderBy string) (string, error) {
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
