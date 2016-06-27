package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/auth"
	"github.com/tecsisa/authorizr/authorizr"
)

const (
	// USER API METHODS
	AddUserMethod             = "AddUser"
	GetUserByExternalIdMethod = "GetUserByExternalId"
	GetListUsersMethod        = "GetListUsers"
	UpdateUserMethod          = "UpdateUser"
	RemoveUserByIdMethod      = "RemoveUserById"
	GetGroupsByUserIdMethod   = "GetGroupsByUserId"

	// GROUP API METHODS
	AddGroupMethod                  = "AddGroup"
	GetGroupByNameMethod            = "GetGroupByName"
	GetListGroupsMethod             = "GetListGroups"
	UpdateGroupMethod               = "UpdateGroup"
	RemoveGroupMethod               = "RemoveGroup"
	AddMemberMethod                 = "AddMember"
	RemoveMemberMethod              = "RemoveMember"
	ListMembersMethod               = "ListMembers"
	AttachPolicyToGroupMethod       = "AttachPolicyToGroup"
	DetachPolicyToGroupMethod       = "DetachPolicyToGroup"
	ListAttachedGroupPoliciesMethod = "ListAttachedGroupPolicies"

	// POLICY API METHODS
	AddPolicyMethod               = "AddPolicy"
	GetPolicyByNameMethod         = "GetPolicyByName"
	GetListPoliciesMethod         = "GetListPolicies"
	UpdatePolicyMethod            = "UpdatePolicy"
	DeletePolicyMethod            = "DeletePolicy"
	GetPolicyAttachedGroupsMethod = "GetPolicyAttachedGroups"

	// AUTHZ API
	GetUsersAuthorizedMethod             = "GetUsersAuthorized"
	GetGroupsAuthorizedMethod            = "GetGroupsAuthorized"
	GetPoliciesAuthorizedMethod          = "GetPoliciesAuthorized"
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
	testApi.ArgsIn[GetListUsersMethod] = make([]interface{}, 2)
	testApi.ArgsIn[UpdateUserMethod] = make([]interface{}, 3)
	testApi.ArgsIn[RemoveUserByIdMethod] = make([]interface{}, 2)
	testApi.ArgsIn[GetGroupsByUserIdMethod] = make([]interface{}, 2)

	testApi.ArgsIn[AddGroupMethod] = make([]interface{}, 4)
	testApi.ArgsIn[GetGroupByNameMethod] = make([]interface{}, 3)
	testApi.ArgsIn[GetListGroupsMethod] = make([]interface{}, 3)
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
	testApi.ArgsIn[GetListPoliciesMethod] = make([]interface{}, 3)
	testApi.ArgsIn[UpdatePolicyMethod] = make([]interface{}, 6)
	testApi.ArgsIn[DeletePolicyMethod] = make([]interface{}, 3)
	testApi.ArgsIn[GetPolicyAttachedGroupsMethod] = make([]interface{}, 3)

	testApi.ArgsIn[GetUsersAuthorizedMethod] = make([]interface{}, 4)
	testApi.ArgsIn[GetGroupsAuthorizedMethod] = make([]interface{}, 4)
	testApi.ArgsIn[GetPoliciesAuthorizedMethod] = make([]interface{}, 4)
	testApi.ArgsIn[GetAuthorizedExternalResourcesMethod] = make([]interface{}, 3)

	testApi.ArgsOut[AddUserMethod] = make([]interface{}, 2)
	testApi.ArgsOut[GetUserByExternalIdMethod] = make([]interface{}, 2)
	testApi.ArgsOut[GetListUsersMethod] = make([]interface{}, 2)
	testApi.ArgsOut[UpdateUserMethod] = make([]interface{}, 2)
	testApi.ArgsOut[RemoveUserByIdMethod] = make([]interface{}, 1)
	testApi.ArgsOut[GetGroupsByUserIdMethod] = make([]interface{}, 2)

	testApi.ArgsOut[AddGroupMethod] = make([]interface{}, 2)
	testApi.ArgsOut[GetGroupByNameMethod] = make([]interface{}, 2)
	testApi.ArgsOut[GetListGroupsMethod] = make([]interface{}, 2)
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
	testApi.ArgsOut[GetListPoliciesMethod] = make([]interface{}, 2)
	testApi.ArgsOut[UpdatePolicyMethod] = make([]interface{}, 2)
	testApi.ArgsOut[DeletePolicyMethod] = make([]interface{}, 1)
	testApi.ArgsOut[GetPolicyAttachedGroupsMethod] = make([]interface{}, 2)

	testApi.ArgsOut[GetUsersAuthorizedMethod] = make([]interface{}, 2)
	testApi.ArgsOut[GetGroupsAuthorizedMethod] = make([]interface{}, 2)
	testApi.ArgsOut[GetPoliciesAuthorizedMethod] = make([]interface{}, 2)
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

func (t TestAPI) GetListUsers(authenticatedUser api.AuthenticatedUser, pathPrefix string) ([]string, error) {
	t.ArgsIn[GetListUsersMethod][0] = authenticatedUser
	t.ArgsIn[GetListUsersMethod][1] = pathPrefix
	var externalIDs []string
	if t.ArgsOut[GetListUsersMethod][0] != nil {
		externalIDs = t.ArgsOut[GetListUsersMethod][0].([]string)
	}
	var err error
	if t.ArgsOut[GetListUsersMethod][1] != nil {
		err = t.ArgsOut[GetListUsersMethod][1].(error)
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
	return nil, nil
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

func (t TestAPI) GetListGroups(authenticatedUser api.AuthenticatedUser, org string, pathPrefix string) ([]api.GroupIdentity, error) {
	return nil, nil
}

func (t TestAPI) UpdateGroup(authenticatedUser api.AuthenticatedUser, org string, groupName string, newName string, newPath string) (*api.Group, error) {
	return nil, nil
}

func (t TestAPI) RemoveGroup(authenticatedUser api.AuthenticatedUser, org string, name string) error {
	return nil
}

func (t TestAPI) AddMember(authenticatedUser api.AuthenticatedUser, userID string, groupName string, org string) error {
	return nil
}

func (t TestAPI) RemoveMember(authenticatedUser api.AuthenticatedUser, userID string, groupName string, org string) error {
	return nil
}

func (t TestAPI) ListMembers(authenticatedUser api.AuthenticatedUser, org string, groupName string) ([]string, error) {
	return nil, nil
}

func (t TestAPI) AttachPolicyToGroup(authenticatedUser api.AuthenticatedUser, org string, groupName string, policyName string) error {
	return nil
}

func (t TestAPI) DetachPolicyToGroup(authenticatedUser api.AuthenticatedUser, org string, groupName string, policyName string) error {
	return nil
}

func (t TestAPI) ListAttachedGroupPolicies(authenticatedUser api.AuthenticatedUser, org string, groupName string) ([]api.PolicyIdentity, error) {
	return nil, nil
}

// POLICY API

func (t TestAPI) AddPolicy(authenticatedUser api.AuthenticatedUser, name string, path string, org string, statements []api.Statement) (*api.Policy, error) {
	return nil, nil
}

func (t TestAPI) GetPolicyByName(authenticatedUser api.AuthenticatedUser, org string, policyName string) (*api.Policy, error) {
	return nil, nil
}

func (t TestAPI) GetListPolicies(authenticatedUser api.AuthenticatedUser, org string, pathPrefix string) ([]api.PolicyIdentity, error) {
	return nil, nil
}

func (t TestAPI) UpdatePolicy(authenticatedUser api.AuthenticatedUser, org string, policyName string, newName string, newPath string,
	newStatements []api.Statement) (*api.Policy, error) {
	return nil, nil
}

func (t TestAPI) DeletePolicy(authenticatedUser api.AuthenticatedUser, org string, name string) error {
	return nil
}

func (t TestAPI) GetPolicyAttachedGroups(authenticatedUser api.AuthenticatedUser, org string, policyName string) ([]api.GroupIdentity, error) {
	return nil, nil
}

// AUTHZ API

func (t TestAPI) GetUsersAuthorized(authenticatedUser api.AuthenticatedUser, resourceUrn string, action string, users []api.User) ([]api.User, error) {
	return nil, nil
}

func (t TestAPI) GetGroupsAuthorized(authenticatedUser api.AuthenticatedUser, resourceUrn string, action string, groups []api.Group) ([]api.Group, error) {
	return nil, nil
}

func (t TestAPI) GetPoliciesAuthorized(authenticatedUser api.AuthenticatedUser, resourceUrn string, action string, policies []api.Policy) ([]api.Policy, error) {
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
