package api

import log "github.com/Sirupsen/logrus"

// Authorizr API struct with repositories
type AuthAPI struct {
	UserRepo   UserRepo
	GroupRepo  GroupRepo
	PolicyRepo PolicyRepo
	logger     log.Logger
}

// User repository that contains all user operations for this domain
type UserRepo interface {
	// This method get a user with specified External ID.
	// If user exists, it will return the user with error param as nil
	// If user doesn't exists, it will return the error code database.USER_NOT_FOUND
	// If there is an error, it will return error param with associated error message
	// and error code database.INTERNAL_ERROR
	GetUserByExternalID(id string) (*User, error)

	// This method store a user.
	// If there are a problem inserting user it will return an database.Error error
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
	GetPoliciesAttached(groupID string) ([]Policy, error)
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
	GetAllPolicyGroupRelation(policyID string) ([]Group, error)
}
