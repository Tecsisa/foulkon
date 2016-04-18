package api

import (
	"fmt"
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
func (u *UsersAPI) GetUserById(id uint64) (User, error) {
	user, err := u.UserRepo.GetUserByID(id)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	return user, err
}

// Retrieve users that has path
func (u *UsersAPI) GetListUsers(path string) ([]User, error) {
	users, err := u.UserRepo.GetUsersByPath(path)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	return users, err
}

// Add an User to database
func (u *UsersAPI) AddUser(user User) error {
	return u.UserRepo.AddUser(user)
}

// Remove user with this id
func (u *UsersAPI) RemoveUserById(id uint64) error {
	err := u.UserRepo.RemoveUser(id)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	return err
}

// Get groups for an user
func (u *UsersAPI) GetGroupsByUserId(id uint64) ([]Group, error) {
	groups, err := u.UserRepo.GetGroupsByUserID(id)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	return groups, err
}
