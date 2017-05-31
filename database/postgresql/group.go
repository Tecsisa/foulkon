package postgresql

import (
	"fmt"
	"time"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/database"
)

// GROUP REPOSITORY IMPLEMENTATION

func (pr PostgresRepo) AddGroup(group api.Group) (*api.Group, error) {
	// Create group model
	groupDB := &Group{
		ID:       group.ID,
		Name:     group.Name,
		Path:     group.Path,
		CreateAt: group.CreateAt.UnixNano(),
		UpdateAt: group.UpdateAt.UnixNano(),
		Urn:      group.Urn,
		Org:      group.Org,
	}

	// Store group
	err := pr.Dbmap.Create(groupDB).Error

	// Error handling
	if err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	return dbGroupToAPIGroup(groupDB), nil
}

func (pr PostgresRepo) GetGroupByName(org string, name string) (*api.Group, error) {
	group := &Group{}
	query := pr.Dbmap.Where("org like ? AND name like ?", org, name).First(group)

	// Check if group exists
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

	return dbGroupToAPIGroup(group), nil
}

func (pr PostgresRepo) GetGroupById(id string) (*api.Group, error) {
	group := &Group{}
	query := pr.Dbmap.Where("id like ?", id).First(group)

	// Check if group exists
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

	return dbGroupToAPIGroup(group), nil
}

func (pr PostgresRepo) GetGroupsFiltered(filter *api.Filter) ([]api.Group, int, error) {
	var total int
	groups := []Group{}
	query := pr.Dbmap

	if len(filter.Org) > 0 {
		query = query.Where("org like ? ", filter.Org)
	}
	if len(filter.PathPrefix) > 0 {
		query = query.Where("path like ? ", filter.PathPrefix+"%")
	}
	if len(filter.OrderBy) > 0 {
		query = query.Order(filter.OrderBy)
	}

	// Error handling
	if err := query.Find(&groups).Count(&total).Offset(filter.Offset).Limit(filter.Limit).Find(&groups).Error; err != nil {
		return nil, total, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Transform users for API
	var apiGroups []api.Group
	if groups != nil {
		apiGroups = make([]api.Group, len(groups), cap(groups))
		for i, g := range groups {
			apiGroups[i] = *dbGroupToAPIGroup(&g)
		}
	}

	// No data to return
	return apiGroups, total, nil
}

func (pr PostgresRepo) UpdateGroup(group api.Group) (*api.Group, error) {
	groupDB := Group{
		ID:       group.ID,
		Name:     group.Name,
		Path:     group.Path,
		CreateAt: group.CreateAt.UTC().UnixNano(),
		UpdateAt: group.UpdateAt.UTC().UnixNano(),
		Urn:      group.Urn,
		Org:      group.Org,
	}

	// Update group
	query := pr.Dbmap.Model(&Group{ID: group.ID}).Updates(groupDB)

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

	return &group, nil
}

func (pr PostgresRepo) RemoveGroup(id string) error {
	transaction := pr.Dbmap.Begin()

	// Delete group
	transaction.Where("id like ?", id).Delete(&Group{})
	if err := transaction.Error; err != nil {
		transaction.Rollback()
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Delete all group relations
	transaction.Where("group_id like ?", id).Delete(&GroupUserRelation{})
	if err := transaction.Error; err != nil {
		transaction.Rollback()
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}

	}
	// Delete all policy relations
	transaction.Where("group_id like ?", id).Delete(&GroupPolicyRelation{})
	if err := transaction.Error; err != nil {
		transaction.Rollback()
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	transaction.Commit()
	return nil
}

func (pr PostgresRepo) AddMember(userID string, groupID string) error {
	// Create relation
	relation := &GroupUserRelation{
		UserID:   userID,
		GroupID:  groupID,
		CreateAt: time.Now().UTC().UnixNano(),
	}

	// Store relation
	err := pr.Dbmap.Create(relation).Error

	// Error handling
	if err != nil {
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	return nil
}

func (pr PostgresRepo) RemoveMember(userID string, groupID string) error {
	err := pr.Dbmap.Where("user_id like ? AND group_id like ?", userID, groupID).Delete(&GroupUserRelation{}).Error

	// Error handling
	if err != nil {
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}
	return nil
}

func (pr PostgresRepo) IsMemberOfGroup(userID string, groupID string) (bool, error) {
	relation := GroupUserRelation{}
	query := pr.Dbmap.Where("user_id like ? AND group_id like ?", userID, groupID).First(&relation)

	// Check if relation exists
	if query.RecordNotFound() {
		return false, nil
	}

	// Error Handling
	if err := query.Error; err != nil {
		return false, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	return true, nil
}

func (pr PostgresRepo) GetGroupMembers(groupID string, filter *api.Filter) ([]api.UserGroupRelation, int, error) {
	var total int
	members := []GroupUserRelation{}
	query := pr.Dbmap.Where("group_id like ?", groupID)

	if len(filter.OrderBy) > 0 {
		query = query.Order(filter.OrderBy)
	}

	// Error handling
	if err := query.Find(&members).Count(&total).Offset(filter.Offset).Limit(filter.Limit).Find(&members).Error; err != nil {
		return nil, total, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	var membersList []api.UserGroupRelation
	// Transform relations to API domain
	if members != nil {
		membersList = make([]api.UserGroupRelation, len(members), cap(members))
		for i, m := range members {
			user, err := pr.GetUserByID(m.UserID)

			// Error handling
			if err != nil {
				return nil, total, &database.Error{
					Code:    database.INTERNAL_ERROR,
					Message: err.Error(),
				}
			}

			membersList[i] = &GroupUser{
				User:     user,
				CreateAt: time.Unix(0, m.CreateAt).UTC(),
			}
		}
	}

	return membersList, total, nil
}

func (pr PostgresRepo) AttachPolicy(groupID string, policyID string) error {
	// Create relation
	relation := &GroupPolicyRelation{
		GroupID:  groupID,
		PolicyID: policyID,
		CreateAt: time.Now().UTC().UnixNano(),
	}

	// Store relation
	err := pr.Dbmap.Create(relation).Error

	// Error handling
	if err != nil {
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	return nil
}

func (pr PostgresRepo) DetachPolicy(groupID string, policyID string) error {
	// Remove relation
	err := pr.Dbmap.Where("group_id like ? AND policy_id like ?", groupID, policyID).Delete(&GroupPolicyRelation{}).Error

	// Error handling
	if err != nil {
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	return nil
}

func (pr PostgresRepo) IsAttachedToGroup(groupID string, policyID string) (bool, error) {
	relation := GroupPolicyRelation{}
	query := pr.Dbmap.Where("group_id like ? AND policy_id like ?", groupID, policyID).First(&relation)

	// Check if relation exists
	if query.RecordNotFound() {
		return false, nil
	}

	// Error Handling
	if err := query.Error; err != nil {
		return false, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	return true, nil
}

func (pr PostgresRepo) GetAttachedPolicies(groupID string, filter *api.Filter) ([]api.PolicyGroupRelation, int, error) {
	var total int
	relations := []GroupPolicyRelation{}
	query := pr.Dbmap.Where("group_id like ?", groupID)

	if len(filter.OrderBy) > 0 {
		query = query.Order(filter.OrderBy)
	}

	// Error Handling
	if err := query.Find(&relations).Count(&total).Offset(filter.Offset).Limit(filter.Limit).Find(&relations).Error; err != nil {
		return nil, total, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}
	var policies []api.PolicyGroupRelation
	// Transform relations to API domain
	if relations != nil {
		policies = make([]api.PolicyGroupRelation, len(relations), cap(relations))
		for i, r := range relations {
			policy, err := pr.GetPolicyById(r.PolicyID)
			// Error handling
			if err != nil {
				return nil, total, &database.Error{
					Code:    database.INTERNAL_ERROR,
					Message: err.Error(),
				}
			}

			policies[i] = &PolicyGroup{
				Policy:   policy,
				CreateAt: time.Unix(0, r.CreateAt).UTC(),
			}
		}
	}

	return policies, total, nil
}

// PRIVATE HELPER METHODS

// Transform a Group retrieved from db into a group for API
func dbGroupToAPIGroup(groupdb *Group) *api.Group {
	return &api.Group{
		ID:       groupdb.ID,
		Name:     groupdb.Name,
		Path:     groupdb.Path,
		CreateAt: time.Unix(0, groupdb.CreateAt).UTC(),
		UpdateAt: time.Unix(0, groupdb.UpdateAt).UTC(),
		Urn:      groupdb.Urn,
		Org:      groupdb.Org,
	}
}
