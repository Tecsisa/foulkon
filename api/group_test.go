package api

import (
	"github.com/tecsisa/authorizr/database"
	"reflect"
	"testing"
)

func TestAddGroup(t *testing.T) {
	testcases := map[string]struct {
		authUser AuthenticatedUser
		name     string
		org      string
		path     string

		getGroupsByUserIDResult   []Group
		getPoliciesAttachedResult []Policy
		getUserByExternalIDResult *User

		getGroupByName *Group
		expectedGroup  *Group
		wantError      *Error

		addGroupMethodErr            error
		getGroupByNameMethodErr      error
		getUserByExternalIDMethodErr error
	}{
		"OKCase": {
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
		"InvalidName": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			name: "*%~#@|",
			path: "/example/",
			org:  "org1",
			expectedGroup: &Group{
				ID:   "543210",
				Name: "*%~#@|",
				Org:  "org1",
				Path: "/example/",
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
			name: "group1",
			path: "/**%%/*123",
			org:  "org1",
			expectedGroup: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/**%%/*123",
			},
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"AlreadyExists": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			name: "group1",
			org:  "org1",
			path: "/example/",
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
			},
			expectedGroup: &Group{
				ID:   "543210",
				Name: "group1",
				Org:  "org1",
				Path: "/example/",
			},
			wantError: &Error{
				Code: GROUP_ALREADY_EXIST,
			},
		},
		"NotAuthenticatedUserExist": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			name: "group1",
			org:  "org1",
			path: "/example/",
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
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
			name: "group1",
			org:  "org1",
			path: "/test/asd/",
			getGroupByNameMethodErr: &database.Error{
				Code: database.GROUP_NOT_FOUND,
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
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"addGroupDBErr": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			name: "group1",
			org:  "org1",
			path: "/example/",
			getGroupByNameMethodErr: &database.Error{
				Code: database.GROUP_NOT_FOUND,
			},
			addGroupMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
		"getGroupDBErr": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			name: "group1",
			org:  "org1",
			path: "/example/",
			getGroupByNameMethodErr: &database.Error{
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
		authUser AuthenticatedUser
		name     string
		org      string
		path     string

		getGroupsByUserIDResult   []Group
		getPoliciesAttachedResult []Policy
		getUserByExternalIDResult *User

		expectedGroup *Group
		wantError     *Error

		getUserByExternalIDMethodErr error
		getGroupByNameMethodErr      error
	}{
		"OKCase": {
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
		},
		"InvalidName": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			name: "*%~#@|",
			path: "/example/",
			org:  "org1",
			expectedGroup: &Group{
				ID:   "543210",
				Name: "*%~#@|",
				Org:  "org1",
				Path: "/example/",
			},
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"GroupNotFound": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			name: "group1",
			org:  "org1",
			path: "/example/",
			getGroupByNameMethodErr: &database.Error{
				Code: database.GROUP_NOT_FOUND,
			},
			wantError: &Error{
				Code: GROUP_BY_ORG_AND_NAME_NOT_FOUND,
			},
		},
		"GetGroupDBErr": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			name: "group1",
			org:  "org1",
			path: "/example/",
			getGroupByNameMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
		"NotAuthenticatedUserExist": {
			authUser: AuthenticatedUser{
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
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
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
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
	}

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for x, testcase := range testcases {
		testRepo.ArgsOut[GetGroupByNameMethod][0] = testcase.expectedGroup
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
		authUser   AuthenticatedUser
		org        string
		pathPrefix string

		getGroupsByUserIDResult   []Group
		getPoliciesAttachedResult []Policy
		getUserByExternalIDResult *User

		expectedGroups                []GroupIdentity
		getGroupsFilteredMethodResult []Group
		wantError                     *Error

		getUserByExternalIDMethodErr error
		getGroupsFilteredMethodErr   error
	}{
		"OkTestCaseAdmin": {
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
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
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
			getGroupsFilteredMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
		"ErrorCaseNoPermissions": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "org1",
			pathPrefix: "/path/",
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
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
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
