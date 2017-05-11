package api

import "fmt"

const (
	// Generic API error codes
	UNKNOWN_API_ERROR            = "UnknownApiError"
	INVALID_PARAMETER_ERROR      = "InvalidParameterError"
	UNAUTHORIZED_RESOURCES_ERROR = "UnauthorizedResourcesError"

	// Authentication API error code
	AUTHENTICATION_API_ERROR = "AuthenticationApiError"

	// User API error codes
	USER_BY_EXTERNAL_ID_NOT_FOUND = "UserWithExternalIDNotFound"
	USER_ALREADY_EXIST            = "UserAlreadyExist"

	// Group API error codes
	GROUP_BY_ORG_AND_NAME_NOT_FOUND = "GroupWithOrgAndNameNotFound"
	GROUP_ALREADY_EXIST             = "GroupAlreadyExist"

	// GroupMembers error codes
	USER_IS_ALREADY_A_MEMBER_OF_GROUP = "UserIsAlreadyAMemberOfGroup"
	USER_IS_NOT_A_MEMBER_OF_GROUP     = "UserIsNotAMemberOfGroup"

	// GroupPolicies error codes
	POLICY_IS_ALREADY_ATTACHED_TO_GROUP = "PolicyIsAlreadyAttachedToGroup"
	POLICY_IS_NOT_ATTACHED_TO_GROUP     = "PolicyIsNotAttachedToGroup"

	// Policy API error codes
	POLICY_ALREADY_EXIST             = "PolicyAlreadyExist"
	POLICY_BY_ORG_AND_NAME_NOT_FOUND = "PolicyWithOrgAndNameNotFound"

	// Proxy resources API error codes
	PROXY_RESOURCE_ALREADY_EXIST             = "ProxyResourceAlreadyExist"
	PROXY_RESOURCE_BY_ORG_AND_NAME_NOT_FOUND = "ProxyResourceWithOrgAndNameNotFound"
	PROXY_RESOURCES_ROUTES_CONFLICT          = "ProxyResourcesRoutesConflict"

	// Regex error
	REGEX_NO_MATCH = "RegexNoMatch"
)

type Error struct {
	Code    string `json:"code, omitempty"`
	Message string `json:"message, omitempty"`
}

func (e Error) Error() string {
	return fmt.Sprintf("Code: %v, Message: %v", e.Code, e.Message)
}
