package postgresql

import (
	"fmt"
	"strings"
	"time"

	"github.com/satori/go.uuid"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/database"
)

func (p PostgresRepo) GetPolicyById(id string) (*api.Policy, error) {
	policy := &Policy{}
	query := p.Dbmap.Where("id like ?", id).First(&policy)

	// Check if policy exist
	if query.RecordNotFound() {
		return nil, &database.Error{
			Code:    database.POLICY_NOT_FOUND,
			Message: fmt.Sprintf("Policy with id %v not found", id),
		}
	}

	// Error Handling
	if err := query.Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Retrieve associated statements
	statements := []Statement{}
	query = p.Dbmap.Where("policy_id like ?", policy.ID).Find(&statements)
	// Error Handling
	if err := query.Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Create API policy
	policyApi := policyDBToPolicyAPI(policy)
	policyApi.Statements = statementsDBToStatetmentsAPI(statements)

	// Return group
	return policyApi, nil
}

func (p PostgresRepo) GetPolicyByName(org string, name string) (*api.Policy, error) {
	policy := &Policy{}
	query := p.Dbmap.Where("org like ? AND name like ?", org, name).First(policy)

	// Check if policy exist
	if query.RecordNotFound() {
		return nil, &database.Error{
			Code:    database.POLICY_NOT_FOUND,
			Message: fmt.Sprintf("Policy with organization %v and name %v not found", org, name),
		}
	}

	// Error Handling
	if err := query.Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Retrieve associated statements
	statements := []Statement{}
	query = p.Dbmap.Where("policy_id like ?", policy.ID).Find(&statements)
	// Error Handling
	if err := query.Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Create API policy
	policyApi := policyDBToPolicyAPI(policy)
	policyApi.Statements = statementsDBToStatetmentsAPI(statements)

	// Return policy
	return policyApi, nil
}

func (p PostgresRepo) AddPolicy(policy api.Policy) (*api.Policy, error) {
	// Create policy model
	policyDB := &Policy{
		ID:       policy.ID,
		Name:     policy.Name,
		Path:     policy.Path,
		CreateAt: time.Now().UTC().UnixNano(),
		Urn:      policy.Urn,
		Org:      policy.Org,
	}

	transaction := p.Dbmap.Begin()

	// Create policy
	transaction.Create(policyDB)

	// Create statements
	for _, statementApi := range *policy.Statements {
		// Create statement model
		statementDB := &Statement{
			ID:        uuid.NewV4().String(),
			PolicyID:  policy.ID,
			Effect:    statementApi.Effect,
			Action:    stringArrayToSplitedString(statementApi.Action),
			Resources: stringArrayToSplitedString(statementApi.Resources),
		}
		transaction.Create(statementDB)
	}

	// Error handling
	if err := transaction.Error; err != nil {
		transaction.Rollback()
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	} else {
		transaction.Commit()
	}

	// Create API policy
	policyApi := policyDBToPolicyAPI(policyDB)
	policyApi.Statements = policy.Statements

	return policyApi, nil
}

func (p PostgresRepo) GetPoliciesFiltered(org string, pathPrefix string) ([]api.Policy, error) {
	policies := []Policy{}
	query := p.Dbmap
	if len(org) > 0 {
		query = query.Where("org like ?", org)
	}
	if len(pathPrefix) > 0 {
		query = query.Where("path like ?", pathPrefix+"%")
	}

	// Error handling
	if err := query.Find(&policies).Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Transform policies for API
	if policies != nil {
		apiPolicies := make([]api.Policy, len(policies), cap(policies))

		for i, pol := range policies {
			policy := policyDBToPolicyAPI(&pol)

			// Retrieve associated statements
			statements := []Statement{}
			query = p.Dbmap.Where("policy_id like ?", policy.ID).Find(&statements)
			// Error Handling
			if err := query.Error; err != nil {
				return nil, &database.Error{
					Code:    database.INTERNAL_ERROR,
					Message: err.Error(),
				}
			}

			policy.Statements = statementsDBToStatetmentsAPI(statements)

			// Assign policy
			apiPolicies[i] = *policy
		}

		return apiPolicies, nil
	}

	// No data to return
	return nil, nil
}

func (p PostgresRepo) UpdatePolicy(policy api.Policy, name string, path string, urn string, statements []api.Statement) (*api.Policy, error) {
	// Create policy to update
	policyUpdated := Policy{
		Name: name,
		Path: path,
		Urn:  urn,
	}

	policyDB := Policy{
		ID:       policy.ID,
		Name:     policy.Name,
		Path:     policy.Path,
		CreateAt: policy.CreateAt.UTC().UnixNano(),
		Urn:      policy.Urn,
		Org:      policy.Org,
	}

	transaction := p.Dbmap.Begin()

	// Update policy
	transaction.Model(&policyDB).Update(policyUpdated)

	// Check if policy exist
	if transaction.RecordNotFound() {
		transaction.Rollback()
		return nil, &database.Error{
			Code:    database.POLICY_NOT_FOUND,
			Message: fmt.Sprintf("Policy with name %v not found", name),
		}
	}

	// Clear old statements
	transaction.Where("policy_id like ?", policy.ID).Delete(Statement{})

	// Create new statements
	for _, s := range statements {
		statementDB := &Statement{
			ID:        uuid.NewV4().String(),
			PolicyID:  policy.ID,
			Effect:    s.Effect,
			Action:    stringArrayToSplitedString(s.Action),
			Resources: stringArrayToSplitedString(s.Resources),
		}
		transaction.Create(statementDB)
	}

	// Error Handling
	if err := transaction.Error; err != nil {
		transaction.Rollback()
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	transaction.Commit()

	// Create API policy
	policyApi := policyDBToPolicyAPI(&policyDB)
	policyApi.Statements = &statements

	return policyApi, nil
}

func (p PostgresRepo) RemovePolicy(id string) error {

	transaction := p.Dbmap.Begin()

	// Delete policy relations (group)
	transaction.Where("policy_id like ?", id).Delete(&GroupPolicyRelation{})
	if err := transaction.Error; err != nil {
		transaction.Rollback()
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}
	// Delete policy statements
	transaction.Where("policy_id like ?", id).Delete(&Statement{})
	if err := transaction.Error; err != nil {
		transaction.Rollback()
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}
	//  Delete policy
	transaction.Where("id like ?", id).Delete(&Policy{})
	if err := transaction.Error; err != nil {
		transaction.Rollback()
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	transaction.Commit()
	return nil
}

func (p PostgresRepo) GetAllPolicyGroupRelation(policyID string) ([]api.Group, error) {
	relations := []GroupPolicyRelation{}
	query := p.Dbmap.Where("policy_id like ?", policyID).Find(&relations)

	// Error Handling
	if err := query.Error; err != nil {
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Transform relations to API domain
	if relations != nil {
		groups := make([]api.Group, len(relations), cap(relations))
		for i, r := range relations {
			group, err := p.GetGroupById(r.GroupID)
			// Error handling
			if err != nil {
				return nil, &database.Error{
					Code:    database.INTERNAL_ERROR,
					Message: err.Error(),
				}
			}

			groups[i] = *group
		}

		return groups, nil
	}

	return nil, nil
}

// Private helper methods

// Transform a policy retrieved from db into a policy for API
func policyDBToPolicyAPI(policydb *Policy) *api.Policy {
	return &api.Policy{
		ID:       policydb.ID,
		Name:     policydb.Name,
		Path:     policydb.Path,
		CreateAt: time.Unix(0, policydb.CreateAt).UTC(),
		Urn:      policydb.Urn,
		Org:      policydb.Org,
	}
}

// Transform a list of statements from db into API statements
func statementsDBToStatetmentsAPI(statements []Statement) *[]api.Statement {
	statementsApi := make([]api.Statement, len(statements), cap(statements))
	for i, s := range statements {
		statementsApi[i] = api.Statement{
			Action:    strings.Split(s.Action, ";"),
			Effect:    s.Effect,
			Resources: strings.Split(s.Resources, ";"),
		}
	}

	return &statementsApi
}

// Transform an array of strings into a separated string values
func stringArrayToSplitedString(array []string) string {
	stringVal := ""
	for _, s := range array {
		if len(stringVal) == 0 {
			stringVal = s
		} else {
			stringVal = stringVal + ";" + s
		}
	}

	return stringVal
}
