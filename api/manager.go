package api

// User repository that contains all user operations for this domain
type UserRepo interface {
	GetUserByID(id uint64) (User, error)
	AddUser(User) error
	GetUsersByPath(org string, path string) ([]User, error)
	GetGroupsByUserID(id uint64) ([]Group, error)
	RemoveUser(id uint64) error
}

// Group repository that contains all user operations for this domain
type GroupRepo interface {
}

// Policy repository that contains all user operations for this domain
type PolicyRepo interface {
}
