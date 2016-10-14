package main

import (
	"flag"

	internalhttp "github.com/Tecsisa/foulkon/http"
	"net/http"
)

type GetUserCommand struct {
	Meta
}

func (c *GetUserCommand) Run(args []string) int {
	flagSet := flag.NewFlagSet("user get", flag.ExitOnError)
	id := flagSet.String("id", "", "External user ID")
	flagSet.Parse(args)

	req := c.prepareRequest("GET", internalhttp.USER_ROOT_URL+"/"+*id, nil, nil)
	return c.makeRequest(req, http.StatusOK, true)
}

type GetAllUsersCommand struct {
	Meta
}

func (c *GetAllUsersCommand) Run(args []string) int {
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

	req := c.prepareRequest("GET", internalhttp.USER_ROOT_URL, nil, queryParams)
	return c.makeRequest(req, http.StatusOK, true)
}

type GetUserGroupsCommand struct {
	Meta
}

func (c *GetUserGroupsCommand) Run(args []string) int {
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

	req := c.prepareRequest("GET", internalhttp.USER_ROOT_URL+"/"+*id+"/groups", nil, queryParams)
	return c.makeRequest(req, http.StatusOK, true)
}

type CreateUserCommand struct {
	Meta
}

func (c *CreateUserCommand) Run(args []string) int {
	flagSet := flag.NewFlagSet("user-create", flag.ExitOnError)
	externalId := flagSet.String("id", "", "User's external identifier")
	path := flagSet.String("path", "", "User location")
	flagSet.Parse(args)

	body := map[string]string{
		"externalId": *externalId,
		"path":       *path,
	}

	req := c.prepareRequest("POST", internalhttp.USER_ROOT_URL, body, nil)
	return c.makeRequest(req, http.StatusCreated, true)
}

type DeleteUserCommand struct {
	Meta
}

type UpdateUserCommand struct {
	Meta
}

func (c *UpdateUserCommand) Run(args []string) int {
	flagSet := flag.NewFlagSet("user update", flag.ExitOnError)
	path := flagSet.String("path", "", "User location")
	id := flagSet.String("id", "", "Existing user Id")
	flagSet.Parse(args)

	body := map[string]string{
		"path": *path,
	}

	req := c.prepareRequest("PUT", internalhttp.USER_ROOT_URL+"/"+*id, body, nil)
	return c.makeRequest(req, http.StatusOK, true)
}

func (c *DeleteUserCommand) Run(args []string) int {
	flagSet := flag.NewFlagSet("user get", flag.ExitOnError)
	id := flagSet.String("id", "", "External user ID")
	flagSet.Parse(args)

	req := c.prepareRequest("DELETE", internalhttp.USER_ROOT_URL+"/"+*id, nil, nil)
	return c.makeRequest(req, http.StatusNoContent, false)
}
