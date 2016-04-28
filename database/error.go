package database

import "fmt"

const (
	// Database
	INTERNAL_ERROR = "InternalError"

	// User Codes
	USER_NOT_FOUND     = "UserNotFound"
	USER_ALREADY_EXIST = "UserAlreadyExist"
)

type Error struct {
	Code    string
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("Code: %v, Message: %v", e.Code, e.Message)
}
