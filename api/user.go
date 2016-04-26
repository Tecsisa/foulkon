package api

import (
	"time"
)

// User domain
type User struct {
	ExternalID string    `json:"ExternalID, omitempty"`
	Path       string    `json:"Path, omitempty"`
	Date       time.Time `json:"Date, omitempty"`
	Urn        string    `json:"Urn, omitempty"`
}

// User api
type UsersAPI struct {
	UserRepo UserRepo
}

// Retrieve user by id
func (u *UsersAPI) GetUserById(id string) (*User, error) {
	return u.UserRepo.GetUserByID(id)
}

// Retrieve users that has path
func (u *UsersAPI) GetListUsers(org string, path string) ([]User, error) {
	return u.UserRepo.GetUsersByPath(org, path)
}

// Add an User to database
func (u *UsersAPI) AddUser(user User) (*User, error) {
	return u.UserRepo.AddUser(user)
}

// Remove user with this id
func (u *UsersAPI) RemoveUserById(id string) error {
	return u.UserRepo.RemoveUser(id)
}

// Get groups for an user
func (u *UsersAPI) GetGroupsByUserId(id string) ([]Group, error) {
	return u.UserRepo.GetGroupsByUserID(id)
}
