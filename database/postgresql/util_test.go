package postgresql

import (
	"testing"
	"time"

	"github.com/Tecsisa/foulkon/api"
	"github.com/kylelemons/godebug/pretty"
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
		if diff := pretty.Compare(testcase.expectedResult, testcase.relation.GetUser()); diff != "" {
			t.Errorf("Test %v failed. Received different messages (wanted / received) %v",
				x, diff)
			continue
		}
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
		if diff := pretty.Compare(testcase.expectedResult, testcase.relation.GetGroup()); diff != "" {
			t.Errorf("Test %v failed. Received different messages (wanted / received) %v",
				x, diff)
			continue
		}
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
		if diff := pretty.Compare(testcase.expectedResult, testcase.relation.GetDate()); diff != "" {
			t.Errorf("Test %v failed. Received different messages (wanted / received) %v",
				x, diff)
			continue
		}
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
		if diff := pretty.Compare(testcase.expectedResult, testcase.relation.GetPolicy()); diff != "" {
			t.Errorf("Test %v failed. Received different messages (wanted / received) %v",
				x, diff)
			continue
		}
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
		if diff := pretty.Compare(testcase.expectedResult, testcase.relation.GetGroup()); diff != "" {
			t.Errorf("Test %v failed. Received different messages (wanted / received) %v",
				x, diff)
			continue
		}
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
		if diff := pretty.Compare(testcase.expectedResult, testcase.relation.GetDate()); diff != "" {
			t.Errorf("Test %v failed. Received different messages (wanted / received) %v",
				x, diff)
			continue
		}
	}
}
