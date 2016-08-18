package api

import (
	"testing"

	"github.com/tecsisa/foulkon/database"
)

func TestAuthAPI_AddPolicy(t *testing.T) {
	testcases := map[string]struct {
		requestInfo RequestInfo
		org         string
		policyName  string
		path        string
		statements  []Statement

		getGroupsByUserIDResult   []Group
		getAttachedPoliciesResult []Policy
		getUserByExternalIDResult *User

		addPolicyMethodResult       *Policy
		getPolicyByNameMethodResult *Policy
		wantError                   error

		getPolicyByNameMethodErr error
		addPolicyMethodErr       error
	}{
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			getPolicyByNameMethodErr: &database.Error{
				Code: database.POLICY_NOT_FOUND,
			},
			addPolicyMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
		},
		"ErrorCasePolicyAlreadyExists": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
			wantError: &Error{
				Code:    POLICY_ALREADY_EXIST,
				Message: "Unable to create policy, policy with org 123 and name test already exist",
			},
		},
		"ErrorCaseBadName": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "**!^#~",
			path:       "/path/",
			statements: []Statement{
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
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name **!^#~",
			},
		},
		"ErrorCaseEmptyActions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "p1",
			path:       "/path/",
			statements: []Statement{
				{
					Effect:  "allow",
					Actions: []string{},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path/"),
					},
				},
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Empty actions",
			},
		},
		"ErrorCaseEmptyResources": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "p1",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{},
				},
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Empty resources",
			},
		},
		"ErrorCaseBadOrgName": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "**!^#~",
			policyName: "p1",
			path:       "/path/",
			statements: []Statement{
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
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org **!^#~",
			},
		},
		"ErrorCaseBadPath": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/**!^#~path/",
			statements: []Statement{
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
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: path /**!^#~path/",
			},
		},
		"ErrorCaseBadStatement": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
				{
					Effect: "idufhefmfcasfluhf",
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
				Message: "Invalid effect: idufhefmfcasfluhf - Only 'allow' and 'deny' accepted",
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/1/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:123:policy/path/test",
			},
		},
		"ErrorCaseDenyResource": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "example",
			policyName: "test",
			path:       "/path/",
			getPolicyByNameMethodErr: &database.Error{
				Code: database.POLICY_NOT_FOUND,
			},
			statements: []Statement{
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
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/1/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policy",
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/"),
							},
						},
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_CREATE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/"),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								POLICY_ACTION_CREATE_POLICY,
							},
							Resources: []string{
								CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
							},
						},
					},
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam:example:policy/path/test",
			},
		},
		"ErrorCaseAddPolicyDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			getPolicyByNameMethodErr: &database.Error{
				Code: database.POLICY_NOT_FOUND,
			},
			addPolicyMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
		"ErrorCaseGetPolicyDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			getPolicyByNameMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
	}

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for x, testcase := range testcases {
		testRepo.ArgsOut[AddPolicyMethod][0] = testcase.addPolicyMethodResult
		testRepo.ArgsOut[AddPolicyMethod][1] = testcase.addPolicyMethodErr
		testRepo.ArgsOut[GetPolicyByNameMethod][0] = testcase.getPolicyByNameMethodResult
		testRepo.ArgsOut[GetPolicyByNameMethod][1] = testcase.getPolicyByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		policy, err := testAPI.AddPolicy(testcase.requestInfo, testcase.policyName, testcase.path, testcase.org, testcase.statements)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.addPolicyMethodResult, policy)
	}
}

func TestAuthAPI_GetPolicyByName(t *testing.T) {
	testcases := map[string]struct {
		requestInfo RequestInfo
		org         string
		policyName  string

		getGroupsByUserIDResult   []Group
		getAttachedPoliciesResult []Policy
		getUserByExternalIDResult *User

		getPolicyByNameMethodResult *Policy
		wantError                   error

		getPolicyByNameMethodErr error
	}{
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
		},
		"ErrorCaseInternalError": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			getPolicyByNameMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
		"ErrorCaseBadPolicyName": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "~#**!",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name ~#**!",
			},
		},
		"ErrorCaseBadOrgName": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "~#**!",
			policyName: "p1",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org ~#**!",
			},
		},
		"ErrorCasePolicyNotFound": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			getPolicyByNameMethodErr: &database.Error{
				Code: database.POLICY_NOT_FOUND,
			},
			wantError: &Error{
				Code: POLICY_BY_ORG_AND_NAME_NOT_FOUND,
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "example",
			policyName: "policyUser",
			getPolicyByNameMethodResult: &Policy{
				ID:   "POLICY-USER-ID",
				Name: "policyUser",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
				Statements: &[]Statement{
					{
						Effect: "deny",
						Actions: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_POLICY, "/path/"),
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/1/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam:example:policy/path/policyUser",
			},
		},
		"ErrorCaseDenyResourceErr": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "example",
			policyName: "policyUser",
			getPolicyByNameMethodResult: &Policy{
				ID:   "POLICY-USER-ID",
				Name: "policyUser",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
				Statements: &[]Statement{
					{
						Effect: "deny",
						Actions: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_POLICY, "/path/"),
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/1/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
							},
						},
					},
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam:example:policy/path/policyUser",
			},
		},
	}

	for x, testcase := range testcases {

		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetPolicyByNameMethod][0] = testcase.getPolicyByNameMethodResult
		testRepo.ArgsOut[GetPolicyByNameMethod][1] = testcase.getPolicyByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		policy, err := testAPI.GetPolicyByName(testcase.requestInfo, testcase.org, testcase.policyName)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.getPolicyByNameMethodResult, policy)
	}
}

func TestAuthAPI_ListPolicies(t *testing.T) {
	testcases := map[string]struct {
		requestInfo RequestInfo
		org         string
		pathPrefix  string

		expectedPolicies []PolicyIdentity

		getGroupsByUserIDResult   []Group
		getAttachedPoliciesResult []Policy
		getUserByExternalIDResult *User
		getUserByExternalIDErr    error

		getPoliciesFilteredMethodResult []Policy
		getPoliciesFilteredMethodErr    error

		wantError error
	}{
		"OkCaseAdmin": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			pathPrefix: "/",
			expectedPolicies: []PolicyIdentity{
				{
					Org:  "example",
					Name: "policyAllowed",
				},
				{
					Org:  "example",
					Name: "policyDenied",
				},
			},
			getPoliciesFilteredMethodResult: []Policy{
				{
					ID:   "PolicyAllowed",
					Name: "policyAllowed",
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyAllowed"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
				{
					ID:   "PolicyDenied",
					Name: "policyDenied",
					Org:  "example",
					Path: "/path2/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path2/", "policyDenied"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
			},
		},
		"OkCaseAdminNoOrg": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "",
			pathPrefix: "/",
			expectedPolicies: []PolicyIdentity{
				{
					Org:  "example",
					Name: "policyAllowed",
				},
			},
			getPoliciesFilteredMethodResult: []Policy{
				{
					ID:   "PolicyAllowed",
					Name: "policyAllowed",
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyAllowed"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
			},
		},
		"OkCaseUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org: "example",
			expectedPolicies: []PolicyIdentity{
				{
					Org:  "example",
					Name: "policyAllowed",
				},
			},
			getPoliciesFilteredMethodResult: []Policy{
				{
					ID:   "PolicyAllowed",
					Name: "policyAllowed",
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyAllowed"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
				{
					ID:   "PolicyDenied",
					Name: "policyDenied",
					Org:  "example",
					Path: "/path2/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path2/", "policyDenied"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/1/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_LIST_POLICIES,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								POLICY_ACTION_LIST_POLICIES,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path2/"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseInvalidPath": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			pathPrefix: "/path*/",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: PathPrefix /path*/",
			},
		},
		"ErrorCaseInvalidOrg": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "!#$$%**^",
			pathPrefix: "/",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org !#$$%**^",
			},
		},
		"ErrorCaseInternalErrorGetPoliciesFiltered": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "",
			pathPrefix: "/path/",
			getPoliciesFilteredMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "123",
			pathPrefix: "/path/",
			getPoliciesFilteredMethodResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Org:  "example",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
			},
			getUserByExternalIDErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
		},
	}

	for x, testcase := range testcases {

		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetPoliciesFilteredMethod][0] = testcase.getPoliciesFilteredMethodResult
		testRepo.ArgsOut[GetPoliciesFilteredMethod][1] = testcase.getPoliciesFilteredMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		policies, err := testAPI.ListPolicies(testcase.requestInfo, testcase.org, testcase.pathPrefix)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.expectedPolicies, policies)
	}
}

func TestAuthAPI_UpdatePolicy(t *testing.T) {
	testcases := map[string]struct {
		requestInfo   RequestInfo
		org           string
		policyName    string
		path          string
		newPolicyName string
		newPath       string
		statements    []Statement
		newStatements []Statement

		getPolicyByNameMethodResult *Policy
		getGroupsByUserIDResult     []Group
		getAttachedPoliciesResult   []Policy
		getUserByExternalIDResult   *User
		updatePolicyMethodResult    *Policy

		wantError error

		getPolicyByNameMethodErr error
		getUserByExternalIDErr   error
		updatePolicyMethodErr    error

		getPolicyByNameMethodSpecialFunc func(string, string) (*Policy, error)
	}{
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
			updatePolicyMethodResult: &Policy{
				ID:   "test2",
				Name: "test2",
				Org:  "123",
				Path: "/path2/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path2/", "test"),
				Statements: &[]Statement{
					{
						Effect: "allow",
						Actions: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_USER, "/path2/"),
						},
					},
				},
			},
		},
		"ErrorCaseInvalidPolicyName": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "**!~#",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name **!~#",
			},
		},
		"ErrorCaseInvalidOrgName": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "**!~#",
			policyName: "p1",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org **!~#",
			},
		},
		"ErrorCaseInvalidNewPolicyName": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "**!~#",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: new name **!~#",
			},
		},
		"ErrorCaseInvalidNewPath": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/**~#!/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: new path /**~#!/",
			},
		},
		"ErrorCaseInvalidNewStatements": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "jblkasdjgp",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid effect: jblkasdjgp - Only 'allow' and 'deny' accepted",
			},
		},
		"ErrorCaseGetPolicyDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
		"ErrorCasePolicyNotFound": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodErr: &database.Error{
				Code: database.POLICY_NOT_FOUND,
			},
			wantError: &Error{
				Code: POLICY_BY_ORG_AND_NAME_NOT_FOUND,
			},
		},
		"ErrorCaseAuthUserNotFound": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
			updatePolicyMethodResult: &Policy{
				ID:   "test2",
				Name: "test2",
				Org:  "123",
				Path: "/path2/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path2/", "test"),
				Statements: &[]Statement{
					{
						Effect: "allow",
						Actions: []string{
							USER_ACTION_GET_USER,
						},
						Resources: []string{
							GetUrnPrefix("", RESOURCE_USER, "/path2/"),
						},
					},
				},
			},
			getUserByExternalIDErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
		},
		"ErrorCaseDenyResource": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_UPDATE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								POLICY_ACTION_UPDATE_POLICY,
							},
							Resources: []string{
								CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
							},
						},
					},
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam:123:policy/path/test",
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "123",
				Path: "/path/",
				Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam:123:policy/path/test",
			},
		},
		"ErrorCaseNewPolicyAlreadyExists": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodSpecialFunc: func(org string, name string) (*Policy, error) {
				if org == "123" && name == "test" {
					return &Policy{
						ID:   "test1",
						Name: "test",
						Org:  "123",
						Path: "/path/",
						Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
					}, nil
				} else {
					return &Policy{
						ID:   "test2",
						Name: "test2",
						Org:  "123",
						Path: "/path/",
						Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
					}, nil
				}
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_UPDATE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
			},
			wantError: &Error{
				Code:    POLICY_ALREADY_EXIST,
				Message: "Policy name: test2 already exists",
			},
		},
		"ErrorCaseNoPermissionsToRetrieveTarget": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodSpecialFunc: func(org string, name string) (*Policy, error) {
				if org == "123" && name == "test" {
					return &Policy{
						ID:   "test1",
						Name: "test",
						Org:  "123",
						Path: "/path/",
						Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
					}, nil
				} else {
					return &Policy{
						ID:   "test2",
						Name: "test2",
						Org:  "123",
						Path: "/path2/",
						Urn:  CreateUrn("123", RESOURCE_POLICY, "/path2/", "test"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									USER_ACTION_GET_USER,
								},
								Resources: []string{
									GetUrnPrefix("", RESOURCE_USER, "/path2/"),
								},
							},
						},
					}, nil
				}
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_UPDATE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam:123:policy/path2/test",
			},
		},
		"ErrorCaseNoPermissionsToUpdateTarget": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodSpecialFunc: func(org string, name string) (*Policy, error) {
				if org == "123" && name == "test" {
					return &Policy{
						ID:   "test1",
						Name: "test",
						Org:  "123",
						Path: "/path/",
						Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
					}, nil
				} else {
					return nil, &database.Error{
						Code: database.POLICY_NOT_FOUND,
					}
				}
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_UPDATE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam:123:policy/path2/test2",
			},
		},
		"ErrorCaseExplicitDenyPermissionsToUpdateTarget": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodSpecialFunc: func(org string, name string) (*Policy, error) {
				if org == "123" && name == "test" {
					return &Policy{
						ID:   "test1",
						Name: "test",
						Org:  "123",
						Path: "/path/",
						Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
					}, nil
				} else {
					return nil, &database.Error{
						Code: database.POLICY_NOT_FOUND,
					}
				}
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_UPDATE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path2/"),
							},
						},
					},
				},
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_UPDATE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path2/"),
							},
						},
					},
				},
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "deny",
							Actions: []string{
								POLICY_ACTION_UPDATE_POLICY,
							},
							Resources: []string{
								CreateUrn("123", RESOURCE_POLICY, "/path2/", "test2"),
							},
						},
					},
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:iam:123:policy/path2/test2",
			},
		},
		"ErrorCaseErrorUpdatingPolicy": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			org:        "123",
			policyName: "test",
			path:       "/path/",
			statements: []Statement{
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
			newPolicyName: "test2",
			newPath:       "/path2/",
			newStatements: []Statement{
				{
					Effect: "allow",
					Actions: []string{
						USER_ACTION_GET_USER,
					},
					Resources: []string{
						GetUrnPrefix("", RESOURCE_USER, "/path2/"),
					},
				},
			},
			getPolicyByNameMethodSpecialFunc: func(org string, name string) (*Policy, error) {
				if org == "123" && name == "test" {
					return &Policy{
						ID:   "test1",
						Name: "test",
						Org:  "123",
						Path: "/path/",
						Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "test"),
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
					}, nil
				} else {
					return nil, &database.Error{
						Code: database.POLICY_NOT_FOUND,
					}
				}
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_UPDATE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path2/"),
							},
						},
					},
				},
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Path: "/path/",
					Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_UPDATE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("123", RESOURCE_POLICY, "/path2/"),
							},
						},
					},
				},
			},
			updatePolicyMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
	}

	testRepo := makeTestRepo()
	testAPI := makeTestAPI(testRepo)

	for x, testcase := range testcases {
		testRepo.ArgsOut[UpdatePolicyMethod][0] = testcase.updatePolicyMethodResult
		testRepo.ArgsOut[UpdatePolicyMethod][1] = testcase.updatePolicyMethodErr
		testRepo.ArgsOut[GetPolicyByNameMethod][0] = testcase.getPolicyByNameMethodResult
		testRepo.ArgsOut[GetPolicyByNameMethod][1] = testcase.getPolicyByNameMethodErr
		testRepo.SpecialFuncs[GetPolicyByNameMethod] = testcase.getPolicyByNameMethodSpecialFunc
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		policy, err := testAPI.UpdatePolicy(testcase.requestInfo, testcase.org, testcase.policyName, testcase.newPolicyName, testcase.newPath, testcase.newStatements)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.updatePolicyMethodResult, policy)
	}
}

func TestAuthAPI_RemovePolicy(t *testing.T) {
	testcases := map[string]struct {
		requestInfo RequestInfo
		org         string
		name        string

		getPolicyByNameMethodResult *Policy
		getPolicyByNameMethodErr    error
		getGroupsByUserIDResult     []Group
		getAttachedPoliciesResult   []Policy
		getUserByExternalIDResult   *User
		getUserByExternalIDErr      error
		deletePolicyErr             error

		wantError error
	}{
		"OkCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:  "example",
			name: "test",
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
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
		},
		"ErrorCaseInvalidName": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:  "123",
			name: "invalid*",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name invalid*",
			},
		},
		"ErrorCaseInvalidOrg": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:  "**!^#$%",
			name: "invalid",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org **!^#$%",
			},
		},
		"ErrorCasePolicyNotExist": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:  "123",
			name: "policy",
			wantError: &Error{
				Code: POLICY_BY_ORG_AND_NAME_NOT_FOUND,
			},
			getPolicyByNameMethodErr: &database.Error{
				Code: database.POLICY_NOT_FOUND,
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:  "example",
			name: "test",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:example:policy/path/test",
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
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
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "123456",
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseNotEnoughPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:  "example",
			name: "test",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:example:policy/path/test",
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
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
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "123456",
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/"),
							},
						},
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_DELETE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/"),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								POLICY_ACTION_DELETE_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseRemoveFail": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:  "example",
			name: "test",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
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
			deletePolicyErr: &database.Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
	}

	for x, testcase := range testcases {

		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[RemovePolicyMethod][0] = testcase.deletePolicyErr
		testRepo.ArgsOut[GetPolicyByNameMethod][0] = testcase.getPolicyByNameMethodResult
		testRepo.ArgsOut[GetPolicyByNameMethod][1] = testcase.getPolicyByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		err := testAPI.RemovePolicy(testcase.requestInfo, testcase.org, testcase.name)
		checkMethodResponse(t, x, testcase.wantError, err, nil, nil)
	}
}

func TestAuthAPI_ListAttachedGroups(t *testing.T) {
	testcases := map[string]struct {
		requestInfo    RequestInfo
		org            string
		policyName     string
		expectedGroups []string

		getGroupsByUserIDResult   []Group
		getAttachedPoliciesResult []Policy
		getUserByExternalIDResult *User

		getAttachedGroupsResult []Group
		getAttachedGroupsErr    error

		getPolicyByNameMethodResult *Policy
		wantError                   error

		getPolicyByNameMethodErr error
	}{
		"OkCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "example",
			policyName: "test",
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
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
			getAttachedGroupsResult: []Group{
				{
					ID:   "Group1",
					Org:  "org1",
					Name: "group1",
				},
				{
					ID:   "Group2",
					Org:  "org2",
					Name: "group2",
				},
			},
			expectedGroups: []string{"group1", "group2"},
		},
		"ErrorCaseInvalidName": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "invalid*",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name invalid*",
			},
		},
		"ErrorCaseInvalidOrg": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "!*^**~$%",
			policyName: "p1",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org !*^**~$%",
			},
		},
		"ErrorCasePolicyNotExist": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "123",
			policyName: "policy",
			wantError: &Error{
				Code: POLICY_BY_ORG_AND_NAME_NOT_FOUND,
			},
			getPolicyByNameMethodErr: &database.Error{
				Code: database.POLICY_NOT_FOUND,
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "example",
			policyName: "test",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:example:policy/path/test",
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
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
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "123456",
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseNotEnoughPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:        "example",
			policyName: "test",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:example:policy/path/test",
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
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
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "123456",
			},
			getGroupsByUserIDResult: []Group{
				{
					ID:   "GROUP-USER-ID",
					Name: "groupUser",
				},
			},
			getAttachedPoliciesResult: []Policy{
				{
					ID:   "POLICY-USER-ID",
					Name: "policyUser",
					Statements: &[]Statement{
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_GET_POLICY,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/"),
							},
						},
						{
							Effect: "allow",
							Actions: []string{
								POLICY_ACTION_LIST_ATTACHED_GROUPS,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/"),
							},
						},
						{
							Effect: "deny",
							Actions: []string{
								POLICY_ACTION_LIST_ATTACHED_GROUPS,
							},
							Resources: []string{
								GetUrnPrefix("example", RESOURCE_POLICY, "/path/"),
							},
						},
					},
				},
			},
		},
		"ErrorCaseGetAttachedPoliciesFail": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:        "example",
			policyName: "test",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getPolicyByNameMethodResult: &Policy{
				ID:   "test1",
				Name: "test",
				Org:  "example",
				Path: "/path/",
				Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "test"),
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
			getAttachedGroupsErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
	}

	for x, testcase := range testcases {

		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetPolicyByNameMethod][0] = testcase.getPolicyByNameMethodResult
		testRepo.ArgsOut[GetPolicyByNameMethod][1] = testcase.getPolicyByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		testRepo.ArgsOut[GetAttachedGroupsMethod][0] = testcase.getAttachedGroupsResult
		testRepo.ArgsOut[GetAttachedGroupsMethod][1] = testcase.getAttachedGroupsErr
		groups, err := testAPI.ListAttachedGroups(testcase.requestInfo, testcase.org, testcase.policyName)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.expectedGroups, groups)
	}
}
