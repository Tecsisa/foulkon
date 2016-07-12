package postgresql

import (
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/database"
)

func TestPostgresRepo_AddUser(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousUser *api.User
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
			},
			expectedResponse: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
			},
		},
		"ErrorCaseUserAlreadyExist": {
			previousUser: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
			},
			userToCreate: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
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
			if err := insertUser(test.previousUser.ID, test.previousUser.ExternalID, test.previousUser.Path,
				test.previousUser.CreateAt.Unix(), test.previousUser.Urn); err != nil {
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
				test.expectedResponse.CreateAt.UnixNano(), test.expectedResponse.Urn, "")
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
		previousUser *api.User
		// Postgres Repo Args
		externalID string
		// Expected result
		expectedResponse *api.User
		expectedError    *database.Error
	}{
		"OkCase": {
			previousUser: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
			},
			externalID: "ExternalID",
			expectedResponse: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
			},
		},
		"ErrorCaseUserNotExist": {
			previousUser: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
			},
			externalID: "NotExist",
			expectedError: &database.Error{
				Code:    database.USER_NOT_FOUND,
				Message: "User with ExternalID NotExist not found",
			},
		},
	}

	for n, test := range testcases {
		// Clean user database
		cleanUserTable()

		// Insert previous data
		if test.previousUser != nil {
			if err := insertUser(test.previousUser.ID, test.previousUser.ExternalID, test.previousUser.Path,
				test.previousUser.CreateAt.UnixNano(), test.previousUser.Urn); err != nil {
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
		previousUser *api.User
		// Postgres Repo Args
		userID string
		// Expected result
		expectedResponse *api.User
		expectedError    *database.Error
	}{
		"OkCase": {
			previousUser: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
			},
			userID: "UserID",
			expectedResponse: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
			},
		},
		"ErrorCaseUserNotExist": {
			previousUser: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "Path",
				Urn:        "urn",
				CreateAt:   now,
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
			if err := insertUser(test.previousUser.ID, test.previousUser.ExternalID, test.previousUser.Path,
				test.previousUser.CreateAt.UnixNano(), test.previousUser.Urn); err != nil {
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

func TestPostgresRepo_UpdateUser(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousUser *api.User
		// Postgres Repo Args
		userToUpdate *api.User
		newPath      string
		newUrn       string
		// Expected result
		expectedResponse *api.User
	}{
		"OkCase": {
			previousUser: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "OldPath",
				Urn:        "Oldurn",
				CreateAt:   now,
			},
			userToUpdate: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "OldPath",
				Urn:        "Oldurn",
				CreateAt:   now,
			},
			newPath: "NewPath",
			newUrn:  "NewUrn",
			expectedResponse: &api.User{
				ID:         "UserID",
				ExternalID: "ExternalID",
				Path:       "NewPath",
				Urn:        "NewUrn",
				CreateAt:   now,
			},
		},
	}

	for n, test := range testcases {
		// Clean user database
		cleanUserTable()

		// Insert previous data
		if test.previousUser != nil {
			if err := insertUser(test.previousUser.ID, test.previousUser.ExternalID, test.previousUser.Path,
				test.previousUser.CreateAt.UnixNano(), test.previousUser.Urn); err != nil {
				t.Errorf("Test %v failed. Unexpected error inserting previous users: %v", n, err)
				continue
			}
		}
		// Call to repository to update an user
		updatedUser, err := repoDB.UpdateUser(*test.userToUpdate, test.newPath, test.newUrn)

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
			test.expectedResponse.CreateAt.UnixNano(), test.expectedResponse.Urn, "")
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

func TestPostgresRepo_GetUsersFiltered(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		previousUsers []api.User
		// Postgres Repo Args
		pathPrefix string
		// Expected result
		expectedResponse []api.User
	}{
		"OkCase1": {
			previousUsers: []api.User{
				api.User{
					ID:         "UserID1",
					ExternalID: "ExternalID1",
					Path:       "Path123",
					Urn:        "urn1",
					CreateAt:   now,
				},
				api.User{
					ID:         "UserID2",
					ExternalID: "ExternalID2",
					Path:       "Path456",
					Urn:        "urn2",
					CreateAt:   now,
				},
			},
			pathPrefix: "Path",
			expectedResponse: []api.User{
				api.User{
					ID:         "UserID1",
					ExternalID: "ExternalID1",
					Path:       "Path123",
					Urn:        "urn1",
					CreateAt:   now,
				},
				api.User{
					ID:         "UserID2",
					ExternalID: "ExternalID2",
					Path:       "Path456",
					Urn:        "urn2",
					CreateAt:   now,
				},
			},
		},
		"OkCase2": {
			previousUsers: []api.User{
				api.User{
					ID:         "UserID1",
					ExternalID: "ExternalID1",
					Path:       "Path123",
					Urn:        "urn1",
					CreateAt:   now,
				},
				api.User{
					ID:         "UserID2",
					ExternalID: "ExternalID2",
					Path:       "Path456",
					Urn:        "urn2",
					CreateAt:   now,
				},
			},
			pathPrefix: "Path123",
			expectedResponse: []api.User{
				api.User{
					ID:         "UserID1",
					ExternalID: "ExternalID1",
					Path:       "Path123",
					Urn:        "urn1",
					CreateAt:   now,
				},
			},
		},
		"OkCase3": {
			previousUsers: []api.User{
				api.User{
					ID:         "UserID1",
					ExternalID: "ExternalID1",
					Path:       "Path123",
					Urn:        "urn1",
					CreateAt:   now,
				},
				api.User{
					ID:         "UserID2",
					ExternalID: "ExternalID2",
					Path:       "Path456",
					Urn:        "urn2",
					CreateAt:   now,
				},
			},
			pathPrefix:       "NoPath",
			expectedResponse: []api.User{},
		},
	}

	for n, test := range testcases {
		// Clean user database
		cleanUserTable()

		// Insert previous data
		if test.previousUsers != nil {
			for _, previousUser := range test.previousUsers {
				if err := insertUser(previousUser.ID, previousUser.ExternalID, previousUser.Path,
					previousUser.CreateAt.UnixNano(), previousUser.Urn); err != nil {
					t.Errorf("Test %v failed. Unexpected error inserting previous users: %v", n, err)
					continue
				}
			}
		}
		// Call to repository to get users
		receivedUsers, err := repoDB.GetUsersFiltered(test.pathPrefix)
		if err != nil {
			t.Errorf("Test %v failed. Unexpected error: %v", n, err)
			continue
		}
		// Check response
		if diff := pretty.Compare(receivedUsers, test.expectedResponse); diff != "" {
			t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
			continue
		}

	}
}

func TestPostgresRepo_GetGroupsByUserID(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// Previous data
		relation *struct {
			user_id       string
			groups        []api.Group
			groupNotFound bool
		}
		// Postgres Repo Args
		userID string
		// Expected result
		expectedResponse []api.Group
		expectedError    *database.Error
	}{
		"OkCase": {
			relation: &struct {
				user_id       string
				groups        []api.Group
				groupNotFound bool
			}{
				user_id: "UserID",
				groups: []api.Group{
					api.Group{
						ID:       "GroupID1",
						Name:     "Name1",
						Path:     "Path1",
						Urn:      "urn1",
						CreateAt: now,
						Org:      "Org",
					},
					api.Group{
						ID:       "GroupID2",
						Name:     "Name2",
						Path:     "Path2",
						Urn:      "urn2",
						CreateAt: now,
						Org:      "Org",
					},
				},
			},
			userID: "UserID",
			expectedResponse: []api.Group{
				api.Group{
					ID:       "GroupID1",
					Name:     "Name1",
					Path:     "Path1",
					Urn:      "urn1",
					CreateAt: now,
					Org:      "Org",
				},
				api.Group{
					ID:       "GroupID2",
					Name:     "Name2",
					Path:     "Path2",
					Urn:      "urn2",
					CreateAt: now,
					Org:      "Org",
				},
			},
		},
		"ErrorCase": {
			relation: &struct {
				user_id       string
				groups        []api.Group
				groupNotFound bool
			}{
				user_id: "UserID",
				groups: []api.Group{
					api.Group{
						ID: "GroupID1",
					},
					api.Group{
						ID: "GroupID2",
					},
				},
				groupNotFound: true,
			},
			userID: "UserID",
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
				if err := insertGroupUserRelation(test.relation.user_id, group.ID); err != nil {
					t.Errorf("Test %v failed. Unexpected error inserting prevoius group user relations: %v", n, err)
					continue
				}
				if !test.relation.groupNotFound {
					if err := insertGroup(group.ID, group.Name, group.Path,
						group.CreateAt.UnixNano(), group.Urn, group.Org); err != nil {
						t.Errorf("Test %v failed. Unexpected error inserting previous data: %v", n, err)
						continue
					}
				}
			}
		}
		// Call to repository to get groups associated
		receivedUsers, err := repoDB.GetGroupsByUserID(test.userID)
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
			if diff := pretty.Compare(receivedUsers, test.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
		}

	}
}
