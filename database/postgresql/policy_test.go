package postgresql

import (
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/database"
)

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
					{
						Effect: "allow",
						Actions: []string{
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
					{
						Effect: "allow",
						Actions: []string{
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
					{
						Effect: "allow",
						Actions: []string{
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
					{
						Effect: "allow",
						Actions: []string{
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
				t.Errorf("Test %v failed. Received different policies number: %v", n, policyNumber)
				continue
			}
			for _, statement := range *test.policy.Statements {
				statementNumber, err := getStatementsCountFiltered(
					"",
					"",
					statement.Effect,
					stringArrayToString(statement.Actions),
					stringArrayToString(statement.Resources))
				if err != nil {
					t.Errorf("Test %v failed. Unexpected error counting statements: %v", n, err)
					continue
				}
				if statementNumber != 1 {
					t.Errorf("Test %v failed. Received different statements number: %v", n, statementNumber)
					continue
				}
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
				{
					ID:        "0123",
					Effect:    "allow",
					PolicyID:  "1234",
					Actions:   api.USER_ACTION_GET_USER,
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
					{
						Effect: "allow",
						Actions: []string{
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
				{
					ID:        "0123",
					Effect:    "allow",
					PolicyID:  "1234",
					Actions:   api.USER_ACTION_GET_USER,
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
					{
						Effect: "allow",
						Actions: []string{
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

func TestPostgresRepo_GetPoliciesFiltered(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		policy     *Policy
		statements []Statement
		// Postgres Repo Args
		org        string
		pathPrefix string
		// Expected result
		expectedResponse []api.Policy
	}{
		"OkCase": {
			org:        "org1",
			pathPrefix: "/path/",
			policy: &Policy{
				ID:       "1234",
				Name:     "test",
				Org:      "org1",
				Path:     "/path/",
				CreateAt: now.UnixNano(),
				Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path/", "test"),
			},
			statements: []Statement{
				{
					ID:        "0123",
					Effect:    "allow",
					PolicyID:  "1234",
					Actions:   api.USER_ACTION_GET_USER,
					Resources: api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
				},
			},
			expectedResponse: []api.Policy{
				{
					ID:       "1234",
					Name:     "test",
					Org:      "org1",
					Path:     "/path/",
					CreateAt: now,
					Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path/", "test"),
					Statements: &[]api.Statement{
						{
							Effect: "allow",
							Actions: []string{
								api.USER_ACTION_GET_USER,
							},
							Resources: []string{
								api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
							},
						},
					},
				},
			},
		},
		"OKCaseNotFound": {
			org:              "org1",
			pathPrefix:       "test",
			expectedResponse: []api.Policy{},
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
		receivedPolicy, err := repoDB.GetPoliciesFiltered(test.org, test.pathPrefix)
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

func TestPostgresRepo_UpdatePolicy(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		previousPolicy *api.Policy
		policy         api.Policy
		name           string
		path           string
		urn            string
		statements     []api.Statement
		// Expected result
		expectedResponse *api.Policy
	}{
		"OkCase": {
			previousPolicy: &api.Policy{
				ID:       "test1",
				Name:     "test",
				Org:      "123",
				Path:     "/path/",
				CreateAt: now,
				Urn:      api.CreateUrn("123", api.RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]api.Statement{
					{
						Effect: "allow",
						Actions: []string{
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
					{
						Effect: "allow",
						Actions: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			name: "newName",
			path: "/newPath/",
			urn:  api.CreateUrn("123", api.RESOURCE_POLICY, "/newPath/", "newName"),
			statements: []api.Statement{
				{
					Effect: "allow",
					Actions: []string{
						api.USER_ACTION_GET_USER,
					},
					Resources: []string{
						api.GetUrnPrefix("123", api.RESOURCE_USER, "/newPath/"),
					},
				},
			},
			expectedResponse: &api.Policy{
				ID:       "test1",
				Name:     "newName",
				Org:      "123",
				Path:     "/newPath/",
				CreateAt: now,
				Urn:      api.CreateUrn("123", api.RESOURCE_POLICY, "/newPath/", "newName"),
				Statements: &[]api.Statement{
					{
						Effect: "allow",
						Actions: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("123", api.RESOURCE_USER, "/newPath/"),
						},
					},
				},
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
		receivedPolicy, err := repoDB.UpdatePolicy(test.policy, test.name, test.path, test.urn, test.statements)
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

func TestPostgresRepo_RemovePolicy(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		previousPolicy *api.Policy
		id             string
		group          *api.Group
	}{
		"OkCase": {
			previousPolicy: &api.Policy{
				ID:       "test1",
				Name:     "test",
				Org:      "123",
				Path:     "/path/",
				CreateAt: now,
				Urn:      api.CreateUrn("123", api.RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]api.Statement{
					{
						Effect: "allow",
						Actions: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			id: "test1",
			group: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "urn",
				CreateAt: now,
				Org:      "Org",
			},
		},
	}

	for n, test := range testcases {
		// Clean policy database
		cleanPolicyTable()
		cleanStatementTable()
		cleanGroupTable()
		cleanGroupPolicyRelationTable()

		// Call to repository to add a policy
		if test.previousPolicy != nil {
			_, err := repoDB.AddPolicy(*test.previousPolicy)
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error: %v", n, err)
				continue
			}
		}
		if test.group != nil {
			err := insertGroup(test.group.ID, test.group.Name, test.group.Path,
				test.group.CreateAt.UnixNano(), test.group.Urn, test.group.Org)
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting group: %v", n, err)
				continue
			}
			err = insertGroupPolicyRelation(test.group.ID, test.previousPolicy.ID)
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting group relation: %v", n, err)
				continue
			}
		}
		err := repoDB.RemovePolicy(test.id)
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error: %v", n, err)
			continue
		}
		// Check database
		policyNumber, err := getPoliciesCountFiltered(test.id, "", "", "", 0, "")
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error counting policies: %v", n, err)
			continue
		}
		if policyNumber != 0 {
			t.Errorf("Test %v failed. Received different policies number: %v", n, policyNumber)
			continue
		}

		statementNumber, err := getStatementsCountFiltered(
			"",
			test.id,
			"",
			"",
			"")
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error counting statements: %v", n, err)
			continue
		}
		if statementNumber != 0 {
			t.Errorf("Test %v failed. Received different statements number: %v", n, statementNumber)
			continue
		}

		groupPolicyRelationNumber, err := getGroupPolicyRelationCount(test.previousPolicy.ID, "")
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error counting relations: %v", n, err)
			continue
		}
		if groupPolicyRelationNumber != 0 {
			t.Errorf("Test %v failed. Received different relations number: %v", n, groupPolicyRelationNumber)
			continue
		}
	}
}

func TestPostgresRepo_GetAttachedGroups(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		previousPolicy   *api.Policy
		group            *api.Group
		expectedResponse []api.Group
	}{
		"OkCase": {
			previousPolicy: &api.Policy{
				ID:       "test1",
				Name:     "test",
				Org:      "123",
				Path:     "/path/",
				CreateAt: now,
				Urn:      api.CreateUrn("123", api.RESOURCE_POLICY, "/path/", "test"),
				Statements: &[]api.Statement{
					{
						Effect: "allow",
						Actions: []string{
							api.USER_ACTION_GET_USER,
						},
						Resources: []string{
							api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			group: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "urn",
				CreateAt: now,
				Org:      "Org",
			},
			expectedResponse: []api.Group{
				{
					ID:       "GroupID",
					Name:     "Name",
					Path:     "Path",
					Urn:      "urn",
					CreateAt: now,
					Org:      "Org",
				},
			},
		},
	}

	for n, test := range testcases {
		// Clean database
		cleanPolicyTable()
		cleanStatementTable()
		cleanGroupTable()
		cleanGroupPolicyRelationTable()

		// Call to repository to add a policy
		_, err := repoDB.AddPolicy(*test.previousPolicy)
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error: %v", n, err)
			continue
		}
		if test.group != nil {
			err := insertGroup(test.group.ID, test.group.Name, test.group.Path,
				test.group.CreateAt.UnixNano(), test.group.Urn, test.group.Org)
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting group: %v", n, err)
				continue
			}
			err = insertGroupPolicyRelation(test.group.ID, test.previousPolicy.ID)
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting group relation: %v", n, err)
				continue
			}
		}

		groups, err := repoDB.GetAttachedGroups(test.previousPolicy.ID)
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error: %v", n, err)
			continue
		}
		// Check response
		if diff := pretty.Compare(groups, test.expectedResponse); diff != "" {
			t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
			continue
		}
	}
}

func Test_dbPolicyToAPIPolicy(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		dbPolicy  *Policy
		apiPolicy *api.Policy
	}{
		"OkCase": {
			dbPolicy: &Policy{
				ID:       "test1",
				Name:     "test",
				Org:      "123",
				Path:     "/path/",
				CreateAt: now.UnixNano(),
				Urn:      api.CreateUrn("123", api.RESOURCE_POLICY, "/path/", "test"),
			},
			apiPolicy: &api.Policy{
				ID:       "test1",
				Name:     "test",
				Org:      "123",
				Path:     "/path/",
				CreateAt: now,
				Urn:      api.CreateUrn("123", api.RESOURCE_POLICY, "/path/", "test"),
			},
		},
	}

	for n, test := range testcases {
		receivedAPIPolicy := dbPolicyToAPIPolicy(test.dbPolicy)
		// Check response
		if diff := pretty.Compare(receivedAPIPolicy, test.apiPolicy); diff != "" {
			t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
			continue
		}
	}
}

func Test_dbStatementsToAPIStatements(t *testing.T) {
	testcases := map[string]struct {
		dbStatements  []Statement
		apiStatements *[]api.Statement
	}{
		"OkCase": {
			dbStatements: []Statement{
				{
					ID:        "0123",
					Effect:    "allow",
					PolicyID:  "1234",
					Actions:   api.USER_ACTION_GET_USER,
					Resources: api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
				},
			},
			apiStatements: &[]api.Statement{
				{
					Effect: "allow",
					Actions: []string{
						api.USER_ACTION_GET_USER,
					},
					Resources: []string{
						api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
					},
				},
			},
		},
		"OkCase2": {
			dbStatements: []Statement{
				{
					ID:        "0123",
					Effect:    "allow",
					PolicyID:  "1234",
					Actions:   api.USER_ACTION_GET_USER,
					Resources: api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
				},
				{
					ID:        "4321",
					Effect:    "deny",
					PolicyID:  "1234",
					Actions:   api.GROUP_ACTION_GET_GROUP + ";" + api.GROUP_ACTION_CREATE_GROUP,
					Resources: api.GetUrnPrefix("", api.RESOURCE_GROUP, "/xxx/") + ";" + api.GetUrnPrefix("", api.RESOURCE_GROUP, "/xxx2/"),
				},
			},
			apiStatements: &[]api.Statement{
				{
					Effect: "allow",
					Actions: []string{
						api.USER_ACTION_GET_USER,
					},
					Resources: []string{
						api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
					},
				},
				{
					Effect: "deny",
					Actions: []string{
						api.GROUP_ACTION_GET_GROUP,
						api.GROUP_ACTION_CREATE_GROUP,
					},
					Resources: []string{
						api.GetUrnPrefix("", api.RESOURCE_GROUP, "/xxx/"),
						api.GetUrnPrefix("", api.RESOURCE_GROUP, "/xxx2/"),
					},
				},
			},
		},
	}

	for n, test := range testcases {
		receivedAPIStatements := dbStatementsToAPIStatements(test.dbStatements)
		// Check response
		if diff := pretty.Compare(receivedAPIStatements, test.apiStatements); diff != "" {
			t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
			continue
		}
	}
}

func Test_stringArrayToString(t *testing.T) {
	testcases := map[string]struct {
		arrayString    []string
		expectedString string
	}{
		"OkCase": {
			arrayString: []string{
				"asd",
				"123",
				"456",
				"zxc",
			},
			expectedString: "asd;123;456;zxc",
		},
		"OkCase2": {
			arrayString:    []string{},
			expectedString: "",
		},
	}

	for n, test := range testcases {
		receivedString := stringArrayToString(test.arrayString)
		// Check response
		if diff := pretty.Compare(receivedString, test.expectedString); diff != "" {
			t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
			continue
		}
	}
}
