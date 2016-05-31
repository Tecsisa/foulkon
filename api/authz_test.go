package api

import (
	"github.com/tecsisa/authorizr/database"
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
				Admin:      false,
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
		funcMap: make(map[string]interface{}),
	}
	api := AuthAPI{
		UserRepo:   testRepo,
		GroupRepo:  testRepo,
		PolicyRepo: testRepo,
	}
	for n, test := range testcases {
		t.Logf("Running test case %v", n)
		// Set test value
		testRepo.funcMap[GetUserByExternalIDMethod] = func(id string) (*User, error) {
			return nil, &database.Error{
				Code:    database.INTERNAL_ERROR,
				Message: "Error",
			}
		}
		// Run test
		authorizedUsers, err := api.GetUsersAuthorized(test.user, test.resourceUrn, test.action, test.users)
		if err != nil {
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
