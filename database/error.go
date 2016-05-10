package database

import "fmt"

const (
	// Database
	INTERNAL_ERROR = "InternalError"

	// User Codes
	USER_NOT_FOUND = "UserNotFound"

	// Group Codes
	GROUP_NOT_FOUND = "GroupNotFound"

	// Group User Relation Codes
	GROUP_USER_RELATION_NOT_FOUND = "GroupUserRelationNotFound"

	// Group Policy Relation Codes
	GROUP_POLICY_RELATION_NOT_FOUND = "GroupPolicyRelationNotFound"

	// Policy Codes
	POLICY_NOT_FOUND = "PolicyNotFound"
)

type Error struct {
	Code    string
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("Code: %v, Message: %v", e.Code, e.Message)
}
