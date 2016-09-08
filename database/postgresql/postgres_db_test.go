package postgresql

import (
	"fmt"
	"os"
	"testing"
	"time"

	"errors"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/database"
	"github.com/kylelemons/godebug/pretty"
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
		if diff := pretty.Compare(err, test.expectedError); diff != "" {
			t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
			continue
		}
	}

}

// Aux methods

func insertUser(id string, externalID string, path string, createAt int64, urn string) error {
	err := repoDB.Dbmap.Exec("INSERT INTO public.users (id, external_id, path, create_at, urn) VALUES (?, ?, ?, ?, ?)",
		id, externalID, path, createAt, urn).Error

	// Error handling
	if err != nil {
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}
	return nil
}

func insertGroupUserRelation(userID string, groupID string) error {
	err := repoDB.Dbmap.Exec("INSERT INTO public.group_user_relations (user_id, group_id) VALUES (?, ?)",
		userID, groupID).Error

	// Error handling
	if err != nil {
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}
	return nil
}

func getUsersCountFiltered(id string, externalID string, path string, createAt int64, urn string, pathPrefix string) (int, error) {
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
	if urn != "" {
		query = query.Where("urn = ?", urn)
	}
	var number int
	if err := query.Count(&number).Error; err != nil {
		return 0, err
	}

	return number, nil
}

func cleanUserTable() error {
	if err := repoDB.Dbmap.Delete(&User{}).Error; err != nil {
		return err
	}
	return nil
}

func cleanGroupUserRelationTable() error {
	if err := repoDB.Dbmap.Delete(&GroupUserRelation{}).Error; err != nil {
		return err
	}
	return nil
}

// GROUP

func insertGroup(id string, name string, path string, createAt int64, urn string, org string) error {
	err := repoDB.Dbmap.Exec("INSERT INTO public.groups (id, name, path, create_at, urn, org) VALUES (?, ?, ?, ?, ?, ?)",
		id, name, path, createAt, urn, org).Error

	// Error handling
	if err != nil {
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}
	return nil
}

func getGroupsCountFiltered(id string, name string, path string, createAt int64, urn string, org string) (int, error) {
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
	if urn != "" {
		query = query.Where("urn = ?", urn)
	}
	if org != "" {
		query = query.Where("org = ?", org)
	}
	var number int
	if err := query.Count(&number).Error; err != nil {
		return 0, err
	}

	return number, nil
}

func getGroupUserRelations(groupID string, userID string) (int, error) {
	query := repoDB.Dbmap.Table(GroupUserRelation{}.TableName())
	if groupID != "" {
		query = query.Where("group_id = ?", groupID)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	var number int
	if err := query.Count(&number).Error; err != nil {
		return 0, err
	}

	return number, nil
}

func cleanGroupTable() error {
	if err := repoDB.Dbmap.Delete(&Group{}).Error; err != nil {
		return err
	}
	return nil
}

func cleanGroupPolicyRelationTable() error {
	if err := repoDB.Dbmap.Delete(&GroupPolicyRelation{}).Error; err != nil {
		return err
	}
	return nil
}

// POLICY

func cleanPolicyTable() error {
	if err := repoDB.Dbmap.Delete(&Policy{}).Error; err != nil {
		return err
	}
	return nil
}

func cleanStatementTable() error {
	if err := repoDB.Dbmap.Delete(&Statement{}).Error; err != nil {
		return err
	}
	return nil
}

func insertPolicy(id string, name string, org string, path string, createAt int64, urn string, statements []Statement) error {
	err := repoDB.Dbmap.Exec("INSERT INTO public.policies (id, name, org, path, create_at, urn) VALUES (?, ?, ?, ?, ?, ?)",
		id, name, org, path, createAt, urn).Error

	// Error handling
	if err != nil {
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	for _, v := range statements {
		err = insertStatements(v.ID, v.PolicyID, v.Actions, v.Effect, v.Resources)
		// Error handling
		if err != nil {
			return &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: err.Error(),
			}
		}
	}

	return nil
}

func insertStatements(id string, policyId string, actions string, effect string, resources string) error {
	err := repoDB.Dbmap.Exec("INSERT INTO public.statements (id, policy_id, effect, actions, resources) VALUES (?, ?, ?, ?, ?)",
		id, policyId, effect, actions, resources).Error

	// Error handling
	if err != nil {
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}

	return nil
}

func getPoliciesCountFiltered(id string, org string, name string, path string, createAt int64, urn string) (int, error) {
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
	if err := query.Count(&number).Error; err != nil {
		return 0, err
	}

	return number, nil
}

func getGroupPolicyRelationCount(policyID string, groupID string) (int, error) {
	query := repoDB.Dbmap.Table(GroupPolicyRelation{}.TableName())
	if policyID != "" {
		query = query.Where("policy_id = ?", policyID)
	}
	if groupID != "" {
		query = query.Where("group_id = ?", groupID)
	}
	var number int
	if err := query.Count(&number).Error; err != nil {
		return 0, err
	}

	return number, nil
}

func insertGroupPolicyRelation(groupID string, policyID string) error {
	err := repoDB.Dbmap.Exec("INSERT INTO public.group_policy_relations (group_id, policy_id) VALUES (?, ?)",
		groupID, policyID).Error

	// Error handling
	if err != nil {
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}
	return nil
}

func getStatementsCountFiltered(id string, policyId string, effect string, actions string, resources string) (int, error) {
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
	if err := query.Count(&number).Error; err != nil {
		return 0, err
	}

	return number, nil
}
