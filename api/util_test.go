package api

import "testing"

func TestCreateUrn(t *testing.T) {
	testcases := map[string]struct {
		org         string
		resource    string
		path        string
		name        string
		expectedUrn string
	}{
		"UserResource": {
			resource:    RESOURCE_USER,
			path:        "/mypath/",
			name:        "user",
			expectedUrn: "urn:iws:iam::user/mypath/user",
		},
		"GroupResource": {
			resource:    RESOURCE_GROUP,
			org:         "org1",
			path:        "/mygrouppath/",
			name:        "group",
			expectedUrn: "urn:iws:iam:org1:group/mygrouppath/group",
		},
		"PolicyResource": {
			resource:    RESOURCE_POLICY,
			org:         "org1",
			path:        "/policypath/",
			name:        "policy",
			expectedUrn: "urn:iws:iam:org1:policy/policypath/policy",
		},
	}

	for x, testcase := range testcases {
		urn := CreateUrn(testcase.org, testcase.resource, testcase.path, testcase.name)
		if urn != testcase.expectedUrn {
			t.Fatalf("Test %v failed. Received different urns (wanted: %v / received: %v)",
				x, testcase.expectedUrn, urn)
		}
	}
}

func TestGetUrnPrefix(t *testing.T) {
	testcases := map[string]struct {
		org         string
		resource    string
		path        string
		expectedUrn string
	}{
		"UserResourcePrefix": {
			resource:    RESOURCE_USER,
			path:        "/mypath/",
			expectedUrn: "urn:iws:iam::user/mypath/*",
		},
		"GroupResourcePrefix": {
			resource:    RESOURCE_GROUP,
			org:         "org1",
			path:        "/mygrouppath",
			expectedUrn: "urn:iws:iam:org1:group/mygrouppath*",
		},
		"PolicyResourcePrefix": {
			resource:    RESOURCE_POLICY,
			org:         "org1",
			path:        "/policypath/",
			expectedUrn: "urn:iws:iam:org1:policy/policypath/*",
		},
	}

	for x, testcase := range testcases {
		urn := GetUrnPrefix(testcase.org, testcase.resource, testcase.path)
		if urn != testcase.expectedUrn {
			t.Fatalf("Test %v failed. Received different urns (wanted: %v / received: %v)",
				x, testcase.expectedUrn, urn)
		}
	}
}
