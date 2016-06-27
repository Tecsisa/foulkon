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
			t.Errorf("Test %v failed. Received different urns (wanted: %v / received: %v)",
				x, testcase.expectedUrn, urn)
			continue
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
			t.Errorf("Test %v failed. Received different urns (wanted: %v / received: %v)",
				x, testcase.expectedUrn, urn)
			continue
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
			t.Errorf("Test %v failed. Received different values (wanted: %v / received: %v)",
				x, testcase.valid, valid)
			continue
		}
	}
}

func TestIsValidName(t *testing.T) {
	testcases := map[string]struct {
		name  string
		valid bool
	}{
		"Case1": {
			name:  "",
			valid: false,
		},
		"Case2": {
			name:  "name.value",
			valid: false,
		},
		"Case3": {
			name:  "@",
			valid: false,
		},
		"Case4": {
			name:  ",",
			valid: false,
		},
		"Case5": {
			name:  "/",
			valid: false,
		},
		"Case6": {
			name:  GetRandomString([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"), MAX_NAME_LENGTH+1),
			valid: false,
		},
		"Case7": {
			name:  "validName",
			valid: true,
		},
	}

	for x, testcase := range testcases {
		valid := IsValidName(testcase.name)
		if valid != testcase.valid {
			t.Errorf("Test %v failed. Received different values (wanted: %v / received: %v)",
				x, testcase.valid, valid)
			continue
		}
	}
}

func TestIsValidPath(t *testing.T) {
	testcases := map[string]struct {
		path  string
		valid bool
	}{
		"Case1": {
			path:  "",
			valid: false,
		},
		"Case2": {
			path:  "/path",
			valid: false,
		},
		"Case3": {
			path:  "path/",
			valid: false,
		},
		"Case4": {
			path:  ",",
			valid: false,
		},
		"Case5": {
			path:  "*",
			valid: false,
		},
		"Case6": {
			path:  "/*",
			valid: false,
		},
		"Case7": {
			path:  "*/",
			valid: false,
		},
		"Case8": {
			path:  "/*/",
			valid: false,
		},
		"Case9": {
			path:  "path.value",
			valid: false,
		},
		"Case10": {
			path:  "/" + GetRandomString([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"), MAX_PATH_LENGTH-1) + "/",
			valid: false,
		},
		"Case11": {
			path:  "/",
			valid: true,
		},
		"Case12": {
			path:  "/path/",
			valid: true,
		},
	}

	for x, testcase := range testcases {
		valid := IsValidPath(testcase.path)
		if valid != testcase.valid {
			t.Errorf("Test %v failed. Received different values (wanted: %v / received: %v)",
				x, testcase.valid, valid)
			continue
		}
	}
}

func TestIsValidEffect(t *testing.T) {
	testcases := map[string]struct {
		// Method args
		effect string
		// Expected results
		wantError *Error
	}{
		"OKCaseAllow": {
			effect: "allow",
		},
		"OKCaseDeny": {
			effect: "deny",
		},
		"ErrorCaseInvalidEffect": {
			effect: "other",
			wantError: &Error{
				Code: REGEX_NO_MATCH,
			},
		},
	}

	for x, testcase := range testcases {
		err := IsValidEffect(testcase.effect)
		if testcase.wantError != nil {
			apiError, ok := err.(*Error)
			if !ok || apiError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", x, err)
				continue
			}
			if apiError.Code != testcase.wantError.Code {
				t.Errorf("Test %v failed. Got error %v, expected %v", x, apiError, testcase.wantError.Code)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Error: %v", x, err)
				continue
			}
		}
	}
}

func TestIsValidAction(t *testing.T) {
	testcases := map[string]struct {
		// Method args
		actions []string
		// Expected results
		wantError *Error
	}{
		"OKCaseValidAction": {
			actions: []string{
				"iam:operation",
				"iam:*",
			},
		},
		"ErrorCaseMalformedAction": {
			actions: []string{
				"iam:",
			},
			wantError: &Error{
				Code: REGEX_NO_MATCH,
			},
		},
		"ErrorCaseMalformedAction2": {
			actions: []string{
				"*",
			},
			wantError: &Error{
				Code: REGEX_NO_MATCH,
			},
		},
		"ErrorCaseInvalidCharacters": {
			actions: []string{
				"iam:**",
			},
			wantError: &Error{
				Code: REGEX_NO_MATCH,
			},
		},
		"ErrorCaseInvalidCharacters2": {
			actions: []string{
				"iam::",
			},
			wantError: &Error{
				Code: REGEX_NO_MATCH,
			},
		},
		"ErrorCaseMaxLengthExceeded": {
			actions: []string{
				GetRandomString([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ:*"), MAX_ACTION_LENGTH+1),
			},
			wantError: &Error{
				Code: REGEX_NO_MATCH,
			},
		},
	}

	for x, testcase := range testcases {
		err := IsValidAction(testcase.actions)
		if testcase.wantError != nil {
			apiError, ok := err.(*Error)
			if !ok || apiError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", x, err)
				continue
			}
			if apiError.Code != testcase.wantError.Code {
				t.Errorf("Test %v failed. Got error %v, expected %v", x, apiError, testcase.wantError.Code)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Error: %v", x, err)
				continue
			}
		}
	}
}

func TestIsValidStatement(t *testing.T) {
	testcases := map[string]struct {
		// Method args
		Statements *[]Statement
		// Expected results
		wantError *Error
	}{
		"OKCase": {
			Statements: &[]Statement{
				Statement{
					Effect: "allow",
					Action: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
		},
		"ErrorCaseInvalidEffect": {
			Statements: &[]Statement{
				Statement{
					Effect: "FAILallowZ",
					Action: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			wantError: &Error{
				Code: REGEX_NO_MATCH,
			},
		},
		"ErrorCaseInvalidAction": {
			Statements: &[]Statement{
				Statement{
					Effect: "allow",
					Action: []string{
						"fail***",
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			wantError: &Error{
				Code: REGEX_NO_MATCH,
			},
		},
		"ErrorCaseInvalidResource": {
			Statements: &[]Statement{
				Statement{
					Effect: "allow",
					Action: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/***"),
					},
				},
			},
			wantError: &Error{
				Code: REGEX_NO_MATCH,
			},
		},
	}

	for x, testcase := range testcases {
		err := IsValidStatement(testcase.Statements)
		if testcase.wantError != nil {
			apiError, ok := err.(*Error)
			if !ok || apiError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", x, err)
				continue
			}
			if apiError.Code != testcase.wantError.Code {
				t.Errorf("Test %v failed. Got error %v, expected %v", x, apiError, testcase.wantError.Code)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Error: %v", x, err)
				continue
			}
		}
	}
}

func TestIsValidResources(t *testing.T) {
	testcases := map[string]struct {
		// Method args
		Resources []string
		// Expected results
		wantError *Error
	}{
		"OKCase1block": {
			Resources: []string{
				"*",
			},
		},
		"ErrorCase1block": {
			Resources: []string{
				"fail",
			},
			wantError: &Error{
				Code: REGEX_NO_MATCH,
			},
		},
		"OKCase2block": {
			Resources: []string{
				"urn:*",
			},
		},
		"ErrorCase2blockNotUrn": {
			Resources: []string{
				"fail:asd",
			},
			wantError: &Error{
				Code: REGEX_NO_MATCH,
			},
		},
		"ErrorCase2block": {
			Resources: []string{
				"urn:fail",
			},
			wantError: &Error{
				Code: REGEX_NO_MATCH,
			},
		},
		"OKCase3block": {
			Resources: []string{
				"urn:iws:*",
			},
		},
		"ErrorCase3blockBadString": {
			Resources: []string{
				"urn:iws***:something",
			},
			wantError: &Error{
				Code: REGEX_NO_MATCH,
			},
		},
		"ErrorCase3block": {
			Resources: []string{
				"urn:iws:fail",
			},
			wantError: &Error{
				Code: REGEX_NO_MATCH,
			},
		},
		"OKCase4block": {
			Resources: []string{
				"urn:iws:iam:*",
			},
		},
		"ErrorCase4blockBadString": {
			Resources: []string{
				"urn:iws:some***:fail",
			},
			wantError: &Error{
				Code: REGEX_NO_MATCH,
			},
		},
		"ErrorCase4block": {
			Resources: []string{
				"urn:iws:iam:fail",
			},
			wantError: &Error{
				Code: REGEX_NO_MATCH,
			},
		},
		"OKCase5block": {
			Resources: []string{
				"urn:iws:iam::sysadmins/user1",
			},
		},
		"ErrorCase5blockBadString": {
			Resources: []string{
				"urn:iws:iam:some***:fail",
			},
			wantError: &Error{
				Code: REGEX_NO_MATCH,
			},
		},
		"ErrorCase5block": {
			Resources: []string{
				"urn:iws:iam:org1:fail**!^_#",
			},
			wantError: &Error{
				Code: REGEX_NO_MATCH,
			},
		},
		"ErrorCaseBadResource": {
			Resources: []string{
				"urn:iws:iam:org1:fail:fail:fail",
			},
			wantError: &Error{
				Code: INVALID_PARAMETER_ERROR,
			},
		},
	}

	for x, testcase := range testcases {
		err := IsValidResources(testcase.Resources)
		if testcase.wantError != nil {
			apiError, ok := err.(*Error)
			if !ok || apiError == nil {
				t.Errorf("Test %v failed. Unexpected data retrieved from error: %v", x, err)
				continue
			}
			if apiError.Code != testcase.wantError.Code {
				t.Errorf("Test %v failed. Got error %v, expected %v", x, apiError, testcase.wantError.Code)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Test %v failed. Error: %v", x, err)
				continue
			}
		}
	}
}
