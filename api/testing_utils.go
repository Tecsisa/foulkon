package api

const (
	GetUserByExternalIDMethod       = "GetUserByExternalID"
	AddUserMethod                   = "AddUser"
	UpdateUserMethod                = "UpdateUser"
	GetUsersFilteredMethod          = "GetUsersFiltered"
	GetGroupsByUserIDMethod         = "GetGroupsByUserID"
	RemoveUserMethod                = "RemoveUser"
	GetGroupByNameMethod            = "GetGroupByName"
	IsMemberOfGroupMethod           = "IsMemberOfGroup"
	GetGroupMembersMethod           = "GetGroupMembers"
	IsAttachedToGroupMethod         = "IsAttachedToGroup"
	GetPoliciesAttachedMethod       = "GetPoliciesAttached"
	GetGroupsFilteredMethod         = "GetGroupsFiltered"
	RemoveGroupMethod               = "RemoveGroup"
	AddGroupMethod                  = "AddGroup"
	AddMemberMethod                 = "AddMember"
	RemoveMemberMethod              = "RemoveMember"
	UpdateGroupMethod               = "UpdateGroup"
	AttachPolicyMethod              = "AttachPolicy"
	DetachPolicyMethod              = "DetachPolicy"
	GetPolicyByNameMethod           = "GetPolicyByName"
	AddPolicyMethod                 = "AddPolicy"
	UpdatePolicyMethod              = "UpdatePolicy"
	RemovePolicyMethod              = "RemovePolicy"
	GetPoliciesFilteredMethod       = "GetPoliciesFiltered"
	GetAllPolicyGroupRelationMethod = "GetAllPolicyGroupRelation"
)

// Test repo that implements all manager interfaces
type TestRepo struct {
	ArgsIn  map[string][]interface{}
	ArgsOut map[string][]interface{}
}

//////////////////
// User repo
//////////////////
func (t TestRepo) GetUserByExternalID(id string) (*User, error) {
	t.ArgsIn[GetUserByExternalIDMethod][0] = id
	var err error
	if t.ArgsOut[GetUserByExternalIDMethod][1] != nil {
		err = t.ArgsOut[GetUserByExternalIDMethod][1].(error)
	}
	return t.ArgsOut[GetUserByExternalIDMethod][0].(*User), err
}

func (t TestRepo) AddUser(user User) (*User, error) {
	return nil, nil
}

func (t TestRepo) UpdateUser(user User, newPath string, newUrn string) (*User, error) {
	return nil, nil
}

func (t TestRepo) GetUsersFiltered(pathPrefix string) ([]User, error) {
	return nil, nil
}

func (t TestRepo) GetGroupsByUserID(id string) ([]Group, error) {
	return nil, nil
}

func (t TestRepo) RemoveUser(id string) error {
	return nil
}

//////////////////
// Group repo
//////////////////
func (t TestRepo) GetGroupByName(org string, name string) (*Group, error) {
	return nil, nil
}
func (t TestRepo) IsMemberOfGroup(userID string, groupID string) (bool, error) {
	return false, nil
}
func (t TestRepo) GetGroupMembers(groupID string) ([]User, error) {
	return nil, nil
}
func (t TestRepo) IsAttachedToGroup(groupID string, policyID string) (bool, error) {
	return false, nil
}
func (t TestRepo) GetPoliciesAttached(groupID string) ([]Policy, error) {
	return nil, nil
}
func (t TestRepo) GetGroupsFiltered(org string, pathPrefix string) ([]Group, error) {
	return nil, nil
}
func (t TestRepo) RemoveGroup(id string) error {
	return nil
}

func (t TestRepo) AddGroup(group Group) (*Group, error) {
	return nil, nil
}
func (t TestRepo) AddMember(userID string, groupID string) error {
	return nil
}
func (t TestRepo) RemoveMember(userID string, groupID string) error {
	return nil
}
func (t TestRepo) UpdateGroup(group Group, newName string, newPath string, newUrn string) (*Group, error) {
	return nil, nil
}
func (t TestRepo) AttachPolicy(groupID string, policyID string) error {
	return nil
}
func (t TestRepo) DetachPolicy(groupID string, policyID string) error {
	return nil
}

//////////////////
// Policy repo
//////////////////

func (t TestRepo) GetPolicyByName(org string, name string) (*Policy, error) {
	return nil, nil
}

func (t TestRepo) AddPolicy(policy Policy) (*Policy, error) {
	return nil, nil
}

func (t TestRepo) UpdatePolicy(policy Policy, newName string, newPath string, newUrn string, newStatements []Statement) (*Policy, error) {
	return nil, nil
}

func (t TestRepo) RemovePolicy(id string) error {
	return nil
}

func (t TestRepo) GetPoliciesFiltered(org string, pathPrefix string) ([]Policy, error) {
	return nil, nil
}

func (t TestRepo) GetAllPolicyGroupRelation(policyID string) ([]Group, error) {
	return nil, nil
}
