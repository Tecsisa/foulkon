package api

// User repository that contains all user operations for this domain
type UserRepo interface {
	GetUserByID(id string) (*User, error)
	AddUser(User) (*User, error)
	GetUsersFiltered(path string) ([]User, error)
	GetGroupsByUserID(id string) ([]Group, error)
	RemoveUser(id string) error
}

// Group repository that contains all user operations for this domain
type GroupRepo interface {
	GetGroupByID(id string) (*Group, error)
	AddGroup(Group) (*Group, error)
	GetGroupsByPath(org string, path string) ([]Group, error)
	GetUsersByGroupID(id string) ([]Group, error)
	RemoveGroup(id string) error
}

// Policy repository that contains all user operations for this domain
type PolicyRepo interface {
}
