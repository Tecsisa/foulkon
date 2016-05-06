package api

import (
	"fmt"
	"github.com/tecsisa/authorizr/database"
	"time"
)

// Policy domain
type Policy struct {
	ID         string       `json:"ID, omitempty"`
	Name       string       `json:"Name, omitempty"`
	Path       string       `json:"Path, omitempty"`
	Org        string       `json:"Org, omitempty"`
	CreateAt   time.Time    `json:"CreateAt, omitempty"`
	Urn        string       `json:"Urn, omitempty"`
	Statements *[]Statement `json:"Statements, omitempty"`
}

type Statement struct {
	Effect    string   `json:"Effect, omitempty"`
	Action    []string `json:"Action, omitempty"`
	Resources []string `json:"Resources, omitempty"`
}

type PoliciesAPI struct {
	PolicyRepo PolicyRepo
}

func (p *PoliciesAPI) GetPolicies(path string) string {
	return path
}

func (p *PoliciesAPI) AddPolicy(policy Policy) (*Policy, error) {
	// Check if policy already exist
	policyDB, err := p.PolicyRepo.GetPolicyByName(policy.Org, policy.Name)

	// If policy exist it can't create it
	if policyDB != nil {
		return nil, &Error{
			Code:    POLICY_ALREADY_EXIST,
			Message: fmt.Sprintf("Unable to create policy, policy with org %v and name %v already exist", policy.Org, policy.Name),
		}
	}

	// Create policy
	policyCreated, err := p.PolicyRepo.AddPolicy(policy)

	// Check if there is an unexpected error in DB
	if err != nil {
		//Transform to DB error
		dbError := err.(*database.Error)
		return nil, &Error{
			Code:    UNKNOWN_API_ERROR,
			Message: dbError.Message,
		}
	}

	// Return policy created
	return policyCreated, nil
}
