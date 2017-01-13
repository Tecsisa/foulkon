package api

import (
	"testing"

	"github.com/Tecsisa/foulkon/database"
	"github.com/stretchr/testify/assert"
)

func TestProxyAPI_GetProxyResources(t *testing.T) {
	testcases := map[string]struct {
		wantError error

		getProxyResourcesMethod []ProxyResource
		getProxyResourcesErr    error
	}{
		"OkCase": {
			getProxyResourcesMethod: []ProxyResource{
				{
					Resource: ResourceEntity{
						Host:   "host",
						Path:   "/path",
						Method: "Method",
						Urn:    "urn",
						Action: "action",
					},
				},
			},
		},
		"ErrorCaseInternalError": {
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getProxyResourcesErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
	}

	for n, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeProxyTestAPI(testRepo)

		testRepo.ArgsOut[GetProxyResourcesMethod][0] = testcase.getProxyResourcesMethod
		testRepo.ArgsOut[GetProxyResourcesMethod][2] = testcase.getProxyResourcesErr

		resources, err := testAPI.GetProxyResources()
		checkMethodResponse(t, n, testcase.wantError, err, testcase.getProxyResourcesMethod, resources)
	}
}

func TestWorkerAPI_GetProxyResourceByName(t *testing.T) {
	testcases := map[string]struct {
		// API Method args
		requestInfo       RequestInfo
		org               string
		name              string
		proxyResourceName string
		// Expected result
		expectedProxyResource *ProxyResource
		wantError             error
		// Manager Results
		getUserByExternalIDResult          *User
		getGroupsByUserIDResult            []TestUserGroupRelation
		getAttachedPoliciesResult          []TestPolicyGroupRelation
		getProxyResourceByNameMethodResult *ProxyResource
		// Manager Errors
		getUserByExternalIDMethodErr    error
		getProxyResourceByNameMethodErr error
	}{
		"OKCaseAdmin": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			name: "pr",
			org:  "org",
			expectedProxyResource: &ProxyResource{
				ID:   "12345",
				Name: "pr",
				Org:  "org",
				Path: "/path/",
				Resource: ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
				},
			},
			getProxyResourceByNameMethodResult: &ProxyResource{
				ID:   "12345",
				Name: "pr",
				Org:  "org",
				Path: "/path/",
				Resource: ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
				},
			},
		},
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			name: "pr",
			org:  "org",
			expectedProxyResource: &ProxyResource{
				ID:   "12345",
				Name: "pr",
				Org:  "org",
				Path: "/path/",
				Resource: ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									PROXY_ACTION_GET_PROXY_RESOURCE,
								},
								Resources: []string{
									GetUrnPrefix("org", RESOURCE_PROXY, "/path/"),
								},
							},
						},
					},
				},
			},
			getProxyResourceByNameMethodResult: &ProxyResource{
				ID:   "12345",
				Name: "pr",
				Org:  "org",
				Path: "/path/",
				Resource: ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
				},
			},
		},
		"ErrorCaseInvalidName": {
			name: "*%~#@|",
			org:  "org",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name *%~#@|",
			},
		},
		"ErrorCaseInvalidOrg": {
			name: "pr",
			org:  "*%~#@|",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org *%~#@|",
			},
		},
		"ErrorCaseProxyResourceNotFound": {
			name: "pr",
			org:  "org",
			wantError: &Error{
				Code: PROXY_RESOURCE_BY_ORG_AND_NAME_NOT_FOUND,
			},
			getProxyResourceByNameMethodErr: &database.Error{
				Code: database.PROXY_RESOURCE_NOT_FOUND,
			},
		},
		"ErrorCaseGetProxyResourceDBErr": {
			name: "pr",
			org:  "org",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getProxyResourceByNameMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
		"ErrorCaseNotAuthenticatedUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "pr",
			org:  "org",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
			getProxyResourceByNameMethodResult: &ProxyResource{
				ID:   "12345",
				Name: "pr",
				Org:  "org",
				Path: "/path/",
				Resource: ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
				},
			},
		},
		"ErrorCaseUnauthorizedResource": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "pr",
			org:  "org",
			expectedProxyResource: &ProxyResource{
				ID:   "12345",
				Name: "pr",
				Org:  "org",
				Path: "/path/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org:proxy/path/asd/pr",
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									PROXY_ACTION_GET_PROXY_RESOURCE,
								},
								Resources: []string{
									GetUrnPrefix("org", RESOURCE_PROXY, "/path/"),
								},
							},
							{
								Effect: "deny",
								Actions: []string{
									PROXY_ACTION_GET_PROXY_RESOURCE,
								},
								Resources: []string{
									GetUrnPrefix("org", RESOURCE_PROXY, "/path/asd"),
								},
							},
						},
					},
				},
			},
			getProxyResourceByNameMethodResult: &ProxyResource{
				ID:   "12345",
				Name: "pr",
				Org:  "org",
				Path: "/path/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/asd/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
				},
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "pr",
			org:  "org",
			expectedProxyResource: &ProxyResource{
				ID:   "12345",
				Name: "pr",
				Org:  "org",
				Path: "/path/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org:proxy/path/asd/pr",
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:         "POLICY-USER-ID",
						Name:       "policyUser",
						Path:       "/path/",
						Urn:        CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{},
					},
				},
			},
			getProxyResourceByNameMethodResult: &ProxyResource{
				ID:   "12345",
				Name: "pr",
				Org:  "org",
				Path: "/path/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/asd/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
				},
			},
		},
	}

	for n, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetProxyResourceByNameMethod][0] = testcase.getProxyResourceByNameMethodResult
		testRepo.ArgsOut[GetProxyResourceByNameMethod][1] = testcase.getProxyResourceByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult

		group, err := testAPI.GetProxyResourceByName(testcase.requestInfo, testcase.org, testcase.name)
		checkMethodResponse(t, n, testcase.wantError, err, testcase.expectedProxyResource, group)
	}
}

func TestWorkerAPI_RemoveProxyResource(t *testing.T) {
	testcases := map[string]struct {
		// API method args
		requestInfo RequestInfo
		name        string
		org         string
		// Expected Result
		wantError error
		// Manager Results
		getUserByExternalIDResult          *User
		getGroupsByUserIDResult            []TestUserGroupRelation
		getAttachedPoliciesResult          []TestPolicyGroupRelation
		getProxyResourceByNameMethodResult *ProxyResource
		// API Errors
		getUserByExternalIDMethodErr error
		getProxyResourceByMethodErr  error
		removeProxyResourceMethodErr error
		getGroupsByUserIDError       error
	}{
		"OKCaseAdmin": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			name: "pr",
			org:  "org",
			getProxyResourceByNameMethodResult: &ProxyResource{
				ID:   "12345",
				Name: "pr",
				Org:  "org",
				Path: "/path/",
				Resource: ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
			},
		},
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "pr",
			org:  "org",
			getProxyResourceByNameMethodResult: &ProxyResource{
				ID:   "12345",
				Name: "pr",
				Org:  "org",
				Path: "/path/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Org:  "org",
						Path: "/example/",
						Urn:  CreateUrn("org", RESOURCE_POLICY, "/example/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									PROXY_ACTION_DELETE_RESOURCE,
									PROXY_ACTION_GET_PROXY_RESOURCE,
								},
								Resources: []string{
									GetUrnPrefix("org", RESOURCE_PROXY, "/path/"),
								},
							},
						},
					},
				},
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/example/",
						Org:  "org",
						Urn:  CreateUrn("org", RESOURCE_GROUP, "/example/", "group"),
					},
				},
			},
		},
		"ErrorCaseInvalidName": {
			name: "invalid*",
			org:  "org",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name invalid*",
			},
		},
		"ErrorCaseInvalidOrg": {
			name: "pr",
			org:  "**^!$%&",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org **^!$%&",
			},
		},
		"ErrorCaseProxyResourceNotFound": {
			name: "pr",
			org:  "org",
			wantError: &Error{
				Code: PROXY_RESOURCE_BY_ORG_AND_NAME_NOT_FOUND,
			},
			getProxyResourceByMethodErr: &database.Error{
				Code: database.PROXY_RESOURCE_NOT_FOUND,
			},
		},
		"ErrorCaseUnauthorizedUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:  "org",
			name: "pr",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
			getProxyResourceByNameMethodResult: &ProxyResource{
				ID:   "12345",
				Name: "pr",
				Org:  "org",
				Path: "/path/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
				},
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrorCaseImplicitUnauthorizedDeleteProxyResource": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "pr",
			org:  "org",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org:proxy/path/pr",
			},
			getProxyResourceByNameMethodResult: &ProxyResource{
				ID:   "12345",
				Name: "pr",
				Org:  "org",
				Path: "/path/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("org", RESOURCE_USER, "/example/", "123456"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/example/",
						Org:  "org",
						Urn:  CreateUrn("org", RESOURCE_GROUP, "/example/", "group"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Org:  "org",
						Path: "/example/",
						Urn:  CreateUrn("org", RESOURCE_POLICY, "/example/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									PROXY_ACTION_GET_PROXY_RESOURCE,
								},
								Resources: []string{
									GetUrnPrefix("org", RESOURCE_PROXY, ""),
								},
							},
						},
					},
				},
			},
		},
		"ErrorCaseExplicitUnauthorizedDeleteProxyResource": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "pr",
			org:  "org",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org:proxy/path/pr",
			},
			getProxyResourceByNameMethodResult: &ProxyResource{
				ID:   "12345",
				Name: "pr",
				Org:  "org",
				Path: "/path/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("org", RESOURCE_USER, "/example/", "123456"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/example/",
						Org:  "org",
						Urn:  CreateUrn("org", RESOURCE_GROUP, "/example/", "group1"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Org:  "org",
						Path: "/example/",
						Urn:  CreateUrn("org", RESOURCE_POLICY, "/example/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									PROXY_ACTION_DELETE_RESOURCE,
									PROXY_ACTION_GET_PROXY_RESOURCE,
								},
								Resources: []string{
									GetUrnPrefix("org", RESOURCE_PROXY, ""),
								},
							},
							{
								Effect: "deny",
								Actions: []string{
									PROXY_ACTION_DELETE_RESOURCE,
								},
								Resources: []string{
									GetUrnPrefix("org", RESOURCE_PROXY, "/path/pr"),
								},
							},
						},
					},
				},
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "pr",
			org:  "org",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org:proxy/path/pr",
			},
			getProxyResourceByNameMethodResult: &ProxyResource{
				ID:   "12345",
				Name: "pr",
				Org:  "org",
				Path: "/path/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("org", RESOURCE_USER, "/example/", "123456"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/example/",
						Org:  "org",
						Urn:  CreateUrn("org", RESOURCE_GROUP, "/example/", "group1"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:         "POLICY-USER-ID",
						Name:       "policyUser",
						Org:        "org1",
						Path:       "/example/",
						Urn:        CreateUrn("org1", RESOURCE_POLICY, "/example/", "policyUser"),
						Statements: &[]Statement{},
					},
				},
			},
		},
		"ErrorCaseDeleteProxyResourceDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			name: "pr",
			org:  "org",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getProxyResourceByNameMethodResult: &ProxyResource{
				ID:   "12345",
				Name: "pr",
				Org:  "org",
				Path: "/path/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
			},
			removeProxyResourceMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
	}

	for n, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetProxyResourceByNameMethod][0] = testcase.getProxyResourceByNameMethodResult
		testRepo.ArgsOut[GetProxyResourceByNameMethod][1] = testcase.getProxyResourceByMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][1] = testcase.getGroupsByUserIDError
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		testRepo.ArgsOut[RemoveProxyResourceMethod][0] = testcase.removeProxyResourceMethodErr

		err := testAPI.RemoveProxyResource(testcase.requestInfo, testcase.org, testcase.name)
		checkMethodResponse(t, n, testcase.wantError, err, nil, nil)
	}
}

func TestWorkerAPI_UpdateProxyResource(t *testing.T) {
	testcases := map[string]struct {
		// API Method args
		requestInfo RequestInfo
		org         string
		name        string
		newName     string
		newPath     string
		newResource ResourceEntity
		// Expected Result
		expectedProxyResource *ProxyResource
		wantError             error
		// Manager Results
		getProxyResourceByNameResult            *ProxyResource
		getGroupsByUserIDResult                 []TestUserGroupRelation
		getAttachedPoliciesResult               []TestPolicyGroupRelation
		getUserByExternalIDResult               *User
		updateProxyResourceResult               *ProxyResource
		getProxyResourceByNameMethodSpecialFunc func(string, string) (*ProxyResource, error)
		// API Errors
		getProxyResourceByNameMethodErr error
		getUserByExternalIDMethodErr    error
		updateProxyResourceMethodErr    error
	}{
		"OKCaseAdmin": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:     "org",
			name:    "pr",
			newName: "newPr",
			newPath: "/new/",
			newResource: ResourceEntity{
				Host:   "http://new.com",
				Path:   "/new",
				Method: "POST",
				Action: "new:get",
				Urn:    "urn:ews:example:instance1:resource/get",
			},
			expectedProxyResource: &ProxyResource{
				ID:   "12345",
				Name: "newPr",
				Org:  "org",
				Path: "/new/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://new.com",
					Path:   "/new",
					Method: "POST",
					Action: "new:get",
					Urn:    "urn:ews:example:new:resource/get",
				},
			},
			getProxyResourceByNameResult: &ProxyResource{
				ID:   "12345",
				Name: "pr",
				Org:  "org",
				Path: "/path/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
					Urn:    "urn:ews:example:instance1:resource/get",
				},
			},
			updateProxyResourceResult: &ProxyResource{
				ID:   "12345",
				Name: "newPr",
				Org:  "org",
				Path: "/new/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://new.com",
					Path:   "/new",
					Method: "POST",
					Action: "new:get",
					Urn:    "urn:ews:example:new:resource/get",
				},
			},
		},
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:     "org",
			name:    "pr",
			newName: "newPr",
			newPath: "/new/",
			newResource: ResourceEntity{
				Host:   "http://new.com",
				Path:   "/new",
				Method: "POST",
				Action: "new:get",
				Urn:    "urn:ews:example:instance1:resource/get",
			},
			expectedProxyResource: &ProxyResource{
				ID:   "12345",
				Name: "newPr",
				Org:  "org",
				Path: "/new/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://new.com",
					Path:   "/new",
					Method: "POST",
					Action: "new:get",
					Urn:    "urn:ews:example:new:resource/get",
				},
			},
			getProxyResourceByNameResult: &ProxyResource{
				ID:   "12345",
				Name: "pr",
				Org:  "org",
				Path: "/path/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
					Urn:    "urn:ews:example:instance1:resource/get",
				},
			},
			updateProxyResourceResult: &ProxyResource{
				ID:   "12345",
				Name: "newPr",
				Org:  "org",
				Path: "/new/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://new.com",
					Path:   "/new",
					Method: "POST",
					Action: "new:get",
					Urn:    "urn:ews:example:new:resource/get",
				},
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Org:  "org",
						Path: "/path/",
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Org:  "org",
						Path: "/path/",
						Urn:  CreateUrn("org", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									PROXY_ACTION_GET_PROXY_RESOURCE,
									PROXY_ACTION_UPDATE_RESOURCE,
								},
								Resources: []string{
									GetUrnPrefix("org", RESOURCE_PROXY, ""),
								},
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
		},
		"ErrorCaseInvalidName": {
			org:     "org",
			newName: "%$%&&",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: new name %$%&&",
			},
		},
		"ErrorCaseInvalidPath": {
			org:     "org",
			newName: "pr",
			newPath: "/$",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: new path /$",
			},
		},
		"ErrorCaseInvalidResource": {
			org:     "org",
			newName: "pr",
			newPath: "/path/",
			newResource: ResourceEntity{
				Host: "invalid",
			},
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "Invalid parameter host, value: invalid",
			},
		},
		"ErrorCaseInvalidOrg": {
			org:     "$^**!",
			name:    "pr",
			newName: "newName",
			newPath: "/new/",
			newResource: ResourceEntity{
				Host:   "http://new.com",
				Path:   "/new",
				Method: "POST",
				Action: "new:get",
				Urn:    "urn:ews:example:new:resource/get",
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org $^**!",
			},
		},
		"ErrorCaseProxyResourceNotFound": {
			org:     "org",
			name:    "pr",
			newName: "newName",
			newPath: "/new/",
			newResource: ResourceEntity{
				Host:   "http://new.com",
				Path:   "/new",
				Method: "POST",
				Action: "new:get",
				Urn:    "urn:ews:example:new:resource/get",
			},
			wantError: &Error{
				Code: PROXY_RESOURCE_BY_ORG_AND_NAME_NOT_FOUND,
			},
			getProxyResourceByNameMethodErr: &database.Error{
				Code: database.PROXY_RESOURCE_NOT_FOUND,
			},
		},
		"ErrorCaseUnauthorizedUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:     "org",
			name:    "pr",
			newName: "newName",
			newPath: "/new/",
			newResource: ResourceEntity{
				Host:   "http://new.com",
				Path:   "/new",
				Method: "POST",
				Action: "new:get",
				Urn:    "urn:ews:example:new:resource/get",
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
			getProxyResourceByNameResult: &ProxyResource{
				ID:   "12345",
				Name: "pr",
				Org:  "org",
				Path: "/path/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
					Urn:    "urn:ews:example:instance1:resource/get",
				},
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrorCaseProxyResourceAlreadyExist": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:     "org",
			name:    "pr",
			newName: "newName",
			newPath: "/new/",
			newResource: ResourceEntity{
				Host:   "http://new.com",
				Path:   "/new",
				Method: "POST",
				Action: "new:get",
				Urn:    "urn:ews:example:new:resource/get",
			},
			wantError: &Error{
				Code:    PROXY_RESOURCE_ALREADY_EXIST,
				Message: "Proxy resource name: newName already exists",
			},
			getProxyResourceByNameMethodSpecialFunc: func(org string, name string) (*ProxyResource, error) {
				if org == "org" && name == "pr" {
					return &ProxyResource{
						ID:   "12345",
						Name: "pr",
						Org:  "org",
						Path: "/path/",
						Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
						Resource: ResourceEntity{
							Host:   "http://example.com",
							Path:   "/path",
							Method: "GET",
							Action: "example:get",
							Urn:    "urn:ews:example:instance1:resource/get",
						},
					}, nil
				}
				return &ProxyResource{
					ID:   "123456",
					Name: "pr2",
					Org:  "org2",
					Path: "/path2/",
					Urn:  CreateUrn("org2", RESOURCE_PROXY, "/path2/", "pr2"),
					Resource: ResourceEntity{
						Host:   "http://example2.com",
						Path:   "/path2",
						Method: "GET2",
						Action: "example:get2",
						Urn:    "urn:ews:example:instance2:resource/get",
					},
				}, nil
			},
		},
		"ErrorCaseGetProxyResourceDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:     "org",
			name:    "pr",
			newName: "newName",
			newPath: "/new/",
			newResource: ResourceEntity{
				Host:   "http://new.com",
				Path:   "/new",
				Method: "POST",
				Action: "new:get",
				Urn:    "urn:ews:example:new:resource/get",
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getProxyResourceByNameMethodSpecialFunc: func(org string, name string) (*ProxyResource, error) {
				if org == "org" && name == "pr" {
					return &ProxyResource{
						ID:   "12345",
						Name: "pr",
						Org:  "org",
						Path: "/path/",
						Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
						Resource: ResourceEntity{
							Host:   "http://example.com",
							Path:   "/path",
							Method: "GET",
							Action: "example:get",
							Urn:    "urn:ews:example:instance1:resource/get",
						},
					}, nil
				}

				return nil, &database.Error{
					Code: database.INTERNAL_ERROR,
				}
			},
		},
		"ErrorCaseUnauthorizedUpdateProxyResource": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:     "org",
			name:    "pr",
			newName: "newName",
			newPath: "/new/",
			newResource: ResourceEntity{
				Host:   "http://new.com",
				Path:   "/new",
				Method: "POST",
				Action: "new:get",
				Urn:    "urn:ews:example:new:resource/get",
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org:proxy/path/pr",
			},
			expectedProxyResource: &ProxyResource{
				ID:   "12345",
				Name: "newPr",
				Org:  "org",
				Path: "/new/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://new.com",
					Path:   "/new",
					Method: "POST",
					Action: "new:get",
					Urn:    "urn:ews:example:new:resource/get",
				},
			},
			getProxyResourceByNameResult: &ProxyResource{
				ID:   "12345",
				Name: "pr",
				Org:  "org",
				Path: "/path/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
					Urn:    "urn:ews:example:instance1:resource/get",
				},
			},
			updateProxyResourceResult: &ProxyResource{
				ID:   "12345",
				Name: "newPr",
				Org:  "org",
				Path: "/new/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://new.com",
					Path:   "/new",
					Method: "POST",
					Action: "new:get",
					Urn:    "urn:ews:example:new:resource/get",
				},
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Org:  "org",
						Path: "/path/",
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Org:  "org",
						Path: "/path/",
						Urn:  CreateUrn("org", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									PROXY_ACTION_GET_PROXY_RESOURCE,
								},
								Resources: []string{
									GetUrnPrefix("org", RESOURCE_PROXY, ""),
								},
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
		},
		"ErrorCaseDenyUpdateProxyGroup": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:     "org",
			name:    "pr",
			newName: "newName",
			newPath: "/new/",
			newResource: ResourceEntity{
				Host:   "http://new.com",
				Path:   "/new",
				Method: "POST",
				Action: "new:get",
				Urn:    "urn:ews:example:new:resource/get",
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org:proxy/path/pr",
			},
			expectedProxyResource: &ProxyResource{
				ID:   "12345",
				Name: "newPr",
				Org:  "org",
				Path: "/new/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://new.com",
					Path:   "/new",
					Method: "POST",
					Action: "new:get",
					Urn:    "urn:ews:example:new:resource/get",
				},
			},
			getProxyResourceByNameResult: &ProxyResource{
				ID:   "12345",
				Name: "pr",
				Org:  "org",
				Path: "/path/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
					Urn:    "urn:ews:example:instance1:resource/get",
				},
			},
			updateProxyResourceResult: &ProxyResource{
				ID:   "12345",
				Name: "newPr",
				Org:  "org",
				Path: "/new/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://new.com",
					Path:   "/new",
					Method: "POST",
					Action: "new:get",
					Urn:    "urn:ews:example:new:resource/get",
				},
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Org:  "org",
						Path: "/path/",
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Org:  "org",
						Path: "/path/",
						Urn:  CreateUrn("org", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									PROXY_ACTION_UPDATE_RESOURCE,
									PROXY_ACTION_GET_PROXY_RESOURCE,
								},
								Resources: []string{
									GetUrnPrefix("org", RESOURCE_PROXY, ""),
								},
							}, {
								Effect: "deny",
								Actions: []string{
									PROXY_ACTION_UPDATE_RESOURCE,
								},
								Resources: []string{
									GetUrnPrefix("org", RESOURCE_PROXY, "/path/"),
								},
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
		},
		"ErrorCaseNoPermissionsToUpdateTarget": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:     "org",
			name:    "pr",
			newName: "newName",
			newPath: "/new/",
			newResource: ResourceEntity{
				Host:   "http://new.com",
				Path:   "/new",
				Method: "POST",
				Action: "new:get",
				Urn:    "urn:ews:example:new:resource/get",
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org:proxy/new/newName",
			},
			getProxyResourceByNameMethodSpecialFunc: func(org string, name string) (*ProxyResource, error) {
				if org == "org" && name == "pr" {
					return &ProxyResource{
						ID:   "12345",
						Name: "pr",
						Org:  "org",
						Path: "/path/",
						Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
						Resource: ResourceEntity{
							Host:   "http://example.com",
							Path:   "/path",
							Method: "GET",
							Action: "example:get",
							Urn:    "urn:ews:example:instance1:resource/get",
						},
					}, nil
				}
				return nil, &database.Error{
					Code: database.PROXY_RESOURCE_NOT_FOUND,
				}
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Org:  "org",
						Path: "/path/",
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Org:  "org",
						Path: "/path/",
						Urn:  CreateUrn("org", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									PROXY_ACTION_UPDATE_RESOURCE,
									PROXY_ACTION_GET_PROXY_RESOURCE,
								},
								Resources: []string{
									GetUrnPrefix("org", RESOURCE_PROXY, "/path/"),
								},
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
		},
		"ErrorCaseDenyToUpdateTarget": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:     "org",
			name:    "pr",
			newName: "newName",
			newPath: "/new/",
			newResource: ResourceEntity{
				Host:   "http://new.com",
				Path:   "/new",
				Method: "POST",
				Action: "new:get",
				Urn:    "urn:ews:example:new:resource/get",
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org:proxy/new/newName",
			},
			getProxyResourceByNameMethodSpecialFunc: func(org string, name string) (*ProxyResource, error) {
				if org == "org" && name == "pr" {
					return &ProxyResource{
						ID:   "12345",
						Name: "pr",
						Org:  "org",
						Path: "/path/",
						Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
						Resource: ResourceEntity{
							Host:   "http://example.com",
							Path:   "/path",
							Method: "GET",
							Action: "example:get",
							Urn:    "urn:ews:example:instance1:resource/get",
						},
					}, nil
				}
				return nil, &database.Error{
					Code: database.PROXY_RESOURCE_NOT_FOUND,
				}
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Org:  "org",
						Path: "/path/",
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Org:  "org",
						Path: "/path/",
						Urn:  CreateUrn("org", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									PROXY_ACTION_UPDATE_RESOURCE,
									PROXY_ACTION_GET_PROXY_RESOURCE,
								},
								Resources: []string{
									GetUrnPrefix("org", RESOURCE_PROXY, ""),
								},
							}, {
								Effect: "deny",
								Actions: []string{
									PROXY_ACTION_UPDATE_RESOURCE,
								},
								Resources: []string{
									GetUrnPrefix("org", RESOURCE_PROXY, "/new/"),
								},
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
		},
		"ErrorCaseNoPermission": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			org:     "org",
			name:    "pr",
			newName: "newName",
			newPath: "/new/",
			newResource: ResourceEntity{
				Host:   "http://new.com",
				Path:   "/new",
				Method: "POST",
				Action: "new:get",
				Urn:    "urn:ews:example:new:resource/get",
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org:proxy/path/pr",
			},
			expectedProxyResource: &ProxyResource{
				ID:   "12345",
				Name: "newPr",
				Org:  "org",
				Path: "/new/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://new.com",
					Path:   "/new",
					Method: "POST",
					Action: "new:get",
					Urn:    "urn:ews:example:new:resource/get",
				},
			},
			getProxyResourceByNameResult: &ProxyResource{
				ID:   "12345",
				Name: "pr",
				Org:  "org",
				Path: "/path/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
					Urn:    "urn:ews:example:instance1:resource/get",
				},
			},
			updateProxyResourceResult: &ProxyResource{
				ID:   "12345",
				Name: "newPr",
				Org:  "org",
				Path: "/new/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://new.com",
					Path:   "/new",
					Method: "POST",
					Action: "new:get",
					Urn:    "urn:ews:example:new:resource/get",
				},
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Org:  "org",
						Path: "/path/",
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:         "POLICY-USER-ID",
						Name:       "policyUser",
						Org:        "org",
						Path:       "/path/",
						Urn:        CreateUrn("org", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{},
					},
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
		},
		"ErrorCaseUpdateGroupDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			org:     "org",
			name:    "pr",
			newName: "newName",
			newPath: "/new/",
			newResource: ResourceEntity{
				Host:   "http://new.com",
				Path:   "/new",
				Method: "POST",
				Action: "new:get",
				Urn:    "urn:ews:example:new:resource/get",
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getProxyResourceByNameResult: &ProxyResource{
				ID:   "12345",
				Name: "pr",
				Org:  "org",
				Path: "/path/",
				Urn:  CreateUrn("org", RESOURCE_PROXY, "/path/", "pr"),
				Resource: ResourceEntity{
					Host:   "http://example.com",
					Path:   "/path",
					Method: "GET",
					Action: "example:get",
					Urn:    "urn:ews:example:instance1:resource/get",
				},
			},
			updateProxyResourceMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
	}

	for n, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[UpdateProxyResourceMethod][0] = testcase.updateProxyResourceResult
		testRepo.ArgsOut[UpdateProxyResourceMethod][1] = testcase.updateProxyResourceMethodErr
		testRepo.ArgsOut[GetProxyResourceByNameMethod][0] = testcase.getProxyResourceByNameResult
		testRepo.ArgsOut[GetProxyResourceByNameMethod][1] = testcase.getProxyResourceByNameMethodErr
		testRepo.SpecialFuncs[GetProxyResourceByNameMethod] = testcase.getProxyResourceByNameMethodSpecialFunc
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult

		proxyResource, err := testAPI.UpdateProxyResource(testcase.requestInfo, testcase.org, testcase.name, testcase.newName, testcase.newPath, testcase.newResource)
		checkMethodResponse(t, n, testcase.wantError, err, testcase.expectedProxyResource, proxyResource)
	}
}

func TestWorkerAPI_ListProxyResources(t *testing.T) {
	testcases := map[string]struct {
		// API Method args
		requestInfo RequestInfo
		filter      *Filter
		// Expected result
		expectedProxyResources []ProxyResourceIdentity
		totalResult            int
		wantError              error
		// Manager Results
		getGroupsByUserIDResult   []TestUserGroupRelation
		getAttachedPoliciesResult []TestPolicyGroupRelation
		getUserByExternalIDResult *User
		getProxyResourcesMethod   []ProxyResource
		// Manager Errors
		getUserByExternalIDErr error
		getProxyResourcesErr   error
	}{
		"OkCaseAdmin": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			filter: &Filter{
				Org: "org",
			},
			expectedProxyResources: []ProxyResourceIdentity{
				{
					Name: "name",
					Org:  "org",
				},
				{
					Name: "name2",
					Org:  "org",
				},
			},
			totalResult: 1,
			getProxyResourcesMethod: []ProxyResource{
				{
					ID:   "ID",
					Name: "name",
					Path: "path",
					Org:  "org",
					Resource: ResourceEntity{
						Host:   "host",
						Path:   "/path",
						Method: "Method",
						Urn:    "urnr",
						Action: "action",
					},
					Urn: "urn",
				},
				{
					ID:   "ID2",
					Name: "name2",
					Path: "path2",
					Org:  "org",
					Resource: ResourceEntity{
						Host:   "host2",
						Path:   "/path2",
						Method: "Method2",
						Urn:    "urnr2",
						Action: "action2",
					},
					Urn: "urn2",
				},
			},
		},
		"OkCaseUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			filter: &testFilter,
			expectedProxyResources: []ProxyResourceIdentity{
				{
					Name: "name",
					Org:  "org",
				},
			},
			totalResult: 1,
			getProxyResourcesMethod: []ProxyResource{
				{
					ID:   "ID",
					Name: "name",
					Path: "/path/",
					Org:  "org",
					Resource: ResourceEntity{
						Host:   "host",
						Path:   "/path",
						Method: "Method",
						Urn:    "urnr",
						Action: "action",
					},
					Urn: CreateUrn("org", RESOURCE_PROXY, "/path/", "name"),
				},
				{
					ID:   "ID2",
					Name: "name2",
					Path: "/path2/",
					Org:  "org",
					Resource: ResourceEntity{
						Host:   "host2",
						Path:   "/path2",
						Method: "Method2",
						Urn:    "urnr2",
						Action: "action2",
					},
					Urn: CreateUrn("org", RESOURCE_PROXY, "/path2/", "name2"),
				},
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "1234"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/path/1/",
						Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Org:  "example",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									PROXY_ACTION_LIST_RESOURCES,
								},
								Resources: []string{
									GetUrnPrefix("org", RESOURCE_PROXY, "/path/"),
								},
							},
							{
								Effect: "deny",
								Actions: []string{
									PROXY_ACTION_LIST_RESOURCES,
								},
								Resources: []string{
									GetUrnPrefix("org", RESOURCE_PROXY, "/path2/"),
								},
							},
						},
					},
				},
			},
		},
		"ErrorCaseMaxLimitSize": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			filter: &Filter{
				Limit: 10000,
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: limit 10000, max limit allowed: 1000",
			},
		},
		"ErrorCaseInvalidPath": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			filter: &Filter{
				PathPrefix: "/path*/ /*",
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: pathPrefix /path*/ /*",
			},
		},
		"ErrorCaseInvalidOrg": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			filter: &Filter{
				Org: "!#$$%**^",
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org !#$$%**^",
			},
		},
		"ErrorCaseInternalErrorGetProxyResource": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			filter: &Filter{
				PathPrefix: "/path/",
			},
			getProxyResourcesErr: &database.Error{
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
			filter: &Filter{
				PathPrefix: "/path/",
			},
			getProxyResourcesMethod: []ProxyResource{
				{
					ID:   "ID",
					Name: "name",
					Path: "/path/",
					Org:  "org",
					Resource: ResourceEntity{
						Host:   "host",
						Path:   "/path",
						Method: "Method",
						Urn:    "urnr",
						Action: "action",
					},
					Urn: CreateUrn("org", RESOURCE_PROXY, "/path/", "name"),
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

	for n, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetProxyResourcesMethod][0] = testcase.getProxyResourcesMethod
		testRepo.ArgsOut[GetProxyResourcesMethod][1] = testcase.totalResult
		testRepo.ArgsOut[GetProxyResourcesMethod][2] = testcase.getProxyResourcesErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult

		resources, total, err := testAPI.ListProxyResources(testcase.requestInfo, testcase.filter)
		checkMethodResponse(t, n, testcase.wantError, err, testcase.expectedProxyResources, resources)
		assert.Equal(t, testcase.totalResult, total, "Error in test case %v", n)
	}
}

func TestWorkerAPI_AddProxyResource(t *testing.T) {
	testcases := map[string]struct {
		// API Method args
		requestInfo RequestInfo
		name        string
		org         string
		path        string
		resource    ResourceEntity
		// Expected results
		expectedProxyResource *ProxyResource
		wantError             error
		// Manager Results
		getUserByExternalIDResult *User
		getGroupsByUserIDResult   []TestUserGroupRelation
		getAttachedPoliciesResult []TestPolicyGroupRelation
		getProxyResource          *ProxyResource
		addMemberMethodResult     *Group
		// Manager Errors
		getProxyResourceMethodErr    error
		getUserByExternalIDMethodErr error
		addProxyResourceMethodErr    error
	}{
		"OKCaseAdmin": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			name: "name",
			org:  "org",
			path: "/example/",
			resource: ResourceEntity{
				Host:   "http://host.com",
				Path:   "/path",
				Method: "GET",
				Urn:    "urn:ews:example:instance1:resource/get",
				Action: "action",
			},
			expectedProxyResource: &ProxyResource{
				ID:   "ID",
				Name: "name",
				Path: "/example/",
				Org:  "org",
				Resource: ResourceEntity{
					Host:   "http://host.com",
					Path:   "/path",
					Method: "GET",
					Urn:    "urn:ews:example:instance1:resource/get",
					Action: "action",
				},
				Urn: "urn",
			},
			getProxyResourceMethodErr: &database.Error{
				Code: database.PROXY_RESOURCE_NOT_FOUND,
			},
		},
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "name",
			org:  "org",
			path: "/example/",
			resource: ResourceEntity{
				Host:   "http://host.com",
				Path:   "/path",
				Method: "GET",
				Urn:    "urn:ews:example:instance1:resource/get",
				Action: "action",
			},
			expectedProxyResource: &ProxyResource{
				ID:   "ID",
				Name: "name",
				Path: "/example/",
				Org:  "org",
				Resource: ResourceEntity{
					Host:   "http://host.com",
					Path:   "/path",
					Method: "GET",
					Urn:    "urn:ews:example:instance1:resource/get",
					Action: "action",
				},
				Urn: "urn",
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									PROXY_ACTION_CREATE_RESOURCE,
								},
								Resources: []string{
									GetUrnPrefix("org", RESOURCE_PROXY, "/example/"),
								},
							},
						},
					},
				},
			},
			getProxyResourceMethodErr: &database.Error{
				Code: database.PROXY_RESOURCE_NOT_FOUND,
			},
		},
		"ErrorCaseInvalidName": {
			name: "*%~#@|",
			org:  "org1",
			path: "/example/",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name *%~#@|",
			},
		},
		"ErrorCaseInvalidOrg": {
			name: "pr",
			org:  "*%~#@|",
			path: "/example/",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org *%~#@|",
			},
		},
		"ErrorCaseInvalidPath": {
			name: "pr",
			org:  "org1",
			path: "/**%%/*123",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: path /**%%/*123",
			},
		},
		"ErrorCaseInvalidResource": {
			name: "pr",
			org:  "org1",
			path: "/path/",
			resource: ResourceEntity{
				Host:   "http://host.com",
				Path:   "invalid",
				Method: "GET",
				Urn:    "urn:ews:example:instance1:resource/get",
				Action: "action",
			},
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "Invalid parameter path_resource, value: invalid",
			},
		},
		"ErrorProxyResourceAlreadyExists": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			name: "name1",
			org:  "org1",
			path: "/example/",
			resource: ResourceEntity{
				Host:   "http://host.com",
				Path:   "/path",
				Method: "GET",
				Urn:    "urn:ews:example:instance1:resource/get",
				Action: "action",
			},
			wantError: &Error{
				Code:    PROXY_RESOURCE_ALREADY_EXIST,
				Message: "Unable to create proxy resource, proxy resource with org org1 and name name1 already exist",
			},
			getUserByExternalIDResult: &User{
				ID:         "543210",
				ExternalID: "1234",
				Path:       "/path/",
			},
		},
		"ErrorCaseNotAuthenticatedUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "name1",
			org:  "org1",
			path: "/example/",
			resource: ResourceEntity{
				Host:   "http://host.com",
				Path:   "/path",
				Method: "GET",
				Urn:    "urn:ews:example:instance1:resource/get",
				Action: "action",
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrorCaseUnauthorizedResource": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "name1",
			org:  "org1",
			path: "/test/asd/",
			resource: ResourceEntity{
				Host:   "http://host.com",
				Path:   "/path",
				Method: "GET",
				Urn:    "urn:ews:example:instance1:resource/get",
				Action: "action",
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:proxy/test/asd/name1",
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									PROXY_ACTION_CREATE_RESOURCE,
								},
								Resources: []string{
									GetUrnPrefix("org1", RESOURCE_PROXY, "/test/"),
								},
							},
							{
								Effect: "deny",
								Actions: []string{
									PROXY_ACTION_CREATE_RESOURCE,
								},
								Resources: []string{
									GetUrnPrefix("org1", RESOURCE_PROXY, "/test/asd"),
								},
							},
						},
					},
				},
			},
			getProxyResourceMethodErr: &database.Error{
				Code: database.PROXY_RESOURCE_NOT_FOUND,
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "name1",
			org:  "org1",
			path: "/test/asd/",
			resource: ResourceEntity{
				Host:   "http://host.com",
				Path:   "/path",
				Method: "GET",
				Urn:    "urn:ews:example:instance1:resource/get",
				Action: "action",
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:iam:org1:proxy/test/asd/name1",
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_GROUP, "/path/", "groupUser"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:         "POLICY-USER-ID",
						Name:       "policyUser",
						Path:       "/path/",
						Urn:        CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{},
					},
				},
			},
			getProxyResourceMethodErr: &database.Error{
				Code: database.PROXY_RESOURCE_NOT_FOUND,
			},
		},
		"ErrorCaseAddProxyResourceDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			name: "name1",
			org:  "org1",
			path: "/example/",
			resource: ResourceEntity{
				Host:   "http://host.com",
				Path:   "/path",
				Method: "GET",
				Urn:    "urn:ews:example:instance1:resource/get",
				Action: "action",
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getProxyResourceMethodErr: &database.Error{
				Code: database.PROXY_RESOURCE_NOT_FOUND,
			},
			addProxyResourceMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
		"ErrorCaseGetProxyResourceDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			name: "name1",
			org:  "org1",
			path: "/example/",
			resource: ResourceEntity{
				Host:   "http://host.com",
				Path:   "/path",
				Method: "GET",
				Urn:    "urn:ews:example:instance1:resource/get",
				Action: "action",
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getProxyResourceMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
	}

	for x, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetProxyResourceByNameMethod][0] = testcase.getProxyResource
		testRepo.ArgsOut[GetProxyResourceByNameMethod][1] = testcase.getProxyResourceMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		testRepo.ArgsOut[AddProxyResourceMethod][0] = testcase.expectedProxyResource
		testRepo.ArgsOut[AddProxyResourceMethod][1] = testcase.addProxyResourceMethodErr

		proxyResource, err := testAPI.AddProxyResource(testcase.requestInfo, testcase.name, testcase.org, testcase.path, testcase.resource)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.expectedProxyResource, proxyResource)
	}
}
