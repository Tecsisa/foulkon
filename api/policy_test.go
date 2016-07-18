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
		getAttachedPoliciesResult []Policy
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
		"BadOrgName": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "~#**!",
			policyName: "p1",
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
					{
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
				{
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
					{
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
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
						{
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
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		policy, err := testAPI.GetPolicyByName(testcase.authUser, testcase.org, testcase.policyName)
		if testcase.wantError != nil {
			apiError, ok := err.(*Error)
			if !ok || apiError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", x, err)
				continue
			}
			if apiError.Code != testcase.wantError.Code {
				t.Errorf("Test %v failed. Got error %v, expected %v", x, apiError, testcase.wantError.Code)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Error: %v", x, err)
				continue
			} else {
				if !reflect.DeepEqual(policy, testcase.getPolicyByNameMethodResult) {
					t.Errorf("Test %v failed. Received different policies (wanted:%v / received:%v)",
						x, testcase.getPolicyByNameMethodResult, policy)
					continue
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
		getAttachedPoliciesResult []Policy
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
		"PolicyAlreadyExists": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"BadOrgName": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "**!^#~",
			policyName: "p1",
			path:       "/path/",
			statements: []Statement{
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
				{
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
					Name: "policy",
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/"),
							},
						},
						{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_CREATE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/"),
							},
						},
						{
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
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		policy, err := testAPI.AddPolicy(testcase.authUser, testcase.policyName, testcase.path, testcase.org, testcase.statements)
		if testcase.wantError != nil {
			apiError, ok := err.(*Error)
			if !ok || apiError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", x, err)
				continue
			}
			if apiError.Code != testcase.wantError.Code {
				t.Errorf("Test %v failed. Got error %v, expected %v", x, apiError, testcase.wantError.Code)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Error: %v", x, err)
				continue
			} else {
				if !reflect.DeepEqual(policy, testcase.addPolicyMethodResult) {
					t.Errorf("Test %v failed. Received different policies (wanted:%v / received:%v)",
						x, testcase.addPolicyMethodResult, policy)
					continue
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
		getAttachedPoliciesResult   []Policy
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
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
			updatePolicyMethodResult: &Policy{
				ID:   "test2",
				Name: "test2",
				Org:  "123",
				Path: "/path2/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path2/", "test"),
				Statements: &[]Statement{
					{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
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
		"InvalidOrgName": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "**!~#",
			policyName: "p1",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
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
			newPolicyName: "**!~#",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
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
			newPolicyName: "test2",
			newPath:       "/**~#!/",
			newStatements: []Statement{
				{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
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
			updatePolicyMethodResult: &Policy{
				ID:   "test2",
				Name: "test2",
				Org:  "123",
				Path: "/path2/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path2/", "test"),
				Statements: &[]Statement{
					{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
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
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
						{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_UPDATE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
						{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
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
					}, nil
				} else {
					return &Policy{
						ID:   "test2",
						Name: "test2",
						Org:  "123",
						Path: "/path/",
						Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
						{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
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
					}, nil
				} else {
					return &Policy{
						ID:   "test2",
						Name: "test2",
						Org:  "123",
						Path: "/path2/",
						Urn:  CreateUrn("123", RESOURCE_POLICY, "/path2/", "test"),
						Statements: &[]Statement{
							{
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
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
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
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
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
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
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
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
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
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
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
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
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
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
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
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
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
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
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
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		policy, err := testAPI.UpdatePolicy(testcase.authUser, testcase.org, testcase.policyName, testcase.newPolicyName, testcase.newPath, testcase.newStatements)
		if testcase.wantError != nil {
			apiError, ok := err.(*Error)
			if !ok || apiError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", x, err)
				continue
			}
			if apiError.Code != testcase.wantError.Code {
				t.Errorf("Test %v failed. Got error %v, expected %v", x, apiError, testcase.wantError.Code)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Error: %v", x, err)
				continue
			} else {
				if !reflect.DeepEqual(policy, testcase.updatePolicyMethodResult) {
					t.Errorf("Test %v failed. Received different policies (wanted:%v / received:%v)",
						x, testcase.updatePolicyMethodResult, policy)
					continue
				}
			}
		}
	}
}

func TestAuthAPI_GetPolicyList(t *testing.T) {
	testcases := map[string]struct {
		authUser   AuthenticatedUser
		org        string
		pathPrefix string

		expectedPolicies []PolicyIdentity

		getGroupsByUserIDResult   []Group
		getAttachedPoliciesResult []Policy
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
		"ErrorCaseInvalidOrg": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "!#$$%**^",
			pathPrefix: "/",
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"ErrorCaseInternalErrorGetPoliciesFiltered": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "",
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
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
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
				{
					Org:  "example",
					Name: "policyAllowed",
				},
				{
					Org:  "example",
					Name: "policyDenied",
				},
			},
			getPoliciesFilteredMethodResult: []Policy{
				{
					ID:   "PolicyAllowed",
					Name: "policyAllowed",
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyAllowed"),
					Statements: &[]Statement{
						{
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
				{
					ID:   "PolicyDenied",
					Name: "policyDenied",
					Org:  "example",
					Path: "/path2/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path2/", "policyDenied"),
					Statements: &[]Statement{
						{
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
		"OkTestCaseAdminNoOrg": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "",
			pathPrefix: "/",
			expectedPolicies: []PolicyIdentity{
				{
					Org:  "example",
					Name: "policyAllowed",
				},
			},
			getPoliciesFilteredMethodResult: []Policy{
				{
					ID:   "PolicyAllowed",
					Name: "policyAllowed",
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyAllowed"),
					Statements: &[]Statement{
						{
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
				{
					Org:  "example",
					Name: "policyAllowed",
				},
			},
			getPoliciesFilteredMethodResult: []Policy{
				{
					ID:   "PolicyAllowed",
					Name: "policyAllowed",
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyAllowed"),
					Statements: &[]Statement{
						{
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
				{
					ID:   "PolicyDenied",
					Name: "policyDenied",
					Org:  "example",
					Path: "/path2/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path2/", "policyDenied"),
					Statements: &[]Statement{
						{
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
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_LIST_POLICIES,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
						{
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
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		policies, err := testAPI.GetPolicyList(testcase.authUser, testcase.org, testcase.pathPrefix)
		if testcase.wantError != nil {
			apiError, ok := err.(*Error)
			if !ok || apiError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", x, err)
				continue
			}
			if apiError.Code != testcase.wantError.Code {
				t.Errorf("Test %v failed. Got error %v, expected %v", x, apiError, testcase.wantError.Code)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Error: %v", x, err)
				continue
			} else {
				if !reflect.DeepEqual(policies, testcase.expectedPolicies) {
					t.Errorf("Test %v failed. Received different policies (wanted:%v / received:%v)",
						x, testcase.expectedPolicies, policies)
					continue
				}
			}
		}
	}
}

func TestDeletePolicy(t *testing.T) {
	testcases := map[string]struct {
		authUser AuthenticatedUser
		org      string
		name     string

		getPolicyByNameMethodResult *Policy
		getPolicyByNameMethodErr    error
		getGroupsByUserIDResult     []Group
		getAttachedPoliciesResult   []Policy
		getUserByExternalIDResult   *User
		getUserByExternalIDErr      error
		deletePolicyErr             error

		wantError *Error
	}{
		"ErrorCaseInvalidName": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:  "123",
			name: "invalid*",
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"ErrorCaseInvalidOrg": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:  "**!^#$%",
			name: "invalid",
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"ErrorCasePolicyNotExist": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:  "123",
			name: "policy",
			wantError: &Error{
				Code: POLICY_BY_ORG_AND_NAME_NOT_FOUND,
			},
			getPolicyByNameMethodErr: &database.Error{
				Code: database.POLICY_NOT_FOUND,
			},
		},
		"ErrorCaseNoPermissions": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			org:  "example",
			name: "test",
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
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
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "123456",
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Statements: &[]Statement{
						{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseNotEnoughPermissions": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			org:  "example",
			name: "test",
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
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
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "123456",
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Statements: &[]Statement{
						{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/"),
							},
						},
						{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_DELETE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/"),
							},
						},
						{
							Effect: "deny",
							Action: []string{
								POLICY_ACTION_DELETE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseRemoveFail": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:  "example",
			name: "test",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
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
			deletePolicyErr: &database.Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
		"OkTestCase": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:  "example",
			name: "test",
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
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

		testRepo.ArgsOut[RemovePolicyMethod][0] = testcase.deletePolicyErr
		testRepo.ArgsOut[GetPolicyByNameMethod][0] = testcase.getPolicyByNameMethodResult
		testRepo.ArgsOut[GetPolicyByNameMethod][1] = testcase.getPolicyByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		err := testAPI.DeletePolicy(testcase.authUser, testcase.org, testcase.name)
		if testcase.wantError != nil {
			apiError, ok := err.(*Error)
			if !ok || apiError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", x, err)
				continue
			}
			if apiError.Code != testcase.wantError.Code {
				t.Errorf("Test %v failed. Got error %v, expected %v", x, apiError, testcase.wantError.Code)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Error: %v", x, err)
				continue
			}
		}
	}
}

func TestAuthAPI_GetAttachedGroups(t *testing.T) {
	testcases := map[string]struct {
		authUser       AuthenticatedUser
		org            string
		policyName     string
		expectedGroups []string

		getGroupsByUserIDResult   []Group
		getAttachedPoliciesResult []Policy
		getUserByExternalIDResult *User

		getAttachedGroupsResult []Group
		getAttachedGroupsErr    error

		getPolicyByNameMethodResult *Policy
		wantError                   *Error

		getPolicyByNameMethodErr error
	}{
		"ErrorCaseInvalidName": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "invalid*",
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"ErrorCaseInvalidOrg": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "!*^**~$%",
			policyName: "p1",
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"ErrorCasePolicyNotExist": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "policy",
			wantError: &Error{
				Code: POLICY_BY_ORG_AND_NAME_NOT_FOUND,
			},
			getPolicyByNameMethodErr: &database.Error{
				Code: database.POLICY_NOT_FOUND,
			},
		},
		"ErrorCaseNoPermissions": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "example",
			policyName: "test",
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
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
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "123456",
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Statements: &[]Statement{
						{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseNotEnoughPermissions": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "example",
			policyName: "test",
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
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
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "123456",
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Statements: &[]Statement{
						{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/"),
							},
						},
						{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_LIST_ATTACHED_GROUPS,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/"),
							},
						},
						{
							Effect: "deny",
							Action: []string{
								POLICY_ACTION_LIST_ATTACHED_GROUPS,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseGetAttachedPoliciesFail": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "example",
			policyName: "test",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
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
			getAttachedGroupsErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
		"OkTestCase": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "example",
			policyName: "test",
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
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
			getAttachedGroupsResult: []Group{
				{
					ID:   "Group1",
					Org:  "org1",
					Name: "group1",
				},
				{
					ID:   "Group2",
					Org:  "org2",
					Name: "group2",
				},
			},
			expectedGroups: []string{"group1", "group2"},
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
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		testRepo.ArgsOut[GetAttachedGroupsMethod][0] = testcase.getAttachedGroupsResult
		testRepo.ArgsOut[GetAttachedGroupsMethod][1] = testcase.getAttachedGroupsErr
		groups, err := testAPI.GetAttachedGroups(testcase.authUser, testcase.org, testcase.policyName)
		if testcase.wantError != nil {
			apiError, ok := err.(*Error)
			if !ok || apiError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", x, err)
				continue
			}
			if apiError.Code != testcase.wantError.Code {
				t.Errorf("Test %v failed. Got error %v, expected %v", x, apiError, testcase.wantError.Code)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Error: %v", x, err)
				continue
			} else {
				if !reflect.DeepEqual(groups, testcase.expectedGroups) {
					t.Errorf("Test %v failed. Received different groups (wanted:%v / received:%v)",
						x, testcase.expectedGroups, groups)
					continue
				}
			}
		}
	}
}
