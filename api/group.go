package api

import (
	"fmt"
	"github.com/tecsisa/authorizr/database"
	"time"
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

type GroupsAPI struct {
	GroupRepo GroupRepo
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
