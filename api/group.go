package api

import (
	"fmt"
	"time"

	"github.com/satori/go.uuid"
	"github.com/tecsisa/authorizr/database"
)

// Group domain
type Group struct {
	ID       string    `json:"ID, omitempty"`
	Name     string    `json:"Name, omitempty"`
	Path     string    `json:"Path, omitempty"`
	CreateAt time.Time `json:"CreateAt, omitempty"`
	Urn      string    `json:"Urn, omitempty"`
	Org      string    `json:"Org, omitempty"`
}

func (g Group) GetUrn() string {
	return g.Urn
}

type GroupMembers struct {
	Group Group  `json:"Group, omitempty"`
	Users []User `json:"Users, omitempty"`
}

type GroupPolicies struct {
	Group    Group    `json:"Group, omitempty"`
	Policies []Policy `json:"Policies, omitempty"`
}

// Add an Group to database if not exist
func (api *AuthAPI) AddGroup(authenticatedUser AuthenticatedUser, org string, name string, path string) (*Group, error) {
	// Validate name
	if !IsValidName(name) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: Name %v", name),
		}
	}
	// Validate path
	if !IsValidPath(path) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: Path %v", path),
		}
	}
	// Create group
	group := createGroup(org, name, path)

	// Check restrictions
	groupsFiltered, err := api.GetGroupsAuthorized(authenticatedUser, group.Urn, GROUP_ACTION_CREATE_GROUP, []Group{group})
	if err != nil {
		return nil, err
	}

	// Check if we have our user authorized
	if len(groupsFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with external ID %v is not allowed to access to resource %v",
				authenticatedUser.Identifier, group.Urn),
		}
	}

	// Check if group already exist
	_, err = api.GroupRepo.GetGroupByName(org, name)

	// Check if group could be retrieved
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		switch dbError.Code {
		// Group doesn't exist in DB
		case database.GROUP_NOT_FOUND:
			// Create group
			groupCreated, err := api.GroupRepo.AddGroup(group)

			// Check if there is an unexpected error in DB
			if err != nil {
				//Transform to DB error
				dbError := err.(*database.Error)
				return nil, &Error{
					Code:    UNKNOWN_API_ERROR,
					Message: dbError.Message,
				}
			}

			// Return group created
			return groupCreated, nil
		default: // Unexpected error
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	} else {
		return nil, &Error{
			Code:    GROUP_ALREADY_EXIST,
			Message: fmt.Sprintf("Unable to create group, group with org %v and name %v already exist", org, name),
		}
	}

}

//  Add a new member into an existing group
func (api *AuthAPI) AddMember(authenticatedUser AuthenticatedUser, userID string, groupName string, org string) error {
	// Call repo to retrieve the group
	groupDB, err := api.GetGroupByName(authenticatedUser, org, groupName)
	if err != nil {
		return err
	}

	// Check restrictions
	groupsFiltered, err := api.GetGroupsAuthorized(authenticatedUser, groupDB.Urn, GROUP_ACTION_ADD_MEMBER, []Group{*groupDB})
	if err != nil {
		return err
	}

	// Check if we have our user authorized
	if len(groupsFiltered) < 1 {
		return &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with external ID %v is not allowed to access to resource %v",
				authenticatedUser.Identifier, groupDB.Urn),
		}
	}

	// Call repo to retrieve the user
	userDB, err := api.GetUserByExternalId(authenticatedUser, userID)
	if err != nil {
		return err
	}

	// Call repo to retrieve the GroupUserRelation
	groupMembers, err := api.GroupRepo.GetGroupUserRelation(userDB.ID, groupDB.ID)

	// Error handling
	if groupMembers != nil {
		return &Error{
			Code:    USER_IS_ALREADY_A_MEMBER_OF_GROUP,
			Message: fmt.Sprintf("User: %v is already a member of Group: %v", userID, groupName),
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

	return nil
}

//  Remove a member from a group
func (api *AuthAPI) RemoveMember(authenticatedUser AuthenticatedUser, userID string, groupName string, org string) error {
	// Call repo to retrieve the group
	groupDB, err := api.GetGroupByName(authenticatedUser, org, groupName)
	if err != nil {
		return err
	}

	// Check restrictions
	groupsFiltered, err := api.GetGroupsAuthorized(authenticatedUser, groupDB.Urn, GROUP_ACTION_REMOVE_MEMBER, []Group{*groupDB})
	if err != nil {
		return err
	}

	// Check if we have our user authorized
	if len(groupsFiltered) < 1 {
		return &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with external ID %v is not allowed to access to resource %v",
				authenticatedUser.Identifier, groupDB.Urn),
		}
	}

	// Call repo to retrieve the user
	userDB, err := api.GetUserByExternalId(authenticatedUser, userID)
	if err != nil {
		return err
	}

	// Call repo to retrieve the GroupUserRelation
	_, err = api.GroupRepo.GetGroupUserRelation(userDB.ID, groupDB.ID)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		// Relation doesn't exist in DB
		switch dbError.Code {
		case database.GROUP_USER_RELATION_NOT_FOUND:
			return &Error{
				Code:    USER_IS_NOT_A_MEMBER_OF_GROUP,
				Message: dbError.Message,
			}
		default: // Unexpected error
			return &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
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

	return nil
}

// List members of a group
func (api *AuthAPI) ListMembers(authenticatedUser AuthenticatedUser, org string, groupName string) (*GroupMembers, error) {
	// Call repo to retrieve the group
	group, err := api.GetGroupByName(authenticatedUser, org, groupName)
	if err != nil {
		return nil, err
	}

	// Check restrictions
	groupsFiltered, err := api.GetGroupsAuthorized(authenticatedUser, group.Urn, GROUP_ACTION_LIST_MEMBERS, []Group{*group})
	if err != nil {
		return nil, err
	}

	// Check if we have our user authorized
	if len(groupsFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with external ID %v is not allowed to access to resource %v",
				authenticatedUser.Identifier, group.Urn),
		}
	}

	// Get Members
	members, err := api.GroupRepo.GetAllGroupUserRelation(group.ID)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	// Return members
	return members, nil
}

// Remove group
func (api *AuthAPI) RemoveGroup(authenticatedUser AuthenticatedUser, org string, name string) error {
	// Call repo to retrieve the group
	group, err := api.GetGroupByName(authenticatedUser, org, name)
	if err != nil {
		return err
	}

	// Check restrictions
	groupsFiltered, err := api.GetGroupsAuthorized(authenticatedUser, group.Urn, GROUP_ACTION_DELETE_GROUP, []Group{*group})
	if err != nil {
		return err
	}

	// Check if we have our user authorized
	if len(groupsFiltered) < 1 {
		return &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with external ID %v is not allowed to access to resource %v",
				authenticatedUser.Identifier, group.Urn),
		}
	}

	// Remove group with given org and name
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

	return nil
}

func (api *AuthAPI) GetGroupByName(authenticatedUser AuthenticatedUser, org string, name string) (*Group, error) {
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
	groupsFiltered, err := api.GetGroupsAuthorized(authenticatedUser, group.Urn, GROUP_ACTION_GET_GROUP, []Group{*group})
	if err != nil {
		return nil, err
	}

	// Check if we have our user authorized
	if len(groupsFiltered) > 0 {
		groupsFiltered := groupsFiltered[0]
		return &groupsFiltered, nil
	} else {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with external ID %v is not allowed to access to resource %v",
				authenticatedUser.Identifier, group.Urn),
		}
	}

}

func (api *AuthAPI) GetListGroups(authenticatedUser AuthenticatedUser, org string, pathPrefix string) ([]Group, error) {
	// Call repo to retrieve the groups
	groups, err := api.GroupRepo.GetGroupsFiltered(org, pathPrefix)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	// Check restrictions to list
	urnPrefix := GetUrnPrefix(org, RESOURCE_GROUP, pathPrefix)
	groupsFiltered, err := api.GetGroupsAuthorized(authenticatedUser, urnPrefix, GROUP_ACTION_LIST_GROUPS, groups)
	if err != nil {
		return nil, err
	}

	// Return groups
	return groupsFiltered, nil
}

// Update Group to database if exist
func (api *AuthAPI) UpdateGroup(authenticatedUser AuthenticatedUser, org string, groupName string, newName string, newPath string) (*Group, error) {
	// Validate name
	if !IsValidName(newName) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: Name %v", newName),
		}
	}
	// Validate path
	if !IsValidPath(newPath) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: Path %v", newPath),
		}
	}

	// Call repo to retrieve the group
	group, err := api.GetGroupByName(authenticatedUser, org, groupName)
	if err != nil {
		return nil, err
	}

	// Check restrictions
	groupsFiltered, err := api.GetGroupsAuthorized(authenticatedUser, group.Urn, GROUP_ACTION_UPDATE_GROUP, []Group{*group})
	if err != nil {
		return nil, err
	}

	// Check if we have our user authorized
	if len(groupsFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with external ID %v is not allowed to access to resource %v",
				authenticatedUser.Identifier, group.Urn),
		}
	}

	// Check if group with newName exist
	_, err = api.GetGroupByName(authenticatedUser, org, newName)

	if err == nil {
		// Group already exists
		return nil, &Error{
			Code:    GROUP_ALREADY_EXIST,
			Message: fmt.Sprintf("Group name: %v already exists", newName),
		}
	}

	if err != nil {
		apiError := err.(*Error)
		switch apiError.Code {
		case UNAUTHORIZED_RESOURCES_ERROR, UNKNOWN_API_ERROR:
			return nil, err
		default: //Do nothing
		}
	}

	// Get Group updated
	groupToUpdate := createGroup(org, newName, newPath)

	// Check restrictions
	groupsFiltered, err = api.GetGroupsAuthorized(authenticatedUser, groupToUpdate.Urn, GROUP_ACTION_UPDATE_GROUP, []Group{groupToUpdate})
	if err != nil {
		return nil, err
	}

	// Check if we have our user authorized
	if len(groupsFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with external ID %v is not allowed to access to resource %v",
				authenticatedUser.Identifier, groupToUpdate.Urn),
		}
	}

	// Update group
	group, err = api.GroupRepo.UpdateGroup(*group, newName, newPath, groupToUpdate.Urn)

	// Check if there is an unexpected error in DB
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	return group, nil

}

func (api *AuthAPI) AttachPolicyToGroup(authenticatedUser AuthenticatedUser, org string, groupName string, policyName string) error {
	// Check if group exist
	group, err := api.GetGroupByName(authenticatedUser, org, groupName)
	if err != nil {
		return err
	}

	// Check restrictions
	groupsFiltered, err := api.GetGroupsAuthorized(authenticatedUser, group.Urn, GROUP_ACTION_ATTACH_GROUP_POLICY, []Group{*group})
	if err != nil {
		return err
	}

	// Check if we have our user authorized
	if len(groupsFiltered) < 1 {
		return &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with external ID %v is not allowed to access to resource %v",
				authenticatedUser.Identifier, group.Urn),
		}
	}

	// Check if policy exist
	policy, err := api.GetPolicy(authenticatedUser, org, policyName)
	if err != nil {
		return err
	}

	// Check if exist this relation
	groupPolicies, err := api.GroupRepo.GetGroupPolicyRelation(group.ID, policy.ID)

	if groupPolicies != nil {
		// Unexpected error
		return &Error{
			Code:    POLICY_IS_ALREADY_ATTACHED_TO_GROUP,
			Message: fmt.Sprintf("Policy: %v is already attached to Group: %v", policy.ID, group.ID),
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

	return nil
}

func (api *AuthAPI) DetachPolicyToGroup(authenticatedUser AuthenticatedUser, org string, groupName string, policyName string) error {
	// Check if group exist
	group, err := api.GetGroupByName(authenticatedUser, org, groupName)
	if err != nil {
		return err
	}

	// Check restrictions
	groupsFiltered, err := api.GetGroupsAuthorized(authenticatedUser, group.Urn, GROUP_ACTION_DETACH_GROUP_POLICY, []Group{*group})
	if err != nil {
		return err
	}

	// Check if we have our user authorized
	if len(groupsFiltered) < 1 {
		return &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with external ID %v is not allowed to access to resource %v",
				authenticatedUser.Identifier, group.Urn),
		}
	}

	// Check if policy exist
	policy, err := api.GetPolicy(authenticatedUser, org, policyName)
	if err != nil {
		return err
	}

	// Check if exist this relation
	_, err = api.GroupRepo.GetGroupPolicyRelation(group.ID, policy.ID)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		// Relation doesn't exist in DB
		switch dbError.Code {
		case database.GROUP_POLICY_RELATION_NOT_FOUND:
			return &Error{
				Code:    POLICY_IS_NOT_ATTACHED_TO_GROUP,
				Message: dbError.Message,
			}
		default: // Unexpected error
			return &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
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

	return nil
}

func (api *AuthAPI) ListAttachedGroupPolicies(authenticatedUser AuthenticatedUser, org string, groupName string) (*GroupPolicies, error) {
	// Check if group exist
	group, err := api.GetGroupByName(authenticatedUser, org, groupName)
	if err != nil {
		return nil, err
	}

	// Check restrictions
	groupsFiltered, err := api.GetGroupsAuthorized(authenticatedUser, group.Urn, GROUP_ACTION_LIST_ATTACHED_GROUP_POLICIES, []Group{*group})
	if err != nil {
		return nil, err
	}

	// Check if we have our user authorized
	if len(groupsFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with external ID %v is not allowed to access to resource %v",
				authenticatedUser.Identifier, group.Urn),
		}
	}

	// Call repo to retrieve the GroupPolicyRelations
	groupPolicies, err := api.GroupRepo.GetAllGroupPolicyRelation(group.ID)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	// Return policies
	return groupPolicies, nil
}

// Private helper methods

func createGroup(org string, name string, path string) Group {
	urn := CreateUrn(org, RESOURCE_GROUP, path, name)
	group := Group{
		ID:   uuid.NewV4().String(),
		Name: name,
		Path: path,
		Urn:  urn,
		Org:  org,
	}

	return group
}
