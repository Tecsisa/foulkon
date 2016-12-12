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
						Url:    "/url",
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
						Url:    "/url",
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
						Url:    "/url2",
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
						Url:    "/url",
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
						Url:    "/url2",
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
		}, "ErrorCaseNoPermissions": {
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
						Url:    "/url",
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
				Url:    "/url",
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
					Url:    "/url",
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
				Url:    "/url",
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
					Url:    "/url",
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
			name: "n1",
			org:  "*%~#@|",
			path: "/example/",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: org *%~#@|",
			},
		},
		"ErrorCaseInvalidPath": {
			name: "group1",
			org:  "org1",
			path: "/**%%/*123",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: path /**%%/*123",
			},
		},
		"ErrorCaseInvalidResource": {
			name: "group1",
			org:  "org1",
			path: "/path/",
			resource: ResourceEntity{
				Host:   "http://host.com",
				Url:    "invalid",
				Method: "GET",
				Urn:    "urn:ews:example:instance1:resource/get",
				Action: "action",
			},
			wantError: &Error{
				Code:    REGEX_NO_MATCH,
				Message: "No regex match in parameter: invalid",
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
				Url:    "/url",
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
				Url:    "/url",
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
				Url:    "/url",
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
				Url:    "/url",
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
				Url:    "/url",
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
				Url:    "/url",
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
