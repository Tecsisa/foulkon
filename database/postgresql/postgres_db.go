package postgresql

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Tecsisa/foulkon/api"
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

	// construct a gorm DbMap
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

	// Create tables if not exist
	err = db.AutoMigrate(&User{}, &Group{}, &Policy{}, &Statement{}, &GroupUserRelation{}, &GroupPolicyRelation{},
		&ProxyResource{}, &OidcProvider{}, &OidcClient{}).Error
	if err != nil {
		return nil, err
	}

	return db, nil
}

// User table
type User struct {
	ID         string `gorm:"primary_key"`
	ExternalID string `gorm:"not null;unique"`
	Path       string `gorm:"not null"`
	CreateAt   int64  `gorm:"not null"`
	UpdateAt   int64  `gorm:"not null"`
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
	UpdateAt int64  `gorm:"not null"`
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
	UpdateAt int64  `gorm:"not null"`
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
	UserID   string `gorm:"primary_key"`
	GroupID  string `gorm:"primary_key"`
	CreateAt int64  `gorm:"not null"`
}

// GroupUserRelation's table name
func (GroupUserRelation) TableName() string {
	return "group_user_relations"
}

// Group Policy table
type GroupPolicyRelation struct {
	GroupID  string `gorm:"primary_key"`
	PolicyID string `gorm:"primary_key"`
	CreateAt int64  `gorm:"not null"`
}

// GroupPolicyRelation's table name
func (GroupPolicyRelation) TableName() string {
	return "group_policy_relations"
}

func (pr PostgresRepo) OrderByValidColumns(action string) []string {
	switch action {
	case api.USER_ACTION_LIST_USERS:
		return []string{"path", "external_id", "create_at", "update_at", "urn"}
	case api.USER_ACTION_LIST_GROUPS_FOR_USER:
		return []string{"create_at"}
	case api.GROUP_ACTION_LIST_GROUPS:
		return []string{"name", "path", "org", "create_at", "update_at", "urn"}
	case api.GROUP_ACTION_LIST_MEMBERS:
		return []string{"create_at"}
	case api.GROUP_ACTION_LIST_ATTACHED_GROUP_POLICIES:
		return []string{"create_at"}
	case api.POLICY_ACTION_LIST_POLICIES:
		return []string{"name", "path", "org", "create_at", "update_at", "urn"}
	case api.POLICY_ACTION_LIST_ATTACHED_GROUPS:
		return []string{"create_at"}
	case api.PROXY_ACTION_LIST_RESOURCES:
		return []string{"name", "path", "org", "host", "path_resource", "method",
			"urn_resource", "urn", "action", "create_at", "update_at"}
	case api.AUTH_OIDC_ACTION_LIST_PROVIDERS:
		return []string{"name", "path", "create_at", "update_at", "urn"}
	default:
		return nil
	}
}

// ProxyResource table
type ProxyResource struct {
	ID           string `gorm:"primary_key"`
	Name         string `gorm:"not null"`
	Org          string `gorm:"not null"`
	Path         string `gorm:"not null"`
	Host         string `gorm:"not null;unique_index:idx_resource"`
	PathResource string `gorm:"not null;unique_index:idx_resource"`
	Method       string `gorm:"not null;unique_index:idx_resource"`
	UrnResource  string `gorm:"not null;unique_index:idx_resource"`
	Urn          string `gorm:"not null"`
	Action       string `gorm:"not null;unique_index:idx_resource"`
	CreateAt     int64  `gorm:"not null"`
	UpdateAt     int64  `gorm:"not null"`
}

// ProxyResource's table name
func (ProxyResource) TableName() string {
	return "proxy_resources"
}

// Auth OIDC Provider table
type OidcProvider struct {
	ID        string `gorm:"primary_key"`
	Name      string `gorm:"not null"`
	Path      string `gorm:"not null"`
	Urn       string `gorm:"not null;unique"`
	CreateAt  int64  `gorm:"not null"`
	UpdateAt  int64  `gorm:"not null"`
	IssuerURL string `gorm:"not null"`
}

// OidcProvider's table name
func (OidcProvider) TableName() string {
	return "oidc_providers"
}

// Auth OIDC Client table
type OidcClient struct {
	ID             string `gorm:"primary_key"`
	OidcProviderID string `gorm:"not null;unique_index:idx_oidc_client"`
	Name           string `gorm:"not null;unique_index:idx_oidc_client"`
}

// OidcClient's table name
func (OidcClient) TableName() string {
	return "oidc_clients"
}
