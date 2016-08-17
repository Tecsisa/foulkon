package api

import (
	"fmt"

	"github.com/Sirupsen/logrus"
)

const (
	// Generic API error codes
	UNKNOWN_API_ERROR            = "UnknownApiError"
	INVALID_PARAMETER_ERROR      = "InvalidParameterError"
	UNAUTHORIZED_RESOURCES_ERROR = "UnauthorizedResourcesError"

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

func LogErrorMessage(logger *logrus.Logger, requestInfo RequestInfo, err *Error) {
	logger.WithFields(logrus.Fields{
		"requestID": requestInfo.RequestID,
		"code":      err.Code,
	}).Error(err.Message)
}
