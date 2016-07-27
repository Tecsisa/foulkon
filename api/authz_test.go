package api

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/tecsisa/authorizr/database"
)

func TestGetAuthorizedUsers(t *testing.T) {
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
				{
					ID:  "654321",
					Urn: CreateUrn("", RESOURCE_USER, "/path/", "user1"),
				},
			},
			usersAuthorized: []User{
				{
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

		authorizedUsers, err := testAPI.GetAuthorizedUsers(test.authUser, test.resourceUrn, test.action, test.usersToAuthorize)
		if test.wantError != nil {
			if apiError := err.(*Error); test.wantError.Code != apiError.Code {
				t.Errorf("Test %v failed. Received different error codes (wanted:%v / received:%v)", n,
					test.wantError.Code, apiError.Code)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Error: %v", n, err)
				continue
			}

			if !test.authUser.Admin {
				// Check received authenticated user in method GetUserByExternalID
				if testRepo.ArgsIn[GetUserByExternalIDMethod][0] != test.authUser.Identifier {
					t.Errorf("Test %v failed. Received different user identifiers (wanted:%v / received:%v)",
						n, test.authUser.Identifier, testRepo.ArgsIn[GetUserByExternalIDMethod][0])
					continue
				}
			}

			// Check result
			if diff := pretty.Compare(authorizedUsers, test.usersAuthorized); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
		}
	}
}

func TestGetAuthorizedGroups(t *testing.T) {
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
				{
					ID:  "654321",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
				},
			},
			groupsAuthorized: []Group{
				{
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

		authorizedGroups, err := testAPI.GetAuthorizedGroups(test.authUser, test.resourceUrn, test.action, test.groupsToAuthorize)
		if test.wantError != nil {
			if apiError := err.(*Error); test.wantError.Code != apiError.Code {
				t.Errorf("Test %v failed. Received different error codes (wanted:%v / received:%v)", n,
					test.wantError.Code, apiError.Code)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Error: %v", n, err)
				continue
			}

			if !test.authUser.Admin {
				// Check received authenticated user in method GetUserByExternalID
				if testRepo.ArgsIn[GetUserByExternalIDMethod][0] != test.authUser.Identifier {
					t.Errorf("Test %v failed. Received different user identifiers (wanted:%v / received:%v)",
						n, test.authUser.Identifier, testRepo.ArgsIn[GetUserByExternalIDMethod][0])
					continue
				}
			}

			// Check result
			if diff := pretty.Compare(authorizedGroups, test.groupsAuthorized); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
		}
	}
}

func TestGetAuthorizedPolicies(t *testing.T) {
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
				{
					ID:  "654321",
					Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
				},
			},
			policiesAuthorized: []Policy{
				{
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

		authorizedPolicies, err := testAPI.GetAuthorizedPolicies(test.authUser, test.resourceUrn, test.action, test.policiesToAuthorize)
		if test.wantError != nil {
			if apiError := err.(*Error); test.wantError.Code != apiError.Code {
				t.Errorf("Test %v failed. Received different error codes (wanted:%v / received:%v)", n,
					test.wantError.Code, apiError.Code)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Error: %v", n, err)
				continue
			}

			if !test.authUser.Admin {
				// Check received authenticated user in method GetUserByExternalID
				if testRepo.ArgsIn[GetUserByExternalIDMethod][0] != test.authUser.Identifier {
					t.Errorf("Test %v failed. Received different user identifiers (wanted:%v / received:%v)",
						n, test.authUser.Identifier, testRepo.ArgsIn[GetUserByExternalIDMethod][0])
					continue
				}
			}

			// Check result
			if diff := pretty.Compare(authorizedPolicies, test.policiesAuthorized); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
		}
	}
}

func TestGetAuthorizedExternalResources(t *testing.T) {
	testcases := map[string]struct {
		// Authenticated user
		authUser AuthenticatedUser
		// Resource urns that user wants to access
		resourceUrns []string
		// Action to do
		action string
		// Expected allowed resources
		expectedResources []string
		// Error to compare when we expect an error
		wantError *Error
		// GetUserByExternalID Method Out Arguments
		getUserByExternalIDResult *User
		getUserByExternalIDError  error
		// GetGroupsByUserID Method Out Arguments
		getGroupsByUserIDResult []Group
		getGroupsByUserIDError  error
		// GetAttachedPolicies Method Out Arguments
		getAttachedPoliciesResult []Policy
		getAttachedPoliciesError  error
	}{
		"ErrortestCaseInvalidAction": {
			action: "valid::Action",
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"ErrortestCaseInvalidResource": {
			action: "product:DoSomething",
			resourceUrns: []string{
				"urn:invalid/resource:resource",
			},
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"ErrortestCaseInvalidResourceWithPrefix": {
			action: "product:DoSomething",
			resourceUrns: []string{
				"urn:*",
			},
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"ErrortestCaseEmptyResources": {
			action: "product:DoSomething",
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"ErrortestCaseGetRestrictions": {
			action: "product:DoSomething",
			resourceUrns: []string{
				CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
			},
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
			getUserByExternalIDError: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrortestCaseActionPrefix": {
			action: "product:DoPrefix*",
			resourceUrns: []string{
				CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
			},
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
		"OktestCaseFullUrnAllow": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			resourceUrns: []string{
				CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
				CreateUrn("example", RESOURCE_POLICY, "/path/", "policy2"),
				CreateUrn("example1", RESOURCE_POLICY, "/path/", "policy3"),
			},
			action: POLICY_ACTION_GET_POLICY,
			expectedResources: []string{
				CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
				CreateUrn("example", RESOURCE_POLICY, "/path/", "policy2"),
			},
			getUserByExternalIDResult: &User{
				ID:  "123456",
				Urn: CreateUrn("", RESOURCE_USER, "/path/", "user1"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:  "GROUP-USER-ID",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:  "POLICY-USER-ID",
					Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
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
		"OktestCaseFullUrnDeny": {
			authUser: AuthenticatedUser{
				Identifier: "123456",
				Admin:      false,
			},
			action: POLICY_ACTION_GET_POLICY,
			resourceUrns: []string{
				CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
				CreateUrn("example", RESOURCE_POLICY, "/path/", "policy2"),
				CreateUrn("example1", RESOURCE_POLICY, "/path/", "policy3"),
			},
			expectedResources: []string{},
			getUserByExternalIDResult: &User{
				ID:  "123456",
				Urn: CreateUrn("", RESOURCE_USER, "/path/", "user1"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:  "GROUP-USER-ID",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:  "POLICY-USER-ID",
					Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "deny",
							Action: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
						{
							Effect: "allow",
							Action: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/path2"),
								GetUrnPrefix("example2", RESOURCE_POLICY, "/path/path2"),
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
			resourceUrns: []string{
				"urn:ews:product:instance:resource/path1/resourceAllow",
				"urn:ews:product:instance:resource/path2/resourceAllow",
				"urn:ews:product:instance:resource/path1/resourceDeny",
				"urn:ews:product:instance:resource/path2/resourceDeny",
			},
			action: "product:DoAction",
			expectedResources: []string{
				"urn:ews:product:instance:resource/path1/resourceAllow",
				"urn:ews:product:instance:resource/path2/resourceAllow",
			},
			getUserByExternalIDResult: &User{
				ID:  "123456",
				Urn: CreateUrn("", RESOURCE_USER, "/path/", "user1"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:  "GROUP-USER-ID",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:  "POLICY-USER-ID",
					Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Action: []string{
								"product:DoAction",
							},
							Resources: []string{
								"urn:ews:product:instance:resource/path1/resourceAllow",
								"urn:ews:product:instance:resource/path2/resourceAllow",
								"urn:ews:product:instance:resource/path1*",
								"urn:ews:product:instance:resource/path2*",
							},
						},
						{
							Effect: "deny",
							Action: []string{
								"product:DoAction",
							},
							Resources: []string{
								"urn:ews:product:instance:resource/path1/resourceDeny",
								"urn:ews:product:instance:resource/path2/resourceDeny",
								"urn:ews:product:instance:resource/path3*",
								"urn:ews:product:instance:resource/path4*",
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

		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = test.getAttachedPoliciesResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][1] = test.getAttachedPoliciesError

		resources, err := testAPI.GetAuthorizedExternalResources(test.authUser, test.action, test.resourceUrns)
		if test.wantError != nil {
			if apiError := err.(*Error); test.wantError.Code != apiError.Code {
				t.Errorf("Test %v failed. Received different error codes (wanted:%v / received:%v)", n,
					test.wantError.Code, apiError.Code)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Error: %v", n, err)
				continue
			}

			if !test.authUser.Admin {
				// Check received authenticated user in method GetUserByExternalID
				if testRepo.ArgsIn[GetUserByExternalIDMethod][0] != test.authUser.Identifier {
					t.Errorf("Test %v failed. Received different user identifiers (wanted:%v / received:%v)",
						n, test.authUser.Identifier, testRepo.ArgsIn[GetUserByExternalIDMethod][0])
					continue
				}
			}

			// Check result
			if diff := pretty.Compare(resources, test.expectedResources); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
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
		// GetAttachedPolicies Method Out Arguments
		getAttachedPoliciesResult []Policy
		getAttachedPoliciesError  error
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
				{
					ID:  "GROUP-USER-ID",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:  "POLICY-USER-ID",
					Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
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
				{
					ID:  "GROUP-USER-ID",
					Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:  "POLICY-USER-ID",
					Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
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

		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = test.getAttachedPoliciesResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][1] = test.getAttachedPoliciesError

		authorizedResources, err := testAPI.getAuthorizedResources(test.authUser, test.resourceUrn, test.action, test.resourcesToAuthorize)
		if test.wantError != nil {
			if apiError := err.(*Error); test.wantError.Code != apiError.Code {
				t.Errorf("Test %v failed. Received different error codes (wanted:%v / received:%v)", n,
					test.wantError.Code, apiError.Code)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Error: %v", n, err)
				continue
			}

			if !test.authUser.Admin {
				// Check received authenticated user in method GetUserByExternalID
				if testRepo.ArgsIn[GetUserByExternalIDMethod][0] != test.authUser.Identifier {
					t.Errorf("Test %v failed. Received different user identifiers (wanted:%v / received:%v)",
						n, test.authUser.Identifier, testRepo.ArgsIn[GetUserByExternalIDMethod][0])
					continue
				}
			}

			// Check result
			if diff := pretty.Compare(authorizedResources, test.resourcesAuthorized); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
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
		// GetAttachedPolicies Method Out Arguments
		getAttachedPoliciesResult []Policy
		getAttachedPoliciesError  error
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
				{
					ID: "GroupID",
				},
			},
			getAttachedPoliciesError: &database.Error{
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
				{
					ID: "GROUP-USER-ID",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:  "POLICY-USER-ID",
					Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
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
						{
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
				{
					ID: "GROUP-USER-ID",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:  "POLICY-USER-ID",
					Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
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
						{
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

		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = test.getAttachedPoliciesResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][1] = test.getAttachedPoliciesError

		restrictions, err := testAPI.getRestrictions(test.authUserID, test.action, test.resourceUrn)
		if test.wantError != nil {
			if apiError := err.(*Error); test.wantError.Code != apiError.Code {
				t.Errorf("Test %v failed. Received different error codes (wanted:%v / received:%v)", n,
					test.wantError.Code, apiError.Code)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Error: %v", n, err)
				continue
			}

			if testRepo.ArgsIn[GetUserByExternalIDMethod][0] != test.authUserID {
				t.Errorf("Test %v failed. Received different user identifiers (wanted:%v / received:%v)",
					n, test.authUserID, testRepo.ArgsIn[GetUserByExternalIDMethod][0])
				continue
			}

			if param := testRepo.ArgsIn[GetGroupsByUserIDMethod][0]; test.authUserID != "" && param != test.authUserID {
				t.Errorf("Test %v failed. Received different user identifiers (wanted:%v / received:%v)",
					n, test.authUserID, testRepo.ArgsIn[GetGroupsByUserIDMethod][0])
				continue
			}

			if param := testRepo.ArgsIn[GetAttachedPoliciesMethod][0]; test.getGroupsByUserIDResult != nil &&
				param != test.getGroupsByUserIDResult[0].ID {
				t.Errorf("Test %v failed. Received different user identifiers (wanted:%v / received:%v)",
					n, test.authUserID, testRepo.ArgsIn[GetAttachedPoliciesMethod][0])
				continue
			}

			// Check result
			if diff := pretty.Compare(restrictions, test.expectedRestrictions); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
		}
	}
}

func TestGetGroupsByUser(t *testing.T) {
	testcases := map[string]struct {
		// User ID to retrieve its groups
		userID string
		// Expected Groups
		expectedGroups []Group
		// Error to compare when we expect an error
		wantError *Error
		// GetGroupsByUserID Method Out Arguments
		getGroupsByUserIDResult []Group
		getGroupsByUserIDError  error
	}{
		"OktestCase": {
			userID: "UserID",
			expectedGroups: []Group{
				{
					ID: "GROUP-USER-ID1",
				},
				{
					ID: "GROUP-USER-ID2",
				},
			},
			getGroupsByUserIDResult: []Group{
				{
					ID: "GROUP-USER-ID1",
				},
				{
					ID: "GROUP-USER-ID2",
				},
			},
		},
		"ErrortestCase": {
			userID: "UserID",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getGroupsByUserIDError: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
	}

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for n, test := range testcases {

		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = test.getGroupsByUserIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][1] = test.getGroupsByUserIDError

		groups, err := testAPI.getGroupsByUser(test.userID)
		if test.wantError != nil {
			if apiError := err.(*Error); test.wantError.Code != apiError.Code {
				t.Errorf("Test %v failed. Received different error codes (wanted:%v / received:%v)", n,
					test.wantError.Code, apiError.Code)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Error: %v", n, err)
				continue
			}

			if param := testRepo.ArgsIn[GetGroupsByUserIDMethod][0]; param != test.userID {
				t.Errorf("Test %v failed. Received different user identifiers (wanted:%v / received:%v)",
					n, test.userID, testRepo.ArgsIn[GetGroupsByUserIDMethod][0])
				continue
			}

			// Check result
			if diff := pretty.Compare(groups, test.expectedGroups); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
		}
	}
}

func TestGetPoliciesByGroups(t *testing.T) {
	testcases := map[string]struct {
		groups           []Group
		expectedPolicies []Policy
		// Error to compare when we expect an error
		wantError *Error
		// GetAttachedPolicies Method Out Arguments
		getAttachedPoliciesResult []Policy
		getAttachedPoliciesError  error
	}{
		"OktestCaseEmptyGroups": {
			groups: []Group{},
		},
		"OktestCaseNilGroups": {},
		"OktestCaseNoPoliciesForGroups": {
			groups: []Group{
				{
					ID: "GroupID1",
				},
				{
					ID: "GroupID2",
				},
			},
			expectedPolicies:          []Policy{},
			getAttachedPoliciesResult: nil,
		},
		"OktestCase": {
			groups: []Group{
				{
					ID: "GroupID1",
				},
				{
					ID: "GroupID2",
				},
			},
			expectedPolicies: []Policy{
				{
					ID: "PolicyID",
				},
				{
					ID: "PolicyID",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID: "PolicyID",
				},
			},
		},
		"ErrortestCase": {
			groups: []Group{
				{
					ID: "GroupID1",
				},
				{
					ID: "GroupID2",
				},
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getAttachedPoliciesError: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
	}

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for n, test := range testcases {

		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = test.getAttachedPoliciesResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][1] = test.getAttachedPoliciesError

		policies, err := testAPI.getPoliciesByGroups(test.groups)
		if test.wantError != nil {
			if apiError := err.(*Error); test.wantError.Code != apiError.Code {
				t.Errorf("Test %v failed. Received different error codes (wanted:%v / received:%v)", n,
					test.wantError.Code, apiError.Code)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Error: %v", n, err)
				continue
			}

			if param := testRepo.ArgsIn[GetAttachedPoliciesMethod][0]; len(test.groups) > 0 &&
				param != test.groups[len(test.groups)-1].ID {
				t.Errorf("Test %v failed. Received different group identifiers (wanted:%v / received:%v)",
					n, test.groups[len(test.groups)-1].ID, testRepo.ArgsIn[GetAttachedPoliciesMethod][0])
				continue
			}

			// Check result
			if diff := pretty.Compare(policies, test.expectedPolicies); diff != "" {
				t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
				continue
			}
		}
	}
}

func TestGetStatementsByRequestedAction(t *testing.T) {
	testcases := map[string]struct {
		// Policies to retrieve its statements according to an action
		policies []Policy
		action   string
		// Expected data
		expectedStatements []Statement
	}{
		"OktestCaseEmptyPolicies": {
			policies: []Policy{},
			action:   "action",
		},
		"OktestCaseNilPolicies": {
			action: "action",
		},
		"OktestCaseFilteredStatements": {
			policies: []Policy{
				{
					ID: "PolicyID1Contained",
					Statements: &[]Statement{
						{
							Effect: "allow",
							Action: []string{
								"act*", "noaction",
							},
							Resources: []string{
								CreateUrn("example", RESOURCE_GROUP, "/path1/", "groupAllow"),
								CreateUrn("example", RESOURCE_GROUP, "/path2/", "groupAllow"),
								GetUrnPrefix("example", RESOURCE_GROUP, "/path1/"),
								GetUrnPrefix("example", RESOURCE_GROUP, "/path2/"),
							},
						},
						{
							Effect: "deny",
							Action: []string{
								"noaction", "action",
							},
							Resources: []string{
								CreateUrn("example", RESOURCE_POLICY, "/pathdeny/", "policydeny"),
							},
						},
						{
							Effect: "allow",
							Action: []string{
								"noaction",
							},
							Resources: []string{
								CreateUrn("example", RESOURCE_POLICY, "/nocontained/", "policy"),
							},
						},
					},
				},
				{
					ID: "PolicyID2NoContained",
					Statements: &[]Statement{
						{
							Effect: "allow",
							Action: []string{
								"noact*",
							},
							Resources: []string{
								CreateUrn("example", RESOURCE_GROUP, "/path1/", "groupAllow"),
								CreateUrn("example", RESOURCE_GROUP, "/path2/", "groupAllow"),
								GetUrnPrefix("example", RESOURCE_GROUP, "/path1/"),
								GetUrnPrefix("example", RESOURCE_GROUP, "/path2/"),
							},
						},
					},
				},
			},
			action: "action",
			expectedStatements: []Statement{
				{
					Effect: "allow",
					Action: []string{
						"act*", "noaction",
					},
					Resources: []string{
						CreateUrn("example", RESOURCE_GROUP, "/path1/", "groupAllow"),
						CreateUrn("example", RESOURCE_GROUP, "/path2/", "groupAllow"),
						GetUrnPrefix("example", RESOURCE_GROUP, "/path1/"),
						GetUrnPrefix("example", RESOURCE_GROUP, "/path2/"),
					},
				},
				{
					Effect: "deny",
					Action: []string{
						"noaction", "action",
					},
					Resources: []string{
						CreateUrn("example", RESOURCE_POLICY, "/pathdeny/", "policydeny"),
					},
				},
			},
		},
	}

	for n, test := range testcases {

		statements := getStatementsByRequestedAction(test.policies, test.action)

		// Check result
		if diff := pretty.Compare(statements, test.expectedStatements); diff != "" {
			t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
			continue
		}
	}
}

func TestCleanRepeatedRestrictions(t *testing.T) {

}

func TestIsActionContained(t *testing.T) {
	testcases := map[string]struct {
		actionRequested  string
		statementActions []string
		expectedResponse bool
	}{
		"OktestCaseActionContainedWithPrefix": {
			actionRequested: "action",
			statementActions: []string{
				"*",
				"noAction",
				"noAction2",
			},
			expectedResponse: true,
		},
		"OktestCaseActionContainedWithPrefix2": {
			actionRequested: "action",
			statementActions: []string{
				"noac*",
				"action*",
				"noaction",
			},
			expectedResponse: true,
		},
		"OktestCaseActionContainedWithoutPrefix": {
			actionRequested: "action",
			statementActions: []string{
				"example1",
				"example2",
				"action",
			},
			expectedResponse: true,
		},
		"OktestCaseNoActionContainedWithPrefix": {
			actionRequested: "action",
			statementActions: []string{
				"actn*",
				"actions*",
				"noaction*",
			},
			expectedResponse: false,
		},
		"OktestCaseNoActionContainedWithoutPrefix": {
			actionRequested: "action",
			statementActions: []string{
				"actions",
				"actio",
				"acti",
			},
			expectedResponse: false,
		},
	}

	for n, test := range testcases {

		isContained := isActionContained(test.actionRequested, test.statementActions)

		// Check result
		if test.expectedResponse != isContained {
			t.Errorf("Test %v failed. Received different responses (wanted:%v / received:%v)",
				n, test.expectedResponse, isContained)
			continue
		}
	}
}

func TestIsResourceContained(t *testing.T) {
	testcases := map[string]struct {
		resource         string
		resourcePrefix   string
		expectedResponse bool
	}{
		"OktestCaseContainedWithRoot": {
			resource:         "resource",
			resourcePrefix:   "*",
			expectedResponse: true,
		},
		"OktestCaseContainedWithPrefix": {
			resource:         "resource",
			resourcePrefix:   "res*",
			expectedResponse: true,
		},
		"OktestCaseNoContainedWithPrefix": {
			resource:         "resource",
			resourcePrefix:   "nores*",
			expectedResponse: false,
		},
	}

	for n, test := range testcases {

		isContained := isContainedOrEqual(test.resource, test.resourcePrefix)

		// Check result
		if test.expectedResponse != isContained {
			t.Errorf("Test %v failed. Received different values (wanted:%v / received:%v)",
				n, test.expectedResponse, isContained)
			continue
		}
	}
}

func TestIsFullUrn(t *testing.T) {
	testcases := map[string]struct {
		resource         string
		expectedResponse bool
	}{
		"OktestCaseIsFullUrn": {
			resource:         "resource",
			expectedResponse: true,
		},
		"OktestCaseIsNotFullUrn": {
			resource:         "resource*",
			expectedResponse: false,
		},
	}

	for n, test := range testcases {

		isFullUrn := isFullUrn(test.resource)

		// Check result
		if test.expectedResponse != isFullUrn {
			t.Errorf("Test %v failed. Received different values (wanted:%v / received:%v)",
				n, test.expectedResponse, isFullUrn)
			continue
		}
	}
}

func TestInsertRestriction(t *testing.T) {
	testcases := map[string]struct {
		resource struct {
			isAllow   bool
			isFullUrn bool
			urn       string
		}
		restrictions         *Restrictions
		expectedRestrictions *Restrictions
	}{
		"AllowFullUrn1": {
			resource: struct {
				isAllow   bool
				isFullUrn bool
				urn       string
			}{
				isAllow:   true,
				isFullUrn: true,
				urn:       "asd:/path/asd",
			},
			restrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{"asd:/*"},
				DeniedFullUrns:     []string{},
			},
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{"asd:/*"},
				DeniedFullUrns:     []string{},
			},
		},
		"AllowFullUrn2": {
			resource: struct {
				isAllow   bool
				isFullUrn bool
				urn       string
			}{
				isAllow:   true,
				isFullUrn: true,
				urn:       "asd:/path/asd",
			},
			restrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{"asd:/path/asd"},
			},
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{"asd:/path/asd"},
			},
		},
		"AllowFullUrn3": {
			resource: struct {
				isAllow   bool
				isFullUrn bool
				urn       string
			}{
				isAllow:   true,
				isFullUrn: true,
				urn:       "asd:/path/asd",
			},
			restrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"asd:/*"},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{},
			},
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"asd:/*"},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{},
			},
		},
		"AllowFullUrn4": {
			resource: struct {
				isAllow   bool
				isFullUrn bool
				urn       string
			}{
				isAllow:   true,
				isFullUrn: true,
				urn:       "asd:/path/asd",
			},
			restrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"asd:/path2/*"},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{},
			},
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"asd:/path2/*"},
				AllowedFullUrns:    []string{"asd:/path/asd"},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{},
			},
		},
		"AllowFullUrn5": {
			resource: struct {
				isAllow   bool
				isFullUrn bool
				urn       string
			}{
				isAllow:   true,
				isFullUrn: true,
				urn:       "asd:/path/asd",
			},
			restrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"asd:/path2/*"},
				AllowedFullUrns:    []string{"asd:/path/asd"},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{"asd:/path3/zxc"},
			},
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"asd:/path2/*"},
				AllowedFullUrns:    []string{"asd:/path/asd"},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{"asd:/path3/zxc"},
			},
		},

		"AllowPrefix1": {
			resource: struct {
				isAllow   bool
				isFullUrn bool
				urn       string
			}{
				isAllow:   true,
				isFullUrn: false,
				urn:       "asd:/path/*",
			},
			restrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"asd:/path2/*"},
				AllowedFullUrns:    []string{"asd:/path/asd"},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{"asd:/path3/zxc"},
			},
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"asd:/path2/*", "asd:/path/*"},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{"asd:/path3/zxc"},
			},
		},
		"AllowPrefix2": {
			resource: struct {
				isAllow   bool
				isFullUrn bool
				urn       string
			}{
				isAllow:   true,
				isFullUrn: false,
				urn:       "asd:/path/*",
			},
			restrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{},
			},
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"asd:/path/*"},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{},
			},
		},
		"AllowPrefix3": {
			resource: struct {
				isAllow   bool
				isFullUrn bool
				urn       string
			}{
				isAllow:   true,
				isFullUrn: false,
				urn:       "asd:/path/*",
			},
			restrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"asd:/*"},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{},
			},
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"asd:/*"},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{},
			},
		},
		"AllowPrefix4": {
			resource: struct {
				isAllow   bool
				isFullUrn bool
				urn       string
			}{
				isAllow:   true,
				isFullUrn: false,
				urn:       "asd:/path/*",
			},
			restrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{"asd:/path/asd1", "asd:/path/asd2"},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{},
			},
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"asd:/path/*"},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{},
			},
		},
		"AllowPrefix5": {
			resource: struct {
				isAllow   bool
				isFullUrn bool
				urn       string
			}{
				isAllow:   true,
				isFullUrn: false,
				urn:       "asd:/path/*",
			},
			restrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"zxc:/path2/*"},
				AllowedFullUrns:    []string{"zxc:/path/asd"},
				DeniedUrnPrefixes:  []string{"asd:/*"},
				DeniedFullUrns:     []string{"asd:/path3/asd"},
			},
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"zxc:/path2/*"},
				AllowedFullUrns:    []string{"zxc:/path/asd"},
				DeniedUrnPrefixes:  []string{"asd:/*"},
				DeniedFullUrns:     []string{"asd:/path3/asd"},
			},
		},

		"DenyFullUrn1": {
			resource: struct {
				isAllow   bool
				isFullUrn bool
				urn       string
			}{
				isAllow:   false,
				isFullUrn: true,
				urn:       "asd:/path/asd",
			},
			restrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{"asd:/*"},
				DeniedFullUrns:     []string{},
			},
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{"asd:/*"},
				DeniedFullUrns:     []string{},
			},
		},
		"DenyFullUrn2": {
			resource: struct {
				isAllow   bool
				isFullUrn bool
				urn       string
			}{
				isAllow:   false,
				isFullUrn: true,
				urn:       "asd:/path/asd",
			},
			restrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{"asd:/path/asd"},
			},
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{"asd:/path/asd"},
			},
		},
		"DenyFullUrn3": {
			resource: struct {
				isAllow   bool
				isFullUrn bool
				urn       string
			}{
				isAllow:   false,
				isFullUrn: true,
				urn:       "asd:/path/asd",
			},
			restrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"asd:/*"},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{},
			},
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"asd:/*"},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{"asd:/path/asd"},
			},
		},
		"DenyFullUrn4": {
			resource: struct {
				isAllow   bool
				isFullUrn bool
				urn       string
			}{
				isAllow:   false,
				isFullUrn: true,
				urn:       "asd:/path/asd",
			},
			restrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"asd:/path2/*"},
				AllowedFullUrns:    []string{"asd:/path/asd"},
				DeniedUrnPrefixes:  []string{"asd:/path3/*"},
				DeniedFullUrns:     []string{},
			},
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"asd:/path2/*"},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{"asd:/path3/*"},
				DeniedFullUrns:     []string{"asd:/path/asd"},
			},
		},
		"DenyFullUrn5": {
			resource: struct {
				isAllow   bool
				isFullUrn bool
				urn       string
			}{
				isAllow:   false,
				isFullUrn: true,
				urn:       "asd:/path/asd",
			},
			restrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"asd:/path2/*"},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{"asd:/path/*"},
				DeniedFullUrns:     []string{"asd:/path3/asd"},
			},
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"asd:/path2/*"},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{"asd:/path/*"},
				DeniedFullUrns:     []string{"asd:/path3/asd"},
			},
		},

		"DenyPrefix1": {
			resource: struct {
				isAllow   bool
				isFullUrn bool
				urn       string
			}{
				isAllow:   false,
				isFullUrn: false,
				urn:       "asd:/path/*",
			},
			restrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"asd:/path2/*"},
				AllowedFullUrns:    []string{"asd:/path/asd"},
				DeniedUrnPrefixes:  []string{"asd:/path2/*"},
				DeniedFullUrns:     []string{"asd:/path3/zxc"},
			},
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"asd:/path2/*"},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{"asd:/path2/*", "asd:/path/*"},
				DeniedFullUrns:     []string{"asd:/path3/zxc"},
			},
		},
		"DenyPrefix2": {
			resource: struct {
				isAllow   bool
				isFullUrn bool
				urn       string
			}{
				isAllow:   false,
				isFullUrn: false,
				urn:       "asd:/path/*",
			},
			restrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{},
			},
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{"asd:/path/*"},
				DeniedFullUrns:     []string{},
			},
		},
		"DenyPrefix3": {
			resource: struct {
				isAllow   bool
				isFullUrn bool
				urn       string
			}{
				isAllow:   false,
				isFullUrn: false,
				urn:       "asd:/path/*",
			},
			restrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{"asd:/*"},
				DeniedFullUrns:     []string{},
			},
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{"asd:/*"},
				DeniedFullUrns:     []string{},
			},
		},
		"DenyPrefix4": {
			resource: struct {
				isAllow   bool
				isFullUrn bool
				urn       string
			}{
				isAllow:   false,
				isFullUrn: false,
				urn:       "asd:/path/*",
			},
			restrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{"asd:/path/asd1", "asd:/path/asd2"},
			},
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{"asd:/path/*"},
				DeniedFullUrns:     []string{},
			},
		},
		"DenyPrefix5": {
			resource: struct {
				isAllow   bool
				isFullUrn bool
				urn       string
			}{
				isAllow:   false,
				isFullUrn: false,
				urn:       "asd:/path/*",
			},
			restrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"zxc:/path2/*"},
				AllowedFullUrns:    []string{"zxc:/path/asd"},
				DeniedUrnPrefixes:  []string{"asd:/*"},
				DeniedFullUrns:     []string{"asd:/path3/asd"},
			},
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{"zxc:/path2/*"},
				AllowedFullUrns:    []string{"zxc:/path/asd"},
				DeniedUrnPrefixes:  []string{"asd:/*"},
				DeniedFullUrns:     []string{"asd:/path3/asd"},
			},
		},
	}

	for n, test := range testcases {

		test.restrictions.insertRestriction(test.resource.isAllow, test.resource.isFullUrn, test.resource.urn)

		// Check result
		if diff := pretty.Compare(test.restrictions, test.expectedRestrictions); diff != "" {
			t.Errorf("Test %v failed. Received different values (wanted / received): %v",
				n, diff)
			continue
		}
	}
}

func TestGetRestrictionsWhenResourceRequestedIsPrefix(t *testing.T) {
	testcases := map[string]struct {
		statements           []Statement
		resource             string
		expectedRestrictions *Restrictions
	}{
		"OktestCaseEmptyStatement": {
			statements: []Statement{},
			resource:   "resource*",
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{},
			},
		},
		"OktestCaseIsNotFullUrn": {
			resource: "resource*",
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{},
			},
		},
		"OktestCaseStatementResourcePrefix": {
			statements: []Statement{
				{
					Effect: "allow",
					Action: []string{
						USER_ACTION_GET_USER, USER_ACTION_CREATE_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path1/"),
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
						GetUrnPrefix("", RESOURCE_USER, "/"),
						GetUrnPrefix("", RESOURCE_USER, "/path"),
					},
				},
				{
					Effect: "deny",
					Action: []string{
						USER_ACTION_DELETE_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path1/"),
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
						GetUrnPrefix("", RESOURCE_USER, "/"),
						GetUrnPrefix("", RESOURCE_USER, "/path"),
					},
				},
			},
			resource: GetUrnPrefix("", RESOURCE_USER, "/path"),
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes: []string{
					GetUrnPrefix("", RESOURCE_USER, "/"),
				},
				DeniedFullUrns: []string{},
			},
		},
		"OktestCaseStatementResourceIsFull": {
			statements: []Statement{
				{
					Effect: "allow",
					Action: []string{
						USER_ACTION_GET_USER, USER_ACTION_CREATE_USER,
					},
					Resources: []string{
						CreateUrn("", RESOURCE_USER, "/path/", "userAllowed"),
					},
				},
				{
					Effect: "deny",
					Action: []string{
						USER_ACTION_DELETE_USER,
					},
					Resources: []string{
						CreateUrn("", RESOURCE_USER, "/path/", "userDenied"),
					},
				},
			},
			resource: GetUrnPrefix("", RESOURCE_USER, "/path"),
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns: []string{
					CreateUrn("", RESOURCE_USER, "/path/", "userAllowed"),
				},
				DeniedUrnPrefixes: []string{},
				DeniedFullUrns: []string{
					CreateUrn("", RESOURCE_USER, "/path/", "userDenied"),
				},
			},
		},
	}

	for n, test := range testcases {

		restrictions := getRestrictions(test.statements, test.resource, isFullUrn(test.resource))

		// Check result
		if diff := pretty.Compare(restrictions, test.expectedRestrictions); diff != "" {
			t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
			continue
		}
	}
}

func TestGetRestrictionsWhenResourceRequestedIsFullUrn(t *testing.T) {
	testcases := map[string]struct {
		statements           []Statement
		resource             string
		expectedRestrictions *Restrictions
	}{
		"OktestCaseEmptyStatement": {
			statements: []Statement{},
			resource:   "resource",
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{},
			},
		},
		"OktestCaseIsNotFullUrn": {
			resource: "resource",
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{},
			},
		},
		"OktestCaseStatementResourcePrefix": {
			statements: []Statement{
				{
					Effect: "allow",
					Action: []string{
						USER_ACTION_GET_USER, USER_ACTION_CREATE_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
						GetUrnPrefix("", RESOURCE_USER, "/"),
					},
				},
				{
					Effect: "deny",
					Action: []string{
						USER_ACTION_DELETE_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
						GetUrnPrefix("", RESOURCE_USER, "/"),
					},
				},
			},
			resource: CreateUrn("", RESOURCE_USER, "/path/", "user"),
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes: []string{
					GetUrnPrefix("", RESOURCE_USER, "/"),
				},
				DeniedFullUrns: []string{},
			},
		},
		"OktestCaseStatementResourceIsFull": {
			statements: []Statement{
				{
					Effect: "allow",
					Action: []string{
						USER_ACTION_GET_USER, USER_ACTION_CREATE_USER,
					},
					Resources: []string{
						CreateUrn("", RESOURCE_USER, "/path/", "user"),
						CreateUrn("", RESOURCE_USER, "/path/", "user2"),
					},
				},
				{
					Effect: "deny",
					Action: []string{
						USER_ACTION_DELETE_USER,
					},
					Resources: []string{
						CreateUrn("", RESOURCE_USER, "/path/", "user"),
						CreateUrn("", RESOURCE_USER, "/path/", "user2"),
					},
				},
			},
			resource: CreateUrn("", RESOURCE_USER, "/path/", "user"),
			expectedRestrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns: []string{
					CreateUrn("", RESOURCE_USER, "/path/", "user"),
				},
			},
		},
	}

	for n, test := range testcases {

		restrictions := getRestrictions(test.statements, test.resource, isFullUrn(test.resource))

		// Check result
		if diff := pretty.Compare(restrictions, test.expectedRestrictions); diff != "" {
			t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
			continue
		}
	}
}

func TestFilterResources(t *testing.T) {
	testcases := map[string]struct {
		resources         []Resource
		restrictions      *Restrictions
		expectedResources []Resource
	}{
		"OktestCaseEmptyResources": {
			resources: []Resource{},
			restrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{},
			},
			expectedResources: []Resource{},
		},
		"OktestCaseAllow": {
			resources: []Resource{
				User{
					Urn: CreateUrn("", RESOURCE_USER, "/path/", "user"),
				},
			},
			restrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns: []string{
					CreateUrn("", RESOURCE_USER, "/path/", "user"),
				},
				DeniedUrnPrefixes: []string{},
				DeniedFullUrns:    []string{},
			},
			expectedResources: []Resource{
				User{
					Urn: CreateUrn("", RESOURCE_USER, "/path/", "user"),
				},
			},
		},
		"OktestCaseDeny": {
			resources: []Resource{
				User{
					Urn: CreateUrn("", RESOURCE_USER, "/path/", "user"),
				},
			},
			restrictions: &Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns: []string{
					CreateUrn("", RESOURCE_USER, "/path/", "user"),
				},
				DeniedUrnPrefixes: []string{
					GetUrnPrefix("", RESOURCE_USER, "/path"),
				},
				DeniedFullUrns: []string{},
			},
			expectedResources: []Resource{},
		},
	}

	for n, test := range testcases {

		filteredResources := filterResources(test.resources, test.restrictions)

		// Check result
		if diff := pretty.Compare(filteredResources, test.expectedResources); diff != "" {
			t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
			continue
		}
	}
}

func TestIsAllowedResource(t *testing.T) {
	testcases := map[string]struct {
		resource     Resource
		restrictions Restrictions
		expectedData bool
	}{
		"OktestCaseNoRestrictions": {
			resource: User{
				Urn: CreateUrn("", RESOURCE_USER, "/path/", "user"),
			},
			restrictions: Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns:     []string{},
			},
			expectedData: false,
		},
		"OktestCaseDeniedByUrnPrefix": {
			resource: User{
				Urn: CreateUrn("", RESOURCE_USER, "/path/", "user"),
			},
			restrictions: Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes: []string{
					GetUrnPrefix("", RESOURCE_USER, "/path"),
				},
				DeniedFullUrns: []string{},
			},
			expectedData: false,
		},
		"OktestCaseDeniedByFullUrn": {
			resource: User{
				Urn: CreateUrn("", RESOURCE_USER, "/path/", "user"),
			},
			restrictions: Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns:    []string{},
				DeniedUrnPrefixes:  []string{},
				DeniedFullUrns: []string{
					CreateUrn("", RESOURCE_USER, "/path/", "user"),
				},
			},
			expectedData: false,
		},
		"OktestCaseAllowedByUrnPrefix": {
			resource: User{
				Urn: CreateUrn("", RESOURCE_USER, "/path/", "user"),
			},
			restrictions: Restrictions{
				AllowedUrnPrefixes: []string{
					GetUrnPrefix("", RESOURCE_USER, "/path"),
				},
				AllowedFullUrns:   []string{},
				DeniedUrnPrefixes: []string{},
				DeniedFullUrns:    []string{},
			},
			expectedData: true,
		},
		"OktestCaseAllowedByFullUrn": {
			resource: User{
				Urn: CreateUrn("", RESOURCE_USER, "/path/", "user"),
			},
			restrictions: Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns: []string{
					CreateUrn("", RESOURCE_USER, "/path/", "user"),
				},
				DeniedUrnPrefixes: []string{},
				DeniedFullUrns:    []string{},
			},
			expectedData: true,
		},
		"OktestCaseConflictDenyAndAllow": {
			resource: User{
				Urn: CreateUrn("", RESOURCE_USER, "/path/", "user"),
			},
			restrictions: Restrictions{
				AllowedUrnPrefixes: []string{},
				AllowedFullUrns: []string{
					CreateUrn("", RESOURCE_USER, "/path/", "user"),
				},
				DeniedUrnPrefixes: []string{
					GetUrnPrefix("", RESOURCE_USER, "/path"),
				},
				DeniedFullUrns: []string{},
			},
			expectedData: false,
		},
		"OktestCaseAllowedWithPrefixAndFull": {
			resource: User{
				Urn: CreateUrn("", RESOURCE_USER, "/path/", "user"),
			},
			restrictions: Restrictions{
				AllowedUrnPrefixes: []string{
					GetUrnPrefix("", RESOURCE_USER, "/path"),
				},
				AllowedFullUrns: []string{
					CreateUrn("", RESOURCE_USER, "/path/", "user"),
				},
				DeniedUrnPrefixes: []string{},
				DeniedFullUrns:    []string{},
			},
			expectedData: true,
		},
	}

	for n, test := range testcases {

		response := isAllowedResource(test.resource, test.restrictions)

		// Check result
		if diff := pretty.Compare(response, test.expectedData); diff != "" {
			t.Errorf("Test %v failed. Received different responses (received/wanted) %v", n, diff)
			continue
		}
	}
}
