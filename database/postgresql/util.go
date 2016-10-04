package postgresql

import (
	"time"

	"github.com/Tecsisa/foulkon/api"
)

// GroupUser struct contains (Group-User) relationship
type GroupUser struct {
	User     *api.User
	Group    *api.Group
	CreateAt time.Time
}

// GetUser returns a member of a GroupUser relation
func (gu GroupUser) GetUser() *api.User {
	return gu.User
}

// GetGroup returns a Group of a GroupUser relation
func (gu GroupUser) GetGroup() *api.Group {
	return gu.Group
}

// GetDate returns the date when the relation was created
func (gu GroupUser) GetDate() time.Time {
	return gu.CreateAt
}

// PolicyGroup struct contains (Policy-Group) relationship
type PolicyGroup struct {
	Group    *api.Group
	Policy   *api.Policy
	CreateAt time.Time
}

// GetGroup returns a Group of a PolicyGroup relation
func (pg PolicyGroup) GetGroup() *api.Group {
	return pg.Group
}

// GetPolicy returns a Policy of a PolicyGroup relation
func (pg PolicyGroup) GetPolicy() *api.Policy {
	return pg.Policy
}

// GetDate returns the date when the relation was created
func (pg PolicyGroup) GetDate() time.Time {
	return pg.CreateAt
}
