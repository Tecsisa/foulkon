package api

import (
	"reflect"
	"testing"
	"time"

	"github.com/tecsisa/authorizr/database"
)

func TestGetUsersAuthorized(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// User authenticated
		user AuthenticatedUser
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
		result *User
		err    error
	}{
		"OKtestCaseAdmin": {
			user: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			resourceUrn: CreateUrn("", RESOURCE_USER, "/path/", "user1"),
			action:      USER_ACTION_GET_USER,
			usersToAuthorize: []User{
				User{
					ID:         "654321",
					ExternalID: "user1",
					Path:       "/path/",
					CreateAt:   now,
					Urn:        CreateUrn("", RESOURCE_USER, "/path/", "user1"),
				},
			},
			usersAuthorized: []User{
				User{
					ID:         "654321",
					ExternalID: "user1",
					Path:       "/path/",
					CreateAt:   now,
					Urn:        CreateUrn("", RESOURCE_USER, "/path/", "user1"),
				},
			},
		},
		"OKtestCaseAdminWithEmptyResources": {
			user: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			resourceUrn:      CreateUrn("", RESOURCE_USER, "/path/", "user1"),
			action:           USER_ACTION_GET_USER,
			usersToAuthorize: []User{},
			usersAuthorized:  []User{},
		},
		"ErrortestCaseAuthenticatedUserNotExist": {
			user: AuthenticatedUser{
				Identifier: "notAdminUser",
				Admin:      false,
			},
			resourceUrn: CreateUrn("", RESOURCE_USER, "/path/", "user1"),
			action:      USER_ACTION_GET_USER,
			wantError: &Error{
				Code: UNAUTHORIZED_RESOURCES_ERROR,
			},
			result: nil,
			err: &database.Error{
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
		// Init out resources for method GetUserByExternalID
		testRepo.ArgsIn[GetUserByExternalIDMethod] = make([]interface{}, 1)

		testRepo.ArgsOut[GetUserByExternalIDMethod] = make([]interface{}, 2)
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = test.result
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = test.err

		authorizedUsers, err := api.GetUsersAuthorized(test.user, test.resourceUrn, test.action, test.usersToAuthorize)
		if test.wantError != nil {
			apiError := err.(*Error)
			if test.wantError.Code != apiError.Code {
				t.Fatalf("Unexpected error code in test case %v, want: %v, received: %v", n, test.wantError.Code, apiError.Code)
			}
		} else {
			if err != nil {
				t.Fatalf("Unexpected error in test case %v, error: %v", n, err)
			}

			// Check result
			if !reflect.DeepEqual(authorizedUsers, test.usersAuthorized) {
				t.Fatalf("Struct are not equal received [%v], wanted [%v]", authorizedUsers, test.usersAuthorized)
			}
		}
	}
}

func TestGroupsAuthorized(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// User authenticated
		user AuthenticatedUser
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
		result *User
		err    error
	}{
		"OKtestCaseAdmin": {
			user: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			resourceUrn: CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
			action:      GROUP_ACTION_GET_GROUP,
			groupsToAuthorize: []Group{
				Group{
					ID:       "654321",
					Name:     "group1",
					Path:     "/path/",
					CreateAt: now,
					Urn:      CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
				},
			},
			groupsAuthorized: []Group{
				Group{
					ID:       "654321",
					Name:     "group1",
					Path:     "/path/",
					CreateAt: now,
					Urn:      CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
				},
			},
		},
		"OKtestCaseAdminWithEmptyResources": {
			user: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			resourceUrn:       CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
			action:            GROUP_ACTION_GET_GROUP,
			groupsToAuthorize: []Group{},
			groupsAuthorized:  []Group{},
		},
		"ErrortestCaseAuthenticatedUserNotExist": {
			user: AuthenticatedUser{
				Identifier: "notAdminUser",
				Admin:      false,
			},
			resourceUrn: CreateUrn("example", RESOURCE_GROUP, "/path/", "group1"),
			action:      GROUP_ACTION_GET_GROUP,
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			result: nil,
			err: &database.Error{
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
		// Init out resources for method GetUserByExternalID
		testRepo.ArgsIn[GetUserByExternalIDMethod] = make([]interface{}, 1)

		testRepo.ArgsOut[GetUserByExternalIDMethod] = make([]interface{}, 2)
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = test.result
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = test.err

		authorizedUsers, err := api.GetGroupsAuthorized(test.user, test.resourceUrn, test.action, test.groupsToAuthorize)
		if test.wantError != nil {
			apiError := err.(*Error)
			if test.wantError.Code != apiError.Code {
				t.Fatalf("Unexpected error code in test case %v, want: %v, received: %v", n, test.wantError.Code, apiError.Code)
			}
		} else {
			if err != nil {
				t.Fatalf("Unexpected error in test case %v, error: %v", n, err)
			}

			// Check result
			if !reflect.DeepEqual(authorizedUsers, test.groupsAuthorized) {
				t.Fatalf("Struct are not equal received [%v], wanted [%v]", authorizedUsers, test.groupsAuthorized)
			}
		}
	}
}

func TestGetPoliciesAuthorized(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		// User authenticated
		user AuthenticatedUser
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
		result *User
		err    error
	}{
		"OKtestCaseAdmin": {
			user: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			resourceUrn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
			action:      POLICY_ACTION_GET_POLICY,
			policiesToAuthorize: []Policy{
				Policy{
					ID:       "654321",
					Name:     "group1",
					Path:     "/path/",
					CreateAt: now,
					Urn:      CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
				},
			},
			policiesAuthorized: []Policy{
				Policy{
					ID:       "654321",
					Name:     "group1",
					Path:     "/path/",
					CreateAt: now,
					Urn:      CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
				},
			},
		},
		"OKtestCaseAdminWithEmptyResources": {
			user: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			resourceUrn:         CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
			action:              POLICY_ACTION_GET_POLICY,
			policiesToAuthorize: []Policy{},
			policiesAuthorized:  []Policy{},
		},
		"ErrortestCaseAuthenticatedUserNotExist": {
			user: AuthenticatedUser{
				Identifier: "notAdminUser",
				Admin:      false,
			},
			resourceUrn: CreateUrn("example", RESOURCE_POLICY, "/path/", "policy1"),
			action:      POLICY_ACTION_GET_POLICY,
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			result: nil,
			err: &database.Error{
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
		// Init out resources for method GetUserByExternalID
		testRepo.ArgsIn[GetUserByExternalIDMethod] = make([]interface{}, 1)

		testRepo.ArgsOut[GetUserByExternalIDMethod] = make([]interface{}, 2)
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = test.result
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = test.err

		authorizedUsers, err := api.GetPoliciesAuthorized(test.user, test.resourceUrn, test.action, test.policiesToAuthorize)
		if test.wantError != nil {
			apiError := err.(*Error)
			if test.wantError.Code != apiError.Code {
				t.Fatalf("Unexpected error code in test case %v, want: %v, received: %v", n, test.wantError.Code, apiError.Code)
			}
		} else {
			if err != nil {
				t.Fatalf("Unexpected error in test case %v, error: %v", n, err)
			}

			// Check result
			if !reflect.DeepEqual(authorizedUsers, test.policiesAuthorized) {
				t.Fatalf("Struct are not equal received [%v], wanted [%v]", authorizedUsers, test.policiesAuthorized)
			}
		}
	}
}

func TestGetEffectByUserActionResource(t *testing.T) {

}

// Test for aux methods of Authorizr

func TestGetAuthorizedResources(t *testing.T) {

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
