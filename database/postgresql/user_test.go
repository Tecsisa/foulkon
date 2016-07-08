package postgresql

import (
	"fmt"
	"github.com/kylelemons/godebug/pretty"
	"github.com/tecsisa/authorizr/api"
	"github.com/tecsisa/authorizr/database"
	"testing"
	"time"
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
			insertUser(test.previousUser.ID, test.previousUser.ExternalID, test.previousUser.Path,
				test.previousUser.CreateAt.Unix(), test.previousUser.Urn)
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
			userNumber, err := getUsersCountFiltered(test.userToCreate.ID, test.userToCreate.ExternalID, test.userToCreate.Path,
				test.userToCreate.CreateAt.UnixNano(), test.userToCreate.Urn)
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
				Message: fmt.Sprint("User with ExternalID NotExist not found"),
			},
		},
	}

	for n, test := range testcases {
		// Clean user database
		cleanUserTable()

		// Insert previous data
		if test.previousUser != nil {
			insertUser(test.previousUser.ID, test.previousUser.ExternalID, test.previousUser.Path,
				test.previousUser.CreateAt.UnixNano(), test.previousUser.Urn)
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
			// Check database
			userNumber, err := getUsersCountFiltered("", test.externalID, "", 0, "")
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
				Message: fmt.Sprint("User with id NotExist not found"),
			},
		},
	}

	for n, test := range testcases {
		// Clean user database
		cleanUserTable()

		// Insert previous data
		if test.previousUser != nil {
			insertUser(test.previousUser.ID, test.previousUser.ExternalID, test.previousUser.Path,
				test.previousUser.CreateAt.UnixNano(), test.previousUser.Urn)
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
			// Check database
			userNumber, err := getUsersCountFiltered("", test.userID, "", 0, "")
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
