package api

import (
	"fmt"
	"regexp"
)

const (
	RESOURCE_GROUP  = "group"
	RESOURCE_USER   = "user"
	RESOURCE_POLICY = "policy"

	MAX_EXTERNAL_ID_LENGTH = 128
	MAX_PATH_LENGTH        = 512
)

func CreateUrn(org string, resource string, path string, name string) string {
	switch resource {
	case RESOURCE_USER:
		return fmt.Sprintf("urn:iws:iam:user%v%v", path, name)
	default:
		return fmt.Sprintf("urn:iws:iam:%v:%v%v%v", org, resource, path, name)
	}
}

func IsValidUserExternalID(externalID string) bool {
	r, _ := regexp.Compile(`^[\w+=,.@-]+$`)
	return r.MatchString(externalID) && len(externalID) < MAX_EXTERNAL_ID_LENGTH
}

func IsValidPath(path string) bool {
	r, _ := regexp.Compile(`^/$|^/[\w+/]+\w+/$`)
	r2, _ := regexp.Compile(`/{2,}`)
	return r.MatchString(path) && !r2.MatchString(path) && len(path) < MAX_PATH_LENGTH
}
