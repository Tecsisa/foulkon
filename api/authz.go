package api

import (
	"fmt"
	"strings"

	"github.com/tecsisa/authorizr/database"
)

// Interface that all resource have to implement
type Resource interface {
	// This method return urn associated to resource
	GetUrn() string
}

// User struct that define what role is
type AuthenticatedUser struct {
	Identifier string
	Admin      bool
}

type Restrictions struct {
	AllowedUrnPrefixes []string
	AllowedFullUrns    []string
	DeniedUrnPrefixes  []string
	DeniedFullUrns     []string
}

// Return authorized users according to
func (api *AuthAPI) GetUsersAuthorized(user AuthenticatedUser, resourceUrn string, action string, users []User) ([]User, error) {
	resourcesToAuthorize := []Resource{}
	for _, usr := range users {
		resourcesToAuthorize = append(resourcesToAuthorize, usr)
	}
	resources, err := api.getAuthorizedResources(user, resourceUrn, action, resourcesToAuthorize)
	if err != nil {
		return nil, err
	}
	usersFiltered := []User{}
	for _, res := range resources {
		usersFiltered = append(usersFiltered, res.(User))
	}
	return usersFiltered, nil
}

func (api *AuthAPI) GetGroupsAuthorized(user AuthenticatedUser, resourceUrn string, action string, groups []Group) ([]Group, error) {
	resourcesToAuthorize := []Resource{}
	for _, group := range groups {
		resourcesToAuthorize = append(resourcesToAuthorize, group)
	}
	resources, err := api.getAuthorizedResources(user, resourceUrn, action, resourcesToAuthorize)
	if err != nil {
		return nil, err
	}
	groupsFiltered := []Group{}
	for _, res := range resources {
		groupsFiltered = append(groupsFiltered, res.(Group))
	}
	return groupsFiltered, nil
}

func (api *AuthAPI) GetPoliciesAuthorized(user AuthenticatedUser, resourceUrn string, action string, policies []Policy) ([]Policy, error) {
	resourcesToAuthorize := []Resource{}
	for _, policy := range policies {
		resourcesToAuthorize = append(resourcesToAuthorize, policy)
	}
	resources, err := api.getAuthorizedResources(user, resourceUrn, action, resourcesToAuthorize)
	if err != nil {
		return nil, err
	}
	policiesFiltered := []Policy{}
	for _, res := range resources {
		policiesFiltered = append(policiesFiltered, res.(Policy))
	}
	return policiesFiltered, nil
}

// Private Helper Methods

// This method use authenticated user to retrieve its restrictions and apply it to a resource URN (could be a prefix)
// and retrieve filtered resources
func (api *AuthAPI) getAuthorizedResources(user AuthenticatedUser, resourceUrn string, action string, resources []Resource) ([]Resource, error) {

	// If user is an admin return all resources without restriction
	if user.Admin {
		return resources, nil
	}

	// Check authorization for this user
	restrictions, err := api.getRestrictions(user.Identifier, action, resourceUrn)
	if err != nil {
		return nil, err
	}

	// Check if there are some restrictions for this urn resource
	if len(restrictions.AllowedFullUrns) < 1 && len(restrictions.AllowedUrnPrefixes) < 1 {
		return nil, &Error{
			Code:    UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with external ID %v is not allowed to access to resource %v", user.Identifier, resourceUrn),
		}
	}

	// Filter resources
	resourcesFiltered := filterResources(resources, restrictions)

	return resourcesFiltered, nil
}

func (api *AuthAPI) getRestrictions(externalID string, action string, resource string) (*Restrictions, error) {
	// Get user if exist
	user, err := api.UserRepo.GetUserByExternalID(externalID)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		switch dbError.Code {
		case database.USER_NOT_FOUND:
			return nil, &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: fmt.Sprintf("User authenticated with external ID %v not found. It can't be possible retrieve its permission", externalID),
			}
		default:
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	}

	// Get groups for this user
	groups, err := api.getGroupsByUser(user.ID)

	// Error handling
	if err != nil {
		return nil, err
	}

	// Get policies by groups
	policies, err := api.getPoliciesByGroups(groups)

	// Error handling
	if err != nil {
		return nil, err
	}

	// Retrieve statements for action requested for these policies
	statements, err := api.getStatementsByRequestedAction(policies, action)

	// Error handling
	if err != nil {
		return nil, err
	}

	// Retrieve restrictions restrictions
	var authResources *Restrictions
	if isFullUrn(resource) {
		authResources = getRestrictionsWhenResourceRequestedIsFullUrn(statements, resource)
	} else {
		authResources = getRestrictionsWhenResourceRequestedIsPrefix(statements, resource)
	}

	// Clean up repeated resources
	return cleanRepeatedRestrictions(authResources), nil
}

func (api *AuthAPI) getGroupsByUser(userID string) ([]Group, error) {
	// Get group relations by user
	groups, err := api.UserRepo.GetGroupsByUserID(userID)

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

func (api *AuthAPI) getPoliciesByGroups(groups []Group) ([]Policy, error) {
	// Retrieve per each group its attached policies
	if groups == nil || len(groups) < 1 {
		return nil, nil
	}

	// Create a empty slice
	policies := []Policy{}

	for _, group := range groups {
		// Retrieve policies for this group
		policyRelations, err := api.GroupRepo.GetAllGroupPolicyRelation(group.ID)

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

func (api *AuthAPI) getStatementsByRequestedAction(policies []Policy, actionRequested string) ([]Statement, error) {
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

func cleanRepeatedRestrictions(authResources *Restrictions) *Restrictions {
	// TODO rsoleto: Falta implementar
	return authResources
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

func getRestrictionsWhenResourceRequestedIsPrefix(statements []Statement, resource string) *Restrictions {
	authResources := &Restrictions{
		AllowedUrnPrefixes: []string{},
		AllowedFullUrns:    []string{},
		DeniedUrnPrefixes:  []string{},
		DeniedFullUrns:     []string{},
	}
	if statements != nil || len(statements) > 0 {
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
	}

	return authResources
}

func getRestrictionsWhenResourceRequestedIsFullUrn(statements []Statement, resource string) *Restrictions {
	authResources := &Restrictions{
		AllowedUrnPrefixes: []string{},
		AllowedFullUrns:    []string{},
		DeniedUrnPrefixes:  []string{},
		DeniedFullUrns:     []string{},
	}
	if statements != nil || len(statements) > 0 {
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
	}

	return authResources
}

func filterResources(resources []Resource, restrictions *Restrictions) []Resource {
	filteredResource := []Resource{}
	for _, r := range resources {
		if isAllowedResource(r, *restrictions) {
			filteredResource = append(filteredResource, r)
		}
	}

	return filteredResource
}

// Check if resource is allowed or not
func isAllowedResource(resource Resource, restrictions Restrictions) bool {
	allowed := false
	denied := false
	// Check deny restrictions
	if len(restrictions.DeniedUrnPrefixes) > 0 {
		for _, restriction := range restrictions.DeniedUrnPrefixes {
			if isResourceContained(resource.GetUrn(), restriction) {
				denied = true
				break
			}
		}
	}
	if len(restrictions.DeniedFullUrns) > 0 && !denied {
		for _, restriction := range restrictions.DeniedFullUrns {
			if resource.GetUrn() == restriction {
				denied = true
				break
			}
		}
	}

	// Check allow restrictions
	if len(restrictions.AllowedUrnPrefixes) > 0 && !denied {
		for _, restriction := range restrictions.AllowedUrnPrefixes {
			if isResourceContained(resource.GetUrn(), restriction) {
				allowed = true
				break
			}
		}
	}
	if len(restrictions.AllowedFullUrns) > 0 && !denied && !allowed {
		for _, restriction := range restrictions.AllowedFullUrns {
			if resource.GetUrn() == restriction {
				allowed = true
				break
			}
		}
	}

	// If it is allowed and not denied
	if allowed && !denied {
		return true
	} else {
		return false
	}
}
