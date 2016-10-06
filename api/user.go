package api

import (
	"fmt"
	"time"

	"github.com/Tecsisa/foulkon/database"
	"github.com/satori/go.uuid"
)

// TYPE DEFINITIONS

// User domain
type User struct {
	ID         string    `json:"id, omitempty"`
	ExternalID string    `json:"externalId, omitempty"`
	Path       string    `json:"path, omitempty"`
	Urn        string    `json:"urn, omitempty"`
	CreateAt   time.Time `json:"createAt, omitempty"`
	UpdateAt   time.Time `json:"updateAt, omitempty"`
}

type UserGroups struct {
	Org      string    `json:"org, omitempty"`
	Name     string    `json:"name, omitempty"`
	CreateAt time.Time `json:"joined, omitempty"`
}

func (u User) String() string {
	return fmt.Sprintf("[id: %v, externalId: %v, path: %v, urn: %v, createAt: %v]",
		u.ID, u.ExternalID, u.Path, u.Urn, u.CreateAt.Format("2006-01-02 15:04:05 MST"))
}

func (u User) GetUrn() string {
	return u.Urn
}

// USER API IMPLEMENTATION

func (api AuthAPI) AddUser(requestInfo RequestInfo, externalId string, path string) (*User, error) {
	// Validate fields
	if !IsValidUserExternalID(externalId) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: externalId %v", externalId),
		}
	}
	if !IsValidPath(path) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: path %v", path),
		}
	}

	user := createUser(externalId, path)

	// Check restrictions
	usersFiltered, err := api.GetAuthorizedUsers(requestInfo, user.Urn, USER_ACTION_CREATE_USER, []User{user})
	if err != nil {
		return nil, err
	}
	if len(usersFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, user.Urn),
		}
	}

	// Check if user already exists
	_, err = api.UserRepo.GetUserByExternalID(externalId)

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
			LogOperation(api.Logger, requestInfo, fmt.Sprintf("User created %+v", createdUser))
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
			Message: fmt.Sprintf("Unable to create user, user with externalId %v already exist", externalId),
		}
	}
}

func (api AuthAPI) GetUserByExternalID(requestInfo RequestInfo, externalId string) (*User, error) {
	if !IsValidUserExternalID(externalId) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: externalId %v", externalId),
		}
	}
	// Retrieve user from DB
	user, err := api.UserRepo.GetUserByExternalID(externalId)

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
		}
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	// Check restrictions
	filteredUsers, err := api.GetAuthorizedUsers(requestInfo, user.Urn, USER_ACTION_GET_USER, []User{*user})
	if err != nil {
		return nil, err
	}

	if len(filteredUsers) > 0 {
		filteredUser := filteredUsers[0]
		return &filteredUser, nil
	}
	return nil, &Error{
		Code: UNAUTHORIZED_RESOURCES_ERROR,
		Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
			requestInfo.Identifier, user.Urn),
	}
}

func (api AuthAPI) ListUsers(requestInfo RequestInfo, filter *Filter) ([]string, int, error) {
	// Check parameters
	var total int
	orderByValidColumns := api.UserRepo.OrderByValidColumns(USER_ACTION_LIST_USERS)
	err := validateFilter(filter, orderByValidColumns)
	if err != nil {
		return nil, total, err
	}

	// Retrieve users with specified path prefix
	users, total, err := api.UserRepo.GetUsersFiltered(filter)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, total, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	// Check restrictions
	urnPrefix := GetUrnPrefix("", RESOURCE_USER, filter.PathPrefix)
	usersFiltered, err := api.GetAuthorizedUsers(requestInfo, urnPrefix, USER_ACTION_LIST_USERS, users)
	if err != nil {
		return nil, total, err
	}

	// Return user IDs
	externalIds := []string{}
	for _, u := range usersFiltered {
		externalIds = append(externalIds, u.ExternalID)
	}

	return externalIds, total, nil
}

func (api AuthAPI) UpdateUser(requestInfo RequestInfo, externalId string, newPath string) (*User, error) {
	if !IsValidPath(newPath) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: path %v", newPath),
		}
	}

	// Call repo to retrieve the user
	oldUser, err := api.GetUserByExternalID(requestInfo, externalId)
	if err != nil {
		return nil, err
	}

	// Check restrictions
	usersFiltered, err := api.GetAuthorizedUsers(requestInfo, oldUser.Urn, USER_ACTION_UPDATE_USER, []User{*oldUser})
	if err != nil {
		return nil, err
	}
	if len(usersFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, oldUser.Urn),
		}
	}

	auxUser := User{
		Urn: CreateUrn("", RESOURCE_USER, newPath, externalId),
	}

	// Check restrictions
	usersFiltered, err = api.GetAuthorizedUsers(requestInfo, auxUser.Urn, USER_ACTION_GET_USER, []User{auxUser})
	if err != nil {
		return nil, err
	}
	if len(usersFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, auxUser.Urn),
		}
	}

	user := User{
		ID:         oldUser.ID,
		ExternalID: oldUser.ExternalID,
		Path:       newPath,
		CreateAt:   oldUser.CreateAt,
		UpdateAt:   time.Now().UTC(),
		Urn:        auxUser.Urn,
	}

	updatedUser, err := api.UserRepo.UpdateUser(user)

	// Check unexpected DB error
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	LogOperation(api.Logger, requestInfo, fmt.Sprintf("User updated from %+v to %+v", oldUser, updatedUser))
	return updatedUser, nil

}

func (api AuthAPI) RemoveUser(requestInfo RequestInfo, externalId string) error {
	// Call repo to retrieve the user
	user, err := api.GetUserByExternalID(requestInfo, externalId)
	if err != nil {
		return err
	}

	// Check restrictions
	usersFiltered, err := api.GetAuthorizedUsers(requestInfo, user.Urn, USER_ACTION_DELETE_USER, []User{*user})
	if err != nil {
		return err
	}
	if len(usersFiltered) < 1 {
		return &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, user.Urn),
		}
	}

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
	LogOperation(api.Logger, requestInfo, fmt.Sprintf("User deleted %+v", user))
	return nil
}

func (api AuthAPI) ListGroupsByUser(requestInfo RequestInfo, filter *Filter) ([]UserGroups, int, error) {
	// Check parameters
	var total int
	orderByValidColumns := api.UserRepo.OrderByValidColumns(USER_ACTION_LIST_GROUPS_FOR_USER)
	err := validateFilter(filter, orderByValidColumns)
	if err != nil {
		return nil, total, err
	}

	// Call repo to retrieve the user
	user, err := api.GetUserByExternalID(requestInfo, filter.ExternalID)
	if err != nil {
		return nil, total, err
	}

	// Check restrictions
	usersFiltered, err := api.GetAuthorizedUsers(requestInfo, user.Urn, USER_ACTION_LIST_GROUPS_FOR_USER, []User{*user})
	if err != nil {
		return nil, total, err
	}
	if len(usersFiltered) < 1 {
		return nil, total, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, user.Urn),
		}
	}

	// Call group repo to retrieve groups associated to user
	groups, total, err := api.UserRepo.GetGroupsByUserID(user.ID, filter)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, total, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}
	// Transform to identifiers
	groupIDs := []UserGroups{}
	for _, g := range groups {
		groupIDs = append(groupIDs, UserGroups{
			Org:      g.GetGroup().Org,
			Name:     g.GetGroup().Name,
			CreateAt: g.GetDate(),
		})
	}

	return groupIDs, total, nil
}

// PRIVATE HELPER METHODS

func createUser(externalId string, path string) User {
	urn := CreateUrn("", RESOURCE_USER, path, externalId)
	user := User{
		ID:         uuid.NewV4().String(),
		ExternalID: externalId,
		Path:       path,
		CreateAt:   time.Now().UTC(),
		UpdateAt:   time.Now().UTC(),
		Urn:        urn,
	}

	return user
}
