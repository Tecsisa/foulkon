package postgresql

import (
	"fmt"
	"strings"
	"time"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/database"
	"github.com/satori/go.uuid"
)

// POLICY REPOSITORY IMPLEMENTATION

func (p PostgresRepo) AddPolicy(policy api.Policy) (*api.Policy, error) {
	// Create policy model
	policyDB := &Policy{
		ID:       policy.ID,
		Name:     policy.Name,
		Path:     policy.Path,
		CreateAt: policy.CreateAt.UnixNano(),
		UpdateAt: policy.UpdateAt.UnixNano(),
		Urn:      policy.Urn,
		Org:      policy.Org,
	}

	transaction := p.Dbmap.Begin()

	// Create policy
	if err := transaction.Create(policyDB).Error; err != nil {
		transaction.Rollback()
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Create statements
	for _, statementApi := range *policy.Statements {
		// Create statement model
		statementDB := &Statement{
			ID:        uuid.NewV4().String(),
			PolicyID:  policy.ID,
			Effect:    statementApi.Effect,
			Actions:   stringArrayToString(statementApi.Actions),
			Resources: stringArrayToString(statementApi.Resources),
		}
		if err := transaction.Create(statementDB).Error; err != nil {
			transaction.Rollback()
			return nil, &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: err.Error(),
			}
		}
	}

	transaction.Commit()

	// Create API policy
	policyApi := dbPolicyToAPIPolicy(policyDB)
	policyApi.Statements = policy.Statements

	return policyApi, nil
}

func (p PostgresRepo) GetPolicyByName(org string, name string) (*api.Policy, error) {
	policy := &Policy{}
	query := p.Dbmap.Where("org like ? AND name like ?", org, name).First(policy)

	// Check if policy exists
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
	policyApi := dbPolicyToAPIPolicy(policy)
	policyApi.Statements = dbStatementsToAPIStatements(statements)

	return policyApi, nil
}

func (p PostgresRepo) GetPolicyById(id string) (*api.Policy, error) {
	policy := &Policy{}
	query := p.Dbmap.Where("id like ?", id).First(&policy)

	// Check if policy exists
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
	policyApi := dbPolicyToAPIPolicy(policy)
	policyApi.Statements = dbStatementsToAPIStatements(statements)

	return policyApi, nil
}

func (p PostgresRepo) GetPoliciesFiltered(filter *api.Filter) ([]api.Policy, int, error) {
	var total int
	policies := []Policy{}
	query := p.Dbmap

	if len(filter.Org) > 0 {
		query = query.Where("org like ?", filter.Org)
	}
	if len(filter.PathPrefix) > 0 {
		query = query.Where("path like ?", filter.PathPrefix+"%")
	}

	// Error handling
	if err := query.Find(&policies).Count(&total).Offset(filter.Offset).Limit(filter.Limit).Find(&policies).Error; err != nil {
		return nil, total, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Transform policies for API
	var apiPolicies []api.Policy
	if policies != nil {
		apiPolicies = make([]api.Policy, len(policies), cap(policies))

		for i, pol := range policies {
			policy := dbPolicyToAPIPolicy(&pol)

			// Retrieve associated statements
			statements := []Statement{}
			query = p.Dbmap.Where("policy_id like ?", policy.ID).Find(&statements)
			// Error Handling
			if err := query.Error; err != nil {
				return nil, total, &database.Error{
					Code:    database.INTERNAL_ERROR,
					Message: err.Error(),
				}
			}

			policy.Statements = dbStatementsToAPIStatements(statements)

			// Assign policy
			apiPolicies[i] = *policy
		}

	}

	return apiPolicies, total, nil
}

func (p PostgresRepo) UpdatePolicy(policy api.Policy) (*api.Policy, error) {

	policyDB := Policy{
		ID:       policy.ID,
		Name:     policy.Name,
		Path:     policy.Path,
		CreateAt: policy.CreateAt.UTC().UnixNano(),
		UpdateAt: policy.UpdateAt.UTC().UnixNano(),
		Urn:      policy.Urn,
		Org:      policy.Org,
	}

	transaction := p.Dbmap.Begin()

	// Update policy
	if err := transaction.Model(&Policy{ID: policy.ID}).Update(policyDB).Error; err != nil {
		transaction.Rollback()
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Clear old statements
	if err := transaction.Where("policy_id like ?", policy.ID).Delete(Statement{}).Error; err != nil {
		transaction.Rollback()
		return nil, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	// Create new statements
	for _, s := range *policy.Statements {
		statementDB := &Statement{
			ID:        uuid.NewV4().String(),
			PolicyID:  policy.ID,
			Effect:    s.Effect,
			Actions:   stringArrayToString(s.Actions),
			Resources: stringArrayToString(s.Resources),
		}
		if err := transaction.Create(statementDB).Error; err != nil {
			transaction.Rollback()
			return nil, &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: err.Error(),
			}
		}
	}

	transaction.Commit()

	return &policy, nil
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

func (p PostgresRepo) GetAttachedGroups(policyID string, filter *api.Filter) ([]api.PolicyGroupRelation, int, error) {
	var total int
	relations := []GroupPolicyRelation{}
	query := p.Dbmap.Where("policy_id like ?", policyID).Find(&relations).Count(&total).Offset(filter.Offset).Limit(filter.Limit).Find(&relations)

	// Error Handling
	if err := query.Error; err != nil {
		return nil, total, &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	var groups []api.PolicyGroupRelation
	// Transform relations to API domain
	if relations != nil {
		groups = make([]api.PolicyGroupRelation, len(relations), cap(relations))
		for i, r := range relations {
			group, err := p.GetGroupById(r.GroupID)
			// Error handling
			if err != nil {
				return nil, total, &database.Error{
					Code:    database.INTERNAL_ERROR,
					Message: err.Error(),
				}
			}

			groups[i] = PolicyGroup{
				Group:    group,
				CreateAt: time.Unix(0, r.CreateAt).UTC(),
			}
		}
	}

	return groups, total, nil
}

// PRIVATE HELPER METHODS

// Transform a policy retrieved from db into a policy for API
func dbPolicyToAPIPolicy(policydb *Policy) *api.Policy {
	return &api.Policy{
		ID:       policydb.ID,
		Name:     policydb.Name,
		Path:     policydb.Path,
		CreateAt: time.Unix(0, policydb.CreateAt).UTC(),
		UpdateAt: time.Unix(0, policydb.UpdateAt).UTC(),
		Urn:      policydb.Urn,
		Org:      policydb.Org,
	}
}

// Transform a list of statements from db into API statements
func dbStatementsToAPIStatements(statements []Statement) *[]api.Statement {
	statementsApi := make([]api.Statement, len(statements), cap(statements))
	for i, s := range statements {
		statementsApi[i] = api.Statement{
			Actions:   strings.Split(s.Actions, ";"),
			Effect:    s.Effect,
			Resources: strings.Split(s.Resources, ";"),
		}
	}

	return &statementsApi
}

// Transform an array of strings into a semicolon-separated string
func stringArrayToString(array []string) string {
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
