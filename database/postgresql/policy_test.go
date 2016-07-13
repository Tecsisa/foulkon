package postgresql

import (
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/database"
)

func TestPostgresRepo_GetPolicyById(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		policy     *Policy
		statements []Statement
		// Postgres Repo Args
		id string
		// Expected result
		expectedResponse *api.Policy
		expectedError    *database.Error
	}{
		"OkCase": {
			id: "1234",
			policy: &Policy{
				ID:       "1234",
				Name:     "test",
				Org:      "org1",
				Path:     "/path/",
				CreateAt: now.UnixNano(),
				Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path/", "test"),
			},
			statements: []Statement{
				Statement{
					ID:        "0123",
					Effect:    "allow",
					PolicyID:  "1234",
					Action:    api.USER_ACTION_GET_USER,
					Resources: api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
				},
			},
			expectedResponse: &api.Policy{
				ID:       "1234",
				Name:     "test",
				Org:      "org1",
				Path:     "/path/",
				CreateAt: now,
				Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]api.Statement{
					api.Statement{
						Effect: "allow",
						Action: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
		},
		"ErrorCaseNotFound": {
			id: "1234",
			expectedError: &database.Error{
				Code:    database.POLICY_NOT_FOUND,
				Message: "Policy with id 1234 not found",
			},
		},
	}

	for n, test := range testcases {
		// Clean policy database
		cleanPolicyTable()
		cleanStatementTable()

		// Insert previous data
		if test.policy != nil {
			err := insertPolicy(test.policy.ID, test.policy.Name, test.policy.Org, test.policy.Path, test.policy.CreateAt, test.policy.Urn, test.statements)
			if err != nil {
				t.Errorf("Test %v failed. Error inserting policy/statements: %v", n, err)
			}
		}
		// Call to repository to get a policy
		receivedPolicy, err := repoDB.GetPolicyById(test.id)
		if test.expectedError != nil {
			dbError, ok := err.(*database.Error)
			if !ok || dbError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", n, err)
				continue
			}
			if diff := pretty.Compare(dbError, test.expectedError); diff != "" {
				t.Errorf("Test %v failed. Received different error response (received/wanted) %v", n, diff)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error: %v", n, err)
				continue
			}
			// Check response
			if diff := pretty.Compare(receivedPolicy, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
		}

	}
}

func TestPostgresRepo_GetPolicyByName(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		policy     *Policy
		statements []Statement
		// Postgres Repo Args
		org  string
		name string
		// Expected result
		expectedResponse *api.Policy
		expectedError    *database.Error
	}{
		"OkCase": {
			org:  "org1",
			name: "test",
			policy: &Policy{
				ID:       "1234",
				Name:     "test",
				Org:      "org1",
				Path:     "/path/",
				CreateAt: now.UnixNano(),
				Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path/", "test"),
			},
			statements: []Statement{
				Statement{
					ID:        "0123",
					Effect:    "allow",
					PolicyID:  "1234",
					Action:    api.USER_ACTION_GET_USER,
					Resources: api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
				},
			},
			expectedResponse: &api.Policy{
				ID:       "1234",
				Name:     "test",
				Org:      "org1",
				Path:     "/path/",
				CreateAt: now,
				Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]api.Statement{
					api.Statement{
						Effect: "allow",
						Action: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
		},
		"ErrorCaseNotFound": {
			org:  "org1",
			name: "test",
			expectedError: &database.Error{
				Code:    database.POLICY_NOT_FOUND,
				Message: "Policy with organization org1 and name test not found",
			},
		},
	}

	for n, test := range testcases {
		// Clean policy database
		cleanPolicyTable()
		cleanStatementTable()

		// Insert previous data
		if test.policy != nil {
			err := insertPolicy(test.policy.ID, test.policy.Name, test.policy.Org, test.policy.Path, test.policy.CreateAt, test.policy.Urn, test.statements)
			if err != nil {
				t.Errorf("Test %v failed. Error inserting policy/statements: %v", n, err)
			}
		}
		// Call to repository to get a policy
		receivedPolicy, err := repoDB.GetPolicyByName(test.org, test.name)
		if test.expectedError != nil {
			dbError, ok := err.(*database.Error)
			if !ok || dbError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", n, err)
				continue
			}
			if diff := pretty.Compare(dbError, test.expectedError); diff != "" {
				t.Errorf("Test %v failed. Received different error response (received/wanted) %v", n, diff)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error: %v", n, err)
				continue
			}
			// Check response
			if diff := pretty.Compare(receivedPolicy, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
		}

	}
}

func TestPostgresRepo_AddPolicy(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		previousPolicy *api.Policy
		policy         api.Policy
		// Expected result
		expectedResponse *api.Policy
		expectedError    *database.Error
	}{
		"OkCase": {
			policy: api.Policy{
				ID:       "test1",
				Name:     "test",
				Org:      "123",
				Path:     "/path/",
				CreateAt: now,
				Urn:      api.CreateUrn("123", api.RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]api.Statement{
					api.Statement{
						Effect: "allow",
						Action: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			expectedResponse: &api.Policy{
				ID:       "test1",
				Name:     "test",
				Org:      "123",
				Path:     "/path/",
				CreateAt: now,
				Urn:      api.CreateUrn("123", api.RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]api.Statement{
					api.Statement{
						Effect: "allow",
						Action: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
		},
		"ErrorCaseAlreadyExists": {
			previousPolicy: &api.Policy{
				ID:       "test1",
				Name:     "test",
				Org:      "123",
				Path:     "/path/",
				CreateAt: now,
				Urn:      api.CreateUrn("123", api.RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]api.Statement{
					api.Statement{
						Effect: "allow",
						Action: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			policy: api.Policy{
				ID:       "test1",
				Name:     "test",
				Org:      "123",
				Path:     "/path/",
				CreateAt: now,
				Urn:      api.CreateUrn("123", api.RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]api.Statement{
					api.Statement{
						Effect: "allow",
						Action: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			expectedError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "pq: duplicate key value violates unique constraint \"policies_pkey\"",
			},
		},
	}

	for n, test := range testcases {
		// Clean policy database
		cleanPolicyTable()
		cleanStatementTable()

		// Call to repository to add a policy
		if test.previousPolicy != nil {
			_, err := repoDB.AddPolicy(*test.previousPolicy)
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error: %v", n, err)
				continue
			}
		}
		receivedPolicy, err := repoDB.AddPolicy(test.policy)
		if test.expectedError != nil {
			dbError, ok := err.(*database.Error)
			if !ok || dbError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", n, err)
				continue
			}
			if diff := pretty.Compare(dbError, test.expectedError); diff != "" {
				t.Errorf("Test %v failed. Received different error response (received/wanted) %v", n, diff)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error: %v", n, err)
				continue
			}
			// Check response
			if diff := pretty.Compare(receivedPolicy, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
			// Check database
			policyNumber, err := getPoliciesCountFiltered(test.policy.ID, test.policy.Org, test.policy.Name, test.policy.Path, test.policy.CreateAt.UnixNano(), test.policy.Urn)
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error counting policies: %v", n, err)
				continue
			}
			if policyNumber != 1 {
				t.Fatalf("Test %v failed. Received different policies number: %v", n, policyNumber)
				continue
			}
			for _, statement := range *test.policy.Statements {
				statementNumber, err := getStatementsCountFiltered(
					"",
					"",
					statement.Effect,
					stringArrayToString(statement.Action),
					stringArrayToString(statement.Resources))
				if err != nil {
					t.Errorf("Test %v failed. Unexpected error counting statements: %v", n, err)
					continue
				}
				if statementNumber != 1 {
					t.Fatalf("Test %v failed. Received different statements number: %v", n, statementNumber)
					continue
				}
			}
		}

	}
}
