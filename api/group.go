package api

import (
	"fmt"
	"time"

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

type GroupsAPI struct {
	GroupRepo  GroupRepo
	UserRepo   UserRepo
	PolicyRepo PolicyRepo
}

// Add an Group to database if not exist
func (g *GroupsAPI) AddGroup(group Group) (*Group, error) {
	// Check if group already exist
	groupDB, err := g.GroupRepo.GetGroupByName(group.Org, group.Name)

	// If group exist it can't create it
	if groupDB != nil {
		return nil, &Error{
			Code:    GROUP_ALREADY_EXIST,
			Message: fmt.Sprintf("Unable to create group, group with org %v and name %v already exist", groupDB.Org, groupDB.Name),
		}
	}
	// Create group
	groupCreated, err := g.GroupRepo.AddGroup(group)

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
}

//
func (g *GroupsAPI) AddMember(userID string, groupName string, org string) error {
	// Call repo to retrieve the group
	groupDB, err := g.GroupRepo.GetGroupByName(org, groupName)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		// Group doesn't exist in DB
		if dbError.Code == database.GROUP_NOT_FOUND {
			return &Error{
				Code:    GROUP_BY_ORG_AND_NAME_NOT_FOUND,
				Message: dbError.Message,
			}
		} else { // Unexpected error
			return &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	}

	// Call repo to retrieve the user
	userDB, err := g.UserRepo.GetUserByExternalID(userID)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		// User doesn't exist in DB
		if dbError.Code == database.USER_NOT_FOUND {
			return &Error{
				Code:    USER_BY_EXTERNAL_ID_NOT_FOUND,
				Message: dbError.Message,
			}
		} else { // Unexpected error
			return &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	}

	// Call repo to retrieve the GroupUserRelation
	groupMembers, err := g.GroupRepo.GetGroupUserRelation(userDB.ID, groupDB.ID)

	// Error handling
	if groupMembers != nil {
		return &Error{
			Code:    USER_IS_ALREADY_A_MEMBER_OF_GROUP,
			Message: fmt.Sprintf("User: %v is already a member of Group: %v", userID, groupName),
		}
	}

	// Add Member
	err = g.GroupRepo.AddMember(*userDB, *groupDB)

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

// Remove group
func (g *GroupsAPI) RemoveGroup(org string, name string) error {
	// Remove group with given org and name
	err := g.GroupRepo.RemoveGroup(org, name)

	if err != nil {
		//Transform to DB error
		dbError := err.(database.Error)
		// If group doesn't exist
		if dbError.Code == database.GROUP_NOT_FOUND {
			return &Error{
				Code:    GROUP_BY_ORG_AND_NAME_NOT_FOUND,
				Message: dbError.Message,
			}
		} else { // Unexpected error
			return &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	}

	return nil
}

func (g *GroupsAPI) GetGroupByName(org string, name string) (*Group, error) {
	// Call repo to retrieve the group
	group, err := g.GroupRepo.GetGroupByName(org, name)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		// Group doesn't exist in DB
		if dbError.Code == database.GROUP_NOT_FOUND {
			return nil, &Error{
				Code:    GROUP_BY_ORG_AND_NAME_NOT_FOUND,
				Message: dbError.Message,
			}
		} else { // Unexpected error
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	}

	// Return group
	return group, nil

}

func (g *GroupsAPI) GetGroupById(id string) (*Group, error) {
	// Call repo to retrieve the group
	group, err := g.GroupRepo.GetGroupById(id)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		// Group doesn't exist in DB
		if dbError.Code == database.GROUP_NOT_FOUND {
			return nil, &Error{
				Code:    GROUP_BY_ID_NOT_FOUND,
				Message: dbError.Message,
			}
		} else { // Unexpected error
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	}

	// Return group
	return group, nil

}

func (g *GroupsAPI) GetListGroups(org string, pathPrefix string) ([]Group, error) {
	// Call repo to retrieve the groups
	groups, err := g.GroupRepo.GetGroupsFiltered(org, pathPrefix)

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

func (g *GroupsAPI) UpdateGroup(org string, groupName string, newName string, newPath string) (*Group, error) {
	// Call repo to retrieve the group
	groupDB, err := g.GroupRepo.GetGroupByName(org, groupName)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		// Group doesn't exist in DB
		if dbError.Code == database.GROUP_NOT_FOUND {
			return nil, &Error{
				Code:    GROUP_BY_ORG_AND_NAME_NOT_FOUND,
				Message: dbError.Message,
			}
		} else { // Unexpected error
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	}

	// Get Urn
	urn := CreateUrn(org, RESOURCE_GROUP, newPath, newName)

	// Update group
	group, err := g.GroupRepo.UpdateGroup(*groupDB, newName, newPath, urn)

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

func (g *GroupsAPI) AttachPolicyToGroup(org string, groupName string, policyName string) error {
	// Check if group exist
	group, err := g.GetGroupByName(org, groupName)

	// Error handling
	if err != nil {
		return err
	}

	// Check if policy exist
	policy, err := g.PolicyRepo.GetPolicyByName(org, policyName)

	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		if dbError.Code == database.POLICY_NOT_FOUND {
			return &Error{
				Code:    POLICY_BY_ORG_AND_NAME_NOT_FOUND,
				Message: dbError.Message,
			}
		} else { // Unexpected error
			return &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	}

	// Check if exist this relation
	groupPolicies, err := g.GroupRepo.GetGroupPolicyRelation(group.ID, policy.ID)

	if groupPolicies != nil {
		// Unexpected error
		return &Error{
			Code:    POLICY_IS_ALREADY_ATTACHED_TO_GROUP,
			Message: fmt.Sprintf("Policy: %v is already attached to Group: %v", policy.ID, group.ID),
		}
	}

	// Attach Policy to Group
	err = g.GroupRepo.AttachPolicy(*group, *policy)

	if err != nil {
		dbError := err.(*database.Error)
		return &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	return nil
}
