package api

import (
	"fmt"
	"time"

	"github.com/Tecsisa/foulkon/database"
	"github.com/satori/go.uuid"
)

// TYPE DEFINITIONS

// Group domain
type Group struct {
	ID       string    `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Path     string    `json:"path,omitempty"`
	Org      string    `json:"org,omitempty"`
	Urn      string    `json:"urn,omitempty"`
	CreateAt time.Time `json:"createAt,omitempty"`
	UpdateAt time.Time `json:"updateAt,omitempty"`
}

func (g Group) String() string {
	return fmt.Sprintf("[id: %v, name: %v, path: %v, org: %v, urn: %v, createAt: %v]",
		g.ID, g.Name, g.Path, g.Org, g.Urn, g.CreateAt.Format("2006-01-02 15:04:05 MST"))
}

func (g Group) GetUrn() string {
	return g.Urn
}

// Group identifier to retrieve them from DB
type GroupIdentity struct {
	Org  string `json:"org,omitempty"`
	Name string `json:"name,omitempty"`
}

type GroupMembers struct {
	User     string    `json:"user,omitempty"`
	CreateAt time.Time `json:"joined,omitempty"`
}

type GroupPolicies struct {
	Policy   string    `json:"policy,omitempty"`
	CreateAt time.Time `json:"attached,omitempty"`
}

// GROUP API IMPLEMENTATION

func (api WorkerAPI) AddGroup(requestInfo RequestInfo, org string, name string, path string) (*Group, error) {
	// Validate fields
	if !IsValidName(name) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: name %v", name),
		}
	}
	if !IsValidOrg(org) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: org %v", org),
		}
	}
	if !IsValidPath(path) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: path %v", path),
		}
	}

	group := createGroup(org, name, path)

	// Check restrictions
	groupsFiltered, err := api.GetAuthorizedGroups(requestInfo, group.Urn, GROUP_ACTION_CREATE_GROUP, []Group{group})
	if err != nil {
		return nil, err
	}
	if len(groupsFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, group.Urn),
		}
	}

	// Check if group already exists
	_, err = api.GroupRepo.GetGroupByName(org, name)

	// Check if group could be retrieved
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		switch dbError.Code {
		// Group doesn't exist in DB, so we can create it
		case database.GROUP_NOT_FOUND:
			// Create group
			createdGroup, err := api.GroupRepo.AddGroup(group)

			// Check if there is an unexpected error in DB
			if err != nil {
				//Transform to DB error
				dbError := err.(*database.Error)
				return nil, &Error{
					Code:    UNKNOWN_API_ERROR,
					Message: dbError.Message,
				}
			}
			LogOperation(requestInfo.RequestID, requestInfo.Identifier, fmt.Sprintf("Group created %+v", createdGroup))
			return createdGroup, nil
		default: // Unexpected error
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	} else {
		return nil, &Error{
			Code:    GROUP_ALREADY_EXIST,
			Message: fmt.Sprintf("Unable to create group, group with org %v and name %v already exists", org, name),
		}
	}

}

func (api WorkerAPI) GetGroupByName(requestInfo RequestInfo, org string, name string) (*Group, error) {
	// Validate fields
	if !IsValidName(name) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: name %v", name),
		}
	}
	if !IsValidOrg(org) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: org %v", org),
		}
	}

	// Call repo to retrieve the group
	group, err := api.GroupRepo.GetGroupByName(org, name)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		// Group doesn't exist in DB
		switch dbError.Code {
		case database.GROUP_NOT_FOUND:
			return nil, &Error{
				Code:    GROUP_BY_ORG_AND_NAME_NOT_FOUND,
				Message: dbError.Message,
			}
		default: // Unexpected error
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	}

	// Check restrictions
	groupsFiltered, err := api.GetAuthorizedGroups(requestInfo, group.Urn, GROUP_ACTION_GET_GROUP, []Group{*group})
	if err != nil {
		return nil, err
	}

	// Check if we have our user authorized
	if len(groupsFiltered) > 0 {
		groupsFiltered := groupsFiltered[0]
		return &groupsFiltered, nil
	}
	return nil, &Error{
		Code: UNAUTHORIZED_RESOURCES_ERROR,
		Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
			requestInfo.Identifier, group.Urn),
	}
}

func (api WorkerAPI) ListGroups(requestInfo RequestInfo, filter *Filter) ([]GroupIdentity, int, error) {
	// Validate fields
	var total int
	orderByValidColumns := api.GroupRepo.OrderByValidColumns(GROUP_ACTION_LIST_GROUPS)
	err := validateFilter(filter, orderByValidColumns)
	if err != nil {
		return nil, total, err
	}

	// Call repo to retrieve the groups
	groups, total, err := api.GroupRepo.GetGroupsFiltered(filter)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, total, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	// Check restrictions to list
	var urnPrefix string
	if len(filter.Org) == 0 {
		urnPrefix = "*"
	} else {
		urnPrefix = GetUrnPrefix(filter.Org, RESOURCE_GROUP, filter.PathPrefix)
	}
	filteredGroups, err := api.GetAuthorizedGroups(requestInfo, urnPrefix, GROUP_ACTION_LIST_GROUPS, groups)
	if err != nil {
		return nil, total, err
	}

	// Transform to identifiers
	groupIDs := []GroupIdentity{}
	for _, g := range filteredGroups {
		groupIDs = append(groupIDs, GroupIdentity{
			Org:  g.Org,
			Name: g.Name,
		})
	}

	return groupIDs, total, nil
}

func (api WorkerAPI) UpdateGroup(requestInfo RequestInfo, org string, name string, newName string, newPath string) (*Group, error) {
	// Validate fields
	if !IsValidName(newName) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: new name %v", newName),
		}
	}
	if !IsValidPath(newPath) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: new path %v", newPath),
		}
	}

	// Call repo to retrieve the old group
	oldGroup, err := api.GetGroupByName(requestInfo, org, name)
	if err != nil {
		return nil, err
	}

	// Check restrictions
	groupsFiltered, err := api.GetAuthorizedGroups(requestInfo, oldGroup.Urn, GROUP_ACTION_UPDATE_GROUP, []Group{*oldGroup})
	if err != nil {
		return nil, err
	}
	if len(groupsFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, oldGroup.Urn),
		}
	}

	// Check if a group with "newName" already exists
	newGroup, err := api.GetGroupByName(requestInfo, org, newName)

	if err == nil && oldGroup.ID != newGroup.ID {
		// Group already exists
		return nil, &Error{
			Code:    GROUP_ALREADY_EXIST,
			Message: fmt.Sprintf("Group name: %v already exists", newName),
		}
	}

	if err != nil {
		if apiError := err.(*Error); apiError.Code != GROUP_BY_ORG_AND_NAME_NOT_FOUND {
			return nil, err
		}
	}

	auxGroup := Group{
		Urn: CreateUrn(org, RESOURCE_GROUP, newPath, newName),
	}

	// Check restrictions
	groupsFiltered, err = api.GetAuthorizedGroups(requestInfo, auxGroup.Urn, GROUP_ACTION_UPDATE_GROUP, []Group{auxGroup})
	if err != nil {
		return nil, err
	}
	if len(groupsFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, auxGroup.Urn),
		}
	}

	// Update group
	group := Group{
		ID:       oldGroup.ID,
		Name:     newName,
		Path:     newPath,
		Org:      oldGroup.Org,
		Urn:      auxGroup.Urn,
		CreateAt: oldGroup.CreateAt,
		UpdateAt: time.Now().UTC(),
	}

	updatedGroup, err := api.GroupRepo.UpdateGroup(group)

	// Check unexpected DB error
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	LogOperation(requestInfo.RequestID, requestInfo.Identifier, fmt.Sprintf("Group updated from %+v to %+v", oldGroup, updatedGroup))
	return updatedGroup, nil

}

func (api WorkerAPI) RemoveGroup(requestInfo RequestInfo, org string, name string) error {
	// Call repo to retrieve the group
	group, err := api.GetGroupByName(requestInfo, org, name)
	if err != nil {
		return err
	}

	// Check restrictions
	groupsFiltered, err := api.GetAuthorizedGroups(requestInfo, group.Urn, GROUP_ACTION_DELETE_GROUP, []Group{*group})
	if err != nil {
		return err
	}
	if len(groupsFiltered) < 1 {
		return &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, group.Urn),
		}
	}

	err = api.GroupRepo.RemoveGroup(group.ID)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	LogOperation(requestInfo.RequestID, requestInfo.Identifier, fmt.Sprintf("Group deleted %v", group))
	return nil
}

func (api WorkerAPI) AddMember(requestInfo RequestInfo, externalId string, name string, org string) error {
	// Call repo to retrieve the group
	groupDB, err := api.GetGroupByName(requestInfo, org, name)
	if err != nil {
		return err
	}

	// Check restrictions
	groupsFiltered, err := api.GetAuthorizedGroups(requestInfo, groupDB.Urn, GROUP_ACTION_ADD_MEMBER, []Group{*groupDB})
	if err != nil {
		return err
	}
	if len(groupsFiltered) < 1 {
		return &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, groupDB.Urn),
		}
	}

	// Call repo to retrieve the user
	userDB, err := api.GetUserByExternalID(requestInfo, externalId)
	if err != nil {
		return err
	}

	// Call repo to retrieve the GroupUserRelation
	isMember, err := api.GroupRepo.IsMemberOfGroup(userDB.ID, groupDB.ID)
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	// Error handling
	if isMember {
		return &Error{
			Code:    USER_IS_ALREADY_A_MEMBER_OF_GROUP,
			Message: fmt.Sprintf("User: %v is already a member of Group: %v", externalId, name),
		}
	}

	// Add Member
	err = api.GroupRepo.AddMember(userDB.ID, groupDB.ID)

	// Check if there is an unexpected error in DB
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}
	LogOperation(requestInfo.RequestID, requestInfo.Identifier, fmt.Sprintf("Member %+v added to group %+v", userDB, groupDB))
	return nil
}

func (api WorkerAPI) RemoveMember(requestInfo RequestInfo, externalId string, name string, org string) error {
	// Call repo to retrieve the group
	groupDB, err := api.GetGroupByName(requestInfo, org, name)
	if err != nil {
		return err
	}

	// Check restrictions
	groupsFiltered, err := api.GetAuthorizedGroups(requestInfo, groupDB.Urn, GROUP_ACTION_REMOVE_MEMBER, []Group{*groupDB})
	if err != nil {
		return err
	}
	if len(groupsFiltered) < 1 {
		return &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, groupDB.Urn),
		}
	}

	// Call repo to retrieve the user
	userDB, err := api.GetUserByExternalID(requestInfo, externalId)
	if err != nil {
		return err
	}

	// Call repo to check if user is a member of group
	isMember, err := api.GroupRepo.IsMemberOfGroup(userDB.ID, groupDB.ID)
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	if !isMember {
		return &Error{
			Code: USER_IS_NOT_A_MEMBER_OF_GROUP,
			Message: fmt.Sprintf("User with externalId %v is not a member of group with org %v and name %v",
				userDB.ExternalID, groupDB.Org, groupDB.Name),
		}
	}

	// Remove Member
	err = api.GroupRepo.RemoveMember(userDB.ID, groupDB.ID)

	// Check if there is an unexpected error in DB
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	LogOperation(requestInfo.RequestID, requestInfo.Identifier, fmt.Sprintf("Member %+v removed from group %+v", userDB, groupDB))
	return nil
}

func (api WorkerAPI) ListMembers(requestInfo RequestInfo, filter *Filter) ([]GroupMembers, int, error) {
	// Validate fields
	var total int
	orderByValidColumns := api.UserRepo.OrderByValidColumns(GROUP_ACTION_LIST_MEMBERS)
	err := validateFilter(filter, orderByValidColumns)
	if err != nil {
		return nil, total, err
	}

	// Call repo to retrieve the group
	group, err := api.GetGroupByName(requestInfo, filter.Org, filter.GroupName)
	if err != nil {
		return nil, total, err
	}

	// Check restrictions
	groupsFiltered, err := api.GetAuthorizedGroups(requestInfo, group.Urn, GROUP_ACTION_LIST_MEMBERS, []Group{*group})
	if err != nil {
		return nil, total, err
	}
	if len(groupsFiltered) < 1 {
		return nil, total, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, group.Urn),
		}
	}

	// Get Members
	users, total, err := api.GroupRepo.GetGroupMembers(group.ID, filter)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, total, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	members := []GroupMembers{}
	if users != nil {
		members = make([]GroupMembers, len(users), cap(users))
		for i, m := range users {
			members[i] = GroupMembers{
				User:     m.GetUser().ExternalID,
				CreateAt: m.GetDate(),
			}
		}
	}

	return members, total, nil
}

func (api WorkerAPI) AttachPolicyToGroup(requestInfo RequestInfo, org string, name string, policyName string) error {

	// Check if group exists
	group, err := api.GetGroupByName(requestInfo, org, name)
	if err != nil {
		return err
	}

	// Check restrictions
	groupsFiltered, err := api.GetAuthorizedGroups(requestInfo, group.Urn, GROUP_ACTION_ATTACH_GROUP_POLICY, []Group{*group})
	if err != nil {
		return err
	}
	if len(groupsFiltered) < 1 {
		return &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, group.Urn),
		}
	}

	// Check if policy exists
	policy, err := api.GetPolicyByName(requestInfo, org, policyName)
	if err != nil {
		return err
	}

	// Check existing relationship
	isAttached, err := api.GroupRepo.IsAttachedToGroup(group.ID, policy.ID)
	if err != nil {
		dbError := err.(*database.Error)
		return &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	if isAttached {
		// Unexpected error
		return &Error{
			Code:    POLICY_IS_ALREADY_ATTACHED_TO_GROUP,
			Message: fmt.Sprintf("Policy: %v is already attached to Group: %v", policy.Name, group.Name),
		}
	}

	// Attach Policy to Group
	err = api.GroupRepo.AttachPolicy(group.ID, policy.ID)

	if err != nil {
		dbError := err.(*database.Error)
		return &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	LogOperation(requestInfo.RequestID, requestInfo.Identifier, fmt.Sprintf("Policy %+v attached to group %+v", policy, group))
	return nil
}

func (api WorkerAPI) DetachPolicyToGroup(requestInfo RequestInfo, org string, name string, policyName string) error {

	// Check if group exists
	group, err := api.GetGroupByName(requestInfo, org, name)
	if err != nil {
		return err
	}

	// Check restrictions
	groupsFiltered, err := api.GetAuthorizedGroups(requestInfo, group.Urn, GROUP_ACTION_DETACH_GROUP_POLICY, []Group{*group})
	if err != nil {
		return err
	}
	if len(groupsFiltered) < 1 {
		return &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, group.Urn),
		}
	}

	// Check if policy exists
	policy, err := api.GetPolicyByName(requestInfo, org, policyName)
	if err != nil {
		return err
	}

	// Check existing relationship
	isAttached, err := api.GroupRepo.IsAttachedToGroup(group.ID, policy.ID)
	if err != nil {
		dbError := err.(*database.Error)
		return &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	if !isAttached {
		return &Error{
			Code: POLICY_IS_NOT_ATTACHED_TO_GROUP,
			Message: fmt.Sprintf("Policy with org %v and name %v is not attached to group with org %v and name %v",
				policy.Org, policy.Name, group.Org, group.Name),
		}

	}

	// Detach Policy to Group
	err = api.GroupRepo.DetachPolicy(group.ID, policy.ID)

	if err != nil {
		dbError := err.(*database.Error)
		return &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	LogOperation(requestInfo.RequestID, requestInfo.Identifier, fmt.Sprintf("Policy %+v detached from group %+v", policy, group))
	return nil
}

func (api WorkerAPI) ListAttachedGroupPolicies(requestInfo RequestInfo, filter *Filter) ([]GroupPolicies, int, error) {
	// Validate fields
	var total int
	orderByValidColumns := api.UserRepo.OrderByValidColumns(GROUP_ACTION_LIST_ATTACHED_GROUP_POLICIES)
	err := validateFilter(filter, orderByValidColumns)
	if err != nil {
		return nil, total, err
	}

	// Check if group exists
	group, err := api.GetGroupByName(requestInfo, filter.Org, filter.GroupName)
	if err != nil {
		return nil, total, err
	}

	// Check restrictions
	groupsFiltered, err := api.GetAuthorizedGroups(requestInfo, group.Urn, GROUP_ACTION_LIST_ATTACHED_GROUP_POLICIES, []Group{*group})
	if err != nil {
		return nil, total, err
	}
	if len(groupsFiltered) < 1 {
		return nil, total, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, group.Urn),
		}
	}

	// Call repo to retrieve the GroupPolicyRelations
	attachedPolicies, total, err := api.GroupRepo.GetAttachedPolicies(group.ID, filter)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, total, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	policies := []GroupPolicies{}
	if attachedPolicies != nil {
		policies = make([]GroupPolicies, len(attachedPolicies), cap(attachedPolicies))
		for i, m := range attachedPolicies {
			policies[i] = GroupPolicies{
				Policy:   m.GetPolicy().Name,
				CreateAt: m.GetDate(),
			}
		}
	}

	return policies, total, nil
}

// PRIVATE HELPER METHODS

func createGroup(org string, name string, path string) Group {
	urn := CreateUrn(org, RESOURCE_GROUP, path, name)
	group := Group{
		ID:       uuid.NewV4().String(),
		Name:     name,
		Path:     path,
		CreateAt: time.Now().UTC(),
		UpdateAt: time.Now().UTC(),
		Urn:      urn,
		Org:      org,
	}

	return group
}
