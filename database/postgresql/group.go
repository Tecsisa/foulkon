package postgresql

import (
	"fmt"
	"time"

	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/database"
)

func (g PostgresRepo) GetGroupByName(org string, name string) (*api.Group, error) {
	group := &Group{}
	query := g.Dbmap.Where("org like ? AND name like ?", org, name).First(group)

	// Check if group exist
	if query.RecordNotFound() {
		return nil, &database.Error{
			Code:    database.GROUP_NOT_FOUND,
			Message: fmt.Sprintf("Group with organization %v and name %v not found", org, name),
		}
	}

	// Error Handling
	if err := query.Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Return group
	return groupDBToGroupAPI(group), nil
}

func (g PostgresRepo) GetGroupById(id string) (*api.Group, error) {
	group := &Group{}
	query := g.Dbmap.Where("id like ?", id).First(group)

	// Check if group exist
	if query.RecordNotFound() {
		return nil, &database.Error{
			Code:    database.GROUP_NOT_FOUND,
			Message: fmt.Sprintf("Group with id %v not found", id),
		}
	}

	// Error Handling
	if err := query.Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Return group
	return groupDBToGroupAPI(group), nil
}

func (g PostgresRepo) AddGroup(group api.Group) (*api.Group, error) {

	// Create group model
	groupDB := &Group{
		ID:       group.ID,
		Name:     group.Name,
		Path:     group.Path,
		CreateAt: time.Now().UTC().UnixNano(),
		Urn:      group.Urn,
		Org:      group.Org,
	}

	// Store group
	err := g.Dbmap.Create(groupDB).Error

	// Error handling
	if err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	return groupDBToGroupAPI(groupDB), nil
}

func (g PostgresRepo) AddMember(user api.User, group api.Group) error {

	// Create relation
	relation := &GroupUserRelation{
		UserID:  user.ID,
		GroupID: group.ID,
	}

	// Store relation
	err := g.Dbmap.Create(relation).Error

	// Error handling
	if err != nil {
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	return nil
}

func (g PostgresRepo) GetListGroups(org string) ([]api.Group, error) {
	groups := []Group{}
	query := g.Dbmap.Where("org like ?", org)

	// Error handling
	if err := query.Find(&groups).Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Transform users for API
	if groups != nil {
		apigroups := make([]api.Group, len(groups), cap(groups))
		for i, g := range groups {
			apigroups[i] = *groupDBToGroupAPI(&g)
		}
		return apigroups, nil
	}

	// No data to return
	return nil, nil
}

func (g PostgresRepo) GetGroupUserRelation(user api.User, group api.Group) (*api.GroupMembers, error) {
	relation := GroupUserRelation{}
	query := g.Dbmap.Where("user_id like ? AND group_id like ?", user.ID, group.ID).First(&relation)

	// Check if relation exist
	if query.RecordNotFound() {
		return nil, &database.Error{
			Code:    database.GROUP_USER_RELATION_NOT_FOUND,
			Message: fmt.Sprintf("Relation doesn't exist with UserID %v and groupID %v", user.ID, group.ID),
		}
	}

	// Error Handling
	if err := query.Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Return GroupMembers
	return &api.GroupMembers{
		Group: group,
		Users: []api.User{user},
	}, nil
}

func (g PostgresRepo) RemoveGroup(org string, name string) error {
	// Retrieve group with this org and name
	group, err := g.GetGroupByName(org, name)

	// Go to delete group
	if group != nil {
		err = g.Dbmap.Delete(&group).Error
		// Error handling
		if err != nil {
			return database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: err.Error(),
			}
		}
		return nil
	}

	// Return error if group isn't found
	return err
}

// Transform a Group retrieved from db into a group for API
func groupDBToGroupAPI(groupdb *Group) *api.Group {
	return &api.Group{
		ID:       groupdb.ID,
		Name:     groupdb.Name,
		Path:     groupdb.Path,
		CreateAt: time.Unix(0, groupdb.CreateAt).UTC(),
		Urn:      groupdb.Urn,
		Org:      groupdb.Org,
	}
}
