package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Tecsisa/foulkon/api"
)

// statusFromError attempts to convert API errors to appropriate GRPC status
// responses (analogous to processHttpResponse, http/handler.go)
func statusFromError(err error) error {
	if err, ok := err.(*api.Error); ok {
		var code codes.Code
		switch err.Code {
		case api.USER_ALREADY_EXIST,
			api.GROUP_ALREADY_EXIST, // TODO we only do users so far
			api.USER_IS_ALREADY_A_MEMBER_OF_GROUP,
			api.PROXY_RESOURCE_ALREADY_EXIST,
			api.POLICY_IS_ALREADY_ATTACHED_TO_GROUP,
			api.POLICY_ALREADY_EXIST,
			api.PROXY_RESOURCES_ROUTES_CONFLICT,
			api.AUTH_OIDC_PROVIDER_ALREADY_EXIST:
			code = codes.AlreadyExists
		case api.UNAUTHORIZED_RESOURCES_ERROR:
			code = codes.PermissionDenied
		case api.USER_BY_EXTERNAL_ID_NOT_FOUND, api.GROUP_BY_ORG_AND_NAME_NOT_FOUND,
			api.USER_IS_NOT_A_MEMBER_OF_GROUP, api.POLICY_IS_NOT_ATTACHED_TO_GROUP,
			api.POLICY_BY_ORG_AND_NAME_NOT_FOUND, api.PROXY_RESOURCE_BY_ORG_AND_NAME_NOT_FOUND,
			api.AUTH_OIDC_PROVIDER_BY_NAME_NOT_FOUND:
			code = codes.NotFound
		case api.INVALID_PARAMETER_ERROR, api.REGEX_NO_MATCH:
			// Unexpected input in validation parameters
			code = codes.InvalidArgument
		default:
			code = codes.Internal
		}
		return status.Error(code, err.Message)
	}
	return err
}
