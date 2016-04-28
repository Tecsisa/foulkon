package api

// User repository that contains all user operations for this domain
type UserRepo interface {
	// This method get a user with specified External ID.
	// If user exists, it will return the user with error param as nil
	// If user doesn't exists, it will return the error code database.Err
	// If there is an error, it will return error param with associated error, bool param as false and user as nil
	GetUserByExternalID(id string) (*User, error)

	// This method store a user.
	// If there is an user with this external ID, it will return an error
	// If user doesn't exists, it will return the bool param as false and other params as nil
	// If there is an error, it will return error param with associated error, bool param as false and user as nil
	AddUser(User) (*User, error)

	GetUsersFiltered(pathPrefix string) ([]User, error)
	GetGroupsByUserID(id string) ([]Group, error)
	RemoveUser(id string) error
}

// Group repository that contains all user operations for this domain
type GroupRepo interface {
}

// Policy repository that contains all user operations for this domain
type PolicyRepo interface {
}
