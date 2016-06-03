package api

import (
	"reflect"
	"testing"

	"github.com/tecsisa/authorizr/database"
)

func TestGetUserByExternalId(t *testing.T) {
	testcases := map[string]struct {
		authUser   AuthenticatedUser
		externalID string

		getGroupsByUserIDResult   []Group
		getPoliciesAttachedResult []Policy
		getUserByExternalIDResult *User

		expectedUser *User
		wantError    *Error

		getUserByExternalIDMethodErr error
	}{
		"OKCase": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "123",
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
			externalID: "123",
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
			externalID: "000",
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
			externalID: "111",
			wantError: &Error{
				Code: USER_BY_EXTERNAL_ID_NOT_FOUND,
			},
			getUserByExternalIDMethodErr: &database.Error{
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
			externalID: "1234",
			getUserByExternalIDMethodErr: &database.Error{
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
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetPoliciesAttachedMethod][0] = testcase.getPoliciesAttachedResult
		user, err := testAPI.GetUserByExternalId(testcase.authUser, testcase.externalID)
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
		authUser   AuthenticatedUser
		externalID string
		path       string

		getUserByExternalIDResult *User
		getGroupsByUserIDResult   []Group
		getPoliciesAttachedResult []Policy

		expectedUser *User
		wantError    *Error

		addUserMethodErr             error
		getUserByExternalIDMethodErr error
	}{
		"OKCase": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "1234",
			path:       "/example/",
			expectedUser: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/example/",
			},
			getUserByExternalIDMethodErr: &database.Error{
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
			path:       "/example/",
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
			path: "/example/",
		},
		"NoAuth": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			externalID: "1234",
			path:       "/example/",
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code:    database.USER_NOT_FOUND,
				Message: "User not found",
			},
		},
		"ErrorUnauthorizedResource": {
			authUser: AuthenticatedUser{
				Identifier: "000",
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
								USER_ACTION_CREATE_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/test/"),
							},
						},
						Statement{
							Effect: "deny",
							Action: []string{
								USER_ACTION_CREATE_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/test/asd"),
							},
						},
					},
				},
			},
			externalID: "1234",
			path:       "/test/asd/",
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"InvalidExtID": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "*%~#@|",
			path:       "/example/",
			expectedUser: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/example/",
			},
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"InvalidPath": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "012",
			path:       "/**%%/*123",
			expectedUser: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/example/",
			},
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"DBErr1": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code:    database.USER_NOT_FOUND,
				Message: "User not found",
			},
			externalID: "12",
			path:       "/example/",
			addUserMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
		"DBErr2": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "User not found",
			},
			externalID: "12",
			path:       "/example/",
			addUserMethodErr: &database.Error{
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
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetPoliciesAttachedMethod][0] = testcase.getPoliciesAttachedResult
		testRepo.ArgsOut[AddUserMethod][0] = testcase.expectedUser
		testRepo.ArgsOut[AddUserMethod][1] = testcase.addUserMethodErr
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
		authUser   AuthenticatedUser
		externalID string
		newPath    string

		expectedUser *User

		getGroupsByUserIDResult   []Group
		getPoliciesAttachedResult []Policy

		wantError *Error

		updateUserMethodErr          error
		getUserByExternalIDMethodErr error
	}{
		"OKCase": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "1234",
			newPath:    "/example/",
			expectedUser: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/example/",
			},
		},
		"DBErr": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "123456",
			newPath:    "/example/",
			expectedUser: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/example/",
			},
			updateUserMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
		"InvalidExtID": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "*%~#@|",
			newPath:    "/example/",
			expectedUser: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/example/",
			},
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"InvalidPath": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "012",
			newPath:    "/**%%/*123",
			expectedUser: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/example/",
			},
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
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
			newPath: "/example/",
		},
		"NoAuth": {
			authUser: AuthenticatedUser{
				Identifier: "1234",
				Admin:      false,
			},
			externalID: "1234",
			expectedUser: &User{
				ID:         "543210",
				ExternalID: "12345",
				Path:       "/example/",
			},
			newPath: "/example/",
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"UpdateNotAllowed": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			newPath: "/newpath/",
			expectedUser: &User{
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
							Effect: "allow",
							Action: []string{
								USER_ACTION_UPDATE_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
						Statement{
							Effect: "deny",
							Action: []string{
								USER_ACTION_UPDATE_USER,
							},
							Resources: []string{
								CreateUrn("", RESOURCE_USER, "/path/", "000"),
							},
						},
					},
				},
			},
			externalID: "000",
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"ZeroPermissions": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			newPath: "/newpath/",
			expectedUser: &User{
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
					},
				},
			},
			externalID: "000",
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"GetNewPathNotAllowed": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			newPath: "/newpath/",
			expectedUser: &User{
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
							Effect: "allow",
							Action: []string{
								USER_ACTION_UPDATE_USER,
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
								CreateUrn("", RESOURCE_USER, "/newpath/", "000"),
							},
						},
					},
				},
			},
			externalID: "000",
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"NewPathNotAllowed": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			newPath: "/newpath/",
			expectedUser: &User{
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
							Effect: "allow",
							Action: []string{
								USER_ACTION_UPDATE_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
						Statement{
							Effect: "allow",
							Action: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/newpath/"),
							},
						},
						Statement{
							Effect: "deny",
							Action: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								CreateUrn("", RESOURCE_USER, "/newpath/", "000"),
							},
						},
					},
				},
			},
			externalID: "000",
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
	}

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for x, testcase := range testcases {
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.expectedUser
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetPoliciesAttachedMethod][0] = testcase.getPoliciesAttachedResult
		testRepo.ArgsOut[UpdateUserMethod][0] = testcase.expectedUser
		testRepo.ArgsOut[UpdateUserMethod][1] = testcase.updateUserMethodErr
		user, err := testAPI.UpdateUser(testcase.authUser, testcase.externalID, testcase.newPath)
		if testcase.wantError != nil {
			if errCode := err.(*Error).Code; errCode != testcase.wantError.Code {
				t.Fatalf("Test %v failed. Got error %v, expected %v", x, errCode, testcase.wantError.Code)
			}
		} else {
			if err != nil {
				t.Fatalf("Test %v failed: %v", x, err)
			} else {
				if testcase.expectedUser.ExternalID != user.ExternalID {
					t.Fatalf("Test %v failed. Received different users (wanted:%v / received:%v)", x, testcase.expectedUser.ExternalID, user.ExternalID)
				}
			}
		}
	}

}
