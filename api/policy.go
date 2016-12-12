package api

import (
	"fmt"
	"time"

	"github.com/Tecsisa/foulkon/database"
	"github.com/satori/go.uuid"
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
	UpdateAt   time.Time    `json:"updateAt, omitempty"`
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

type PolicyGroups struct {
	Group    string    `json:"group, omitempty"`
	CreateAt time.Time `json:"attached, omitempty"`
}

func (s Statement) String() string {
	return fmt.Sprintf("[effect: %v, actions: %v, resources: %v]", s.Effect, s.Actions, s.Resources)
}

// POLICY API IMPLEMENTATION

func (api WorkerAPI) AddPolicy(requestInfo RequestInfo, name string, path string, org string, statements []Statement) (*Policy, error) {
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

			LogOperation(requestInfo.RequestID, requestInfo.Identifier, fmt.Sprintf("Policy created %+v", createdPolicy))
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

func (api WorkerAPI) GetPolicyByName(requestInfo RequestInfo, org string, policyName string) (*Policy, error) {
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
		}
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
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
	}
	return nil, &Error{
		Code: UNAUTHORIZED_RESOURCES_ERROR,
		Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
			requestInfo.Identifier, policy.Urn),
	}
}

func (api WorkerAPI) ListPolicies(requestInfo RequestInfo, filter *Filter) ([]PolicyIdentity, int, error) {
	// Validate fields
	var total int
	orderByValidColumns := api.PolicyRepo.OrderByValidColumns(POLICY_ACTION_LIST_POLICIES)
	err := validateFilter(filter, orderByValidColumns)
	if err != nil {
		return nil, total, err
	}

	// Call repo to retrieve the policies
	policies, total, err := api.PolicyRepo.GetPoliciesFiltered(filter)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, total, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	// Check restrictions to list
	var urnPrefix string
	if len(filter.Org) == 0 {
		urnPrefix = "*"
	} else {
		urnPrefix = GetUrnPrefix(filter.Org, RESOURCE_POLICY, filter.PathPrefix)
	}
	policiesFiltered, err := api.GetAuthorizedPolicies(requestInfo, urnPrefix, POLICY_ACTION_LIST_POLICIES, policies)
	if err != nil {
		return nil, total, err
	}

	policyIDs := []PolicyIdentity{}
	for _, p := range policiesFiltered {
		policyIDs = append(policyIDs, PolicyIdentity{
			Org:  p.Org,
			Name: p.Name,
		})
	}

	return policyIDs, total, nil
}

func (api WorkerAPI) UpdatePolicy(requestInfo RequestInfo, org string, policyName string, newName string, newPath string,
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

	// Call repo to retrieve the old policy
	oldPolicy, err := api.GetPolicyByName(requestInfo, org, policyName)
	if err != nil {
		return nil, err
	}

	// Check restrictions
	policiesFiltered, err := api.GetAuthorizedPolicies(requestInfo, oldPolicy.Urn, POLICY_ACTION_UPDATE_POLICY, []Policy{*oldPolicy})
	if err != nil {
		return nil, err
	}
	if len(policiesFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, oldPolicy.Urn),
		}
	}

	// Check if policy with "newName" exists
	targetPolicy, err := api.GetPolicyByName(requestInfo, org, newName)

	if err == nil && targetPolicy.ID != oldPolicy.ID {
		// Policy already exists
		return nil, &Error{
			Code:    POLICY_ALREADY_EXIST,
			Message: fmt.Sprintf("Policy name: %v already exists", newName),
		}
	}

	if err != nil {
		if apiError := err.(*Error); apiError.Code != POLICY_BY_ORG_AND_NAME_NOT_FOUND {
			return nil, err
		}
	}

	auxPolicy := Policy{
		Urn: CreateUrn(org, RESOURCE_POLICY, newPath, newName),
	}

	// Check restrictions
	policiesFiltered, err = api.GetAuthorizedPolicies(requestInfo, auxPolicy.Urn, POLICY_ACTION_UPDATE_POLICY, []Policy{auxPolicy})
	if err != nil {
		return nil, err
	}
	if len(policiesFiltered) < 1 {
		return nil, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, auxPolicy.Urn),
		}
	}

	policy := Policy{
		ID:         oldPolicy.ID,
		Name:       newName,
		Path:       newPath,
		Org:        oldPolicy.Org,
		Urn:        auxPolicy.Urn,
		CreateAt:   oldPolicy.CreateAt,
		UpdateAt:   time.Now().UTC(),
		Statements: &newStatements,
	}

	// Update policy
	updatedPolicy, err := api.PolicyRepo.UpdatePolicy(policy)

	// Check unexpected DB error
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	LogOperation(requestInfo.RequestID, requestInfo.Identifier, fmt.Sprintf("Policy updated from %+v to %+v", oldPolicy, updatedPolicy))
	return updatedPolicy, nil
}

func (api WorkerAPI) RemovePolicy(requestInfo RequestInfo, org string, name string) error {

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

	LogOperation(requestInfo.RequestID, requestInfo.Identifier, fmt.Sprintf("Policy deleted %+v", policy))
	return nil
}

func (api WorkerAPI) ListAttachedGroups(requestInfo RequestInfo, filter *Filter) ([]PolicyGroups, int, error) {
	// Validate fields
	var total int
	orderByValidColumns := api.UserRepo.OrderByValidColumns(POLICY_ACTION_LIST_ATTACHED_GROUPS)
	err := validateFilter(filter, orderByValidColumns)
	if err != nil {
		return nil, total, err
	}

	// Call repo to retrieve the policy
	policy, err := api.GetPolicyByName(requestInfo, filter.Org, filter.PolicyName)
	if err != nil {
		return nil, total, err
	}

	// Check restrictions
	policiesFiltered, err := api.GetAuthorizedPolicies(requestInfo, policy.Urn, POLICY_ACTION_LIST_ATTACHED_GROUPS, []Policy{*policy})
	if err != nil {
		return nil, total, err
	}
	if len(policiesFiltered) < 1 {
		return nil, total, &Error{
			Code: UNAUTHORIZED_RESOURCES_ERROR,
			Message: fmt.Sprintf("User with externalId %v is not allowed to access to resource %v",
				requestInfo.Identifier, policy.Urn),
		}
	}

	// Call repo to retrieve the attached groups
	attachedGroups, total, err := api.PolicyRepo.GetAttachedGroups(policy.ID, filter)

	// Error handling
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, total, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	groups := []PolicyGroups{}
	if attachedGroups != nil {
		groups = make([]PolicyGroups, len(attachedGroups), cap(attachedGroups))
		for i, m := range attachedGroups {
			groups[i] = PolicyGroups{
				Group:    m.GetGroup().Name,
				CreateAt: m.GetDate(),
			}
		}
	}

	return groups, total, nil
}

// PRIVATE HELPER METHODS

func createPolicy(name string, path string, org string, statements *[]Statement) Policy {
	urn := CreateUrn(org, RESOURCE_POLICY, path, name)
	policy := Policy{
		ID:         uuid.NewV4().String(),
		Name:       name,
		Path:       path,
		CreateAt:   time.Now().UTC(),
		UpdateAt:   time.Now().UTC(),
		Org:        org,
		Urn:        urn,
		Statements: statements,
	}

	return policy
}
