package postgresql

import (
	"fmt"
	"time"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/database"
)

// USER REPOSITORY IMPLEMENTATION

func (u PostgresRepo) AddUser(user api.User) (*api.User, error) {

	// Create user model
	userDB := &User{
		ID:         user.ID,
		ExternalID: user.ExternalID,
		Path:       user.Path,
		CreateAt:   user.CreateAt.UnixNano(),
		UpdateAt:   user.UpdateAt.UnixNano(),
		Urn:        user.Urn,
	}

	// Store user
	err := u.Dbmap.Create(userDB).Error

	// Error handling
	if err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	return dbUserToAPIUser(userDB), nil
}

func (u PostgresRepo) GetUserByExternalID(id string) (*api.User, error) {
	user := &User{}
	query := u.Dbmap.Where("external_id like ?", id).First(user)

	// Check if user exists
	if query.RecordNotFound() {
		return nil, &database.Error{
			Code:    database.USER_NOT_FOUND,
			Message: fmt.Sprintf("User with externalId %v not found", id),
		}
	}

	// Error Handling
	if err := query.Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	return dbUserToAPIUser(user), nil
}

func (u PostgresRepo) GetUserByID(id string) (*api.User, error) {
	user := &User{}
	query := u.Dbmap.Where("id like ?", id).First(user)

	// Check if user exists
	if query.RecordNotFound() {
		return nil, &database.Error{
			Code:    database.USER_NOT_FOUND,
			Message: fmt.Sprintf("User with id %v not found", id),
		}
	}

	// Error Handling
	if err := query.Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	return dbUserToAPIUser(user), nil
}

func (u PostgresRepo) GetUsersFiltered(filter *api.Filter) ([]api.User, int, error) {
	var total int
	users := []User{}
	query := u.Dbmap

	// Check if path is filled, else it doesn't use it to filter
	if len(filter.PathPrefix) > 0 {
		query = query.Where("path like ?", filter.PathPrefix+"%")
	}

	// Error handling
	if err := query.Find(&users).Count(&total).Offset(filter.Offset).Limit(filter.Limit).Find(&users).Error; err != nil {
		return nil, total, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Transform users for API
	var apiusers []api.User
	if users != nil {
		apiusers = make([]api.User, len(users), cap(users))
		for i, u := range users {
			apiusers[i] = *dbUserToAPIUser(&u)
		}
	}

	return apiusers, total, nil
}

func (u PostgresRepo) UpdateUser(user api.User) (*api.User, error) {

	userDB := User{
		ID:         user.ID,
		ExternalID: user.ExternalID,
		Path:       user.Path,
		CreateAt:   user.CreateAt.UnixNano(),
		UpdateAt:   user.UpdateAt.UnixNano(),
		Urn:        user.Urn,
	}

	// Update user
	query := u.Dbmap.Model(&User{ID: user.ID}).Updates(userDB)

	// Error Handling
	if err := query.Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	return &user, nil
}

func (u PostgresRepo) RemoveUser(id string) error {
	transaction := u.Dbmap.Begin()
	// Delete user
	transaction.Where("id like ?", id).Delete(&User{})

	// Error handling
	if err := transaction.Error; err != nil {
		transaction.Rollback()
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	//  delete all user relations
	transaction.Where("user_id like ?", id).Delete(&GroupUserRelation{})

	// Error handling
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

func (u PostgresRepo) GetGroupsByUserID(id string, filter *api.Filter) ([]api.Group, int, error) {
	var total int
	relations := []GroupUserRelation{}
	query := u.Dbmap.Where("user_id like ?", id).Find(&relations).Count(&total).Offset(filter.Offset).Limit(filter.Limit).Find(&relations)

	// Error Handling
	if err := query.Error; err != nil {
		return nil, total, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	var apiGroups []api.Group
	// Transform relations to API domain
	if relations != nil {
		apiGroups = make([]api.Group, len(relations), cap(relations))
		for i, r := range relations {
			group, err := u.GetGroupById(r.GroupID)
			// Error handling
			if err != nil {
				return nil, total, &database.Error{
					Code:    database.INTERNAL_ERROR,
					Message: err.Error(),
				}
			}
			apiGroups[i] = *group
		}
	}

	return apiGroups, total, nil
}

// PRIVATE HELPER METHODS

// Transform a user retrieved from db into a user for API
func dbUserToAPIUser(userdb *User) *api.User {
	return &api.User{
		ID:         userdb.ID,
		ExternalID: userdb.ExternalID,
		Path:       userdb.Path,
		CreateAt:   time.Unix(0, userdb.CreateAt).UTC(),
		UpdateAt:   time.Unix(0, userdb.UpdateAt).UTC(),
		Urn:        userdb.Urn,
	}
}
