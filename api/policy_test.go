package api

import (
	"reflect"
	"testing"

	"github.com/tecsisa/authorizr/database"
)

func TestGetPolicyByName(t *testing.T) {
	testcases := map[string]struct {
		authUser   AuthenticatedUser
		org        string
		policyName string

		getGroupsByUserIDResult   []Group
		getPoliciesAttachedResult []Policy
		getUserByExternalIDResult *User

		getPolicyByNameMethodResult *Policy
		wantError                   *Error

		getPolicyByNameMethodErr error
	}{
		"OKCase": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
		"InternalError": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			getPolicyByNameMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
		"BadPolicyName": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "~#**!",
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"PolicyNotFound": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			getPolicyByNameMethodErr: &database.Error{
				Code: database.POLICY_NOT_FOUND,
			},
			wantError: &Error{
				Code: POLICY_BY_ORG_AND_NAME_NOT_FOUND,
			},
		},
		"NoPermissions": {
			authUser: AuthenticatedUser{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "example",
			policyName: "policyUser",
			getPolicyByNameMethodResult: &Policy{
				ID:   "POLICY-USER-ID",
				Name: "policyUser",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
				Statements: &[]Statement{
					Statement{
						Effect: "deny",
						Action: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_POLICY, "/path/"),
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
			getGroupsByUserIDResult: []Group{
				Group{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/1/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"DenyResourceErr": {
			authUser: AuthenticatedUser{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "example",
			policyName: "policyUser",
			getPolicyByNameMethodResult: &Policy{
				ID:   "POLICY-USER-ID",
				Name: "policyUser",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
				Statements: &[]Statement{
					Statement{
						Effect: "deny",
						Action: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_POLICY, "/path/"),
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
			getGroupsByUserIDResult: []Group{
				Group{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/1/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getPoliciesAttachedResult: []Policy{
				Policy{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
						Statement{
							Effect: "deny",
							Action: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
							},
						},
					},
				},
			},
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		// this has to be fixed (issue #68)
		//"BadOrgName": {
		//	authUser: AuthenticatedUser{
		//		Identifier: "123456",
		//		Admin:      true,
		//	},
		//	org:        "123~#**!",
		//	policyName: "test",
		//	wantError: &Error{
		//		Code: INVALID_PARAMETER_ERROR,
		//	},
		//},
	}

	for x, testcase := range testcases {

		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetPolicyByNameMethod][0] = testcase.getPolicyByNameMethodResult
		testRepo.ArgsOut[GetPolicyByNameMethod][1] = testcase.getPolicyByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetPoliciesAttachedMethod][0] = testcase.getPoliciesAttachedResult
		policy, err := testAPI.GetPolicyByName(testcase.authUser, testcase.org, testcase.policyName)
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
				if !reflect.DeepEqual(policy, testcase.getPolicyByNameMethodResult) {
					t.Fatalf("Test %v failed. Received different policies (wanted:%v / received:%v)",
						x, testcase.getPolicyByNameMethodResult, policy)
				}
			}
		}
	}
}

func TestAddPolicy(t *testing.T) {
	testcases := map[string]struct {
		authUser   AuthenticatedUser
		org        string
		policyName string
		path       string
		statements []Statement

		getGroupsByUserIDResult   []Group
		getPoliciesAttachedResult []Policy
		getUserByExternalIDResult *User

		addPolicyMethodResult       *Policy
		getPolicyByNameMethodResult *Policy
		wantError                   *Error

		getPolicyByNameMethodErr error
		addPolicyMethodErr       error
	}{
		"OKCase": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			getPolicyByNameMethodErr: &database.Error{
				Code: database.POLICY_NOT_FOUND,
			},
			addPolicyMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
		"PolicyAlreadyExists": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
			wantError: &Error{
				Code: POLICY_ALREADY_EXIST,
			},
		},
		"BadName": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "**!^#~",
			path:       "/path/",
			statements: []Statement{
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
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"BadPath": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/**!^#~path/",
			statements: []Statement{
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
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"BadStatement": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
				Statement{
					Effect: "idufhefmfcasfluhf",
					Action: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"NoPermissions": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"DenyResource": {
			authUser: AuthenticatedUser{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "example",
			policyName: "test",
			path:       "/path/",
			getPolicyByNameMethodErr: &database.Error{
				Code: database.POLICY_NOT_FOUND,
			},
			statements: []Statement{
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
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getPoliciesAttachedResult: []Policy{
				Policy{
					ID:   "POLICY-USER-ID",
					Name: "policy",
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/"),
							},
						},
						Statement{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_CREATE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/"),
							},
						},
						Statement{
							Effect: "deny",
							Action: []string{
								POLICY_ACTION_CREATE_POLICY,
							},
							Resources: []string{
								CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
							},
						},
					},
				},
			},
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"AddPolicyDBErr": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			getPolicyByNameMethodErr: &database.Error{
				Code: database.POLICY_NOT_FOUND,
			},
			addPolicyMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
		"GetPolicyDBErr": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			getPolicyByNameMethodErr: &database.Error{
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
		testRepo.ArgsOut[AddPolicyMethod][0] = testcase.addPolicyMethodResult
		testRepo.ArgsOut[AddPolicyMethod][1] = testcase.addPolicyMethodErr
		testRepo.ArgsOut[GetPolicyByNameMethod][0] = testcase.getPolicyByNameMethodResult
		testRepo.ArgsOut[GetPolicyByNameMethod][1] = testcase.getPolicyByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetPoliciesAttachedMethod][0] = testcase.getPoliciesAttachedResult
		policy, err := testAPI.AddPolicy(testcase.authUser, testcase.policyName, testcase.path, testcase.org, testcase.statements)
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
				if !reflect.DeepEqual(policy, testcase.addPolicyMethodResult) {
					t.Fatalf("Test %v failed. Received different policies (wanted:%v / received:%v)",
						x, testcase.addPolicyMethodResult, policy)
				}
			}
		}
	}
}

func TestUpdatePolicy(t *testing.T) {
	testcases := map[string]struct {
		authUser      AuthenticatedUser
		org           string
		policyName    string
		path          string
		newPolicyName string
		newPath       string
		statements    []Statement
		newStatements []Statement

		getPolicyByNameMethodResult *Policy
		getGroupsByUserIDResult     []Group
		getPoliciesAttachedResult   []Policy
		getUserByExternalIDResult   *User
		updatePolicyMethodResult    *Policy

		wantError *Error

		getPolicyByNameMethodErr error
		getUserByExternalIDErr   error
		updatePolicyMethodErr    error

		getPolicyByNameMethodSpecialFunc func(string, string) (*Policy, error)
	}{
		"OKCase": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				Statement{
					Effect: "allow",
					Action: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
			updatePolicyMethodResult: &Policy{
				ID:   "test2",
				Name: "test2",
				Org:  "123",
				Path: "/path2/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path2/", "test"),
				Statements: &[]Statement{
					Statement{
						Effect: "allow",
						Action: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_USER, "/path2/"),
						},
					},
				},
			},
		},
		"InvalidPolicyName": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "**!~#",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				Statement{
					Effect: "allow",
					Action: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"InvalidNewPolicyName": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "**!~#",
			newPath:       "/path2/",
			newStatements: []Statement{
				Statement{
					Effect: "allow",
					Action: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"InvalidNewPath": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/**~#!/",
			newStatements: []Statement{
				Statement{
					Effect: "allow",
					Action: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"InvalidNewStatements": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				Statement{
					Effect: "jblkasdjgp",
					Action: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"GetPolicyDBErr": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				Statement{
					Effect: "allow",
					Action: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
		"PolicyNotFound": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				Statement{
					Effect: "allow",
					Action: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodErr: &database.Error{
				Code: database.POLICY_NOT_FOUND,
			},
			wantError: &Error{
				Code: POLICY_BY_ORG_AND_NAME_NOT_FOUND,
			},
		},
		"AuthUserNotFound": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				Statement{
					Effect: "allow",
					Action: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
			updatePolicyMethodResult: &Policy{
				ID:   "test2",
				Name: "test2",
				Org:  "123",
				Path: "/path2/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path2/", "test"),
				Statements: &[]Statement{
					Statement{
						Effect: "allow",
						Action: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_USER, "/path2/"),
						},
					},
				},
			},
			getUserByExternalIDErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"DenyResource": {
			authUser: AuthenticatedUser{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				Statement{
					Effect: "allow",
					Action: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
						Statement{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_UPDATE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
						Statement{
							Effect: "deny",
							Action: []string{
								POLICY_ACTION_UPDATE_POLICY,
							},
							Resources: []string{
								CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
							},
						},
					},
				},
			},
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"NoPermissions": {
			authUser: AuthenticatedUser{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				Statement{
					Effect: "allow",
					Action: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
			},
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"NewPolicyAlreadyExists": {
			authUser: AuthenticatedUser{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				Statement{
					Effect: "allow",
					Action: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodSpecialFunc: func(org string, name string) (*Policy, error) {
				if org == "123" && name == "test" {
					return &Policy{
						ID:   "test1",
						Name: "test",
						Org:  "123",
						Path: "/path/",
						Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
					}, nil
				} else {
					return &Policy{
						ID:   "test2",
						Name: "test2",
						Org:  "123",
						Path: "/path/",
						Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
					}, nil
				}
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
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
						Statement{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_UPDATE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
			},
			wantError: &Error{
				Code: POLICY_ALREADY_EXIST,
			},
		},
		"NoPermissionsToRetrieveTarget": {
			authUser: AuthenticatedUser{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				Statement{
					Effect: "allow",
					Action: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodSpecialFunc: func(org string, name string) (*Policy, error) {
				if org == "123" && name == "test" {
					return &Policy{
						ID:   "test1",
						Name: "test",
						Org:  "123",
						Path: "/path/",
						Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
					}, nil
				} else {
					return &Policy{
						ID:   "test2",
						Name: "test2",
						Org:  "123",
						Path: "/path2/",
						Urn:  CreateUrn("123", RESOURCE_POLICY, "/path2/", "test"),
						Statements: &[]Statement{
							Statement{
								Effect: "allow",
								Action: []string{
									USER_ACTION_GET_USER,
								},
								Resources: []string{
									GetUrnPrefix("", RESOURCE_USER, "/path2/"),
								},
							},
						},
					}, nil
				}
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
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
				Policy{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_UPDATE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
			},
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"NoPermissionsToUpdateTarget": {
			authUser: AuthenticatedUser{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				Statement{
					Effect: "allow",
					Action: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodSpecialFunc: func(org string, name string) (*Policy, error) {
				if org == "123" && name == "test" {
					return &Policy{
						ID:   "test1",
						Name: "test",
						Org:  "123",
						Path: "/path/",
						Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
					}, nil
				} else {
					return nil, &database.Error{
						Code: database.POLICY_NOT_FOUND,
					}
				}
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
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
				Policy{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_UPDATE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
			},
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"ExplicitDenyPermissionsToUpdateTarget": {
			authUser: AuthenticatedUser{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				Statement{
					Effect: "allow",
					Action: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodSpecialFunc: func(org string, name string) (*Policy, error) {
				if org == "123" && name == "test" {
					return &Policy{
						ID:   "test1",
						Name: "test",
						Org:  "123",
						Path: "/path/",
						Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
					}, nil
				} else {
					return nil, &database.Error{
						Code: database.POLICY_NOT_FOUND,
					}
				}
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
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
				Policy{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_UPDATE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
				Policy{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path2/"),
							},
						},
					},
				},
				Policy{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_UPDATE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path2/"),
							},
						},
					},
				},
				Policy{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "deny",
							Action: []string{
								POLICY_ACTION_UPDATE_POLICY,
							},
							Resources: []string{
								CreateUrn("123", RESOURCE_POLICY, "/path2/", "test2"),
							},
						},
					},
				},
			},
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"ErrorUpdatingPolicy": {
			authUser: AuthenticatedUser{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				Statement{
					Effect: "allow",
					Action: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodSpecialFunc: func(org string, name string) (*Policy, error) {
				if org == "123" && name == "test" {
					return &Policy{
						ID:   "test1",
						Name: "test",
						Org:  "123",
						Path: "/path/",
						Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
					}, nil
				} else {
					return nil, &database.Error{
						Code: database.POLICY_NOT_FOUND,
					}
				}
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
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
				Policy{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_UPDATE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
				Policy{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path2/"),
							},
						},
					},
				},
				Policy{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_UPDATE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path2/"),
							},
						},
					},
				},
			},
			updatePolicyMethodErr: &database.Error{
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
		testRepo.ArgsOut[UpdatePolicyMethod][0] = testcase.updatePolicyMethodResult
		testRepo.ArgsOut[UpdatePolicyMethod][1] = testcase.updatePolicyMethodErr
		testRepo.ArgsOut[GetPolicyByNameMethod][0] = testcase.getPolicyByNameMethodResult
		testRepo.ArgsOut[GetPolicyByNameMethod][1] = testcase.getPolicyByNameMethodErr
		testRepo.SpecialFuncs[GetPolicyByNameMethod] = testcase.getPolicyByNameMethodSpecialFunc
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetPoliciesAttachedMethod][0] = testcase.getPoliciesAttachedResult
		policy, err := testAPI.UpdatePolicy(testcase.authUser, testcase.org, testcase.policyName, testcase.newPolicyName, testcase.newPath, testcase.newStatements)
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
				if !reflect.DeepEqual(policy, testcase.updatePolicyMethodResult) {
					t.Fatalf("Test %v failed. Received different policies (wanted:%v / received:%v)",
						x, testcase.updatePolicyMethodResult, policy)
				}
			}
		}
	}
}

func TestGetListPolicies(t *testing.T) {
	testcases := map[string]struct {
		authUser   AuthenticatedUser
		org        string
		pathPrefix string

		expectedPolicies []PolicyIdentity

		getGroupsByUserIDResult   []Group
		getPoliciesAttachedResult []Policy
		getUserByExternalIDResult *User
		getUserByExternalIDErr    error

		getPoliciesFilteredMethodResult []Policy
		getPoliciesFilteredMethodErr    error

		wantError *Error
	}{
		"ErrorCaseInvalidPath": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			pathPrefix: "/path*/",
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"ErrorCaseInternalErrorGetPoliciesFiltered": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			pathPrefix: "/path/",
			getPoliciesFilteredMethodErr: &database.Error{
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
			org:        "123",
			pathPrefix: "/path/",
			getPoliciesFilteredMethodResult: []Policy{
				Policy{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
			},
			getUserByExternalIDErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"OkTestCaseAdmin": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			pathPrefix: "/",
			expectedPolicies: []PolicyIdentity{
				PolicyIdentity{
					Org:  "example",
					Name: "policyAllowed",
				},
				PolicyIdentity{
					Org:  "example",
					Name: "policyDenied",
				},
			},
			getPoliciesFilteredMethodResult: []Policy{
				Policy{
					ID:   "PolicyAllowed",
					Name: "policyAllowed",
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyAllowed"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
				Policy{
					ID:   "PolicyDenied",
					Name: "policyDenied",
					Org:  "example",
					Path: "/path2/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path2/", "policyDenied"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
			},
		},
		"OkTestCaseUser": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			org: "example",
			expectedPolicies: []PolicyIdentity{
				PolicyIdentity{
					Org:  "example",
					Name: "policyAllowed",
				},
			},
			getPoliciesFilteredMethodResult: []Policy{
				Policy{
					ID:   "PolicyAllowed",
					Name: "policyAllowed",
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyAllowed"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
				Policy{
					ID:   "PolicyDenied",
					Name: "policyDenied",
					Org:  "example",
					Path: "/path2/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path2/", "policyDenied"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
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
			getGroupsByUserIDResult: []Group{
				Group{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/1/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getPoliciesAttachedResult: []Policy{
				Policy{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_LIST_POLICIES,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
						Statement{
							Effect: "deny",
							Action: []string{
								POLICY_ACTION_LIST_POLICIES,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path2/"),
							},
						},
					},
				},
			},
		},
		// this has to be fixed (issue #68)
		//"BadOrgName": {
		//	authUser: AuthenticatedUser{
		//		Identifier: "123456",
		//		Admin:      true,
		//	},
		//	org:        "123~#**!",
		//	policyName: "test",
		//	wantError: &Error{
		//		Code: INVALID_PARAMETER_ERROR,
		//	},
		//},
	}

	for x, testcase := range testcases {

		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetPoliciesFilteredMethod][0] = testcase.getPoliciesFilteredMethodResult
		testRepo.ArgsOut[GetPoliciesFilteredMethod][1] = testcase.getPoliciesFilteredMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetPoliciesAttachedMethod][0] = testcase.getPoliciesAttachedResult
		policies, err := testAPI.GetListPolicies(testcase.authUser, testcase.org, testcase.pathPrefix)
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
				if !reflect.DeepEqual(policies, testcase.expectedPolicies) {
					t.Fatalf("Test %v failed. Received different policies (wanted:%v / received:%v)",
						x, testcase.expectedPolicies, policies)
				}
			}
		}
	}
}
