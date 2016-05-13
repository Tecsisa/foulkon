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

func (g PostgresRepo) UpdateGroup(group api.Group, newName string, newPath string, urn string) (*api.Group, error) {

	// Create new group
	groupUpdated := Group{
		Name: newName,
		Path: newPath,
		Urn:  urn,
	}

	groupDB := Group{
		ID:       group.ID,
		Name:     group.Name,
		Path:     group.Path,
		CreateAt: group.CreateAt.UTC().UnixNano(),
		Urn:      group.Urn,
		Org:      group.Org,
	}

	// Update group
	query := g.Dbmap.Model(&groupDB).Update(groupUpdated)

	// Check if group exist
	if query.RecordNotFound() {
		return nil, &database.Error{
			Code:    database.GROUP_NOT_FOUND,
			Message: fmt.Sprintf("Group with name %v not found", group.Name),
		}
	}

	// Error Handling
	if err := query.Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	return groupDBToGroupAPI(&groupDB), nil
}

func (g PostgresRepo) GetGroupsFiltered(org string, pathPrefix string) ([]api.Group, error) {
	groups := []Group{}
	query := g.Dbmap
	if len(org) > 0 {
		query = g.Dbmap.Where("org like ? ", org)
	}
	if len(pathPrefix) > 0 {
		query = g.Dbmap.Where("path like ? ", pathPrefix+"%")
	}
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

func (g PostgresRepo) GetGroupUserRelation(userID string, groupID string) (*api.GroupMembers, error) {
	relation := GroupUserRelation{}
	query := g.Dbmap.Where("user_id like ? AND group_id like ?", userID, groupID).First(&relation)

	// Check if relation exist
	if query.RecordNotFound() {
		return nil, &database.Error{
			Code:    database.GROUP_USER_RELATION_NOT_FOUND,
			Message: fmt.Sprintf("Relation doesn't exist with UserID %v and groupID %v", userID, groupID),
		}
	}

	// Error Handling
	if err := query.Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Retrieve group
	group, err := g.GetGroupById(groupID)
	// Error Handling
	if err != nil {
		return nil, err
	}

	// Retrieve user
	user, err := g.GetUserByID(userID)
	// Error Handling
	if err != nil {
		return nil, err
	}

	// Return GroupMembers
	return &api.GroupMembers{
		Group: *group,
		Users: []api.User{*user},
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

func (g PostgresRepo) GetGroupPolicyRelation(groupID string, policyID string) (*api.GroupPolicies, error) {
	relation := GroupPolicyRelation{}
	query := g.Dbmap.Where("group_id like ? AND policy_id like ?", groupID, policyID).First(&relation)

	// Check if relation exist
	if query.RecordNotFound() {
		return nil, &database.Error{
			Code:    database.GROUP_POLICY_RELATION_NOT_FOUND,
			Message: fmt.Sprintf("Relation doesn't exist with GroupID %v and PolicyID %v", groupID, policyID),
		}
	}

	// Error Handling
	if err := query.Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Retrieve group
	group, err := g.GetGroupById(groupID)
	// Error Handling
	if err != nil {
		return nil, err
	}

	// Retrieve policy
	policy, err := g.GetPolicyById(policyID)
	// Error Handling
	if err != nil {
		return nil, err
	}

	// Return GroupPolicies
	return &api.GroupPolicies{
		Group:    *group,
		Policies: []api.Policy{*policy},
	}, nil
}

func (g PostgresRepo) AttachPolicy(group api.Group, policy api.Policy) error {
	// Create relation
	relation := &GroupPolicyRelation{
		GroupID:  group.ID,
		PolicyID: policy.ID,
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
