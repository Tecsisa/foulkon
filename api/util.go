package api

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	// Resource types
	RESOURCE_GROUP  = "group"
	RESOURCE_USER   = "user"
	RESOURCE_POLICY = "policy"

	// Constraints
	MAX_EXTERNAL_ID_LENGTH = 128
	MAX_NAME_LENGTH        = 128
	MAX_PATH_LENGTH        = 512

	// Actions

	// User actions
	USER_ACTION_CREATE_USER          = "iam:CreateUser"
	USER_ACTION_DELETE_USER          = "iam:DeleteUser"
	USER_ACTION_GET_USER             = "iam:GetUser"
	USER_ACTION_LIST_USERS           = "iam:ListUsers"
	USER_ACTION_UPDATE_USER          = "iam:UpdateUser"
	USER_ACTION_LIST_GROUPS_FOR_USER = "iam:ListGroupsForUser"
	USER_ACTION_LIST_ORG_USERS       = "iam:ListOrgUsers"

	// Group actions
	GROUP_ACTION_CREATE_GROUP                 = "iam:CreateGroup"
	GROUP_ACTION_DELETE_GROUP                 = "iam:DeleteGroup"
	GROUP_ACTION_GET_GROUP                    = "iam:GetGroup"
	GROUP_ACTION_LIST_GROUPS                  = "iam:ListGroups"
	GROUP_ACTION_UPDATE_GROUP                 = "iam:UpdateGroup"
	GROUP_ACTION_LIST_MEMBERS                 = "iam:ListMembers"
	GROUP_ACTION_ADD_MEMBER                   = "iam:AddMember"
	GROUP_ACTION_REMOVE_MEMBER                = "iam:RemoveMember"
	GROUP_ACTION_ATTACH_GROUP_POLICY          = "iam:AttachGroupPolicy"
	GROUP_ACTION_DETACH_GROUP_POLICY          = "iam:DetachGroupPolicy"
	GROUP_ACTION_LIST_ATTACHED_GROUP_POLICIES = "iam:ListAttachedGroupPolicies"
	GROUP_ACTION_LIST_ALL_GROUPS              = "iam:ListAllGroups"

	// Policy actions
	POLICY_ACTION_CREATE_POLICY        = "iam:CreatePolicy"
	POLICY_ACTION_DELETE_POLICY        = "iam:DeletePolicy"
	POLICY_ACTION_UPDATE_POLICY        = "iam:UpdatePolicy"
	POLICY_ACTION_GET_POLICY           = "iam:GetPolicy"
	POLICY_ACTION_LIST_ATTACHED_GROUPS = "iam:ListAttachedGroups"
	POLICY_ACTION_LIST_POLICIES        = "iam:ListPolicies"
	POLICY_ACTION_LIST_ALL_POLICIES    = "iam:ListAllPolicies"
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
	r, _ := regexp.Compile(`^[\w+.@\-]+$`)
	return r.MatchString(externalID) && len(externalID) < MAX_EXTERNAL_ID_LENGTH
}

// this func validates group and policy names
func IsValidName(name string) bool {
	r, _ := regexp.Compile(`^[\w\-]+$`)
	return r.MatchString(name) && len(name) < MAX_NAME_LENGTH
}

func IsValidPath(path string) bool {
	r, _ := regexp.Compile(`^\*$|^/$|^/[\w+/\-]+\w+/$`)
	r2, _ := regexp.Compile(`[/:]{2,}`)
	return r.MatchString(path) && !r2.MatchString(path) && len(path) < MAX_PATH_LENGTH
}

func IsValidEffect(effect string) bool {
	if strings.ToLower(effect) == "allow" || strings.ToLower(effect) == "deny" {
		return true
	}
	return false
}

func IsValidAction(actions []string) bool {
	r, _ := regexp.Compile(`^[\w\-:]+[\w-*]+$`)
	r2, _ := regexp.Compile(`[*:]{2,}`)
	for _, action := range actions {
		if !r.MatchString(action) || r2.MatchString(action) || len(action) > MAX_NAME_LENGTH {
			return false
		}
	}
	return true
}

func IsValidResource(resources []string) bool {
	r, _ := regexp.Compile(`^[\w+.@\-/:]+[\w+.@\-*]+$`)
	r2, _ := regexp.Compile(`[/*:]{2,}`)
	for _, resource := range resources {
		if !r.MatchString(resource) || r2.MatchString(resource) || len(resource) > MAX_PATH_LENGTH {
			return false
		}
	}
	return true
}

func IsValidStatement(statements *[]Statement) bool {
	for _, statement := range *statements {
		if !IsValidEffect(statement.Effect) || !IsValidAction(statement.Action) || !IsValidResource(statement.Resources) {
			return false
		}
	}
	return true
}
