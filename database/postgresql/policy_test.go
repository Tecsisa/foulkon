package postgresql

import (
	"testing"
	"time"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/database"
	"github.com/stretchr/testify/assert"
)

func TestPostgresRepo_AddPolicy(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		previousPolicy *Policy
		statements     []Statement
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
				UpdateAt: now,
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
				UpdateAt: now,
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
			previousPolicy: &Policy{
				ID:       "test1",
				Name:     "test",
				Org:      "123",
				Path:     "/path/",
				CreateAt: now.UnixNano(),
				UpdateAt: now.UnixNano(),
				Urn:      api.CreateUrn("123", api.RESOURCE_POLICY, "/path/", "test"),
			},
			statements: []Statement{
				{
					ID:        "test1",
					PolicyID:  "test1",
					Effect:    "allow",
					Actions:   api.USER_ACTION_GET_USER,
					Resources: api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
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
		cleanPolicyTable(t, n)
		cleanStatementTable(t, n)

		// Call to repository to add a policy
		if test.previousPolicy != nil {
			insertPolicy(t, n, *test.previousPolicy, test.statements)
		}
		receivedPolicy, err := repoDB.AddPolicy(test.policy)
		if test.expectedError != nil {
			dbError, _ := err.(*database.Error)
			assert.Equal(t, test.expectedError, dbError, "Error in test case %v", n)
		} else {
			assert.Nil(t, err, "Error in test case %v", n)
			// Check response
			assert.Equal(t, test.expectedResponse, receivedPolicy, "Error in test case %v", n)
			// Check database
			policyNumber := getPoliciesCountFiltered(t, n, test.policy.ID, test.policy.Org, test.policy.Name, test.policy.Path, test.policy.CreateAt.UnixNano(), test.policy.Urn)
			assert.Equal(t, 1, policyNumber, "Error in test case %v", n)

			for _, statement := range *test.policy.Statements {
				statementNumber := getStatementsCountFiltered(
					t,
					n,
					"",
					"",
					statement.Effect,
					stringArrayToString(statement.Actions),
					stringArrayToString(statement.Resources))
				assert.Equal(t, 1, statementNumber, "Error in test case %v", n)
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
				UpdateAt: now.UnixNano(),
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
				UpdateAt: now,
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
		cleanPolicyTable(t, n)
		cleanStatementTable(t, n)

		// Insert previous data
		if test.policy != nil {
			insertPolicy(t, n, *test.policy, test.statements)
		}
		// Call to repository to get a policy
		receivedPolicy, err := repoDB.GetPolicyByName(test.org, test.name)
		if test.expectedError != nil {
			dbError, _ := err.(*database.Error)
			assert.Equal(t, test.expectedError, dbError, "Error in test case %v", n)
		} else {
			assert.Nil(t, err, "Error in test case %v", n)
			// Check response
			assert.Equal(t, test.expectedResponse, receivedPolicy, "Error in test case %v", n)
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
				UpdateAt: now.UnixNano(),
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
				UpdateAt: now,
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
		cleanPolicyTable(t, n)
		cleanStatementTable(t, n)

		// Insert previous data
		if test.policy != nil {
			insertPolicy(t, n, *test.policy, test.statements)
		}
		// Call to repository to get a policy
		receivedPolicy, err := repoDB.GetPolicyById(test.id)
		if test.expectedError != nil {
			dbError, _ := err.(*database.Error)
			assert.Equal(t, test.expectedError, dbError, "Error in test case %v", n)
		} else {
			assert.Nil(t, err, "Error in test case %v", n)
			// Check response
			assert.Equal(t, test.expectedResponse, receivedPolicy, "Error in test case %v", n)
		}

	}
}

func TestPostgresRepo_GetPoliciesFiltered(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		policies   []Policy
		statements []Statement
		// Postgres Repo Args
		filter *api.Filter
		// Expected result
		expectedResponse []api.Policy
	}{
		"OkCase": {
			filter: &api.Filter{
				PathPrefix: "/",
				Org:        "org1",
				Offset:     0,
				Limit:      20,
				OrderBy:    "name desc",
			},
			policies: []Policy{
				{
					ID:       "111",
					Name:     "test1",
					Org:      "org1",
					Path:     "/path1/",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path1/", "test1"),
				}, {
					ID:       "222",
					Name:     "test2",
					Org:      "org1",
					Path:     "/path2/",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path2/", "test2"),
				},
			},
			statements: []Statement{
				{
					ID:        "1",
					Effect:    "allow",
					PolicyID:  "111",
					Actions:   api.USER_ACTION_GET_USER,
					Resources: api.GetUrnPrefix("", api.RESOURCE_USER, "/path1/"),
				},
				{
					ID:        "2",
					Effect:    "allow",
					PolicyID:  "222",
					Actions:   api.USER_ACTION_GET_USER,
					Resources: api.GetUrnPrefix("", api.RESOURCE_USER, "/path2/"),
				},
			},
			expectedResponse: []api.Policy{
				{
					ID:       "222",
					Name:     "test2",
					Org:      "org1",
					Path:     "/path2/",
					CreateAt: now,
					UpdateAt: now,
					Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path2/", "test2"),
					Statements: &[]api.Statement{
						{
							Effect: "allow",
							Actions: []string{
								api.USER_ACTION_GET_USER,
							},
							Resources: []string{
								api.GetUrnPrefix("", api.RESOURCE_USER, "/path2/"),
							},
						},
					},
				},
				{
					ID:       "111",
					Name:     "test1",
					Org:      "org1",
					Path:     "/path1/",
					CreateAt: now,
					UpdateAt: now,
					Urn:      api.CreateUrn("org1", api.RESOURCE_POLICY, "/path1/", "test1"),
					Statements: &[]api.Statement{
						{
							Effect: "allow",
							Actions: []string{
								api.USER_ACTION_GET_USER,
							},
							Resources: []string{
								api.GetUrnPrefix("", api.RESOURCE_USER, "/path1/"),
							},
						},
					},
				},
			},
		},
		"OKCaseNotFound": {
			filter: &api.Filter{
				PathPrefix: "test",
				Org:        "org1",
				Offset:     0,
				Limit:      20,
			},
			expectedResponse: []api.Policy{},
		},
	}

	for n, test := range testcases {
		// Clean policy database
		cleanPolicyTable(t, n)
		cleanStatementTable(t, n)

		// Insert previous data
		for i, policy := range test.policies {
			var statement []Statement = []Statement{test.statements[i]}
			insertPolicy(t, n, policy, statement)
		}
		// Call to repository to get a policy
		receivedPolicy, total, err := repoDB.GetPoliciesFiltered(test.filter)
		assert.Nil(t, err, "Error in test case %v", n)

		// Check total
		assert.Equal(t, len(test.expectedResponse), total, "Error in test case %v", n)

		// Check response
		assert.Equal(t, test.expectedResponse, receivedPolicy, "Error in test case %v", n)
	}
}

func TestPostgresRepo_UpdatePolicy(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		previousPolicies   []Policy
		previousStatements []Statement
		policy             *api.Policy
		// Expected result
		expectedResponse *api.Policy
	}{
		"OkCase": {
			previousPolicies: []Policy{
				{
					ID:       "test1",
					Name:     "test",
					Org:      "123",
					Path:     "/path/",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Urn:      api.CreateUrn("123", api.RESOURCE_POLICY, "/path/", "test"),
				},
			},
			previousStatements: []Statement{
				{
					ID:        "1",
					PolicyID:  "111",
					Effect:    "allow",
					Actions:   api.USER_ACTION_GET_USER,
					Resources: api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
				},
			},
			policy: &api.Policy{
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
		cleanPolicyTable(t, n)
		cleanStatementTable(t, n)

		// Call to repository to add a policy
		if test.previousPolicies != nil {
			for _, p := range test.previousPolicies {
				insertPolicy(t, n, p, test.previousStatements)
			}
		}
		receivedPolicy, err := repoDB.UpdatePolicy(*test.policy)
		assert.Nil(t, err, "Error in test case %v", n)

		// Check response
		assert.Equal(t, test.expectedResponse, receivedPolicy, "Error in test case %v", n)
	}
}

func TestPostgresRepo_RemovePolicy(t *testing.T) {
	type relation struct {
		policyID      string
		groupID       string
		createAt      int64
		groupNotFound bool
	}
	type policyData struct {
		policy     Policy
		statements []Statement
	}
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousPolicies []policyData
		relations        []relation
		// Postgres Repo Args
		policyToDelete string
	}{
		"OkCase": {
			previousPolicies: []policyData{
				{
					policy: Policy{
						ID:       "test1",
						Name:     "test1",
						Org:      "123",
						Path:     "/path/",
						CreateAt: now.UnixNano(),
						UpdateAt: now.UnixNano(),
						Urn:      api.CreateUrn("123", api.RESOURCE_POLICY, "/path/", "test1"),
					},
					statements: []Statement{
						{
							ID:        "test1",
							PolicyID:  "test1",
							Effect:    "allow",
							Actions:   api.USER_ACTION_GET_USER,
							Resources: api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
				{
					policy: Policy{
						ID:       "test2",
						Name:     "test2",
						Org:      "123",
						Path:     "/path/",
						CreateAt: now.UnixNano(),
						UpdateAt: now.UnixNano(),
						Urn:      api.CreateUrn("123", api.RESOURCE_POLICY, "/path/", "test2"),
					},
					statements: []Statement{
						{
							ID:        "test2",
							PolicyID:  "test2",
							Effect:    "allow",
							Actions:   api.USER_ACTION_GET_USER,
							Resources: api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
						},
					},
				},
			},
			relations: []relation{
				{
					policyID: "test1",
					groupID:  "GroupID",
					createAt: now.UnixNano(),
				},
				{
					policyID: "test2",
					groupID:  "GroupID2",
					createAt: now.UnixNano(),
				},
			},
			policyToDelete: "test1",
		},
	}

	for n, test := range testcases {
		// Clean policy database
		cleanPolicyTable(t, n)
		cleanStatementTable(t, n)
		cleanGroupTable(t, n)
		cleanGroupPolicyRelationTable(t, n)

		// insert previous policy
		if test.previousPolicies != nil {
			for _, p := range test.previousPolicies {
				insertPolicy(t, n, p.policy, p.statements)
			}
		}
		if test.relations != nil {
			for _, rel := range test.relations {
				insertGroupPolicyRelation(t, n, rel.groupID, rel.policyID, rel.createAt)
			}
		}
		err := repoDB.RemovePolicy(test.policyToDelete)
		assert.Nil(t, err, "Error in test case %v", n)

		// Check database
		policyNumber := getPoliciesCountFiltered(t, n, test.policyToDelete, "", "", "", 0, "")
		assert.Equal(t, 0, policyNumber, "Error in test case %v", n)

		statementNumber := getStatementsCountFiltered(
			t,
			n,
			"",
			test.policyToDelete,
			"",
			"",
			"")
		assert.Equal(t, 0, statementNumber, "Error in test case %v", n)

		// Check total policy number
		totalPolicyNumber := getPoliciesCountFiltered(t, n, "", "", "", "", 0, "")
		assert.Equal(t, 1, totalPolicyNumber, "Error in test case %v", n)

		totalStatementNumber := getStatementsCountFiltered(t, n, "", "", "", "", "")
		assert.Equal(t, 1, totalStatementNumber, "Error in test case %v", n)

		totalGroupPolicyRelationNumber := getGroupPolicyRelationCount(t, n, "", "")
		assert.Equal(t, 1, totalGroupPolicyRelationNumber, "Error in test case %v", n)
	}
}

func TestPostgresRepo_GetAttachedGroups(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		previousPolicy   *Policy
		statements       []Statement
		filter           *api.Filter
		group            []Group
		createAt         []int64
		expectedResponse []*PolicyGroup
	}{
		"OkCase": {
			previousPolicy: &Policy{
				ID:       "test1",
				Name:     "test",
				Org:      "123",
				Path:     "/path/",
				CreateAt: now.UnixNano(),
				UpdateAt: now.UnixNano(),
				Urn:      api.CreateUrn("123", api.RESOURCE_POLICY, "/path/", "test"),
			},
			statements: []Statement{
				{
					ID:        "test1",
					PolicyID:  "test1",
					Effect:    "allow",
					Actions:   api.USER_ACTION_GET_USER,
					Resources: api.GetUrnPrefix("", api.RESOURCE_USER, "/path/"),
				},
			},
			filter: &api.Filter{
				OrderBy: "create_at desc",
			},
			group: []Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path1",
					Urn:      "urn1",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Org:      "Org1",
				},
				{
					ID:       "GroupID2",
					Name:     "Name2",
					Path:     "Path2",
					Urn:      "urn2",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Org:      "Org2",
				},
			},
			createAt: []int64{now.UnixNano() - 1, now.UnixNano()},
			expectedResponse: []*PolicyGroup{
				{
					Group: &api.Group{
						ID:       "GroupID2",
						Name:     "Name2",
						Path:     "Path2",
						Urn:      "urn2",
						CreateAt: now,
						UpdateAt: now,
						Org:      "Org2",
					},
					CreateAt: now,
				},
				{
					Group: &api.Group{
						ID:       "GroupID1",
						Name:     "Name1",
						Path:     "Path1",
						Urn:      "urn1",
						CreateAt: now,
						UpdateAt: now,
						Org:      "Org1",
					},
					CreateAt: now.Add(-1),
				},
			},
		},
	}

	for n, test := range testcases {
		// Clean database
		cleanPolicyTable(t, n)
		cleanStatementTable(t, n)
		cleanGroupTable(t, n)
		cleanGroupPolicyRelationTable(t, n)

		// Call to repository to add a policy
		insertPolicy(t, n, *test.previousPolicy, test.statements)
		for i, group := range test.group {
			insertGroup(t, n, group)
			insertGroupPolicyRelation(t, n, group.ID, test.previousPolicy.ID, test.createAt[i])
		}

		groups, total, err := repoDB.GetAttachedGroups(test.previousPolicy.ID, test.filter)
		assert.Nil(t, err, "Error in test case %v", n)

		// Check total
		assert.Equal(t, len(test.expectedResponse), total, "Error in test case %v", n)

		// Check response
		for i, r := range groups {
			assert.Equal(t, test.expectedResponse[i].GetGroup(), r.GetGroup(), "Error in test case %v", n)
			assert.Equal(t, test.expectedResponse[i].GetPolicy(), r.GetPolicy(), "Error in test case %v", n)
			assert.Equal(t, test.expectedResponse[i].GetDate(), r.GetDate(), "Error in test case %v", n)
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
				UpdateAt: now.UnixNano(),
				Urn:      api.CreateUrn("123", api.RESOURCE_POLICY, "/path/", "test"),
			},
			apiPolicy: &api.Policy{
				ID:       "test1",
				Name:     "test",
				Org:      "123",
				Path:     "/path/",
				CreateAt: now,
				UpdateAt: now,
				Urn:      api.CreateUrn("123", api.RESOURCE_POLICY, "/path/", "test"),
			},
		},
	}

	for n, test := range testcases {
		receivedAPIPolicy := dbPolicyToAPIPolicy(test.dbPolicy)
		// Check response
		assert.Equal(t, test.apiPolicy, receivedAPIPolicy, "Error in test case %v", n)
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
		assert.Equal(t, test.apiStatements, receivedAPIStatements, "Error in test case %v", n)
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
		assert.Equal(t, test.expectedString, receivedString, "Error in test case %v", n)
	}
}
