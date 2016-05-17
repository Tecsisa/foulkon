package api

import "fmt"

const (
	// Generic API error codes
	UNKNOWN_API_ERROR       = "UnknownApiError"
	MISSING_PARAMETER_ERROR = "MissingParameterError"
	INVALID_PARAMETER_ERROR = "InvalidParameterError"

	// User API error codes
	USER_BY_ID_NOT_FOUND          = "UserWithIDNotFound"
	USER_BY_EXTERNAL_ID_NOT_FOUND = "UserWithExternalIDNotFound"
	USER_ALREADY_EXIST            = "UserAlreadyExist"

	// Group API error codes
	GROUP_BY_ID_NOT_FOUND           = "GroupWithIdNotFound"
	GROUP_BY_ORG_AND_NAME_NOT_FOUND = "GroupWithOrgAndNameNotFound"
	GROUP_ALREADY_EXIST             = "GroupAlreadyExist"

	// GroupMembers error codes
	USER_IS_ALREADY_A_MEMBER_OF_GROUP = "UserIsAlreadyAMemberOfGroup"
	USER_IS_NOT_A_MEMBER_OF_GROUP     = "UserIsNotAMemberOfGroup"

	// GroupPolicies error codes
	POLICY_IS_ALREADY_ATTACHED_TO_GROUP = "PolicyIsAlreadyAttachedToGroup"

	// Policy API error codes
	POLICY_ALREADY_EXIST             = "PolicyAlreadyExist"
	POLICY_BY_ORG_AND_NAME_NOT_FOUND = "PolicyWithOrgAndNameNotFound"
)

type Error struct {
	Code    string
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("Code: %v, Message: %v", e.Code, e.Message)
}
