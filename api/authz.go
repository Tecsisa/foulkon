package api

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Tecsisa/foulkon/database"
)

// TYPE DEFINITIONS

type RequestInfo struct {
	Identifier string
	Admin      bool
	RequestID  string
}

type EffectRestriction struct {
	Effect       string        `json:"effect, omitempty"`
	Restrictions *Restrictions `json:"restrictions, omitempty"`
}

type Restrictions struct {
	AllowedUrnPrefixes []string `json:"allowedUrnPrefixes, omitempty"`
	AllowedFullUrns    []string `json:"allowedFullUrns, omitempty"`
	DeniedUrnPrefixes  []string `json:"deniedUrnPrefixes, omitempty"`
	DeniedFullUrns     []string `json:"deniedFullUrns, omitempty"`
}

type ExternalResource struct {
	Urn string `json:"urn, omitempty"`
}

func (e ExternalResource) GetUrn() string {
	return e.Urn
}

// AUTHZ API IMPLEMENTATION

// GetAuthorizedUsers returns authorized users for specified resource+action
func (api AuthAPI) GetAuthorizedUsers(requestInfo RequestInfo, resourceUrn string, action string, users []User) ([]User, error) {
	resourcesToAuthorize := []Resource{}
	for _, usr := range users {
		resourcesToAuthorize = append(resourcesToAuthorize, usr)
	}
	resources, err := api.getAuthorizedResources(requestInfo, resourceUrn, action, resourcesToAuthorize)
	if err != nil {
		return nil, err
	}
	usersFiltered := []User{}
	for _, res := range resources {
		usersFiltered = append(usersFiltered, res.(User))
	}
	return usersFiltered, nil
}

// GetAuthorizedGroups returns authorized users for specified user combined with resource+action
func (api AuthAPI) GetAuthorizedGroups(requestInfo RequestInfo, resourceUrn string, action string, groups []Group) ([]Group, error) {
	resourcesToAuthorize := []Resource{}
	for _, group := range groups {
		resourcesToAuthorize = append(resourcesToAuthorize, group)
	}
	resources, err := api.getAuthorizedResources(requestInfo, resourceUrn, action, resourcesToAuthorize)
	if err != nil {
		return nil, err
	}
	groupsFiltered := []Group{}
	for _, res := range resources {
		groupsFiltered = append(groupsFiltered, res.(Group))
	}
	return groupsFiltered, nil
}

// GetAuthorizedPolicies returns authorized policies for specified user combined with resource+action
func (api AuthAPI) GetAuthorizedPolicies(requestInfo RequestInfo, resourceUrn string, action string, policies []Policy) ([]Policy, error) {
	resourcesToAuthorize := []Resource{}
	for _, policy := range policies {
		resourcesToAuthorize = append(resourcesToAuthorize, policy)
	}
	resources, err := api.getAuthorizedResources(requestInfo, resourceUrn, action, resourcesToAuthorize)
	if err != nil {
		return nil, err
	}
	policiesFiltered := []Policy{}
	for _, res := range resources {
		policiesFiltered = append(policiesFiltered, res.(Policy))
	}
	return policiesFiltered, nil
}

// GetAuthorizedExternalResources returns the resources where the specified user has the action granted
func (api AuthAPI) GetAuthorizedExternalResources(requestInfo RequestInfo, action string, resources []string) ([]string, error) {
	// Validate parameters
	if err := AreValidActions([]string{action}); err != nil {
		// Transform to API error
		apiError := err.(*Error)
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: apiError.Message,
		}
	}
	if len(resources) < 1 || len(resources) > MAX_RESOURCE_NUMBER {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter Resources. Resources can't be empty or bigger than %v elements", MAX_RESOURCE_NUMBER),
		}
	}
	externalResources := []Resource{}
	for _, res := range resources {
		if !isFullUrn(res) {
			return nil, &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: fmt.Sprintf("Invalid parameter resource %v. Urn prefixes are not allowed here", res),
			}
		}
		if err := AreValidResources([]string{res}); err != nil {
			// Transform to API error
			apiError := err.(*Error)
			return nil, &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: apiError.Message,
			}
		}
		externalResources = append(externalResources, ExternalResource{Urn: res})
	}
	if strings.Contains(action, "*") {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter action %v. Action parameter can't be a prefix", action),
		}
	}

	allowedUrns, err := api.getAuthorizedResources(requestInfo, "urn:*", action, externalResources)
	if err != nil {
		return nil, err
	}

	response := []string{}
	for _, res := range allowedUrns {
		response = append(response, res.GetUrn())
	}

	return response, nil
}

// PRIVATE HELPER METHODS

// getAuthorizedResources retrieves filtered resources where the authenticated user has permissions
func (api AuthAPI) getAuthorizedResources(requestInfo RequestInfo, resourceUrn string, action string, resources []Resource) ([]Resource, error) {
	// If user is an admin return all resources without restriction
	if requestInfo.Admin {
		return resources, nil
	}

	// Check authorization for this user
	restrictions, err := api.getRestrictions(requestInfo.Identifier, action, resourceUrn)
	if err != nil {
		return nil, err
	}

	api.Logger.Debugf("Restrictions: %v", *restrictions)

	// Check if there are some restrictions for this urn resource
	if len(restrictions.AllowedFullUrns) < 1 && len(restrictions.AllowedUrnPrefixes) < 1 {
		return nil, &Error{
			Code:    UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v", requestInfo.Identifier, resourceUrn),
		}
	}

	// Filter resources
	resourcesFiltered := filterResources(resources, restrictions)

	return resourcesFiltered, nil
}

// Get restrictions for this action and full resource or prefix resource, attached to this authenticated user
func (api AuthAPI) getRestrictions(externalID string, action string, resource string) (*Restrictions, error) {
	// Get user if exists
	user, err := api.UserRepo.GetUserByExternalID(externalID)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		switch dbError.Code {
		case database.USER_NOT_FOUND:
			return nil, &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: fmt.Sprintf("Authenticated user with externalId %v not found. Unable to retrieve permissions.", externalID),
			}
		default:
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	}

	groups, err := api.getGroupsByUser(user.ID)
	if err != nil {
		return nil, err
	}

	policies, err := api.getPoliciesByGroups(groups)
	if err != nil {
		return nil, err
	}

	// Retrieve valid statements
	statements := getStatementsByRequestedAction(policies, action)

	// Retrieve restrictions
	var authResources *Restrictions
	authResources = getRestrictions(statements, resource, isFullUrn(resource))

	return authResources, nil
}

func (api AuthAPI) getGroupsByUser(userID string) ([]Group, error) {
	groups, _, err := api.UserRepo.GetGroupsByUserID(userID, &Filter{})
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	return groups, nil
}

// Retrieve policies attached to a slice of groups
func (api AuthAPI) getPoliciesByGroups(groups []Group) ([]Policy, error) {
	if groups == nil || len(groups) < 1 {
		return nil, nil
	}

	// Create an empty slice
	policies := []Policy{}

	// Retrieve per each group its attached policies
	for _, group := range groups {
		// Retrieve policies for this group
		policiesAttached, _, err := api.GroupRepo.GetAttachedPolicies(group.ID, &Filter{})
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

	return policies, nil
}

// Filter a slice of statements for a specified action
func getStatementsByRequestedAction(policies []Policy, requestedAction string) []Statement {
	// Check received policies
	if policies == nil || len(policies) < 1 {
		return nil
	}

	statements := []Statement{}
	for _, policy := range policies {
		for _, statement := range *policy.Statements {
			if isActionContained(requestedAction, statement.Actions) {
				statements = append(statements, statement)
			}
		}
	}

	return statements
}

// Returns true if an action is contained inside a slice of statements
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

// Returns true if a resource is contained in a prefix
func isContainedOrEqual(resource string, resourcePrefix string) bool {
	prefix := strings.Trim(resourcePrefix, "*")
	if len(prefix) < 1 {
		return true
	}
	return strings.HasPrefix(resource, prefix)
}

func isFullUrn(resource string) bool {
	return !strings.ContainsAny(resource, "*")
}

// Insert restriction with filtering and cleaning
func (r *Restrictions) insertRestriction(allow bool, fullUrn bool, resource string) {

	// HELPER FUNCS
	// delete an element from a slice, given an index
	deleteElementFunc := func(i int, slice []string) []string {
		if len(slice) > 1 {
			slice = append(slice[:i], slice[i+1:]...)
		} else {
			slice = []string{}
		}
		return slice
	}

	// skip resource insertion if contained in a given slice
	skip := func(slice []string) bool {
		for _, urn := range slice {
			if isContainedOrEqual(resource, urn) {
				return true
			}
		}
		return false
	}

	if allow {
		if fullUrn {
			// if urn is already contained wherever, skip
			if skip(r.DeniedUrnPrefixes) || skip(r.AllowedUrnPrefixes) || skip(r.DeniedFullUrns) || skip(r.AllowedFullUrns) {
				return
			}

			r.AllowedFullUrns = append(r.AllowedFullUrns, resource)

		} else { // urnPrefix
			// if urnPrefix is already contained in any denied prefixes, skip
			if skip(r.DeniedUrnPrefixes) {
				return
			}

			// if urnPrefix is already contained in any allowed prefixes, skip
			for i, allowPrefix := range r.AllowedUrnPrefixes {
				if isContainedOrEqual(resource, allowPrefix) {
					return
				}
				// if urnPrefix contains other prefixes already inserted, delete them
				if isContainedOrEqual(allowPrefix, resource) {
					r.AllowedUrnPrefixes = deleteElementFunc(i, r.AllowedUrnPrefixes)
				}
			}
			for i, allowUrn := range r.AllowedFullUrns {
				// if urnPrefix contains fullUrns already inserted, delete them
				if isContainedOrEqual(allowUrn, resource) {
					r.AllowedFullUrns = deleteElementFunc(i, r.AllowedFullUrns)
				}
			}

			r.AllowedUrnPrefixes = append(r.AllowedUrnPrefixes, resource)

		}
	} else { // deny
		if fullUrn {
			// if urn is already contained in denied restrictions, skip
			if skip(r.DeniedUrnPrefixes) || skip(r.DeniedFullUrns) {
				return
			}

			// if urn is already allowed, delete it
			for i, allowUrn := range r.AllowedFullUrns {
				if isContainedOrEqual(resource, allowUrn) {
					r.AllowedFullUrns = deleteElementFunc(i, r.AllowedFullUrns)
				}
			}

			r.DeniedFullUrns = append(r.DeniedFullUrns, resource)

		} else { // urnPrefix
			for i, denyPrefix := range r.DeniedUrnPrefixes {
				// if denyPrefix is contained in prefixes already inserted, skip
				if isContainedOrEqual(resource, denyPrefix) {
					return
				}
				// if denyPrefix contains prefixes already inserted, delete them
				if isContainedOrEqual(denyPrefix, resource) {
					r.DeniedUrnPrefixes = deleteElementFunc(i, r.DeniedUrnPrefixes)
				}
			}

			// create waitGroup in order to launch 3 concurrent goroutines to delete duplicate restrictions
			var wg sync.WaitGroup
			wg.Add(3)

			go func() {
				defer wg.Done()
				for i, allowPrefix := range r.AllowedUrnPrefixes {
					// if denyPrefix contains allowed prefixes already inserted, delete them
					if resource == allowPrefix {
						r.AllowedUrnPrefixes = deleteElementFunc(i, r.AllowedUrnPrefixes)
					}
				}

			}()
			go func() {
				defer wg.Done()
				for i, allowUrn := range r.AllowedFullUrns {

					// if denyPrefix contains allowed full urns already inserted, delete them
					if isContainedOrEqual(allowUrn, resource) {
						r.AllowedFullUrns = deleteElementFunc(i, r.AllowedFullUrns)
					}
				}
			}()
			go func() {
				defer wg.Done()
				for i, denyUrn := range r.DeniedFullUrns {
					// if denyPrefix contains denied full urns already inserted, delete them
					if isContainedOrEqual(denyUrn, resource) {
						r.DeniedFullUrns = deleteElementFunc(i, r.DeniedFullUrns)
					}
				}
			}()

			wg.Wait()
			r.DeniedUrnPrefixes = append(r.DeniedUrnPrefixes, resource)
		}
	}
}

// Retrieve restrictions for a specified resource according to the statements
func getRestrictions(statements []Statement, resource string, resourceIsFullUrn bool) *Restrictions {
	restrictions := &Restrictions{
		AllowedUrnPrefixes: []string{},
		AllowedFullUrns:    []string{},
		DeniedUrnPrefixes:  []string{},
		DeniedFullUrns:     []string{},
	}
	if statements != nil || len(statements) > 0 {
		for _, statement := range statements {
			for _, statementResource := range statement.Resources {
				// Append resource to allowed or denied resources, if the resource URN is not a prefix (full URN), and is contained inside the passed resource.
				// Else, it means that resource is a prefix, so we have to check if the passed resource contains it or vice versa.
				statementIsFullUrn := isFullUrn(statementResource)
				statementIsAllow := statement.Effect == "allow"

				if !resourceIsFullUrn {
					if isContainedOrEqual(statementResource, resource) || isContainedOrEqual(resource, statementResource) {
						restrictions.insertRestriction(statementIsAllow, statementIsFullUrn, statementResource)
					}
				} else {
					// Insert restriction if resource is contained in statements
					if isContainedOrEqual(resource, statementResource) {
						restrictions.insertRestriction(statementIsAllow, statementIsFullUrn, statementResource)
					}
				}
			}
		}
	}

	return restrictions
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
			if isContainedOrEqual(resource.GetUrn(), restriction) {
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
			if isContainedOrEqual(resource.GetUrn(), restriction) {
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

	return allowed && !denied
}
