package api

import "fmt"

const (
	// Generic API error codes
	UNKNOWN_API_ERROR = "UnknownApiError"

	// User API error codes
	USER_BY_ID_NOT_FOUND          = "UserWithIDNotFound"
	USER_BY_EXTERNAL_ID_NOT_FOUND = "UserWithExternalIDNotFound"
	USER_ALREADY_EXIST            = "UserAlreadyExist"

	// Group API error codes
	GROUP_BY_ID_NOT_FOUND           = "GroupWithIdNotFound"
	GROUP_BY_ORG_AND_NAME_NOT_FOUND = "GroupWithOrgAndNameNotFound"
	GROUP_ALREADY_EXIST             = "GroupAlreadyExist"

	// GroupMembers error codes
	USER_ALREADY_IS_A_MEMBER_OF_GROUP = "UserAlreadyIsAMemberOfGroup"

	// Policy API error codes
	POLICY_ALREADY_EXIST = "PolicyAlreadyExist"
)

type Error struct {
	Code    string
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("Code: %v, Message: %v", e.Code, e.Message)
}
