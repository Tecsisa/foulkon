package client

import (
	"flag"

	internalhttp "github.com/Tecsisa/foulkon/http"
)

type GetUserCommand struct {
	Meta
}

func (c *GetUserCommand) Run(args []string) (string, error) {
	flagSet := flag.NewFlagSet("user get", flag.ExitOnError)
	id := flagSet.String("id", "", "External user ID")
	flagSet.Parse(args)

	req, err := c.prepareRequest("GET", internalhttp.USER_ROOT_URL+"/"+*id, nil, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

type GetAllUsersCommand struct {
	Meta
}

func (c *GetAllUsersCommand) Run(args []string) (string, error) {
	flagSet := flag.NewFlagSet("user get-all", flag.ExitOnError)
	offset := flagSet.String("offset", "", "The offset of the items returned")
	limit := flagSet.String("limit", "", "The maximum number of items in the response")
	orderBy := flagSet.String("order-by", "", "order data by field")
	pathPrefix := flagSet.String("path-prefix", "", "search starts from this path")
	flagSet.Parse(args)

	queryParams := map[string]string{
		"Offset":     *offset,
		"Limit":      *limit,
		"OrderBy":    *orderBy,
		"PathPrefix": *pathPrefix,
	}

	req, err := c.prepareRequest("GET", internalhttp.USER_ROOT_URL, nil, queryParams)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

type GetUserGroupsCommand struct {
	Meta
}

func (c *GetUserGroupsCommand) Run(args []string) (string, error) {
	flagSet := flag.NewFlagSet("user groups", flag.ExitOnError)
	id := flagSet.String("id", "", "External user ID")
	offset := flagSet.String("offset", "", "The offset of the items returned")
	limit := flagSet.String("limit", "", "The maximum number of items in the response")
	orderBy := flagSet.String("order-by", "", "order data by field")
	flagSet.Parse(args)

	queryParams := map[string]string{
		"Offset":  *offset,
		"Limit":   *limit,
		"OrderBy": *orderBy,
	}
	flagSet.Parse(args)

	req, err := c.prepareRequest("GET", internalhttp.USER_ROOT_URL+"/"+*id+"/groups", nil, queryParams)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

type CreateUserCommand struct {
	Meta
}

func (c *CreateUserCommand) Run(args []string) (string, error) {
	flagSet := flag.NewFlagSet("user-create", flag.ExitOnError)
	externalId := flagSet.String("id", "", "User's external identifier")
	path := flagSet.String("path", "", "User location")
	flagSet.Parse(args)

	body := map[string]string{
		"externalId": *externalId,
		"path":       *path,
	}

	req, err := c.prepareRequest("POST", internalhttp.USER_ROOT_URL, body, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

type UpdateUserCommand struct {
	Meta
}

func (c *UpdateUserCommand) Run(args []string) (string, error) {
	flagSet := flag.NewFlagSet("user update", flag.ExitOnError)
	id := flagSet.String("id", "", "Existing user Id")
	path := flagSet.String("path", "", "User location")
	flagSet.Parse(args)

	body := map[string]string{
		"path": *path,
	}

	req, err := c.prepareRequest("PUT", internalhttp.USER_ROOT_URL+"/"+*id, body, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}

type DeleteUserCommand struct {
	Meta
}

func (c *DeleteUserCommand) Run(args []string) (string, error) {
	flagSet := flag.NewFlagSet("user get", flag.ExitOnError)
	id := flagSet.String("id", "", "External user ID")
	flagSet.Parse(args)

	req, err := c.prepareRequest("DELETE", internalhttp.USER_ROOT_URL+"/"+*id, nil, nil)
	if err != nil {
		return "", err
	}
	return c.makeRequest(req)
}
