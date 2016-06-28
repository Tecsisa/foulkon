package api

import log "github.com/Sirupsen/logrus"

// Authorizr API that implements API interfaces using repositories
type AuthAPI struct {
	UserRepo   UserRepo
	GroupRepo  GroupRepo
	PolicyRepo PolicyRepo
	Logger     log.Logger
}

// API INTERFACES

type UserAPI interface {
	AddUser(authenticatedUser AuthenticatedUser, externalID string, path string) (*User, error)
	GetUserByExternalId(authenticatedUser AuthenticatedUser, id string) (*User, error)
	GetUserList(authenticatedUser AuthenticatedUser, pathPrefix string) ([]string, error)
	UpdateUser(authenticatedUser AuthenticatedUser, externalID string, newPath string) (*User, error)
	RemoveUserById(authenticatedUser AuthenticatedUser, id string) error
	GetGroupsByUserId(authenticatedUser AuthenticatedUser, id string) ([]GroupIdentity, error)
}

type GroupAPI interface {
	AddGroup(authenticatedUser AuthenticatedUser, org string, name string, path string) (*Group, error)
	GetGroupByName(authenticatedUser AuthenticatedUser, org string, name string) (*Group, error)
	GetGroupList(authenticatedUser AuthenticatedUser, org string, pathPrefix string) ([]GroupIdentity, error)
	UpdateGroup(authenticatedUser AuthenticatedUser, org string, groupName string, newName string, newPath string) (*Group, error)
	RemoveGroup(authenticatedUser AuthenticatedUser, org string, name string) error

	AddMember(authenticatedUser AuthenticatedUser, userID string, groupName string, org string) error
	RemoveMember(authenticatedUser AuthenticatedUser, userID string, groupName string, org string) error
	ListMembers(authenticatedUser AuthenticatedUser, org string, groupName string) ([]string, error)

	AttachPolicyToGroup(authenticatedUser AuthenticatedUser, org string, groupName string, policyName string) error
	DetachPolicyToGroup(authenticatedUser AuthenticatedUser, org string, groupName string, policyName string) error
	ListAttachedGroupPolicies(authenticatedUser AuthenticatedUser, org string, groupName string) ([]PolicyIdentity, error)
}

type PolicyAPI interface {
	AddPolicy(authenticatedUser AuthenticatedUser, name string, path string, org string, statements []Statement) (*Policy, error)
	GetPolicyByName(authenticatedUser AuthenticatedUser, org string, policyName string) (*Policy, error)
	GetPolicyList(authenticatedUser AuthenticatedUser, org string, pathPrefix string) ([]PolicyIdentity, error)
	UpdatePolicy(authenticatedUser AuthenticatedUser, org string, policyName string, newName string, newPath string,
		newStatements []Statement) (*Policy, error)
	DeletePolicy(authenticatedUser AuthenticatedUser, org string, name string) error
	GetAttachedGroups(authenticatedUser AuthenticatedUser, org string, policyName string) ([]GroupIdentity, error)
}

type AuthzAPI interface {
	GetAuthorizedUsers(authenticatedUser AuthenticatedUser, resourceUrn string, action string, users []User) ([]User, error)
	GetAuthorizedGroups(authenticatedUser AuthenticatedUser, resourceUrn string, action string, groups []Group) ([]Group, error)
	GetAuthorizedPolicies(authenticatedUser AuthenticatedUser, resourceUrn string, action string, policies []Policy) ([]Policy, error)
	GetAuthorizedExternalResources(authenticatedUser AuthenticatedUser, action string, resources []string) ([]string, error)
}

// REPOSITORY INTERFACES

// User repository that contains all user operations for this domain
type UserRepo interface {
	GetUserByExternalID(id string) (*User, error)

	AddUser(user User) (*User, error)
	UpdateUser(user User, newPath string, newUrn string) (*User, error)

	GetUsersFiltered(pathPrefix string) ([]User, error)
	GetGroupsByUserID(id string) ([]Group, error)
	RemoveUser(id string) error
}

// Group repository that contains all user operations for this domain
type GroupRepo interface {
	GetGroupByName(org string, name string) (*Group, error)
	IsMemberOfGroup(userID string, groupID string) (bool, error)
	GetGroupMembers(groupID string) ([]User, error)
	IsAttachedToGroup(groupID string, policyID string) (bool, error)
	GetAttachedPolicies(groupID string) ([]Policy, error)
	GetGroupsFiltered(org string, pathPrefix string) ([]Group, error)
	RemoveGroup(id string) error

	AddGroup(group Group) (*Group, error)
	AddMember(userID string, groupID string) error
	RemoveMember(userID string, groupID string) error
	UpdateGroup(group Group, newName string, newPath string, newUrn string) (*Group, error)
	AttachPolicy(groupID string, policyID string) error
	DetachPolicy(groupID string, policyID string) error
}

// Policy repository that contains all user operations for this domain
type PolicyRepo interface {
	GetPolicyByName(org string, name string) (*Policy, error)
	AddPolicy(policy Policy) (*Policy, error)
	UpdatePolicy(policy Policy, newName string, newPath string, newUrn string, newStatements []Statement) (*Policy, error)
	RemovePolicy(id string) error
	GetPoliciesFiltered(org string, pathPrefix string) ([]Policy, error)
	GetAttachedGroups(policyID string) ([]Group, error)
}
