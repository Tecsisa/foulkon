package api

import (
	"testing"

	"reflect"

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
				Path:       "/users/test/",
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
				Path:       "/users/test/",
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

//func TestGetListUsers(t *testing.T) {
//	testcases := map[string]struct {
//		// User authenticated
//		authUser AuthenticatedUser
//		// Resource urn that user wants to access
//		pathPrefix string
//		// Error to compare when we expect an error
//		Users       []*User
//		expectedIDs []string
//		wantError   *Error
//		// GetUserByExternalID Method Out Arguments
//		err error
//	}{
//		"OKCase": {
//			authUser: AuthenticatedUser{
//				Identifier: "123456",
//				Admin:      true,
//			},
//			pathPrefix: "/",
//			Users: []*User{
//				{
//					ID:         "543210",
//					ExternalID: "123",
//					Path:       "urn:bla:bla:/tecsisa/test",
//				},
//				{
//					ID:         "543211",
//					ExternalID: "124",
//					Path:       "urn:bla:bla:/tecsisa/test2",
//				},
//			},
//			expectedIDs: []string{"123", "124"},
//		},
//		"NoIDpassed": {
//			wantError: &Error{
//				Code: INVALID_PARAMETER_ERROR,
//			},
//		},
//	}
//
//	testRepo := makeTestRepo()
//	testAPI := makeTestAPI(testRepo)
//
//	for x, testcase := range testcases {
//
//		testRepo.ArgsOut[GetListUsers][0] = testcase.expectedIDs
//		testRepo.ArgsOut[GetListUsers][1] = testcase.err
//		users, err := testAPI.GetListUsers(testcase.authUser, testcase.pathPrefix)
//		if testcase.wantError != nil {
//			if errCode := err.(*Error).Code; errCode != testcase.wantError.Code {
//				t.Fatalf("Test %v failed. Got error %v, expected %v", x, errCode, testcase.wantError.Code)
//			}
//		} else {
//			if err != nil {
//				t.Fatalf("Test %v failed: %v", x, err)
//			} else {
//				if !reflect.DeepEqual(testcase.expectedIDs, users) {
//					for _, x := range users {
//						println(x)
//					}
//					t.Fatalf("Test %v failed. Received different users", x)
//				}
//			}
//		}
//	}
//
//}

func TestAddUser(t *testing.T) {
	testcases := map[string]struct {
		// User authenticated
		authUser AuthenticatedUser
		// Resource urn that user wants to access
		externalID string
		path       string
		// Error to compare when we expect an error
		expectedUser *User
		wantError    *Error
		// GetUserByExternalID Method Out Arguments
		err    error
		errGet error
	}{
		"OKCase": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "1234",
			path:       "/tecsisa/",
			expectedUser: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/tecsisa/",
			},
			errGet: &database.Error{
				Code:    database.USER_NOT_FOUND,
				Message: "User not found",
			},
		},
		"AlreadyExists": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "1234",
			path:       "/tecsisa/",
			wantError: &Error{
				Code:    USER_ALREADY_EXIST,
				Message: "User already exists",
			},
		},
		"Nopath": {
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
			externalID: "1234",
		},
		"NoID": {
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
			path: "/tecsisa/",
		},
		"NoAuth": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			externalID: "1234",
			path:       "/tecsisa/",
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
			errGet: &database.Error{
				Code:    database.USER_NOT_FOUND,
				Message: "User not found",
			},
		},
	}

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for x, testcase := range testcases {

		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.errGet
		testRepo.ArgsOut[AddUserMethod][0] = testcase.expectedUser
		testRepo.ArgsOut[AddUserMethod][1] = testcase.err
		user, err := testAPI.AddUser(testcase.authUser, testcase.externalID, testcase.path)
		if testcase.wantError != nil {
			if errCode := err.(*Error).Code; errCode != testcase.wantError.Code {
				t.Fatalf("Test %v failed. Got error %v, expected %v", x, errCode, testcase.wantError.Code)
			}
		} else {
			if err != nil {
				t.Fatalf("Test %v failed: %v", x, err)
			} else {
				if !reflect.DeepEqual(testcase.expectedUser, user) {
					t.Fatalf("Test %v failed. Received different users", x)
				}
			}
		}
	}

}
