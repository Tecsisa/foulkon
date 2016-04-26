package postgresql

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/tecsisa/authorizr/api"
)

type PostgresRepo struct {
	Dbmap *gorm.DB
}

func InitDb(datasourcename string) (*gorm.DB, error) {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := gorm.Open("postgres", datasourcename)
	if err != nil {
		return nil, err
	}

	// construct a gorp DbMap
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(5 * time.Minute)

	// Check connection
	err = db.DB().Ping()
	if err != nil {
		return nil, err
	}

	// Create tables if not exist
	err = db.AutoMigrate(&User{}).Error
	if err != nil {
		return nil, err
	}

	// Activate sql logger
	db.LogMode(true)

	return db, nil
}

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
	err := u.Dbmap.Where("external_id = ?", id).Find(user).Error
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

func (u PostgresRepo) GetUsersByPath(org string, path string) ([]api.User, error) {
	users := []User{}
	err := u.Dbmap.Where("org like ? AND path like ?", org, path+"%").Find(&users).Error

	// Error Handling
	if err != nil {
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
