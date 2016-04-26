package postgresql

import (
	"errors"
	"time"

	"github.com/tecsisa/authorizr/api"
)

// User database
type User struct {
	ID         int    `gorm:"primary_key"`
	ExternalID string `gorm:"not null;unique"`
	Path       string `gorm:"not null"`
	CreateDate int64  `gorm:"not null"`
	Urn        string `gorm:"not null;unique"`
}

// set User's table name to be `profiles
func (User) TableName() string {
	return "users"
}

func (u PostgresRepo) GetUserByID(id string) (*api.User, error) {
	user := &User{}
	err := u.Dbmap.Where("external_id like ?", id).Find(user).Error
	// Error Handling
	if err != nil {
		return nil, err
	}
	if user != nil {
		return userDBToUserAPI(user), nil
	}
	return nil, nil
}

func (u PostgresRepo) AddUser(user api.User) (*api.User, error) {
	userDB := &User{
		ExternalID: user.ExternalID,
		Path:       user.Path,
		CreateDate: time.Now().UTC().UnixNano(),
		Urn:        user.Urn,
	}

	err := u.Dbmap.Create(userDB).Error
	if err != nil {
		return nil, err
	}
	return userDBToUserAPI(userDB), nil
}

func (u PostgresRepo) GetUsersFiltered(path string) ([]api.User, error) {
	users := []User{}
	query := u.Dbmap
	if len(path) > 0 {
		query = query.Where("path like ?", path+"%")
	}

	if err := query.Where("name = ?", "jinzhu").First(&users).Error; err != nil {
		return nil, err
	}

	if users != nil {
		apiusers := make([]api.User, len(users), cap(users))
		for i, u := range users {
			apiusers[i] = *userDBToUserAPI(&u)
		}
		return apiusers, nil
	}

	return nil, nil
}

func (u PostgresRepo) GetGroupsByUserID(id string) ([]api.Group, error) {
	return nil, nil
}

func (u PostgresRepo) RemoveUser(id string) error {
	user := &User{}
	err := u.Dbmap.Where("external_id = ?", id).Find(user).Error
	// Error Handling
	if err != nil {
		return err
	}

	if user != nil {
		return u.Dbmap.Delete(&user).Error
	} else {
		return errors.New("User not found")
	}
}

// Transform a user retrieved from db into a user for API
func userDBToUserAPI(userdb *User) *api.User {
	return &api.User{
		ExternalID: userdb.ExternalID,
		Path:       userdb.Path,
		Date:       time.Unix(0, userdb.CreateDate).UTC(),
		Urn:        userdb.Urn,
	}
}
