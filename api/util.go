package api

import (
	"fmt"
)

const (
	RESOURCE_GROUP  = "group"
	RESOURCE_USER   = "user"
	RESOURCE_POLICY = "policy"
)

func CreateUrn(org string, resource string, path string, name string) string {
	switch resource {
	case RESOURCE_USER:
		return fmt.Sprintf("urn:iws:iam:user/%v%v", path, name)
	default:
		return fmt.Sprintf("urn:iws:iam:%v:%v/%v%v", org, resource, path, name)
	}
}
