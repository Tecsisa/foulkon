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

type GroupMembers struct {
	Group Group  `json:"Group, omitempty"`
	Users []User `json:"Users, omitempty"`
}

type GroupPolicies struct {
	Group    Group    `json:"Group, omitempty"`
	Policies []Policy `json:"Policies, omitempty"`
}

// Add an Group to database if not exist
func (api *AuthAPI) AddGroup(org string, name string, path string) (*Group, error) {
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

	// Check if group already exist
	_, err := api.GroupRepo.GetGroupByName(org, name)

	// Check if group could be retrieved
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		switch dbError.Code {
		// Group doesn't exist in DB
		case database.GROUP_NOT_FOUND:
			// Create group
			groupCreated, err := api.GroupRepo.AddGroup(createGroup(org, name, path))

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
func (api *AuthAPI) AddMember(userID string, groupName string, org string) error {
	// Call repo to retrieve the group
	groupDB, err := api.getGroupByName(org, groupName)

	// Error handling
	if err != nil {
		return err
	}

	// Call repo to retrieve the user
	userDB, err := api.UserRepo.GetUserByExternalID(userID)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		// User doesn't exist in DB
		switch dbError.Code {
		case database.USER_NOT_FOUND:
			return &Error{
				Code:    USER_BY_EXTERNAL_ID_NOT_FOUND,
				Message: dbError.Message,
			}
		default: // Unexpected error
			return &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}

		}
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
func (api *AuthAPI) RemoveMember(userID string, groupName string, org string) error {
	// Call repo to retrieve the group
	groupDB, err := api.getGroupByName(org, groupName)

	// Error handling
	if err != nil {
		return err
	}

	// Call repo to retrieve the user
	userDB, err := api.UserRepo.GetUserByExternalID(userID)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		// User doesn't exist in DB
		switch dbError.Code {
		case database.USER_NOT_FOUND:
			return &Error{
				Code:    USER_BY_EXTERNAL_ID_NOT_FOUND,
				Message: dbError.Message,
			}
		default: // Unexpected error
			return &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
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
func (api *AuthAPI) ListMembers(org string, groupName string) (*GroupMembers, error) {
	// Call repo to retrieve the group
	group, err := api.getGroupByName(org, groupName)

	// Error handling
	if err != nil {
		return nil, err
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
func (api *AuthAPI) RemoveGroup(org string, name string) error {
	// Call repo to retrieve the group
	group, err := api.getGroupByName(org, name)

	// Error handling
	if err != nil {
		return err
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

func (api *AuthAPI) GetGroupByName(org string, name string) (*Group, error) {
	// Call repo to retrieve the group
	group, err := api.getGroupByName(org, name)

	// Error handling
	if err != nil {
		return nil, err
	}

	// Return group
	return group, nil

}

func (api *AuthAPI) GetGroupById(id string) (*Group, error) {
	// Call repo to retrieve the group
	group, err := api.GroupRepo.GetGroupById(id)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		// Group doesn't exist in DB
		switch dbError.Code {
		case database.GROUP_NOT_FOUND:
			return nil, &Error{
				Code:    GROUP_BY_ID_NOT_FOUND,
				Message: dbError.Message,
			}
		default: // Unexpected error
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}

		}
	}

	// Return group
	return group, nil

}

func (api *AuthAPI) GetListGroups(org string, pathPrefix string) ([]Group, error) {
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

	// Return groups
	return groups, nil
}

// Update Group to database if exist
func (api *AuthAPI) UpdateGroup(org string, groupName string, newName string, newPath string) (*Group, error) {
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
	groupDB, err := api.getGroupByName(org, groupName)

	// Error handling
	if err != nil {
		return nil, err
	}

	// Check if group with newName exist
	_, err = api.getGroupByName(org, newName)

	if err == nil {
		// Group already exists
		return nil, &Error{
			Code:    GROUP_ALREADY_EXIST,
			Message: fmt.Sprintf("Group name: %v already exists", newName),
		}
	}

	// Get Urn
	urn := CreateUrn(org, RESOURCE_GROUP, newPath, newName)

	// Update group
	group, err := api.GroupRepo.UpdateGroup(*groupDB, newName, newPath, urn)

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

func (api *AuthAPI) AttachPolicyToGroup(org string, groupName string, policyName string) error {
	// Check if group exist
	group, err := api.getGroupByName(org, groupName)

	// Error handling
	if err != nil {
		return err
	}

	// Check if policy exist
	policy, err := api.PolicyRepo.GetPolicyByName(org, policyName)

	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		switch dbError.Code {
		case database.POLICY_NOT_FOUND:
			return &Error{
				Code:    POLICY_BY_ORG_AND_NAME_NOT_FOUND,
				Message: dbError.Message,
			}
		default: // Unexpected error
			return &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
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

func (api *AuthAPI) DetachPolicyToGroup(org string, groupName string, policyName string) error {
	// Check if group exist
	group, err := api.getGroupByName(org, groupName)

	// Error handling
	if err != nil {
		return err
	}

	// Check if policy exist
	policy, err := api.PolicyRepo.GetPolicyByName(org, policyName)

	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		switch dbError.Code {
		case database.POLICY_NOT_FOUND:
			return &Error{
				Code:    POLICY_BY_ORG_AND_NAME_NOT_FOUND,
				Message: dbError.Message,
			}
		default: // Unexpected error
			return &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
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

func (api *AuthAPI) ListAttachedGroupPolicies(org string, groupName string) (*GroupPolicies, error) {
	// Check if group exist
	group, err := api.getGroupByName(org, groupName)

	// Error handling
	if err != nil {
		return nil, err
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

// This method gets the group by name and organization
func (api *AuthAPI) getGroupByName(org string, name string) (*Group, error) {
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

	// Return group
	return group, nil

}
