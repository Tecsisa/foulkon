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
		expectedUser              *User
		wantError                 *Error
		getGroupsByUserIDResult   []Group
		getPoliciesAttachedResult []Policy
		getUserByExternalIDResult *User

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
		"ErrorAuthUserNotExist": {
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
		"ErrorUnauthorizedResource": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			getUserByExternalIDResult: &User{
				ID:         "000",
				ExternalID: "000",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "000"),
			},
			getGroupsByUserIDResult: []Group{
				Group{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getPoliciesAttachedResult: []Policy{
				Policy{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
						Statement{
							Effect: "deny",
							Action: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								CreateUrn("", RESOURCE_USER, "/path/", "000"),
							},
						},
					},
				},
			},
			id: "000",
			expectedUser: &User{
				ID:         "000",
				ExternalID: "000",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "000"),
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
		"InternalError": {
			id: "1234",
			err: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
	}

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for x, testcase := range testcases {
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.expectedUser
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.err
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetPoliciesAttachedMethod][0] = testcase.getPoliciesAttachedResult
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
					t.Fatalf("Test %v failed. Received different users (wanted:%v / received:%v)", x, testcase.expectedUser.ExternalID, user.ExternalID)
				}
			}
		}
	}

}

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

func TestUpdateUser(t *testing.T) {
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
				if !reflect.DeepEqual(testcase.expectedUser.ExternalID, user.ExternalID) {
					t.Fatalf("Test %v failed. Received different users (wanted:%v / received:%v)", x)
				}
			}
		}
	}

}
