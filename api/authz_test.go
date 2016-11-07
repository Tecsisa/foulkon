package api

import (
	"testing"

	"fmt"

	"github.com/Tecsisa/foulkon/database"
	"github.com/stretchr/testify/assert"
)

func TestGetAuthorizedUsers(t *testing.T) {
	testcases := map[string]struct {
		// Authenticated user
		requestInfo RequestInfo
		// Resource urn that user wants to access
		resourceUrn string
		// Action to do
		action string
		// Resources received from db that system has to authorize
		usersToAuthorize []User
		// Resources authorized by method
		usersAuthorized []User
		// Error to compare when we expect an error
		wantError error
		// GetUserByExternalID Method Out Arguments
		getUserByExternalIDResult *User
		getUserByExternalIDError  error
	}{
		"OKtestCaseAdmin": {
			requestInfo: RequestInfo{
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
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			resourceUrn:      CreateUrn("", RESOURCE_USER, "/path/", "user1"),
			action:           USER_ACTION_GET_USER,
			usersToAuthorize: []User{},
			usersAuthorized:  []User{},
		},
		"ErrortestCaseAuthenticatedUserNotExist": {
			requestInfo: RequestInfo{
				Identifier: "notAdminUser",
				Admin:      false,
			},
			resourceUrn: CreateUrn("", RESOURCE_USER, "/path/", "user1"),
			action:      USER_ACTION_GET_USER,
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId notAdminUser not found. Unable to retrieve permissions.",
			},
			getUserByExternalIDResult: nil,
			getUserByExternalIDError: &database.Error{
				Code:    database.USER_NOT_FOUND,
				Message: "Error",
			},
		},
	}

	for n, test := range testcases {

		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = test.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = test.getUserByExternalIDError

		authorizedUsers, err := testAPI.GetAuthorizedUsers(test.requestInfo, test.resourceUrn, test.action, test.usersToAuthorize)
		checkMethodResponse(t, n, test.wantError, err, test.usersAuthorized, authorizedUsers)
		if !test.requestInfo.Admin {
			// Check received authenticated user in method GetUserByExternalID
			assert.Equal(t, test.requestInfo.Identifier, testRepo.ArgsIn[GetUserByExternalIDMethod][0], "Error in test case %v", n)
		}
	}
}

func TestGetAuthorizedGroups(t *testing.T) {
	testcases := map[string]struct {
		// Authenticated user
		requestInfo RequestInfo
		// Resource urn that user wants to access
		resourceUrn string
		// Action to do
		action string
		// Resources received from db that system has to authorize
		groupsToAuthorize []Group
		// Resources authorized by method
		groupsAuthorized []Group
		// Error to compare when we expect an error
		wantError error
		// GetUserByExternalID Method Out Arguments
		getUserByExternalIDResult *User
		getUserByExternalIDError  error
	}{
		"OKtestCaseAdmin": {
			requestInfo: RequestInfo{
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
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			resourceUrn:       CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
			action:            GROUP_ACTION_GET_GROUP,
			groupsToAuthorize: []Group{},
			groupsAuthorized:  []Group{},
		},
		"ErrortestCaseAuthenticatedUserNotExist": {
			requestInfo: RequestInfo{
				Identifier: "notAdminUser",
				Admin:      false,
			},
			resourceUrn: CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
			action:      GROUP_ACTION_GET_GROUP,
			wantError: &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: "Error",
			},
			getUserByExternalIDResult: nil,
			getUserByExternalIDError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "Error",
			},
		},
	}

	for n, test := range testcases {

		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = test.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = test.getUserByExternalIDError

		authorizedGroups, err := testAPI.GetAuthorizedGroups(test.requestInfo, test.resourceUrn, test.action, test.groupsToAuthorize)
		checkMethodResponse(t, n, test.wantError, err, test.groupsAuthorized, authorizedGroups)
		if !test.requestInfo.Admin {
			// Check received authenticated user in method GetUserByExternalID
			assert.Equal(t, test.requestInfo.Identifier, testRepo.ArgsIn[GetUserByExternalIDMethod][0], "Error in test case %v", n)
		}
	}
}

func TestGetAuthorizedPolicies(t *testing.T) {
	testcases := map[string]struct {
		// Authenticated user
		requestInfo RequestInfo
		// Resource urn that user wants to access
		resourceUrn string
		// Action to do
		action string
		// Resources received from db that system has to authorize
		policiesToAuthorize []Policy
		// Resources authorized by method
		policiesAuthorized []Policy
		// Error to compare when we expect an error
		wantError error
		// GetUserByExternalID Method Out Arguments
		getUserByExternalIDResult *User
		getUserByExternalIDError  error
	}{
		"OKtestCaseAdmin": {
			requestInfo: RequestInfo{
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
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			resourceUrn:         CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
			action:              POLICY_ACTION_GET_POLICY,
			policiesToAuthorize: []Policy{},
			policiesAuthorized:  []Policy{},
		},
		"ErrortestCaseDatabaseError": {
			requestInfo: RequestInfo{
				Identifier: "USER-AUTHENTICATED",
				Admin:      false,
			},
			resourceUrn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
			action:      POLICY_ACTION_GET_POLICY,
			wantError: &Error{
				Code:    UNKNOWN_API_ERROR,
				Message: "Error",
			},
			getUserByExternalIDResult: nil,
			getUserByExternalIDError: &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "Error",
			},
		},
	}

	for n, test := range testcases {

		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = test.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = test.getUserByExternalIDError

		authorizedPolicies, err := testAPI.GetAuthorizedPolicies(test.requestInfo, test.resourceUrn, test.action, test.policiesToAuthorize)
		checkMethodResponse(t, n, test.wantError, err, test.policiesAuthorized, authorizedPolicies)
		if !test.requestInfo.Admin {
			// Check received authenticated user in method GetUserByExternalID
			assert.Equal(t, test.requestInfo.Identifier, testRepo.ArgsIn[GetUserByExternalIDMethod][0], "Error in test case %v", n)
		}
	}
}

func TestGetAuthorizedExternalResources(t *testing.T) {
	testcases := map[string]struct {
		// Authenticated user
		requestInfo RequestInfo
		// Resource urns that user wants to access
		resourceUrns []string
		// Action to do
		action string
		// Expected allowed resources
		expectedResources []string
		// Error to compare when we expect an error
		wantError error
		// GetUserByExternalID Method Out Arguments
		getUserByExternalIDResult *User
		getUserByExternalIDError  error
		// GetGroupsByUserID Method Out Arguments
		getGroupsByUserIDResult []TestUserGroupRelation
		getGroupsByUserIDError  error
		// GetAttachedPolicies Method Out Arguments
		getAttachedPoliciesResult []TestPolicyGroupRelation
		getAttachedPoliciesError  error
	}{
		"ErrortestCaseInvalidAction": {
			requestInfo: RequestInfo{
				Admin: true,
			},
			action: "valid::Action",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "No regex match in action: valid::Action",
			},
		},
		"ErrortestCaseInvalidResource": {
			requestInfo: RequestInfo{
				Admin: true,
			},
			action: "product:DoSomething",
			resourceUrns: []string{
				"urn:invalid/resource:resource",
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "No regex match in resource: urn:invalid/resource:resource",
			},
		},
		"ErrortestCaseInvalidResourceWithPrefix": {
			requestInfo: RequestInfo{
				Admin: true,
			},
			action: "product:DoSomething",
			resourceUrns: []string{
				"urn:*",
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter resource urn:*. Urn prefixes are not allowed here",
			},
		},
		"ErrortestCaseEmptyResources": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			action: "product:DoSomething",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: fmt.Sprintf("Invalid parameter Resources. Resources can't be empty or bigger than %v elements", MAX_RESOURCE_NUMBER),
			},
		},
		"ErrortestCaseMaxResourcesExceed": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			action:       "product:DoSomething",
			resourceUrns: getResources(MAX_RESOURCE_NUMBER+1, "urn:iws:iam:org:genericresource/pathname"),
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: fmt.Sprintf("Invalid parameter Resources. Resources can't be empty or bigger than %v elements", MAX_RESOURCE_NUMBER),
			},
		},
		"ErrortestCaseGetRestrictions": {
			action: "product:DoSomething",
			resourceUrns: []string{
				CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId  not found. Unable to retrieve permissions.",
			},
			getUserByExternalIDError: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrortestCaseActionPrefix": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			action: "product:DoPrefix*",
			resourceUrns: []string{
				CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter action product:DoPrefix*. Action parameter can't be a prefix",
			},
		},
		"ErrortestCaseNoAllowedUrns": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			action: POLICY_ACTION_GET_POLICY,
			resourceUrns: []string{
				CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
				CreateUrn("example", RESOURCE_POLICY, "/path/", "policy2"),
				CreateUrn("example1", RESOURCE_POLICY, "/path/", "policy3"),
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to any resource",
			},
			getUserByExternalIDResult: &User{
				ID:  "123456",
				Urn: CreateUrn("", RESOURCE_USER, "/path/", "user1"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:  "GROUP-USER-ID",
						Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:  "POLICY-USER-ID",
						Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "deny",
								Actions: []string{
									POLICY_ACTION_GET_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
								},
							},
							{
								Effect: "allow",
								Actions: []string{
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
		},
		"OktestCaseFullUrnAllow": {
			requestInfo: RequestInfo{
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
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:  "GROUP-USER-ID",
						Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:  "POLICY-USER-ID",
						Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
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
		},
		"OktestCaseFullUrnDeny": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			action: POLICY_ACTION_GET_POLICY,
			resourceUrns: []string{
				CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
				CreateUrn("example", RESOURCE_POLICY, "/path/", "policy2"),
				CreateUrn("example1", RESOURCE_POLICY, "/path/", "policy3"),
				CreateUrn("example2", RESOURCE_POLICY, "/path/path2/", "policy3"),
			},
			expectedResources: []string{CreateUrn("example2", RESOURCE_POLICY, "/path/path2/", "policy3")},
			getUserByExternalIDResult: &User{
				ID:  "123456",
				Urn: CreateUrn("", RESOURCE_USER, "/path/", "user1"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:  "GROUP-USER-ID",
						Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:  "POLICY-USER-ID",
						Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "deny",
								Actions: []string{
									POLICY_ACTION_GET_POLICY,
								},
								Resources: []string{
									GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
								},
							},
							{
								Effect: "allow",
								Actions: []string{
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
		},
		"OktestCaseWithRestrictions": {
			requestInfo: RequestInfo{
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
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:  "GROUP-USER-ID",
						Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:  "POLICY-USER-ID",
						Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
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
								Actions: []string{
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
		},
	}

	for n, test := range testcases {

		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = test.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = test.getUserByExternalIDError

		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = test.getGroupsByUserIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][1] = test.getGroupsByUserIDError

		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = test.getAttachedPoliciesResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][1] = test.getAttachedPoliciesError

		resources, err := testAPI.GetAuthorizedExternalResources(test.requestInfo, test.action, test.resourceUrns)
		checkMethodResponse(t, n, test.wantError, err, test.expectedResources, resources)
		if !test.requestInfo.Admin {
			// Check received authenticated user in method GetUserByExternalID
			assert.Equal(t, test.requestInfo.Identifier, testRepo.ArgsIn[GetUserByExternalIDMethod][0], "Error in test case %v", n)
		}
	}
}

// Test for aux methods of Foulkon

func TestGetAuthorizedResources(t *testing.T) {
	testcases := map[string]struct {
		// Authenticated user
		requestInfo RequestInfo
		// Resource urn that user wants to access
		resourceUrn string
		// Action to do
		action string
		// Resources received from db that system has to authorize
		resourcesToAuthorize []Resource
		// Resources authorized by method
		resourcesAuthorized []Resource
		// Error to compare when we expect an error
		wantError error
		// GetUserByExternalID Method Out Arguments
		getUserByExternalIDResult *User
		getUserByExternalIDError  error
		// GetGroupsByUserID Method Out Arguments
		getGroupsByUserIDResult []TestUserGroupRelation
		getGroupsByUserIDError  error
		// GetAttachedPolicies Method Out Arguments
		getAttachedPoliciesResult []TestPolicyGroupRelation
		getAttachedPoliciesError  error
	}{
		"OKtestCaseAdmin": {
			requestInfo: RequestInfo{
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
			requestInfo: RequestInfo{
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
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
		},
		"ErrortestCaseNotAllowedResources": {
			requestInfo: RequestInfo{
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
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:example:group/path/group1",
			},
		},
		"OKtestCaseResourcesFiltered": {
			requestInfo: RequestInfo{
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
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:  "GROUP-USER-ID",
						Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:  "POLICY-USER-ID",
						Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
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
		},
		"OKtestCaseResourcesFilteredReturnEmpty": {
			// This test case checks if user has access to groups in /path2/ prefix, but there are groups
			// only in /path/, so we expect a empty slice of groups authorized
			requestInfo: RequestInfo{
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
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:  "GROUP-USER-ID",
						Urn: CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:  "POLICY-USER-ID",
						Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
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
		},
	}

	for n, test := range testcases {

		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = test.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = test.getUserByExternalIDError

		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = test.getGroupsByUserIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][1] = test.getGroupsByUserIDError

		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = test.getAttachedPoliciesResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][1] = test.getAttachedPoliciesError

		authorizedResources, err := testAPI.getAuthorizedResources(test.requestInfo, test.resourceUrn, test.action, test.resourcesToAuthorize)
		checkMethodResponse(t, n, test.wantError, err, test.resourcesAuthorized, authorizedResources)
		if !test.requestInfo.Admin {
			// Check received authenticated user in method GetUserByExternalID
			assert.Equal(t, test.requestInfo.Identifier, testRepo.ArgsIn[GetUserByExternalIDMethod][0], "Error in test case %v", n)
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
		wantError error
		// GetUserByExternalID Method Out Arguments
		getUserByExternalIDResult *User
		getUserByExternalIDError  error
		// GetGroupsByUserID Method Out Arguments
		getGroupsByUserIDResult []TestUserGroupRelation
		getGroupsByUserIDError  error
		// GetAttachedPolicies Method Out Arguments
		getAttachedPoliciesResult []TestPolicyGroupRelation
		getAttachedPoliciesError  error
	}{
		"ErrortestCaseGetUserAuthenticatedNotFound": {
			authUserID:  "NotFound",
			resourceUrn: "urn:resource",
			action:      USER_ACTION_GET_USER,
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId NotFound not found. Unable to retrieve permissions.",
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
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID: "GroupID",
					},
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
				AllowedFullUrns:   []string{},
				DeniedFullUrns:    []string{},
				DeniedUrnPrefixes: []string{},
			},
			getUserByExternalIDResult: &User{
				ID: "AuthUserID",
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID: "GROUP-USER-ID",
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:  "POLICY-USER-ID",
						Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
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
								Actions: []string{
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
				AllowedFullUrns: []string{},
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
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID: "GROUP-USER-ID",
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:  "POLICY-USER-ID",
						Urn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
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
								Actions: []string{
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
			}},
	}

	for n, test := range testcases {

		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = test.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = test.getUserByExternalIDError

		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = test.getGroupsByUserIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][2] = test.getGroupsByUserIDError

		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = test.getAttachedPoliciesResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][2] = test.getAttachedPoliciesError

		restrictions, err := testAPI.getRestrictions(test.authUserID, test.action, test.resourceUrn)
		checkMethodResponse(t, n, test.wantError, err, test.expectedRestrictions, restrictions)
		if test.wantError == nil {
			assert.Equal(t, test.authUserID, testRepo.ArgsIn[GetUserByExternalIDMethod][0], "Error in test case %v", n)
			assert.Equal(t, test.authUserID, testRepo.ArgsIn[GetGroupsByUserIDMethod][0], "Error in test case %v", n)
			if test.getGroupsByUserIDResult != nil {
				assert.Equal(t, test.getGroupsByUserIDResult[0].Group.ID, testRepo.ArgsIn[GetAttachedPoliciesMethod][0], "Error in test case %v", n)
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
		wantError error
		// GetGroupsByUserID Method Out Arguments
		getGroupsByUserIDResult []TestUserGroupRelation
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
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID: "GROUP-USER-ID1",
					},
				},
				{
					Group: &Group{
						ID: "GROUP-USER-ID2",
					},
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

	for n, test := range testcases {

		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = test.getGroupsByUserIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][2] = test.getGroupsByUserIDError

		groups, err := testAPI.getGroupsByUser(test.userID)
		checkMethodResponse(t, n, test.wantError, err, test.expectedGroups, groups)
		assert.Equal(t, test.userID, testRepo.ArgsIn[GetGroupsByUserIDMethod][0], "Error in test case %v", n)
	}
}

func TestGetPoliciesByGroups(t *testing.T) {
	testcases := map[string]struct {
		groups           []Group
		expectedPolicies []Policy
		// Error to compare when we expect an error
		wantError error
		// GetAttachedPolicies Method Out Arguments
		getAttachedPoliciesResult []TestPolicyGroupRelation
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
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{

						ID: "PolicyID",
					},
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

	for n, test := range testcases {

		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = test.getAttachedPoliciesResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][2] = test.getAttachedPoliciesError

		policies, err := testAPI.getPoliciesByGroups(test.groups)
		checkMethodResponse(t, n, test.wantError, err, test.expectedPolicies, policies)
		if test.wantError == nil && len(test.groups) > 0 {
			assert.Equal(t, testRepo.ArgsIn[GetAttachedPoliciesMethod][0], test.groups[len(test.groups)-1].ID, "Error in test case %v", n)
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
							Actions: []string{
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
							Actions: []string{
								"noaction", "action",
							},
							Resources: []string{
								CreateUrn("example", RESOURCE_POLICY, "/pathdeny/", "policydeny"),
							},
						},
						{
							Effect: "allow",
							Actions: []string{
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
							Actions: []string{
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
					Actions: []string{
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
					Actions: []string{
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
		checkMethodResponse(t, n, nil, nil, test.expectedStatements, statements)
	}
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
		checkMethodResponse(t, n, nil, nil, test.expectedResponse, isContained)
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
		checkMethodResponse(t, n, nil, nil, test.expectedResponse, isContained)
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
		checkMethodResponse(t, n, nil, nil, test.expectedResponse, isFullUrn)
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
		checkMethodResponse(t, n, nil, nil, test.expectedRestrictions, test.restrictions)
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
					Actions: []string{
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
					Actions: []string{
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
					Actions: []string{
						USER_ACTION_GET_USER, USER_ACTION_CREATE_USER,
					},
					Resources: []string{
						CreateUrn("", RESOURCE_USER, "/path/", "userAllowed"),
					},
				},
				{
					Effect: "deny",
					Actions: []string{
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
		checkMethodResponse(t, n, nil, nil, test.expectedRestrictions, restrictions)
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
					Actions: []string{
						USER_ACTION_GET_USER, USER_ACTION_CREATE_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
						GetUrnPrefix("", RESOURCE_USER, "/"),
					},
				},
				{
					Effect: "deny",
					Actions: []string{
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
					Actions: []string{
						USER_ACTION_GET_USER, USER_ACTION_CREATE_USER,
					},
					Resources: []string{
						CreateUrn("", RESOURCE_USER, "/path/", "user"),
						CreateUrn("", RESOURCE_USER, "/path/", "user2"),
					},
				},
				{
					Effect: "deny",
					Actions: []string{
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
		checkMethodResponse(t, n, nil, nil, test.expectedRestrictions, restrictions)
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
		checkMethodResponse(t, n, nil, nil, test.expectedResources, filteredResources)
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
		checkMethodResponse(t, n, nil, nil, test.expectedData, response)
	}
}
