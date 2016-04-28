package api

import "fmt"

const (
	// Generic API error codes
	UNKNOWN_API_ERROR = "UnknownApiError"

	// User API error codes
	USER_BY_EXTERNAL_ID_NOT_FOUND = "UserWithExternalIDNotFound"
	USER_ALREADY_EXIST            = "UserAlreadyExist"
)

type Error struct {
	Code    string
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("Code: %v, Message: %v", e.Code, e.Message)
}
