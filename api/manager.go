package api

import log "github.com/Sirupsen/logrus"

// TYPE DEFINITIONS

// Interface that all resource types have to implement
type Resource interface {
	// This method must return resource URN
	GetUrn() string
}

// Authorizr API that implements API interfaces using repositories
type AuthAPI struct {
	UserRepo   UserRepo
	GroupRepo  GroupRepo
	PolicyRepo PolicyRepo
	Logger     log.Logger
}

// API INTERFACES WITH AUTHORIZATION

type UserAPI interface {
	// Store user in database. Throw error when parameters are invalid,
	// user already exists or unexpected error happen.
	AddUser(requestInfo RequestInfo, externalId string, path string) (*User, error)

	// Retrieve user from database. Throw error when parameter is invalid,
	// user doesn't exist or unexpected error happen.
	GetUserByExternalID(requestInfo RequestInfo, externalId string) (*User, error)

	// Retrieve user identifiers from database filtered by pathPrefix (optional parameter). Throw error
	// if pathPrefix is invalid or unexpected error happen.
	ListUsers(requestInfo RequestInfo, pathPrefix string) ([]string, error)

	// Update user stored in database with new pathPrefix. Throw error if the input parameters
	// are invalid, user doesn't exist or unexpected error happen.
	UpdateUser(requestInfo RequestInfo, externalId string, newPath string) (*User, error)

	// Remove user stored in database with its group relationships.
	// Throw error if externalId parameter is invalid, user doesn't exist or unexpected error happen.
	RemoveUser(requestInfo RequestInfo, externalId string) error

	// Retrieve groups that belongs to the user. Throw error if externalId parameter is invalid, user
	// doesn't exist or unexpected error happen.
	ListGroupsByUser(requestInfo RequestInfo, externalId string) ([]GroupIdentity, error)
}

type GroupAPI interface {
	// Store group in database. Throw error when the input parameters are invalid,
	// the group already exist or unexpected error happen.
	AddGroup(requestInfo RequestInfo, org string, name string, path string) (*Group, error)

	// Retrieve group from database. Throw error when the input parameters are invalid,
	// group doesn't exist or unexpected error happen.
	GetGroupByName(requestInfo RequestInfo, org string, name string) (*Group, error)

	// Retrieve group identifiers from database filtered by org and pathPrefix parameters. These input parameters are optional.
	// Throw error if the input parameters are invalid or unexpected error happen.
	ListGroups(requestInfo RequestInfo, org string, pathPrefix string) ([]GroupIdentity, error)

	// Update group stored in database with new name and pathPrefix.
	// Throw error if the input parameters are invalid, group to update doesn't exist,
	// target group already exist or unexpected error happen.
	UpdateGroup(requestInfo RequestInfo, org string, groupName string, newName string, newPath string) (*Group, error)

	// Remove group stored in database with its user and policy relationships.
	// Throw error if the input parameters are invalid, the group doesn't exist or unexpected error happen.
	RemoveGroup(requestInfo RequestInfo, org string, name string) error

	// Add new member to group. Throw error if the input parameters are invalid, user doesn't exist,
	// group doesn't exist, user is already a member of the group or unexpected error happen.
	AddMember(requestInfo RequestInfo, externalId string, groupName string, org string) error

	// Remove member from group. Throw error if the input parameters are invalid, user doesn't exist,
	// group doesn't exist, user isn't a member of the group or unexpected error happen.
	RemoveMember(requestInfo RequestInfo, externalId string, groupName string, org string) error

	// List user identifiers that belong to the group. Throw error if the input parameters are invalid,
	// group doesn't exist or unexpected error happen.
	ListMembers(requestInfo RequestInfo, org string, groupName string) ([]string, error)

	// Attach policy to group. Throw error if the input parameters are invalid, policy doesn't exist,
	// group doesn't exist, policy is already attached to the group or unexpected error happen.
	AttachPolicyToGroup(requestInfo RequestInfo, org string, groupName string, policyName string) error

	// Detach policy from group. Throw error if the input parameters are invalid, policy doesn't exist,
	// group doesn't exist, policy isn't attached to the group or unexpected error happen.
	DetachPolicyToGroup(requestInfo RequestInfo, org string, groupName string, policyName string) error

	// Retrieve name of policies that are attached to the group. Throw error if the input parameters are invalid,
	// group doesn't exist or unexpected error happen.
	ListAttachedGroupPolicies(requestInfo RequestInfo, org string, groupName string) ([]string, error)
}

type PolicyAPI interface {
	// Store policy in database. Throw error when the input parameters are invalid,
	// the policy already exist or unexpected error happen.
	AddPolicy(requestInfo RequestInfo, name string, path string, org string, statements []Statement) (*Policy, error)

	// Retrieve policy from database. Throw error when the input parameters are invalid,
	// policy doesn't exist or unexpected error happen.
	GetPolicyByName(requestInfo RequestInfo, org string, name string) (*Policy, error)

	// Retrieve policy identifiers from database filtered by org and pathPrefix parameters. These input parameters are optional.
	// Throw error if the input parameters are invalid or unexpected error happen.
	ListPolicies(requestInfo RequestInfo, org string, pathPrefix string) ([]PolicyIdentity, error)

	// Update policy stored in database with new name, new pathPrefix and new statements.
	// It overrides older statements. Throw error if the input parameters are invalid,
	// policy to update doesn't exist, target policy already exist or unexpected error happen.
	UpdatePolicy(requestInfo RequestInfo, org string, name string, newName string, newPath string,
		newStatements []Statement) (*Policy, error)

	// Remove policy stored in database with its groups relationships.
	// Throw error if the input parameters are invalid, the policy doesn't exist or unexpected error happen.
	RemovePolicy(requestInfo RequestInfo, org string, name string) error

	// Retrieve name of groups that are attached to the policy. Throw error if the input parameters are invalid,
	// policy doesn't exist or unexpected error happen.
	ListAttachedGroups(requestInfo RequestInfo, org string, name string) ([]string, error)
}

type AuthzAPI interface {
	// Retrieve list of authorized user resources filtered according to the input parameters. Throw error
	// if requestInfo doesn't exist, requestInfo doesn't have access to any resources or unexpected error happen.
	GetAuthorizedUsers(requestInfo RequestInfo, resourceUrn string, action string, users []User) ([]User, error)

	// Retrieve list of authorized group resources filtered according to the input parameters. Throw error
	// if requestInfo doesn't exist, requestInfo doesn't have access to any resources or unexpected error happen.
	GetAuthorizedGroups(requestInfo RequestInfo, resourceUrn string, action string, groups []Group) ([]Group, error)

	// Retrieve list of authorized policies resources filtered according to the input parameters. Throw error
	// if requestInfo doesn't exist, requestInfo doesn't have access to any resources or unexpected error happen.
	GetAuthorizedPolicies(requestInfo RequestInfo, resourceUrn string, action string, policies []Policy) ([]Policy, error)

	// Retrieve list of authorized external resources filtered according to the input parameters. Throw error
	// if requestInfo doesn't exist, requestInfo doesn't have access to any resources or unexpected error happen.
	GetAuthorizedExternalResources(requestInfo RequestInfo, action string, resources []string) ([]string, error)
}

// REPOSITORY INTERFACES

// User repository that contains all database operations
type UserRepo interface {
	// Store user in database if there aren't errors.
	AddUser(user User) (*User, error)

	// Retrieve user from database if it exists. Otherwise it throws an error.
	GetUserByExternalID(id string) (*User, error)

	// Retrieve user list from database filtered by pathPrefix optional parameter. Throw error
	// if there are problems with database.
	GetUsersFiltered(pathPrefix string) ([]User, error)

	// Update user stored in database with new pathPrefix. Throw error if the database restrictions
	// are not satisfied or unexpected error happen.
	UpdateUser(user User, newPath string, newUrn string) (*User, error)

	// Remove user stored in database with its group relationships.
	// Throw error if there are problems during transactions.
	RemoveUser(id string) error

	// Retrieve groups that belong to the user. Throw error
	// if there are problems with database.
	GetGroupsByUserID(id string) ([]Group, error)
}

// Group repository that contains all database operations
type GroupRepo interface {
	// Store group in database if there aren't errors.
	AddGroup(group Group) (*Group, error)

	// Retrieve group from database if it exists. Otherwise it throws an error.
	GetGroupByName(org string, name string) (*Group, error)

	// Retrieve groups from database filtered by org and pathPrefix optional parameters. Throw error
	// if there are problems with database.
	GetGroupsFiltered(org string, pathPrefix string) ([]Group, error)

	// Update group stored in database with new name and pathPrefix.
	// Throw error if there are problems with database.
	UpdateGroup(group Group, newName string, newPath string, newUrn string) (*Group, error)

	// Remove group stored in database with its user and policy relationships.
	// Throw error if there are problems during transactions.
	RemoveGroup(groupID string) error

	// Add new member to group. It doesn't check restrictions about existence of group or user. It throws
	// errors if there are problems with database.
	AddMember(userID string, groupID string) error

	// Remove member from group. It doesn't check restrictions about existence of group or user. It throws
	// errors if there are problems with database.
	RemoveMember(userID string, groupID string) error

	// Check if user is member of group. It returns true if at least one relation exists. It throws
	// errors if there are problems with database.
	IsMemberOfGroup(userID string, groupID string) (bool, error)

	// Retrieve users that belong to the group. Throw error if there are problems with database.
	GetGroupMembers(groupID string) ([]User, error)

	// Attach policy to group. It doesn't check restrictions about existence of group or policy. It throws
	// errors if there are problems with database.
	AttachPolicy(groupID string, policyID string) error

	// Detach policy from group. It doesn't check restrictions about existence of group or policy. It throws
	// errors if there are problems with database.
	DetachPolicy(groupID string, policyID string) error

	// Check if policy is attached to group. It returns true if at least one relation exists. It throws
	// errors if there are problems with database.
	IsAttachedToGroup(groupID string, policyID string) (bool, error)

	// Retrieve policies that are attached to the group. Throw error if there are problems with database.
	GetAttachedPolicies(groupID string) ([]Policy, error)
}

// Policy repository that contains all database operations
type PolicyRepo interface {
	// Store policy in database if there aren't errors.
	AddPolicy(policy Policy) (*Policy, error)

	// Retrieve policy from database if it exists. Otherwise it throws an error.
	GetPolicyByName(org string, name string) (*Policy, error)

	// Retrieve policies from database filtered by org and pathPrefix optional parameters. Throw error
	// if there are problems with database.
	GetPoliciesFiltered(org string, pathPrefix string) ([]Policy, error)

	// Update policy stored in database with new name and pathPrefix. Also it overrides statements.
	// Throw error if there are problems with database.
	UpdatePolicy(policy Policy, newName string, newPath string, newUrn string, newStatements []Statement) (*Policy, error)

	// Remove policy stored in database with its groups relationships.
	// Throw error if there are problems during transactions.
	RemovePolicy(id string) error

	// Retrieve groups that are attached to the policy. Throw error if there are problems with database.
	GetAttachedGroups(policyID string) ([]Group, error)
}
