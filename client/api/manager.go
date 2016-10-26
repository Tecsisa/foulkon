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

type GroupAPI interface {
	GetGroup(organizationId, groupName string) (string, error)
	GetAllGroups(pathPrefix, offset, limit, orderBy string) (string, error)
	GetGroupsByOrg(organizationId, pathPrefix, offset, limit, orderBy string) (string, error)
	CreateGroup(organizationId, groupName, path string) (string, error)
	UpdateGroup(organizationId, groupName, newName, newPath string) (string, error)
	DeleteGroup(organizationId, groupName string) (string, error)
	GetGroupPolicies(organizationId, groupName, offset, limit, orderBy string) (string, error)
	AttachPolicyToGroup(organizationId, groupName, policyName string) (string, error)
	DetachPolicyFromGroup(organizationId, groupName, policyName string) (string, error)
	GetGroupMembers(organizationId, groupName, pathPrefix, offset, limit, orderBy string) (string, error)
	AddMemberToGroup(organizationId, groupName, userName string) (string, error)
	RemoveMemberFromGroup(organizationId, groupName, userName string) (string, error)
}

type ClientAPI struct {
	Address     string
	requestInfo map[string]string
	http.Client
}
