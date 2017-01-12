package postgresql

import (
	"fmt"
	"os"
	"testing"
	"time"

	"errors"

	"github.com/Tecsisa/foulkon/api"
	"github.com/stretchr/testify/assert"
)

var repoDB PostgresRepo
var testFilter = &api.Filter{
	PathPrefix: "",
	Offset:     0,
	Limit:      0,
}

func TestMain(m *testing.M) {
	// Wait for DB
	time.Sleep(3 * time.Second)
	// Retrieve db connector to run test
	dbmap, err := InitDb("postgres://postgres:password@localhost:54320/postgres?sslmode=disable", "5", "20", "300")
	if err != nil {
		fmt.Fprintln(os.Stderr, "There was an error starting connector", err)
		os.Exit(1)
	}
	repoDB = PostgresRepo{
		Dbmap: dbmap,
	}

	result := m.Run()

	os.Exit(result)
}

func TestInitDb(t *testing.T) {
	// Retrieve db connector to run test
	testcases := map[string]struct {
		idle          string
		max           string
		ttl           string
		expectedError error
	}{
		"ErrorCaseInvalidIdleParam": {
			idle:          "asd",
			max:           "20",
			ttl:           "300",
			expectedError: errors.New("Invalid postgresql idleConns param: asd"),
		},
		"ErrorCaseInvalidMaxParam": {
			idle:          "5",
			max:           "asd",
			ttl:           "300",
			expectedError: errors.New("Invalid postgresql maxOpenConns param: asd"),
		},
		"ErrorCaseInvalidTTLParam": {
			idle:          "5",
			max:           "20",
			ttl:           "asd",
			expectedError: errors.New("Invalid postgresql connTTL param: asd"),
		},
	}

	for n, test := range testcases {
		_, err := InitDb("postgres://postgres:password@localhost:54320/postgres?sslmode=disable", test.idle, test.max, test.ttl)
		assert.Equal(t, test.expectedError, err, "Error in test case %v", n)
	}

}

func TestPostgresRepo_OrderByValidColumns(t *testing.T) {
	testcases := map[string]struct {
		action          string
		expectedColumns []string
	}{
		"OkCaseAction-" + api.USER_ACTION_LIST_USERS: {
			action:          api.USER_ACTION_LIST_USERS,
			expectedColumns: []string{"path", "external_id", "create_at", "update_at", "urn"},
		},
		"OkCaseAction-" + api.USER_ACTION_LIST_GROUPS_FOR_USER: {
			action:          api.USER_ACTION_LIST_GROUPS_FOR_USER,
			expectedColumns: []string{"create_at"},
		},
		"OkCaseAction-" + api.GROUP_ACTION_LIST_GROUPS: {
			action:          api.GROUP_ACTION_LIST_GROUPS,
			expectedColumns: []string{"name", "path", "org", "create_at", "update_at", "urn"},
		},
		"OkCaseAction-" + api.GROUP_ACTION_LIST_MEMBERS: {
			action:          api.GROUP_ACTION_LIST_MEMBERS,
			expectedColumns: []string{"create_at"},
		},
		"OkCaseAction-" + api.GROUP_ACTION_LIST_ATTACHED_GROUP_POLICIES: {
			action:          api.GROUP_ACTION_LIST_ATTACHED_GROUP_POLICIES,
			expectedColumns: []string{"create_at"},
		},
		"OkCaseAction-" + api.POLICY_ACTION_LIST_POLICIES: {
			action:          api.POLICY_ACTION_LIST_POLICIES,
			expectedColumns: []string{"name", "path", "org", "create_at", "update_at", "urn"},
		},
		"OkCaseAction-" + api.POLICY_ACTION_LIST_ATTACHED_GROUPS: {
			action:          api.POLICY_ACTION_LIST_ATTACHED_GROUPS,
			expectedColumns: []string{"create_at"},
		},
		"OkCaseAction-" + api.PROXY_ACTION_LIST_RESOURCES: {
			action: api.PROXY_ACTION_LIST_RESOURCES,
			expectedColumns: []string{"name", "path", "org", "host", "path_resource", "method",
				"urn_resource", "urn", "action", "create_at", "update_at"},
		},
		"OkCaseOtherActions": {
			action:          "other",
			expectedColumns: nil,
		},
	}

	for n, test := range testcases {
		validColumns := PostgresRepo{}.OrderByValidColumns(test.action)
		assert.Equal(t, test.expectedColumns, validColumns, "Error in test case %v", n)
	}
}

// Aux methods

func insertUser(t *testing.T, testcase string, user User) {
	err := repoDB.Dbmap.Exec("INSERT INTO public.users (id, external_id, path, create_at, update_at, urn) VALUES (?, ?, ?, ?, ?, ?)",
		user.ID, user.ExternalID, user.Path, user.CreateAt, user.UpdateAt, user.Urn).Error

	// Error handling
	assert.Nil(t, err, "Error in test case %v", testcase)
}

func insertGroupUserRelation(t *testing.T, testcase string, userID string, groupID string, createAt int64) {
	err := repoDB.Dbmap.Exec("INSERT INTO public.group_user_relations (user_id, group_id, create_at) VALUES (?, ?, ?)",
		userID, groupID, createAt).Error

	// Error handling
	assert.Nil(t, err, "Error in test case %v", testcase)
}

func getUsersCountFiltered(t *testing.T, testcase string,
	id string, externalID string, path string, createAt int64, updateAt int64, urn string, pathPrefix string) int {
	query := repoDB.Dbmap.Table(User{}.TableName())
	if id != "" {
		query = query.Where("id = ?", id)
	}
	if externalID != "" {
		query = query.Where("external_id = ?", externalID)
	}
	if path != "" {
		query = query.Where("path = ?", path)
	}
	if pathPrefix != "" {
		query = query.Where("path like ?", pathPrefix+"%")
	}
	if createAt != 0 {
		query = query.Where("create_at = ?", createAt)
	}
	if updateAt != 0 {
		query = query.Where("update_at = ?", updateAt)
	}
	if urn != "" {
		query = query.Where("urn = ?", urn)
	}
	var number int
	err := query.Count(&number).Error
	assert.Nil(t, err, "Error in test case %v", testcase)

	return number
}

func cleanUserTable(t *testing.T, testcase string) {
	err := repoDB.Dbmap.Delete(&User{}).Error
	assert.Nil(t, err, "Error in test case %v", testcase)
}

func cleanGroupUserRelationTable(t *testing.T, testcase string) {
	err := repoDB.Dbmap.Delete(&GroupUserRelation{}).Error
	assert.Nil(t, err, "Error in test case %v", testcase)
}

// GROUP

func insertGroup(t *testing.T, testcase string, group Group) {
	err := repoDB.Dbmap.Exec("INSERT INTO public.groups (id, name, path, create_at, update_at, urn, org) VALUES (?, ?, ?, ?, ?, ?, ?)",
		group.ID, group.Name, group.Path, group.CreateAt, group.UpdateAt, group.Urn, group.Org).Error

	assert.Nil(t, err, "Error in test case %v", testcase)
}

func getGroupsCountFiltered(t *testing.T, testcase string,
	id string, name string, path string, createAt int64, updateAt int64, urn string, org string) int {
	query := repoDB.Dbmap.Table(Group{}.TableName())
	if id != "" {
		query = query.Where("id = ?", id)
	}
	if name != "" {
		query = query.Where("name = ?", name)
	}
	if path != "" {
		query = query.Where("path = ?", path)
	}
	if createAt != 0 {
		query = query.Where("create_at = ?", createAt)
	}
	if updateAt != 0 {
		query = query.Where("update_at = ?", updateAt)
	}
	if urn != "" {
		query = query.Where("urn = ?", urn)
	}
	if org != "" {
		query = query.Where("org = ?", org)
	}
	var number int
	err := query.Count(&number).Error
	assert.Nil(t, err, "Error in test case %v", testcase)

	return number
}

func getGroupUserRelations(t *testing.T, testcase string, groupID string, userID string) int {
	query := repoDB.Dbmap.Table(GroupUserRelation{}.TableName())
	if groupID != "" {
		query = query.Where("group_id = ?", groupID)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	var number int
	err := query.Count(&number).Error
	assert.Nil(t, err, "Error in test case %v", testcase)

	return number
}

func cleanGroupTable(t *testing.T, testcase string) {
	err := repoDB.Dbmap.Delete(&Group{}).Error
	assert.Nil(t, err, "Error in test case %v", testcase)
}

func cleanGroupPolicyRelationTable(t *testing.T, testcase string) {
	err := repoDB.Dbmap.Delete(&GroupPolicyRelation{}).Error
	assert.Nil(t, err, "Error in test case %v", testcase)
}

// POLICY

func cleanPolicyTable(t *testing.T, testcase string) {
	err := repoDB.Dbmap.Delete(&Policy{}).Error
	assert.Nil(t, err, "Error in test case %v", testcase)
}

func cleanStatementTable(t *testing.T, testcase string) {
	err := repoDB.Dbmap.Delete(&Statement{}).Error
	assert.Nil(t, err, "Error in test case %v", testcase)
}

func insertPolicy(t *testing.T, testcase string, policy Policy, statements []Statement) {
	err := repoDB.Dbmap.Exec("INSERT INTO public.policies (id, name, org, path, create_at, update_at, urn) VALUES (?, ?, ?, ?, ?, ?, ?)",
		policy.ID, policy.Name, policy.Org, policy.Path, policy.CreateAt, policy.UpdateAt, policy.Urn).Error

	// Error handling
	assert.Nil(t, err, "Error in test case %v", testcase)

	for _, v := range statements {
		v.PolicyID = policy.ID
		insertStatements(t, testcase, v)
		// Error handling
		assert.Nil(t, err, "Error in test case %v", testcase)
	}
}

func insertStatements(t *testing.T, testcase string, statement Statement) {
	err := repoDB.Dbmap.Exec("INSERT INTO public.statements (id, policy_id, effect, actions, resources) VALUES (?, ?, ?, ?, ?)",
		statement.ID, statement.PolicyID, statement.Effect, statement.Actions, statement.Resources).Error

	// Error handling
	assert.Nil(t, err, "Error in test case %v", testcase)
}

func getPoliciesCountFiltered(t *testing.T, testcase string,
	id string, org string, name string, path string, createAt int64, urn string) int {
	query := repoDB.Dbmap.Table(Policy{}.TableName())
	if id != "" {
		query = query.Where("id = ?", id)
	}
	if org != "" {
		query = query.Where("org = ?", org)
	}
	if path != "" {
		query = query.Where("path = ?", path)
	}
	if name != "" {
		query = query.Where("name = ?", name)
	}
	if createAt != 0 {
		query = query.Where("create_at = ?", createAt)
	}
	if urn != "" {
		query = query.Where("urn = ?", urn)
	}
	var number int
	err := query.Count(&number).Error
	assert.Nil(t, err, "Error in test case %v", testcase)

	return number
}

func getGroupPolicyRelationCount(t *testing.T, testcase string, policyID string, groupID string) int {
	query := repoDB.Dbmap.Table(GroupPolicyRelation{}.TableName())
	if policyID != "" {
		query = query.Where("policy_id = ?", policyID)
	}
	if groupID != "" {
		query = query.Where("group_id = ?", groupID)
	}
	var number int
	err := query.Count(&number).Error
	assert.Nil(t, err, "Error in test case %v", testcase)

	return number
}

func insertGroupPolicyRelation(t *testing.T, testcase string, groupID string, policyID string, createAt int64) {
	err := repoDB.Dbmap.Exec("INSERT INTO public.group_policy_relations (group_id, policy_id, create_at) VALUES (?, ?, ?)",
		groupID, policyID, createAt).Error

	// Error handling
	assert.Nil(t, err, "Error in test case %v", testcase)
}

func getStatementsCountFiltered(t *testing.T, testcase string,
	id string, policyId string, effect string, actions string, resources string) int {
	query := repoDB.Dbmap.Table(Statement{}.TableName())
	if id != "" {
		query = query.Where("id = ?", id)
	}
	if policyId != "" {
		query = query.Where("policy_id = ?", policyId)
	}
	if effect != "" {
		query = query.Where("effect = ?", effect)
	}
	if actions != "" {
		query = query.Where("actions = ?", actions)
	}
	if resources != "" {
		query = query.Where("resources = ?", resources)
	}
	var number int
	err := query.Count(&number).Error
	assert.Nil(t, err, "Error in test case %v", testcase)

	return number
}

// PROXY

func cleanProxyResourcesTable(t *testing.T, testcase string) {
	err := repoDB.Dbmap.Delete(&ProxyResource{}).Error
	assert.Nil(t, err, "Error in test case %v", testcase)
}

func insertProxyResource(t *testing.T, testcase string, pr ProxyResource) {
	err := repoDB.Dbmap.Exec("INSERT INTO public.proxy_resources (id, name, org, path, host, path_resource, method, urn_resource, "+
		"urn, action, create_at, update_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		pr.ID, pr.Name, pr.Org, pr.Path, pr.Host, pr.PathResource, pr.Method, pr.UrnResource, pr.Urn, pr.Action, pr.CreateAt, pr.UpdateAt).Error

	// Error handling
	assert.Nil(t, err, "Error in testcase %v", testcase)
}

func getProxyResourcesCountFiltered(t *testing.T, testcase string, id string,
	name string, org string, path string, urn string, createAt int64, updateAt int64) int {
	query := repoDB.Dbmap.Table(ProxyResource{}.TableName())
	if id != "" {
		query = query.Where("id = ?", id)
	}
	if org != "" {
		query = query.Where("org = ?", org)
	}
	if path != "" {
		query = query.Where("path = ?", path)
	}
	if name != "" {
		query = query.Where("name = ?", name)
	}
	if urn != "" {
		query = query.Where("urn = ?", urn)
	}
	if createAt != 0 {
		query = query.Where("create_at = ?", createAt)
	}
	if updateAt != 0 {
		query = query.Where("update_at = ?", updateAt)
	}
	var number int
	err := query.Count(&number).Error
	assert.Nil(t, err, "Error in test case %v", testcase)

	return number
}
