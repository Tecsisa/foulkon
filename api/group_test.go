package api

import (
	"github.com/tecsisa/authorizr/database"
	"reflect"
	"testing"
)

func TestAddGroup(t *testing.T) {
	testcases := map[string]struct {
		// API Method args
		authUser AuthenticatedUser
		name     string
		org      string
		path     string
		// Expected results
		expectedGroup *Group
		wantError     *Error
		// Manager Results
		getUserByExternalIDResult *User
		getGroupsByUserIDResult   []Group
		getPoliciesAttachedResult []Policy
		getGroupByName            *Group
		addMemberMethodResult     *Group
		// API Errors
		getGroupByNameMethodErr      error
		getUserByExternalIDMethodErr error
		addGroupMethodErr            error
	}{
		"OKCaseAdmin": {
			authUser: AuthenticatedUser{
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
		"ErrorCaseInvalidName": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			name: "*%~#@|",
			org:  "org1",
			path: "/example/",
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"ErrorCaseInvalidPath": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			name: "group1",
			org:  "org1",
			path: "/**%%/*123",
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"ErrorCaseGroupAlreadyExists": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			name: "group1",
			org:  "org1",
			path: "/example/",
			wantError: &Error{
				Code: GROUP_ALREADY_EXIST,
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
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			name: "group1",
			org:  "org1",
			path: "/example/",
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrorCaseUnauthorizedResource": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			name: "group1",
			org:  "org1",
			path: "/test/asd/",
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
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
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								GROUP_ACTION_CREATE_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, "/test/"),
							},
						},
						Statement{
							Effect: "deny",
							Action: []string{
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
		"ErrorCaseAddGroupDBErr": {
			authUser: AuthenticatedUser{
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
			authUser: AuthenticatedUser{
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

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for x, testcase := range testcases {
		testRepo.ArgsOut[GetGroupByNameMethod][0] = testcase.getGroupByName
		testRepo.ArgsOut[GetGroupByNameMethod][1] = testcase.getGroupByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetPoliciesAttachedMethod][0] = testcase.getPoliciesAttachedResult
		testRepo.ArgsOut[AddGroupMethod][0] = testcase.expectedGroup
		testRepo.ArgsOut[AddGroupMethod][1] = testcase.addGroupMethodErr
		group, err := testAPI.AddGroup(testcase.authUser, testcase.org, testcase.name, testcase.path)
		if testcase.wantError != nil {
			apiError, ok := err.(*Error)
			if !ok || apiError == nil {
				t.Fatalf("Test %v failed. Unexpected data retrieved from error: %v", x, err)
			}
			if apiError.Code != testcase.wantError.Code {
				t.Fatalf("Test %v failed. Got error %v, expected %v", x, apiError, testcase.wantError.Code)
			}
		} else {
			if err != nil {
				t.Fatalf("Test %v failed: %v", x, err)
			} else {
				if !reflect.DeepEqual(testcase.expectedGroup, group) {
					t.Fatalf("Test %v failed. Received different groups", x)
				}
			}
		}
	}

}

func TestGetGroupByName(t *testing.T) {
	testcases := map[string]struct {
		// API Method args
		authUser AuthenticatedUser
		name     string
		org      string
		path     string
		// Expected result
		expectedGroup *Group
		wantError     *Error
		// Manager Results
		getUserByExternalIDResult  *User
		getGroupsByUserIDResult    []Group
		getPoliciesAttachedResult  []Policy
		getGroupByNameMethodResult *Group
		// API Errors
		getUserByExternalIDMethodErr error
		getGroupByNameMethodErr      error
	}{
		"OKCaseAdmin": {
			authUser: AuthenticatedUser{
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
		"ErrorCaseInvalidName": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			name: "*%~#@|",
			org:  "org1",
			path: "/example/",
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"ErrorCaseGroupNotFound": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
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
			authUser: AuthenticatedUser{
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
		"ErrorCaseNotAuthenticatedUser": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			name: "group1",
			org:  "org1",
			path: "/example/",
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
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
			authUser: AuthenticatedUser{
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
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
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
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, "/test/"),
							},
						},
						Statement{
							Effect: "deny",
							Action: []string{
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
	}

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for x, testcase := range testcases {
		testRepo.ArgsOut[GetGroupByNameMethod][0] = testcase.getGroupByNameMethodResult
		testRepo.ArgsOut[GetGroupByNameMethod][1] = testcase.getGroupByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetPoliciesAttachedMethod][0] = testcase.getPoliciesAttachedResult
		group, err := testAPI.GetGroupByName(testcase.authUser, testcase.org, testcase.name)
		if testcase.wantError != nil {
			apiError, ok := err.(*Error)
			if !ok || apiError == nil {
				t.Fatalf("Test %v failed. Unexpected data retrieved from error: %v", x, err)
			}
			if apiError.Code != testcase.wantError.Code {
				t.Fatalf("Test %v failed. Got error %v, expected %v", x, apiError, testcase.wantError.Code)
			}
		} else {
			if err != nil {
				t.Fatalf("Test %v failed: %v", x, err)
			} else {
				if !reflect.DeepEqual(testcase.expectedGroup, group) {
					t.Fatalf("Test %v failed. Received different groups", x)
				}
			}
		}
	}
}

func TestGetListGroups(t *testing.T) {
	testcases := map[string]struct {
		// API Method args
		authUser   AuthenticatedUser
		org        string
		pathPrefix string
		// Expected result
		expectedGroups []GroupIdentity
		wantError      *Error
		// Manager Results
		getGroupsFilteredMethodResult []Group
		getGroupsByUserIDResult       []Group
		getPoliciesAttachedResult     []Policy
		getUserByExternalIDResult     *User
		// API Errors
		getUserByExternalIDMethodErr error
		getGroupsFilteredMethodErr   error
	}{
		"OkCaseAdmin": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "org1",
			pathPrefix: "/",
			expectedGroups: []GroupIdentity{
				GroupIdentity{
					Org:  "org1",
					Name: "group1",
				},
				GroupIdentity{
					Org:  "org2",
					Name: "group2",
				},
			},
			getGroupsFilteredMethodResult: []Group{
				Group{
					Name: "group1",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "group1"),
				},
				Group{
					Name: "group2",
					Org:  "org2",
					Path: "/path2/",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path2/", "group2"),
				},
			},
		},
		"OkTestCaseUser": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			org: "org1",
			expectedGroups: []GroupIdentity{
				GroupIdentity{
					Org:  "org1",
					Name: "group1",
				},
			},
			getGroupsFilteredMethodResult: []Group{
				Group{
					Name: "group1",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", "group1"),
				},
				Group{
					Name: "group2",
					Org:  "org2",
					Path: "/path2/",
					Urn:  CreateUrn("org2", RESOURCE_GROUP, "/path2/", "group2"),
				},
			},
			getGroupsByUserIDResult: []Group{
				Group{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/1/",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", ""),
				},
			},
			getPoliciesAttachedResult: []Policy{
				Policy{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
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
		"ErrorCaseInvalidPath": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "org1",
			pathPrefix: "/example/das",
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"ErrorCaseInternalErrorGetGroupsFiltered": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "org1",
			pathPrefix: "/path/",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getGroupsFilteredMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
		"ErrorCaseNoPermissions": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "org1",
			pathPrefix: "/path/",
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
			getGroupsFilteredMethodResult: []Group{
				Group{
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
	}

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for x, testcase := range testcases {

		testRepo.ArgsOut[GetGroupsFilteredMethod][0] = testcase.getGroupsFilteredMethodResult
		testRepo.ArgsOut[GetGroupsFilteredMethod][1] = testcase.getGroupsFilteredMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetPoliciesAttachedMethod][0] = testcase.getPoliciesAttachedResult

		groups, err := testAPI.GetListGroups(testcase.authUser, testcase.org, testcase.pathPrefix)
		if testcase.wantError != nil {
			apiError, ok := err.(*Error)
			if !ok || apiError == nil {
				t.Fatalf("Test %v failed. Unexpected data retrieved from error: %v", x, err)
			}
			if apiError.Code != testcase.wantError.Code {
				t.Fatalf("Test %v failed. Got error %v, expected %v", x, apiError, testcase.wantError.Code)
			}
		} else {
			if err != nil {
				t.Fatalf("Test %v failed. Error: %v", x, err)
			} else {
				if !reflect.DeepEqual(groups, testcase.expectedGroups) {
					t.Fatalf("Test %v failed. Received different policies (wanted:%v / received:%v)",
						x, testcase.expectedGroups, groups)
				}
			}
		}
	}
}

func TestAddMember(t *testing.T) {
	testcases := map[string]struct {
		// API Method args
		authUser  AuthenticatedUser
		userID    string
		org       string
		groupName string
		// Expected result
		wantError *Error
		// Manager Results
		getGroupsByUserIDResult   []Group
		getPoliciesAttachedResult []Policy
		getUserByExternalIDResult *User
		getGroupByNameResult      *Group
		isMemberOfGroupResult     bool
		// API Errors
		getUserByExternalIDMethodErr error
		getGroupByNameMethodErr      error
		addMemberMethodErr           error
		isMemberOfGroupMethodErr     error
	}{
		"OkCaseAdmin": {
			authUser: AuthenticatedUser{
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
		"ErrorCaseInvalidExternalID": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			userID: "d*%$",
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"ErrorCaseInvalidGroupName": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			userID:    "12345",
			groupName: "d*%$",
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"ErrorCaseGroupNotFound": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			userID:    "12345",
			groupName: "group1",
			wantError: &Error{
				Code: GROUP_BY_ORG_AND_NAME_NOT_FOUND,
			},
			getGroupByNameMethodErr: &database.Error{
				Code: database.GROUP_NOT_FOUND,
			},
		},
		"ErrorCaseUnauthorizedUser": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			userID:    "12345",
			groupName: "group1",
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
			getGroupsByUserIDResult: []Group{
				Group{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Org:  "org1",
					Path: "/path/1/",
				},
			},
			getPoliciesAttachedResult: []Policy{
				Policy{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "deny",
							Action: []string{
								GROUP_ACTION_ADD_MEMBER,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, ""),
							},
						},
						Statement{
							Effect: "allow",
							Action: []string{
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
		"ErrorCaseNoPermissions": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			userID:    "12345",
			groupName: "group1",
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
			getGroupsByUserIDResult: []Group{
				Group{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Org:  "org1",
					Path: "/path/1/",
					Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", ""),
				},
			},
			getPoliciesAttachedResult: []Policy{
				Policy{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "org1",
					Path: "/path/",
					Urn:  CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "deny",
							Action: []string{
								GROUP_ACTION_ADD_MEMBER,
							},
							Resources: []string{
								GetUrnPrefix("org1", RESOURCE_GROUP, ""),
							},
						},
						Statement{
							Effect: "allow",
							Action: []string{
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
				Name: "groupUser",
				Org:  "org1",
				Path: "/path/1/",
				Urn:  CreateUrn("org1", RESOURCE_GROUP, "/path/", ""),
			},
		},
		"ErrorCaseUserNotFound": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			userID:    "12345",
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
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			userID:    "12345",
			groupName: "group1",
			wantError: &Error{
				Code: USER_IS_ALREADY_A_MEMBER_OF_GROUP,
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
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			userID:    "12345",
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
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			userID:    "12345",
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

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for x, testcase := range testcases {

		testRepo.ArgsOut[AddMemberMethod][0] = testcase.addMemberMethodErr
		testRepo.ArgsOut[GetGroupByNameMethod][0] = testcase.getGroupByNameResult
		testRepo.ArgsOut[GetGroupByNameMethod][1] = testcase.getGroupByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetPoliciesAttachedMethod][0] = testcase.getPoliciesAttachedResult
		testRepo.ArgsOut[IsMemberOfGroupMethod][0] = testcase.isMemberOfGroupResult
		testRepo.ArgsOut[IsMemberOfGroupMethod][1] = testcase.isMemberOfGroupMethodErr

		err := testAPI.AddMember(testcase.authUser, testcase.userID, testcase.groupName, testcase.org)
		if testcase.wantError != nil {
			apiError, ok := err.(*Error)
			if !ok || apiError == nil {
				t.Fatalf("Test %v failed. Unexpected data retrieved from error: %v", x, err)
			}
			if apiError.Code != testcase.wantError.Code {
				t.Fatalf("Test %v failed. Got error %v, expected %v", x, apiError, testcase.wantError.Code)
			}
		} else {
			if err != nil {
				t.Fatalf("Test %v failed. Error: %v", x, err)
			}
		}
	}
}
