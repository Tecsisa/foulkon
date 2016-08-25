package postgresql

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" //GORM needs to import the lib/pq driver
)

type PostgresRepo struct {
	Dbmap *gorm.DB
}

func InitDb(datasourcename string, idleConns string, maxOpenConns string, connTTL string) (*gorm.DB, error) {
	// connect to db using GORM - github.com/jinzhu/gorm
	db, err := gorm.Open("postgres", datasourcename)
	if err != nil {
		return nil, err
	}

	// construct a gorp DbMap
	idle, err := strconv.Atoi(idleConns)
	if err != nil {
		return nil, fmt.Errorf("Invalid postgresql idleConns param: %v", idleConns)
	}
	maxOpen, err := strconv.Atoi(maxOpenConns)
	if err != nil {
		return nil, fmt.Errorf("Invalid postgresql maxOpenConns param: %v", maxOpenConns)
	}
	ttl, err := strconv.Atoi(connTTL)
	if err != nil {
		return nil, fmt.Errorf("Invalid postgresql connTTL param: %v", connTTL)
	}
	db.DB().SetMaxIdleConns(idle)
	db.DB().SetMaxOpenConns(maxOpen)
	db.DB().SetConnMaxLifetime(time.Duration(ttl) * time.Second)

	// Check connection
	err = db.DB().Ping()
	if err != nil {
		return nil, err
	}

	// Create tables if not exist =
	err = db.AutoMigrate(&User{}, &Group{}, &Policy{}, &Statement{}, &GroupUserRelation{}, &GroupPolicyRelation{}).Error
	if err != nil {
		return nil, err
	}

	// TODO:
	// Activate sql logger
	//db.LogMode(true)

	return db, nil
}

// User table
type User struct {
	ID         string `gorm:"primary_key"`
	ExternalID string `gorm:"not null;unique"`
	Path       string `gorm:"not null"`
	CreateAt   int64  `gorm:"not null"`
	Urn        string `gorm:"not null;unique"`
}

// User's table name
func (User) TableName() string {
	return "users"
}

// Group table
type Group struct {
	ID       string `gorm:"primary_key"`
	Name     string `gorm:"not null"`
	Path     string `gorm:"not null"`
	Org      string `gorm:"not null"`
	CreateAt int64  `gorm:"not null"`
	Urn      string `gorm:"not null;unique"`
}

// Group's table name
func (Group) TableName() string {
	return "groups"
}

// Policy table
type Policy struct {
	ID       string `gorm:"primary_key"`
	Name     string `gorm:"not null"`
	Path     string `gorm:"not null"`
	Org      string `gorm:"not null"`
	CreateAt int64  `gorm:"not null"`
	Urn      string `gorm:"not null;unique"`
}

// Policy's table name
func (Policy) TableName() string {
	return "policies"
}

// Statement table
type Statement struct {
	ID        string `gorm:"primary_key"`
	PolicyID  string `gorm:"not null"`
	Effect    string `gorm:"not null"`
	Actions   string `gorm:"not null"`
	Resources string `gorm:"not null"`
}

// Statement's table name
func (Statement) TableName() string {
	return "statements"
}

// Group-Users Relationship
type GroupUserRelation struct {
	UserID  string `gorm:"primary_key"`
	GroupID string `gorm:"primary_key"`
}

// GroupUserRelation's table name
func (GroupUserRelation) TableName() string {
	return "group_user_relations"
}

// Group Policy table
type GroupPolicyRelation struct {
	GroupID  string `gorm:"primary_key"`
	PolicyID string `gorm:"primary_key"`
}

// GroupPolicyRelation's table name
func (GroupPolicyRelation) TableName() string {
	return "group_policy_relations"
}
