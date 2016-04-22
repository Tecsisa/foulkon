package api

// User repository that contains all user operations for this domain
type UserRepo interface {
	GetUserByID(id string) (*User, error)
	AddUser(User) (*User, error)
	GetUsersByPath(org string, path string) ([]User, error)
	GetGroupsByUserID(id string) ([]Group, error)
	RemoveUser(id string) error
}

// Group repository that contains all user operations for this domain
type GroupRepo interface {
}

// Policy repository that contains all user operations for this domain
type PolicyRepo interface {
}
