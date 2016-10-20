package api

import internalhttp "github.com/Tecsisa/foulkon/http"

func (c *ClientAPI) GetUser(externalId string) (string, error) {
	req, err := c.prepareRequest("GET", internalhttp.USER_ROOT_URL+"/"+externalId, nil, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) GetAllUsers(pathPrefix, offset, limit, orderBy string) (string, error) {
	urlParams := map[string]string{
		"pathPrefix": pathPrefix,
		"offset":     offset,
		"limit":      limit,
		"orderBy":    orderBy,
	}
	req, err := c.prepareRequest("GET", internalhttp.USER_ROOT_URL, nil, urlParams)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) GetUserGroups(externalId, offset, limit, orderBy string) (string, error) {
	urlParams := map[string]string{
		"offset":  offset,
		"limit":   limit,
		"orderBy": orderBy,
	}
	req, err := c.prepareRequest("GET", internalhttp.USER_ROOT_URL+"/"+externalId+"/groups", nil, urlParams)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) CreateUser(externalId, path string) (string, error) {
	body := map[string]string{
		"externalId": externalId,
		"path":       path,
	}

	req, err := c.prepareRequest("POST", internalhttp.USER_ROOT_URL, body, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) UpdateUser(externalId, path string) (string, error) {
	body := map[string]string{
		"path": path,
	}

	req, err := c.prepareRequest("PUT", internalhttp.USER_ROOT_URL+"/"+externalId, body, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

func (c *ClientAPI) DeleteUser(externalId string) (string, error) {
	req, err := c.prepareRequest("DELETE", internalhttp.USER_ROOT_URL+"/"+externalId, nil, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}
