package api

import (
	"time"
)

// User domain
type User struct {
	Id   uint64    `json:"ID, omitempty"`
	Name string    `json:"Name, omitempty"`
	Path string    `json:"Path, omitempty"`
	Date time.Time `json:"Date, omitempty"`
	Urn  string    `json:"Urn, omitempty"`
	Org  string    `json:"Org, omitempty"`
}

// User api
type UsersAPI struct {
	UserRepo UserRepo
}

// Retrieve user by id
func (u *UsersAPI) GetUserById(id uint64) (*User, error) {
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
func (u *UsersAPI) RemoveUserById(id uint64) error {
	return u.UserRepo.RemoveUser(id)
}

// Get groups for an user
func (u *UsersAPI) GetGroupsByUserId(id uint64) ([]Group, error) {
	return u.UserRepo.GetGroupsByUserID(id)
}
