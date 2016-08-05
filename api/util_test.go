package api

import (
	"fmt"
	"testing"
)

func TestCreateUrn(t *testing.T) {
	testcases := map[string]struct {
		org         string
		resource    string
		path        string
		name        string
		expectedUrn string
	}{
		"OkCaseUserResource": {
			resource:    RESOURCE_USER,
			path:        "/mypath/",
			name:        "user",
			expectedUrn: "urn:iws:iam::user/mypath/user",
		},
		"OkCaseGroupResource": {
			resource:    RESOURCE_GROUP,
			org:         "org1",
			path:        "/mygrouppath/",
			name:        "group",
			expectedUrn: "urn:iws:iam:org1:group/mygrouppath/group",
		},
		"OkCasePolicyResource": {
			resource:    RESOURCE_POLICY,
			org:         "org1",
			path:        "/policypath/",
			name:        "policy",
			expectedUrn: "urn:iws:iam:org1:policy/policypath/policy",
		},
	}

	for x, testcase := range testcases {
		urn := CreateUrn(testcase.org, testcase.resource, testcase.path, testcase.name)
		checkMethodResponse(t, x, nil, nil, testcase.expectedUrn, urn)
	}
}

func TestGetUrnPrefix(t *testing.T) {
	testcases := map[string]struct {
		org         string
		resource    string
		path        string
		expectedUrn string
	}{
		"OkCaseUserResourcePrefix": {
			resource:    RESOURCE_USER,
			path:        "/mypath/",
			expectedUrn: "urn:iws:iam::user/mypath/*",
		},
		"OkCaseGroupResourcePrefix": {
			resource:    RESOURCE_GROUP,
			org:         "org1",
			path:        "/mygrouppath",
			expectedUrn: "urn:iws:iam:org1:group/mygrouppath*",
		},
		"OkCasePolicyResourcePrefix": {
			resource:    RESOURCE_POLICY,
			org:         "org1",
			path:        "/policypath/",
			expectedUrn: "urn:iws:iam:org1:policy/policypath/*",
		},
	}

	for x, testcase := range testcases {
		urn := GetUrnPrefix(testcase.org, testcase.resource, testcase.path)
		checkMethodResponse(t, x, nil, nil, testcase.expectedUrn, urn)
	}
}

func TestIsValidUserExternalID(t *testing.T) {
	testcases := map[string]struct {
		externalID string
		valid      bool
	}{
		"OkCaseEmpty": {
			externalID: "",
			valid:      false,
		},
		"OkCaseFullPrefix": {
			externalID: "*",
			valid:      false,
		},
		"OkCaseSlash1": {
			externalID: "/",
			valid:      false,
		},
		"OkCaseSlash2": {
			externalID: "something/",
			valid:      false,
		},
		"OkCasePrefix1": {
			externalID: "prefix*",
			valid:      false,
		},
		"OkCasePrefix2": {
			externalID: "pre*fix",
			valid:      false,
		},
		"OkCaseComma": {
			externalID: "comma,",
			valid:      false,
		},
		"OkCaseMaxLimitExceed": {
			externalID: GetRandomString([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"), MAX_EXTERNAL_ID_LENGTH+1),
			valid:      false,
		},
		"OkCaseLetter": {
			externalID: "good",
			valid:      true,
		},
		"OkCaseNumber": {
			externalID: "123456",
			valid:      true,
		},
		"OkCaseFullExample": {
			externalID: "example-of-user123@email.com-that-is-valid",
			valid:      true,
		},
	}

	for x, testcase := range testcases {
		valid := IsValidUserExternalID(testcase.externalID)
		checkMethodResponse(t, x, nil, nil, testcase.valid, valid)
	}
}

func TestIsValidOrg(t *testing.T) {
	testcases := map[string]struct {
		org   string
		valid bool
	}{
		"OkCaseEmpty": {
			org:   "",
			valid: false,
		},
		"OkCaseInvalidDot": {
			org:   "name.value",
			valid: false,
		},
		"OkCaseInvalid@": {
			org:   "@",
			valid: false,
		},
		"OkCaseInvalidComma": {
			org:   ",",
			valid: false,
		},
		"OkCaseInvalidSlash": {
			org:   "/",
			valid: false,
		},
		"OkCaseMaxLimitExceed": {
			org:   GetRandomString([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"), MAX_NAME_LENGTH+1),
			valid: false,
		},
		"OkCase": {
			org:   "validName",
			valid: true,
		},
	}

	for x, testcase := range testcases {
		valid := IsValidOrg(testcase.org)
		checkMethodResponse(t, x, nil, nil, testcase.valid, valid)
	}
}

func TestIsValidName(t *testing.T) {
	testcases := map[string]struct {
		name  string
		valid bool
	}{
		"OkCaseEmpty": {
			name:  "",
			valid: false,
		},
		"OkCaseInvalidDot": {
			name:  "name.value",
			valid: false,
		},
		"OkCaseInvalid@": {
			name:  "@",
			valid: false,
		},
		"OkCaseInvalidComma": {
			name:  ",",
			valid: false,
		},
		"OkCaseInvalidSlash": {
			name:  "/",
			valid: false,
		},
		"OkCaseMaxLimitExceed": {
			name:  GetRandomString([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"), MAX_NAME_LENGTH+1),
			valid: false,
		},
		"OkCase": {
			name:  "validName",
			valid: true,
		},
	}

	for x, testcase := range testcases {
		valid := IsValidName(testcase.name)
		checkMethodResponse(t, x, nil, nil, testcase.valid, valid)
	}
}

func TestIsValidPath(t *testing.T) {
	testcases := map[string]struct {
		path  string
		valid bool
	}{
		"OkCaseEmpty": {
			path:  "",
			valid: false,
		},
		"OkCaseMalformedPath1": {
			path:  "/path",
			valid: false,
		},
		"OkCaseMalformedPath2": {
			path:  "path/",
			valid: false,
		},
		"OkCaseInvalidComma": {
			path:  ",",
			valid: false,
		},
		"OkCaseInvalidPrefix": {
			path:  "*",
			valid: false,
		},
		"OkCaseInvalidPrefix2": {
			path:  "/*",
			valid: false,
		},
		"OkCaseInvalidPrefix3": {
			path:  "*/",
			valid: false,
		},
		"OkCaseInvalidPrefix4": {
			path:  "/*/",
			valid: false,
		},
		"OkCaseInvalidDot": {
			path:  "path.value",
			valid: false,
		},
		"OkCaseMaxLimitExceed": {
			path:  "/" + GetRandomString([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"), MAX_PATH_LENGTH-1) + "/",
			valid: false,
		},
		"OkCaseRoot": {
			path:  "/",
			valid: true,
		},
		"OkCaseWithPath": {
			path:  "/path/",
			valid: true,
		},
	}

	for x, testcase := range testcases {
		valid := IsValidPath(testcase.path)
		checkMethodResponse(t, x, nil, nil, testcase.valid, valid)
	}
}

func TestIsValidEffect(t *testing.T) {
	testcases := map[string]struct {
		// Method args
		effect string
		// Expected results
		wantError error
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
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid effect: other - Only 'allow' and 'deny' accepted",
			},
		},
	}

	for x, testcase := range testcases {
		err := IsValidEffect(testcase.effect)
		checkMethodResponse(t, x, testcase.wantError, err, nil, nil)
	}
}

func TestAreValidActions(t *testing.T) {
	randomString := GetRandomString([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ:*"), MAX_ACTION_LENGTH+1)
	testcases := map[string]struct {
		// Method args
		actions []string
		// Expected results
		wantError error
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
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in action: iam:",
			},
		},
		"ErrorCaseMalformedAction2": {
			actions: []string{
				"*",
			},
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in action: *",
			},
		},
		"ErrorCaseInvalidCharacters": {
			actions: []string{
				"iam:**",
			},
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in action: iam:**",
			},
		},
		"ErrorCaseInvalidCharacters2": {
			actions: []string{
				"iam::",
			},
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in action: iam::",
			},
		},
		"ErrorCaseMaxLengthExceeded": {
			actions: []string{
				randomString,
			},
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: fmt.Sprintf("No regex match in action: %v", randomString),
			},
		},
	}

	for x, testcase := range testcases {
		err := AreValidActions(testcase.actions)
		checkMethodResponse(t, x, testcase.wantError, err, nil, nil)
	}
}

func TestAreValidStatements(t *testing.T) {
	testcases := map[string]struct {
		// Method args
		Statements *[]Statement
		// Expected results
		wantError error
	}{
		"OKCase": {
			Statements: &[]Statement{
				{
					Effect: "allow",
					Actions: []string{
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
				{
					Effect: "FAILallowZ",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid effect: FAILallowZ - Only 'allow' and 'deny' accepted",
			},
		},
		"ErrorCaseInvalidAction": {
			Statements: &[]Statement{
				{
					Effect: "allow",
					Actions: []string{
						"fail***",
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in action: fail***",
			},
		},
		"ErrorCaseInvalidResource": {
			Statements: &[]Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/***"),
					},
				},
			},
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in resource: urn:iws:iam::user/path/****",
			},
		},
	}

	for x, testcase := range testcases {
		err := AreValidStatements(testcase.Statements)
		checkMethodResponse(t, x, testcase.wantError, err, nil, nil)
	}
}

func TestAreValidResources(t *testing.T) {
	testcases := map[string]struct {
		// Method args
		Resources []string
		// Expected results
		wantError error
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
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in resource: fail",
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
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in resource: fail:asd",
			},
		},
		"ErrorCase2block": {
			Resources: []string{
				"urn:fail",
			},
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in resource: urn:fail",
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
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in resource: urn:iws***:something",
			},
		},
		"ErrorCase3block": {
			Resources: []string{
				"urn:iws:fail",
			},
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in resource: urn:iws:fail",
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
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in resource: urn:iws:some***:fail",
			},
		},
		"ErrorCase4block": {
			Resources: []string{
				"urn:iws:iam:fail",
			},
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in resource: urn:iws:iam:fail",
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
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in resource: urn:iws:iam:some***:fail",
			},
		},
		"ErrorCase5block": {
			Resources: []string{
				"urn:iws:iam:org1:fail**!^_#",
			},
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in resource: urn:iws:iam:org1:fail**!^_#",
			},
		},
		"ErrorCaseBadResource": {
			Resources: []string{
				"urn:iws:iam:org1:fail:fail:fail",
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid resource definition: urn:iws:iam:org1:fail:fail:fail",
			},
		},
	}

	for x, testcase := range testcases {
		err := AreValidResources(testcase.Resources)
		checkMethodResponse(t, x, testcase.wantError, err, nil, nil)
	}
}
