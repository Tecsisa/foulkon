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
	// Stores an user in database. Throw error when parameters are invalid,
	// user already exists or unexpected error happen.
	AddUser(authenticatedUser AuthenticatedUser, externalId string, path string) (*User, error)

	// Retrieve an user from database. Throw error when parameter is invalid,
	// user doesn't exist or unexpected error happen.
	GetUserByExternalID(authenticatedUser AuthenticatedUser, externalId string) (*User, error)

	// Retrieve the user identifiers from database filtered by pathPrefix (optional parameter). Throw error
	// if pathPrefix is invalid or unexpected error happen.
	ListUsers(authenticatedUser AuthenticatedUser, pathPrefix string) ([]string, error)

	// Update an user stored in database with a new pathPrefix. Throw error if the input parameters
	// are invalid, user doesn't exist or unexpected error happen.
	UpdateUser(authenticatedUser AuthenticatedUser, externalId string, newPath string) (*User, error)

	// Remove an user stored in database with its group relationships.
	// Throw error if externalId parameter is invalid, user doesn't exist or unexpected error happen.
	RemoveUser(authenticatedUser AuthenticatedUser, externalId string) error

	// Retrieve groups that belongs to the user. Throw error if externalId parameter is invalid, user
	// doesn't exist or unexpected error happen.
	ListGroupsByUser(authenticatedUser AuthenticatedUser, externalId string) ([]GroupIdentity, error)
}

type GroupAPI interface {
	// Store a group in database. Throw error when the input parameters are invalid,
	// the group already exist or unexpected error happen.
	AddGroup(authenticatedUser AuthenticatedUser, org string, name string, path string) (*Group, error)

	// Retrieve a group from database. Throw error when the input parameters are invalid,
	// group doesn't exist or unexpected error happen.
	GetGroupByName(authenticatedUser AuthenticatedUser, org string, name string) (*Group, error)

	// Retrieve group identifiers from database filtered by org and pathPrefix parameters. These input parameters are optional.
	// Throw error if the input parameters are invalid or unexpected error happen.
	ListGroups(authenticatedUser AuthenticatedUser, org string, pathPrefix string) ([]GroupIdentity, error)

	// Update a group stored in database with a new name and pathPrefix.
	// Throw error if the input parameters are invalid, group to update doesn't exist,
	// target group already exist or unexpected error happen.
	UpdateGroup(authenticatedUser AuthenticatedUser, org string, groupName string, newName string, newPath string) (*Group, error)

	// Remove a group stored in database with its user and policy relationships.
	// Throw error if the input parameters are invalid, the group doesn't exist or unexpected error happen.
	RemoveGroup(authenticatedUser AuthenticatedUser, org string, name string) error

	// Add a new member to a group. Throw error if the input parameters are invalid, user doesn't exist,
	// group doesn't exist, user is already a member of the group or unexpected error happen.
	AddMember(authenticatedUser AuthenticatedUser, externalId string, groupName string, org string) error

	// Remove a member from a group. Throw error if the input parameters are invalid, user doesn't exist,
	// group doesn't exist, user isn't a member of the group or unexpected error happen.
	RemoveMember(authenticatedUser AuthenticatedUser, externalId string, groupName string, org string) error

	// List the user identifiers that belong to the group. Throw error if the input parameters are invalid,
	// group doesn't exist or unexpected error happen.
	ListMembers(authenticatedUser AuthenticatedUser, org string, groupName string) ([]string, error)

	// Attach a policy to a group. Throw error if the input parameters are invalid, policy doesn't exist,
	// group doesn't exist, policy is already attached to the group or unexpected error happen.
	AttachPolicyToGroup(authenticatedUser AuthenticatedUser, org string, groupName string, policyName string) error

	// Detach a policy from a group. Throw error if the input parameters are invalid, policy doesn't exist,
	// group doesn't exist, policy isn't attached to the group or unexpected error happen.
	DetachPolicyToGroup(authenticatedUser AuthenticatedUser, org string, groupName string, policyName string) error

	// Retrieve the name of policies that are attached to the group. Throw error if the input parameters are invalid,
	// group doesn't exist or unexpected error happen.
	ListAttachedGroupPolicies(authenticatedUser AuthenticatedUser, org string, groupName string) ([]string, error)
}

type PolicyAPI interface {
	// Store a policy in database. Throw error when the input parameters are invalid,
	// the policy already exist or unexpected error happen.
	AddPolicy(authenticatedUser AuthenticatedUser, name string, path string, org string, statements []Statement) (*Policy, error)

	// Retrieve a policy from database. Throw error when the input parameters are invalid,
	// policy doesn't exist or unexpected error happen.
	GetPolicyByName(authenticatedUser AuthenticatedUser, org string, name string) (*Policy, error)

	// Retrieve policy identifiers from database filtered by org and pathPrefix parameters. These input parameters are optional.
	// Throw error if the input parameters are invalid or unexpected error happen.
	ListPolicies(authenticatedUser AuthenticatedUser, org string, pathPrefix string) ([]PolicyIdentity, error)

	// Update a policy stored in database with a new name, new pathPrefix and new statements.
	// It overrides older statements. Throw error if the input parameters are invalid,
	// policy to update doesn't exist, target policy already exist or unexpected error happen.
	UpdatePolicy(authenticatedUser AuthenticatedUser, org string, name string, newName string, newPath string,
		newStatements []Statement) (*Policy, error)

	// Remove a policy stored in database with its groups relationships.
	// Throw error if the input parameters are invalid, the policy doesn't exist or unexpected error happen.
	RemovePolicy(authenticatedUser AuthenticatedUser, org string, name string) error

	// Retrieve the name of groups that are attached to the policy. Throw error if the input parameters are invalid,
	// policy doesn't exist or unexpected error happen.
	ListAttachedGroups(authenticatedUser AuthenticatedUser, org string, name string) ([]string, error)
}

type AuthzAPI interface {
	// Retrieve a list of authorized user resources filtered according to the input parameters. Throw error
	// if authenticatedUser doesn't exist, authenticatedUser doesn't have access to any resources or unexpected error happen.
	GetAuthorizedUsers(authenticatedUser AuthenticatedUser, resourceUrn string, action string, users []User) ([]User, error)

	// Retrieve a list of authorized group resources filtered according to the input parameters. Throw error
	// if authenticatedUser doesn't exist, authenticatedUser doesn't have access to any resources or unexpected error happen.
	GetAuthorizedGroups(authenticatedUser AuthenticatedUser, resourceUrn string, action string, groups []Group) ([]Group, error)

	// Retrieve a list of authorized policies resources filtered according to the input parameters. Throw error
	// if authenticatedUser doesn't exist, authenticatedUser doesn't have access to any resources or unexpected error happen.
	GetAuthorizedPolicies(authenticatedUser AuthenticatedUser, resourceUrn string, action string, policies []Policy) ([]Policy, error)

	// Retrieve a list of authorized external resources filtered according to the input parameters. Throw error
	// if authenticatedUser doesn't exist, authenticatedUser doesn't have access to any resources or unexpected error happen.
	GetAuthorizedExternalResources(authenticatedUser AuthenticatedUser, action string, resources []string) ([]string, error)
}

// REPOSITORY INTERFACES

// User repository that contains all database operations
type UserRepo interface {
	// Store an user in database if there aren't errors.
	AddUser(user User) (*User, error)

	// Retrieve an user from database if it exists. Otherwise it throws an error.
	GetUserByExternalID(id string) (*User, error)

	// Retrieve the user list from database filtered by pathPrefix optional parameter. Throw error
	// if there are problems with database.
	GetUsersFiltered(pathPrefix string) ([]User, error)

	// Update an user stored in database with a new pathPrefix. Throw error if the database restrictions
	// are not satisfied or unexpected error happen.
	UpdateUser(user User, newPath string, newUrn string) (*User, error)

	// Remove an user stored in database with its group relationships.
	// Throw error if there are problems during transactions.
	RemoveUser(id string) error

	// Retrieve groups that belong to the user. Throw error
	// if there are problems with database.
	GetGroupsByUserID(id string) ([]Group, error)
}

// Group repository that contains all database operations
type GroupRepo interface {
	// Store a group in database if there aren't errors.
	AddGroup(group Group) (*Group, error)

	// Retrieve a group from database if it exists. Otherwise it throws an error.
	GetGroupByName(org string, name string) (*Group, error)

	// Retrieve groups from database filtered by org and pathPrefix optional parameters. Throw error
	// if there are problems with database.
	GetGroupsFiltered(org string, pathPrefix string) ([]Group, error)

	// Update a group stored in database with a new name and pathPrefix.
	// Throw error if there are problems with database.
	UpdateGroup(group Group, newName string, newPath string, newUrn string) (*Group, error)

	// Remove a group stored in database with its user and policy relationships.
	// Throw error if there are problems during transactions.
	RemoveGroup(groupID string) error

	// Add a new member to a group. It doesn't check restrictions about existence of group or user. It throws
	// errors if there are problems with database.
	AddMember(userID string, groupID string) error

	// Remove a member from a group. It doesn't check restrictions about existence of group or user. It throws
	// errors if there are problems with database.
	RemoveMember(userID string, groupID string) error

	// Check if the user is member of the group. It returns true if at least one relation exists. It throws
	// errors if there are problems with database.
	IsMemberOfGroup(userID string, groupID string) (bool, error)

	// Retrieve users that belong to the group. Throw error if there are problems with database.
	GetGroupMembers(groupID string) ([]User, error)

	// Attach a policy to a group. It doesn't check restrictions about existence of group or policy. It throws
	// errors if there are problems with database.
	AttachPolicy(groupID string, policyID string) error

	// Detach a policy from a group. It doesn't check restrictions about existence of group or policy. It throws
	// errors if there are problems with database.
	DetachPolicy(groupID string, policyID string) error

	// Check if the policy is attached to the group. It returns true if at least one relation exists. It throws
	// errors if there are problems with database.
	IsAttachedToGroup(groupID string, policyID string) (bool, error)

	// Retrieve the policies that are attached to the group. Throw error if there are problems with database.
	GetAttachedPolicies(groupID string) ([]Policy, error)
}

// Policy repository that contains all database operations
type PolicyRepo interface {
	// Store a policy in database if there aren't errors.
	AddPolicy(policy Policy) (*Policy, error)

	// Retrieve a policy from database if it exists. Otherwise it throws an error.
	GetPolicyByName(org string, name string) (*Policy, error)

	// Retrieve policies from database filtered by org and pathPrefix optional parameters. Throw error
	// if there are problems with database.
	GetPoliciesFiltered(org string, pathPrefix string) ([]Policy, error)

	// Update a policy stored in database with a new name and pathPrefix. Also it overrides statements.
	// Throw error if there are problems with database.
	UpdatePolicy(policy Policy, newName string, newPath string, newUrn string, newStatements []Statement) (*Policy, error)

	// Remove a policy stored in database with its groups relationships.
	// Throw error if there are problems during transactions.
	RemovePolicy(id string) error

	// Retrieve the groups that are attached to the policy. Throw error if there are problems with database.
	GetAttachedGroups(policyID string) ([]Group, error)
}
