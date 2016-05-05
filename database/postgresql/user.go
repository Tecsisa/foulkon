package postgresql

import (
	"fmt"
	"time"

	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/database"
)

func (u PostgresRepo) GetUserByExternalID(id string) (*api.User, error) {
	user := &User{}
	query := u.Dbmap.Where("external_id like ?", id).First(user)

	// Check if user exist
	if query.RecordNotFound() {
		return nil, &database.Error{
			Code:    database.USER_NOT_FOUND,
			Message: fmt.Sprintf("User with ExternalID %v not found", id),
		}
	}

	// Error Handling
	if err := query.Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Return user
	return userDBToUserAPI(user), nil
}

func (u PostgresRepo) GetUserByID(id string) (*api.User, error) {
	user := &User{}
	query := u.Dbmap.Where("id like ?", id).First(user)

	// Check if user exist
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

	// Return user
	return userDBToUserAPI(user), nil
}

func (u PostgresRepo) AddUser(user api.User) (*api.User, error) {

	// Create user model
	userDB := &User{
		ID:         user.ID,
		ExternalID: user.ExternalID,
		Path:       user.Path,
		CreateAt:   time.Now().UTC().UnixNano(),
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

	return userDBToUserAPI(userDB), nil
}

func (u PostgresRepo) GetUsersFiltered(pathPrefix string) ([]api.User, error) {
	users := []User{}
	query := u.Dbmap

	// Check if path is filled, else it doesn't use it to filter
	if len(pathPrefix) > 0 {
		query = query.Where("path like ?", pathPrefix+"%")
	}

	// Error handling
	if err := query.Find(&users).Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Transform users for API
	if users != nil {
		apiusers := make([]api.User, len(users), cap(users))
		for i, u := range users {
			apiusers[i] = *userDBToUserAPI(&u)
		}
		return apiusers, nil
	}

	// No data to return
	return nil, nil
}

func (u PostgresRepo) GetGroupsByUserID(id string) ([]api.Group, error) {
	return nil, nil
}

func (u PostgresRepo) RemoveUser(id string) error {
	// Retrieve user with this external id
	user, err := u.GetUserByExternalID(id)

	// Go to delete user
	if user != nil {
		err = u.Dbmap.Delete(&user).Error
		// Error handling
		if err != nil {
			return database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: err.Error(),
			}
		}
		return nil
	}

	// Return error if user isn't found
	return err
}

// Transform a user retrieved from db into a user for API
func userDBToUserAPI(userdb *User) *api.User {
	return &api.User{
		ID:         userdb.ID,
		ExternalID: userdb.ExternalID,
		Path:       userdb.Path,
		CreateAt:   time.Unix(0, userdb.CreateAt).UTC(),
		Urn:        userdb.Urn,
	}
}
