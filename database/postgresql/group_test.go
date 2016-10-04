package postgresql

import (
	"testing"
	"time"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/database"
	"github.com/kylelemons/godebug/pretty"
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
		cleanGroupTable()

		// Insert previous data
		if test.previousGroup != nil {
			err := insertGroup(*test.previousGroup)
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting previous data: %v", n, err)
				continue
			}
		}
		// Call to repository to store group
		storedGroup, err := repoDB.AddGroup(*test.groupToCreate)
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
			if diff := pretty.Compare(storedGroup, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
			// Check database
			groupNumber, err := getGroupsCountFiltered(test.groupToCreate.ID, test.groupToCreate.Name, test.groupToCreate.Path,
				test.groupToCreate.CreateAt.UnixNano(), test.groupToCreate.UpdateAt.UnixNano(), test.groupToCreate.Urn, test.groupToCreate.Org)
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error counting groups: %v", n, err)
				continue
			}
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
		cleanGroupTable()

		// Insert previous data
		if test.previousGroup != nil {
			err := insertGroup(*test.previousGroup)
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting previous data: %v", n, err)
				continue
			}
		}

		// Call to repository to get group
		receivedGroup, err := repoDB.GetGroupByName(test.org, test.name)
		if test.expectedError != nil {
			dbError, ok := err.(*database.Error)
			if !ok || dbError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", n, err)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error: %v", n, err)
				continue
			}
			// Check response
			if diff := pretty.Compare(receivedGroup, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
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
		cleanGroupTable()

		// Insert previous data
		if test.previousGroup != nil {
			err := insertGroup(*test.previousGroup)
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting previous data: %v", n, err)
				continue
			}
		}

		// Call to repository to get group
		receivedGroup, err := repoDB.GetGroupById(test.groupID)
		if test.expectedError != nil {
			dbError, ok := err.(*database.Error)
			if !ok || dbError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", n, err)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error: %v", n, err)
				continue
			}
			// Check response
			if diff := pretty.Compare(receivedGroup, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
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
		cleanGroupTable()

		// Insert previous data
		if test.previousGroups != nil {
			for _, previousGroup := range test.previousGroups {
				if err := insertGroup(previousGroup); err != nil {
					t.Errorf("Test %v failed. Unexpected error inserting previous groups: %v", n, err)
					continue
				}
			}
		}
		// Call to repository to get groups
		receivedGroups, total, err := repoDB.GetGroupsFiltered(test.filter)
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error: %v", n, err)
			continue
		}
		// Check total
		if total != len(test.expectedResponse) {
			t.Errorf("Test %v failed. Received different total elements: %v", n, total)
			continue
		}
		// Check response
		if diff := pretty.Compare(receivedGroups, test.expectedResponse); diff != "" {
			t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
			continue
		}
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
		cleanGroupTable()

		// Insert previous data
		if test.previousGroups != nil {
			for _, previousGroup := range test.previousGroups {
				err := insertGroup(previousGroup)
				if err != nil {
					t.Errorf("Test %v failed. Unexpected error inserting previous data: %v", n, err)
					continue
				}
			}
		}

		// Call to repository to update group
		updatedGroup, err := repoDB.UpdateGroup(*test.groupToUpdate)
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
			if diff := pretty.Compare(updatedGroup, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
			// Check database
			groupNumber, err := getGroupsCountFiltered(test.expectedResponse.ID, test.expectedResponse.Name, test.expectedResponse.Path,
				test.expectedResponse.CreateAt.UnixNano(), test.expectedResponse.UpdateAt.UnixNano(), test.expectedResponse.Urn,
				test.expectedResponse.Org)
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error counting groups: %v", n, err)
				continue
			}
			if groupNumber != 1 {
				t.Fatalf("Test %v failed. Received different group number: %v", n, groupNumber)
				continue
			}
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
		cleanGroupTable()
		cleanGroupUserRelationTable()
		cleanGroupPolicyRelationTable()

		// Insert previous data
		if test.previousGroups != nil {
			for _, g := range test.previousGroups {
				if err := insertGroup(g); err != nil {
					t.Errorf("Test %v failed. Unexpected error inserting previous group %v: %v", n, g.ID, err)
					continue
				}
			}
		}
		if test.userRelations != nil {
			for _, rel := range test.userRelations {
				if err := insertGroupUserRelation(rel.userID, rel.groupID, rel.CreateAt); err != nil {
					t.Errorf("Test %v failed. Unexpected error inserting previous group user relations: %v", n, err)
					continue
				}
			}
		}
		if test.policyRelations != nil {
			for _, rel := range test.policyRelations {
				if err := insertGroupPolicyRelation(rel.groupID, rel.policyID, rel.CreateAt); err != nil {
					t.Errorf("Test %v failed. Unexpected error inserting previous group policy relations: %v", n, err)
					continue
				}
			}
		}
		// Call to repository to remove group
		err := repoDB.RemoveGroup(test.groupToDelete)

		// Check database
		groupNumber, err := getGroupsCountFiltered(test.groupToDelete, "", "", 0, 0, "", "")
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error counting groups: %v", n, err)
			continue
		}
		if groupNumber != 0 {
			t.Errorf("Test %v failed. Received different group number: %v", n, groupNumber)
			continue
		}

		// Check total groups
		totalGroupNumber, err := getGroupsCountFiltered("", "", "", 0, 0, "", "")
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error counting total groups: %v", n, err)
			continue
		}
		if totalGroupNumber != 1 {
			t.Errorf("Test %v failed. Received different total group number: %v", n, totalGroupNumber)
			continue
		}

		// Check group user relations
		relations, err := getGroupUserRelations(test.groupToDelete, "")
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error counting group user relations: %v", n, err)
			continue
		}
		if relations != 0 {
			t.Errorf("Test %v failed. Received different group user relation number: %v", n, relations)
			continue
		}

		// Check total group user relations
		totalRelations, err := getGroupUserRelations("", "")
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error counting total group user relations: %v", n, err)
			continue
		}
		if totalRelations != 1 {
			t.Errorf("Test %v failed. Received different total group user relation number: %v", n, totalRelations)
			continue
		}

		// Check group policy relations
		relations, err = getGroupPolicyRelationCount("", test.groupToDelete)
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error counting group policy relations: %v", n, err)
			continue
		}
		if relations != 0 {
			t.Errorf("Test %v failed. Received different group policy relation number: %v", n, relations)
			continue
		}

		// Check total group policy relations
		totalRelations, err = getGroupPolicyRelationCount("", "")
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error counting total group policy relations: %v", n, err)
			continue
		}
		if totalRelations != 1 {
			t.Errorf("Test %v failed. Received different total group policy relation number: %v", n, totalRelations)
			continue
		}
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
				Message: "pq: null value in column user_id violates not-null constraint",
			},
		},
	}

	for n, test := range testcases {
		// Clean GroupUserRelation database
		cleanGroupUserRelationTable()

		// Call to repository to store member
		err := repoDB.AddMember(test.userID, test.groupID)
		if test.expectedError != nil {
			dbError, ok := err.(*database.Error)
			if !ok || dbError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", n, err)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error: %v", n, err)
				continue
			}

			// Check database
			relations, err := getGroupUserRelations(test.groupID, test.userID)
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error counting relations: %v", n, err)
				continue
			}
			if relations != 1 {
				t.Errorf("Test %v failed. Received different relations number: %v", n, relations)
				continue
			}
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
		cleanGroupUserRelationTable()

		// Insert previous data
		if test.relation != nil {
			if err := insertGroupUserRelation(test.relation.userID, test.relation.groupID, test.relation.createAt); err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting previous group user relations: %v", n, err)
				continue
			}
		}

		// Call to repository to remove member
		err := repoDB.RemoveMember(test.userID, test.groupID)

		if err != nil {
			t.Errorf("Test %v failed. Unexpected error: %v", n, err)
			continue
		}

		// Check database
		relations, err := getGroupUserRelations(test.groupID, test.userID)
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error counting relations: %v", n, err)
			continue
		}
		if relations != 0 {
			t.Errorf("Test %v failed. Received different relations number: %v", n, relations)
			continue
		}
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
		cleanGroupUserRelationTable()

		// Insert previous data
		if test.relation != nil {
			if err := insertGroupUserRelation(test.relation.userID, test.relation.groupID, test.relation.createAt); err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting previous group user relations: %v", n, err)
				continue
			}
		}

		isMember, err := repoDB.IsMemberOfGroup(test.member, test.group)

		if err != nil {
			t.Errorf("Test %v failed. Unexpected error: %v", n, err)
			continue
		}
		// Check response
		if diff := pretty.Compare(isMember, test.isMember); diff != "" {
			t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
			continue
		}
	}
}

func TestPostgresRepo_GetGroupMembers(t *testing.T) {
	type relations struct {
		users        []User
		groupID      string
		createAt     int64
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
		expectedResponse []GroupUser
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
				createAt: now.UnixNano(),
			},
			groupID: "GroupID",
			filter:  testFilter,
			expectedResponse: []GroupUser{
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
		"ErrorCase": {
			relations: &relations{
				users: []User{
					{
						ID: "UserID1",
					},
				},
				groupID:      "GroupID",
				userNotFound: true,
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
		cleanUserTable()
		cleanGroupUserRelationTable()

		// Insert previous data
		if test.relations != nil {
			for _, user := range test.relations.users {
				if err := insertGroupUserRelation(user.ID, test.relations.groupID, test.relations.createAt); err != nil {
					t.Errorf("Test %v failed. Unexpected error inserting previous group user relations: %v", n, err)
					continue
				}
				if !test.relations.userNotFound {
					if err := insertUser(user); err != nil {
						t.Errorf("Test %v failed. Unexpected error inserting previous data: %v", n, err)
						continue
					}
				}
			}

		}

		receivedUsers, total, err := repoDB.GetGroupMembers(test.groupID, test.filter)
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
			// Check total
			if total != len(test.expectedResponse) {
				t.Errorf("Test %v failed. Received different total elements: %v", n, total)
				continue
			}
			// Check response
			if diff := pretty.Compare(receivedUsers, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
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
				Message: "pq: null value in column user_id violates not-null constraint",
			},
		},
	}

	for n, test := range testcases {
		// Clean GroupPolicyRelation database
		cleanGroupPolicyRelationTable()

		// Call to repository to attach policy
		err := repoDB.AttachPolicy(test.groupID, test.policyID)
		if test.expectedError != nil {
			dbError, ok := err.(*database.Error)
			if !ok || dbError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", n, err)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error: %v", n, err)
				continue
			}

			// Check database
			relations, err := getGroupPolicyRelationCount(test.policyID, test.groupID)
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error counting relations: %v", n, err)
				continue
			}
			if relations != 1 {
				t.Errorf("Test %v failed. Received different relations number: %v", n, relations)
				continue
			}
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
		cleanGroupPolicyRelationTable()

		// Insert previous data
		if test.relation != nil {
			if err := insertGroupPolicyRelation(test.relation.groupID, test.relation.policyID, test.relation.createAt); err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting previous group policy relations: %v", n, err)
				continue
			}
		}

		// Call to repository to detach policy
		err := repoDB.DetachPolicy(test.groupID, test.policyID)

		if err != nil {
			t.Errorf("Test %v failed. Unexpected error: %v", n, err)
			continue
		}

		// Check database
		relations, err := getGroupPolicyRelationCount(test.policyID, test.groupID)
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error counting relations: %v", n, err)
			continue
		}
		if relations != 0 {
			t.Errorf("Test %v failed. Received different relations number: %v", n, relations)
			continue
		}
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
		cleanGroupPolicyRelationTable()

		// Insert previous data
		if test.relation != nil {
			if err := insertGroupPolicyRelation(test.relation.groupID, test.relation.policyID, test.relation.createAt); err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting previous group policy relations: %v", n, err)
				continue
			}
		}

		// Call repository to check if policy is attached to group
		result, err := repoDB.IsAttachedToGroup(test.groupID, test.policyID)

		if err != nil {
			t.Errorf("Test %v failed. Unexpected error: %v", n, err)
			continue
		}

		if result != test.expectedResult {
			t.Errorf("Test %v failed. Received %v, expected %v", n, result, test.expectedResult)
			continue
		}
	}
}

func TestPostgresRepo_GetAttachedPolicies(t *testing.T) {
	type relations struct {
		policies       []Policy
		groupID        string
		createAt       int64
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
		expectedResponse []PolicyGroup
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
				createAt: now.UnixNano(),
			},
			statements: []Statement{},
			groupID:    "GroupID",
			filter:     testFilter,
			expectedResponse: []PolicyGroup{
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
					CreateAt: now,
				},
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
				createAt:       now.UnixNano(),
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
		cleanPolicyTable()
		cleanGroupPolicyRelationTable()

		// Insert previous data
		if test.relations != nil {
			for _, policy := range test.relations.policies {
				if err := insertGroupPolicyRelation(test.relations.groupID, policy.ID, test.relations.createAt); err != nil {
					t.Errorf("Test %v failed. Unexpected error inserting previous group policy relations: %v", n, err)
					continue
				}
				if !test.relations.policyNotFound {
					if err := insertPolicy(policy, test.statements); err != nil {
						t.Errorf("Test %v failed. Unexpected error inserting previous data: %v", n, err)
						continue
					}
				}
			}

		}

		receivedPolicies, total, err := repoDB.GetAttachedPolicies(test.groupID, test.filter)
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
			// Check total
			if total != len(test.expectedResponse) {
				t.Errorf("Test %v failed. Received different total elements: %v", n, total)
				continue
			}
			// Check response
			if diff := pretty.Compare(receivedPolicies, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
		}
	}
}
