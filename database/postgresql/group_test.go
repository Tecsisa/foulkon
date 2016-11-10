package postgresql

import (
	"testing"
	"time"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/database"

	"github.com/stretchr/testify/assert"
)

func TestPostgresRepo_AddGroup(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousGroup *Group
		// Postgres Repo Args
		groupToCreate *api.Group
		// Expected result
		expectedResponse *api.Group
		expectedError    *database.Error
	}{
		"OkCase": {
			groupToCreate: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "urn",
				CreateAt: now,
				UpdateAt: now,
				Org:      "Org",
			},
			expectedResponse: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "urn",
				CreateAt: now,
				UpdateAt: now,
				Org:      "Org",
			},
		},
		"ErrorCasegroupAlreadyExist": {
			previousGroup: &Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "urn",
				CreateAt: now.UnixNano(),
				UpdateAt: now.UnixNano(),
				Org:      "Org",
			},
			groupToCreate: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "urn",
				CreateAt: now,
				UpdateAt: now,
				Org:      "Org",
			},
			expectedError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "pq: duplicate key value violates unique constraint \"groups_pkey\"",
			},
		},
	}

	for n, test := range testcases {
		// Clean user database
		cleanGroupTable(t, n)

		// Insert previous data
		if test.previousGroup != nil {
			insertGroup(t, n, *test.previousGroup)
		}
		// Call to repository to store group
		storedGroup, err := repoDB.AddGroup(*test.groupToCreate)
		if test.expectedError != nil {
			dbError, _ := err.(*database.Error)
			assert.Equal(t, dbError, test.expectedError, "Error in test case %v", n)
		} else {
			assert.Nil(t, err, "Error in test case %v", n)
			// Check response
			assert.Equal(t, storedGroup, test.expectedResponse, "Error in test case %v", n)
			// Check database
			groupNumber := getGroupsCountFiltered(t, n, test.groupToCreate.ID, test.groupToCreate.Name, test.groupToCreate.Path,
				test.groupToCreate.CreateAt.UnixNano(), test.groupToCreate.UpdateAt.UnixNano(), test.groupToCreate.Urn, test.groupToCreate.Org)

			if groupNumber != 1 {
				t.Errorf("Test %v failed. Received different group number: %v", n, groupNumber)
				continue
			}
		}
	}
}

func TestPostgresRepo_GetGroupByName(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousGroup *Group
		// Postgres Repo Args
		org  string
		name string
		// Expected result
		expectedResponse *api.Group
		expectedError    *database.Error
	}{
		"OkCase": {
			previousGroup: &Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Urn",
				CreateAt: now.UnixNano(),
				UpdateAt: now.UnixNano(),
				Org:      "Org",
			},
			org:  "Org",
			name: "Name",
			expectedResponse: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Urn",
				CreateAt: now,
				UpdateAt: now,
				Org:      "Org",
			},
		},
		"ErrorCaseGroupNotExist": {
			previousGroup: &Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Urn",
				CreateAt: now.UnixNano(),
				UpdateAt: now.UnixNano(),
				Org:      "Org",
			},
			org:  "Org",
			name: "NotExist",
			expectedError: &database.Error{
				Code:    database.GROUP_NOT_FOUND,
				Message: "Group with organization Org and name NotExist not found",
			},
		},
	}

	for n, test := range testcases {
		// Clean group database
		cleanGroupTable(t, n)

		// Insert previous data
		if test.previousGroup != nil {
			insertGroup(t, n, *test.previousGroup)
		}

		// Call to repository to get group
		receivedGroup, err := repoDB.GetGroupByName(test.org, test.name)
		if test.expectedError != nil {
			dbError, _ := err.(*database.Error)
			assert.Equal(t, test.expectedError, dbError, "Error in test case %v", n)
		} else {
			assert.Nil(t, err, "Error in test case %v", n)
			// Check response
			assert.Equal(t, receivedGroup, test.expectedResponse, "Error in test case %v", n)
		}
	}
}

func TestPostgresRepo_GetGroupById(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousGroup *Group
		// Postgres Repo Args
		groupID string
		// Expected result
		expectedResponse *api.Group
		expectedError    *database.Error
	}{
		"OkCase": {
			previousGroup: &Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Urn",
				CreateAt: now.UnixNano(),
				UpdateAt: now.UnixNano(),
				Org:      "Org",
			},
			groupID: "GroupID",
			expectedResponse: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Urn",
				CreateAt: now,
				UpdateAt: now,
				Org:      "Org",
			},
		},
		"ErrorCaseGroupNotExist": {
			previousGroup: &Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Urn",
				CreateAt: now.UnixNano(),
				UpdateAt: now.UnixNano(),
				Org:      "Org",
			},
			groupID: "NotExist",
			expectedError: &database.Error{
				Code:    database.GROUP_NOT_FOUND,
				Message: "Group with id NotExist not found",
			},
		},
	}

	for n, test := range testcases {
		// Clean group database
		cleanGroupTable(t, n)

		// Insert previous data
		if test.previousGroup != nil {
			insertGroup(t, n, *test.previousGroup)
		}

		// Call to repository to get group
		receivedGroup, err := repoDB.GetGroupById(test.groupID)
		if test.expectedError != nil {
			dbError, _ := err.(*database.Error)
			assert.Equal(t, test.expectedError, dbError, "Error in test case %v", n)
		} else {
			assert.Nil(t, err, "Error in test case %v", n)
			// Check response
			assert.Equal(t, receivedGroup, test.expectedResponse, "Error in test case %v", n)
		}
	}
}

func TestPostgresRepo_GetGroupsFiltered(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousGroups []Group
		// Postgres Repo Args
		org    string
		filter *api.Filter
		// Expected result
		expectedResponse []api.Group
	}{
		"OkCasePathPrefix1": {
			previousGroups: []Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Org:      "Org1",
				},
				{
					ID:       "GroupID2",
					Name:     "Name2",
					Path:     "Path456",
					Urn:      "urn2",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Org:      "Org2",
				},
			},
			filter: &api.Filter{
				PathPrefix: "Path",
				Offset:     0,
				Limit:      20,
			},
			expectedResponse: []api.Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now,
					UpdateAt: now,
					Org:      "Org1",
				},
				{
					ID:       "GroupID2",
					Name:     "Name2",
					Path:     "Path456",
					Urn:      "urn2",
					CreateAt: now,
					UpdateAt: now,
					Org:      "Org2",
				},
			},
		},
		"OkCasePathPrefix2": {
			previousGroups: []Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Org:      "Org1",
				},
				{
					ID:       "GroupID2",
					Name:     "Name2",
					Path:     "Path456",
					Urn:      "urn2",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Org:      "Org2",
				},
			},
			filter: &api.Filter{
				PathPrefix: "Path123",
				Offset:     0,
				Limit:      20,
			},
			expectedResponse: []api.Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now,
					UpdateAt: now,
					Org:      "Org1",
				},
			},
		},
		"OkCasePathPrefix3": {
			previousGroups: []Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Org:      "Org1",
				},
				{
					ID:       "GroupID2",
					Name:     "Name2",
					Path:     "Path456",
					Urn:      "urn2",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Org:      "Org2",
				},
			},
			filter: &api.Filter{
				PathPrefix: "NoPath",
				Offset:     0,
				Limit:      20,
			},
			expectedResponse: []api.Group{},
		},
		"OkCaseGetByOrg": {
			previousGroups: []Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Org:      "Org1",
				},
				{
					ID:       "GroupID2",
					Name:     "Name2",
					Path:     "Path456",
					Urn:      "urn2",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Org:      "Org2",
				},
			},
			filter: &api.Filter{
				Org: "Org1",
			},
			expectedResponse: []api.Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now,
					UpdateAt: now,
					Org:      "Org1",
				},
			},
		},
		"OkCaseOrderBy": {
			previousGroups: []Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Org:      "Org1",
				},
				{
					ID:       "GroupID2",
					Name:     "Name2",
					Path:     "Path456",
					Urn:      "urn2",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Org:      "Org2",
				},
			},
			filter: &api.Filter{
				OrderBy: "path desc",
			},
			expectedResponse: []api.Group{
				{
					ID:       "GroupID2",
					Name:     "Name2",
					Path:     "Path456",
					Urn:      "urn2",
					CreateAt: now,
					UpdateAt: now,
					Org:      "Org2",
				},
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now,
					UpdateAt: now,
					Org:      "Org1",
				},
			},
		},
		"OkCaseGetByOrgAndPathPrefix": {
			previousGroups: []Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Org:      "Org1",
				},
				{
					ID:       "GroupID2",
					Name:     "Name2",
					Path:     "Path456",
					Urn:      "urn2",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Org:      "Org2",
				},
			},
			filter: &api.Filter{
				PathPrefix: "Path123",
				Org:        "Org1",
				Offset:     0,
				Limit:      20,
			},
			expectedResponse: []api.Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now,
					UpdateAt: now,
					Org:      "Org1",
				},
			},
		},
		"OkCaseWithoutParams": {
			previousGroups: []Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Org:      "Org1",
				},
				{
					ID:       "GroupID2",
					Name:     "Name2",
					Path:     "Path456",
					Urn:      "urn2",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Org:      "Org2",
				},
			},
			filter: testFilter,
			expectedResponse: []api.Group{
				{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path123",
					Urn:      "urn1",
					CreateAt: now,
					UpdateAt: now,
					Org:      "Org1",
				},
				{
					ID:       "GroupID2",
					Name:     "Name2",
					Path:     "Path456",
					Urn:      "urn2",
					CreateAt: now,
					UpdateAt: now,
					Org:      "Org2",
				},
			},
		},
	}

	for n, test := range testcases {
		// Clean group database
		cleanGroupTable(t, n)

		// Insert previous data
		if test.previousGroups != nil {
			for _, previousGroup := range test.previousGroups {
				insertGroup(t, n, previousGroup)
			}
		}
		// Call to repository to get groups
		receivedGroups, total, err := repoDB.GetGroupsFiltered(test.filter)
		assert.Nil(t, err, "Error in test case %v", n)

		// Check total
		assert.Equal(t, total, len(test.expectedResponse), "Error in test case %v", n)

		// Check response
		assert.Equal(t, receivedGroups, test.expectedResponse, "Error in test case %v", n)
	}
}

func TestPostgresRepo_UpdateGroup(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousGroups []Group
		// Postgres Repo Args
		groupToUpdate *api.Group
		// Expected result
		expectedResponse *api.Group
		expectedError    *database.Error
	}{
		"OkCase": {
			previousGroups: []Group{
				{
					ID:       "GroupID",
					Name:     "Name",
					Path:     "Path",
					Urn:      "Urn",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Org:      "Org",
				},
			},
			groupToUpdate: &api.Group{
				ID:       "GroupID",
				Name:     "NewName",
				Path:     "NewPath",
				Urn:      "NewUrn",
				CreateAt: now,
				UpdateAt: now,
				Org:      "Org",
			},
			expectedResponse: &api.Group{
				ID:       "GroupID",
				Name:     "NewName",
				Path:     "NewPath",
				Urn:      "NewUrn",
				CreateAt: now,
				UpdateAt: now,
				Org:      "Org",
			},
		},
		"ErrorCaseDuplicateUrn": {
			previousGroups: []Group{
				{
					ID:       "GroupID",
					Name:     "Name",
					Path:     "Path",
					Urn:      "Urn",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Org:      "Org",
				},
				{
					ID:       "GroupID2",
					Name:     "Name2",
					Path:     "Path2",
					Urn:      "Fail",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Org:      "Org2",
				},
			},
			groupToUpdate: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Fail",
				CreateAt: now,
				Org:      "Org",
			},
			expectedError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "pq: duplicate key value violates unique constraint \"groups_urn_key\"",
			},
		},
	}

	for n, test := range testcases {
		// Clean group database
		cleanGroupTable(t, n)

		// Insert previous data
		if test.previousGroups != nil {
			for _, previousGroup := range test.previousGroups {
				insertGroup(t, n, previousGroup)
			}
		}

		// Call to repository to update group
		updatedGroup, err := repoDB.UpdateGroup(*test.groupToUpdate)
		if test.expectedError != nil {
			dbError, _ := err.(*database.Error)
			assert.Equal(t, test.expectedError, dbError, "Error in test case %v", n)
		} else {
			assert.Nil(t, err, "Error in test case %v", n)
			// Check response
			assert.Equal(t, updatedGroup, test.expectedResponse, "Error in test case %v", n)
			// Check database
			groupNumber := getGroupsCountFiltered(t, n, test.expectedResponse.ID, test.expectedResponse.Name, test.expectedResponse.Path,
				test.expectedResponse.CreateAt.UnixNano(), test.expectedResponse.UpdateAt.UnixNano(), test.expectedResponse.Urn,
				test.expectedResponse.Org)
			assert.Equal(t, 1, groupNumber, "Error in test case %v", n)
		}
	}
}

func TestPostgresRepo_RemoveGroup(t *testing.T) {
	type userRelation struct {
		userID   string
		groupID  string
		CreateAt int64
	}
	type policyRelation struct {
		policyID string
		groupID  string
		CreateAt int64
	}
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousGroups  []Group
		userRelations   []userRelation
		policyRelations []policyRelation
		// Postgres Repo Args
		groupToDelete string
	}{
		"OkCase": {
			previousGroups: []Group{
				{
					ID:       "GroupID",
					Name:     "Name",
					Path:     "Path",
					Urn:      "Urn",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Org:      "Org",
				},
				{
					ID:       "GroupID2",
					Name:     "Name",
					Path:     "Path",
					Urn:      "Urn2",
					CreateAt: now.UnixNano(),
					UpdateAt: now.UnixNano(),
					Org:      "Org",
				},
			},
			userRelations: []userRelation{
				{
					userID:   "UserID",
					groupID:  "GroupID",
					CreateAt: now.UnixNano(),
				},
				{
					userID:   "UserID1",
					groupID:  "GroupID2",
					CreateAt: now.UnixNano(),
				},
			},
			policyRelations: []policyRelation{
				{
					policyID: "policyID",
					groupID:  "GroupID",
					CreateAt: now.UnixNano(),
				},
				{
					policyID: "policyID1",
					groupID:  "GroupID2",
					CreateAt: now.UnixNano(),
				},
			},
			groupToDelete: "GroupID",
		},
	}

	for n, test := range testcases {
		cleanGroupTable(t, n)
		cleanGroupUserRelationTable(t, n)
		cleanGroupPolicyRelationTable(t, n)

		// Insert previous data
		if test.previousGroups != nil {
			for _, g := range test.previousGroups {
				insertGroup(t, n, g)
			}

		}
		if test.userRelations != nil {
			for _, rel := range test.userRelations {
				insertGroupUserRelation(t, n, rel.userID, rel.groupID, rel.CreateAt)
			}
		}
		if test.policyRelations != nil {
			for _, rel := range test.policyRelations {
				insertGroupPolicyRelation(t, n, rel.groupID, rel.policyID, rel.CreateAt)
			}
		}
		// Call to repository to remove group
		err := repoDB.RemoveGroup(test.groupToDelete)
		assert.Nil(t, err, "Error in test case %v", n)

		// Check database
		groupNumber := getGroupsCountFiltered(t, n, test.groupToDelete, "", "", 0, 0, "", "")
		assert.Equal(t, 0, groupNumber, "Error in test case %v", n)

		// Check total groups
		totalGroupNumber := getGroupsCountFiltered(t, n, "", "", "", 0, 0, "", "")
		assert.Equal(t, 1, totalGroupNumber, "Error in test case %v", n)

		// Check group user relations
		relations := getGroupUserRelations(t, n, test.groupToDelete, "")
		assert.Equal(t, 0, relations, "Error in test case %v", n)

		// Check total group user relations
		totalRelations := getGroupUserRelations(t, n, "", "")
		assert.Equal(t, 1, totalRelations, "Error in test case %v", n)

		// Check group policy relations
		relations = getGroupPolicyRelationCount(t, n, "", test.groupToDelete)
		assert.Equal(t, 0, relations, "Error in test case %v", n)

		// Check total group policy relations
		totalRelations = getGroupPolicyRelationCount(t, n, "", "")
		assert.Equal(t, 1, totalRelations, "Error in test case %v", n)
	}
}

func TestPostgresRepo_AddMember(t *testing.T) {
	testcases := map[string]struct {
		// Postgres Repo Args
		userID  string
		groupID string
		// Expected result
		expectedError *database.Error
	}{
		"OkCase": {
			userID:  "UserID",
			groupID: "GroupID",
		},
		"ErrorCaseInternalError": {
			groupID: "GroupID",
			expectedError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "pq: null value in column \"user_id\" violates not-null constraint",
			},
		},
	}

	for n, test := range testcases {
		// Clean GroupUserRelation database
		cleanGroupUserRelationTable(t, n)

		// Call to repository to store member
		err := repoDB.AddMember(test.userID, test.groupID)
		if test.expectedError != nil {
			dbError, _ := err.(*database.Error)
			assert.Equal(t, test.expectedError, dbError, "Error in test case %v", n)
		} else {
			assert.Nil(t, err, "Error in test case %v", n)

			// Check database
			relations := getGroupUserRelations(t, n, test.groupID, test.userID)
			assert.Equal(t, 1, relations, "Error in test case %v", n)
		}
	}
}

func TestPostgresRepo_RemoveMember(t *testing.T) {
	type relation struct {
		userID   string
		groupID  string
		createAt int64
	}
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		relation *relation
		// Postgres Repo Args
		userID  string
		groupID string
		// Expected result
		expectedError *database.Error
	}{
		"OkCase": {
			relation: &relation{
				userID:   "UserID",
				groupID:  "GroupID",
				createAt: now.UnixNano(),
			},
			userID:  "UserID",
			groupID: "GroupID",
		},
	}

	for n, test := range testcases {
		// Clean GroupUserRelation database
		cleanGroupUserRelationTable(t, n)

		// Insert previous data
		if test.relation != nil {
			insertGroupUserRelation(t, n, test.relation.userID, test.relation.groupID, test.relation.createAt)
		}

		// Call to repository to remove member
		err := repoDB.RemoveMember(test.userID, test.groupID)
		assert.Nil(t, err, "Error in test case %v", n)

		// Check database
		relations := getGroupUserRelations(t, n, test.groupID, test.userID)
		assert.Equal(t, 0, relations, "Error in test case %v", n)
	}
}

func TestPostgresRepo_IsMemberOfGroup(t *testing.T) {
	type relation struct {
		userID   string
		groupID  string
		createAt int64
	}
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		relation *relation
		// Postgres Repo Args
		group  string
		member string
		// Expected result
		isMember bool
	}{
		"OkCaseIsMember": {
			relation: &relation{
				userID:   "UserID",
				groupID:  "GroupID",
				createAt: now.UnixNano(),
			},
			group:    "GroupID",
			member:   "UserID",
			isMember: true,
		},
		"OkCaseIsNotMember": {
			group:    "GroupID",
			member:   "UserID",
			isMember: false,
		},
	}

	for n, test := range testcases {
		cleanGroupUserRelationTable(t, n)

		// Insert previous data
		if test.relation != nil {
			insertGroupUserRelation(t, n, test.relation.userID, test.relation.groupID, test.relation.createAt)
		}

		isMember, err := repoDB.IsMemberOfGroup(test.member, test.group)
		assert.Nil(t, err, "Error in test case %v", n)

		// Check response
		assert.Equal(t, isMember, test.isMember, "Error in test case %v", n)
	}
}

func TestPostgresRepo_GetGroupMembers(t *testing.T) {
	type relations struct {
		users        []User
		groupID      string
		createAt     []int64
		userNotFound bool
	}
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		relations *relations
		// Postgres Repo Args
		groupID string
		filter  *api.Filter
		// Expected result
		expectedResponse []*GroupUser
		expectedError    *database.Error
	}{
		"OkCase": {
			relations: &relations{
				users: []User{
					{
						ID:         "UserID1",
						ExternalID: "ExternalID1",
						Path:       "Path",
						Urn:        "urn1",
						CreateAt:   now.UnixNano(),
						UpdateAt:   now.UnixNano(),
					},
					{
						ID:         "UserID2",
						ExternalID: "ExternalID2",
						Path:       "Path",
						Urn:        "urn2",
						CreateAt:   now.UnixNano(),
						UpdateAt:   now.UnixNano(),
					},
				},
				groupID:  "GroupID",
				createAt: []int64{now.UnixNano(), now.UnixNano()},
			},
			groupID: "GroupID",
			filter:  testFilter,
			expectedResponse: []*GroupUser{
				{
					User: &api.User{
						ID:         "UserID1",
						ExternalID: "ExternalID1",
						Path:       "Path",
						Urn:        "urn1",
						CreateAt:   now,
						UpdateAt:   now,
					},
					CreateAt: now,
				},
				{
					User: &api.User{
						ID:         "UserID2",
						ExternalID: "ExternalID2",
						Path:       "Path",
						Urn:        "urn2",
						CreateAt:   now,
						UpdateAt:   now,
					},
					CreateAt: now,
				},
			},
		},
		"OkCaseOrderBy": {
			relations: &relations{
				users: []User{
					{
						ID:         "UserID1",
						ExternalID: "ExternalID1",
						Path:       "Path",
						Urn:        "urn1",
						CreateAt:   now.UnixNano(),
						UpdateAt:   now.UnixNano(),
					},
					{
						ID:         "UserID2",
						ExternalID: "ExternalID2",
						Path:       "Path",
						Urn:        "urn2",
						CreateAt:   now.UnixNano(),
						UpdateAt:   now.UnixNano(),
					},
				},
				groupID:  "GroupID",
				createAt: []int64{now.UnixNano(), now.Add(-1).UnixNano()},
			},
			groupID: "GroupID",
			filter: &api.Filter{
				OrderBy: "create_at desc",
			},
			expectedResponse: []*GroupUser{
				{
					User: &api.User{
						ID:         "UserID1",
						ExternalID: "ExternalID1",
						Path:       "Path",
						Urn:        "urn1",
						CreateAt:   now,
						UpdateAt:   now,
					},
					CreateAt: now,
				},
				{
					User: &api.User{
						ID:         "UserID2",
						ExternalID: "ExternalID2",
						Path:       "Path",
						Urn:        "urn2",
						CreateAt:   now,
						UpdateAt:   now,
					},
					CreateAt: now.Add(-1),
				},
			},
		},
		"ErrorCase": {
			relations: &relations{
				users: []User{
					{
						ID: "UserID1",
					},
				},
				groupID:      "GroupID",
				userNotFound: true,
				createAt:     []int64{now.UnixNano()},
			},
			groupID: "GroupID",
			filter:  testFilter,
			expectedError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "Code: UserNotFound, Message: User with id UserID1 not found",
			},
		},
	}

	for n, test := range testcases {
		cleanUserTable(t, n)
		cleanGroupUserRelationTable(t, n)

		// Insert previous data
		if test.relations != nil {
			for x, user := range test.relations.users {
				insertGroupUserRelation(t, n, user.ID, test.relations.groupID, test.relations.createAt[x])
				if !test.relations.userNotFound {
					insertUser(t, n, user)
				}
			}

		}

		receivedUsers, total, err := repoDB.GetGroupMembers(test.groupID, test.filter)
		if test.expectedError != nil {
			dbError, _ := err.(*database.Error)
			assert.Equal(t, test.expectedError, dbError, "Error in test case %v", n)
		} else {
			assert.Nil(t, err, "Error in test case %v", n)
			// Check total
			assert.Equal(t, len(test.expectedResponse), total, "Error in test case %v", n)

			// Check response
			for i, r := range receivedUsers {
				assert.Equal(t, test.expectedResponse[i].GetGroup(), r.GetGroup(), "Error in test case %v", n)
				assert.Equal(t, test.expectedResponse[i].GetUser(), r.GetUser(), "Error in test case %v", n)
				assert.Equal(t, test.expectedResponse[i].GetDate(), r.GetDate(), "Error in test case %v", n)
			}
		}
	}
}

func TestPostgresRepo_AttachPolicy(t *testing.T) {
	testcases := map[string]struct {
		// Postgres Repo Args
		policyID string
		groupID  string
		// Expected result
		expectedError *database.Error
	}{
		"OkCase": {
			policyID: "PolicyID",
			groupID:  "GroupID",
		},
		"ErrorCaseInternalError": {
			expectedError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "pq: null value in column \"group_id\" violates not-null constraint",
			},
		},
	}

	for n, test := range testcases {
		// Clean GroupPolicyRelation database
		cleanGroupPolicyRelationTable(t, n)

		// Call to repository to attach policy
		err := repoDB.AttachPolicy(test.groupID, test.policyID)
		if test.expectedError != nil {
			dbError, _ := err.(*database.Error)
			assert.Equal(t, test.expectedError, dbError, "Error in test case %v", n)
		} else {
			assert.Nil(t, err, "Error in test case %v", n)
			// Check database
			relations := getGroupPolicyRelationCount(t, n, test.policyID, test.groupID)
			assert.Equal(t, 1, relations, "Error in test case %v", n)
		}
	}
}

func TestPostgresRepo_DetachPolicy(t *testing.T) {
	type relation struct {
		policyID string
		groupID  string
		createAt int64
	}
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		relation *relation
		// Postgres Repo Args
		policyID string
		groupID  string
		// Expected result
		expectedError *database.Error
	}{
		"OkCase": {
			relation: &relation{
				policyID: "PolicyID",
				groupID:  "GroupID",
				createAt: now.UnixNano(),
			},
			policyID: "PolicyID",
			groupID:  "GroupID",
		},
	}

	for n, test := range testcases {
		// Clean GroupPolicyRelation database
		cleanGroupPolicyRelationTable(t, n)

		// Insert previous data
		if test.relation != nil {
			insertGroupPolicyRelation(t, n, test.relation.groupID, test.relation.policyID, test.relation.createAt)
		}

		// Call to repository to detach policy
		err := repoDB.DetachPolicy(test.groupID, test.policyID)
		assert.Nil(t, err, "Error in test case %v", n)

		// Check database
		relations := getGroupPolicyRelationCount(t, n, test.policyID, test.groupID)
		assert.Equal(t, 0, relations, "Error in test case %v", n)
	}
}

func TestPostgresRepo_IsAttachedToGroup(t *testing.T) {
	type relation struct {
		groupID  string
		policyID string
		createAt int64
	}
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		relation *relation
		// Postgres Repo Args
		groupID  string
		policyID string
		// Expected result
		expectedResult bool
	}{
		"OkCase": {
			relation: &relation{
				groupID:  "GroupID",
				policyID: "PolicyID",
				createAt: now.UnixNano(),
			},
			groupID:        "GroupID",
			policyID:       "PolicyID",
			expectedResult: true,
		},
		"OkCaseNotFound": {
			relation: &relation{
				groupID:  "GroupID",
				policyID: "PolicyID",
				createAt: now.UnixNano(),
			},
			groupID:        "GroupID",
			policyID:       "PolicyIDXXXXXXX",
			expectedResult: false,
		},
	}

	for n, test := range testcases {
		// Clean GroupPolicyRelation database
		cleanGroupPolicyRelationTable(t, n)

		// Insert previous data
		if test.relation != nil {
			insertGroupPolicyRelation(t, n, test.relation.groupID, test.relation.policyID, test.relation.createAt)
		}

		// Call repository to check if policy is attached to group
		result, err := repoDB.IsAttachedToGroup(test.groupID, test.policyID)
		assert.Nil(t, err, "Error in test case %v", n)
		assert.Equal(t, test.expectedResult, result, "Error in test case %v", n)
	}
}

func TestPostgresRepo_GetAttachedPolicies(t *testing.T) {
	type relations struct {
		policies       []Policy
		groupID        string
		createAt       []int64
		policyNotFound bool
	}
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		relations  *relations
		statements []Statement
		// Postgres Repo Args
		groupID string
		filter  *api.Filter
		// Expected result
		expectedResponse []*PolicyGroup
		expectedError    *database.Error
	}{
		"OkCase": {
			relations: &relations{
				policies: []Policy{
					{
						ID:       "PolicyID1",
						Name:     "Name1",
						Org:      "org1",
						Path:     "/path/",
						CreateAt: now.UnixNano(),
						UpdateAt: now.UnixNano(),
						Urn:      "Urn1",
					},
					{
						ID:       "PolicyID2",
						Name:     "Name2",
						Org:      "org1",
						Path:     "/path/",
						CreateAt: now.UnixNano(),
						UpdateAt: now.UnixNano(),
						Urn:      "Urn2",
					},
				},
				groupID:  "GroupID",
				createAt: []int64{now.UnixNano() - 1, now.UnixNano()},
			},
			statements: []Statement{},
			groupID:    "GroupID",
			filter: &api.Filter{
				OrderBy: "create_at desc",
			},
			expectedResponse: []*PolicyGroup{
				{
					Policy: &api.Policy{
						ID:         "PolicyID2",
						Name:       "Name2",
						Org:        "org1",
						Path:       "/path/",
						CreateAt:   now,
						UpdateAt:   now,
						Urn:        "Urn2",
						Statements: &[]api.Statement{},
					},
					CreateAt: now,
				},
				{
					Policy: &api.Policy{
						ID:         "PolicyID1",
						Name:       "Name1",
						Org:        "org1",
						Path:       "/path/",
						CreateAt:   now,
						UpdateAt:   now,
						Urn:        "Urn1",
						Statements: &[]api.Statement{},
					},
					CreateAt: now.Add(-1),
				},
			},
		},
		"ErrorCase": {
			relations: &relations{
				policies: []Policy{
					{
						ID:       "PolicyID1",
						Name:     "Name1",
						Org:      "org1",
						Path:     "/path/",
						CreateAt: now.UnixNano(),
						UpdateAt: now.UnixNano(),
						Urn:      "Urn1",
					},
					{
						ID:       "PolicyID2",
						Name:     "Name2",
						Org:      "org1",
						Path:     "/path/",
						CreateAt: now.UnixNano(),
						UpdateAt: now.UnixNano(),
						Urn:      "Urn2",
					},
				},
				groupID:        "GroupID",
				createAt:       []int64{now.UnixNano(), now.UnixNano()},
				policyNotFound: true,
			},
			statements: []Statement{},
			groupID:    "GroupID",
			filter:     testFilter,
			expectedError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "Code: PolicyNotFound, Message: Policy with id PolicyID1 not found",
			},
		},
	}

	for n, test := range testcases {
		cleanPolicyTable(t, n)
		cleanGroupPolicyRelationTable(t, n)

		// Insert previous data
		if test.relations != nil {
			for i, policy := range test.relations.policies {
				insertGroupPolicyRelation(t, n, test.relations.groupID, policy.ID, test.relations.createAt[i])
				if !test.relations.policyNotFound {
					insertPolicy(t, n, policy, test.statements)
				}
			}

		}

		receivedPolicies, total, err := repoDB.GetAttachedPolicies(test.groupID, test.filter)
		if test.expectedError != nil {
			dbError, _ := err.(*database.Error)
			assert.Equal(t, test.expectedError, dbError, "Error in test case %v", n)
		} else {
			assert.Nil(t, err, "Error in test case %v", n)

			// Check total
			assert.Equal(t, len(test.expectedResponse), total, "Error in test case %v", n)

			// Check response
			for i, r := range receivedPolicies {
				assert.Equal(t, test.expectedResponse[i].GetGroup(), r.GetGroup(), "Error in test case %v", n)
				assert.Equal(t, test.expectedResponse[i].GetPolicy(), r.GetPolicy(), "Error in test case %v", n)
				assert.Equal(t, test.expectedResponse[i].GetDate(), r.GetDate(), "Error in test case %v", n)
			}
		}
	}
}
