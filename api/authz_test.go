package api

import (
	"reflect"
	"testing"
	"time"
)

func TestGetUsersAuthorized(t *testing.T) {
	testcases := map[string]struct {
		// User authenticated
		user AuthenticatedUser
		// Resource urn that user wants to access
		resourceUrn string
		// Action to do
		action string
		// Resources received from db that system has to authorize
		users []User
	}{
		"OKtestCaseAdmin": {
			user: AuthenticatedUser{
				Identifier: "123456",
				Admin:      true,
			},
			resourceUrn: CreateUrn("", RESOURCE_USER, "/path/", "user1"),
			action:      USER_ACTION_GET_USER,
			users: []User{
				User{
					ID:         "654321",
					ExternalID: "",
					Path:       "/path/",
					CreateAt:   time.Now().UTC(),
					Urn:        CreateUrn("", RESOURCE_USER, "/path/", "user1"),
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
		t.Logf("Running test case %v", n)
		testRepo.ArgsOut[GetUserByExternalIDMethod] = make([]interface{}, 2)
		testRepo.ArgsIn[GetUserByExternalIDMethod] = make([]interface{}, 2)
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = &User{}
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = nil

		// Run test
		authorizedUsers, err := api.GetUsersAuthorized(test.user, test.resourceUrn, test.action, test.users)
		if err != nil {
			t.Errorf("Userid: %v", testRepo.ArgsIn[GetUserByExternalIDMethod][0].(string))
			t.Errorf("Unexpected error in test case %v, error: %v", n, err)
		}

		// Check result
		if !reflect.DeepEqual(authorizedUsers, test.users) {
			t.Errorf("Struct are not equal received [%v], wanted [%v]", authorizedUsers, test.users)
		}
	}
}

func TestGroupsAuthorized(t *testing.T) {

}

func TestGetPoliciesAuthorized(t *testing.T) {

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
