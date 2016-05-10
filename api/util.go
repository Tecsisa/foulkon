package api

import (
	"fmt"
)

const (
	RESOURCE_GROUP  = "group"
	RESOURCE_USER   = "user"
	RESOURCE_POLICY = "policy"
)

func CreateUrn(org string, resource string, path string) string {
	switch resource {
	case RESOURCE_USER:
		return fmt.Sprintf("urn:iws:iam:user/%v", path)
	default:
		return fmt.Sprintf("urn:iws:iam:%v:%v/%v", org, resource, path)
	}
}
