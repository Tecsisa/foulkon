package api

import (
	"fmt"
	"github.com/tecsisa/authorizr/database"
	"strings"
)

// This Repo contains all Repositories that manages the domain
type Repo struct {
	UserRepo   UserRepo
	GroupRepo  GroupRepo
	PolicyRepo PolicyRepo
}

// User repository that contains all user operations for this domain
type UserRepo interface {
	// This method get a user with specified External ID.
	// If user exists, it will return the user with error param as nil
	// If user doesn't exists, it will return the error code database.USER_NOT_FOUND
	// If there is an error, it will return error param with associated error message
	// and error code database.INTERNAL_ERROR
	GetUserByExternalID(id string) (*User, error)
	GetUserByID(id string) (*User, error)

	// This method store a user.
	// If there are a problem inserting user it will return an database.Error error
	AddUser(user User) (*User, error)
	UpdateUser(user User, newPath string, newUrn string) (*User, error)

	GetUsersFiltered(pathPrefix string) ([]User, error)
	GetGroupsByUserID(id string) ([]Group, error)
	RemoveUser(id string) error
}

// Group repository that contains all user operations for this domain
type GroupRepo interface {
	GetGroupById(id string) (*Group, error)

	GetGroupByName(org string, name string) (*Group, error)
	GetGroupUserRelation(userID string, groupID string) (*GroupMembers, error)
	GetAllGroupUserRelation(groupID string) (*GroupMembers, error)
	GetGroupPolicyRelation(groupID string, policyID string) (*GroupPolicies, error)
	GetAllGroupPolicyRelation(groupID string) (*GroupPolicies, error)
	GetGroupsFiltered(org string, pathPrefix string) ([]Group, error)
	RemoveGroup(id string) error

	AddGroup(group Group) (*Group, error)
	AddMember(userID string, groupID string) error
	RemoveMember(userID string, groupID string) error
	UpdateGroup(group Group, newName string, newPath string, newUrn string) (*Group, error)
	AttachPolicy(groupID string, policyID string) error
}

// Policy repository that contains all user operations for this domain
type PolicyRepo interface {
	GetPolicyById(id string) (*Policy, error)
	GetPolicyByName(org string, name string) (*Policy, error)
	AddPolicy(policy Policy) (*Policy, error)
	UpdatePolicy(policy Policy, newName string, newPath string, newUrn string, newStatements []Statement) (*Policy, error)
	DeletePolicy(id string) error
	GetPoliciesFiltered(org string, pathPrefix string) ([]Policy, error)
}

type AuthResources struct {
	AllowedUrnPrefixes []string
	AllowedFullUrns    []string
	DeniedUrnPrefixes  []string
	DeniedFullUrns     []string
}

func (repo *Repo) Authorize(externalID string, action string, resource string) (*AuthResources, error) {
	// Get user if exist
	user, err := repo.UserRepo.GetUserByExternalID(externalID)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	// Get groups for this user
	groups, err := repo.getGroupsByUser(user.ID)

	// Error handling
	if err != nil {
		return nil, err
	}

	// Get policies by groups
	policies, err := repo.getPoliciesByGroups(groups)

	// Error handling
	if err != nil {
		return nil, err
	}

	// Retrieve statements for action requested for these policies
	statements, err := repo.getStatementsByRequestedAction(policies, action)

	// Error handling
	if err != nil {
		return nil, err
	}

	// Retrieve restrictions restrictions
	var authResources *AuthResources
	if isFullUrn(resource) {
		authResources, err = getRestrictionsWhenResourceRequestedIsFullUrn(statements, resource)
	} else {
		authResources, err = getRestrictionsWhenResourceRequestedIsPrefix(statements, resource)
	}

	// Error handling
	if err != nil {
		return nil, err
	}

	// Clean up repeated resources
	return cleanRepeatedRestrictions(authResources)
}

func (repo *Repo) getGroupsByUser(userID string) ([]Group, error) {
	// Get group relations by user
	groups, err := repo.UserRepo.GetGroupsByUserID(userID)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	// Return groups
	return groups, nil
}

func (repo *Repo) getPoliciesByGroups(groups []Group) ([]Policy, error) {
	// Retrieve per each group its attached policies
	if groups == nil || len(groups) < 1 {
		return nil, nil
	}

	// Create a empty slice
	policies := []Policy{}

	for _, group := range groups {
		// Retrieve policies for this group
		policyRelations, err := repo.GroupRepo.GetAllGroupPolicyRelation(group.ID)

		// Error handling
		if err != nil {
			//Transform to DB error
			dbError := err.(*database.Error)
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}

		for _, policy := range policyRelations.Policies {
			policies = append(policies, policy)
		}
	}

	// Return filled policies
	return policies, nil
}

func (repo *Repo) getStatementsByRequestedAction(policies []Policy, actionRequested string) ([]Statement, error) {
	// Check received policies
	if policies == nil || len(policies) < 1 {
		return nil, nil
	}

	// Retrieve statements related to action requested
	statements := []Statement{}
	for _, policy := range policies {
		// Policy always has an statement at least
		for _, statement := range *policy.Statements {
			// Check if there are an action related to requested
			if isActionContained(actionRequested, statement.Action) {
				statements = append(statements, statement)
			}
		}
	}

	// Return statements
	return statements, nil
}

func cleanRepeatedRestrictions(authResources *AuthResources) (*AuthResources, error) {
	// TODO rsoleto: Falta implementar
	return authResources, nil
}

func isActionContained(actionRequested string, statementActions []string) bool {
	match := false
	for _, statementAction := range statementActions {
		// Prefixes
		if strings.ContainsAny(statementAction, "*") {
			value := strings.Trim(statementAction, "*")
			if len(value) < 1 || strings.HasPrefix(actionRequested, value) {
				match = true
				break
			}
		} else if strings.Compare(actionRequested, statementAction) == 0 {
			match = true
			break
		}
	}

	return match
}

func isResourceContained(resource string, resourcePrefix string) bool {
	prefix := strings.Trim(resourcePrefix, "*")
	if len(prefix) < 1 {
		return true
	} else {
		return strings.HasPrefix(resource, prefix)
	}
}

func isFullUrn(resource string) bool {
	if strings.ContainsAny(resource, "*") {
		return false
	} else {
		return true
	}
}

func getRestrictionsWhenResourceRequestedIsPrefix(statements []Statement, resource string) (*AuthResources, error) {
	authResources := &AuthResources{
		AllowedUrnPrefixes: []string{},
		AllowedFullUrns:    []string{},
		DeniedUrnPrefixes:  []string{},
		DeniedFullUrns:     []string{},
	}
	for _, statement := range statements {
		for _, statementResource := range statement.Resources {
			// If is full URN the statement of resource, we need to check if is a sub resource
			if isFullUrn(statementResource) && isResourceContained(statementResource, resource) {
				if statement.Effect == "allow" {
					authResources.AllowedFullUrns = append(authResources.AllowedFullUrns, statementResource)
				} else {
					authResources.DeniedFullUrns = append(authResources.DeniedFullUrns, statementResource)
				}
			} else {
				// We have two prefixes, now we have to decide which is shorter,
				// and then if shorter contains other resource
				switch {
				case len(statementResource) > len(resource):
					if isResourceContained(statementResource, resource) {
						if statement.Effect == "allow" {
							authResources.AllowedUrnPrefixes = append(authResources.AllowedUrnPrefixes, statementResource)
						} else {
							authResources.DeniedUrnPrefixes = append(authResources.DeniedUrnPrefixes, statementResource)
						}
					}
				case len(resource) > len(statementResource):
					if isResourceContained(resource, statementResource) {
						if statement.Effect == "allow" {
							authResources.AllowedUrnPrefixes = append(authResources.AllowedUrnPrefixes, statementResource)
						} else {
							authResources.DeniedUrnPrefixes = append(authResources.DeniedUrnPrefixes, statementResource)
						}
					}
				case resource == statementResource:
					if statement.Effect == "allow" {
						authResources.AllowedUrnPrefixes = append(authResources.AllowedUrnPrefixes, statementResource)
					} else {
						authResources.DeniedUrnPrefixes = append(authResources.DeniedUrnPrefixes, statementResource)
					}
				default: //Do nothing
				}
			}

		}
	}

	if len(authResources.AllowedUrnPrefixes) < 1 && len(authResources.AllowedFullUrns) < 1 {
		return nil, &Error{
			Code:    UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("There aren't authorized resources"),
		}
	} else {
		return authResources, nil
	}

}

func getRestrictionsWhenResourceRequestedIsFullUrn(statements []Statement, resource string) (*AuthResources, error) {
	authResources := &AuthResources{
		AllowedUrnPrefixes: []string{},
		AllowedFullUrns:    []string{},
		DeniedUrnPrefixes:  []string{},
		DeniedFullUrns:     []string{},
	}
	for _, statement := range statements {
		for _, statementResource := range statement.Resources {
			switch {
			case isFullUrn(statementResource) && statementResource == resource:
				if statement.Effect == "allow" {
					authResources.AllowedFullUrns = append(authResources.AllowedFullUrns, statementResource)
				} else {
					authResources.DeniedFullUrns = append(authResources.DeniedFullUrns, statementResource)
				}
			case !isFullUrn(statementResource) && isResourceContained(resource, statementResource):
				if statement.Effect == "allow" {
					authResources.AllowedUrnPrefixes = append(authResources.AllowedUrnPrefixes, statementResource)
				} else {
					authResources.DeniedUrnPrefixes = append(authResources.DeniedUrnPrefixes, statementResource)
				}
			default: //Do nothing
			}
		}
	}

	if len(authResources.AllowedUrnPrefixes) < 1 && len(authResources.AllowedFullUrns) < 1 {
		return nil, &Error{
			Code:    UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("There aren't authorized resources"),
		}
	} else {
		return authResources, nil
	}
}
