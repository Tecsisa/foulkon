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

func (u User) GetUrn() string {
	return u.Urn
}

// Retrieve user by external id
func (api AuthAPI) GetUserByExternalId(authenticatedUser AuthenticatedUser, id string) (*User, error) {
	if !IsValidUserExternalID(id) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: ExternalID %v", id),
		}
	}
	// Retrieve user from DB
	user, err := api.UserRepo.GetUserByExternalID(id)

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

	// Check restrictions
	filteredUsers, err := api.GetAuthorizedUsers(authenticatedUser, user.Urn, USER_ACTION_GET_USER, []User{*user})
	if err != nil {
		return nil, err
	}
	if len(filteredUsers) > 0 {
		filteredUser := filteredUsers[0]
		return &filteredUser, nil
	} else {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with external ID %v is not allowed to access to resource %v",
				authenticatedUser.Identifier, user.Urn),
		}
	}

}

func (api AuthAPI) GetUserList(authenticatedUser AuthenticatedUser, pathPrefix string) ([]string, error) {
	// Check parameters
	if len(pathPrefix) > 0 && !IsValidPath(pathPrefix) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: PathPrefix %v", pathPrefix),
		}
	}

	if len(pathPrefix) == 0 {
		pathPrefix = "/"
	}

	// Retrieve users with specified path prefix
	users, err := api.UserRepo.GetUsersFiltered(pathPrefix)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	// Check restrictions
	urnPrefix := GetUrnPrefix("", RESOURCE_USER, pathPrefix)
	usersFiltered, err := api.GetAuthorizedUsers(authenticatedUser, urnPrefix, USER_ACTION_LIST_USERS, users)
	if err != nil {
		return nil, err
	}

	// Return user IDs
	externalIDs := []string{}
	for _, u := range usersFiltered {
		externalIDs = append(externalIDs, u.ExternalID)
	}

	return externalIDs, nil
}

func (api AuthAPI) AddUser(authenticatedUser AuthenticatedUser, externalID string, path string) (*User, error) {
	// Validate fields
	if !IsValidUserExternalID(externalID) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: ExternalID %v", externalID),
		}
	}
	if !IsValidPath(path) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: Path %v", path),
		}
	}

	user := createUser(externalID, path)

	// Check restrictions
	usersFiltered, err := api.GetAuthorizedUsers(authenticatedUser, user.Urn, USER_ACTION_CREATE_USER, []User{user})
	if err != nil {
		return nil, err
	}
	if len(usersFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with external ID %v is not allowed to access to resource %v",
				authenticatedUser.Identifier, user.Urn),
		}
	}

	// Check if user already exists
	_, err = api.UserRepo.GetUserByExternalID(externalID)

	// Check if user could be retrieved
	if err != nil {
		// Transform to DB error
		dbError := err.(*database.Error)
		// User doesn't exist in DB
		switch dbError.Code {
		case database.USER_NOT_FOUND:
			// Create user
			createdUser, err := api.UserRepo.AddUser(user)

			// Check unexpected DB error
			if err != nil {
				//Transform to DB error
				dbError := err.(*database.Error)
				return nil, &Error{
					Code:    UNKNOWN_API_ERROR,
					Message: dbError.Message,
				}
			}

			// Return user created
			return createdUser, nil
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

func (api AuthAPI) UpdateUser(authenticatedUser AuthenticatedUser, externalID string, newPath string) (*User, error) {
	// Validate fields
	if !IsValidUserExternalID(externalID) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: ExternalID %v", externalID),
		}
	}
	if !IsValidPath(newPath) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: Path %v", newPath),
		}
	}

	// Call repo to retrieve the user
	userDB, err := api.GetUserByExternalId(authenticatedUser, externalID)
	if err != nil {
		return nil, err
	}

	// Check restrictions
	usersFiltered, err := api.GetAuthorizedUsers(authenticatedUser, userDB.Urn, USER_ACTION_UPDATE_USER, []User{*userDB})
	if err != nil {
		return nil, err
	}
	if len(usersFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with external ID %v is not allowed to access to resource %v",
				authenticatedUser.Identifier, userDB.Urn),
		}
	}

	userToUpdate := createUser(externalID, newPath)

	// Check restrictions
	usersFiltered, err = api.GetAuthorizedUsers(authenticatedUser, userToUpdate.Urn, USER_ACTION_GET_USER, []User{userToUpdate})
	if err != nil {
		return nil, err
	}
	if len(usersFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with external ID %v is not allowed to access to resource %v",
				authenticatedUser.Identifier, userToUpdate.Urn),
		}
	}

	// Update user
	user, err := api.UserRepo.UpdateUser(*userDB, newPath, userToUpdate.Urn)

	// Check unexpected DB error
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

func (api AuthAPI) RemoveUserById(authenticatedUser AuthenticatedUser, id string) error {
	// Call repo to retrieve the user
	user, err := api.GetUserByExternalId(authenticatedUser, id)
	if err != nil {
		return err
	}

	// Check restrictions
	usersFiltered, err := api.GetAuthorizedUsers(authenticatedUser, user.Urn, USER_ACTION_DELETE_USER, []User{*user})
	if err != nil {
		return err
	}
	if len(usersFiltered) < 1 {
		return &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with external ID %v is not allowed to access to resource %v",
				authenticatedUser.Identifier, user.Urn),
		}
	}

	// Remove user with given id
	err = api.UserRepo.RemoveUser(user.ID)

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

func (api AuthAPI) GetGroupsByUserId(authenticatedUser AuthenticatedUser, id string) ([]GroupIdentity, error) {
	// Call repo to retrieve the user
	user, err := api.GetUserByExternalId(authenticatedUser, id)
	if err != nil {
		return nil, err
	}

	// Check restrictions
	usersFiltered, err := api.GetAuthorizedUsers(authenticatedUser, user.Urn, USER_ACTION_LIST_GROUPS_FOR_USER, []User{*user})
	if err != nil {
		return nil, err
	}
	if len(usersFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with external ID %v is not allowed to access to resource %v",
				authenticatedUser.Identifier, user.Urn),
		}
	}

	// Call group repo to retrieve groups associated to user
	groups, err := api.UserRepo.GetGroupsByUserID(user.ID)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	// Transform to identifiers
	groupIDs := []GroupIdentity{}
	for _, g := range groups {
		groupIDs = append(groupIDs, GroupIdentity{
			Org:  g.Org,
			Name: g.Name,
		})
	}

	return groupIDs, nil
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
