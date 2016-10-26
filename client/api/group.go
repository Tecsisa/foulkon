package api

import internalhttp "github.com/Tecsisa/foulkon/http"

func (c *ClientAPI) GetGroup(organizationId, groupName string) (string, error) {
	req, err := c.prepareRequest("GET", internalhttp.API_VERSION_1+"/organizations/"+organizationId+"/groups/"+groupName, nil, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) GetAllGroups(pathPrefix, offset, limit, orderBy string) (string, error) {
	urlParams := map[string]string{
		"PathPrefix": pathPrefix,
		"Offset":     offset,
		"Limit":      limit,
		"OrderBy":    orderBy,
	}
	req, err := c.prepareRequest("GET", internalhttp.API_VERSION_1+"/groups", nil, urlParams)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) GetGroupsByOrg(organizationId, pathPrefix, offset, limit, orderBy string) (string, error) {
	urlParams := map[string]string{
		"PathPrefix": pathPrefix,
		"Offset":     offset,
		"Limit":      limit,
		"OrderBy":    orderBy,
	}
	req, err := c.prepareRequest("GET", internalhttp.API_VERSION_1+"/organizations/"+organizationId+"/groups", nil, urlParams)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) CreateGroup(organizationId, groupName, path string) (string, error) {
	body := map[string]interface{}{
		"name": groupName,
		"path": path,
	}
	req, err := c.prepareRequest("POST", internalhttp.API_VERSION_1+"/organizations/"+organizationId+"/groups", body, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) UpdateGroup(organizationId, groupName, newName, newPath string) (string, error) {
	body := map[string]interface{}{
		"name": newName,
		"path": newPath,
	}
	req, err := c.prepareRequest("PUT", internalhttp.API_VERSION_1+"/organizations/"+organizationId+"/groups/"+groupName, body, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) DeleteGroup(organizationId, groupName string) (string, error) {
	req, err := c.prepareRequest("DELETE", internalhttp.API_VERSION_1+"/organizations/"+organizationId+"/groups/"+groupName, nil, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) GetGroupPolicies(organizationId, groupName, offset, limit, orderBy string) (string, error) {
	urlParams := map[string]string{
		"Offset":  offset,
		"Limit":   limit,
		"OrderBy": orderBy,
	}
	req, err := c.prepareRequest("GET", internalhttp.API_VERSION_1+"/organizations/"+organizationId+"/groups/"+groupName+"/policies", nil, urlParams)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) AttachPolicyToGroup(organizationId, groupName, policyName string) (string, error) {
	req, err := c.prepareRequest("POST", internalhttp.API_VERSION_1+"/organizations/"+organizationId+"/groups/"+groupName+"/policies/"+policyName, nil, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) DetachPolicyFromGroup(organizationId, groupName, policyName string) (string, error) {
	req, err := c.prepareRequest("DELETE", internalhttp.API_VERSION_1+"/organizations/"+organizationId+"/groups/"+groupName+"/policies/"+policyName, nil, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) GetGroupMembers(organizationId, groupName, pathPrefix, offset, limit, orderBy string) (string, error) {
	urlParams := map[string]string{
		"PathPrefix": pathPrefix,
		"Offset":     offset,
		"Limit":      limit,
		"OrderBy":    orderBy,
	}
	req, err := c.prepareRequest("GET", internalhttp.API_VERSION_1+"/organizations/"+organizationId+"/groups/"+groupName+"/users", nil, urlParams)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) AddMemberToGroup(organizationId, groupName, userName string) (string, error) {
	req, err := c.prepareRequest("POST", internalhttp.API_VERSION_1+"/organizations/"+organizationId+"/groups/"+groupName+"/users/"+userName, nil, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) RemoveMemberFromGroup(organizationId, groupName, userName string) (string, error) {
	req, err := c.prepareRequest("DELETE", internalhttp.API_VERSION_1+"/organizations/"+organizationId+"/groups/"+groupName+"/users/"+userName, nil, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}
