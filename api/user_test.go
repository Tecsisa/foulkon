package api

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/tecsisa/authorizr/database"
)

func TestAuthAPI_AddUser(t *testing.T) {
	testcases := map[string]struct {
		authUser   AuthenticatedUser
		externalID string
		path       string

		getUserByExternalIDResult *User
		getGroupsByUserIDResult   []Group
		getAttachedPoliciesResult []Policy

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
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Action: []string{
								USER_ACTION_CREATE_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/test/"),
							},
						},
						{
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
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		testRepo.ArgsOut[AddUserMethod][0] = testcase.expectedUser
		testRepo.ArgsOut[AddUserMethod][1] = testcase.addUserMethodErr
		user, err := testAPI.AddUser(testcase.authUser, testcase.externalID, testcase.path)
		if testcase.wantError != nil {
			if errCode := err.(*Error).Code; errCode != testcase.wantError.Code {
				t.Errorf("Test %v failed. Got error %v, expected %v", x, errCode, testcase.wantError.Code)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed: %v", x, err)
				continue
			} else {
				if diff := pretty.Compare(user, testcase.expectedUser); diff != "" {
					t.Errorf("Test %v failed. Received different responses (received/wanted) %v", x, diff)
					continue
				}
			}
		}
	}

}

func TestAuthAPI_GetUserByExternalID(t *testing.T) {
	testcases := map[string]struct {
		authUser   AuthenticatedUser
		externalID string

		getGroupsByUserIDResult   []Group
		getAttachedPoliciesResult []Policy
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
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Action: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
						{
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
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		user, err := testAPI.GetUserByExternalID(testcase.authUser, testcase.externalID)
		if testcase.wantError != nil {
			if errCode := err.(*Error).Code; errCode != testcase.wantError.Code {
				t.Errorf("Test %v failed. Got error %v, expected %v", x, errCode, testcase.wantError.Code)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed", x)
				continue
			} else {
				if testcase.expectedUser.ExternalID != user.ExternalID {
					t.Errorf("Test %v failed. Received different users (wanted:%v / received:%v)",
						x, testcase.expectedUser.ExternalID, user.ExternalID)
					continue
				}
			}
		}
	}

}

func TestAuthAPI_ListUsers(t *testing.T) {
	testcases := map[string]struct {
		authUser   AuthenticatedUser
		pathPrefix string

		getUsersFilteredMethodResult []User
		getGroupsByUserIDResult      []Group
		getAttachedPoliciesResult    []Policy
		getUserByExternalIDResult    *User

		expectedResult []string
		wantError      *Error

		GetUsersFilteredMethodErr    error
		getUserByExternalIDMethodErr error
	}{
		"OKCase": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			pathPrefix: "",
			getUsersFilteredMethodResult: []User{
				{
					ID:         "123",
					ExternalID: "123",
					Path:       "/example/test/",
				},
				{
					ID:         "321",
					ExternalID: "321",
					Path:       "/example/test2/",
				},
			},
			expectedResult: []string{"123", "321"},
		},
		"GetUserExtDBErr": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			pathPrefix: "/example/",
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
		"InvalidPath": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			pathPrefix: "/^*$**~#!/",
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"FilterUsersDBErr": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			pathPrefix: "/example/",
			GetUsersFilteredMethodErr: &database.Error{
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
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		testRepo.ArgsOut[GetUsersFilteredMethod][0] = testcase.getUsersFilteredMethodResult
		testRepo.ArgsOut[GetUsersFilteredMethod][1] = testcase.GetUsersFilteredMethodErr
		users, err := testAPI.ListUsers(testcase.authUser, testcase.pathPrefix)
		if testcase.wantError != nil {
			if errCode := err.(*Error).Code; errCode != testcase.wantError.Code {
				t.Errorf("Test %v failed. Got error %v, expected %v",
					x, errCode, testcase.wantError.Code)
				continue
			}
		} else {
			if diff := pretty.Compare(users, testcase.expectedResult); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", x, diff)
				continue
			}
		}
	}

}

func TestAuthAPI_UpdateUser(t *testing.T) {
	testcases := map[string]struct {
		authUser   AuthenticatedUser
		externalID string
		newPath    string

		expectedUser *User

		getGroupsByUserIDResult   []Group
		getAttachedPoliciesResult []Policy

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
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Action: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
						{
							Effect: "allow",
							Action: []string{
								USER_ACTION_UPDATE_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
						{
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
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
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
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Action: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
						{
							Effect: "allow",
							Action: []string{
								USER_ACTION_UPDATE_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
						{
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
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Action: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
						{
							Effect: "allow",
							Action: []string{
								USER_ACTION_UPDATE_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
						{
							Effect: "allow",
							Action: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/newpath/"),
							},
						},
						{
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
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		testRepo.ArgsOut[UpdateUserMethod][0] = testcase.expectedUser
		testRepo.ArgsOut[UpdateUserMethod][1] = testcase.updateUserMethodErr
		user, err := testAPI.UpdateUser(testcase.authUser, testcase.externalID, testcase.newPath)
		if testcase.wantError != nil {
			if errCode := err.(*Error).Code; errCode != testcase.wantError.Code {
				t.Errorf("Test %v failed. Got error %v, expected %v", x, errCode, testcase.wantError.Code)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed: %v", x, err)
				continue
			} else {
				if diff := pretty.Compare(user, testcase.expectedUser); diff != "" {
					t.Errorf("Test %v failed. Received different responses (received/wanted) %v", x, diff)
					continue
				}
			}
		}
	}

}

func TestAuthAPI_RemoveUser(t *testing.T) {
	testcases := map[string]struct {
		authUser   AuthenticatedUser
		externalID string

		expectedUser *User

		getGroupsByUserIDResult   []Group
		getAttachedPoliciesResult []Policy

		wantError *Error

		removeUserMethodErr error
	}{
		"OKCase": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "1234",
			expectedUser: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/example/",
			},
		},
		"InvalidExtID": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "*%~#@|",
			expectedUser: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/example/",
			},
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"NoID": {
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
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
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"DeleteNotAllowedInPath": {
			authUser: AuthenticatedUser{
				Identifier: "1234",
				Admin:      false,
			},
			externalID: "1234",
			expectedUser: &User{
				ID:         "1234",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
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
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"DeleteUserNotAllowed": {
			authUser: AuthenticatedUser{
				Identifier: "1234",
				Admin:      false,
			},
			externalID: "1234",
			expectedUser: &User{
				ID:         "1234",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Action: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
						{
							Effect: "allow",
							Action: []string{
								USER_ACTION_DELETE_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
						{
							Effect: "deny",
							Action: []string{
								USER_ACTION_DELETE_USER,
							},
							Resources: []string{
								CreateUrn("", RESOURCE_USER, "/path/", "1234"),
							},
						},
					},
				},
			},
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"DBErr": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "123456",
			expectedUser: &User{
				ID:         "543210",
				ExternalID: "123456",
				Path:       "/example/",
			},
			removeUserMethodErr: &database.Error{
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
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		testRepo.ArgsOut[RemoveUserMethod][0] = testcase.removeUserMethodErr
		err := testAPI.RemoveUser(testcase.authUser, testcase.externalID)
		if testcase.wantError != nil {
			if errCode := err.(*Error).Code; errCode != testcase.wantError.Code {
				t.Errorf("Test %v failed. Got error %v, expected %v",
					x, errCode, testcase.wantError.Code)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed: %v", x, err)
				continue
			}
		}
	}
}

func TestAuthAPI_ListGroupsByUser(t *testing.T) {
	testcases := map[string]struct {
		authUser         AuthenticatedUser
		externalID       string
		wantError        *Error
		expectedResponse []GroupIdentity

		getUserByExternalIDResult *User

		getGroupsByUserIDResult   []Group
		getAttachedPoliciesResult []Policy

		getGroupsByUserIDErr         error
		getUserByExternalIDMethodErr error
	}{
		"OKCase": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "1234",
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/example/",
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP1",
					Org:  "org1",
					Name: "groupUser1",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser1"),
				},
				{
					ID:   "GROUP2",
					Org:  "org2",
					Name: "groupUser2",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser1"),
				},
			},
			expectedResponse: []GroupIdentity{
				{
					Org:  "org1",
					Name: "groupUser1",
				},
				{
					Org:  "org2",
					Name: "groupUser2",
				},
			},
		},
		"AuthUserWithoutPermissions": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			externalID: "1234",
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"getGroupsDBErr": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "1234",
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/example/",
			},
			getGroupsByUserIDErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
		"GetUserExtIDDBErr": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "1234",
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP1",
					Name: "groupUser1",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser1"),
				},
				{
					ID:   "GROUP2",
					Name: "groupUser2",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser1"),
				},
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
		"DenyResourceErr": {
			authUser: AuthenticatedUser{
				Identifier: "1234",
				Admin:      false,
			},
			externalID: "1234",
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/1/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Action: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								CreateUrn("", RESOURCE_USER, "/path/", "1234"),
							},
						},
						{
							Effect: "allow",
							Action: []string{
								USER_ACTION_LIST_GROUPS_FOR_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
						{
							Effect: "deny",
							Action: []string{
								USER_ACTION_LIST_GROUPS_FOR_USER,
							},
							Resources: []string{
								CreateUrn("", RESOURCE_USER, "/path/", "1234"),
							},
						},
					},
				},
			},
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"UnauthorizedListGroupsErr": {
			authUser: AuthenticatedUser{
				Identifier: "1234",
				Admin:      false,
			},
			externalID: "12345",
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/1/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
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
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
	}

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for x, testcase := range testcases {
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][1] = testcase.getGroupsByUserIDErr
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		groups, err := testAPI.ListGroupsByUser(testcase.authUser, testcase.externalID)
		if testcase.wantError != nil {
			if errCode := err.(*Error).Code; errCode != testcase.wantError.Code {
				t.Errorf("Test %v failed. Got error %v, expected %v",
					x, errCode, testcase.wantError.Code)
				continue
			}
		} else {
			if diff := pretty.Compare(groups, testcase.expectedResponse); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", x, diff)
				continue
			}
		}
	}

}
