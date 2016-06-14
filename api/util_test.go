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

func TestIsValidUserExternalID(t *testing.T) {
	testcases := map[string]struct {
		externalID string
		valid      bool
	}{
		"Case1": {
			externalID: "",
			valid:      false,
		},
		"Case2": {
			externalID: "*",
			valid:      false,
		},
		"Case3": {
			externalID: "/",
			valid:      false,
		},
		"Case4": {
			externalID: "something/",
			valid:      false,
		},
		"Case5": {
			externalID: "prefix*",
			valid:      false,
		},
		"Case6": {
			externalID: "pre*fix",
			valid:      false,
		},
		"Case7": {
			externalID: "comma,",
			valid:      false,
		},
		"Case8": {
			externalID: GetRandomString([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"), MAX_EXTERNAL_ID_LENGTH+1),
			valid:      false,
		},
		"Case9": {
			externalID: "good",
			valid:      true,
		},
		"Case10": {
			externalID: "123456",
			valid:      true,
		},
		"Case11": {
			externalID: "example-of-user123@email.com-that-is-valid",
			valid:      true,
		},
	}

	for x, testcase := range testcases {
		valid := IsValidUserExternalID(testcase.externalID)
		if valid != testcase.valid {
			t.Fatalf("Test %v failed. Received different values (wanted: %v / received: %v)",
				x, testcase.valid, valid)
		}
	}
}
