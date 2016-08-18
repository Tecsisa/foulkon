package api

import (
	"testing"

	"github.com/tecsisa/foulkon/database"
)

func TestAuthAPI_AddGroup(t *testing.T) {
	testcases := map[string]struct {
		// API Method args
		requestInfo RequestInfo
		name        string
		org         string
		path        string
		// Expected results
		expectedGroup *Group
		wantError     error
		// Manager Results
		getUserByExternalIDResult *User
		getGroupsByUserIDResult   []Group
		getAttachedPoliciesResult []Policy
		getGroupByName            *Group
		addMemberMethodResult     *Group
		// Manager Errors
		getGroupByNameMethodErr      error
		getUserByExternalIDMethodErr error
		addGroupMethodErr            error
	}{
		"OKCaseAdmin": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			name: "group1",
			org:  "org1",
			path: "/example/",
			expectedGroup: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/example/",
			},
			getGroupByNameMethodErr: &database.Error{
				Code: database.GROUP_NOT_FOUND,
			},
		},
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "group1",
			org:  "org1",
			path: "/example/",
			expectedGroup: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/example/",
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
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
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_CREATE_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, "/example/"),
							},
						},
					},
				},
			},
			getGroupByNameMethodErr: &database.Error{
				Code: database.GROUP_NOT_FOUND,
			},
		},
		"ErrorCaseInvalidName": {
			name: "*%~#@|",
			org:  "org1",
			path: "/example/",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name *%~#@|",
			},
		},
		"ErrorCaseInvalidOrg": {
			name: "n1",
			org:  "*%~#@|",
			path: "/example/",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org *%~#@|",
			},
		},
		"ErrorCaseInvalidPath": {
			name: "group1",
			org:  "org1",
			path: "/**%%/*123",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: path /**%%/*123",
			},
		},
		"ErrorCaseGroupAlreadyExists": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			name: "group1",
			org:  "org1",
			path: "/example/",
			wantError: &Error{
				Code:    GROUP_ALREADY_EXIST,
				Message: "Unable to create group, group with org org1 and name group1 already exists",
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
			},
			addMemberMethodResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/example/",
			},
		},
		"ErrorCaseNotAuthenticatedUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "group1",
			org:  "org1",
			path: "/example/",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrorCaseUnauthorizedResource": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "group1",
			org:  "org1",
			path: "/test/asd/",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/test/asd/group1",
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
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
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_CREATE_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, "/test/"),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								GROUP_ACTION_CREATE_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, "/test/asd"),
							},
						},
					},
				},
			},
			getGroupByNameMethodErr: &database.Error{
				Code: database.GROUP_NOT_FOUND,
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "group1",
			org:  "org1",
			path: "/test/asd/",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/test/asd/group1",
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
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
					Urn:        CreateUrn("example", RESOURCE_GROUP, "/path/", "policyUser"),
					Statements: &[]Statement{},
				},
			},
			getGroupByNameMethodErr: &database.Error{
				Code: database.GROUP_NOT_FOUND,
			},
		},
		"ErrorCaseAddGroupDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			name: "group1",
			org:  "org1",
			path: "/example/",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getGroupByNameMethodErr: &database.Error{
				Code: database.GROUP_NOT_FOUND,
			},
			addGroupMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
		"ErrorCaseGetGroupDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			name: "group1",
			org:  "org1",
			path: "/example/",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getGroupByNameMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
	}

	for x, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetGroupByNameMethod][0] = testcase.getGroupByName
		testRepo.ArgsOut[GetGroupByNameMethod][1] = testcase.getGroupByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		testRepo.ArgsOut[AddGroupMethod][0] = testcase.expectedGroup
		testRepo.ArgsOut[AddGroupMethod][1] = testcase.addGroupMethodErr

		group, err := testAPI.AddGroup(testcase.requestInfo, testcase.org, testcase.name, testcase.path)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.expectedGroup, group)
	}
}

func TestAuthAPI_GetGroupByName(t *testing.T) {
	testcases := map[string]struct {
		// API Method args
		requestInfo RequestInfo
		name        string
		org         string
		path        string
		// Expected result
		expectedGroup *Group
		wantError     error
		// Manager Results
		getUserByExternalIDResult  *User
		getGroupsByUserIDResult    []Group
		getAttachedPoliciesResult  []Policy
		getGroupByNameMethodResult *Group
		// Manager Errors
		getUserByExternalIDMethodErr error
		getGroupByNameMethodErr      error
	}{
		"OKCaseAdmin": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			name: "group1",
			org:  "org1",
			path: "/example/",
			expectedGroup: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/example/",
			},
			getGroupByNameMethodResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/example/",
			},
		},
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "group1",
			org:  "org1",
			path: "/example/",
			expectedGroup: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/test/asd/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/test/asd/", "group1"),
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
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
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, "/test/"),
							},
						},
					},
				},
			},
			getGroupByNameMethodResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/test/asd/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/test/asd/", "group1"),
			},
		},
		"ErrorCaseInvalidName": {
			name: "*%~#@|",
			org:  "org1",
			path: "/example/",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name *%~#@|",
			},
		},
		"ErrorCaseInvalidOrg": {
			name: "n1",
			org:  "*%~#@|",
			path: "/example/",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org *%~#@|",
			},
		},
		"ErrorCaseGroupNotFound": {
			name: "group1",
			org:  "org1",
			path: "/example/",
			wantError: &Error{
				Code: GROUP_BY_ORG_AND_NAME_NOT_FOUND,
			},
			getGroupByNameMethodErr: &database.Error{
				Code: database.GROUP_NOT_FOUND,
			},
		},
		"ErrorCaseGetGroupDBErr": {
			name: "group1",
			org:  "org1",
			path: "/example/",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getGroupByNameMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
		"ErrorCaseNotAuthenticatedUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "group1",
			org:  "org1",
			path: "/example/",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
			getGroupByNameMethodResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/example/",
			},
		},
		"ErrorCaseUnauthorizedResource": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "group1",
			org:  "org1",
			path: "/test/asd/",
			expectedGroup: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/test/asd/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/test/asd/", "group1"),
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/test/asd/group1",
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
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
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, "/test/"),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, "/test/asd"),
							},
						},
					},
				},
			},
			getGroupByNameMethodResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/test/asd/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/test/asd/", "group1"),
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "group1",
			org:  "org1",
			path: "/test/asd/",
			expectedGroup: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/test/asd/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/test/asd/", "group1"),
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/test/asd/group1",
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
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
					Urn:        CreateUrn("example", RESOURCE_GROUP, "/path/", "policyUser"),
					Statements: &[]Statement{},
				},
			},
			getGroupByNameMethodResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/test/asd/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/test/asd/", "group1"),
			},
		},
	}

	for x, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetGroupByNameMethod][0] = testcase.getGroupByNameMethodResult
		testRepo.ArgsOut[GetGroupByNameMethod][1] = testcase.getGroupByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult

		group, err := testAPI.GetGroupByName(testcase.requestInfo, testcase.org, testcase.name)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.expectedGroup, group)
	}
}

func TestAuthAPI_ListGroups(t *testing.T) {
	testcases := map[string]struct {
		// API Method args
		requestInfo RequestInfo
		org         string
		pathPrefix  string
		// Expected result
		expectedGroups []GroupIdentity
		wantError      error
		// Manager Results
		getGroupsFilteredMethodResult []Group
		getGroupsByUserIDResult       []Group
		getAttachedPoliciesResult     []Policy
		getUserByExternalIDResult     *User
		// Manager Errors
		getUserByExternalIDMethodErr error
		getGroupsFilteredMethodErr   error
	}{
		"OKCaseAdmin": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "org1",
			pathPrefix: "/",
			expectedGroups: []GroupIdentity{
				{
					Org:  "org1",
					Name: "group1",
				},
			},
			getGroupsFilteredMethodResult: []Group{
				{
					Name: "group1",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "group1"),
				},
			},
		},
		"OKCaseAdminNoGroup": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			pathPrefix: "/",
			expectedGroups: []GroupIdentity{
				{
					Org:  "org1",
					Name: "group1",
				},
				{
					Org:  "org2",
					Name: "group2",
				},
			},
			getGroupsFilteredMethodResult: []Group{
				{
					Name: "group1",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "group1"),
				},
				{
					Name: "group2",
					Org:  "org2",
					Path: "/path2/",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path2/", "group2"),
				},
			},
		},
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org: "org1",
			expectedGroups: []GroupIdentity{
				{
					Org:  "org1",
					Name: "group1",
				},
			},
			getGroupsFilteredMethodResult: []Group{
				{
					Name: "group1",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "group1"),
				},
				{
					Name: "group2",
					Org:  "org2",
					Path: "/path2/",
					Urn:  CreateUrn("org2", RESOURCE_GROUP, "/path2/", "group2"),
				},
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/1/",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_LIST_GROUPS,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, ""),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
		"ErrorCaseInvalidOrg": {
			org:        "%org1",
			pathPrefix: "/example/das/",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org %org1",
			},
		},
		"ErrorCaseInvalidPath": {
			org:        "org1",
			pathPrefix: "/example/das",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: PathPrefix /example/das",
			},
		},
		"ErrorCaseInternalErrorGetGroupsFiltered": {
			org:        "org1",
			pathPrefix: "/path/",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getGroupsFilteredMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
		"ErrorCaseNotAuthenticatedUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "org1",
			pathPrefix: "/path/",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
			getGroupsFilteredMethodResult: []Group{
				{
					Name: "group1",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "group1"),
				},
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrorCaseUnauthorizedResource": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org: "org1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/*",
			},
			getGroupsFilteredMethodResult: []Group{
				{
					Name: "group1",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "group1"),
				},
				{
					Name: "group2",
					Org:  "org2",
					Path: "/path2/",
					Urn:  CreateUrn("org2", RESOURCE_GROUP, "/path2/", "group2"),
				},
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/1/",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_LIST_GROUPS,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, ""),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								GROUP_ACTION_LIST_GROUPS,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, ""),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org: "org1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/*",
			},
			getGroupsFilteredMethodResult: []Group{
				{
					Name: "group1",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "group1"),
				},
				{
					Name: "group2",
					Org:  "org2",
					Path: "/path2/",
					Urn:  CreateUrn("org2", RESOURCE_GROUP, "/path2/", "group2"),
				},
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/1/",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:         "POLICY-USER-ID",
					Name:       "policyUser",
					Org:        "org1",
					Path:       "/path/",
					Urn:        CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
	}

	for x, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetGroupsFilteredMethod][0] = testcase.getGroupsFilteredMethodResult
		testRepo.ArgsOut[GetGroupsFilteredMethod][1] = testcase.getGroupsFilteredMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult

		groups, err := testAPI.ListGroups(testcase.requestInfo, testcase.org, testcase.pathPrefix)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.expectedGroups, groups)
	}
}

func TestAuthAPI_UpdateGroup(t *testing.T) {
	testcases := map[string]struct {
		requestInfo  RequestInfo
		org          string
		groupName    string
		newGroupName string
		newPath      string
		// Expected result
		expectedGroup *Group
		wantError     error
		// Manager Results
		getGroupByNameResult            *Group
		getGroupMembersResult           []User
		getGroupsByUserIDResult         []Group
		getAttachedPoliciesResult       []Policy
		getUserByExternalIDResult       *User
		updateGroupResult               *Group
		getGroupByNameMethodSpecialFunc func(string, string) (*Group, error)
		// API Errors
		getGroupByNameMethodErr      error
		getUserByExternalIDMethodErr error
		updateGroupMethodErr         error
	}{
		"OKCaseAdmin": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:          "123",
			groupName:    "group1",
			newGroupName: "newName",
			newPath:      "/new/",
			expectedGroup: &Group{
				ID:   "12345",
				Name: "newName",
				Org:  "123",
				Path: "/new/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/new/", "test"),
			},
			getGroupByNameResult: &Group{
				ID:   "12345",
				Name: "group1",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "test"),
			},
			updateGroupResult: &Group{
				ID:   "12345",
				Name: "newName",
				Org:  "123",
				Path: "/new/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/new/", "test"),
			},
		},
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:          "org1",
			groupName:    "group1",
			newGroupName: "newName",
			newPath:      "/new/",
			expectedGroup: &Group{
				ID:   "12345",
				Name: "newName",
				Org:  "org1",
				Path: "/new/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/new/", "test"),
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "group1",
				Org:  "org1",
				Path: "/path/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "group1"),
			},
			updateGroupResult: &Group{
				ID:   "12345",
				Name: "newName",
				Org:  "org1",
				Path: "/new/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/new/", "test"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Org:  "org1",
					Path: "/path/1/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_GET_GROUP,
								GROUP_ACTION_UPDATE_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, ""),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
		"ErrorCaseInvalidName": {
			org:          "123",
			newGroupName: "%$%&&",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: new name %$%&&",
			},
		},
		"ErrorCaseInvalidPath": {
			org:          "123",
			newGroupName: "group1",
			newPath:      "/$",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: new path /$",
			},
		},
		"ErrorCaseInvalidOrg": {
			org:          "$^**!",
			groupName:    "group1",
			newGroupName: "newName",
			newPath:      "/new/",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org $^**!",
			},
		},
		"ErrorCaseGroupNotFound": {
			org:          "123",
			groupName:    "group1",
			newGroupName: "newName",
			newPath:      "/new/",
			wantError: &Error{
				Code: GROUP_BY_ORG_AND_NAME_NOT_FOUND,
			},
			getGroupByNameMethodErr: &database.Error{
				Code: database.GROUP_NOT_FOUND,
			},
		},
		"ErrorCaseUnauthorizedUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:          "123",
			groupName:    "group1",
			newGroupName: "newName",
			newPath:      "/new/",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "groupUser",
				Org:  "org1",
				Path: "/path/1/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "groupUser"),
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrorCaseGroupAlreadyExist": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:          "123",
			groupName:    "group1",
			newGroupName: "newName",
			newPath:      "/new/",
			wantError: &Error{
				Code:    GROUP_ALREADY_EXIST,
				Message: "Group name: newName already exists",
			},
			getGroupByNameMethodSpecialFunc: func(org string, name string) (*Group, error) {
				if org == "123" && name == "group1" {
					return &Group{
						ID:   "GROUP-USER-ID",
						Name: "group1",
						Org:  "org1",
						Path: "/new/",
						Urn:  CreateUrn("org1", RESOURCE_GROUP, "/new/", "group1"),
					}, nil
				} else {
					return &Group{
						ID:   "GROUP-USER-ID2",
						Name: name,
						Org:  org,
						Path: "/sdada/",
						Urn:  CreateUrn("org1", RESOURCE_GROUP, "/sdada/", name),
					}, nil
				}
			},
		},
		"ErrorCaseGetGroupDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:          "123",
			groupName:    "group1",
			newGroupName: "newName",
			newPath:      "/new/",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getGroupByNameMethodSpecialFunc: func(org string, name string) (*Group, error) {
				if org == "123" && name == "group1" {
					return &Group{
						ID:   "GROUP-USER-ID",
						Name: "group1",
						Org:  "123",
						Path: "/new/",
						Urn:  CreateUrn("org1", RESOURCE_GROUP, "/new/", "group1"),
					}, nil
				} else {
					return nil, &database.Error{
						Code: database.INTERNAL_ERROR,
					}
				}
			},
		},
		"ErrorCaseUnauthorizedUpdateGroup": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:          "123",
			groupName:    "group1",
			newGroupName: "newName",
			newPath:      "/new/",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/path/groupUser",
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "groupUser",
				Org:  "org1",
				Path: "/path/1/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "groupUser"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Org:  "org1",
					Path: "/path/1/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, ""),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
		"ErrorCaseDenyUpdateGroup": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:          "org1",
			groupName:    "group1",
			newGroupName: "newName",
			newPath:      "/new/",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/path/group1",
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "group1",
				Org:  "org1",
				Path: "/path/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "group1"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Org:  "org1",
					Path: "/path/1/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_GET_GROUP,
								GROUP_ACTION_UPDATE_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, ""),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								GROUP_ACTION_UPDATE_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, "/path"),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
		"ErrorCaseNoPermissionsToUpdateTarget": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:          "123",
			groupName:    "group1",
			newGroupName: "newName",
			newPath:      "/new/",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:123:group/new/newName",
			},
			getGroupByNameMethodSpecialFunc: func(org string, name string) (*Group, error) {
				if org == "123" && name == "group1" {
					return &Group{
						ID:   "GROUP-USER-ID",
						Name: "group1",
						Org:  "123",
						Path: "/path/",
						Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "group1"),
					}, nil
				} else {
					return nil, &database.Error{
						Code: database.GROUP_NOT_FOUND,
					}
				}
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "group1",
					Org:  "123",
					Path: "/new/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "123",
					Path: "/path/",
					Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_GET_GROUP,
								GROUP_ACTION_UPDATE_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_GROUP, "/path/"),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
		"ErrorCaseDenyToUpdateTarget": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:          "123",
			groupName:    "group1",
			newGroupName: "newName",
			newPath:      "/new/",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:123:group/new/newName",
			},
			getGroupByNameMethodSpecialFunc: func(org string, name string) (*Group, error) {
				if org == "123" && name == "group1" {
					return &Group{
						ID:   "GROUP-USER-ID",
						Name: "group1",
						Org:  "123",
						Path: "/path/",
						Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "group1"),
					}, nil
				} else {
					return nil, &database.Error{
						Code: database.GROUP_NOT_FOUND,
					}
				}
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "group1",
					Org:  "123",
					Path: "/new/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "123",
					Path: "/path/",
					Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_GET_GROUP,
								GROUP_ACTION_UPDATE_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_GROUP, ""),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								GROUP_ACTION_UPDATE_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_GROUP, "/new/"),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
		"ErrorCaseNoPermission": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:          "org1",
			groupName:    "group1",
			newGroupName: "newName",
			newPath:      "/new/",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/path/group1",
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "group1",
				Org:  "org1",
				Path: "/path/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "group1"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Org:  "org1",
					Path: "/path/1/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:         "POLICY-USER-ID",
					Name:       "policyUser",
					Org:        "org1",
					Path:       "/path/",
					Urn:        CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
		"ErrorCaseUpdateGroupDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:          "123",
			groupName:    "group1",
			newGroupName: "newName",
			newPath:      "/new/",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getGroupByNameResult: &Group{
				ID:   "12345",
				Name: "group1",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "test"),
			},
			updateGroupMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
	}

	for x, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[UpdateGroupMethod][0] = testcase.updateGroupResult
		testRepo.ArgsOut[UpdateGroupMethod][1] = testcase.updateGroupMethodErr
		testRepo.ArgsOut[GetGroupByNameMethod][0] = testcase.getGroupByNameResult
		testRepo.ArgsOut[GetGroupByNameMethod][1] = testcase.getGroupByNameMethodErr
		testRepo.SpecialFuncs[GetGroupByNameMethod] = testcase.getGroupByNameMethodSpecialFunc
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult

		group, err := testAPI.UpdateGroup(testcase.requestInfo, testcase.org, testcase.groupName, testcase.newGroupName, testcase.newPath)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.expectedGroup, group)
	}
}

func TestAuthAPI_RemoveGroup(t *testing.T) {
	testcases := map[string]struct {
		//API method args
		requestInfo RequestInfo
		name        string
		org         string
		// Expected result
		wantError error
		// Manager Results
		getUserByExternalIDResult  *User
		getGroupsByUserIDResult    []Group
		getAttachedPoliciesResult  []Policy
		getGroupByNameMethodResult *Group
		// API Errors
		getUserByExternalIDMethodErr error
		getGroupByNameMethodErr      error
		removeGroupMethodErr         error
		getGroupsByUserIDError       error
	}{
		"OKCaseAdminUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			name: "group1",
			org:  "org1",
			getGroupByNameMethodResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/example/",
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
			},
		},
		"OkCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "group1",
			org:  "org1",
			getGroupByNameMethodResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/example/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/example/", "group1"),
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("org1", RESOURCE_USER, "/example/", "123456"),
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/example/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/example/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_DELETE_GROUP,
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, ""),
							},
						},
					},
				},
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/example/",
					Org:  "org1",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/example/", "group1"),
				},
			},
		},
		"ErrorCaseInvalidName": {
			name: "invalid*",
			org:  "org1",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name invalid*",
			},
		},
		"ErrorCaseInvalidOrg": {
			name: "n1",
			org:  "**^!$%&",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org **^!$%&",
			},
		},
		"ErrorCaseGroupNotFound": {
			name: "group1",
			org:  "org1",
			wantError: &Error{
				Code: GROUP_BY_ORG_AND_NAME_NOT_FOUND,
			},
			getGroupByNameMethodErr: &database.Error{
				Code: database.GROUP_NOT_FOUND,
			},
		},
		"ErrorCaseUnauthorizedUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:  "123",
			name: "group1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
			getGroupByNameMethodResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "groupUser",
				Org:  "org1",
				Path: "/path/1/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "groupUser"),
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrorCaseImplicitUnauthorizedDeleteGroup": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "group1",
			org:  "org1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/example/group1",
			},
			getGroupByNameMethodResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/example/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/example/", "group1"),
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("org1", RESOURCE_USER, "/example/", "123456"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/example/",
					Org:  "org1",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/example/", "group1"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/example/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/example/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, ""),
							},
						},
					},
				},
			},
		},
		"ErrorCaseExplicitUnauthorizedDeleteGroup": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "group1",
			org:  "org1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/example/group1",
			},
			getGroupByNameMethodResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/example/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/example/", "group1"),
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("org1", RESOURCE_USER, "/example/", "123456"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/example/",
					Org:  "org1",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/example/", "group1"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/example/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/example/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_DELETE_GROUP,
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, ""),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								GROUP_ACTION_DELETE_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, "/example/group1"),
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
			name: "group1",
			org:  "org1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/example/group1",
			},
			getGroupByNameMethodResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/example/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/example/", "group1"),
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("org1", RESOURCE_USER, "/example/", "123456"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/example/",
					Org:  "org1",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/example/", "group1"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:         "POLICY-USER-ID",
					Name:       "policyUser",
					Org:        "org1",
					Path:       "/example/",
					Urn:        CreateUrn("org1", RESOURCE_POLICY, "/example/", "policyUser"),
					Statements: &[]Statement{},
				},
			},
		},
		"ErrorCaseDeleteGroupDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			name: "group1",
			org:  "org1",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getGroupByNameMethodResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/example/",
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
			},
			removeGroupMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
	}

	for x, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetGroupByNameMethod][0] = testcase.getGroupByNameMethodResult
		testRepo.ArgsOut[GetGroupByNameMethod][1] = testcase.getGroupByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][1] = testcase.getGroupsByUserIDError
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		testRepo.ArgsOut[RemoveGroupMethod][0] = testcase.removeGroupMethodErr

		err := testAPI.RemoveGroup(testcase.requestInfo, testcase.org, testcase.name)
		checkMethodResponse(t, x, testcase.wantError, err, nil, nil)
	}
}

func TestAuthAPI_AddMember(t *testing.T) {
	testcases := map[string]struct {
		// API Method args
		requestInfo RequestInfo
		userID      string
		org         string
		groupName   string
		// Expected result
		wantError error
		// Manager Results
		getGroupsByUserIDResult   []Group
		getAttachedPoliciesResult []Policy
		getUserByExternalIDResult *User
		getGroupByNameResult      *Group
		isMemberOfGroupResult     bool
		// Manager Errors
		getUserByExternalIDMethodErr error
		getGroupByNameMethodErr      error
		addMemberMethodErr           error
		isMemberOfGroupMethodErr     error
	}{
		"OkCaseAdmin": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			userID:    "12345",
			org:       "org1",
			groupName: "group1",
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "12345",
				Path:       "/test/asd/",
			},
			getGroupByNameResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/test/asd/",
			},
			isMemberOfGroupResult: false,
		},
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			userID:    "12345",
			org:       "org1",
			groupName: "group1",
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								"iam:*",
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, ""),
								GetUrnPrefix("", RESOURCE_USER, ""),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "group1",
				Org:  "org1",
				Path: "/path/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "group1"),
			},
			isMemberOfGroupResult: false,
		},
		"ErrorCaseInvalidExternalID": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			userID:    "d*%$",
			groupName: "group",
			org:       "org1",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: externalId d*%$",
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "group1",
				Org:  "org1",
				Path: "/path/1/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/1/", "group1"),
			},
		},
		"ErrorCaseInvalidOrg": {
			userID:    "12345",
			groupName: "group1",
			org:       "!^**$%&",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org !^**$%&",
			},
		},
		"ErrorCaseInvalidGroupName": {
			userID:    "12345",
			org:       "org1",
			groupName: "d*%$",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name d*%$",
			},
		},
		"ErrorCaseGroupNotFound": {
			userID:    "12345",
			org:       "org1",
			groupName: "group1",
			wantError: &Error{
				Code: GROUP_BY_ORG_AND_NAME_NOT_FOUND,
			},
			getGroupByNameMethodErr: &database.Error{
				Code: database.GROUP_NOT_FOUND,
			},
		},
		"ErrorCaseUnauthorizedUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			userID:    "12345",
			org:       "org1",
			groupName: "group1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "groupUser",
				Org:  "org1",
				Path: "/path/1/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "groupUser"),
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrorCaseUnauthorizedResource": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			userID:    "12345",
			org:       "org1",
			groupName: "group1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/path/1/group1",
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Org:  "org1",
					Path: "/path/1/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, ""),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "group1",
				Org:  "org1",
				Path: "/path/1/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/1/", "group1"),
			},
		},
		"ErrorCaseDenyAddMember": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			userID:    "12345",
			org:       "org1",
			groupName: "group1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/path/group1",
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "deny",
							Actions: []string{
								GROUP_ACTION_ADD_MEMBER,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, "/path/"),
							},
						},
						{
							Effect: "allow",
							Actions: []string{
								"iam:*",
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, ""),
								GetUrnPrefix("", RESOURCE_USER, ""),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "group1",
				Org:  "org1",
				Path: "/path/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "group1"),
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			userID:    "12345",
			org:       "org1",
			groupName: "group1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/path/group1",
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:         "POLICY-USER-ID",
					Name:       "policyUser",
					Org:        "org1",
					Path:       "/path/",
					Urn:        CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "group1",
				Org:  "org1",
				Path: "/path/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "group1"),
			},
		},
		"ErrorCaseUserNotFound": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			userID:    "12345",
			org:       "org1",
			groupName: "group1",
			wantError: &Error{
				Code: USER_BY_EXTERNAL_ID_NOT_FOUND,
			},
			getGroupByNameResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/test/asd/",
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrorCaseIsAlreadyMember": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			userID:    "12345",
			org:       "org1",
			groupName: "group1",
			wantError: &Error{
				Code:    USER_IS_ALREADY_A_MEMBER_OF_GROUP,
				Message: "User: 12345 is already a member of Group: group1",
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "12345",
				Path:       "/test/asd/",
			},
			getGroupByNameResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/test/asd/",
			},
			isMemberOfGroupResult: true,
		},
		"ErrorCaseIsMemberDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			userID:    "12345",
			org:       "org1",
			groupName: "group1",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "12345",
				Path:       "/test/asd/",
			},
			getGroupByNameResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/test/asd/",
			},
			isMemberOfGroupMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
		"ErrorCaseAddMemberDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			userID:    "12345",
			org:       "org1",
			groupName: "group1",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "12345",
				Path:       "/test/asd/",
			},
			getGroupByNameResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/test/asd/",
			},
			isMemberOfGroupResult: false,
			addMemberMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
	}

	for x, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[AddMemberMethod][0] = testcase.addMemberMethodErr
		testRepo.ArgsOut[GetGroupByNameMethod][0] = testcase.getGroupByNameResult
		testRepo.ArgsOut[GetGroupByNameMethod][1] = testcase.getGroupByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		testRepo.ArgsOut[IsMemberOfGroupMethod][0] = testcase.isMemberOfGroupResult
		testRepo.ArgsOut[IsMemberOfGroupMethod][1] = testcase.isMemberOfGroupMethodErr

		err := testAPI.AddMember(testcase.requestInfo, testcase.userID, testcase.groupName, testcase.org)
		checkMethodResponse(t, x, testcase.wantError, err, nil, nil)
	}
}

func TestAuthAPI_RemoveMember(t *testing.T) {
	testcases := map[string]struct {
		// API Method args
		requestInfo RequestInfo
		userID      string
		groupName   string
		org         string
		// Expected result
		wantError error
		// Manager Results
		getGroupByNameResult      *Group
		getUserByExternalIDResult *User
		getGroupsByUserIDResult   []Group
		getAttachedPoliciesResult []Policy
		isMemberOfGroupResult     bool
		// Manager Errors
		getGroupByNameMethodErr      error
		getUserByExternalIDMethodErr error
		isMemberOfGroupMethodErr     error
		removeMemberMethodErr        error
	}{
		"OkCaseAdmin": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			userID:    "12345",
			groupName: "group1",
			org:       "org1",
			getGroupByNameResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/test/",
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "12345",
				Path:       "/test/",
			},
			isMemberOfGroupResult: true,
		},
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			userID:    "12345",
			groupName: "group1",
			org:       "org1",
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "groupUser",
				Org:  "org1",
				Path: "/path/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "groupUser"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Org:  "org1",
					Path: "/path/1/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_REMOVE_MEMBER,
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, ""),
							},
						},
						{
							Effect: "allow",
							Actions: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, ""),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			isMemberOfGroupResult: true,
		},
		"ErrorCaseInvalidExternalID": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			userID:    "$%&",
			org:       "org1",
			groupName: "group1",
			getGroupByNameResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/test/",
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: externalId $%&",
			},
		},
		"ErrorCaseInvalidOrg": {
			userID:    "12345",
			org:       "$**^%&!",
			groupName: "group1",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org $**^%&!",
			},
		},
		"ErrorCaseInvalidName": {
			userID:    "12345",
			org:       "org1",
			groupName: "$%&",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name $%&",
			},
		},
		"ErrorCaseGroupNotFound": {
			userID:    "12345",
			org:       "org1",
			groupName: "group1",
			wantError: &Error{
				Code: GROUP_BY_ORG_AND_NAME_NOT_FOUND,
			},
			getGroupByNameMethodErr: &database.Error{
				Code: database.GROUP_NOT_FOUND,
			},
		},
		"ErrorCaseUnauthorizedUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			userID:    "12345",
			org:       "org1",
			groupName: "group1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "groupUser",
				Org:  "org1",
				Path: "/path/1/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "groupUser"),
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrorCaseUnauthorizedRemoveMember": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			userID:    "12345",
			groupName: "group1",
			org:       "org1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/path/group1",
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "group1",
				Org:  "org1",
				Path: "/path/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "group1"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Org:  "org1",
					Path: "/path/1/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, ""),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
		"ErrorCaseUnauthorizedGetUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			userID:    "1234",
			groupName: "group1",
			org:       "org1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam::user/path/1234",
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "groupUser",
				Org:  "org1",
				Path: "/path/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "groupUser"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Org:  "org1",
					Path: "/path/1/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_GET_GROUP,
								GROUP_ACTION_REMOVE_MEMBER,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, ""),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								USER_ACTION_GET_USER,
							},
							Resources: []string{
								GetUrnPrefix("", RESOURCE_USER, ""),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
		"ErrorCaseDenyRemoveMember": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			userID:    "12345",
			groupName: "group1",
			org:       "org1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/path/groupUser",
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "groupUser",
				Org:  "org1",
				Path: "/path/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "groupUser"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Org:  "org1",
					Path: "/path/1/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_GET_GROUP,
								GROUP_ACTION_REMOVE_MEMBER,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, ""),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								GROUP_ACTION_REMOVE_MEMBER,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, "/path/"),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			userID:    "12345",
			groupName: "group1",
			org:       "org1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/path/groupUser",
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "groupUser",
				Org:  "org1",
				Path: "/path/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "groupUser"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Org:  "org1",
					Path: "/path/1/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:         "POLICY-USER-ID",
					Name:       "policyUser",
					Org:        "org1",
					Path:       "/path/",
					Urn:        CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
		"ErrorCaseUserNotFound": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			userID:    "12345",
			groupName: "group1",
			org:       "org1",
			wantError: &Error{
				Code: USER_BY_EXTERNAL_ID_NOT_FOUND,
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "groupUser",
				Org:  "org1",
				Path: "/path/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "groupUser"),
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrorCaseIsMemberDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			userID:    "12345",
			groupName: "group1",
			org:       "org1",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "groupUser",
				Org:  "org1",
				Path: "/path/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "groupUser"),
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			isMemberOfGroupMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
		"ErrorCaseIsNotMember": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			userID:    "1234",
			groupName: "group1",
			org:       "org1",
			wantError: &Error{
				Code:    USER_IS_NOT_A_MEMBER_OF_GROUP,
				Message: "User with externalId 1234 is not a member of group with org org1 and name groupUser",
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "groupUser",
				Org:  "org1",
				Path: "/path/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "groupUser"),
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			isMemberOfGroupResult: false,
		},
		"ErrorCaseRemoveMemberDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			userID:    "12345",
			groupName: "group1",
			org:       "org1",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "groupUser",
				Org:  "org1",
				Path: "/path/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "groupUser"),
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			isMemberOfGroupResult: true,
			removeMemberMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
	}

	for x, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetGroupByNameMethod][0] = testcase.getGroupByNameResult
		testRepo.ArgsOut[GetGroupByNameMethod][1] = testcase.getGroupByNameMethodErr
		testRepo.ArgsOut[IsMemberOfGroupMethod][0] = testcase.isMemberOfGroupResult
		testRepo.ArgsOut[IsMemberOfGroupMethod][1] = testcase.isMemberOfGroupMethodErr
		testRepo.ArgsOut[RemoveMemberMethod][0] = testcase.removeMemberMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult

		err := testAPI.RemoveMember(testcase.requestInfo, testcase.userID, testcase.groupName, testcase.org)
		checkMethodResponse(t, x, testcase.wantError, err, nil, nil)
	}
}

func TestAuthAPI_ListMembers(t *testing.T) {
	testcases := map[string]struct {
		// API Method args
		requestInfo RequestInfo
		org         string
		groupName   string
		// Expected result
		expectedMembers []string
		wantError       error
		// Manager Results
		getGroupByNameResult      *Group
		getGroupMembersResult     []User
		getGroupsByUserIDResult   []Group
		getAttachedPoliciesResult []Policy
		getUserByExternalIDResult *User
		// API Errors
		getGroupByNameMethodErr      error
		getUserByExternalIDMethodErr error
		getGroupMembersMethodErr     error
	}{
		"OkCaseAdmin": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:       "org1",
			groupName: "group1",
			expectedMembers: []string{
				"member1",
				"member2",
			},
			getGroupByNameResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/test/",
			},
			getGroupMembersResult: []User{
				{
					ID:         "12345",
					ExternalID: "member1",
					Path:       "/test/",
				},
				{
					ID:         "123456",
					ExternalID: "member2",
					Path:       "/test/",
				},
			},
		},
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:       "org1",
			groupName: "group1",
			expectedMembers: []string{
				"member1",
				"member2",
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "groupUser",
				Org:  "org1",
				Path: "/path/1/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "groupUser"),
			},
			getGroupMembersResult: []User{
				{
					ID:         "12345",
					ExternalID: "member1",
					Path:       "/test/",
				},
				{
					ID:         "123456",
					ExternalID: "member2",
					Path:       "/test/",
				},
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Org:  "org1",
					Path: "/path/1/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_LIST_MEMBERS,
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, "/path/"),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
		"ErrorCaseInvalidName": {
			org:       "org1",
			groupName: "*%$",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name *%$",
			},
		},
		"ErrorCaseInvalidOrg": {
			org:       "!^**$%&",
			groupName: "g1",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org !^**$%&",
			},
		},
		"ErrorCaseGroupNotFound": {
			org:       "org1",
			groupName: "group1",
			getGroupByNameMethodErr: &database.Error{
				Code: database.GROUP_NOT_FOUND,
			},
			wantError: &Error{
				Code: GROUP_BY_ORG_AND_NAME_NOT_FOUND,
			},
		},
		"ErrorCaseNotAuthenticatedUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:       "org1",
			groupName: "group1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
			getGroupByNameResult: &Group{
				Name: "group1",
				Org:  "org1",
				Path: "/path/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "group1"),
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrorCaseUnauthorizedResource": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:       "org1",
			groupName: "group1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/path/groupUser",
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "groupUser",
				Org:  "org1",
				Path: "/path/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "groupUser"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Org:  "org1",
					Path: "/path/1/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, ""),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
		"ErrorCaseDenyListMembers": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:       "org1",
			groupName: "group1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/path/groupUser",
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "groupUser",
				Org:  "org1",
				Path: "/path/1/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "groupUser"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Org:  "org1",
					Path: "/path/1/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "deny",
							Actions: []string{
								GROUP_ACTION_LIST_MEMBERS,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, "/path/"),
							},
						},
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_LIST_MEMBERS,
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, ""),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:       "org1",
			groupName: "group1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/path/groupUser",
			},
			getGroupByNameResult: &Group{
				ID:   "GROUP-USER-ID",
				Name: "groupUser",
				Org:  "org1",
				Path: "/path/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "groupUser"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Org:  "org1",
					Path: "/path/1/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:         "POLICY-USER-ID",
					Name:       "policyUser",
					Org:        "org1",
					Path:       "/path/",
					Urn:        CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
		"ErrorCaseListMembersDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:       "org1",
			groupName: "group1",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getGroupByNameResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/test/",
			},
			getGroupMembersMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
	}

	for x, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetGroupByNameMethod][0] = testcase.getGroupByNameResult
		testRepo.ArgsOut[GetGroupByNameMethod][1] = testcase.getGroupByNameMethodErr
		testRepo.ArgsOut[GetGroupMembersMethod][0] = testcase.getGroupMembersResult
		testRepo.ArgsOut[GetGroupMembersMethod][1] = testcase.getGroupMembersMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult

		members, err := testAPI.ListMembers(testcase.requestInfo, testcase.org, testcase.groupName)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.expectedMembers, members)
	}
}

func TestAuthAPI_AttachPolicyToGroup(t *testing.T) {
	testcases := map[string]struct {
		requestInfo RequestInfo
		org         string
		groupName   string
		policyName  string
		// Expected result
		wantError error
		// Manager Results
		getGroupByNameResult      *Group
		getPolicyByNameResult     *Policy
		getUserByExternalIDResult *User
		getGroupsByUserIDResult   []Group
		getAttachedPoliciesResult []Policy
		isAttachedToGroupResult   bool
		// API Errors
		getGroupByNameMethodErr      error
		getPolicyByNameMethodErr     error
		getUserByExternalIDMethodErr error
		isAttachedToGroupMethodErr   error
		attachPolicyMethodErr        error
	}{
		"OkCaseAdmin": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			groupName:  "group1",
			policyName: "policy1",
			getGroupByNameResult: &Group{
				ID:   "12345",
				Name: "group1",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "test"),
			},
			getPolicyByNameResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
			isAttachedToGroupResult: false,
		},
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "123",
			groupName:  "group1",
			policyName: "policy1",
			getGroupByNameResult: &Group{
				ID:   "12345",
				Name: "group1",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "test"),
			},
			getPolicyByNameResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "group1",
					Org:  "123",
					Path: "/path/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "123",
					Path: "/path/",
					Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_GET_GROUP,
								GROUP_ACTION_ATTACH_GROUP_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_GROUP, "/path/"),
							},
						},
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			isAttachedToGroupResult: false,
		},
		"ErrorCaseInvalidGroupName": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:       "123",
			groupName: "$%·",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name $%·",
			},
		},
		"ErrorCaseInvalidOrg": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:       "$%&!",
			groupName: "g1",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org $%&!",
			},
		},
		"ErrorCaseInvalidPolicyName": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			groupName:  "group1",
			policyName: "$·%",
			getGroupByNameResult: &Group{
				ID:   "12345",
				Name: "group1",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "group1"),
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name $·%",
			},
		},
		"ErrorCaseGroupNotFound": {
			org:        "123",
			groupName:  "group1",
			policyName: "policy1",
			wantError: &Error{
				Code: GROUP_BY_ORG_AND_NAME_NOT_FOUND,
			},
			getGroupByNameMethodErr: &database.Error{
				Code: database.GROUP_NOT_FOUND,
			},
		},
		"ErrorCaseNotAuthenticatedUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "123",
			groupName:  "group1",
			policyName: "policy1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
			getGroupByNameResult: &Group{
				Name: "group1",
				Org:  "org1",
				Path: "/path/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "group1"),
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrorCaseUnauthorizedToAttach": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "123",
			groupName:  "group1",
			policyName: "policy1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:123:group/path/group1",
			},
			getGroupByNameResult: &Group{
				ID:   "12345",
				Name: "group1",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "group1"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "group1",
					Org:  "123",
					Path: "/path/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "123",
					Path: "/path/",
					Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_GROUP, ""),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
		"ErrorCaseDenyToAttach": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "123",
			groupName:  "group1",
			policyName: "policy1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:123:group/path/group1",
			},
			getGroupByNameResult: &Group{
				ID:   "12345",
				Name: "group1",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "group1"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "group1",
					Org:  "123",
					Path: "/path/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "123",
					Path: "/path/",
					Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_GET_GROUP,
								GROUP_ACTION_ATTACH_GROUP_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_GROUP, ""),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								GROUP_ACTION_ATTACH_GROUP_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_GROUP, "/path/"),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "123",
			groupName:  "group1",
			policyName: "policy1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:123:group/path/group1",
			},
			getGroupByNameResult: &Group{
				ID:   "12345",
				Name: "group1",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "group1"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "group1",
					Org:  "123",
					Path: "/path/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:         "POLICY-USER-ID",
					Name:       "policyUser",
					Org:        "123",
					Path:       "/path/",
					Urn:        CreateUrn("123", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
		"ErrorCaseIsAttachedDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			groupName:  "group1",
			policyName: "policy1",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getGroupByNameResult: &Group{
				ID:   "12345",
				Name: "group1",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "group1"),
			},
			getPolicyByNameResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
			isAttachedToGroupMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
		"ErrorCasePolicyIsAlreadyAttached": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			groupName:  "group1",
			policyName: "policy1",
			wantError: &Error{
				Code:    POLICY_IS_ALREADY_ATTACHED_TO_GROUP,
				Message: "Policy: test is already attached to Group: group1",
			},
			getGroupByNameResult: &Group{
				ID:   "12345",
				Name: "group1",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "group1"),
			},
			getPolicyByNameResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
			isAttachedToGroupResult: true,
		},
		"ErrorCasePolicyNotFound": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			groupName:  "group1",
			policyName: "policy1",
			wantError: &Error{
				Code: POLICY_BY_ORG_AND_NAME_NOT_FOUND,
			},
			getGroupByNameResult: &Group{
				ID:   "12345",
				Name: "group1",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "group1"),
			},
			getPolicyByNameMethodErr: &database.Error{
				Code: database.POLICY_NOT_FOUND,
			},
		},
		"ErrorCaseAttachPolicyDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			groupName:  "group1",
			policyName: "policy1",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getGroupByNameResult: &Group{
				ID:   "12345",
				Name: "group1",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "group1"),
			},
			getPolicyByNameResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
			isAttachedToGroupResult: false,
			attachPolicyMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
	}

	for x, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetGroupByNameMethod][0] = testcase.getGroupByNameResult
		testRepo.ArgsOut[GetGroupByNameMethod][1] = testcase.getGroupByNameMethodErr
		testRepo.ArgsOut[GetPolicyByNameMethod][0] = testcase.getPolicyByNameResult
		testRepo.ArgsOut[GetPolicyByNameMethod][1] = testcase.getPolicyByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		testRepo.ArgsOut[IsAttachedToGroupMethod][0] = testcase.isAttachedToGroupResult
		testRepo.ArgsOut[IsAttachedToGroupMethod][1] = testcase.isAttachedToGroupMethodErr
		testRepo.ArgsOut[AttachPolicyMethod][0] = testcase.attachPolicyMethodErr

		err := testAPI.AttachPolicyToGroup(testcase.requestInfo, testcase.org, testcase.groupName, testcase.policyName)
		checkMethodResponse(t, x, testcase.wantError, err, nil, nil)
	}
}

func TestAuthAPI_DetachPolicyToGroup(t *testing.T) {
	testcases := map[string]struct {
		requestInfo RequestInfo
		org         string
		groupName   string
		policyName  string
		// Expected result
		wantError error
		// Manager Results
		getGroupByNameResult      *Group
		getPolicyByNameResult     *Policy
		getUserByExternalIDResult *User
		getGroupsByUserIDResult   []Group
		getAttachedPoliciesResult []Policy
		isAttachedToGroupResult   bool
		// API Errors
		getGroupByNameMethodErr      error
		getPolicyByNameMethodErr     error
		getUserByExternalIDMethodErr error
		isAttachedToGroupMethodErr   error
		detachPolicyMethodErr        error
	}{
		"OkCaseAdmin": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			groupName:  "group1",
			policyName: "policy1",
			getGroupByNameResult: &Group{
				ID:   "12345",
				Name: "group1",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "group1"),
			},
			getPolicyByNameResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
			isAttachedToGroupResult: true,
		},
		"OkCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "123",
			groupName:  "group1",
			policyName: "policy1",
			getGroupByNameResult: &Group{
				ID:   "12345",
				Name: "group1",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "group1"),
			},
			getPolicyByNameResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "group1",
					Org:  "123",
					Path: "/path/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "123",
					Path: "/path/",
					Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_GET_GROUP,
								GROUP_ACTION_DETACH_GROUP_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_GROUP, ""),
							},
						},
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, ""),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			isAttachedToGroupResult: true,
		},
		"ErrorCaseInvalidGroupName": {
			org:       "123",
			groupName: "$%·",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name $%·",
			},
		},
		"ErrorCaseInvalidOrg": {
			org:       "$%·",
			groupName: "g1",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org $%·",
			},
		},
		"ErrorCaseInvalidPolicyName": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			groupName:  "group1",
			policyName: "$·%",
			getGroupByNameResult: &Group{
				ID:   "12345",
				Name: "group1",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "group1"),
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name $·%",
			},
		},
		"ErrorCaseGroupNotFound": {
			org:        "123",
			groupName:  "group1",
			policyName: "policy1",
			wantError: &Error{
				Code: GROUP_BY_ORG_AND_NAME_NOT_FOUND,
			},
			getGroupByNameMethodErr: &database.Error{
				Code: database.GROUP_NOT_FOUND,
			},
		},
		"ErrorCaseNotAuthenticatedUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "123",
			groupName:  "group1",
			policyName: "policy1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
			getGroupByNameResult: &Group{
				Name: "group1",
				Org:  "org1",
				Path: "/path/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "group1"),
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrorCaseUnauthorizedToDetach": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "123",
			groupName:  "group1",
			policyName: "policy1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:123:group/path/group1",
			},
			getGroupByNameResult: &Group{
				ID:   "12345",
				Name: "group1",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "group1"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "group1",
					Org:  "123",
					Path: "/path/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "123",
					Path: "/path/",
					Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_GROUP, ""),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
		"ErrorCaseDenyToDetach": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "123",
			groupName:  "group1",
			policyName: "policy1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:123:group/path/group1",
			},
			getGroupByNameResult: &Group{
				ID:   "12345",
				Name: "group1",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "group1"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "group1",
					Org:  "123",
					Path: "/path/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "123",
					Path: "/path/",
					Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_GET_GROUP,
								GROUP_ACTION_DETACH_GROUP_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_GROUP, ""),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								GROUP_ACTION_DETACH_GROUP_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_GROUP, "/path/"),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "123",
			groupName:  "group1",
			policyName: "policy1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:123:group/path/group1",
			},
			getGroupByNameResult: &Group{
				ID:   "12345",
				Name: "group1",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "group1"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "group1",
					Org:  "123",
					Path: "/path/",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:         "POLICY-USER-ID",
					Name:       "policyUser",
					Org:        "123",
					Path:       "/path/",
					Urn:        CreateUrn("123", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
		"ErrorCasePolicyNotFound": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			groupName:  "group1",
			policyName: "policy1",
			wantError: &Error{
				Code: POLICY_BY_ORG_AND_NAME_NOT_FOUND,
			},
			getGroupByNameResult: &Group{
				ID:   "12345",
				Name: "group1",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "group1"),
			},
			getPolicyByNameMethodErr: &database.Error{
				Code: database.POLICY_NOT_FOUND,
			},
		},
		"ErrorCaseIsAttachedDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			groupName:  "group1",
			policyName: "policy1",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getGroupByNameResult: &Group{
				ID:   "12345",
				Name: "group1",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "group1"),
			},
			getPolicyByNameResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
			isAttachedToGroupMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
		"ErrorCasePolicyIsNotAttached": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			groupName:  "group1",
			policyName: "policy1",
			wantError: &Error{
				Code:    POLICY_IS_NOT_ATTACHED_TO_GROUP,
				Message: "Policy with org 123 and name test is not attached to group with org 123 and name group1",
			},
			getGroupByNameResult: &Group{
				ID:   "12345",
				Name: "group1",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "group1"),
			},
			getPolicyByNameResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
			isAttachedToGroupResult: false,
		},
		"ErrorCaseDetachPolicyDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			groupName:  "group1",
			policyName: "policy1",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getGroupByNameResult: &Group{
				ID:   "12345",
				Name: "group1",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_GROUP, "/path/", "group1"),
			},
			getPolicyByNameResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
			isAttachedToGroupResult: true,
			detachPolicyMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
	}

	for x, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetGroupByNameMethod][0] = testcase.getGroupByNameResult
		testRepo.ArgsOut[GetGroupByNameMethod][1] = testcase.getGroupByNameMethodErr
		testRepo.ArgsOut[GetPolicyByNameMethod][0] = testcase.getPolicyByNameResult
		testRepo.ArgsOut[GetPolicyByNameMethod][1] = testcase.getPolicyByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		testRepo.ArgsOut[IsAttachedToGroupMethod][0] = testcase.isAttachedToGroupResult
		testRepo.ArgsOut[IsAttachedToGroupMethod][1] = testcase.isAttachedToGroupMethodErr
		testRepo.ArgsOut[DetachPolicyMethod][0] = testcase.detachPolicyMethodErr

		err := testAPI.DetachPolicyToGroup(testcase.requestInfo, testcase.org, testcase.groupName, testcase.policyName)
		checkMethodResponse(t, x, testcase.wantError, err, nil, nil)
	}
}

func TestAuthAPI_ListAttachedGroupPolicies(t *testing.T) {
	testcases := map[string]struct {
		//API method args
		requestInfo RequestInfo
		name        string
		org         string
		// Expected result
		expectedPolicies []string
		wantError        error
		// Manager Results
		getUserByExternalIDResult  *User
		getGroupsByUserIDResult    []Group
		getAttachedPoliciesResult  []Policy
		getGroupByNameMethodResult *Group
		// API Errors
		getUserByExternalIDMethodErr error
		getAttachedPoliciesErr       error
		getGroupByNameMethodErr      error
	}{
		"OKCaseAdminUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			name: "group1",
			org:  "org1",
			getGroupByNameMethodResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/example/",
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
			},
			expectedPolicies: []string{},
		},
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "group1",
			org:  "org1",
			getGroupByNameMethodResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/example/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/example/", "group1"),
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("org1", RESOURCE_USER, "/example/", "123456"),
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/example/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/example/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_LIST_ATTACHED_GROUP_POLICIES,
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, ""),
							},
						},
					},
				},
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/example/",
					Org:  "org1",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/example/", "groupUser"),
				},
			},
			expectedPolicies: []string{"policyUser"},
		},
		"ErrorCaseInvalidName": {
			name: "invalid*",
			org:  "org1",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name invalid*",
			},
		},
		"ErrorCaseInvalidOrg": {
			name: "n1",
			org:  "!**$%&",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org !**$%&",
			},
		},
		"ErrorCaseGroupNotFound": {
			name: "group1",
			org:  "org1",
			wantError: &Error{
				Code: GROUP_BY_ORG_AND_NAME_NOT_FOUND,
			},
			getGroupByNameMethodErr: &database.Error{
				Code: database.GROUP_NOT_FOUND,
			},
		},
		"ErrorCaseNotAuthenticatedUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "group1",
			org:  "org1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
			getGroupByNameMethodResult: &Group{
				Name: "group1",
				Org:  "org1",
				Path: "/path/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "group1"),
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrorCaseImplicitUnauthorizedListAttachedGroupPolicies": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "group1",
			org:  "org1",
			getGroupByNameMethodResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/example/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/example/", "group1"),
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("org1", RESOURCE_USER, "/example/", "123456"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/example/",
					Org:  "org1",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/example/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/example/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/example/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, ""),
							},
						},
					},
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/example/group1",
			},
		},
		"ErrorCaseExplicitUnauthorizedListAttachedGroupPolicies": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "group1",
			org:  "org1",
			getGroupByNameMethodResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/example/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/example/", "group1"),
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("org1", RESOURCE_USER, "/example/", "123456"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/example/",
					Org:  "org1",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/example/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/example/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/example/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								GROUP_ACTION_LIST_ATTACHED_GROUP_POLICIES,
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, ""),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								GROUP_ACTION_LIST_ATTACHED_GROUP_POLICIES,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, "/example/group1"),
							},
						},
					},
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/example/group1",
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "group1",
			org:  "org1",
			getGroupByNameMethodResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/example/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/example/", "group1"),
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("org1", RESOURCE_USER, "/example/", "123456"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/example/",
					Org:  "org1",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/example/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:         "POLICY-USER-ID",
					Name:       "policyUser",
					Org:        "org1",
					Path:       "/example/",
					Urn:        CreateUrn("org1", RESOURCE_POLICY, "/example/", "policyUser"),
					Statements: &[]Statement{},
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:group/example/group1",
			},
		},
		"ErrorCaseGetPoliciesDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			name: "group1",
			org:  "org1",
			getGroupByNameMethodResult: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/example/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/example/", "group1"),
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("org1", RESOURCE_USER, "/example/", "123456"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/example/",
					Org:  "org1",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/example/", "group1"),
				},
			},
			getAttachedPoliciesErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
	}
	for x, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetGroupByNameMethod][0] = testcase.getGroupByNameMethodResult
		testRepo.ArgsOut[GetGroupByNameMethod][1] = testcase.getGroupByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][1] = testcase.getAttachedPoliciesErr

		policies, err := testAPI.ListAttachedGroupPolicies(testcase.requestInfo, testcase.org, testcase.name)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.expectedPolicies, policies)
	}
}
