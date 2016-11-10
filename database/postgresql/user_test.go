package postgresql

import (
	"testing"
	"time"

	"github.com/Tecsisa/foulkon/api"
	"github.com/Tecsisa/foulkon/database"
	"github.com/stretchr/testify/assert"
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
		cleanUserTable(t, n)

		// Insert previous data
		if test.previousUser != nil {
			insertUser(t, n, *test.previousUser)
		}
		// Call to repository to store an user
		storedUser, err := repoDB.AddUser(*test.userToCreate)
		if test.expectedError != nil {
			dbError, _ := err.(*database.Error)
			assert.Equal(t, test.expectedError, dbError, "Error in test case %v", n)
		} else {
			assert.Nil(t, err, "Error in test case %v", n)

			// Check response
			assert.Equal(t, test.expectedResponse, storedUser, "Error in test case %v", n)
			// Check database
			userNumber := getUsersCountFiltered(t, n, test.expectedResponse.ID, test.expectedResponse.ExternalID, test.expectedResponse.Path,
				test.expectedResponse.CreateAt.UnixNano(), test.expectedResponse.UpdateAt.UnixNano(), test.expectedResponse.Urn, "")
			assert.Equal(t, 1, userNumber, "Error in test case %v", n)
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
		cleanUserTable(t, n)

		// Insert previous data
		if test.previousUser != nil {
			insertUser(t, n, *test.previousUser)
		}
		// Call to repository to get an user
		receivedUser, err := repoDB.GetUserByExternalID(test.externalID)
		if test.expectedError != nil {
			dbError, _ := err.(*database.Error)
			assert.Equal(t, test.expectedError, dbError, "Error in test case %v", n)
		} else {
			assert.Nil(t, err, "Error in test case %v", n)

			// Check response
			assert.Equal(t, test.expectedResponse, receivedUser, "Error in test case %v", n)
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
		cleanUserTable(t, n)

		// Insert previous data
		if test.previousUser != nil {
			insertUser(t, n, *test.previousUser)
		}
		// Call to repository to get an user
		receivedUser, err := repoDB.GetUserByID(test.userID)
		if test.expectedError != nil {
			dbError, _ := err.(*database.Error)
			assert.Equal(t, test.expectedError, dbError, "Error in test case %v", n)
		} else {
			assert.Nil(t, err, "Error in test case %v", n)

			// Check response
			assert.Equal(t, test.expectedResponse, receivedUser, "Error in test case %v", n)
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
				OrderBy:    "urn desc",
			},
			expectedResponse: []api.User{
				{
					ID:         "UserID2",
					ExternalID: "ExternalID2",
					Path:       "Path456",
					Urn:        "urn2",
					CreateAt:   now,
					UpdateAt:   now,
				},
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
		cleanUserTable(t, n)

		// Insert previous data
		if test.previousUsers != nil {
			for _, previousUser := range test.previousUsers {
				insertUser(t, n, previousUser)
			}
		}
		// Call to repository to get users
		receivedUsers, total, err := repoDB.GetUsersFiltered(test.filter)
		assert.Nil(t, err, "Error in test case %v", n)

		// Check total
		assert.Equal(t, len(test.expectedResponse), total, "Error in test case %v", n)

		// Check response
		assert.Equal(t, test.expectedResponse, receivedUsers, "Error in test case %v", n)
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
		cleanUserTable(t, n)

		// Insert previous data
		if test.previousUser != nil {
			insertUser(t, n, *test.previousUser)
		}
		// Call to repository to update an user
		updatedUser, err := repoDB.UpdateUser(*test.userToUpdate)
		assert.Nil(t, err, "Error in test case %v", n)

		// Check response
		assert.Equal(t, test.expectedResponse, updatedUser, "Error in test case %v", n)
		// Check database
		userNumber := getUsersCountFiltered(t, n, test.expectedResponse.ID, test.expectedResponse.ExternalID, test.expectedResponse.Path,
			test.expectedResponse.CreateAt.UnixNano(), test.expectedResponse.UpdateAt.UnixNano(), test.expectedResponse.Urn, "")
		assert.Equal(t, 1, userNumber, "Error in test case %v", n)
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
		cleanUserTable(t, n)
		cleanGroupUserRelationTable(t, n)

		// Insert previous data
		if test.previousUsers != nil {
			for _, usr := range test.previousUsers {
				insertUser(t, n, usr)
			}
		}
		if test.relations != nil {
			for _, rel := range test.relations {
				insertGroupUserRelation(t, n, rel.userID, rel.groupID, rel.createAt)
			}
		}
		// Call to repository to remove user
		err := repoDB.RemoveUser(test.userToDelete)
		assert.Nil(t, err, "Error in test case %v", n)

		// Check database
		userNumber := getUsersCountFiltered(t, n, test.userToDelete, "", "", 0, 0, "", "")
		assert.Equal(t, 0, userNumber, "Error in test case %v", n)

		// Check total users
		totalUserNumber := getUsersCountFiltered(t, n, "", "", "", 0, 0, "", "")
		assert.Equal(t, 1, totalUserNumber, "Error in test case %v", n)

		// Check user deleted relations
		relations := getGroupUserRelations(t, n, "", test.userToDelete)
		assert.Equal(t, 0, relations, "Error in test case %v", n)

		// Check total user relations
		totalRelations := getGroupUserRelations(t, n, "", "")
		assert.Equal(t, 1, totalRelations, "Error in test case %v", n)
	}
}

func TestPostgresRepo_GetGroupsByUserID(t *testing.T) {
	type relation struct {
		userID        string
		groups        []Group
		createAt      []int64
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
		expectedResponse []*GroupUser
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
				createAt: []int64{now.UnixNano(), now.UnixNano()},
			},
			userID: "UserID",
			filter: testFilter,
			expectedResponse: []*GroupUser{
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
		"OkCaseOrderBy": {
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
				createAt: []int64{now.UnixNano() - 1, now.UnixNano()},
			},
			userID: "UserID",
			filter: &api.Filter{
				OrderBy: "create_at desc",
			},
			expectedResponse: []*GroupUser{
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
					CreateAt: now.Add(-1),
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
				createAt:      []int64{now.UnixNano(), now.UnixNano()},
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
		cleanUserTable(t, n)
		cleanGroupTable(t, n)
		cleanGroupUserRelationTable(t, n)

		// Insert previous data
		if test.relation != nil {
			for i, group := range test.relation.groups {
				insertGroupUserRelation(t, n, test.relation.userID, group.ID, test.relation.createAt[i])
				if !test.relation.groupNotFound {
					insertGroup(t, n, group)
				}
			}
		}
		// Call to repository to get groups associated
		receivedUsers, total, err := repoDB.GetGroupsByUserID(test.userID, test.filter)
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
