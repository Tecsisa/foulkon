package api

import (
	"reflect"
	"testing"

	"github.com/tecsisa/authorizr/database"
)

func TestGetUsersAuthorized(t *testing.T) {
	testcases := map[string]struct {
		// Authenticated user
		authUser AuthenticatedUser
		// Resource urn that user wants to access
		resourceUrn string
		// Action to do
		action string
		// Resources received from db that system has to authorize
		usersToAuthorize []User
		// Resources authorized by method
		usersAuthorized []User
		// Error to compare when we expect an error
		wantError *Error
		// GetUserByExternalID Method Out Arguments
		getUserByExternalIDResult *User
		getUserByExternalIDError  error
	}{
		"OKtestCaseAdmin": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			resourceUrn: CreateUrn("", RESOURCE_USER, "/path/", "user1"),
			action:      USER_ACTION_GET_USER,
			usersToAuthorize: []User{
				User{
					ID:  "654321",
					Urn: CreateUrn("", RESOURCE_USER, "/path/", "user1"),
				},
			},
			usersAuthorized: []User{
				User{
					ID:  "654321",
					Urn: CreateUrn("", RESOURCE_USER, "/path/", "user1"),
				},
			},
		},
		"OKtestCaseAdminWithEmptyResources": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			resourceUrn:      CreateUrn("", RESOURCE_USER, "/path/", "user1"),
			action:           USER_ACTION_GET_USER,
			usersToAuthorize: []User{},
			usersAuthorized:  []User{},
		},
		"ErrortestCaseAuthenticatedUserNotExist": {
			authUser: AuthenticatedUser{
				Identifier: "notAdminUser",
				Admin:      false,
			},
			resourceUrn: CreateUrn("", RESOURCE_USER, "/path/", "user1"),
			action:      USER_ACTION_GET_USER,
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
			getUserByExternalIDResult: nil,
			getUserByExternalIDError: &database.Error{
				Code:    database.USER_NOT_FOUND,
				Message: "Error",
			},
		},
	}

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for n, test := range testcases {

		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = test.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = test.getUserByExternalIDError

		authorizedUsers, err := testAPI.GetUsersAuthorized(test.authUser, test.resourceUrn, test.action, test.usersToAuthorize)
		if test.wantError != nil {
			if apiError := err.(*Error); test.wantError.Code != apiError.Code {
				t.Fatalf("Test %v failed. Received different error codes (wanted:%v / received:%v)", n,
					test.wantError.Code, apiError.Code)
			}
		} else {
			if err != nil {
				t.Fatalf("Test %v failed. Error: %v", n, err)
			}

			if !test.authUser.Admin {
				// Check received authenticated user in method GetUserByExternalID
				if testRepo.ArgsIn[GetUserByExternalIDMethod][0] != test.authUser.Identifier {
					t.Fatalf("Test %v failed. Received different user identifiers (wanted:%v / received:%v)",
						n, test.authUser.Identifier, testRepo.ArgsIn[GetUserByExternalIDMethod][0])
				}
			}

			// Check result
			if !reflect.DeepEqual(authorizedUsers, test.usersAuthorized) {
				t.Fatalf("Test %v failed. Received different authorized users (wanted:%v / received:%v)",
					test.usersAuthorized, authorizedUsers)
			}

		}
	}
}

func TestGroupsAuthorized(t *testing.T) {
	testcases := map[string]struct {
		// Authenticated user
		authUser AuthenticatedUser
		// Resource urn that user wants to access
		resourceUrn string
		// Action to do
		action string
		// Resources received from db that system has to authorize
		groupsToAuthorize []Group
		// Resources authorized by method
		groupsAuthorized []Group
		// Error to compare when we expect an error
		wantError *Error
		// GetUserByExternalID Method Out Arguments
		getUserByExternalIDResult *User
		getUserByExternalIDError  error
	}{
		"OKtestCaseAdmin": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			resourceUrn: CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
			action:      GROUP_ACTION_GET_GROUP,
			groupsToAuthorize: []Group{
				Group{
					ID:  "654321",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
				},
			},
			groupsAuthorized: []Group{
				Group{
					ID:  "654321",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
				},
			},
		},
		"OKtestCaseAdminWithEmptyResources": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			resourceUrn:       CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
			action:            GROUP_ACTION_GET_GROUP,
			groupsToAuthorize: []Group{},
			groupsAuthorized:  []Group{},
		},
		"ErrortestCaseAuthenticatedUserNotExist": {
			authUser: AuthenticatedUser{
				Identifier: "notAdminUser",
				Admin:      false,
			},
			resourceUrn: CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
			action:      GROUP_ACTION_GET_GROUP,
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getUserByExternalIDResult: nil,
			getUserByExternalIDError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "Error",
			},
		},
	}

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for n, test := range testcases {

		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = test.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = test.getUserByExternalIDError

		authorizedGroups, err := testAPI.GetGroupsAuthorized(test.authUser, test.resourceUrn, test.action, test.groupsToAuthorize)
		if test.wantError != nil {
			if apiError := err.(*Error); test.wantError.Code != apiError.Code {
				t.Fatalf("Test %v failed. Received different error codes (wanted:%v / received:%v)", n,
					test.wantError.Code, apiError.Code)
			}
		} else {
			if err != nil {
				t.Fatalf("Test %v failed. Error: %v", n, err)
			}

			if !test.authUser.Admin {
				// Check received authenticated user in method GetUserByExternalID
				if testRepo.ArgsIn[GetUserByExternalIDMethod][0] != test.authUser.Identifier {
					t.Fatalf("Test %v failed. Received different user identifiers (wanted:%v / received:%v)",
						n, test.authUser.Identifier, testRepo.ArgsIn[GetUserByExternalIDMethod][0])
				}
			}

			// Check result
			if !reflect.DeepEqual(authorizedGroups, test.groupsAuthorized) {
				t.Fatalf("Test %v failed. Received different authorized groups (wanted:%v / received:%v)",
					test.groupsAuthorized, authorizedGroups)
			}
		}
	}
}

func TestGetPoliciesAuthorized(t *testing.T) {
	testcases := map[string]struct {
		// Authenticated user
		authUser AuthenticatedUser
		// Resource urn that user wants to access
		resourceUrn string
		// Action to do
		action string
		// Resources received from db that system has to authorize
		policiesToAuthorize []Policy
		// Resources authorized by method
		policiesAuthorized []Policy
		// Error to compare when we expect an error
		wantError *Error
		// GetUserByExternalID Method Out Arguments
		getUserByExternalIDResult *User
		getUserByExternalIDError  error
	}{
		"OKtestCaseAdmin": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			resourceUrn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
			action:      POLICY_ACTION_GET_POLICY,
			policiesToAuthorize: []Policy{
				Policy{
					ID:  "654321",
					Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
				},
			},
			policiesAuthorized: []Policy{
				Policy{
					ID:  "654321",
					Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
				},
			},
		},
		"OKtestCaseAdminWithEmptyResources": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			resourceUrn:         CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
			action:              POLICY_ACTION_GET_POLICY,
			policiesToAuthorize: []Policy{},
			policiesAuthorized:  []Policy{},
		},
		"ErrortestCaseDatabaseError": {
			authUser: AuthenticatedUser{
				Identifier: "USER-AUTHENTICATED",
				Admin:      false,
			},
			resourceUrn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
			action:      POLICY_ACTION_GET_POLICY,
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getUserByExternalIDResult: nil,
			getUserByExternalIDError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "Error",
			},
		},
	}

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for n, test := range testcases {

		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = test.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = test.getUserByExternalIDError

		authorizedPolicies, err := testAPI.GetPoliciesAuthorized(test.authUser, test.resourceUrn, test.action, test.policiesToAuthorize)
		if test.wantError != nil {
			if apiError := err.(*Error); test.wantError.Code != apiError.Code {
				t.Fatalf("Test %v failed. Received different error codes (wanted:%v / received:%v)", n,
					test.wantError.Code, apiError.Code)
			}
		} else {
			if err != nil {
				t.Fatalf("Test %v failed. Error: %v", n, err)
			}

			if !test.authUser.Admin {
				// Check received authenticated user in method GetUserByExternalID
				if testRepo.ArgsIn[GetUserByExternalIDMethod][0] != test.authUser.Identifier {
					t.Fatalf("Test %v failed. Received different user identifiers (wanted:%v / received:%v)",
						n, test.authUser.Identifier, testRepo.ArgsIn[GetUserByExternalIDMethod][0])
				}
			}

			// Check result
			if !reflect.DeepEqual(authorizedPolicies, test.policiesAuthorized) {
				t.Fatalf("Test %v failed. Received different authorized policies (wanted:%v / received:%v)",
					authorizedPolicies, test.policiesAuthorized)
			}
		}
	}
}

func TestGetEffectByUserActionResource(t *testing.T) {
	testcases := map[string]struct {
		// Authenticated user
		authUser AuthenticatedUser
		// Resource urn that user wants to access
		resourceUrn string
		// Action to do
		action string
		// Expected restrictions for the resource and action requested
		expectedEffectRestriction *EffectRestriction
		// Error to compare when we expect an error
		wantError *Error
		// GetUserByExternalID Method Out Arguments
		getUserByExternalIDResult *User
		getUserByExternalIDError  error
		// GetGroupsByUserID Method Out Arguments
		getGroupsByUserIDResult []Group
		getGroupsByUserIDError  error
		// GetPoliciesAttached Method Out Arguments
		getPoliciesAttachedResult []Policy
		getPoliciesAttachedError  error
	}{
		"ErrortestCaseInvalidAction": {
			action: "valid::Action",
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"ErrortestCaseInvalidResource": {
			action:      "product:DoSomething",
			resourceUrn: "urn:invalid/resource:resource",
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"ErrortestCaseGetRestrictions": {
			action:      "product:DoSomething",
			resourceUrn: "urn:ews:*",
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
			getUserByExternalIDError: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"OktestCaseFullUrnAllow": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			resourceUrn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
			action:      POLICY_ACTION_GET_POLICY,
			expectedEffectRestriction: &EffectRestriction{
				Effect: "allow",
			},
			getUserByExternalIDResult: &User{
				ID:  "123456",
				Urn: CreateUrn("", RESOURCE_USER, "/path/", "user1"),
			},
			getGroupsByUserIDResult: []Group{
				Group{
					ID:  "GROUP-USER-ID",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getPoliciesAttachedResult: []Policy{
				Policy{
					ID:  "POLICY-USER-ID",
					Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
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
		"OktestCaseFullUrnDeny": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			resourceUrn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
			action:      POLICY_ACTION_GET_POLICY,
			expectedEffectRestriction: &EffectRestriction{
				Effect: "deny",
			},
			getUserByExternalIDResult: &User{
				ID:  "123456",
				Urn: CreateUrn("", RESOURCE_USER, "/path/", "user1"),
			},
			getGroupsByUserIDResult: []Group{
				Group{
					ID:  "GROUP-USER-ID",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getPoliciesAttachedResult: []Policy{
				Policy{
					ID:  "POLICY-USER-ID",
					Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "deny",
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
		"OktestCaseWithRestrictions": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			resourceUrn: "urn:ews:product:instance:resource*",
			action:      "product:DoAction",
			expectedEffectRestriction: &EffectRestriction{
				Restrictions: &Restrictions{
					AllowedFullUrns: []string{
						"urn:ews:product:instance:resource/path1/resourceAllow",
						"urn:ews:product:instance:resource/path2/resourceAllow",
					},
					AllowedUrnPrefixes: []string{
						"urn:ews:product:instance:resource/path1/*",
						"urn:ews:product:instance:resource/path2/*",
					},
					DeniedFullUrns: []string{
						"urn:ews:product:instance:resource/path1/resourceDeny",
						"urn:ews:product:instance:resource/path2/resourceDeny",
					},
					DeniedUrnPrefixes: []string{
						"urn:ews:product:instance:resource/path3/*",
						"urn:ews:product:instance:resource/path4/*",
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:  "123456",
				Urn: CreateUrn("", RESOURCE_USER, "/path/", "user1"),
			},
			getGroupsByUserIDResult: []Group{
				Group{
					ID:  "GROUP-USER-ID",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getPoliciesAttachedResult: []Policy{
				Policy{
					ID:  "POLICY-USER-ID",
					Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								"product:DoAction",
							},
							Resources: []string{
								"urn:ews:product:instance:resource/path1/resourceAllow",
								"urn:ews:product:instance:resource/path2/resourceAllow",
								"urn:ews:product:instance:resource/path1/*",
								"urn:ews:product:instance:resource/path2/*",
							},
						},
						Statement{
							Effect: "deny",
							Action: []string{
								"product:DoAction",
							},
							Resources: []string{
								"urn:ews:product:instance:resource/path1/resourceDeny",
								"urn:ews:product:instance:resource/path2/resourceDeny",
								"urn:ews:product:instance:resource/path3/*",
								"urn:ews:product:instance:resource/path4/*",
							},
						},
					},
				},
			},
		},
	}

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for n, test := range testcases {

		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = test.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = test.getUserByExternalIDError

		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = test.getGroupsByUserIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][1] = test.getGroupsByUserIDError

		testRepo.ArgsOut[GetPoliciesAttachedMethod][0] = test.getPoliciesAttachedResult
		testRepo.ArgsOut[GetPoliciesAttachedMethod][1] = test.getPoliciesAttachedError

		effectRestriction, err := testAPI.GetEffectByUserActionResource(test.authUser, test.action, test.resourceUrn)
		if test.wantError != nil {
			if apiError := err.(*Error); test.wantError.Code != apiError.Code {
				t.Fatalf("Test %v failed. Received different error codes (wanted:%v / received:%v)", n,
					test.wantError.Code, apiError.Code)
			}
		} else {
			if err != nil {
				t.Fatalf("Test %v failed. Error: %v", n, err)
			}

			if !test.authUser.Admin {
				// Check received authenticated user in method GetUserByExternalID
				if testRepo.ArgsIn[GetUserByExternalIDMethod][0] != test.authUser.Identifier {
					t.Fatalf("Test %v failed. Received different user identifiers (wanted:%v / received:%v)",
						n, test.authUser.Identifier, testRepo.ArgsIn[GetUserByExternalIDMethod][0])
				}
			}

			// Check result
			if !reflect.DeepEqual(effectRestriction, test.expectedEffectRestriction) {
				t.Fatalf("Test %v failed. Received different effect restrictions (wanted:%v / received:%v)",
					n, test.expectedEffectRestriction, effectRestriction)
			}
		}
	}
}

// Test for aux methods of Authorizr

func TestGetAuthorizedResources(t *testing.T) {
	testcases := map[string]struct {
		// Authenticated user
		authUser AuthenticatedUser
		// Resource urn that user wants to access
		resourceUrn string
		// Action to do
		action string
		// Resources received from db that system has to authorize
		resourcesToAuthorize []Resource
		// Resources authorized by method
		resourcesAuthorized []Resource
		// Error to compare when we expect an error
		wantError *Error
		// GetUserByExternalID Method Out Arguments
		getUserByExternalIDResult *User
		getUserByExternalIDError  error
		// GetGroupsByUserID Method Out Arguments
		getGroupsByUserIDResult []Group
		getGroupsByUserIDError  error
		// GetPoliciesAttached Method Out Arguments
		getPoliciesAttachedResult []Policy
		getPoliciesAttachedError  error
	}{
		"OKtestCaseAdmin": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			resourceUrn: CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
			action:      GROUP_ACTION_GET_GROUP,
			resourcesToAuthorize: []Resource{
				Group{
					ID:  "654321",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
				},
			},
			resourcesAuthorized: []Resource{
				Group{
					ID:  "654321",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
				},
			},
		},
		"ErrortestCaseGetRestrictions": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			resourceUrn: CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
			action:      GROUP_ACTION_GET_GROUP,
			resourcesToAuthorize: []Resource{
				Group{
					ID:  "654321",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
				},
			},
			resourcesAuthorized: []Resource{
				Group{
					ID:  "654321",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
				},
			},
			getUserByExternalIDError: &database.Error{
				Code:    database.USER_NOT_FOUND,
				Message: "User not found",
			},
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"ErrortestCaseNotAllowedResources": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			resourceUrn: CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
			action:      GROUP_ACTION_GET_GROUP,
			resourcesToAuthorize: []Resource{
				Group{
					ID:  "654321",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
				},
			},
			resourcesAuthorized: []Resource{
				Group{
					ID:  "654321",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
				},
			},
			getUserByExternalIDResult: &User{
				ID:  "123456",
				Urn: CreateUrn("", RESOURCE_USER, "/path/", "user1"),
			},
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
		},
		"OKtestCaseResourcesFiltered": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			resourceUrn: GetUrnPrefix("example", RESOURCE_GROUP, "/path"),
			action:      GROUP_ACTION_GET_GROUP,
			resourcesToAuthorize: []Resource{
				Group{
					ID:  "654321",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
				},
				Group{
					ID:  "UNAUTHORIZED-GROUP-ID",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUnauthorized"),
				},
			},
			resourcesAuthorized: []Resource{
				Group{
					ID:  "654321",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
				},
			},
			getUserByExternalIDResult: &User{
				ID:  "123456",
				Urn: CreateUrn("", RESOURCE_USER, "/path/", "user1"),
			},
			getGroupsByUserIDResult: []Group{
				Group{
					ID:  "GROUP-USER-ID",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getPoliciesAttachedResult: []Policy{
				Policy{
					ID:  "POLICY-USER-ID",
					Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
							},
						},
					},
				},
			},
		},
		"OKtestCaseResourcesFilteredReturnEmpty": {
			// This test case checks if user has access to groups in /path2/ prefix, but there are groups
			// only in /path/, so we expect a empty slice of groups authorized
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			resourceUrn: GetUrnPrefix("example", RESOURCE_GROUP, "/path"),
			action:      GROUP_ACTION_GET_GROUP,
			resourcesToAuthorize: []Resource{
				Group{
					ID:  "654321",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
				},
				Group{
					ID:  "654322",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "group2"),
				},
			},
			resourcesAuthorized: []Resource{},
			getUserByExternalIDResult: &User{
				ID:  "123456",
				Urn: CreateUrn("", RESOURCE_USER, "/path/", "user1"),
			},
			getGroupsByUserIDResult: []Group{
				Group{
					ID:  "GROUP-USER-ID",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getPoliciesAttachedResult: []Policy{
				Policy{
					ID:  "POLICY-USER-ID",
					Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_GROUP, "/path2/"),
							},
						},
					},
				},
			},
		},
	}

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for n, test := range testcases {

		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = test.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = test.getUserByExternalIDError

		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = test.getGroupsByUserIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][1] = test.getGroupsByUserIDError

		testRepo.ArgsOut[GetPoliciesAttachedMethod][0] = test.getPoliciesAttachedResult
		testRepo.ArgsOut[GetPoliciesAttachedMethod][1] = test.getPoliciesAttachedError

		authorizedResources, err := testAPI.getAuthorizedResources(test.authUser, test.resourceUrn, test.action, test.resourcesToAuthorize)
		if test.wantError != nil {
			if apiError := err.(*Error); test.wantError.Code != apiError.Code {
				t.Fatalf("Test %v failed. Received different error codes (wanted:%v / received:%v)", n,
					test.wantError.Code, apiError.Code)
			}
		} else {
			if err != nil {
				t.Fatalf("Test %v failed. Error: %v", n, err)
			}

			if !test.authUser.Admin {
				// Check received authenticated user in method GetUserByExternalID
				if testRepo.ArgsIn[GetUserByExternalIDMethod][0] != test.authUser.Identifier {
					t.Fatalf("Test %v failed. Received different user identifiers (wanted:%v / received:%v)",
						n, test.authUser.Identifier, testRepo.ArgsIn[GetUserByExternalIDMethod][0])
				}
			}

			// Check result
			if !reflect.DeepEqual(authorizedResources, test.resourcesAuthorized) {
				t.Fatalf("Test %v failed. Received different authorized resources (wanted:%v / received:%v)",
					n, test.resourcesAuthorized, authorizedResources)
			}
		}
	}
}

func TestGetRestrictions(t *testing.T) {
	testcases := map[string]struct {
		// Authenticated user identifier
		authUserID string
		// Resource urn that user wants to access
		resourceUrn string
		// Action to do
		action string
		// Expected Restrictions
		expectedRestrictions *Restrictions
		// Error to compare when we expect an error
		wantError *Error
		// GetUserByExternalID Method Out Arguments
		getUserByExternalIDResult *User
		getUserByExternalIDError  error
		// GetGroupsByUserID Method Out Arguments
		getGroupsByUserIDResult []Group
		getGroupsByUserIDError  error
		// GetPoliciesAttached Method Out Arguments
		getPoliciesAttachedResult []Policy
		getPoliciesAttachedError  error
	}{
		"ErrortestCaseGetUserAuthenticatedNotFound": {
			authUserID:  "NotFound",
			resourceUrn: "urn:resource",
			action:      USER_ACTION_GET_USER,
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
			getUserByExternalIDError: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrortestCaseGetUserAuthenticatedInternalError": {
			authUserID:  "InternalError",
			resourceUrn: "urn:resource",
			action:      USER_ACTION_GET_USER,
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getUserByExternalIDError: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
		"ErrortestCaseGetGroupsError": {
			authUserID:  "InternalError",
			resourceUrn: "urn:resource",
			action:      USER_ACTION_GET_USER,
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getUserByExternalIDResult: &User{
				ID: "UserID",
			},
			getGroupsByUserIDError: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
		"ErrortestCaseGetPoliciesError": {
			authUserID:  "InternalError",
			resourceUrn: "urn:resource",
			action:      USER_ACTION_GET_USER,
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getUserByExternalIDResult: &User{
				ID: "UserID",
			},
			getGroupsByUserIDResult: []Group{
				Group{
					ID: "GroupID",
				},
			},
			getPoliciesAttachedError: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
		"OktestCaseEmptyRelationsFullUrn": {
			authUserID:  "AuthUserID",
			resourceUrn: CreateUrn("example", RESOURCE_GROUP, "/path/", "group"),
			action:      USER_ACTION_GET_USER,
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{},
			},
			getUserByExternalIDResult: &User{
				ID: "AuthUserID",
			},
		},
		"OktestCaseEmptyRelationsPrefixUrn": {
			authUserID:  "AuthUserID",
			resourceUrn: GetUrnPrefix("example", RESOURCE_GROUP, "/path/"),
			action:      USER_ACTION_GET_USER,
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{},
			},
			getUserByExternalIDResult: &User{
				ID: "AuthUserID",
			},
		},
		"OktestCaseFullUrn": {
			authUserID:  "AuthUserID",
			resourceUrn: CreateUrn("example", RESOURCE_GROUP, "/path1/", "groupAllow"),
			action:      GROUP_ACTION_GET_GROUP,
			expectedRestrictions: &Restrictions{
				AllowedFullUrns: []string{
					CreateUrn("example", RESOURCE_GROUP, "/path1/", "groupAllow"),
				},
				AllowedUrnPrefixes: []string{
					GetUrnPrefix("example", RESOURCE_GROUP, "/path1/"),
				},
				DeniedFullUrns:    []string{},
				DeniedUrnPrefixes: []string{},
			},
			getUserByExternalIDResult: &User{
				ID: "AuthUserID",
			},
			getGroupsByUserIDResult: []Group{
				Group{
					ID: "GROUP-USER-ID",
				},
			},
			getPoliciesAttachedResult: []Policy{
				Policy{
					ID:  "POLICY-USER-ID",
					Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								CreateUrn("example", RESOURCE_GROUP, "/path1/", "groupAllow"),
								CreateUrn("example", RESOURCE_GROUP, "/path2/", "groupAllow"),
								GetUrnPrefix("example", RESOURCE_GROUP, "/path1/"),
								GetUrnPrefix("example", RESOURCE_GROUP, "/path2/"),
							},
						},
						Statement{
							Effect: "deny",
							Action: []string{
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								CreateUrn("example", RESOURCE_GROUP, "/path1/", "groupDeny"),
								CreateUrn("example", RESOURCE_GROUP, "/path2/", "groupDeny"),
								GetUrnPrefix("example", RESOURCE_GROUP, "/path3/"),
								GetUrnPrefix("example", RESOURCE_GROUP, "/path4/"),
							},
						},
					},
				},
			},
		},
		"OkPrefixUrn": {
			authUserID:  "AuthUserID",
			resourceUrn: GetUrnPrefix("example", RESOURCE_GROUP, "/path"),
			action:      GROUP_ACTION_GET_GROUP,
			expectedRestrictions: &Restrictions{
				AllowedFullUrns: []string{
					CreateUrn("example", RESOURCE_GROUP, "/path1/", "groupAllow"),
					CreateUrn("example", RESOURCE_GROUP, "/path2/", "groupAllow"),
				},
				AllowedUrnPrefixes: []string{
					GetUrnPrefix("example", RESOURCE_GROUP, "/path1/"),
					GetUrnPrefix("example", RESOURCE_GROUP, "/path2/"),
				},
				DeniedFullUrns: []string{
					CreateUrn("example", RESOURCE_GROUP, "/path1/", "groupDeny"),
					CreateUrn("example", RESOURCE_GROUP, "/path2/", "groupDeny"),
				},
				DeniedUrnPrefixes: []string{
					GetUrnPrefix("example", RESOURCE_GROUP, "/path3/"),
					GetUrnPrefix("example", RESOURCE_GROUP, "/path4/"),
				},
			},
			getUserByExternalIDResult: &User{
				ID: "AuthUserID",
			},
			getGroupsByUserIDResult: []Group{
				Group{
					ID: "GROUP-USER-ID",
				},
			},
			getPoliciesAttachedResult: []Policy{
				Policy{
					ID:  "POLICY-USER-ID",
					Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						Statement{
							Effect: "allow",
							Action: []string{
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								CreateUrn("example", RESOURCE_GROUP, "/path1/", "groupAllow"),
								CreateUrn("example", RESOURCE_GROUP, "/path2/", "groupAllow"),
								GetUrnPrefix("example", RESOURCE_GROUP, "/path1/"),
								GetUrnPrefix("example", RESOURCE_GROUP, "/path2/"),
							},
						},
						Statement{
							Effect: "deny",
							Action: []string{
								GROUP_ACTION_GET_GROUP,
							},
							Resources: []string{
								CreateUrn("example", RESOURCE_GROUP, "/path1/", "groupDeny"),
								CreateUrn("example", RESOURCE_GROUP, "/path2/", "groupDeny"),
								GetUrnPrefix("example", RESOURCE_GROUP, "/path3/"),
								GetUrnPrefix("example", RESOURCE_GROUP, "/path4/"),
							},
						},
					},
				},
			}},
	}

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for n, test := range testcases {

		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = test.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = test.getUserByExternalIDError

		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = test.getGroupsByUserIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][1] = test.getGroupsByUserIDError

		testRepo.ArgsOut[GetPoliciesAttachedMethod][0] = test.getPoliciesAttachedResult
		testRepo.ArgsOut[GetPoliciesAttachedMethod][1] = test.getPoliciesAttachedError

		restrictions, err := testAPI.getRestrictions(test.authUserID, test.action, test.resourceUrn)
		if test.wantError != nil {
			if apiError := err.(*Error); test.wantError.Code != apiError.Code {
				t.Fatalf("Test %v failed. Received different error codes (wanted:%v / received:%v)", n,
					test.wantError.Code, apiError.Code)
			}
		} else {
			if err != nil {
				t.Fatalf("Test %v failed. Error: %v", n, err)
			}

			if testRepo.ArgsIn[GetUserByExternalIDMethod][0] != test.authUserID {
				t.Fatalf("Test %v failed. Received different user identifiers (wanted:%v / received:%v)",
					n, test.authUserID, testRepo.ArgsIn[GetUserByExternalIDMethod][0])
			}

			if param := testRepo.ArgsIn[GetGroupsByUserIDMethod][0]; test.authUserID != "" && param != test.authUserID {
				t.Fatalf("Test %v failed. Received different user identifiers (wanted:%v / received:%v)",
					n, test.authUserID, testRepo.ArgsIn[GetGroupsByUserIDMethod][0])
			}

			if param := testRepo.ArgsIn[GetPoliciesAttachedMethod][0]; test.getGroupsByUserIDResult != nil &&
				param != test.getGroupsByUserIDResult[0].ID {
				t.Fatalf("Test %v failed. Received different user identifiers (wanted:%v / received:%v)",
					n, test.authUserID, testRepo.ArgsIn[GetPoliciesAttachedMethod][0])
			}

			// Check result
			if !reflect.DeepEqual(test.expectedRestrictions, restrictions) {
				t.Fatalf("Test %v failed. Received different restrictions (wanted:%v / received:%v)",
					n, test.expectedRestrictions, restrictions)
			}
		}
	}
}

func TestGetGroupsByUser(t *testing.T) {

}

func TestGetPoliciesByGroups(t *testing.T) {

}

func TestGetStatementsByRequestedAction(t *testing.T) {

}

func TestCleanRepeatedRestrictions(t *testing.T) {

}

func TestIsActionContained(t *testing.T) {

}

func TestIsResourceContained(t *testing.T) {

}

func TestIsFullUrn(t *testing.T) {

}

func TestGetRestrictionsWhenResourceRequestedIsPrefix(t *testing.T) {

}

func TestGetRestrictionsWhenResourceRequestedIsFullUrn(t *testing.T) {

}

func TestFilterResources(t *testing.T) {

}

func TestIsAllowedResource(t *testing.T) {

}
