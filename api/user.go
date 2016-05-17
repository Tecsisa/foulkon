package api

import (
	"fmt"
	"time"

	"strings"

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
	UserRepo UserRepo
}

// Retrieve user by external id
func (u *UsersAPI) GetUserByExternalId(id string) (*User, error) {
	// Call repo to retrieve the user
	user, err := u.UserRepo.GetUserByExternalID(id)

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

	// Return user
	return user, nil
}

// Retrieve user by id
func (u *UsersAPI) GetUserByID(id string) (*User, error) {
	// Call repo to retrieve the user
	user, err := u.UserRepo.GetUserByID(id)

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
	users, err := u.UserRepo.GetUsersFiltered(pathPrefix)

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
	// Check parameters
	if len(strings.TrimSpace(externalID)) == 0 ||
		len(strings.TrimSpace(path)) == 0 {
		return nil, &Error{
			Code:    MISSING_PARAMETER_ERROR,
			Message: fmt.Sprintf("There are mising parameters: ExternalID %v, Path %v", externalID, path),
		}
	}

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
	userDB, err := u.UserRepo.GetUserByExternalID(externalID)

	// If user exist it can't create it
	if userDB != nil {
		return nil, &Error{
			Code:    USER_ALREADY_EXIST,
			Message: fmt.Sprintf("Unable to create user, user with ExternalID %v already exist", externalID),
		}
	}

	// Check if there is an unexpected error in DB
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		if dbError.Code != database.USER_NOT_FOUND {
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	}

	// Create user
	user := createUser(externalID, path)
	userCreated, err := u.UserRepo.AddUser(user)

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
}

// Update an User to database if exist
func (u *UsersAPI) UpdateUser(externalID string, newPath string) (*User, error) {
	// Check parameters
	if len(strings.TrimSpace(externalID)) == 0 ||
		len(strings.TrimSpace(newPath)) == 0 {
		return nil, &Error{
			Code:    MISSING_PARAMETER_ERROR,
			Message: fmt.Sprintf("There are mising parameters: ExternalID %v, Path %v", externalID, newPath),
		}
	}

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
	userDB, err := u.UserRepo.GetUserByExternalID(externalID)

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
	user, err := u.UserRepo.UpdateUser(*userDB, newPath, urn)

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
	// Remove user with given external id
	err := u.UserRepo.RemoveUser(id)

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

	return nil
}

// Get groups for an user
func (u *UsersAPI) GetGroupsByUserId(id string) ([]Group, error) {
	return u.UserRepo.GetGroupsByUserID(id)
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
