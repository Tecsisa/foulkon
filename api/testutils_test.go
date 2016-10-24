package api

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/kylelemons/godebug/pretty"
)

const (
	GetUserByExternalIDMethod = "GetUserByExternalID"
	AddUserMethod             = "AddUser"
	UpdateUserMethod          = "UpdateUser"
	GetUsersFilteredMethod    = "GetUsersFiltered"
	GetGroupsByUserIDMethod   = "GetGroupsByUserID"
	RemoveUserMethod          = "RemoveUser"
	GetGroupByNameMethod      = "GetGroupByName"
	IsMemberOfGroupMethod     = "IsMemberOfGroup"
	GetGroupMembersMethod     = "GetGroupMembers"
	IsAttachedToGroupMethod   = "IsAttachedToGroup"
	GetAttachedPoliciesMethod = "GetAttachedPolicies"
	GetGroupsFilteredMethod   = "GetGroupsFiltered"
	RemoveGroupMethod         = "RemoveGroup"
	AddGroupMethod            = "AddGroup"
	AddMemberMethod           = "AddMember"
	RemoveMemberMethod        = "RemoveMember"
	UpdateGroupMethod         = "UpdateGroup"
	AttachPolicyMethod        = "AttachPolicy"
	DetachPolicyMethod        = "DetachPolicy"
	GetPolicyByNameMethod     = "GetPolicyByName"
	AddPolicyMethod           = "AddPolicy"
	UpdatePolicyMethod        = "UpdatePolicy"
	RemovePolicyMethod        = "RemovePolicy"
	GetPoliciesFilteredMethod = "GetPoliciesFiltered"
	GetAttachedGroupsMethod   = "GetAttachedGroups"
	OrderByValidColumnsMethod = "OrderByValidColumns"
	GetProxyResourcesMethod   = "GetProxyResources"
)

// TestRepo that implements all repo manager interfaces
type TestRepo struct {
	ArgsIn       map[string][]interface{}
	ArgsOut      map[string][]interface{}
	SpecialFuncs map[string]interface{}
}

type TestUserGroupRelation struct {
	User     *User
	Group    *Group
	CreateAt time.Time
}

type TestPolicyGroupRelation struct {
	Group    *Group
	Policy   *Policy
	CreateAt time.Time
}

var testFilter = Filter{
	PathPrefix: "",
	Org:        "",
	GroupName:  "",
	PolicyName: "",
	ExternalID: "",
	Offset:     0,
	Limit:      20,
}

// func that initializes the TestRepo
func makeTestRepo() *TestRepo {
	testRepo := &TestRepo{
		ArgsIn:       make(map[string][]interface{}),
		ArgsOut:      make(map[string][]interface{}),
		SpecialFuncs: make(map[string]interface{}),
	}
	testRepo.ArgsIn[GetUserByExternalIDMethod] = make([]interface{}, 1)
	testRepo.ArgsIn[AddUserMethod] = make([]interface{}, 1)
	testRepo.ArgsIn[UpdateUserMethod] = make([]interface{}, 3)
	testRepo.ArgsIn[GetUsersFilteredMethod] = make([]interface{}, 1)
	testRepo.ArgsIn[GetGroupsByUserIDMethod] = make([]interface{}, 2)
	testRepo.ArgsIn[RemoveUserMethod] = make([]interface{}, 1)
	testRepo.ArgsIn[GetGroupByNameMethod] = make([]interface{}, 2)
	testRepo.ArgsIn[IsMemberOfGroupMethod] = make([]interface{}, 2)
	testRepo.ArgsIn[GetGroupMembersMethod] = make([]interface{}, 2)
	testRepo.ArgsIn[IsAttachedToGroupMethod] = make([]interface{}, 2)
	testRepo.ArgsIn[GetAttachedPoliciesMethod] = make([]interface{}, 2)
	testRepo.ArgsIn[GetGroupsFilteredMethod] = make([]interface{}, 1)
	testRepo.ArgsIn[RemoveGroupMethod] = make([]interface{}, 1)
	testRepo.ArgsIn[AddGroupMethod] = make([]interface{}, 1)
	testRepo.ArgsIn[AddMemberMethod] = make([]interface{}, 2)
	testRepo.ArgsIn[RemoveMemberMethod] = make([]interface{}, 2)
	testRepo.ArgsIn[UpdateGroupMethod] = make([]interface{}, 4)
	testRepo.ArgsIn[AttachPolicyMethod] = make([]interface{}, 2)
	testRepo.ArgsIn[DetachPolicyMethod] = make([]interface{}, 2)
	testRepo.ArgsIn[GetPolicyByNameMethod] = make([]interface{}, 2)
	testRepo.ArgsIn[AddPolicyMethod] = make([]interface{}, 1)
	testRepo.ArgsIn[UpdatePolicyMethod] = make([]interface{}, 5)
	testRepo.ArgsIn[RemovePolicyMethod] = make([]interface{}, 1)
	testRepo.ArgsIn[GetPoliciesFilteredMethod] = make([]interface{}, 1)
	testRepo.ArgsIn[GetAttachedGroupsMethod] = make([]interface{}, 2)
	testRepo.ArgsIn[OrderByValidColumnsMethod] = make([]interface{}, 1)

	testRepo.ArgsOut[GetUserByExternalIDMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[AddUserMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[UpdateUserMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[GetUsersFilteredMethod] = make([]interface{}, 3)
	testRepo.ArgsOut[GetGroupsByUserIDMethod] = make([]interface{}, 3)
	testRepo.ArgsOut[RemoveUserMethod] = make([]interface{}, 1)
	testRepo.ArgsOut[GetGroupByNameMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[IsMemberOfGroupMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[GetGroupMembersMethod] = make([]interface{}, 3)
	testRepo.ArgsOut[IsAttachedToGroupMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[GetAttachedPoliciesMethod] = make([]interface{}, 3)
	testRepo.ArgsOut[GetGroupsFilteredMethod] = make([]interface{}, 3)
	testRepo.ArgsOut[RemoveGroupMethod] = make([]interface{}, 1)
	testRepo.ArgsOut[AddGroupMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[AddMemberMethod] = make([]interface{}, 1)
	testRepo.ArgsOut[RemoveMemberMethod] = make([]interface{}, 1)
	testRepo.ArgsOut[UpdateGroupMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[AttachPolicyMethod] = make([]interface{}, 1)
	testRepo.ArgsOut[DetachPolicyMethod] = make([]interface{}, 1)
	testRepo.ArgsOut[GetPolicyByNameMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[AddPolicyMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[UpdatePolicyMethod] = make([]interface{}, 2)
	testRepo.ArgsOut[RemovePolicyMethod] = make([]interface{}, 1)
	testRepo.ArgsOut[GetPoliciesFilteredMethod] = make([]interface{}, 3)
	testRepo.ArgsOut[GetAttachedGroupsMethod] = make([]interface{}, 3)
	testRepo.ArgsOut[OrderByValidColumnsMethod] = make([]interface{}, 1)

	testRepo.ArgsOut[GetProxyResourcesMethod] = make([]interface{}, 2)

	return testRepo
}

func makeTestAPI(testRepo *TestRepo) *WorkerAPI {
	api := &WorkerAPI{
		UserRepo:   testRepo,
		GroupRepo:  testRepo,
		PolicyRepo: testRepo,
	}
	Log = &log.Logger{
		Out:       bytes.NewBuffer([]byte{}),
		Formatter: &log.TextFormatter{},
		Hooks:     make(log.LevelHooks),
		Level:     log.DebugLevel,
	}
	return api
}

func makeProxyTestAPI(testRepo *TestRepo) *ProxyAPI {
	api := &ProxyAPI{
		ProxyRepo: testRepo,
	}
	Log = &log.Logger{
		Out:       bytes.NewBuffer([]byte{}),
		Formatter: &log.TextFormatter{},
		Hooks:     make(log.LevelHooks),
		Level:     log.DebugLevel,
	}
	return api
}

//////////////////////
// UserGroupRelation
//////////////////////

func (t TestUserGroupRelation) GetUser() *User {
	return t.User
}

func (t TestUserGroupRelation) GetGroup() *Group {
	return t.Group
}

func (t TestUserGroupRelation) GetDate() time.Time {
	return t.CreateAt
}

///////////////////////
// PolicyGroupRelation
///////////////////////

func (t TestPolicyGroupRelation) GetPolicy() *Policy {
	return t.Policy
}

func (t TestPolicyGroupRelation) GetGroup() *Group {
	return t.Group
}

func (t TestPolicyGroupRelation) GetDate() time.Time {
	return t.CreateAt
}

//////////////////
// User repo
//////////////////

func (t TestRepo) GetUserByExternalID(id string) (*User, error) {
	t.ArgsIn[GetUserByExternalIDMethod][0] = id
	if specialFunc, ok := t.SpecialFuncs[GetUserByExternalIDMethod].(func(id string) (*User, error)); ok && specialFunc != nil {
		return specialFunc(id)
	}
	var user *User
	if t.ArgsOut[GetUserByExternalIDMethod][0] != nil {
		user = t.ArgsOut[GetUserByExternalIDMethod][0].(*User)
	}
	var err error
	if t.ArgsOut[GetUserByExternalIDMethod][1] != nil {
		err = t.ArgsOut[GetUserByExternalIDMethod][1].(error)
	}
	return user, err
}

func (t TestRepo) AddUser(user User) (*User, error) {
	t.ArgsIn[AddUserMethod][0] = user
	var created *User
	if t.ArgsOut[AddUserMethod][0] != nil {
		created = t.ArgsOut[AddUserMethod][0].(*User)
	}
	var err error
	if t.ArgsOut[AddUserMethod][1] != nil {
		err = t.ArgsOut[AddUserMethod][1].(error)
	}
	return created, err
}

func (t TestRepo) UpdateUser(user User) (*User, error) {
	t.ArgsIn[UpdateUserMethod][0] = user
	var updated *User
	if t.ArgsOut[UpdateUserMethod][0] != nil {
		updated = t.ArgsOut[UpdateUserMethod][0].(*User)
	}
	var err error
	if t.ArgsOut[UpdateUserMethod][1] != nil {
		err = t.ArgsOut[UpdateUserMethod][1].(error)
	}
	return updated, err
}

func (t TestRepo) GetUsersFiltered(filter *Filter) ([]User, int, error) {
	t.ArgsIn[GetUsersFilteredMethod][0] = filter.PathPrefix
	var users []User
	if t.ArgsOut[GetUsersFilteredMethod][0] != nil {
		users = t.ArgsOut[GetUsersFilteredMethod][0].([]User)
	}

	var total int
	if t.ArgsOut[GetUsersFilteredMethod][1] != nil {
		total = t.ArgsOut[GetUsersFilteredMethod][1].(int)
	}
	var err error
	if t.ArgsOut[GetUsersFilteredMethod][2] != nil {
		err = t.ArgsOut[GetUsersFilteredMethod][2].(error)
	}
	return users, total, err
}

func (t TestRepo) GetGroupsByUserID(id string, filter *Filter) ([]UserGroupRelation, int, error) {
	t.ArgsIn[GetGroupsByUserIDMethod][0] = id
	var groups []UserGroupRelation
	if t.ArgsOut[GetGroupsByUserIDMethod][0] != nil {
		testGroups := t.ArgsOut[GetGroupsByUserIDMethod][0].([]TestUserGroupRelation)
		for _, v := range testGroups {
			groups = append(groups, v)
		}
	}

	var total int
	if t.ArgsOut[GetGroupsByUserIDMethod][1] != nil {
		total = t.ArgsOut[GetGroupsByUserIDMethod][1].(int)
	}
	var err error
	if t.ArgsOut[GetGroupsByUserIDMethod][2] != nil {
		err = t.ArgsOut[GetGroupsByUserIDMethod][2].(error)
	}
	return groups, total, err
}

func (t TestRepo) RemoveUser(id string) error {
	t.ArgsIn[RemoveUserMethod][0] = id
	var err error
	if t.ArgsOut[RemoveUserMethod][0] != nil {
		err = t.ArgsOut[RemoveUserMethod][0].(error)
	}
	return err
}

//////////////////
// Group repo
//////////////////

func (t TestRepo) GetGroupByName(org string, name string) (*Group, error) {
	t.ArgsIn[GetGroupByNameMethod][0] = org
	t.ArgsIn[GetGroupByNameMethod][1] = name
	if specialFunc, ok := t.SpecialFuncs[GetGroupByNameMethod].(func(org string, name string) (*Group, error)); ok && specialFunc != nil {
		return specialFunc(org, name)
	}
	var group *Group
	if t.ArgsOut[GetGroupByNameMethod][0] != nil {
		group = t.ArgsOut[GetGroupByNameMethod][0].(*Group)
	}
	var err error
	if t.ArgsOut[GetGroupByNameMethod][1] != nil {
		err = t.ArgsOut[GetGroupByNameMethod][1].(error)
	}
	return group, err
}

func (t TestRepo) IsMemberOfGroup(userID string, groupID string) (bool, error) {
	t.ArgsIn[IsMemberOfGroupMethod][0] = userID
	t.ArgsIn[IsMemberOfGroupMethod][1] = groupID
	var isMember bool
	if t.ArgsOut[IsMemberOfGroupMethod][0] != nil {
		isMember = t.ArgsOut[IsMemberOfGroupMethod][0].(bool)
	}
	var err error
	if t.ArgsOut[IsMemberOfGroupMethod][1] != nil {
		err = t.ArgsOut[IsMemberOfGroupMethod][1].(error)
	}
	return isMember, err
}

func (t TestRepo) GetGroupMembers(groupID string, filter *Filter) ([]UserGroupRelation, int, error) {
	t.ArgsIn[GetGroupMembersMethod][0] = groupID
	var members []UserGroupRelation
	if t.ArgsOut[GetGroupMembersMethod][0] != nil {
		testMembers := t.ArgsOut[GetGroupMembersMethod][0].([]TestUserGroupRelation)
		for _, v := range testMembers {
			members = append(members, v)
		}

	}
	var total int
	if t.ArgsOut[GetGroupMembersMethod][1] != nil {
		total = t.ArgsOut[GetGroupMembersMethod][1].(int)
	}
	var err error
	if t.ArgsOut[GetGroupMembersMethod][2] != nil {
		err = t.ArgsOut[GetGroupMembersMethod][2].(error)
	}
	return members, total, err
}

func (t TestRepo) IsAttachedToGroup(groupID string, policyID string) (bool, error) {
	t.ArgsIn[IsAttachedToGroupMethod][0] = groupID
	t.ArgsIn[IsAttachedToGroupMethod][1] = policyID
	var isAttached bool
	if t.ArgsOut[IsAttachedToGroupMethod][0] != nil {
		isAttached = t.ArgsOut[IsAttachedToGroupMethod][0].(bool)
	}
	var err error
	if t.ArgsOut[IsAttachedToGroupMethod][1] != nil {
		err = t.ArgsOut[IsAttachedToGroupMethod][1].(error)
	}
	return isAttached, err
}

func (t TestRepo) GetAttachedPolicies(groupID string, filter *Filter) ([]PolicyGroupRelation, int, error) {
	t.ArgsIn[GetAttachedPoliciesMethod][0] = groupID
	var policies []PolicyGroupRelation
	if t.ArgsOut[GetAttachedPoliciesMethod][0] != nil {
		testPolicies := t.ArgsOut[GetAttachedPoliciesMethod][0].([]TestPolicyGroupRelation)
		for _, v := range testPolicies {
			policies = append(policies, v)
		}
	}
	var total int
	if t.ArgsOut[GetAttachedPoliciesMethod][1] != nil {
		total = t.ArgsOut[GetAttachedPoliciesMethod][1].(int)
	}

	var err error
	if t.ArgsOut[GetAttachedPoliciesMethod][2] != nil {
		err = t.ArgsOut[GetAttachedPoliciesMethod][2].(error)
	}
	return policies, total, err
}

func (t TestRepo) GetGroupsFiltered(filter *Filter) ([]Group, int, error) {
	t.ArgsIn[GetGroupsFilteredMethod][0] = filter

	var groups []Group
	if t.ArgsOut[GetGroupsFilteredMethod][0] != nil {
		groups = t.ArgsOut[GetGroupsFilteredMethod][0].([]Group)
	}
	var total int
	if t.ArgsOut[GetGroupsFilteredMethod][1] != nil {
		total = t.ArgsOut[GetGroupsFilteredMethod][1].(int)
	}
	var err error
	if t.ArgsOut[GetGroupsFilteredMethod][2] != nil {
		err = t.ArgsOut[GetGroupsFilteredMethod][2].(error)
	}
	return groups, total, err
}
func (t TestRepo) RemoveGroup(id string) error {
	t.ArgsIn[RemoveGroupMethod][0] = id
	var err error
	if t.ArgsOut[RemoveGroupMethod][0] != nil {
		err = t.ArgsOut[RemoveGroupMethod][0].(error)
	}
	return err
}

func (t TestRepo) AddGroup(group Group) (*Group, error) {
	t.ArgsIn[AddGroupMethod][0] = group
	var created *Group
	if t.ArgsOut[AddGroupMethod][0] != nil {
		created = t.ArgsOut[AddGroupMethod][0].(*Group)
	}
	var err error
	if t.ArgsOut[AddGroupMethod][1] != nil {
		err = t.ArgsOut[AddGroupMethod][1].(error)
	}
	return created, err
}

func (t TestRepo) AddMember(userID string, groupID string) error {
	t.ArgsIn[AddMemberMethod][0] = userID
	t.ArgsIn[AddMemberMethod][1] = groupID
	var err error
	if t.ArgsOut[AddMemberMethod][0] != nil {
		err = t.ArgsOut[AddMemberMethod][0].(error)
	}
	return err
}

func (t TestRepo) RemoveMember(userID string, groupID string) error {
	t.ArgsIn[RemoveMemberMethod][0] = userID
	t.ArgsIn[RemoveMemberMethod][1] = groupID
	var err error
	if t.ArgsOut[RemoveMemberMethod][0] != nil {
		err = t.ArgsOut[RemoveMemberMethod][0].(error)
	}
	return err
}

func (t TestRepo) UpdateGroup(group Group) (*Group, error) {
	t.ArgsIn[UpdateGroupMethod][0] = group

	var updated *Group
	if t.ArgsOut[UpdateGroupMethod][0] != nil {
		updated = t.ArgsOut[UpdateGroupMethod][0].(*Group)
	}
	var err error
	if t.ArgsOut[UpdateGroupMethod][1] != nil {
		err = t.ArgsOut[UpdateGroupMethod][1].(error)
	}
	return updated, err
}

func (t TestRepo) AttachPolicy(groupID string, policyID string) error {
	t.ArgsIn[AttachPolicyMethod][0] = groupID
	t.ArgsIn[AttachPolicyMethod][1] = policyID
	var err error
	if t.ArgsOut[AttachPolicyMethod][0] != nil {
		err = t.ArgsOut[AttachPolicyMethod][0].(error)
	}
	return err
}
func (t TestRepo) DetachPolicy(groupID string, policyID string) error {
	t.ArgsIn[DetachPolicyMethod][0] = groupID
	t.ArgsIn[DetachPolicyMethod][1] = policyID
	var err error
	if t.ArgsOut[DetachPolicyMethod][0] != nil {
		err = t.ArgsOut[DetachPolicyMethod][0].(error)
	}
	return err
}

//////////////////
// Policy repo
//////////////////

func (t TestRepo) GetPolicyByName(org string, name string) (*Policy, error) {
	t.ArgsIn[GetPolicyByNameMethod][0] = org
	t.ArgsIn[GetPolicyByNameMethod][1] = name
	if specialFunc, ok := t.SpecialFuncs[GetPolicyByNameMethod].(func(org string, name string) (*Policy, error)); ok && specialFunc != nil {
		return specialFunc(org, name)
	}
	var policy *Policy
	if t.ArgsOut[GetPolicyByNameMethod][0] != nil {
		policy = t.ArgsOut[GetPolicyByNameMethod][0].(*Policy)
	}
	var err error
	if t.ArgsOut[GetPolicyByNameMethod][1] != nil {
		err = t.ArgsOut[GetPolicyByNameMethod][1].(error)
	}
	return policy, err
}

func (t TestRepo) AddPolicy(policy Policy) (*Policy, error) {
	t.ArgsIn[AddPolicyMethod][0] = policy
	var created *Policy
	if t.ArgsOut[AddPolicyMethod][0] != nil {
		created = t.ArgsOut[AddPolicyMethod][0].(*Policy)
	}
	var err error
	if t.ArgsOut[AddPolicyMethod][1] != nil {
		err = t.ArgsOut[AddPolicyMethod][1].(error)
	}
	return created, err
}

func (t TestRepo) UpdatePolicy(policy Policy) (*Policy, error) {
	t.ArgsIn[UpdatePolicyMethod][0] = policy

	var updated *Policy
	if t.ArgsOut[UpdatePolicyMethod][0] != nil {
		updated = t.ArgsOut[UpdatePolicyMethod][0].(*Policy)
	}
	var err error
	if t.ArgsOut[UpdatePolicyMethod][1] != nil {
		err = t.ArgsOut[UpdatePolicyMethod][1].(error)
	}
	return updated, err
}

func (t TestRepo) RemovePolicy(id string) error {
	t.ArgsIn[RemovePolicyMethod][0] = id
	var err error
	if t.ArgsOut[RemovePolicyMethod][0] != nil {
		err = t.ArgsOut[RemovePolicyMethod][0].(error)
	}
	return err
}

func (t TestRepo) GetPoliciesFiltered(filter *Filter) ([]Policy, int, error) {
	t.ArgsIn[GetPoliciesFilteredMethod][0] = filter

	var policies []Policy
	if t.ArgsOut[GetPoliciesFilteredMethod][0] != nil {
		policies = t.ArgsOut[GetPoliciesFilteredMethod][0].([]Policy)
	}
	var total int
	if t.ArgsOut[GetPoliciesFilteredMethod][1] != nil {
		total = t.ArgsOut[GetPoliciesFilteredMethod][1].(int)
	}
	var err error
	if t.ArgsOut[GetPoliciesFilteredMethod][2] != nil {
		err = t.ArgsOut[GetPoliciesFilteredMethod][2].(error)
	}
	return policies, total, err
}

func (t TestRepo) GetAttachedGroups(policyID string, filter *Filter) ([]PolicyGroupRelation, int, error) {
	t.ArgsIn[GetAttachedGroupsMethod][0] = policyID

	var groups []PolicyGroupRelation
	if t.ArgsOut[GetAttachedGroupsMethod][0] != nil {
		testGroups := t.ArgsOut[GetAttachedGroupsMethod][0].([]TestPolicyGroupRelation)
		for _, v := range testGroups {
			groups = append(groups, v)
		}
	}
	var total int
	if t.ArgsOut[GetGroupsByUserIDMethod][1] != nil {
		total = t.ArgsOut[GetAttachedGroupsMethod][1].(int)
	}
	var err error
	if t.ArgsOut[GetAttachedGroupsMethod][2] != nil {
		err = t.ArgsOut[GetAttachedGroupsMethod][2].(error)
	}
	return groups, total, err
}

func (t TestRepo) OrderByValidColumns(action string) []string {
	t.ArgsIn[OrderByValidColumnsMethod][0] = action
	var validColumns []string
	if t.ArgsOut[OrderByValidColumnsMethod][0] != nil {
		validColumns = t.ArgsOut[OrderByValidColumnsMethod][0].([]string)
	}
	return validColumns
}

//////////////
// Proxy repo
//////////////

func (t TestRepo) GetProxyResources() ([]ProxyResource, error) {
	var resources []ProxyResource
	if t.ArgsOut[GetProxyResourcesMethod][0] != nil {
		resources = t.ArgsOut[GetProxyResourcesMethod][0].([]ProxyResource)
	}
	var err error
	if t.ArgsOut[GetProxyResourcesMethod][1] != nil {
		err = t.ArgsOut[GetProxyResourcesMethod][1].(error)
	}
	return resources, err
}

// Private helper methods

func getRandomString(runeValue []rune, n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = runeValue[rand.Intn(len(runeValue))]
	}
	return string(b)
}

func getResources(number int, baseUrn string) []string {
	resources := make([]string, number, number)
	for i := 0; i < number; i++ {
		resources[i] = fmt.Sprintf(baseUrn+"%v", i+1)
	}
	return resources
}

func checkMethodResponse(t *testing.T, testcase string, expectedError error, receivedError error, expectedResponse interface{}, receivedResponse interface{}) {
	if expectedError != nil {
		apiError, ok := receivedError.(*Error)
		if !ok || apiError == nil {
			t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", testcase, receivedError)
			return
		}
		if diff := pretty.Compare(apiError, expectedError); diff != "" {
			t.Errorf("Test %v failed. Received different errors (received/wanted) %v", testcase, diff)
			return
		}
	} else {
		if receivedError != nil {
			t.Errorf("Test %v failed: %v", testcase, receivedError)
			return
		}
		if diff := pretty.Compare(receivedResponse, expectedResponse); diff != "" {
			t.Errorf("Test %v failed. Received different responses (received/wanted) %v", testcase, diff)
			return
		}
	}
}
