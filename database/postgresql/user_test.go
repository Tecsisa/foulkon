package postgresql

import (
	"testing"
	"time"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/database"
	"github.com/kylelemons/godebug/pretty"
)

func TestPostgresRepo_AddUser(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousUser *User
		// Postgres Repo Args
		userToCreate *api.User
		// Expected result
		expectedResponse *api.User
		expectedError    *database.Error
	}{
		"OkCase": {
			userToCreate: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
				UpdateAt:   now,
			},
			expectedResponse: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
				UpdateAt:   now,
			},
		},
		"ErrorCaseUserAlreadyExist": {
			previousUser: &User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now.UnixNano(),
				UpdateAt:   now.UnixNano(),
			},
			userToCreate: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
				UpdateAt:   now,
			},
			expectedError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "pq: duplicate key value violates unique constraint \"users_pkey\"",
			},
		},
	}

	for n, test := range testcases {
		// Clean user database
		cleanUserTable()

		// Insert previous data
		if test.previousUser != nil {
			if err := insertUser(*test.previousUser); err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting previous users: %v", n, err)
				continue
			}
		}
		// Call to repository to store an user
		storedUser, err := repoDB.AddUser(*test.userToCreate)
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
			if diff := pretty.Compare(storedUser, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
			// Check database
			userNumber, err := getUsersCountFiltered(test.expectedResponse.ID, test.expectedResponse.ExternalID, test.expectedResponse.Path,
				test.expectedResponse.CreateAt.UnixNano(), test.expectedResponse.UpdateAt.UnixNano(), test.expectedResponse.Urn, "")
			if err != nil {
				t.Errorf("Test %v failed. Unexpected error counting users: %v", n, err)
				continue
			}
			if userNumber != 1 {
				t.Errorf("Test %v failed. Received different user number: %v", n, userNumber)
				continue
			}

		}

	}
}

func TestPostgresRepo_GetUserByExternalID(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousUser *User
		// Postgres Repo Args
		externalID string
		// Expected result
		expectedResponse *api.User
		expectedError    *database.Error
	}{
		"OkCase": {
			previousUser: &User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now.UnixNano(),
				UpdateAt:   now.UnixNano(),
			},
			externalID: "ExternalID",
			expectedResponse: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
				UpdateAt:   now,
			},
		},
		"ErrorCaseUserNotExist": {
			previousUser: &User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now.UnixNano(),
				UpdateAt:   now.UnixNano(),
			},
			externalID: "NotExist",
			expectedError: &database.Error{
				Code:    database.USER_NOT_FOUND,
				Message: "User with externalId NotExist not found",
			},
		},
	}

	for n, test := range testcases {
		// Clean user database
		cleanUserTable()

		// Insert previous data
		if test.previousUser != nil {
			if err := insertUser(*test.previousUser); err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting previous users: %v", n, err)
				continue
			}
		}
		// Call to repository to get an user
		receivedUser, err := repoDB.GetUserByExternalID(test.externalID)
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
			if diff := pretty.Compare(receivedUser, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
		}

	}
}

func TestPostgresRepo_GetUserByID(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousUser *User
		// Postgres Repo Args
		userID string
		// Expected result
		expectedResponse *api.User
		expectedError    *database.Error
	}{
		"OkCase": {
			previousUser: &User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now.UnixNano(),
				UpdateAt:   now.UnixNano(),
			},
			userID: "UserID",
			expectedResponse: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
				UpdateAt:   now,
			},
		},
		"ErrorCaseUserNotExist": {
			previousUser: &User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now.UnixNano(),
				UpdateAt:   now.UnixNano(),
			},
			userID: "NotExist",
			expectedError: &database.Error{
				Code:    database.USER_NOT_FOUND,
				Message: "User with id NotExist not found",
			},
		},
	}

	for n, test := range testcases {
		// Clean user database
		cleanUserTable()

		// Insert previous data
		if test.previousUser != nil {
			if err := insertUser(*test.previousUser); err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting previous users: %v", n, err)
				continue
			}
		}
		// Call to repository to get an user
		receivedUser, err := repoDB.GetUserByID(test.userID)
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
			if diff := pretty.Compare(receivedUser, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}

		}

	}
}

func TestPostgresRepo_GetUsersFiltered(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousUsers []User
		// Postgres Repo Args
		filter *api.Filter
		// Expected result
		expectedResponse []api.User
	}{
		"OkCase1": {
			previousUsers: []User{
				{
					ID:         "UserID1",
					ExternalID: "ExternalID1",
					Path:       "Path123",
					Urn:        "urn1",
					CreateAt:   now.UnixNano(),
					UpdateAt:   now.UnixNano(),
				},
				{
					ID:         "UserID2",
					ExternalID: "ExternalID2",
					Path:       "Path456",
					Urn:        "urn2",
					CreateAt:   now.UnixNano(),
					UpdateAt:   now.UnixNano(),
				},
			},
			filter: &api.Filter{
				PathPrefix: "Path",
				Offset:     0,
				Limit:      20,
			},
			expectedResponse: []api.User{
				{
					ID:         "UserID1",
					ExternalID: "ExternalID1",
					Path:       "Path123",
					Urn:        "urn1",
					CreateAt:   now,
					UpdateAt:   now,
				},
				{
					ID:         "UserID2",
					ExternalID: "ExternalID2",
					Path:       "Path456",
					Urn:        "urn2",
					CreateAt:   now,
					UpdateAt:   now,
				},
			},
		},
		"OkCase2": {
			previousUsers: []User{
				{
					ID:         "UserID1",
					ExternalID: "ExternalID1",
					Path:       "Path123",
					Urn:        "urn1",
					CreateAt:   now.UnixNano(),
					UpdateAt:   now.UnixNano(),
				},
				{
					ID:         "UserID2",
					ExternalID: "ExternalID2",
					Path:       "Path456",
					Urn:        "urn2",
					CreateAt:   now.UnixNano(),
					UpdateAt:   now.UnixNano(),
				},
			},
			filter: &api.Filter{
				PathPrefix: "Path123",
				Offset:     0,
				Limit:      20,
			},
			expectedResponse: []api.User{
				{
					ID:         "UserID1",
					ExternalID: "ExternalID1",
					Path:       "Path123",
					Urn:        "urn1",
					CreateAt:   now,
					UpdateAt:   now,
				},
			},
		},
		"OkCase3": {
			previousUsers: []User{
				{
					ID:         "UserID1",
					ExternalID: "ExternalID1",
					Path:       "Path123",
					Urn:        "urn1",
					CreateAt:   now.UnixNano(),
					UpdateAt:   now.UnixNano(),
				},
				{
					ID:         "UserID2",
					ExternalID: "ExternalID2",
					Path:       "Path456",
					Urn:        "urn2",
					CreateAt:   now.UnixNano(),
					UpdateAt:   now.UnixNano(),
				},
			},
			filter: &api.Filter{
				PathPrefix: "NoPath",
				Offset:     0,
				Limit:      20,
			},
			expectedResponse: []api.User{},
		},
	}

	for n, test := range testcases {
		// Clean user database
		cleanUserTable()

		// Insert previous data
		if test.previousUsers != nil {
			for _, previousUser := range test.previousUsers {
				if err := insertUser(previousUser); err != nil {
					t.Errorf("Test %v failed. Unexpected error inserting previous users: %v", n, err)
					continue
				}
			}
		}
		// Call to repository to get users
		receivedUsers, total, err := repoDB.GetUsersFiltered(test.filter)
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

func TestPostgresRepo_UpdateUser(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousUser *User
		// Postgres Repo Args
		userToUpdate *api.User
		// Expected result
		expectedResponse *api.User
	}{
		"OkCase": {
			previousUser: &User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "OldPath",
				Urn:        "Oldurn",
				CreateAt:   now.UnixNano(),
				UpdateAt:   now.UnixNano(),
			},
			userToUpdate: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "NewPath",
				Urn:        "NewUrn",
				CreateAt:   now,
				UpdateAt:   now,
			},
			expectedResponse: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "NewPath",
				Urn:        "NewUrn",
				CreateAt:   now,
				UpdateAt:   now,
			},
		},
	}

	for n, test := range testcases {
		// Clean user database
		cleanUserTable()

		// Insert previous data
		if test.previousUser != nil {
			if err := insertUser(*test.previousUser); err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting previous users: %v", n, err)
				continue
			}
		}
		// Call to repository to update an user
		updatedUser, err := repoDB.UpdateUser(*test.userToUpdate)

		if err != nil {
			t.Errorf("Test %v failed. Unexpected error: %v", n, err)
			continue
		}
		// Check response
		if diff := pretty.Compare(updatedUser, test.expectedResponse); diff != "" {
			t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
			continue
		}
		// Check database
		userNumber, err := getUsersCountFiltered(test.expectedResponse.ID, test.expectedResponse.ExternalID, test.expectedResponse.Path,
			test.expectedResponse.CreateAt.UnixNano(), test.expectedResponse.UpdateAt.UnixNano(), test.expectedResponse.Urn, "")
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error counting users: %v", n, err)
			continue
		}
		if userNumber != 1 {
			t.Fatalf("Test %v failed. Received different user number: %v", n, userNumber)
			continue
		}

	}
}

func TestPostgresRepo_RemoveUser(t *testing.T) {
	type relation struct {
		userID        string
		groupID       string
		createAt      int64
		groupNotFound bool
	}
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousUsers []User
		relations     []relation
		// Postgres Repo Args
		userToDelete string
	}{
		"OkCase": {
			previousUsers: []User{
				{
					ID:         "UserID",
					ExternalID: "ExternalID",
					Path:       "OldPath",
					Urn:        "Oldurn",
					CreateAt:   now.UnixNano(),
					UpdateAt:   now.UnixNano(),
				},
				{
					ID:         "UserID2",
					ExternalID: "ExternalID2",
					Path:       "OldPath",
					Urn:        "Oldurn2",
					CreateAt:   now.UnixNano(),
					UpdateAt:   now.UnixNano(),
				},
			},
			relations: []relation{
				{
					userID:   "UserID",
					groupID:  "GroupID",
					createAt: now.UnixNano(),
				},
				{
					userID:   "UserID2",
					groupID:  "GroupID",
					createAt: now.UnixNano(),
				},
			},
			userToDelete: "UserID",
		},
	}

	for n, test := range testcases {
		// Clean user database
		cleanUserTable()
		cleanGroupUserRelationTable()

		// Insert previous data
		if test.previousUsers != nil {
			for _, usr := range test.previousUsers {
				if err := insertUser(usr); err != nil {
					t.Errorf("Test %v failed. Unexpected error inserting previous users: %v", n, err)
					continue
				}
			}
		}
		if test.relations != nil {
			for _, rel := range test.relations {
				if err := insertGroupUserRelation(rel.userID, rel.groupID, rel.createAt); err != nil {
					t.Errorf("Test %v failed. Unexpected error inserting previous group user relations: %v", n, err)
					continue
				}
			}
		}
		// Call to repository to remove user
		err := repoDB.RemoveUser(test.userToDelete)

		// Check database
		userNumber, err := getUsersCountFiltered(test.userToDelete, "", "", 0, 0, "", "")
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error counting users: %v", n, err)
			continue
		}
		if userNumber != 0 {
			t.Errorf("Test %v failed. Received different user number: %v", n, userNumber)
			continue
		}

		// Check total users
		totalUserNumber, err := getUsersCountFiltered("", "", "", 0, 0, "", "")
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error counting total users: %v", n, err)
			continue
		}
		if totalUserNumber != 1 {
			t.Errorf("Test %v failed. Received different total user number: %v", n, totalUserNumber)
			continue
		}

		// Check user deleted relations
		relations, err := getGroupUserRelations("", test.userToDelete)
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error counting group user relations: %v", n, err)
			continue
		}
		if relations != 0 {
			t.Errorf("Test %v failed. Received different group user relation number: %v", n, relations)
			continue
		}

		// Check total user relations
		totalRelations, err := getGroupUserRelations("", "")
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error counting total group user relations: %v", n, err)
			continue
		}
		if totalRelations != 1 {
			t.Errorf("Test %v failed. Received different total group user relation number: %v", n, totalRelations)
			continue
		}

	}
}

func TestPostgresRepo_GetGroupsByUserID(t *testing.T) {
	type relation struct {
		userID        string
		groups        []Group
		createAt      int64
		groupNotFound bool
	}
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		relation *relation
		// Postgres Repo Args
		userID string
		filter *api.Filter
		// Expected result
		expectedResponse []GroupUser
		expectedError    *database.Error
	}{
		"OkCase": {
			relation: &relation{
				userID: "UserID",
				groups: []Group{
					{
						ID:       "GroupID1",
						Name:     "Name1",
						Path:     "Path1",
						Urn:      "urn1",
						CreateAt: now.UnixNano(),
						UpdateAt: now.UnixNano(),
						Org:      "Org",
					},
					{
						ID:       "GroupID2",
						Name:     "Name2",
						Path:     "Path2",
						Urn:      "urn2",
						CreateAt: now.UnixNano(),
						UpdateAt: now.UnixNano(),
						Org:      "Org",
					},
				},
				createAt: now.UnixNano(),
			},
			userID: "UserID",
			filter: testFilter,
			expectedResponse: []GroupUser{
				{
					Group: &api.Group{
						ID:       "GroupID1",
						Name:     "Name1",
						Path:     "Path1",
						Urn:      "urn1",
						CreateAt: now,
						UpdateAt: now,
						Org:      "Org",
					},
					CreateAt: now,
				},
				{
					Group: &api.Group{
						ID:       "GroupID2",
						Name:     "Name2",
						Path:     "Path2",
						Urn:      "urn2",
						CreateAt: now,
						UpdateAt: now,
						Org:      "Org",
					},
					CreateAt: now,
				},
			},
		},
		"ErrorCase": {
			relation: &relation{
				userID: "UserID",
				groups: []Group{
					{
						ID: "GroupID1",
					},
					{
						ID: "GroupID2",
					},
				},
				createAt:      now.UnixNano(),
				groupNotFound: true,
			},
			userID: "UserID",
			filter: testFilter,
			expectedError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "Code: GroupNotFound, Message: Group with id GroupID1 not found",
			},
		},
	}

	for n, test := range testcases {
		// Clean database
		cleanUserTable()
		cleanGroupTable()
		cleanGroupUserRelationTable()

		// Insert previous data
		if test.relation != nil {
			for _, group := range test.relation.groups {
				if err := insertGroupUserRelation(test.relation.userID, group.ID, test.relation.createAt); err != nil {
					t.Errorf("Test %v failed. Unexpected error inserting prevoius group user relations: %v", n, err)
					continue
				}
				if !test.relation.groupNotFound {
					if err := insertGroup(group); err != nil {
						t.Errorf("Test %v failed. Unexpected error inserting previous data: %v", n, err)
						continue
					}
				}
			}
		}
		// Call to repository to get groups associated
		receivedUsers, total, err := repoDB.GetGroupsByUserID(test.userID, test.filter)
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
