package api

import "net/http"

type UserAPI interface {
	GetUser(externalId string) (string, error)

	GetAllUsers(pathPrefix, offset, limit, orderBy string) (string, error)

	GetUserGroups(externalId, offset, limit, orderBy string) (string, error)

	CreateUser(externalId, path string) (string, error)

	UpdateUser(externalId, path string) (string, error)

	DeleteUser(externalId string) (string, error)
}

type ClientAPI struct {
	Address     string
	requestInfo map[string]string
	http.Client
}
