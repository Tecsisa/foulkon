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
			externalID: getRandomString([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"), MAX_EXTERNAL_ID_LENGTH+1),
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
			org:   getRandomString([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"), MAX_NAME_LENGTH+1),
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
			name:  getRandomString([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"), MAX_NAME_LENGTH+1),
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
			path:  "/" + getRandomString([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"), MAX_PATH_LENGTH-1) + "/",
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
	randomString := getRandomString([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ:*"), MAX_ACTION_LENGTH+1)
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
		Resources    []string
		resourceType string
		// Expected results
		wantError error
	}{
		"OKCase1block": {
			Resources: []string{
				"*",
			},
			resourceType: RESOURCE_IAM,
		},
		"ErrorCase1block": {
			Resources: []string{
				"fail",
			},
			resourceType: RESOURCE_IAM,
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in resource: fail",
			},
		},
		"OKCase2block": {
			Resources: []string{
				"urn:*",
			},
			resourceType: RESOURCE_IAM,
		},
		"ErrorCase2blockNotUrn": {
			Resources: []string{
				"fail:asd",
			},
			resourceType: RESOURCE_IAM,
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in resource: fail:asd",
			},
		},
		"ErrorCase2block": {
			Resources: []string{
				"urn:fail",
			},
			resourceType: RESOURCE_IAM,
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in resource: urn:fail",
			},
		},
		"OKCase3block": {
			Resources: []string{
				"urn:iws:*",
			},
			resourceType: RESOURCE_IAM,
		},
		"ErrorCase3blockBadString": {
			Resources: []string{
				"urn:iws***:something",
			},
			resourceType: RESOURCE_IAM,
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in resource: urn:iws***:something",
			},
		},
		"ErrorCase3block": {
			Resources: []string{
				"urn:iws:fail",
			},
			resourceType: RESOURCE_IAM,
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in resource: urn:iws:fail",
			},
		},
		"OKCase4block": {
			Resources: []string{
				"urn:iws:iam:*",
			},
			resourceType: RESOURCE_IAM,
		},
		"ErrorCase4blockBadString": {
			Resources: []string{
				"urn:iws:some***:fail",
			},
			resourceType: RESOURCE_IAM,
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in resource: urn:iws:some***:fail",
			},
		},
		"ErrorCase4block": {
			Resources: []string{
				"urn:iws:iam:fail",
			},
			resourceType: RESOURCE_IAM,
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in resource: urn:iws:iam:fail",
			},
		},
		"OKCase5block": {
			Resources: []string{
				"urn:iws:iam::sysadmins/user1",
			},
			resourceType: RESOURCE_IAM,
		},
		"OKCase5blockExternal": {
			Resources: []string{
				"urn:ews:exam:inst:sysadmins/{admin}",
			},
			resourceType: RESOURCE_EXTERNAL,
		},
		"ErrorCase5blockBadString": {
			Resources: []string{
				"urn:iws:iam:some***:fail",
			},
			resourceType: RESOURCE_IAM,
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in resource: urn:iws:iam:some***:fail",
			},
		},
		"ErrorCase5blockExternal": {
			Resources: []string{
				"urn:ews:exam:inst:sysadmins/{admin",
			},
			resourceType: RESOURCE_EXTERNAL,
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in resource: urn:ews:exam:inst:sysadmins/{admin",
			},
		},
		"ErrorCase5block": {
			Resources: []string{
				"urn:iws:iam:org1:fail**!^_#",
			},
			resourceType: RESOURCE_IAM,
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in resource: urn:iws:iam:org1:fail**!^_#",
			},
		},
		"ErrorCaseBadResource": {
			Resources: []string{
				"urn:iws:iam:org1:fail:fail:fail",
			},
			resourceType: RESOURCE_IAM,
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid resource definition: urn:iws:iam:org1:fail:fail:fail",
			},
		},
	}

	for x, testcase := range testcases {
		err := AreValidResources(testcase.Resources, testcase.resourceType)
		checkMethodResponse(t, x, testcase.wantError, err, nil, nil)
	}
}

func TestValidateFilter(t *testing.T) {
	testcases := map[string]struct {
		// Method args
		filter              *Filter
		OrderByValidColumns []string
		// Expected results
		wantError error
	}{
		"OKCaseAllow": {
			filter: &Filter{
				ExternalID: "123",
				GroupName:  "grp",
				Org:        "org",
				PolicyName: "p1",
				Offset:     2,
				OrderBy:    "name-desc",
			},
			OrderByValidColumns: []string{"name", "test"},
		},
		"ErrorCaseInvalidOrg": {
			filter: &Filter{
				ExternalID: "123",
				GroupName:  "grp",
				Org:        "!*@~#",
				PolicyName: "p1",
				Limit:      10,
				Offset:     2,
				PathPrefix: "/path/",
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org !*@~#",
			},
		},
		"ErrorCaseInvalidPath": {
			filter: &Filter{
				ExternalID: "123",
				GroupName:  "grp",
				Org:        "org",
				PolicyName: "p1",
				Limit:      10,
				Offset:     2,
				PathPrefix: "fail",
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: pathPrefix fail",
			},
		},
		"ErrorCaseInvalidLimit": {
			filter: &Filter{
				ExternalID: "123",
				GroupName:  "grp",
				Org:        "org",
				PolicyName: "p1",
				Limit:      5000,
				Offset:     2,
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: fmt.Sprintf("Invalid parameter: limit 5000, max limit allowed: %v", MAX_LIMIT_SIZE),
			},
		},
		"ErrorCaseInvalidGroupName": {
			filter: &Filter{
				ExternalID: "123",
				GroupName:  "#@!^*",
				Org:        "org",
				PolicyName: "p1",
				Limit:      5000,
				Offset:     2,
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: group #@!^*",
			},
		},
		"ErrorCaseInvalidExtID": {
			filter: &Filter{
				ExternalID: "#@!^*",
				GroupName:  "grp",
				Org:        "org",
				PolicyName: "p1",
				Limit:      5000,
				Offset:     2,
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: externalID #@!^*",
			},
		},
		"ErrorCaseInvalidPolicyName": {
			filter: &Filter{
				ExternalID: "123",
				GroupName:  "grp",
				Org:        "org",
				PolicyName: "#@!^*",
				Limit:      5000,
				Offset:     2,
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: policy #@!^*",
			},
		},
		"ErrorCaseInvalidOrderBy": {
			filter: &Filter{
				ExternalID: "123",
				GroupName:  "grp",
				Org:        "org",
				PolicyName: "asd",
				Limit:      10,
				Offset:     2,
				OrderBy:    "fail-fail",
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: OrderBy fail-fail",
			},
		},
		"ErrorCaseInvalidOrderByColumn": {
			filter: &Filter{
				ExternalID: "123",
				GroupName:  "grp",
				Org:        "org",
				PolicyName: "asd",
				Limit:      10,
				Offset:     2,
				OrderBy:    "xxx-asc",
			},
			OrderByValidColumns: []string{"val1", "val2"},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: OrderBy column xxx",
			},
		},
	}

	for x, testcase := range testcases {
		err := validateFilter(testcase.filter, testcase.OrderByValidColumns)
		checkMethodResponse(t, x, testcase.wantError, err, nil, nil)
	}
}

func TestIsValidProxyResource(t *testing.T) {
	testcases := map[string]struct {
		// Method args
		resource *ResourceEntity
		// Expected results
		wantError error
	}{
		"OKCase": {
			resource: &ResourceEntity{
				Host:   "http://host.com",
				Path:   "/path",
				Method: "GET",
				Urn:    "urn:ews:example:instance1:resource/get",
				Action: "action",
			},
		},
		"ErrorCaseInvalidHost": {
			resource: &ResourceEntity{
				Host: "~32&",
			},
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in parameter: ~32&",
			},
		},
		"ErrorCaseInvalidPath": {
			resource: &ResourceEntity{
				Host: "http://host.com",
				Path: "invalid",
			},
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in parameter: invalid",
			},
		},
		"ErrorCaseInvalidMethod": {
			resource: &ResourceEntity{
				Host:   "http://host.com",
				Path:   "/path",
				Method: "INVALID",
			},
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in parameter: INVALID",
			},
		},
		"ErrorCaseInvalidUrn": {
			resource: &ResourceEntity{
				Host:   "http://host.com",
				Path:   "/path",
				Method: "GET",
				Urn:    "~~&",
			},
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in resource: ~~&",
			},
		},
		"ErrorCaseInvalidAction": {
			resource: &ResourceEntity{
				Host:   "http://host.com",
				Path:   "/path",
				Method: "GET",
				Urn:    "urn:ews:example:instance1:resource/get",
				Action: "iam:",
			},
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in action: iam:",
			},
		},
	}

	for x, testcase := range testcases {
		err := IsValidProxyResource(testcase.resource)
		checkMethodResponse(t, x, testcase.wantError, err, nil, nil)
	}
}
