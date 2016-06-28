package http

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/auth"
	"github.com/tecsisa/authorizr/authorizr"
)

const (
	// USER API METHODS
	AddUserMethod             = "AddUser"
	GetUserByExternalIdMethod = "GetUserByExternalId"
	GetUserListMethod         = "GetUserList"
	UpdateUserMethod          = "UpdateUser"
	RemoveUserByIdMethod      = "RemoveUserById"
	GetGroupsByUserIdMethod   = "GetGroupsByUserId"

	// GROUP API METHODS
	AddGroupMethod                  = "AddGroup"
	GetGroupByNameMethod            = "GetGroupByName"
	GetGroupListMethod              = "GetGroupList"
	UpdateGroupMethod               = "UpdateGroup"
	RemoveGroupMethod               = "RemoveGroup"
	AddMemberMethod                 = "AddMember"
	RemoveMemberMethod              = "RemoveMember"
	ListMembersMethod               = "ListMembers"
	AttachPolicyToGroupMethod       = "AttachPolicyToGroup"
	DetachPolicyToGroupMethod       = "DetachPolicyToGroup"
	ListAttachedGroupPoliciesMethod = "ListAttachedGroupPolicies"

	// POLICY API METHODS
	AddPolicyMethod         = "AddPolicy"
	GetPolicyByNameMethod   = "GetPolicyByName"
	GetPolicyListMethod     = "GetPolicyList"
	UpdatePolicyMethod      = "UpdatePolicy"
	DeletePolicyMethod      = "DeletePolicy"
	GetAttachedGroupsMethod = "GetAttachedGroups"

	// AUTHZ API
	GetUsersAuthorizedMethod             = "GetUsersAuthorized"
	GetAuthorizedGroupsMethod            = "GetAuthorizedGroups"
	GetAuthorizedPoliciesMethod          = "GetAuthorizedPolicies"
	GetAuthorizedExternalResourcesMethod = "GetAuthorizedExternalResources"
)

// Test server used to test handlers
var server *httptest.Server
var testApi *TestAPI

// Test API that implements all api manager interfaces
type TestAPI struct {
	ArgsIn       map[string][]interface{}
	ArgsOut      map[string][]interface{}
	SpecialFuncs map[string]interface{}
}

// Aux connector
type TestConnector struct {
	userID string
}

func (tc TestConnector) Authenticate(h http.Handler) http.Handler {
	return h
}

func (tc TestConnector) RetrieveUserID(r http.Request) string {
	return tc.userID
}

// Main Test that executes at first time and create all necessary data to work
func TestMain(m *testing.M) {
	flag.Parse()

	// Create logger
	logger := &log.Logger{
		Out:       os.Stdout,
		Formatter: &log.TextFormatter{},
		Hooks:     make(log.LevelHooks),
		Level:     log.DebugLevel,
	}

	testApi = makeTestApi()

	// Instantiate Auth Connector
	authConnector := &TestConnector{
		userID: "userID",
	}

	adminUser := "admin"
	adminPassword := "admin"

	// Create authenticator
	authenticator := auth.NewAuthenticator(authConnector, adminUser, adminPassword)

	// Return created core
	worker := &authorizr.Worker{
		Logger:        logger,
		Authenticator: authenticator,
		UserApi:       testApi,
		GroupApi:      testApi,
		PolicyApi:     testApi,
		AuthzApi:      testApi,
	}

	server = httptest.NewServer(WorkerHandlerRouter(worker))

	// Run tests
	result := m.Run()

	// Exit tests.
	os.Exit(result)
}

// func that initializes the TestAPI
func makeTestApi() *TestAPI {
	testApi := &TestAPI{
		ArgsIn:       make(map[string][]interface{}),
		ArgsOut:      make(map[string][]interface{}),
		SpecialFuncs: make(map[string]interface{}),
	}

	testApi.ArgsIn[AddUserMethod] = make([]interface{}, 3)
	testApi.ArgsIn[GetUserByExternalIdMethod] = make([]interface{}, 2)
	testApi.ArgsIn[GetUserListMethod] = make([]interface{}, 2)
	testApi.ArgsIn[UpdateUserMethod] = make([]interface{}, 3)
	testApi.ArgsIn[RemoveUserByIdMethod] = make([]interface{}, 2)
	testApi.ArgsIn[GetGroupsByUserIdMethod] = make([]interface{}, 2)

	testApi.ArgsIn[AddGroupMethod] = make([]interface{}, 4)
	testApi.ArgsIn[GetGroupByNameMethod] = make([]interface{}, 3)
	testApi.ArgsIn[GetGroupListMethod] = make([]interface{}, 3)
	testApi.ArgsIn[UpdateGroupMethod] = make([]interface{}, 5)
	testApi.ArgsIn[RemoveGroupMethod] = make([]interface{}, 3)
	testApi.ArgsIn[AddMemberMethod] = make([]interface{}, 4)
	testApi.ArgsIn[RemoveMemberMethod] = make([]interface{}, 4)
	testApi.ArgsIn[ListMembersMethod] = make([]interface{}, 3)
	testApi.ArgsIn[AttachPolicyToGroupMethod] = make([]interface{}, 4)
	testApi.ArgsIn[DetachPolicyToGroupMethod] = make([]interface{}, 4)
	testApi.ArgsIn[ListAttachedGroupPoliciesMethod] = make([]interface{}, 3)

	testApi.ArgsIn[AddPolicyMethod] = make([]interface{}, 5)
	testApi.ArgsIn[GetPolicyByNameMethod] = make([]interface{}, 3)
	testApi.ArgsIn[GetPolicyListMethod] = make([]interface{}, 3)
	testApi.ArgsIn[UpdatePolicyMethod] = make([]interface{}, 6)
	testApi.ArgsIn[DeletePolicyMethod] = make([]interface{}, 3)
	testApi.ArgsIn[GetAttachedGroupsMethod] = make([]interface{}, 3)

	testApi.ArgsIn[GetUsersAuthorizedMethod] = make([]interface{}, 4)
	testApi.ArgsIn[GetAuthorizedGroupsMethod] = make([]interface{}, 4)
	testApi.ArgsIn[GetAuthorizedPoliciesMethod] = make([]interface{}, 4)
	testApi.ArgsIn[GetAuthorizedExternalResourcesMethod] = make([]interface{}, 3)

	testApi.ArgsOut[AddUserMethod] = make([]interface{}, 2)
	testApi.ArgsOut[GetUserByExternalIdMethod] = make([]interface{}, 2)
	testApi.ArgsOut[GetUserListMethod] = make([]interface{}, 2)
	testApi.ArgsOut[UpdateUserMethod] = make([]interface{}, 2)
	testApi.ArgsOut[RemoveUserByIdMethod] = make([]interface{}, 1)
	testApi.ArgsOut[GetGroupsByUserIdMethod] = make([]interface{}, 2)

	testApi.ArgsOut[AddGroupMethod] = make([]interface{}, 2)
	testApi.ArgsOut[GetGroupByNameMethod] = make([]interface{}, 2)
	testApi.ArgsOut[GetGroupListMethod] = make([]interface{}, 2)
	testApi.ArgsOut[UpdateGroupMethod] = make([]interface{}, 2)
	testApi.ArgsOut[RemoveGroupMethod] = make([]interface{}, 1)
	testApi.ArgsOut[AddMemberMethod] = make([]interface{}, 1)
	testApi.ArgsOut[RemoveMemberMethod] = make([]interface{}, 1)
	testApi.ArgsOut[ListMembersMethod] = make([]interface{}, 2)
	testApi.ArgsOut[AttachPolicyToGroupMethod] = make([]interface{}, 1)
	testApi.ArgsOut[DetachPolicyToGroupMethod] = make([]interface{}, 1)
	testApi.ArgsOut[ListAttachedGroupPoliciesMethod] = make([]interface{}, 2)

	testApi.ArgsOut[AddPolicyMethod] = make([]interface{}, 2)
	testApi.ArgsOut[GetPolicyByNameMethod] = make([]interface{}, 2)
	testApi.ArgsOut[GetPolicyListMethod] = make([]interface{}, 2)
	testApi.ArgsOut[UpdatePolicyMethod] = make([]interface{}, 2)
	testApi.ArgsOut[DeletePolicyMethod] = make([]interface{}, 1)
	testApi.ArgsOut[GetAttachedGroupsMethod] = make([]interface{}, 2)

	testApi.ArgsOut[GetUsersAuthorizedMethod] = make([]interface{}, 2)
	testApi.ArgsOut[GetAuthorizedGroupsMethod] = make([]interface{}, 2)
	testApi.ArgsOut[GetAuthorizedPoliciesMethod] = make([]interface{}, 2)
	testApi.ArgsOut[GetAuthorizedExternalResourcesMethod] = make([]interface{}, 2)

	return testApi
}

// USER API

func (t TestAPI) AddUser(authenticatedUser api.AuthenticatedUser, externalID string, path string) (*api.User, error) {
	t.ArgsIn[AddUserMethod][0] = authenticatedUser
	t.ArgsIn[AddUserMethod][1] = externalID
	t.ArgsIn[AddUserMethod][2] = path
	var user *api.User
	if t.ArgsOut[AddUserMethod][0] != nil {
		user = t.ArgsOut[AddUserMethod][0].(*api.User)
	}
	var err error
	if t.ArgsOut[AddUserMethod][1] != nil {
		err = t.ArgsOut[AddUserMethod][1].(error)
	}
	return user, err
}

func (t TestAPI) GetUserByExternalId(authenticatedUser api.AuthenticatedUser, id string) (*api.User, error) {
	t.ArgsIn[GetUserByExternalIdMethod][0] = authenticatedUser
	t.ArgsIn[GetUserByExternalIdMethod][1] = id
	var user *api.User
	if t.ArgsOut[GetUserByExternalIdMethod][0] != nil {
		user = t.ArgsOut[GetUserByExternalIdMethod][0].(*api.User)
	}
	var err error
	if t.ArgsOut[GetUserByExternalIdMethod][1] != nil {
		err = t.ArgsOut[GetUserByExternalIdMethod][1].(error)
	}
	return user, err
}

func (t TestAPI) GetUserList(authenticatedUser api.AuthenticatedUser, pathPrefix string) ([]string, error) {
	t.ArgsIn[GetUserListMethod][0] = authenticatedUser
	t.ArgsIn[GetUserListMethod][1] = pathPrefix
	var externalIDs []string
	if t.ArgsOut[GetUserListMethod][0] != nil {
		externalIDs = t.ArgsOut[GetUserListMethod][0].([]string)
	}
	var err error
	if t.ArgsOut[GetUserListMethod][1] != nil {
		err = t.ArgsOut[GetUserListMethod][1].(error)
	}
	return externalIDs, err
}

func (t TestAPI) UpdateUser(authenticatedUser api.AuthenticatedUser, externalID string, newPath string) (*api.User, error) {
	t.ArgsIn[UpdateUserMethod][0] = authenticatedUser
	t.ArgsIn[UpdateUserMethod][1] = externalID
	t.ArgsIn[UpdateUserMethod][2] = newPath
	var user *api.User
	if t.ArgsOut[UpdateUserMethod][0] != nil {
		user = t.ArgsOut[UpdateUserMethod][0].(*api.User)
	}
	var err error
	if t.ArgsOut[UpdateUserMethod][1] != nil {
		err = t.ArgsOut[UpdateUserMethod][1].(error)
	}
	return user, err
}

func (t TestAPI) RemoveUserById(authenticatedUser api.AuthenticatedUser, id string) error {
	t.ArgsIn[RemoveUserByIdMethod][0] = authenticatedUser
	t.ArgsIn[RemoveUserByIdMethod][1] = id
	var err error
	if t.ArgsOut[RemoveUserByIdMethod][0] != nil {
		err = t.ArgsOut[RemoveUserByIdMethod][0].(error)
	}
	return err
}

func (t TestAPI) GetGroupsByUserId(authenticatedUser api.AuthenticatedUser, id string) ([]api.GroupIdentity, error) {
	t.ArgsIn[GetGroupsByUserIdMethod][0] = authenticatedUser
	t.ArgsIn[GetGroupsByUserIdMethod][1] = id
	var groups []api.GroupIdentity
	if t.ArgsOut[GetGroupsByUserIdMethod][0] != nil {
		groups = t.ArgsOut[GetGroupsByUserIdMethod][0].([]api.GroupIdentity)
	}
	var err error
	if t.ArgsOut[GetGroupsByUserIdMethod][1] != nil {
		err = t.ArgsOut[GetGroupsByUserIdMethod][1].(error)
	}
	return groups, err
}

// GROUP API

func (t TestAPI) AddGroup(authenticatedUser api.AuthenticatedUser, org string, name string, path string) (*api.Group, error) {
	t.ArgsIn[AddGroupMethod][0] = authenticatedUser
	t.ArgsIn[AddGroupMethod][1] = org
	t.ArgsIn[AddGroupMethod][2] = name
	t.ArgsIn[AddGroupMethod][3] = path
	var group *api.Group
	if t.ArgsOut[AddGroupMethod][0] != nil {
		group = t.ArgsOut[AddGroupMethod][0].(*api.Group)
	}
	var err error
	if t.ArgsOut[AddGroupMethod][1] != nil {
		err = t.ArgsOut[AddGroupMethod][1].(error)
	}
	return group, err
}

func (t TestAPI) GetGroupByName(authenticatedUser api.AuthenticatedUser, org string, name string) (*api.Group, error) {
	t.ArgsIn[GetGroupByNameMethod][0] = authenticatedUser
	t.ArgsIn[GetGroupByNameMethod][1] = org
	t.ArgsIn[GetGroupByNameMethod][2] = name
	var group *api.Group
	if t.ArgsOut[GetGroupByNameMethod][0] != nil {
		group = t.ArgsOut[GetGroupByNameMethod][0].(*api.Group)
	}
	var err error
	if t.ArgsOut[GetGroupByNameMethod][1] != nil {
		err = t.ArgsOut[GetGroupByNameMethod][1].(error)
	}
	return group, err
}

func (t TestAPI) GetGroupList(authenticatedUser api.AuthenticatedUser, org string, pathPrefix string) ([]api.GroupIdentity, error) {
	t.ArgsIn[GetGroupListMethod][0] = authenticatedUser
	t.ArgsIn[GetGroupListMethod][1] = org
	t.ArgsIn[GetGroupListMethod][2] = pathPrefix
	var groups []api.GroupIdentity
	if t.ArgsOut[GetGroupListMethod][0] != nil {
		groups = t.ArgsOut[GetGroupListMethod][0].([]api.GroupIdentity)
	}
	var err error
	if t.ArgsOut[GetGroupListMethod][1] != nil {
		err = t.ArgsOut[GetGroupListMethod][1].(error)
	}
	return groups, err
}

func (t TestAPI) UpdateGroup(authenticatedUser api.AuthenticatedUser, org string, groupName string, newName string, newPath string) (*api.Group, error) {
	t.ArgsIn[UpdateGroupMethod][0] = authenticatedUser
	t.ArgsIn[UpdateGroupMethod][1] = org
	t.ArgsIn[UpdateGroupMethod][2] = groupName
	t.ArgsIn[UpdateGroupMethod][3] = newName
	t.ArgsIn[UpdateGroupMethod][4] = newPath
	var group *api.Group
	if t.ArgsOut[UpdateGroupMethod][0] != nil {
		group = t.ArgsOut[UpdateGroupMethod][0].(*api.Group)
	}
	var err error
	if t.ArgsOut[UpdateGroupMethod][1] != nil {
		err = t.ArgsOut[UpdateGroupMethod][1].(error)
	}
	return group, err
}

func (t TestAPI) RemoveGroup(authenticatedUser api.AuthenticatedUser, org string, name string) error {
	t.ArgsIn[RemoveGroupMethod][0] = authenticatedUser
	t.ArgsIn[RemoveGroupMethod][1] = org
	t.ArgsIn[RemoveGroupMethod][2] = name
	var err error
	if t.ArgsOut[RemoveGroupMethod][0] != nil {
		err = t.ArgsOut[RemoveGroupMethod][0].(error)
	}
	return err
}

func (t TestAPI) AddMember(authenticatedUser api.AuthenticatedUser, userID string, groupName string, org string) error {
	t.ArgsIn[AddMemberMethod][0] = authenticatedUser
	t.ArgsIn[AddMemberMethod][1] = userID
	t.ArgsIn[AddMemberMethod][2] = groupName
	t.ArgsIn[AddMemberMethod][3] = org
	var err error
	if t.ArgsOut[AddMemberMethod][0] != nil {
		err = t.ArgsOut[AddMemberMethod][0].(error)
	}
	return err
}

func (t TestAPI) RemoveMember(authenticatedUser api.AuthenticatedUser, userID string, groupName string, org string) error {
	t.ArgsIn[RemoveMemberMethod][0] = authenticatedUser
	t.ArgsIn[RemoveMemberMethod][1] = userID
	t.ArgsIn[RemoveMemberMethod][2] = groupName
	t.ArgsIn[RemoveMemberMethod][3] = org
	var err error
	if t.ArgsOut[RemoveMemberMethod][0] != nil {
		err = t.ArgsOut[RemoveMemberMethod][0].(error)
	}
	return err
}

func (t TestAPI) ListMembers(authenticatedUser api.AuthenticatedUser, org string, groupName string) ([]string, error) {
	t.ArgsIn[ListMembersMethod][0] = authenticatedUser
	t.ArgsIn[ListMembersMethod][1] = org
	t.ArgsIn[ListMembersMethod][2] = groupName
	var externalIDs []string
	if t.ArgsOut[ListMembersMethod][0] != nil {
		externalIDs = t.ArgsOut[ListMembersMethod][0].([]string)
	}
	var err error
	if t.ArgsOut[ListMembersMethod][1] != nil {
		err = t.ArgsOut[ListMembersMethod][1].(error)
	}
	return externalIDs, err
}

func (t TestAPI) AttachPolicyToGroup(authenticatedUser api.AuthenticatedUser, org string, groupName string, policyName string) error {
	t.ArgsIn[AttachPolicyToGroupMethod][0] = authenticatedUser
	t.ArgsIn[AttachPolicyToGroupMethod][1] = org
	t.ArgsIn[AttachPolicyToGroupMethod][2] = groupName
	t.ArgsIn[AttachPolicyToGroupMethod][3] = policyName
	var err error
	if t.ArgsOut[AttachPolicyToGroupMethod][0] != nil {
		err = t.ArgsOut[AttachPolicyToGroupMethod][0].(error)
	}
	return err
}

func (t TestAPI) DetachPolicyToGroup(authenticatedUser api.AuthenticatedUser, org string, groupName string, policyName string) error {
	t.ArgsIn[DetachPolicyToGroupMethod][0] = authenticatedUser
	t.ArgsIn[DetachPolicyToGroupMethod][1] = org
	t.ArgsIn[DetachPolicyToGroupMethod][2] = groupName
	t.ArgsIn[DetachPolicyToGroupMethod][3] = policyName
	var err error
	if t.ArgsOut[DetachPolicyToGroupMethod][0] != nil {
		err = t.ArgsOut[DetachPolicyToGroupMethod][0].(error)
	}
	return err
}

func (t TestAPI) ListAttachedGroupPolicies(authenticatedUser api.AuthenticatedUser, org string, groupName string) ([]api.PolicyIdentity, error) {
	t.ArgsIn[ListAttachedGroupPoliciesMethod][0] = authenticatedUser
	t.ArgsIn[ListAttachedGroupPoliciesMethod][1] = org
	t.ArgsIn[ListAttachedGroupPoliciesMethod][2] = groupName
	var policies []api.PolicyIdentity
	if t.ArgsOut[ListAttachedGroupPoliciesMethod][0] != nil {
		policies = t.ArgsOut[ListAttachedGroupPoliciesMethod][0].([]api.PolicyIdentity)
	}
	var err error
	if t.ArgsOut[ListAttachedGroupPoliciesMethod][1] != nil {
		err = t.ArgsOut[ListAttachedGroupPoliciesMethod][1].(error)
	}
	return policies, err
}

// POLICY API

func (t TestAPI) AddPolicy(authenticatedUser api.AuthenticatedUser, name string, path string, org string, statements []api.Statement) (*api.Policy, error) {
	t.ArgsIn[AddPolicyMethod][0] = authenticatedUser
	t.ArgsIn[AddPolicyMethod][1] = name
	t.ArgsIn[AddPolicyMethod][2] = path
	t.ArgsIn[AddPolicyMethod][3] = org
	t.ArgsIn[AddPolicyMethod][4] = statements
	var policy *api.Policy
	if t.ArgsOut[AddPolicyMethod][0] != nil {
		policy = t.ArgsOut[AddPolicyMethod][0].(*api.Policy)
	}
	var err error
	if t.ArgsOut[AddPolicyMethod][1] != nil {
		err = t.ArgsOut[AddPolicyMethod][1].(error)
	}
	return policy, err
}

func (t TestAPI) GetPolicyByName(authenticatedUser api.AuthenticatedUser, org string, policyName string) (*api.Policy, error) {
	t.ArgsIn[GetPolicyByNameMethod][0] = authenticatedUser
	t.ArgsIn[GetPolicyByNameMethod][1] = org
	t.ArgsIn[GetPolicyByNameMethod][2] = policyName
	var policy *api.Policy
	if t.ArgsOut[GetPolicyByNameMethod][0] != nil {
		policy = t.ArgsOut[GetPolicyByNameMethod][0].(*api.Policy)
	}
	var err error
	if t.ArgsOut[GetPolicyByNameMethod][1] != nil {
		err = t.ArgsOut[GetPolicyByNameMethod][1].(error)
	}
	return policy, err
}

func (t TestAPI) GetPolicyList(authenticatedUser api.AuthenticatedUser, org string, pathPrefix string) ([]api.PolicyIdentity, error) {
	t.ArgsIn[GetPolicyListMethod][0] = authenticatedUser
	t.ArgsIn[GetPolicyListMethod][1] = org
	t.ArgsIn[GetPolicyListMethod][2] = pathPrefix
	var policies []api.PolicyIdentity
	if t.ArgsOut[GetPolicyListMethod][0] != nil {
		policies = t.ArgsOut[GetPolicyListMethod][0].([]api.PolicyIdentity)
	}
	var err error
	if t.ArgsOut[GetPolicyListMethod][1] != nil {
		err = t.ArgsOut[GetPolicyListMethod][1].(error)
	}
	return policies, err
}

func (t TestAPI) UpdatePolicy(authenticatedUser api.AuthenticatedUser, org string, policyName string, newName string, newPath string,
	newStatements []api.Statement) (*api.Policy, error) {
	t.ArgsIn[UpdatePolicyMethod][0] = authenticatedUser
	t.ArgsIn[UpdatePolicyMethod][1] = org
	t.ArgsIn[UpdatePolicyMethod][2] = policyName
	t.ArgsIn[UpdatePolicyMethod][3] = newName
	t.ArgsIn[UpdatePolicyMethod][4] = newPath
	t.ArgsIn[UpdatePolicyMethod][5] = newStatements

	var policy *api.Policy
	if t.ArgsOut[UpdatePolicyMethod][0] != nil {
		policy = t.ArgsOut[UpdatePolicyMethod][0].(*api.Policy)
	}
	var err error
	if t.ArgsOut[UpdatePolicyMethod][1] != nil {
		err = t.ArgsOut[UpdatePolicyMethod][1].(error)
	}
	return policy, err
}

func (t TestAPI) DeletePolicy(authenticatedUser api.AuthenticatedUser, org string, name string) error {
	t.ArgsIn[DeletePolicyMethod][0] = authenticatedUser
	t.ArgsIn[DeletePolicyMethod][1] = org
	t.ArgsIn[DeletePolicyMethod][2] = name
	var err error
	if t.ArgsOut[DeletePolicyMethod][0] != nil {
		err = t.ArgsOut[DeletePolicyMethod][0].(error)
	}
	return err
}

func (t TestAPI) GetAttachedGroups(authenticatedUser api.AuthenticatedUser, org string, policyName string) ([]api.GroupIdentity, error) {
	t.ArgsIn[GetAttachedGroupsMethod][0] = authenticatedUser
	t.ArgsIn[GetAttachedGroupsMethod][1] = org
	t.ArgsIn[GetAttachedGroupsMethod][2] = policyName
	var groups []api.GroupIdentity
	if t.ArgsOut[GetAttachedGroupsMethod][0] != nil {
		groups = t.ArgsOut[GetAttachedGroupsMethod][0].([]api.GroupIdentity)
	}
	var err error
	if t.ArgsOut[GetAttachedGroupsMethod][1] != nil {
		err = t.ArgsOut[GetAttachedGroupsMethod][1].(error)
	}
	return groups, err
}

// AUTHZ API

func (t TestAPI) GetAuthorizedUsers(authenticatedUser api.AuthenticatedUser, resourceUrn string, action string, users []api.User) ([]api.User, error) {
	return nil, nil
}

func (t TestAPI) GetAuthorizedGroups(authenticatedUser api.AuthenticatedUser, resourceUrn string, action string, groups []api.Group) ([]api.Group, error) {
	return nil, nil
}

func (t TestAPI) GetAuthorizedPolicies(authenticatedUser api.AuthenticatedUser, resourceUrn string, action string, policies []api.Policy) ([]api.Policy, error) {
	return nil, nil
}

func (t TestAPI) GetAuthorizedExternalResources(authenticatedUser api.AuthenticatedUser, action string, resources []string) ([]string, error) {
	t.ArgsIn[GetAuthorizedExternalResourcesMethod][0] = authenticatedUser
	t.ArgsIn[GetAuthorizedExternalResourcesMethod][1] = action
	t.ArgsIn[GetAuthorizedExternalResourcesMethod][2] = resources
	var resourcesToReturn []string
	if t.ArgsOut[GetAuthorizedExternalResourcesMethod][0] != nil {
		resourcesToReturn = t.ArgsOut[GetAuthorizedExternalResourcesMethod][0].([]string)
	}
	var err error
	if t.ArgsOut[GetAuthorizedExternalResourcesMethod][1] != nil {
		err = t.ArgsOut[GetAuthorizedExternalResourcesMethod][1].(error)
	}
	return resourcesToReturn, err
}
