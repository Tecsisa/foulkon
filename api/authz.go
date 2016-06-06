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

type EffectRestriction struct {
	Effect       string        `json:"Effect, omitempty"`
	Restrictions *Restrictions `json:"Restrictions, omitempty"`
}

type Restrictions struct {
	AllowedUrnPrefixes []string
	AllowedFullUrns    []string
	DeniedUrnPrefixes  []string
	DeniedFullUrns     []string
}

// Struct that represent a external resource, with only URN attribute
type ExternalResource struct {
	Urn string
}

func (e ExternalResource) GetUrn() string {
	return e.Urn
}

// Return authorized users according to restrictions associated to authenticated user
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

// Return authorized groups according to restrictions associated to authenticated user
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

// Return authorized policies according to restrictions associated to authenticated user
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

// Get user effect to do the action over the resource
func (api *AuthAPI) GetEffectByUserActionResource(user AuthenticatedUser, action string, resourceUrn string) (*EffectRestriction, error) {
	// Validate parameters
	if err := IsValidAction([]string{action}); err != nil {
		// Transform to API error
		apiError := err.(*Error)
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: apiError.Message,
		}
	}
	if err := IsValidResource([]string{resourceUrn}); err != nil {
		// Transform to API error
		apiError := err.(*Error)
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: apiError.Message,
		}
	}

	// Retrieve restrictions
	effectRestriction := &EffectRestriction{}
	restrictions, err := api.getRestrictions(user.Identifier, action, resourceUrn)
	if err != nil {
		return nil, err
	}

	// Return depending on if resource is a prefix or not
	if isFullUrn(resourceUrn) {
		if (isAllowedResource(ExternalResource{Urn: resourceUrn}, *restrictions)) {
			effectRestriction.Effect = "allow"
		} else {
			effectRestriction.Effect = "deny"
		}
	} else {
		effectRestriction.Restrictions = restrictions
	}

	return effectRestriction, nil
}

// Private Helper Methods

// This method use authenticated user to retrieve its restrictions, apply it to a resource URN (could be a prefix)
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

	api.Logger.Debugf("Restrictions: %v", *restrictions)

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

// Get restrictions for this action and full resource or prefix resource, attached to this authenticated user
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
	statements := getStatementsByRequestedAction(policies, action)

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

// Retrieve groups that user is member
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

// Retrieve policies attached to a slice of groups
func (api *AuthAPI) getPoliciesByGroups(groups []Group) ([]Policy, error) {
	// Retrieve per each group its attached policies
	if groups == nil || len(groups) < 1 {
		return nil, nil
	}

	// Create an empty slice
	policies := []Policy{}

	for _, group := range groups {
		// Retrieve policies for this group
		policiesAttached, err := api.GroupRepo.GetPoliciesAttached(group.ID)

		// Error handling
		if err != nil {
			//Transform to DB error
			dbError := err.(*database.Error)
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}

		for _, policy := range policiesAttached {
			policies = append(policies, policy)
		}
	}

	// Return filled policies
	return policies, nil
}

// Filter a slice of statements for a specified action
func getStatementsByRequestedAction(policies []Policy, actionRequested string) []Statement {
	// Check received policies
	if policies == nil || len(policies) < 1 {
		return nil
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
	return statements
}

// Clean restrictions that are repeated or contained by others (Deny is prior than Allow)
func cleanRepeatedRestrictions(authResources *Restrictions) *Restrictions {
	// TODO rsoleto: Falta implementar
	return authResources
}

// Return if an action is contained inside a slice of statement's actions
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

// Return if a resource is contained by the prefix
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

// Retrieve restrictions for a specified resource prefix according to the statements
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

// Retrieve restrictions for a specified resource according to the statements
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

// Remove resources that are not allowed by the restrictions
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
