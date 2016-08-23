package api

import (
	"fmt"
	"time"

	"github.com/satori/go.uuid"
	"github.com/tecsisa/foulkon/database"
)

// TYPE DEFINITIONS

// Policy domain
type Policy struct {
	ID         string       `json:"id, omitempty"`
	Name       string       `json:"name, omitempty"`
	Path       string       `json:"path, omitempty"`
	Org        string       `json:"org, omitempty"`
	Urn        string       `json:"urn, omitempty"`
	CreateAt   time.Time    `json:"createAt, omitempty"`
	Statements *[]Statement `json:"statements, omitempty"`
}

func (p Policy) String() string {
	return fmt.Sprintf("[id: %v, name: %v, path: %v, org: %v, urn: %v, createAt: %v, statements: %v]",
		p.ID, p.Name, p.Path, p.Org, p.Urn, p.CreateAt.Format("2006-01-02 15:04:05 MST"), p.Statements)
}

func (p Policy) GetUrn() string {
	return p.Urn
}

// Policy identifier to retrieve them from DB
type PolicyIdentity struct {
	Org  string `json:"org, omitempty"`
	Name string `json:"name, omitempty"`
}

type Statement struct {
	Effect    string   `json:"effect, omitempty"`
	Actions   []string `json:"actions, omitempty"`
	Resources []string `json:"resources, omitempty"`
}

func (s Statement) String() string {
	return fmt.Sprintf("[effect: %v, actions: %v, resources: %v]", s.Effect, s.Actions, s.Resources)
}

// POLICY API IMPLEMENTATION

func (api AuthAPI) AddPolicy(requestInfo RequestInfo, name string, path string, org string, statements []Statement) (*Policy, error) {
	// Validate fields
	if !IsValidName(name) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: name %v", name),
		}
	}
	if !IsValidOrg(org) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: org %v", org),
		}
	}
	if !IsValidPath(path) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: path %v", path),
		}

	}
	err := AreValidStatements(&statements)
	if err != nil {
		apiError := err.(*Error)
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: apiError.Message,
		}

	}

	policy := createPolicy(name, path, org, &statements)

	// Check restrictions
	policiesFiltered, err := api.GetAuthorizedPolicies(requestInfo, policy.Urn, POLICY_ACTION_CREATE_POLICY, []Policy{policy})
	if err != nil {
		return nil, err
	}
	if len(policiesFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, policy.Urn),
		}
	}

	// Check if policy already exists
	_, err = api.PolicyRepo.GetPolicyByName(org, name)

	// Check if policy could be retrieved
	if err != nil {
		// Transform to DB error
		dbError := err.(*database.Error)
		switch dbError.Code {
		// Policy doesn't exist in DB
		case database.POLICY_NOT_FOUND:
			// Create policy
			createdPolicy, err := api.PolicyRepo.AddPolicy(policy)

			// Check if there is an unexpected error in DB
			if err != nil {
				//Transform to DB error
				dbError := err.(*database.Error)
				return nil, &Error{
					Code:    UNKNOWN_API_ERROR,
					Message: dbError.Message,
				}
			}

			LogOperation(api.Logger, requestInfo, fmt.Sprintf("Policy created %+v", createdPolicy))
			return createdPolicy, nil
		default: // Unexpected error
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	} else { // Fail if policy exists
		return nil, &Error{
			Code:    POLICY_ALREADY_EXIST,
			Message: fmt.Sprintf("Unable to create policy, policy with org %v and name %v already exist", org, name),
		}
	}
}

func (api AuthAPI) GetPolicyByName(requestInfo RequestInfo, org string, policyName string) (*Policy, error) {
	// Validate fields
	if !IsValidName(policyName) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: name %v", policyName),
		}
	}
	// Validate org
	if !IsValidOrg(org) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: org %v", org),
		}
	}

	// Call repo to retrieve the policy
	policy, err := api.PolicyRepo.GetPolicyByName(org, policyName)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		// Policy doesn't exist in DB
		if dbError.Code == database.POLICY_NOT_FOUND {
			return nil, &Error{
				Code:    POLICY_BY_ORG_AND_NAME_NOT_FOUND,
				Message: dbError.Message,
			}
		} else { // Unexpected error
			return nil, &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: dbError.Message,
			}
		}
	}

	// Check restrictions
	policiesFiltered, err := api.GetAuthorizedPolicies(requestInfo, policy.Urn, POLICY_ACTION_GET_POLICY, []Policy{*policy})
	if err != nil {
		return nil, err
	}
	if len(policiesFiltered) > 0 {
		policyFiltered := policiesFiltered[0]
		return &policyFiltered, nil
	} else {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, policy.Urn),
		}
	}
}

func (api AuthAPI) ListPolicies(requestInfo RequestInfo, org string, pathPrefix string) ([]PolicyIdentity, error) {
	// Validate fields
	if len(org) > 0 && !IsValidOrg(org) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: org %v", org),
		}
	}
	if len(pathPrefix) > 0 && !IsValidPath(pathPrefix) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: PathPrefix %v", pathPrefix),
		}
	}
	if len(pathPrefix) == 0 {
		pathPrefix = "/"
	}

	// Call repo to retrieve the policies
	policies, err := api.PolicyRepo.GetPoliciesFiltered(org, pathPrefix)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	// Check restrictions to list
	var urnPrefix string
	if len(org) == 0 {
		urnPrefix = "*"
	} else {
		urnPrefix = GetUrnPrefix(org, RESOURCE_POLICY, pathPrefix)
	}
	policiesFiltered, err := api.GetAuthorizedPolicies(requestInfo, urnPrefix, POLICY_ACTION_LIST_POLICIES, policies)
	if err != nil {
		return nil, err
	}

	policyIDs := []PolicyIdentity{}
	for _, p := range policiesFiltered {
		policyIDs = append(policyIDs, PolicyIdentity{
			Org:  p.Org,
			Name: p.Name,
		})
	}

	return policyIDs, nil
}

func (api AuthAPI) UpdatePolicy(requestInfo RequestInfo, org string, policyName string, newName string, newPath string,
	newStatements []Statement) (*Policy, error) {
	// Validate fields
	if !IsValidName(newName) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: new name %v", newName),
		}
	}
	if !IsValidPath(newPath) {
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: fmt.Sprintf("Invalid parameter: new path %v", newPath),
		}

	}
	err := AreValidStatements(&newStatements)
	if err != nil {
		apiError := err.(*Error)
		return nil, &Error{
			Code:    INVALID_PARAMETER_ERROR,
			Message: apiError.Message,
		}

	}

	// Call repo to retrieve the policy
	policyDB, err := api.GetPolicyByName(requestInfo, org, policyName)
	if err != nil {
		return nil, err
	}

	// Check restrictions
	policiesFiltered, err := api.GetAuthorizedPolicies(requestInfo, policyDB.Urn, POLICY_ACTION_UPDATE_POLICY, []Policy{*policyDB})
	if err != nil {
		return nil, err
	}
	if len(policiesFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, policyDB.Urn),
		}
	}

	// Check if policy with "newName" exists
	targetPolicy, err := api.GetPolicyByName(requestInfo, org, newName)

	if err == nil && targetPolicy.ID != policyDB.ID {
		// Policy already exists
		return nil, &Error{
			Code:    POLICY_ALREADY_EXIST,
			Message: fmt.Sprintf("Policy name: %v already exists", newName),
		}
	}
	if err != nil {
		if apiError := err.(*Error); apiError.Code == UNAUTHORIZED_RESOURCES_ERROR || apiError.Code == UNKNOWN_API_ERROR {
			return nil, err
		}
	}

	// Get Policy Updated
	policyToUpdate := createPolicy(newName, newPath, org, &newStatements)

	// Check restrictions
	policiesFiltered, err = api.GetAuthorizedPolicies(requestInfo, policyToUpdate.Urn, POLICY_ACTION_UPDATE_POLICY, []Policy{policyToUpdate})
	if err != nil {
		return nil, err
	}
	if len(policiesFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, policyToUpdate.Urn),
		}
	}

	// Update policy
	policy, err := api.PolicyRepo.UpdatePolicy(*policyDB, newName, newPath, policyToUpdate.Urn, newStatements)

	// Check unexpected DB error
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	LogOperation(api.Logger, requestInfo, fmt.Sprintf("Policy updated from %+v to %+v", policyDB, policy))
	return policy, nil
}

func (api AuthAPI) RemovePolicy(requestInfo RequestInfo, org string, name string) error {

	// Call repo to retrieve the policy
	policy, err := api.GetPolicyByName(requestInfo, org, name)
	if err != nil {
		return err
	}

	// Check restrictions
	policiesFiltered, err := api.GetAuthorizedPolicies(requestInfo, policy.Urn, POLICY_ACTION_DELETE_POLICY, []Policy{*policy})
	if err != nil {
		return err
	}
	if len(policiesFiltered) < 1 {
		return &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, policy.Urn),
		}
	}

	err = api.PolicyRepo.RemovePolicy(policy.ID)
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	LogOperation(api.Logger, requestInfo, fmt.Sprintf("Policy deleted %+v", policy))
	return nil
}

func (api AuthAPI) ListAttachedGroups(requestInfo RequestInfo, org string, name string) ([]string, error) {

	// Call repo to retrieve the policy
	policy, err := api.GetPolicyByName(requestInfo, org, name)
	if err != nil {
		return nil, err
	}

	// Check restrictions
	policiesFiltered, err := api.GetAuthorizedPolicies(requestInfo, policy.Urn, POLICY_ACTION_LIST_ATTACHED_GROUPS, []Policy{*policy})
	if err != nil {
		return nil, err
	}
	if len(policiesFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, policy.Urn),
		}
	}

	// Call repo to retrieve the attached groups
	groups, err := api.PolicyRepo.GetAttachedGroups(policy.ID)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	groupNames := []string{}
	for _, g := range groups {
		groupNames = append(groupNames, g.Name)
	}

	return groupNames, nil
}

// PRIVATE HELPER METHODS

func createPolicy(name string, path string, org string, statements *[]Statement) Policy {
	urn := CreateUrn(org, RESOURCE_POLICY, path, name)
	policy := Policy{
		ID:         uuid.NewV4().String(),
		Name:       name,
		Path:       path,
		Org:        org,
		Urn:        urn,
		CreateAt:   time.Now().UTC(),
		Statements: statements,
	}

	return policy
}
