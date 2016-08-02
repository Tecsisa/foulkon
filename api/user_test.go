package api

import (
	"testing"

	"github.com/tecsisa/authorizr/database"
)

func TestAuthAPI_AddUser(t *testing.T) {
	testcases := map[string]struct {
		authUser   AuthenticatedUser
		externalID string
		path       string

		getUserByExternalIDMethodResult *User
		getGroupsByUserIDResult         []Group
		getAttachedPoliciesResult       []Policy

		expectedUser *User
		wantError    error

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
		"ErrorCaseInvalidExtID": {
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
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: externalId *%~#@|",
			},
		},
		"ErrorCaseInvalidPath": {
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
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: path /**%%/*123",
			},
		},
		"ErrorCaseNopath": {
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: path ",
			},
			externalID: "1234",
		},
		"ErrorCaseNoID": {
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: externalId ",
			},
			path: "/example/",
		},
		"ErrorCaseNoAuth": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			externalID: "1234",
			path:       "/example/",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code:    database.USER_NOT_FOUND,
				Message: "User not found",
			},
		},
		"ErrorCaseUnauthorizedResource": {
			authUser: AuthenticatedUser{
				Identifier: "000",
				Admin:      false,
			},
			getUserByExternalIDMethodResult: &User{
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
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 000 is not allowed to access to resource urn:iws:iam::user/test/asd/1234",
			},
		},
		"ErrorCaseUserAlreadyExists": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "1234",
			path:       "/example/",
			wantError: &Error{
				Code:    USER_ALREADY_EXIST,
				Message: "Unable to create user, user with externalId 1234 already exist",
			},
		},
		"ErrorCaseDBErr1": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			getUserByExternalIDMethodResult: &User{
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
		"ErrorCaseDBErr2": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			getUserByExternalIDMethodResult: &User{
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
				Code:    UNKNOWN_API_ERROR,
				Message: "User not found",
			},
		},
	}

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for x, testcase := range testcases {
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDMethodResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		testRepo.ArgsOut[AddUserMethod][0] = testcase.expectedUser
		testRepo.ArgsOut[AddUserMethod][1] = testcase.addUserMethodErr
		user, err := testAPI.AddUser(testcase.authUser, testcase.externalID, testcase.path)
		CheckApiResponse(t, x, testcase.wantError, err, testcase.expectedUser, user)
	}

}

func TestAuthAPI_GetUserByExternalID(t *testing.T) {
	testcases := map[string]struct {
		authUser   AuthenticatedUser
		externalID string

		getGroupsByUserIDMethodResult   []Group
		getAttachedPoliciesMethodResult []Policy
		getUserByExternalIDMethodResult *User

		expectedUser *User
		wantError    error

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
				Urn:        CreateUrn("", RESOURCE_USER, "/users/test/", "123"),
			},
			getUserByExternalIDMethodResult: &User{
				ID:         "543210",
				ExternalID: "123",
				Path:       "/users/test/",
			},
		},
		"ErrorCaseAuthUserNotExist": {
			authUser: AuthenticatedUser{
				Identifier: "notAdminUser",
				Admin:      false,
			},
			externalID: "123",
			expectedUser: &User{
				ID:         "543210",
				ExternalID: "123",
				Path:       "/users/test/",
				Urn:        CreateUrn("", RESOURCE_USER, "/users/test/", "123"),
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId notAdminUser is not allowed to access to resource urn:iws:iam::user/users/test/123",
			},
		},
		"ErrorCaseUnauthorizedResource": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			getUserByExternalIDMethodResult: &User{
				ID:         "000",
				ExternalID: "000",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "000"),
			},
			getGroupsByUserIDMethodResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesMethodResult: []Policy{
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
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam::user/path/000",
			},
		},
		"ErrorCaseUserNotExist": {
			authUser: AuthenticatedUser{
				Identifier: "notAdminUser",
				Admin:      true,
			},
			externalID: "111",
			wantError: &Error{
				Code:    USER_BY_EXTERNAL_ID_NOT_FOUND,
				Message: "User not found",
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code:    database.USER_NOT_FOUND,
				Message: "User not found",
			},
		},
		"ErrorCaseNoID": {
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: externalId ",
			},
		},
		"ErrorCaseInternalError": {
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
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDMethodResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesMethodResult
		user, err := testAPI.GetUserByExternalID(testcase.authUser, testcase.externalID)
		CheckApiResponse(t, x, testcase.wantError, err, testcase.expectedUser, user)
	}

}

func TestAuthAPI_ListUsers(t *testing.T) {
	testcases := map[string]struct {
		authUser   AuthenticatedUser
		pathPrefix string

		getUsersFilteredMethodResult    []User
		getGroupsByUserIDMethodResult   []Group
		getAttachedPoliciesMethodResult []Policy
		getUserByExternalIDMethodResult *User

		expectedResult []string
		wantError      error

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
					Urn:        CreateUrn("", RESOURCE_USER, "/example/test/", "123"),
				},
				{
					ID:         "321",
					ExternalID: "321",
					Path:       "/example/test2/",
					Urn:        CreateUrn("", RESOURCE_USER, "/example/test2/", "321"),
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
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: PathPrefix /^*$**~#!/",
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
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDMethodResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDMethodResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesMethodResult
		testRepo.ArgsOut[GetUsersFilteredMethod][0] = testcase.getUsersFilteredMethodResult
		testRepo.ArgsOut[GetUsersFilteredMethod][1] = testcase.GetUsersFilteredMethodErr
		users, err := testAPI.ListUsers(testcase.authUser, testcase.pathPrefix)
		CheckApiResponse(t, x, testcase.wantError, err, testcase.expectedResult, users)
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

		wantError error

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
				Urn:        CreateUrn("", RESOURCE_USER, "/example/", "1234"),
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
				Urn:        CreateUrn("", RESOURCE_USER, "/example/", "1234"),
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
				Urn:        CreateUrn("", RESOURCE_USER, "/example/", "1234"),
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: externalId *%~#@|",
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
				Urn:        CreateUrn("", RESOURCE_USER, "/example/", "1234"),
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: path /**%%/*123",
			},
		},
		"Nopath": {
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: path ",
			},
			externalID: "1234",
		},
		"NoID": {
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: externalId ",
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
				Urn:        CreateUrn("", RESOURCE_USER, "/example/", "1234"),
			},
			newPath: "/example/",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam::user/example/1234",
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
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam::user/path/000",
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
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam::user/path/000",
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
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam::user/newpath/000",
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
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam::user/newpath/000",
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
		CheckApiResponse(t, x, testcase.wantError, err, testcase.expectedUser, user)
	}

}

func TestAuthAPI_RemoveUser(t *testing.T) {
	testcases := map[string]struct {
		authUser   AuthenticatedUser
		externalID string

		GetUserByExternalIDMethodResult *User

		getGroupsByUserIDResult   []Group
		getAttachedPoliciesResult []Policy

		wantError error

		removeUserMethodErr error
	}{
		"OKCase": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "1234",
			GetUserByExternalIDMethodResult: &User{
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
			GetUserByExternalIDMethodResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/example/",
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: externalId *%~#@|",
			},
		},
		"NoID": {
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: externalId ",
			},
		},
		"NoAuth": {
			authUser: AuthenticatedUser{
				Identifier: "1234",
				Admin:      false,
			},
			externalID: "1234",
			GetUserByExternalIDMethodResult: &User{
				ID:         "543210",
				ExternalID: "12345",
				Path:       "/example/",
				Urn:        CreateUrn("", RESOURCE_USER, "/example/", "12345"),
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam::user/example/12345",
			},
		},
		"DeleteNotAllowedInPath": {
			authUser: AuthenticatedUser{
				Identifier: "1234",
				Admin:      false,
			},
			externalID: "1234",
			GetUserByExternalIDMethodResult: &User{
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
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam::user/path/1234",
			},
		},
		"DeleteUserNotAllowed": {
			authUser: AuthenticatedUser{
				Identifier: "1234",
				Admin:      false,
			},
			externalID: "1234",
			GetUserByExternalIDMethodResult: &User{
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
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam::user/path/1234",
			},
		},
		"DBErr": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "123456",
			GetUserByExternalIDMethodResult: &User{
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
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.GetUserByExternalIDMethodResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		testRepo.ArgsOut[RemoveUserMethod][0] = testcase.removeUserMethodErr
		err := testAPI.RemoveUser(testcase.authUser, testcase.externalID)
		CheckApiResponse(t, x, testcase.wantError, err, nil, nil)
	}
}

func TestAuthAPI_ListGroupsByUser(t *testing.T) {
	testcases := map[string]struct {
		authUser         AuthenticatedUser
		externalID       string
		wantError        error
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
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam::user/path/1234",
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
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
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
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam::user/path/1234",
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
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam::user/path/1234",
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
		CheckApiResponse(t, x, testcase.wantError, err, testcase.expectedResponse, groups)
	}

}
