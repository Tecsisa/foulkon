package postgresql

import (
	"testing"
	"time"

	"github.com/Tecsisa/foulkon/api"
	"github.com/stretchr/testify/assert"
)

func TestGroupUser_GetUser(t *testing.T) {
	testcases := map[string]struct {
		relation       GroupUser
		expectedResult *api.User
	}{
		"OkCase": {
			relation: GroupUser{
				User: &api.User{
					ID:         "ID",
					ExternalID: "user1",
					Path:       "Path",
					Urn:        "urn",
				},
			},
			expectedResult: &api.User{
				ID:         "ID",
				ExternalID: "user1",
				Path:       "Path",
				Urn:        "urn",
			},
		},
	}

	for x, testcase := range testcases {
		assert.Equal(t, testcase.expectedResult, testcase.relation.GetUser(), "Error in test case %v", x)
	}
}

func TestGroupUser_GetGroup(t *testing.T) {
	testcases := map[string]struct {
		relation       GroupUser
		expectedResult *api.Group
	}{
		"OkCase": {
			relation: GroupUser{
				Group: &api.Group{
					ID:   "ID",
					Name: "group1",
					Org:  "org1",
					Path: "Path",
					Urn:  "urn",
				},
			},
			expectedResult: &api.Group{
				ID:   "ID",
				Name: "group1",
				Org:  "org1",
				Path: "Path",
				Urn:  "urn",
			},
		},
	}

	for x, testcase := range testcases {
		assert.Equal(t, testcase.expectedResult, testcase.relation.GetGroup(), "Error in test case %v", x)
	}
}

func TestGroupUser_GetDate(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		relation       GroupUser
		expectedResult time.Time
	}{
		"OkCase": {
			relation: GroupUser{
				CreateAt: now,
			},
			expectedResult: now,
		},
	}

	for x, testcase := range testcases {
		assert.Equal(t, testcase.expectedResult, testcase.relation.GetDate(), "Error in test case %v", x)
	}
}

func TestPolicyGroup_GetPolicy(t *testing.T) {
	testcases := map[string]struct {
		relation       PolicyGroup
		expectedResult *api.Policy
	}{
		"OkCase": {
			relation: PolicyGroup{
				Policy: &api.Policy{
					ID:   "ID",
					Name: "policy1",
					Org:  "org1",
					Path: "Path",
					Urn:  "urn",
				},
			},
			expectedResult: &api.Policy{
				ID:   "ID",
				Name: "policy1",
				Org:  "org1",
				Path: "Path",
				Urn:  "urn",
			},
		},
	}

	for x, testcase := range testcases {
		assert.Equal(t, testcase.expectedResult, testcase.relation.GetPolicy(), "Error in test case %v", x)
	}
}

func TestPolicyGroup_GetGroup(t *testing.T) {
	testcases := map[string]struct {
		relation       PolicyGroup
		expectedResult *api.Group
	}{
		"OkCase": {
			relation: PolicyGroup{
				Group: &api.Group{
					ID:   "ID",
					Name: "group1",
					Org:  "org1",
					Path: "Path",
					Urn:  "urn",
				},
			},
			expectedResult: &api.Group{
				ID:   "ID",
				Name: "group1",
				Org:  "org1",
				Path: "Path",
				Urn:  "urn",
			},
		},
	}

	for x, testcase := range testcases {
		assert.Equal(t, testcase.expectedResult, testcase.relation.GetGroup(), "Error in test case %v", x)
	}
}

func TestPolicyGroup_GetDate(t *testing.T) {
	now := time.Now().UTC()
	testcases := map[string]struct {
		relation       PolicyGroup
		expectedResult time.Time
	}{
		"OkCase": {
			relation: PolicyGroup{
				CreateAt: now,
			},
			expectedResult: now,
		},
	}

	for x, testcase := range testcases {
		assert.Equal(t, testcase.expectedResult, testcase.relation.GetDate(), "Error in test case %v", x)
	}
}
