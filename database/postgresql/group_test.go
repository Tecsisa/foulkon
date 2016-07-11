package postgresql

import (
	"github.com/kylelemons/godebug/pretty"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/database"
	"testing"
	"time"
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
			insertGroup(test.previousGroup.ID, test.previousGroup.Name, test.previousGroup.Path,
				test.previousGroup.CreateAt.UnixNano(), test.previousGroup.Urn, test.previousGroup.Org)
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
				t.Errorf("Test %v failed. Unexpected error counting users: %v", n, err)
				continue
			}
			if groupNumber != 1 {
				t.Errorf("Test %v failed. Received different user number: %v", n, groupNumber)
				continue
			}

		}

	}
}
