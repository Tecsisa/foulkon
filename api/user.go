package api

import (
	"fmt"
	"time"

	"github.com/satori/go.uuid"
	"github.com/tecsisa/authorizr/database"
)

// User domain
type User struct {
	ID         string    `json:"ID, omitempty"`
	ExternalID string    `json:"ExternalID, omitempty"`
	Path       string    `json:"Path, omitempty"`
	CreateAt   time.Time `json:"CreateAt, omitempty"`
	Urn        string    `json:"Urn, omitempty"`
}

// User api
type UsersAPI struct {
	Repo Repo
}

// Retrieve user by external id
func (u *UsersAPI) GetUserByExternalId(id string) (*User, error) {
	// Call repo to retrieve the user
	user, err := u.Repo.UserRepo.GetUserByExternalID(id)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		// User doesn't exist in DB
		if dbError.Code == database.USER_NOT_FOUND {
			return nil, &Error{
				Code:    USER_BY_EXTERNAL_ID_NOT_FOUND,
				Message: dbError.Message,
			}
		} else { // Unexpected error
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	}
	return user, nil

}

// Retrieve user by id
func (u *UsersAPI) GetUserByID(id string) (*User, error) {
	// Call repo to retrieve the user
	user, err := u.Repo.UserRepo.GetUserByID(id)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		// User doesn't exist in DB
		if dbError.Code == database.USER_NOT_FOUND {
			return nil, &Error{
				Code:    USER_BY_ID_NOT_FOUND,
				Message: dbError.Message,
			}
		} else { // Unexpected error
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	}

	// Return user
	return user, nil
}

// Retrieve users that has path
func (u *UsersAPI) GetListUsers(pathPrefix string) ([]User, error) {

	// Retrieve users with specified path prefix
	users, err := u.Repo.UserRepo.GetUsersFiltered(pathPrefix)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	// Return users
	return users, nil
}

// Add an User to database if not exist
func (u *UsersAPI) AddUser(externalID string, path string) (*User, error) {
	// Validate external ID
	if !IsValidUserExternalID(externalID) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: ExternalID %v", externalID),
		}
	}

	// Validate path
	if !IsValidPath(path) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: Path %v", path),
		}
	}

	// Check if user already exist
	_, err := u.Repo.UserRepo.GetUserByExternalID(externalID)

	// Check if user could be retrieved
	if err != nil {
		// Transform to DB error
		dbError := err.(*database.Error)
		// User doesn't exist in DB
		switch dbError.Code {
		case database.USER_NOT_FOUND:
			// Create user
			user := createUser(externalID, path)
			userCreated, err := u.Repo.UserRepo.AddUser(user)

			// Check if there is an unexpected error in DB
			if err != nil {
				//Transform to DB error
				dbError := err.(*database.Error)
				return nil, &Error{
					Code:    UNKNOWN_API_ERROR,
					Message: dbError.Message,
				}
			}

			// Return user created
			return userCreated, nil
		default: // Unexpected error
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	} else {
		return nil, &Error{
			Code:    USER_ALREADY_EXIST,
			Message: fmt.Sprintf("Unable to create user, user with ExternalID %v already exist", externalID),
		}
	}

}

// Update an User to database if exist
func (u *UsersAPI) UpdateUser(externalID string, newPath string) (*User, error) {
	// Validate external ID
	if !IsValidUserExternalID(externalID) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: ExternalID %v", externalID),
		}
	}

	// Validate path
	if !IsValidPath(newPath) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: Path %v", newPath),
		}
	}

	// Call repo to retrieve the user
	userDB, err := u.Repo.UserRepo.GetUserByExternalID(externalID)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		// User doesn't exist in DB
		if dbError.Code == database.USER_NOT_FOUND {
			return nil, &Error{
				Code:    USER_BY_EXTERNAL_ID_NOT_FOUND,
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
	urn := CreateUrn("", RESOURCE_USER, newPath, externalID)

	// Update user
	user, err := u.Repo.UserRepo.UpdateUser(*userDB, newPath, urn)

	// Check if there is an unexpected error in DB
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	return user, nil

}

// Remove user with this id
func (u *UsersAPI) RemoveUserById(id string) error {
	// Call repo to retrieve the user
	user, err := u.Repo.UserRepo.GetUserByExternalID(id)

	if err != nil {
		//Transform to DB error
		dbError := err.(database.Error)
		// If user doesn't exist
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

	// Remove user with given id
	err = u.Repo.UserRepo.RemoveUser(user.ID)

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

// Get groups for an user
func (u *UsersAPI) GetGroupsByUserId(id string) ([]Group, error) {
	return u.Repo.UserRepo.GetGroupsByUserID(id)
}

// Private helper methods

func createUser(externalID string, path string) User {
	urn := CreateUrn("", RESOURCE_USER, path, externalID)
	user := User{
		ID:         uuid.NewV4().String(),
		ExternalID: externalID,
		Path:       path,
		Urn:        urn,
	}

	return user
}
