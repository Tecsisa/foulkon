package api

// User repository that contains all user operations for this domain
type UserRepo interface {
	// This method get a user with specified External ID.
	// If user exists, it will return the user with error param as nil
	// If user doesn't exists, it will return the error code database.USER_NOT_FOUND
	// If there is an error, it will return error param with associated error message
	// and error code database.INTERNAL_ERROR
	GetUserByExternalID(id string) (*User, error)

	// This method store a user.
	// If there are a problem inserting user it will return an database.Error error
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
