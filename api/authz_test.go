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
	testRepo := TestRepo{
		ArgsIn:  make(map[string][]interface{}),
		ArgsOut: make(map[string][]interface{}),
	}
	api := AuthAPI{
		UserRepo:   testRepo,
		GroupRepo:  testRepo,
		PolicyRepo: testRepo,
	}
	for n, test := range testcases {
		// Init resources for method GetUserByExternalID
		testRepo.ArgsIn[GetUserByExternalIDMethod] = make([]interface{}, 1)

		testRepo.ArgsOut[GetUserByExternalIDMethod] = make([]interface{}, 2)
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = test.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = test.getUserByExternalIDError

		authorizedUsers, err := api.GetUsersAuthorized(test.authUser, test.resourceUrn, test.action, test.usersToAuthorize)
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
						n, testRepo.ArgsIn[GetUserByExternalIDMethod][0], test.authUser.Identifier)
				}
			}

			// Check result
			if !reflect.DeepEqual(authorizedUsers, test.usersAuthorized) {
				t.Fatalf("Test %v failed. Received different authorized users (wanted:%v / received:%v)",
					authorizedUsers, test.usersAuthorized)
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
	testRepo := TestRepo{
		ArgsIn:  make(map[string][]interface{}),
		ArgsOut: make(map[string][]interface{}),
	}
	api := AuthAPI{
		UserRepo:   testRepo,
		GroupRepo:  testRepo,
		PolicyRepo: testRepo,
	}
	for n, test := range testcases {
		// Init resources for method GetUserByExternalID
		testRepo.ArgsIn[GetUserByExternalIDMethod] = make([]interface{}, 1)

		testRepo.ArgsOut[GetUserByExternalIDMethod] = make([]interface{}, 2)
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = test.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = test.getUserByExternalIDError

		authorizedGroups, err := api.GetGroupsAuthorized(test.authUser, test.resourceUrn, test.action, test.groupsToAuthorize)
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
						n, testRepo.ArgsIn[GetUserByExternalIDMethod][0], test.authUser.Identifier)
				}
			}

			// Check result
			if !reflect.DeepEqual(authorizedGroups, test.groupsAuthorized) {
				t.Fatalf("Test %v failed. Received different authorized groups (wanted:%v / received:%v)",
					authorizedGroups, test.groupsAuthorized)
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
	testRepo := TestRepo{
		ArgsIn:  make(map[string][]interface{}),
		ArgsOut: make(map[string][]interface{}),
	}
	api := AuthAPI{
		UserRepo:   testRepo,
		GroupRepo:  testRepo,
		PolicyRepo: testRepo,
	}
	for n, test := range testcases {
		// Init resources for method GetUserByExternalID
		testRepo.ArgsIn[GetUserByExternalIDMethod] = make([]interface{}, 1)

		testRepo.ArgsOut[GetUserByExternalIDMethod] = make([]interface{}, 2)
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = test.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = test.getUserByExternalIDError

		authorizedPolicies, err := api.GetPoliciesAuthorized(test.authUser, test.resourceUrn, test.action, test.policiesToAuthorize)
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
						n, testRepo.ArgsIn[GetUserByExternalIDMethod][0], test.authUser.Identifier)
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
	testRepo := TestRepo{
		ArgsIn:  make(map[string][]interface{}),
		ArgsOut: make(map[string][]interface{}),
	}
	api := AuthAPI{
		UserRepo:   testRepo,
		GroupRepo:  testRepo,
		PolicyRepo: testRepo,
	}
	for n, test := range testcases {
		// Init resources for method GetUserByExternalID
		testRepo.ArgsIn[GetUserByExternalIDMethod] = make([]interface{}, 1)

		testRepo.ArgsOut[GetUserByExternalIDMethod] = make([]interface{}, 2)
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = test.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = test.getUserByExternalIDError

		// Init resources for method GetGroupsByUserID
		testRepo.ArgsIn[GetGroupsByUserIDMethod] = make([]interface{}, 1)

		testRepo.ArgsOut[GetGroupsByUserIDMethod] = make([]interface{}, 2)
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = test.getGroupsByUserIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][1] = test.getGroupsByUserIDError

		// Init resources for method GetPoliciesAttached
		testRepo.ArgsIn[GetPoliciesAttachedMethod] = make([]interface{}, 1)

		testRepo.ArgsOut[GetPoliciesAttachedMethod] = make([]interface{}, 2)
		testRepo.ArgsOut[GetPoliciesAttachedMethod][0] = test.getPoliciesAttachedResult
		testRepo.ArgsOut[GetPoliciesAttachedMethod][1] = test.getPoliciesAttachedError

		effectRestriction, err := api.GetEffectByUserActionResource(test.authUser, test.action, test.resourceUrn)
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
						n, testRepo.ArgsIn[GetUserByExternalIDMethod][0], test.authUser.Identifier)
				}
			}

			// Check result
			if !reflect.DeepEqual(effectRestriction, test.expectedEffectRestriction) {
				t.Fatalf("Test %v failed. Received different effect restrictions (wanted:%v / received:%v)",
					n, effectRestriction, test.expectedEffectRestriction)
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
	testRepo := TestRepo{
		ArgsIn:  make(map[string][]interface{}),
		ArgsOut: make(map[string][]interface{}),
	}
	api := AuthAPI{
		UserRepo:   testRepo,
		GroupRepo:  testRepo,
		PolicyRepo: testRepo,
	}
	for n, test := range testcases {
		// Init resources for method GetUserByExternalID
		testRepo.ArgsIn[GetUserByExternalIDMethod] = make([]interface{}, 1)

		testRepo.ArgsOut[GetUserByExternalIDMethod] = make([]interface{}, 2)
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = test.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = test.getUserByExternalIDError

		// Init resources for method GetGroupsByUserID
		testRepo.ArgsIn[GetGroupsByUserIDMethod] = make([]interface{}, 1)

		testRepo.ArgsOut[GetGroupsByUserIDMethod] = make([]interface{}, 2)
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = test.getGroupsByUserIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][1] = test.getGroupsByUserIDError

		// Init resources for method GetPoliciesAttached
		testRepo.ArgsIn[GetPoliciesAttachedMethod] = make([]interface{}, 1)

		testRepo.ArgsOut[GetPoliciesAttachedMethod] = make([]interface{}, 2)
		testRepo.ArgsOut[GetPoliciesAttachedMethod][0] = test.getPoliciesAttachedResult
		testRepo.ArgsOut[GetPoliciesAttachedMethod][1] = test.getPoliciesAttachedError

		authorizedResources, err := api.getAuthorizedResources(test.authUser, test.resourceUrn, test.action, test.resourcesToAuthorize)
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
						n, testRepo.ArgsIn[GetUserByExternalIDMethod][0], test.authUser.Identifier)
				}
			}

			// Check result
			if !reflect.DeepEqual(authorizedResources, test.resourcesAuthorized) {
				t.Fatalf("Test %v failed. Received different authorized resources (wanted:%v / received:%v)",
					n, authorizedResources, test.resourcesAuthorized)
			}
		}
	}
}

func TestGetRestrictions(t *testing.T) {

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
