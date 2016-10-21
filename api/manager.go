package api

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

// TYPE DEFINITIONS

// Resource interface that all resource types have to implement
type Resource interface {
	// This method must return resource URN
	GetUrn() string
}

// UserGroupRelation interface for User-Group relationships
type UserGroupRelation interface {
	GetUser() *User
	GetGroup() *Group
	GetDate() time.Time
}

// PolicyGroupRelation interface for Policy-Group relationships
type PolicyGroupRelation interface {
	GetGroup() *Group
	GetPolicy() *Policy
	GetDate() time.Time
}

// WorkerAPI that implements API interfaces using repositories
type WorkerAPI struct {
	UserRepo   UserRepo
	GroupRepo  GroupRepo
	PolicyRepo PolicyRepo
	ProxyRepo  ProxyRepo
	Logger     *log.Logger
}

// ProxyAPI that implements API interfaces using repositories
type ProxyAPI struct {
	ProxyRepo ProxyRepo
	Logger    *log.Logger
}

// Filter properties for database search
type Filter struct {
	PathPrefix string
	Org        string
	ExternalID string
	PolicyName string
	GroupName  string
	// Pagination
	Offset int
	Limit  int
	// Sorting
	OrderBy string
}

// API INTERFACES WITH AUTHORIZATION

// UserAPI interface
type UserAPI interface {
	// Store user in database. Throw error when parameters are invalid,
	// user already exists or unexpected error happen.
	AddUser(requestInfo RequestInfo, externalId string, path string) (*User, error)

	// Retrieve user from database. Throw error when parameter is invalid,
	// user doesn't exist or unexpected error happen.
	GetUserByExternalID(requestInfo RequestInfo, externalId string) (*User, error)

	// Retrieve user identifiers from database filtered by pathPrefix (optional parameter). Throw error
	// if pathPrefix is invalid or unexpected error happen.
	ListUsers(requestInfo RequestInfo, filter *Filter) ([]string, int, error)

	// Update user stored in database with new pathPrefix. Throw error if the input parameters
	// are invalid, user doesn't exist or unexpected error happen.
	UpdateUser(requestInfo RequestInfo, externalId string, newPath string) (*User, error)

	// Remove user stored in database with its group relationships.
	// Throw error if externalId parameter is invalid, user doesn't exist or unexpected error happen.
	RemoveUser(requestInfo RequestInfo, externalId string) error

	// Retrieve groups that belongs to the user. Throw error if externalId parameter is invalid, user
	// doesn't exist or unexpected error happen.
	ListGroupsByUser(requestInfo RequestInfo, filter *Filter) ([]UserGroups, int, error)
}

// GroupAPI interface
type GroupAPI interface {
	// Store group in database. Throw error when the input parameters are invalid,
	// the group already exist or unexpected error happen.
	AddGroup(requestInfo RequestInfo, org string, name string, path string) (*Group, error)

	// Retrieve group from database. Throw error when the input parameters are invalid,
	// group doesn't exist or unexpected error happen.
	GetGroupByName(requestInfo RequestInfo, org string, name string) (*Group, error)

	// Retrieve group identifiers from database filtered by org and pathPrefix parameters. These input parameters are optional.
	// Throw error if the input parameters are invalid or unexpected error happen.
	ListGroups(requestInfo RequestInfo, filter *Filter) ([]GroupIdentity, int, error)

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
	ListMembers(requestInfo RequestInfo, filter *Filter) ([]GroupMembers, int, error)

	// Attach policy to group. Throw error if the input parameters are invalid, policy doesn't exist,
	// group doesn't exist, policy is already attached to the group or unexpected error happen.
	AttachPolicyToGroup(requestInfo RequestInfo, org string, groupName string, policyName string) error

	// Detach policy from group. Throw error if the input parameters are invalid, policy doesn't exist,
	// group doesn't exist, policy isn't attached to the group or unexpected error happen.
	DetachPolicyToGroup(requestInfo RequestInfo, org string, groupName string, policyName string) error

	// Retrieve policies that are attached to the group. Throw error if the input parameters are invalid,
	// group doesn't exist or unexpected error happen.
	ListAttachedGroupPolicies(requestInfo RequestInfo, filter *Filter) ([]GroupPolicies, int, error)
}

// PolicyAPI interface
type PolicyAPI interface {
	// Store policy in database. Throw error when the input parameters are invalid,
	// the policy already exist or unexpected error happen.
	AddPolicy(requestInfo RequestInfo, name string, path string, org string, statements []Statement) (*Policy, error)

	// Retrieve policy from database. Throw error when the input parameters are invalid,
	// policy doesn't exist or unexpected error happen.
	GetPolicyByName(requestInfo RequestInfo, org string, name string) (*Policy, error)

	// Retrieve policy identifiers from database filtered by org and pathPrefix parameters. These input parameters are optional.
	// Throw error if the input parameters are invalid or unexpected error happen.
	ListPolicies(requestInfo RequestInfo, filter *Filter) ([]PolicyIdentity, int, error)

	// Update policy stored in database with new name, new pathPrefix and new statements.
	// It overrides older statements. Throw error if the input parameters are invalid,
	// policy to update doesn't exist, target policy already exist or unexpected error happen.
	UpdatePolicy(requestInfo RequestInfo, org string, name string, newName string, newPath string,
		newStatements []Statement) (*Policy, error)

	// Remove policy stored in database with its groups relationships.
	// Throw error if the input parameters are invalid, the policy doesn't exist or unexpected error happen.
	RemovePolicy(requestInfo RequestInfo, org string, name string) error

	// Retrieve groups that are attached to the policy. Throw error if the input parameters are invalid,
	// policy doesn't exist or unexpected error happen.
	ListAttachedGroups(requestInfo RequestInfo, filter *Filter) ([]PolicyGroups, int, error)
}

// AuthzAPI interface
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

// ProxyResourcesAPI interface to manage proxy resources
type ProxyResourcesAPI interface {
	// Retrieve list of proxy resources.
	GetProxyResources() ([]ProxyResource, error)
}

// REPOSITORY INTERFACES

// UserRepo contains all database operations
type UserRepo interface {
	// Store user in database if there aren't errors.
	AddUser(user User) (*User, error)

	// Retrieve user from database if it exists. Otherwise it throws an error.
	GetUserByExternalID(id string) (*User, error)

	// Retrieve user list from database filtered by pathPrefix optional parameter. Throw error
	// if there are problems with database.
	GetUsersFiltered(filter *Filter) ([]User, int, error)

	// Update user stored in database with new fields. Throw error if the database restrictions
	// are not satisfied or unexpected error happen.
	UpdateUser(user User) (*User, error)

	// Remove user stored in database with its group relationships.
	// Throw error if there are problems during transactions.
	RemoveUser(id string) error

	// Retrieve groups that belong to the user. Throw error
	// if there are problems with database.
	GetGroupsByUserID(id string, filter *Filter) ([]UserGroupRelation, int, error)

	// OrderByValidColumns returns valid columns that you can use in OrderBy
	OrderByValidColumns(action string) []string
}

// GroupRepo contains all database operations
type GroupRepo interface {
	// Store group in database if there aren't errors.
	AddGroup(group Group) (*Group, error)

	// Retrieve group from database if it exists. Otherwise it throws an error.
	GetGroupByName(org string, name string) (*Group, error)

	// Retrieve groups from database filtered by org and pathPrefix optional parameters. Throw error
	// if there are problems with database.
	GetGroupsFiltered(filter *Filter) ([]Group, int, error)

	// Update group stored in database with new fields.
	// Throw error if there are problems with database.
	UpdateGroup(group Group) (*Group, error)

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
	GetGroupMembers(groupID string, filter *Filter) ([]UserGroupRelation, int, error)

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
	GetAttachedPolicies(groupID string, filter *Filter) ([]PolicyGroupRelation, int, error)

	// OrderByValidColumns returns valid columns that you can use in OrderBy
	OrderByValidColumns(action string) []string
}

// PolicyRepo contains all database operations
type PolicyRepo interface {
	// Store policy in database if there aren't errors.
	AddPolicy(policy Policy) (*Policy, error)

	// Retrieve policy from database if it exists. Otherwise it throws an error.
	GetPolicyByName(org string, name string) (*Policy, error)

	// Retrieve policies from database filtered by org and pathPrefix optional parameters. Throw error
	// if there are problems with database.
	GetPoliciesFiltered(filter *Filter) ([]Policy, int, error)

	// Update policy stored in database with new fields. Also it overrides statements if it has.
	// Throw error if there are problems with database.
	UpdatePolicy(policy Policy) (*Policy, error)

	// Remove policy stored in database with its groups relationships.
	// Throw error if there are problems during transactions.
	RemovePolicy(id string) error

	// Retrieve groups that are attached to the policy. Throw error if there are problems with database.
	GetAttachedGroups(policyID string, filter *Filter) ([]PolicyGroupRelation, int, error)

	// OrderByValidColumns returns valid columns that you can use in OrderBy
	OrderByValidColumns(action string) []string
}

// ProxyRepo contains all database operations
type ProxyRepo interface {
	// Retrieve proxy resources from database. Otherwise it throws an error.
	GetProxyResources() ([]ProxyResource, error)
}
