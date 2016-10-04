package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/auth"
	"github.com/Tecsisa/foulkon/foulkon"
)

const (
	// USER API METHODS
	AddUserMethod             = "AddUser"
	GetUserByExternalIdMethod = "GetUserByExternalId"
	ListUsersMethod           = "ListUsers"
	UpdateUserMethod          = "UpdateUser"
	RemoveUserMethod          = "RemoveUser"
	ListGroupsByUserMethod    = "ListGroupsByUser"

	// GROUP API METHODS
	AddGroupMethod                  = "AddGroup"
	GetGroupByNameMethod            = "GetGroupByName"
	ListGroupsMethod                = "ListGroups"
	UpdateGroupMethod               = "UpdateGroup"
	RemoveGroupMethod               = "RemoveGroup"
	AddMemberMethod                 = "AddMember"
	RemoveMemberMethod              = "RemoveMember"
	ListMembersMethod               = "ListMembers"
	AttachPolicyToGroupMethod       = "AttachPolicyToGroup"
	DetachPolicyToGroupMethod       = "DetachPolicyToGroup"
	ListAttachedGroupPoliciesMethod = "ListAttachedGroupPolicies"

	// POLICY API METHODS
	AddPolicyMethod          = "AddPolicy"
	GetPolicyByNameMethod    = "GetPolicyByName"
	ListPoliciesMethod       = "ListPolicies"
	UpdatePolicyMethod       = "UpdatePolicy"
	RemovePolicyMethod       = "RemovePolicy"
	ListAttachedGroupsMethod = "ListAttachedGroups"

	// AUTHZ API
	GetAuthorizedUsersMethod             = "GetAuthorizedUsers"
	GetAuthorizedGroupsMethod            = "GetAuthorizedGroups"
	GetAuthorizedPoliciesMethod          = "GetAuthorizedPolicies"
	GetAuthorizedExternalResourcesMethod = "GetAuthorizedExternalResources"
)

// Test server used to test handlers
var server *httptest.Server
var proxy *httptest.Server
var testApi *TestAPI
var authConnector *TestConnector
var testFilter = &api.Filter{
	PathPrefix: "",
	Org:        "",
	GroupName:  "",
	PolicyName: "",
	ExternalID: "",
	Offset:     0,
	Limit:      0,
}

// Test API that implements all api manager interfaces
type TestAPI struct {
	ArgsIn       map[string][]interface{}
	ArgsOut      map[string][]interface{}
	SpecialFuncs map[string]interface{}
}

// Aux connector
type TestConnector struct {
	userID          string
	unauthenticated bool
}

func (tc TestConnector) Authenticate(h http.Handler) http.Handler {
	if authConnector.unauthenticated {
		// Reset value
		authConnector.unauthenticated = false
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		})
	}
	return h
}

func (tc TestConnector) RetrieveUserID(r http.Request) string {
	return tc.userID
}

// Main Test that executes at first time and create all necessary data to work
func TestMain(m *testing.M) {
	// Create logger
	logger := &log.Logger{
		Out:       bytes.NewBuffer([]byte{}),
		Formatter: &log.TextFormatter{},
		Hooks:     make(log.LevelHooks),
		Level:     log.DebugLevel,
	}

	testApi = makeTestApi()

	// Instantiate Auth Connector
	authConnector = &TestConnector{
		userID:          "userID",
		unauthenticated: false,
	}

	adminUser := "admin"
	adminPassword := "admin"

	// Create authenticator
	authenticator := auth.NewAuthenticator(authConnector, adminUser, adminPassword)

	// Return created core
	worker := &foulkon.Worker{
		Logger:        logger,
		Authenticator: authenticator,
		UserApi:       testApi,
		GroupApi:      testApi,
		PolicyApi:     testApi,
		AuthzApi:      testApi,
	}

	server = httptest.NewServer(WorkerHandlerRouter(worker))

	proxyCore := &foulkon.Proxy{
		Logger:     logger,
		WorkerHost: server.URL,
		APIResources: []foulkon.APIResource{
			{
				Id:     "resource1",
				Host:   server.URL,
				Url:    USER_ID_URL,
				Method: "GET",
				Urn:    "urn:ews:example:instance1:resource/{userid}",
				Action: "example:user",
			},
			{
				Id:     "hostUnreachable",
				Host:   "fail",
				Url:    "/fail",
				Method: "GET",
				Urn:    "urn:ews:example:instance1:resource/fail",
				Action: "example:fail",
			},
			{
				Id:     "invalidHost",
				Host:   "%&",
				Url:    "/invalid",
				Method: "GET",
				Urn:    "urn:ews:example:instance1:resource/invalid",
				Action: "example:invalid",
			},
			{
				Id:     "invalidUrn",
				Host:   server.URL,
				Url:    "/invalidUrn",
				Method: "GET",
				Urn:    "%&",
				Action: "example:invalid",
			},
			{
				Id:     "urnPrefix",
				Host:   server.URL,
				Url:    "/urnPrefix",
				Method: "GET",
				Urn:    "urn:*",
				Action: "&%",
			},
			{
				Id:     "invalidAction",
				Host:   server.URL,
				Url:    "/invalidAction",
				Method: "GET",
				Urn:    "urn:ews:example:instance1:resource/user",
				Action: "&%",
			},
		},
	}

	proxy = httptest.NewServer(ProxyHandlerRouter(proxyCore))

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
	testApi.ArgsIn[ListUsersMethod] = make([]interface{}, 2)
	testApi.ArgsIn[UpdateUserMethod] = make([]interface{}, 3)
	testApi.ArgsIn[RemoveUserMethod] = make([]interface{}, 2)
	testApi.ArgsIn[ListGroupsByUserMethod] = make([]interface{}, 2)

	testApi.ArgsIn[AddGroupMethod] = make([]interface{}, 4)
	testApi.ArgsIn[GetGroupByNameMethod] = make([]interface{}, 3)
	testApi.ArgsIn[ListGroupsMethod] = make([]interface{}, 2)
	testApi.ArgsIn[UpdateGroupMethod] = make([]interface{}, 5)
	testApi.ArgsIn[RemoveGroupMethod] = make([]interface{}, 3)
	testApi.ArgsIn[AddMemberMethod] = make([]interface{}, 4)
	testApi.ArgsIn[RemoveMemberMethod] = make([]interface{}, 4)
	testApi.ArgsIn[ListMembersMethod] = make([]interface{}, 2)
	testApi.ArgsIn[AttachPolicyToGroupMethod] = make([]interface{}, 4)
	testApi.ArgsIn[DetachPolicyToGroupMethod] = make([]interface{}, 4)
	testApi.ArgsIn[ListAttachedGroupPoliciesMethod] = make([]interface{}, 2)

	testApi.ArgsIn[AddPolicyMethod] = make([]interface{}, 5)
	testApi.ArgsIn[GetPolicyByNameMethod] = make([]interface{}, 3)
	testApi.ArgsIn[ListPoliciesMethod] = make([]interface{}, 2)
	testApi.ArgsIn[UpdatePolicyMethod] = make([]interface{}, 6)
	testApi.ArgsIn[RemovePolicyMethod] = make([]interface{}, 3)
	testApi.ArgsIn[ListAttachedGroupsMethod] = make([]interface{}, 2)

	testApi.ArgsIn[GetAuthorizedUsersMethod] = make([]interface{}, 4)
	testApi.ArgsIn[GetAuthorizedGroupsMethod] = make([]interface{}, 4)
	testApi.ArgsIn[GetAuthorizedPoliciesMethod] = make([]interface{}, 4)
	testApi.ArgsIn[GetAuthorizedExternalResourcesMethod] = make([]interface{}, 3)

	testApi.ArgsOut[AddUserMethod] = make([]interface{}, 2)
	testApi.ArgsOut[GetUserByExternalIdMethod] = make([]interface{}, 2)
	testApi.ArgsOut[ListUsersMethod] = make([]interface{}, 3)
	testApi.ArgsOut[UpdateUserMethod] = make([]interface{}, 2)
	testApi.ArgsOut[RemoveUserMethod] = make([]interface{}, 1)
	testApi.ArgsOut[ListGroupsByUserMethod] = make([]interface{}, 3)

	testApi.ArgsOut[AddGroupMethod] = make([]interface{}, 2)
	testApi.ArgsOut[GetGroupByNameMethod] = make([]interface{}, 2)
	testApi.ArgsOut[ListGroupsMethod] = make([]interface{}, 3)
	testApi.ArgsOut[UpdateGroupMethod] = make([]interface{}, 2)
	testApi.ArgsOut[RemoveGroupMethod] = make([]interface{}, 1)
	testApi.ArgsOut[AddMemberMethod] = make([]interface{}, 1)
	testApi.ArgsOut[RemoveMemberMethod] = make([]interface{}, 1)
	testApi.ArgsOut[ListMembersMethod] = make([]interface{}, 3)
	testApi.ArgsOut[AttachPolicyToGroupMethod] = make([]interface{}, 1)
	testApi.ArgsOut[DetachPolicyToGroupMethod] = make([]interface{}, 1)
	testApi.ArgsOut[ListAttachedGroupPoliciesMethod] = make([]interface{}, 3)

	testApi.ArgsOut[AddPolicyMethod] = make([]interface{}, 2)
	testApi.ArgsOut[GetPolicyByNameMethod] = make([]interface{}, 2)
	testApi.ArgsOut[ListPoliciesMethod] = make([]interface{}, 3)
	testApi.ArgsOut[UpdatePolicyMethod] = make([]interface{}, 2)
	testApi.ArgsOut[RemovePolicyMethod] = make([]interface{}, 1)
	testApi.ArgsOut[ListAttachedGroupsMethod] = make([]interface{}, 3)

	testApi.ArgsOut[GetAuthorizedUsersMethod] = make([]interface{}, 2)
	testApi.ArgsOut[GetAuthorizedGroupsMethod] = make([]interface{}, 2)
	testApi.ArgsOut[GetAuthorizedPoliciesMethod] = make([]interface{}, 2)
	testApi.ArgsOut[GetAuthorizedExternalResourcesMethod] = make([]interface{}, 2)

	return testApi
}

// USER API

func (t TestAPI) AddUser(authenticatedUser api.RequestInfo, externalID string, path string) (*api.User, error) {
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

func (t TestAPI) GetUserByExternalID(authenticatedUser api.RequestInfo, id string) (*api.User, error) {
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

func (t TestAPI) ListUsers(authenticatedUser api.RequestInfo, filter *api.Filter) ([]string, int, error) {
	t.ArgsIn[ListUsersMethod][0] = authenticatedUser
	t.ArgsIn[ListUsersMethod][1] = filter
	var externalIDs []string
	var total int
	if t.ArgsOut[ListUsersMethod][0] != nil {
		externalIDs = t.ArgsOut[ListUsersMethod][0].([]string)
	}
	if t.ArgsOut[ListUsersMethod][1] != nil {
		total = t.ArgsOut[ListUsersMethod][1].(int)
	}
	var err error
	if t.ArgsOut[ListUsersMethod][2] != nil {
		err = t.ArgsOut[ListUsersMethod][2].(error)
	}
	return externalIDs, total, err
}

func (t TestAPI) UpdateUser(authenticatedUser api.RequestInfo, externalID string, newPath string) (*api.User, error) {
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

func (t TestAPI) RemoveUser(authenticatedUser api.RequestInfo, id string) error {
	t.ArgsIn[RemoveUserMethod][0] = authenticatedUser
	t.ArgsIn[RemoveUserMethod][1] = id
	var err error
	if t.ArgsOut[RemoveUserMethod][0] != nil {
		err = t.ArgsOut[RemoveUserMethod][0].(error)
	}
	return err
}

func (t TestAPI) ListGroupsByUser(authenticatedUser api.RequestInfo, filter *api.Filter) ([]api.UserGroups, int, error) {
	t.ArgsIn[ListGroupsByUserMethod][0] = authenticatedUser
	t.ArgsIn[ListGroupsByUserMethod][1] = filter
	var groups []api.UserGroups
	var total int
	if t.ArgsOut[ListGroupsByUserMethod][1] != nil {
		total = t.ArgsOut[ListGroupsByUserMethod][1].(int)
	}
	if t.ArgsOut[ListGroupsByUserMethod][0] != nil {
		groups = t.ArgsOut[ListGroupsByUserMethod][0].([]api.UserGroups)

	}
	var err error
	if t.ArgsOut[ListGroupsByUserMethod][2] != nil {
		err = t.ArgsOut[ListGroupsByUserMethod][2].(error)
	}
	return groups, total, err
}

// GROUP API

func (t TestAPI) AddGroup(authenticatedUser api.RequestInfo, org string, name string, path string) (*api.Group, error) {
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

func (t TestAPI) GetGroupByName(authenticatedUser api.RequestInfo, org string, name string) (*api.Group, error) {
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

func (t TestAPI) ListGroups(authenticatedUser api.RequestInfo, filter *api.Filter) ([]api.GroupIdentity, int, error) {
	t.ArgsIn[ListGroupsMethod][0] = authenticatedUser
	t.ArgsIn[ListGroupsMethod][1] = filter

	var groups []api.GroupIdentity
	var total int
	if t.ArgsOut[ListGroupsMethod][1] != nil {
		total = t.ArgsOut[ListGroupsMethod][1].(int)
	}
	if t.ArgsOut[ListGroupsMethod][0] != nil {
		groups = t.ArgsOut[ListGroupsMethod][0].([]api.GroupIdentity)
	}
	var err error
	if t.ArgsOut[ListGroupsMethod][2] != nil {
		err = t.ArgsOut[ListGroupsMethod][2].(error)
	}
	return groups, total, err
}

func (t TestAPI) UpdateGroup(authenticatedUser api.RequestInfo, org string, groupName string, newName string, newPath string) (*api.Group, error) {
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

func (t TestAPI) RemoveGroup(authenticatedUser api.RequestInfo, org string, name string) error {
	t.ArgsIn[RemoveGroupMethod][0] = authenticatedUser
	t.ArgsIn[RemoveGroupMethod][1] = org
	t.ArgsIn[RemoveGroupMethod][2] = name
	var err error
	if t.ArgsOut[RemoveGroupMethod][0] != nil {
		err = t.ArgsOut[RemoveGroupMethod][0].(error)
	}
	return err
}

func (t TestAPI) AddMember(authenticatedUser api.RequestInfo, userID string, groupName string, org string) error {
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

func (t TestAPI) RemoveMember(authenticatedUser api.RequestInfo, userID string, groupName string, org string) error {
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

func (t TestAPI) ListMembers(authenticatedUser api.RequestInfo, filter *api.Filter) ([]api.GroupMembers, int, error) {
	t.ArgsIn[ListMembersMethod][0] = authenticatedUser
	t.ArgsIn[ListMembersMethod][1] = filter

	var externalIDs []api.GroupMembers
	var total int
	if t.ArgsOut[ListMembersMethod][1] != nil {
		total = t.ArgsOut[ListMembersMethod][1].(int)
	}
	if t.ArgsOut[ListMembersMethod][0] != nil {
		externalIDs = t.ArgsOut[ListMembersMethod][0].([]api.GroupMembers)
	}
	var err error
	if t.ArgsOut[ListMembersMethod][2] != nil {
		err = t.ArgsOut[ListMembersMethod][2].(error)
	}
	return externalIDs, total, err
}

func (t TestAPI) AttachPolicyToGroup(authenticatedUser api.RequestInfo, org string, groupName string, policyName string) error {
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

func (t TestAPI) DetachPolicyToGroup(authenticatedUser api.RequestInfo, org string, groupName string, policyName string) error {
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

func (t TestAPI) ListAttachedGroupPolicies(authenticatedUser api.RequestInfo, filter *api.Filter) ([]api.GroupPolicies, int, error) {
	t.ArgsIn[ListAttachedGroupPoliciesMethod][0] = authenticatedUser
	t.ArgsIn[ListAttachedGroupPoliciesMethod][1] = filter

	var policies []api.GroupPolicies
	var total int
	if t.ArgsOut[ListAttachedGroupPoliciesMethod][1] != nil {
		total = t.ArgsOut[ListAttachedGroupPoliciesMethod][1].(int)
	}
	if t.ArgsOut[ListAttachedGroupPoliciesMethod][0] != nil {
		policies = t.ArgsOut[ListAttachedGroupPoliciesMethod][0].([]api.GroupPolicies)
	}
	var err error
	if t.ArgsOut[ListAttachedGroupPoliciesMethod][2] != nil {
		err = t.ArgsOut[ListAttachedGroupPoliciesMethod][2].(error)
	}
	return policies, total, err
}

// POLICY API

func (t TestAPI) AddPolicy(authenticatedUser api.RequestInfo, name string, path string, org string, statements []api.Statement) (*api.Policy, error) {
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

func (t TestAPI) GetPolicyByName(authenticatedUser api.RequestInfo, org string, policyName string) (*api.Policy, error) {
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

func (t TestAPI) ListPolicies(authenticatedUser api.RequestInfo, filter *api.Filter) ([]api.PolicyIdentity, int, error) {
	t.ArgsIn[ListPoliciesMethod][0] = authenticatedUser
	t.ArgsIn[ListPoliciesMethod][1] = filter

	var policies []api.PolicyIdentity
	var total int
	if t.ArgsOut[ListPoliciesMethod][1] != nil {
		total = t.ArgsOut[ListPoliciesMethod][1].(int)
	}
	if t.ArgsOut[ListPoliciesMethod][0] != nil {
		policies = t.ArgsOut[ListPoliciesMethod][0].([]api.PolicyIdentity)
	}
	var err error
	if t.ArgsOut[ListPoliciesMethod][2] != nil {
		err = t.ArgsOut[ListPoliciesMethod][2].(error)
	}
	return policies, total, err
}

func (t TestAPI) UpdatePolicy(authenticatedUser api.RequestInfo, org string, policyName string, newName string, newPath string,
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

func (t TestAPI) RemovePolicy(authenticatedUser api.RequestInfo, org string, name string) error {
	t.ArgsIn[RemovePolicyMethod][0] = authenticatedUser
	t.ArgsIn[RemovePolicyMethod][1] = org
	t.ArgsIn[RemovePolicyMethod][2] = name
	var err error
	if t.ArgsOut[RemovePolicyMethod][0] != nil {
		err = t.ArgsOut[RemovePolicyMethod][0].(error)
	}
	return err
}

func (t TestAPI) ListAttachedGroups(authenticatedUser api.RequestInfo, filter *api.Filter) ([]api.PolicyGroups, int, error) {
	t.ArgsIn[ListAttachedGroupsMethod][0] = authenticatedUser
	t.ArgsIn[ListAttachedGroupsMethod][1] = filter

	var groups []api.PolicyGroups
	var total int
	if t.ArgsOut[ListAttachedGroupsMethod][1] != nil {
		total = t.ArgsOut[ListAttachedGroupsMethod][1].(int)
	}
	if t.ArgsOut[ListAttachedGroupsMethod][0] != nil {
		groups = t.ArgsOut[ListAttachedGroupsMethod][0].([]api.PolicyGroups)
	}
	var err error
	if t.ArgsOut[ListAttachedGroupsMethod][2] != nil {
		err = t.ArgsOut[ListAttachedGroupsMethod][2].(error)
	}
	return groups, total, err
}

// AUTHZ API

func (t TestAPI) GetAuthorizedUsers(authenticatedUser api.RequestInfo, resourceUrn string, action string, users []api.User) ([]api.User, error) {
	return nil, nil
}

func (t TestAPI) GetAuthorizedGroups(authenticatedUser api.RequestInfo, resourceUrn string, action string, groups []api.Group) ([]api.Group, error) {
	return nil, nil
}

func (t TestAPI) GetAuthorizedPolicies(authenticatedUser api.RequestInfo, resourceUrn string, action string, policies []api.Policy) ([]api.Policy, error) {
	return nil, nil
}

func (t TestAPI) GetAuthorizedExternalResources(authenticatedUser api.RequestInfo, action string, resources []string) ([]string, error) {
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

// Private helper methods

func addQueryParams(filter *api.Filter, r *http.Request) {
	if filter != nil {
		q := r.URL.Query()
		if filter.PathPrefix != "" {
			q.Add("PathPrefix", filter.PathPrefix)
		}
		q.Add("Offset", fmt.Sprintf("%v", filter.Offset))
		q.Add("Limit", fmt.Sprintf("%v", filter.Limit))
		r.URL.RawQuery = q.Encode()
	}
}
