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

		expectedPolicy *Policy
		wantError      *Error

		GetPolicyByNameMethodErr error
	}{
		"OKCase": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			expectedPolicy: &Policy{
				ID:   "test1",
				Name: "test",
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
			GetPolicyByNameMethodErr: &database.Error{
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
			GetPolicyByNameMethodErr: &database.Error{
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
			policyName: "POLICY-USER-ID",
			expectedPolicy: &Policy{
				ID:   "POLICY-USER-ID",
				Name: "policyUser",
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
			policyName: "POLICY-USER-ID",
			expectedPolicy: &Policy{
				ID:   "POLICY-USER-ID",
				Name: "policyUser",
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

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for x, testcase := range testcases {
		testRepo.ArgsOut[GetPolicyByNameMethod][0] = testcase.expectedPolicy
		testRepo.ArgsOut[GetPolicyByNameMethod][1] = testcase.GetPolicyByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetPoliciesAttachedMethod][0] = testcase.getPoliciesAttachedResult
		policy, err := testAPI.GetPolicyByName(testcase.authUser, testcase.org, testcase.policyName)
		if testcase.wantError != nil {
			if errCode := err.(*Error).Code; errCode != testcase.wantError.Code {
				t.Fatalf("Test %v failed. Got error %v, expected %v", x, errCode, testcase.wantError.Code)
			}
		} else {
			if err != nil {
				t.Fatalf("Test %v failed", x)
			} else {
				if !reflect.DeepEqual(policy, testcase.expectedPolicy) {
					t.Fatalf("Test %v failed. Received different policies", x)
				}
			}
		}
	}
}
