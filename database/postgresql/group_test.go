package postgresql

import (
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/database"
)

func TestPostgresRepo_AddGroup(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousGroup *api.Group
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
				Org:      "Org",
			},
			expectedResponse: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "urn",
				CreateAt: now,
				Org:      "Org",
			},
		},
		"ErrorCasegroupAlreadyExist": {
			previousGroup: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "urn",
				CreateAt: now,
				Org:      "Org",
			},
			groupToCreate: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "urn",
				CreateAt: now,
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
			err := insertGroup(test.previousGroup.ID, test.previousGroup.Name, test.previousGroup.Path,
				test.previousGroup.CreateAt.UnixNano(), test.previousGroup.Urn, test.previousGroup.Org)
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
				test.groupToCreate.CreateAt.UnixNano(), test.groupToCreate.Urn, test.groupToCreate.Org)
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
		previousGroup *api.Group
		// Postgres Repo Args
		org  string
		name string
		// Expected result
		expectedResponse *api.Group
		expectedError    *database.Error
	}{
		"OkCase": {
			previousGroup: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Urn",
				CreateAt: now,
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
				Org:      "Org",
			},
		},
		"ErrorCaseGroupNotExist": {
			previousGroup: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Urn",
				CreateAt: now,
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
			err := insertGroup(test.previousGroup.ID, test.previousGroup.Name, test.previousGroup.Path,
				test.previousGroup.CreateAt.UnixNano(), test.previousGroup.Urn, test.previousGroup.Org)
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

func TestPostgresRepo_UpdateGroup(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousGroups []api.Group
		// Postgres Repo Args
		groupToUpdate *api.Group
		newName       string
		newPath       string
		newUrn        string
		// Expected result
		expectedResponse *api.Group
		expectedError    *database.Error
	}{
		"OkCase": {
			previousGroups: []api.Group{
				{
					ID:       "GroupID",
					Name:     "Name",
					Path:     "Path",
					Urn:      "Urn",
					CreateAt: now,
					Org:      "Org",
				},
			},
			groupToUpdate: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Urn",
				CreateAt: now,
				Org:      "Org",
			},
			newName: "NewName",
			newPath: "NewPath",
			newUrn:  "NewUrn",
			expectedResponse: &api.Group{
				ID:       "GroupID",
				Name:     "NewName",
				Path:     "NewPath",
				Urn:      "NewUrn",
				CreateAt: now,
				Org:      "Org",
			},
		},
		"ErrorCaseDuplicateUrn": {
			previousGroups: []api.Group{
				{
					ID:       "GroupID",
					Name:     "Name",
					Path:     "Path",
					Urn:      "Urn",
					CreateAt: now,
					Org:      "Org",
				},
				{
					ID:       "GroupID2",
					Name:     "Name2",
					Path:     "Path2",
					Urn:      "Fail",
					CreateAt: now,
					Org:      "Org2",
				},
			},
			groupToUpdate: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Urn",
				CreateAt: now,
				Org:      "Org",
			},
			newName: "NewName",
			newPath: "NewPath",
			newUrn:  "Fail",
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
				err := insertGroup(previousGroup.ID, previousGroup.Name, previousGroup.Path,
					previousGroup.CreateAt.UnixNano(), previousGroup.Urn, previousGroup.Org)
				if err != nil {
					t.Errorf("Test %v failed. Unexpected error inserting previous data: %v", n, err)
					continue
				}
			}
		}

		// Call to repository to update group
		updatedGroup, err := repoDB.UpdateGroup(*test.groupToUpdate, test.newName, test.newPath, test.newUrn)
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
				test.expectedResponse.CreateAt.UnixNano(), test.expectedResponse.Urn, test.expectedResponse.Org)
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
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousGroup *api.Group
		relation      *struct {
			user_id       string
			group_ids     []string
			groupNotFound bool
		}
		// Postgres Repo Args
		groupToDelete string
	}{
		"OkCase": {
			previousGroup: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Urn",
				CreateAt: now,
				Org:      "Org",
			},
			relation: &struct {
				user_id       string
				group_ids     []string
				groupNotFound bool
			}{
				user_id:   "UserID",
				group_ids: []string{"GroupID"},
			},
			groupToDelete: "GroupID",
		},
	}

	for n, test := range testcases {
		cleanGroupTable()
		cleanGroupUserRelationTable()

		// Insert previous data
		if test.previousGroup != nil {
			if err := insertGroup(test.previousGroup.ID, test.previousGroup.Name, test.previousGroup.Path,
				test.previousGroup.CreateAt.Unix(), test.previousGroup.Urn, test.previousGroup.Org); err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting previous group: %v", n, err)
				continue
			}
		}
		if test.relation != nil {
			for _, id := range test.relation.group_ids {
				if err := insertGroupUserRelation(test.relation.user_id, id); err != nil {
					t.Errorf("Test %v failed. Unexpected error inserting previous group user relations: %v", n, err)
					continue
				}
			}
		}
		// Call to repository to remove group
		err := repoDB.RemoveGroup(test.groupToDelete)

		// Check database
		groupNumber, err := getGroupsCountFiltered(test.groupToDelete, "", "",
			0, "", "")
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error counting groups: %v", n, err)
			continue
		}
		if groupNumber != 0 {
			t.Errorf("Test %v failed. Received different group number: %v", n, groupNumber)
			continue
		}

		relations, err := getGroupUserRelations(test.previousGroup.ID, "")
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
	testcases := map[string]struct {
		// Previous data
		relation *struct {
			user_id  string
			group_id string
		}
		// Postgres Repo Args
		group  string
		member string
		// Expected result
		isMember bool
	}{
		"OkCaseIsMember": {
			relation: &struct {
				user_id  string
				group_id string
			}{
				user_id:  "UserID",
				group_id: "GroupID",
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
			if err := insertGroupUserRelation(test.relation.user_id, test.relation.group_id); err != nil {
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

func TestPostgresRepo_GetGroupById(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousGroup *api.Group
		// Postgres Repo Args
		groupID string
		// Expected result
		expectedResponse *api.Group
		expectedError    *database.Error
	}{
		"OkCase": {
			previousGroup: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Urn",
				CreateAt: now,
				Org:      "Org",
			},
			groupID: "GroupID",
			expectedResponse: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Urn",
				CreateAt: now,
				Org:      "Org",
			},
		},
		"ErrorCaseGroupNotExist": {
			previousGroup: &api.Group{
				ID:       "GroupID",
				Name:     "Name",
				Path:     "Path",
				Urn:      "Urn",
				CreateAt: now,
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
			err := insertGroup(test.previousGroup.ID, test.previousGroup.Name, test.previousGroup.Path,
				test.previousGroup.CreateAt.UnixNano(), test.previousGroup.Urn, test.previousGroup.Org)
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
