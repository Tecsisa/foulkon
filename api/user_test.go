package api

import (
	"testing"

	"github.com/tecsisa/foulkon/database"
)

func TestAuthAPI_AddUser(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		requestInfo RequestInfo
		externalID  string
		path        string
		// Expected result
		expectedUser *User
		wantError    error
		// Manager Results
		getUserByExternalIDMethodResult      *User
		getUserByExternalIDMethodSpecialFunc func(string) (*User, error)
		getGroupsByUserIDResult              []Group
		getAttachedPoliciesResult            []Policy
		// API Errors
		addUserMethodErr             error
		getUserByExternalIDMethodErr error
	}{
		"OKCaseAdmin": {
			requestInfo: RequestInfo{
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
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			externalID: "1234",
			path:       "/example/",
			expectedUser: &User{
				ID:         "USER_ID",
				ExternalID: "1234",
				Path:       "/example/",
				Urn:        CreateUrn("", RESOURCE_USER, "/example/", "1234"),
			},
			getUserByExternalIDMethodSpecialFunc: func(id string) (*User, error) {
				if id == "123456" {
					return &User{
						ID:         "000",
						ExternalID: "000",
						Path:       "/path/",
						Urn:        CreateUrn("", RESOURCE_USER, "/path/", "000"),
					}, nil
				} else {
					return nil, &database.Error{
						Code:    database.USER_NOT_FOUND,
						Message: "User not found",
					}
				}
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
							Actions: []string{
								USER_ACTION_CREATE_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/example/"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseInvalidExtID": {
			externalID: "*%~#@|",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: externalId *%~#@|",
			},
		},
		"ErrorCaseNoID": {
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: externalId ",
			},
		},
		"ErrorCaseInvalidPath": {
			externalID: "012",
			path:       "/**%%/*123",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: path /**%%/*123",
			},
		},
		"ErrorCaseNopath": {
			externalID: "1234",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: path ",
			},
		},
		"ErrorCaseNoAuth": {
			requestInfo: RequestInfo{
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
		"ErrorCaseAddUserNotAllowed": {
			requestInfo: RequestInfo{
				Identifier: "000",
				Admin:      false,
			},
			externalID: "1234",
			path:       "/test/asd/",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 000 is not allowed to access to resource urn:iws:iam::user/test/asd/1234",
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
							Actions: []string{
								USER_ACTION_CREATE_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/test/"),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								USER_ACTION_CREATE_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/test/asd"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "000",
				Admin:      false,
			},
			externalID: "1234",
			path:       "/test/asd/",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 000 is not allowed to access to resource urn:iws:iam::user/test/asd/1234",
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
					ID:         "POLICY-USER-ID",
					Name:       "policyUser",
					Path:       "/path/",
					Urn:        CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{},
				},
			},
		},
		"ErrorCaseAddUserDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "12",
			path:       "/example/",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
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
			addUserMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
		"ErrorCaseGetUserDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "12",
			path:       "/example/",
			wantError: &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: "User not found",
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
		},
		"ErrorCaseUserAlreadyExist": {
			requestInfo: RequestInfo{
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
	}

	for x, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDMethodResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.SpecialFuncs[GetUserByExternalIDMethod] = testcase.getUserByExternalIDMethodSpecialFunc
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		testRepo.ArgsOut[AddUserMethod][0] = testcase.expectedUser
		testRepo.ArgsOut[AddUserMethod][1] = testcase.addUserMethodErr
		user, err := testAPI.AddUser(testcase.requestInfo, testcase.externalID, testcase.path)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.expectedUser, user)
	}

}

func TestAuthAPI_GetUserByExternalID(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		requestInfo RequestInfo
		externalID  string
		// Expected result
		expectedUser *User
		wantError    error
		// Manager Results
		getGroupsByUserIDMethodResult   []Group
		getAttachedPoliciesMethodResult []Policy
		getUserByExternalIDMethodResult *User
		// API Errors
		getUserByExternalIDMethodErr error
	}{
		"OKCaseAdmin": {
			requestInfo: RequestInfo{
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
				Urn:        CreateUrn("", RESOURCE_USER, "/users/test/", "123"),
			},
		},
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			externalID: "000",
			expectedUser: &User{
				ID:         "000",
				ExternalID: "000",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "000"),
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
							Actions: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseInvalidExtID": {
			externalID: "*%~#@|",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: externalId *%~#@|",
			},
		},
		"ErrorCaseNoID": {
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: externalId ",
			},
		},
		"ErrorCaseUserNotFound": {
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
		"ErrorCaseGetUserExtIDDBErr": {
			externalID: "1234",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
		"ErrorCaseNoAuth": {
			requestInfo: RequestInfo{
				Identifier: "notAdminUser",
				Admin:      false,
			},
			externalID: "123",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId notAdminUser is not allowed to access to resource urn:iws:iam::user/users/test/123",
			},
			getUserByExternalIDMethodResult: &User{
				ID:         "543210",
				ExternalID: "123",
				Path:       "/users/test/",
				Urn:        CreateUrn("", RESOURCE_USER, "/users/test/", "123"),
			},
		},
		"ErrorCaseGetUserNotAllowed": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			externalID: "000",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam::user/path/000",
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
							Actions: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								CreateUrn("", RESOURCE_USER, "/path/", "000"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			externalID: "000",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam::user/path/000",
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
					ID:         "POLICY-USER-ID",
					Name:       "policyUser",
					Path:       "/path/",
					Urn:        CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{},
				},
			},
		},
	}

	for x, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDMethodResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDMethodResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesMethodResult
		user, err := testAPI.GetUserByExternalID(testcase.requestInfo, testcase.externalID)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.expectedUser, user)
	}

}

func TestAuthAPI_ListUsers(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		requestInfo RequestInfo
		pathPrefix  string
		// Expected result
		expectedResult []string
		wantError      error
		// Manager Results
		getUsersFilteredMethodResult    []User
		getGroupsByUserIDMethodResult   []Group
		getAttachedPoliciesMethodResult []Policy
		getUserByExternalIDMethodResult *User
		// API Errors
		GetUsersFilteredMethodErr    error
		getUserByExternalIDMethodErr error
	}{
		"OKCaseAdmin": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			pathPrefix:     "",
			expectedResult: []string{"123", "321"},
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
		},
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			pathPrefix:     "",
			expectedResult: []string{"123", "321"},
			getUserByExternalIDMethodResult: &User{
				ID:         "000",
				ExternalID: "000",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "000"),
			},
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
							Actions: []string{
								USER_ACTION_LIST_USERS,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, ""),
							},
						},
					},
				},
			},
		},
		"OKCaseNoResourcesAllowed": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			pathPrefix:     "",
			expectedResult: []string{},
			getUserByExternalIDMethodResult: &User{
				ID:         "000",
				ExternalID: "000",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "000"),
			},
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
							Actions: []string{
								USER_ACTION_LIST_USERS,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, ""),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								USER_ACTION_LIST_USERS,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/example/"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseInvalidPath": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			pathPrefix: "/^*$**~#!/",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: PathPrefix /^*$**~#!/",
			},
		},
		"ErrorCaseNoAuth": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			pathPrefix: "/example/",
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
		},
		"ErrorCaseFilterUsersDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			pathPrefix: "/example/",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			GetUsersFilteredMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			pathPrefix: "",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam::user/*",
			},
			getUserByExternalIDMethodResult: &User{
				ID:         "000",
				ExternalID: "000",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "000"),
			},
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
					ID:         "POLICY-USER-ID",
					Name:       "policyUser",
					Path:       "/path/",
					Urn:        CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{},
				},
			},
		},
		"ErrorCaseListNotAllowed": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			pathPrefix: "",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam::user/*",
			},
			getUserByExternalIDMethodResult: &User{
				ID:         "000",
				ExternalID: "000",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "000"),
			},
			getUsersFilteredMethodResult: []User{
				{
					ID:         "123",
					ExternalID: "123",
					Path:       "/example/test/",
					Urn:        CreateUrn("", RESOURCE_USER, "/example/test/", "123"),
				},
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
							Actions: []string{
								USER_ACTION_LIST_USERS,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/example/test/"),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								USER_ACTION_LIST_USERS,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/example/test/"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseGetUserDbErrInAuth": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			pathPrefix: "/example/",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
	}

	for x, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDMethodResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDMethodResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesMethodResult
		testRepo.ArgsOut[GetUsersFilteredMethod][0] = testcase.getUsersFilteredMethodResult
		testRepo.ArgsOut[GetUsersFilteredMethod][1] = testcase.GetUsersFilteredMethodErr
		users, err := testAPI.ListUsers(testcase.requestInfo, testcase.pathPrefix)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.expectedResult, users)
	}

}

func TestAuthAPI_UpdateUser(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		requestInfo RequestInfo
		externalID  string
		newPath     string
		// Expected result
		expectedUser *User
		wantError    error
		// Manager Results
		getUserByExternalIDMethodResult *User
		getGroupsByUserIDMethodResult   []Group
		getAttachedPoliciesMethodResult []Policy
		// API Errors
		updateUserMethodErr          error
		getUserByExternalIDMethodErr error
	}{
		"OKCaseAdmin": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "1234",
			newPath:    "/example2/",
			expectedUser: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/example2/",
				Urn:        CreateUrn("", RESOURCE_USER, "/example/", "1234"),
			},
			getUserByExternalIDMethodResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/example/",
				Urn:        CreateUrn("", RESOURCE_USER, "/example/", "1234"),
			},
		},
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			externalID: "000",
			newPath:    "/newpath/",
			expectedUser: &User{
				ID:         "000",
				ExternalID: "000",
				Path:       "/newpath/",
				Urn:        CreateUrn("", RESOURCE_USER, "/newpath/", "000"),
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
							Actions: []string{
								USER_ACTION_GET_USER,
								USER_ACTION_UPDATE_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, ""),
							},
						},
					},
				},
			},
		},
		"ErrorCaseInvalidPath": {
			newPath: "/**%%/*123",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: path /**%%/*123",
			},
		},
		"ErrorCaseNoPath": {
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: path ",
			},
		},
		"ErrorCaseInvalidExtID": {
			externalID: "*%~#@|",
			newPath:    "/example/",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: externalId *%~#@|",
			},
		},
		"ErrorCaseNoID": {
			newPath: "/example/",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: externalId ",
			},
		},
		"ErrorCaseNoAuth": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			externalID: "1234",
			newPath:    "/example/",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam::user/example/12345",
			},
			getUserByExternalIDMethodResult: &User{
				ID:         "543210",
				ExternalID: "12345",
				Path:       "/example/",
				Urn:        CreateUrn("", RESOURCE_USER, "/example/", "12345"),
			},
		},
		"ErrorCaseUserNotFound": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			externalID: "1234",
			newPath:    "/example/",
			wantError: &Error{
				Code: USER_BY_EXTERNAL_ID_NOT_FOUND,
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrorCaseGetUserExtIDDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "1234",
			newPath:    "/example/",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
		"ErrorCaseUpdateNotAllowed": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			externalID: "000",
			newPath:    "/newpath/",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam::user/path/000",
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
							Actions: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			externalID: "000",
			newPath:    "/newpath/",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam::user/path/000",
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
							Actions: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
						{
							Effect: "allow",
							Actions: []string{
								USER_ACTION_UPDATE_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								USER_ACTION_UPDATE_USER,
							},
							Resources: []string{
								CreateUrn("", RESOURCE_USER, "/path/", "000"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseGetNewPathNotAllowed": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			externalID: "000",
			newPath:    "/newpath/",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam::user/newpath/000",
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
							Actions: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
						{
							Effect: "allow",
							Actions: []string{
								USER_ACTION_UPDATE_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								CreateUrn("", RESOURCE_USER, "/newpath/", "000"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseNewPathNotAllowed": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			externalID: "000",
			newPath:    "/newpath/",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam::user/newpath/000",
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
							Actions: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
						{
							Effect: "allow",
							Actions: []string{
								USER_ACTION_UPDATE_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
						{
							Effect: "allow",
							Actions: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/newpath/"),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								CreateUrn("", RESOURCE_USER, "/newpath/", "000"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseUpdateUserDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "123456",
			newPath:    "/example/",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getUserByExternalIDMethodResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/example/",
				Urn:        CreateUrn("", RESOURCE_USER, "/example/", "1234"),
			},
			updateUserMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
	}

	for x, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDMethodResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDMethodResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesMethodResult
		testRepo.ArgsOut[UpdateUserMethod][0] = testcase.expectedUser
		testRepo.ArgsOut[UpdateUserMethod][1] = testcase.updateUserMethodErr
		user, err := testAPI.UpdateUser(testcase.requestInfo, testcase.externalID, testcase.newPath)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.expectedUser, user)
	}

}

func TestAuthAPI_RemoveUser(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		requestInfo RequestInfo
		externalID  string
		// Expected result
		wantError error
		// Manager Results
		getUserByExternalIDMethodResult *User
		getGroupsByUserIDResult         []Group
		getAttachedPoliciesResult       []Policy
		// API Errors
		getUserByExternalIDMethodErr error
		removeUserMethodErr          error
	}{
		"OKCaseAdmin": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "1234",
			getUserByExternalIDMethodResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/example/",
			},
		},
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			externalID: "1234",
			getUserByExternalIDMethodResult: &User{
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
							Actions: []string{
								USER_ACTION_GET_USER,
								USER_ACTION_DELETE_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseInvalidExtID": {
			externalID: "*%~#@|",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: externalId *%~#@|",
			},
		},
		"ErrorCaseNoID": {
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: externalId ",
			},
		},
		"ErrorCaseNoAuth": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			externalID: "1234",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam::user/example/12345",
			},
			getUserByExternalIDMethodResult: &User{
				ID:         "543210",
				ExternalID: "12345",
				Path:       "/example/",
				Urn:        CreateUrn("", RESOURCE_USER, "/example/", "12345"),
			},
		},
		"ErrorCaseUserNotFound": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			externalID: "1234",
			wantError: &Error{
				Code: USER_BY_EXTERNAL_ID_NOT_FOUND,
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrorCaseGetUserExtIDDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "1234",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
		"ErrorCaseRemoveNotAllowedInPath": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			externalID: "1234",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam::user/path/1234",
			},
			getUserByExternalIDMethodResult: &User{
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
							Actions: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			externalID: "1234",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam::user/path/1234",
			},
			getUserByExternalIDMethodResult: &User{
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
							Actions: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
						{
							Effect: "allow",
							Actions: []string{
								USER_ACTION_DELETE_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								USER_ACTION_DELETE_USER,
							},
							Resources: []string{
								CreateUrn("", RESOURCE_USER, "/path/", "1234"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseRemoveUserDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "123456",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getUserByExternalIDMethodResult: &User{
				ID:         "543210",
				ExternalID: "123456",
				Path:       "/example/",
			},
			removeUserMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
	}

	for x, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDMethodResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		testRepo.ArgsOut[RemoveUserMethod][0] = testcase.removeUserMethodErr
		err := testAPI.RemoveUser(testcase.requestInfo, testcase.externalID)
		checkMethodResponse(t, x, testcase.wantError, err, nil, nil)
	}
}

func TestAuthAPI_ListGroupsByUser(t *testing.T) {
	testcases := map[string]struct {
		// API Method args
		requestInfo RequestInfo
		externalID  string
		wantError   error
		// Expected result
		expectedResponse []GroupIdentity
		// Manager Results
		getUserByExternalIDMethodResult *User
		getGroupsByUserIDMethodResult   []Group
		getAttachedPoliciesMethodResult []Policy
		// Manager Errors
		getGroupsByUserIDMethodErr   error
		getUserByExternalIDMethodErr error
	}{
		"OKCaseAdmin": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "1234",
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
			getUserByExternalIDMethodResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/example/",
			},
			getGroupsByUserIDMethodResult: []Group{
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
		},
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			externalID: "1234",
			expectedResponse: []GroupIdentity{
				{
					Org:  "org1",
					Name: "group1",
				},
				{
					Org:  "org2",
					Name: "group2",
				},
			},
			getUserByExternalIDMethodResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDMethodResult: []Group{
				{
					ID:   "GROUP-USER-ID-1",
					Name: "group1",
					Path: "/path/1/",
					Org:  "org1",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/1/", "group1"),
				},
				{
					ID:   "GROUP-USER-ID-2",
					Name: "group2",
					Path: "/path/2/",
					Org:  "org2",
					Urn:  CreateUrn("org2", RESOURCE_GROUP, "/path/2/", "group2"),
				},
			},
			getAttachedPoliciesMethodResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								CreateUrn("", RESOURCE_USER, "/path/", "1234"),
							},
						},
						{
							Effect: "allow",
							Actions: []string{
								USER_ACTION_LIST_GROUPS_FOR_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseNoAuth": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			externalID: "1234",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam::user/example/12345",
			},
			getUserByExternalIDMethodResult: &User{
				ID:         "543210",
				ExternalID: "12345",
				Path:       "/example/",
				Urn:        CreateUrn("", RESOURCE_USER, "/example/", "12345"),
			},
		},
		"ErrorCaseUserNotFound": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			externalID: "1234",
			wantError: &Error{
				Code: USER_BY_EXTERNAL_ID_NOT_FOUND,
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrorCaseGetUserExtIDDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "1234",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
		"ErrorCaseGetGroupsDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			externalID: "1234",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getUserByExternalIDMethodResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/example/",
			},
			getGroupsByUserIDMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			externalID: "1234",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam::user/path/1234",
			},
			getUserByExternalIDMethodResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDMethodResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/1/",
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
							Actions: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								CreateUrn("", RESOURCE_USER, "/path/", "1234"),
							},
						},
						{
							Effect: "allow",
							Actions: []string{
								USER_ACTION_LIST_GROUPS_FOR_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								USER_ACTION_LIST_GROUPS_FOR_USER,
							},
							Resources: []string{
								CreateUrn("", RESOURCE_USER, "/path/", "1234"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseUnauthorizedListGroups": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			externalID: "12345",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam::user/path/1234",
			},
			getUserByExternalIDMethodResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDMethodResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/1/",
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
							Actions: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, "/path/"),
							},
						},
					},
				},
			},
		},
	}

	for x, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDMethodResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDMethodResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][1] = testcase.getGroupsByUserIDMethodErr
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesMethodResult
		groups, err := testAPI.ListGroupsByUser(testcase.requestInfo, testcase.externalID)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.expectedResponse, groups)
	}

}
