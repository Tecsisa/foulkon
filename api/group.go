package api

import "time"

// Group domain
type Group struct {
	//Id   uint64    `json:"ID, omitempty"`
	Name string    `json:"Name, omitempty"`
	Path string    `json:"Path, omitempty"`
	Date time.Time `json:"Date, omitempty"`
	Urn  string    `json:"Urn, omitempty"`
	Org  string    `json:"Org, omitempty"`
}

type GroupsAPI struct {
	GroupRepo GroupRepo
}

// Retrieve group by id
func (g *GroupsAPI) GetGroupById(id string) (*Group, error) {
	return g.GroupRepo.GetGroupByID(id)
}

// Retrieve groups that has the path prefix and belongs to org parameter
func (g *GroupsAPI) GetListGroups(org string, path string) ([]Group, error) {
	return g.GroupRepo.GetGroupsByPath(org, path)
}

// Add a group to database
func (g *GroupsAPI) AddGroup(group Group) (*Group, error) {
	return g.GroupRepo.AddGroup(group)
}

// Remove group with this id
func (g *GroupsAPI) RemoveGroupById(id string) error {
	return g.GroupRepo.RemoveGroup(id)
}

// Get users for a group
func (g *GroupsAPI) GetUsersByGroupId(id string) ([]Group, error) {
	return g.GroupRepo.GetUsersByGroupID(id)
}
