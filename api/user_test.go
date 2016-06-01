package api

import (
	"testing"

	"github.com/tecsisa/authorizr/database"
)

func TestGetUserByExternalId(t *testing.T) {
	testcases := map[string]struct {
		// User authenticated
		authUser AuthenticatedUser
		// Resource urn that user wants to access
		id string
		// Error to compare when we expect an error
		expectedUser *User
		wantError    *Error
		// GetUserByExternalID Method Out Arguments
		err error
	}{
		"OKCase": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			id: "123",
			expectedUser: &User{
				ID:         "543210",
				ExternalID: "123",
				Path:       "urn:bla:bla:/users/test",
			},
		},
		"ErrorUserUnauthorized": {
			authUser: AuthenticatedUser{
				Identifier: "notAdminUser",
				Admin:      false,
			},
			id: "123",
			expectedUser: &User{
				ID:         "543210",
				ExternalID: "123",
				Path:       "urn:bla:bla:/users/test",
			},
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"ErrorUserNotExist": {
			authUser: AuthenticatedUser{
				Identifier: "notAdminUser",
				Admin:      true,
			},
			id: "111",
			wantError: &Error{
				Code: USER_BY_EXTERNAL_ID_NOT_FOUND,
			},
			err: &database.Error{
				Code:    database.USER_NOT_FOUND,
				Message: "Error",
			},
		},
		"NoIDpassed": {
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
	}

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for x, testcase := range testcases {
		testRepo.ArgsIn[GetUserByExternalIDMethod] = make([]interface{}, 1)
		testRepo.ArgsOut[GetUserByExternalIDMethod] = make([]interface{}, 2)
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.expectedUser
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.err
		user, err := testAPI.GetUserByExternalId(testcase.authUser, testcase.id)
		if testcase.wantError != nil {
			if errCode := err.(*Error).Code; errCode != testcase.wantError.Code {
				t.Fatalf("Test %v failed. Got error %v, expected %v", x, errCode, testcase.wantError.Code)
			}
		} else {
			if err != nil {
				t.Fatalf("Test %v failed", x)
			} else {
				if testcase.expectedUser.ExternalID != user.ExternalID {
					t.Fatalf("Test %v failed. Got user %v, expected %v", x, testcase.expectedUser.ExternalID, testcase.id)
				}
			}
		}
	}

}
