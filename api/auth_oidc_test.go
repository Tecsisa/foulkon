package api

import (
	"testing"

	"github.com/Tecsisa/foulkon/database"
	"github.com/stretchr/testify/assert"
)

func TestWorkerAPI_AddOidcProvider(t *testing.T) {
	testcases := map[string]struct {
		requestInfo      RequestInfo
		oidcProviderName string
		path             string
		issuerURL        string
		oidcClients      []string

		getGroupsByUserIDResult   []TestUserGroupRelation
		getAttachedPoliciesResult []TestPolicyGroupRelation
		getUserByExternalIDResult *User

		addOidcProviderMethodResult       *OidcProvider
		getOidcProviderByNameMethodResult *OidcProvider
		wantError                         error

		getOidcProviderByNameMethodErr error
		addOidcProviderMethodErr       error
	}{
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			oidcProviderName: "test",
			path:             "/path/",
			issuerURL:        "https://test.com",
			oidcClients: []string{
				"client",
			},
			getOidcProviderByNameMethodErr: &database.Error{
				Code: database.AUTH_OIDC_PROVIDER_NOT_FOUND,
			},
			addOidcProviderMethodResult: &OidcProvider{
				ID:        "test1",
				Name:      "test",
				Path:      "/path/",
				Urn:       CreateUrn("123", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "test"),
				IssuerURL: "https://test.com",
				OidcClients: []OidcClient{
					{
						Name: "client",
					},
				},
			},
		},
		"ErrorCaseOidcProviderAlreadyExists": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			oidcProviderName: "test",
			path:             "/path/",
			issuerURL:        "https://test.com",
			oidcClients: []string{
				"client",
			},
			getOidcProviderByNameMethodResult: &OidcProvider{
				ID:        "test1",
				Name:      "test",
				Path:      "/path/",
				Urn:       CreateUrn("123", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "test"),
				IssuerURL: "https://test.com",
				OidcClients: []OidcClient{
					{
						Name: "client",
					},
				},
			},
			wantError: &Error{
				Code:    AUTH_OIDC_PROVIDER_ALREADY_EXIST,
				Message: "Unable to create OIDC provider, OIDC provider with name test already exist",
			},
		},
		"ErrorCaseBadName": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			oidcProviderName: "**!^#~",
			path:             "/path/",
			issuerURL:        "https://test.com",
			oidcClients: []string{
				"client",
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name **!^#~",
			},
		},
		"ErrorCaseInvalidClientNames": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			oidcProviderName: "test",
			path:             "/path/",
			issuerURL:        "https://test.com",
			oidcClients:      []string{"~$"},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: client name ~$",
			},
		},
		"ErrorCaseBadPath": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			oidcProviderName: "test",
			path:             "*/ /**!^#~path/",
			issuerURL:        "https://test.com",
			oidcClients: []string{
				"client1",
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: path */ /**!^#~path/",
			},
		},
		"ErrorCaseInvalidIssuerUrl": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			oidcProviderName: "test",
			path:             "/path/",
			issuerURL:        "~htt:/pjs://test.com",
			oidcClients:      []string{},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: issuerUrl ~htt:/pjs://test.com",
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			oidcProviderName: "test",
			path:             "/path/",
			issuerURL:        "https://test.com",
			oidcClients: []string{
				"client",
			},
			getOidcProviderByNameMethodErr: &database.Error{
				Code: database.AUTH_OIDC_PROVIDER_NOT_FOUND,
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
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:auth::oidc/path/test",
			},
		},
		"ErrorCaseDenyResource": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			oidcProviderName: "test",
			path:             "/path/",
			issuerURL:        "https://test.com",
			oidcClients: []string{
				"client",
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
						Name: "policy",
						Org:  "example",
						Path: "/path/",
						Urn:  CreateUrn("example", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									AUTH_OIDC_ACTION_GET_PROVIDER,
								},
								Resources: []string{
									GetUrnPrefix("", RESOURCE_AUTH_OIDC_PROVIDER, "/"),
								},
							},
							{
								Effect: "allow",
								Actions: []string{
									AUTH_OIDC_ACTION_CREATE_PROVIDER,
								},
								Resources: []string{
									GetUrnPrefix("", RESOURCE_AUTH_OIDC_PROVIDER, "/"),
								},
							},
							{
								Effect: "deny",
								Actions: []string{
									AUTH_OIDC_ACTION_CREATE_PROVIDER,
								},
								Resources: []string{
									CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "test"),
								},
							},
						},
					},
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:auth::oidc/path/test",
			},
		},
		"ErrorCaseAddOidcProviderErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			oidcProviderName: "test",
			path:             "/path/",
			issuerURL:        "https://test.com",
			oidcClients: []string{
				"client",
			},
			getOidcProviderByNameMethodErr: &database.Error{
				Code: database.AUTH_OIDC_PROVIDER_NOT_FOUND,
			},
			addOidcProviderMethodErr: &database.Error{
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
			oidcProviderName: "test",
			path:             "/path/",
			issuerURL:        "https://test.com",
			oidcClients: []string{
				"client",
			},
			getOidcProviderByNameMethodErr: &database.Error{
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
		testRepo.ArgsOut[AddOidcProviderMethod][0] = testcase.addOidcProviderMethodResult
		testRepo.ArgsOut[AddOidcProviderMethod][1] = testcase.addOidcProviderMethodErr
		testRepo.ArgsOut[GetOidcProviderByNameMethod][0] = testcase.getOidcProviderByNameMethodResult
		testRepo.ArgsOut[GetOidcProviderByNameMethod][1] = testcase.getOidcProviderByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		oidcProvider, err := testAPI.AddOidcProvider(testcase.requestInfo, testcase.oidcProviderName,
			testcase.path, testcase.issuerURL, testcase.oidcClients)
		checkMethodResponse(t, x, testcase.wantError, err, oidcProvider, testcase.addOidcProviderMethodResult)
	}
}

func TestWorkerAPI_GetOidcProviderByName(t *testing.T) {
	testcases := map[string]struct {
		requestInfo      RequestInfo
		oidcProviderName string

		getGroupsByUserIDResult   []TestUserGroupRelation
		getAttachedPoliciesResult []TestPolicyGroupRelation
		getUserByExternalIDResult *User

		getOidcProviderByNameMethodResult *OidcProvider
		wantError                         error

		getOidcProviderByNameMethodErr error
	}{
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			oidcProviderName: "test",
			getOidcProviderByNameMethodResult: &OidcProvider{
				ID:        "test1",
				Name:      "test",
				Path:      "/path/",
				Urn:       CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "test"),
				IssuerURL: "https://test.com",
				OidcClients: []OidcClient{
					{
						Name: "client",
					},
				},
			},
		},
		"ErrorCaseInternalError": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			oidcProviderName: "test",
			getOidcProviderByNameMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
		},
		"ErrorCaseBadOidcProviderName": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			oidcProviderName: "~#**!",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name ~#**!",
			},
		},
		"ErrorCaseOidcProviderNotFound": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			oidcProviderName: "test",
			getOidcProviderByNameMethodErr: &database.Error{
				Code: database.AUTH_OIDC_PROVIDER_NOT_FOUND,
			},
			wantError: &Error{
				Code: AUTH_OIDC_PROVIDER_BY_NAME_NOT_FOUND,
			},
		},
		"ErrorCaseNoPermissions": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			oidcProviderName: "test",
			getOidcProviderByNameMethodResult: &OidcProvider{
				ID:        "test1",
				Name:      "test",
				Path:      "/path/",
				Urn:       CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "test"),
				IssuerURL: "https://test.com",
				OidcClients: []OidcClient{
					{
						Name: "client",
					},
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
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:auth::oidc/path/test",
			},
		},
		"ErrorCaseDenyResourceErr": {
			requestInfo: RequestInfo{
				Identifier: "1234",
				Admin:      false,
			},
			oidcProviderName: "test",
			getOidcProviderByNameMethodResult: &OidcProvider{
				ID:        "test1",
				Name:      "test",
				Path:      "/path/",
				Urn:       CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "test"),
				IssuerURL: "https://test.com",
				OidcClients: []OidcClient{
					{
						Name: "client",
					},
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
									AUTH_OIDC_ACTION_GET_PROVIDER,
								},
								Resources: []string{
									GetUrnPrefix("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/"),
								},
							},
							{
								Effect: "deny",
								Actions: []string{
									AUTH_OIDC_ACTION_GET_PROVIDER,
								},
								Resources: []string{
									CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "test"),
								},
							},
						},
					},
				},
			},
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 1234 is not allowed to access to resource urn:iws:auth::oidc/path/test",
			},
		},
	}

	for x, testcase := range testcases {

		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetOidcProviderByNameMethod][0] = testcase.getOidcProviderByNameMethodResult
		testRepo.ArgsOut[GetOidcProviderByNameMethod][1] = testcase.getOidcProviderByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		oidcProvider, err := testAPI.GetOidcProviderByName(testcase.requestInfo, testcase.oidcProviderName)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.getOidcProviderByNameMethodResult, oidcProvider)
	}
}

func TestWorkerAPI_ListOidcProviders(t *testing.T) {
	testcases := map[string]struct {
		// API Method args
		requestInfo RequestInfo
		filter      *Filter
		// Expected result
		expectedOidcProviders []string
		totalResult           int
		wantError             error
		// Manager Results
		getGroupsByUserIDResult              []TestUserGroupRelation
		getAttachedPoliciesResult            []TestPolicyGroupRelation
		getUserByExternalIDResult            *User
		getOidcProvidersFilteredMethodResult []OidcProvider
		// Manager Errors
		getUserByExternalIDErr            error
		getOidcProvidersFilteredMethodErr error
	}{
		"OkCaseAdmin": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			filter: &Filter{
				Org: "example",
			},
			expectedOidcProviders: []string{
				"oidcProviderAllowed",
				"oidcProviderDenied",
			},
			totalResult: 2,
			getOidcProvidersFilteredMethodResult: []OidcProvider{
				{
					ID:        "oidcProviderAllowed",
					Name:      "oidcProviderAllowed",
					Path:      "/path/",
					Urn:       CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "oidcProviderAllowed"),
					IssuerURL: "https://oidcProviderAllowed.com",
					OidcClients: []OidcClient{
						{
							Name: "client",
						},
					},
				},
				{
					ID:        "oidcProviderDenied",
					Name:      "oidcProviderDenied",
					Path:      "/path/",
					Urn:       CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "oidcProviderDenied"),
					IssuerURL: "https://oidcProviderDenied.com",
					OidcClients: []OidcClient{
						{
							Name: "client",
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
			filter: &Filter{
				PathPrefix: "/path/",
			},
			expectedOidcProviders: []string{
				"oidcProviderAllowed",
			},
			totalResult: 1,
			getOidcProvidersFilteredMethodResult: []OidcProvider{
				{
					ID:        "oidcProviderAllowed",
					Name:      "oidcProviderAllowed",
					Path:      "/path/",
					Urn:       CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "oidcProviderAllowed"),
					IssuerURL: "https://oidcProviderAllowed.com",
					OidcClients: []OidcClient{
						{
							Name: "client",
						},
					},
				},
				{
					ID:        "oidcProviderDenied",
					Name:      "oidcProviderDenied",
					Path:      "/path/",
					Urn:       CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "oidcProviderDenied"),
					IssuerURL: "https://oidcProviderDenied.com",
					OidcClients: []OidcClient{
						{
							Name: "client",
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
									AUTH_OIDC_ACTION_LIST_PROVIDERS,
								},
								Resources: []string{
									GetUrnPrefix("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/"),
								},
							},
							{
								Effect: "deny",
								Actions: []string{
									AUTH_OIDC_ACTION_LIST_PROVIDERS,
								},
								Resources: []string{
									CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "oidcProviderDenied"),
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
				PathPrefix: "/path/",
				Limit:      10000,
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
				Org:        "123",
				Limit:      0,
			},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: pathPrefix /path*/ /*",
			},
		},
		"ErrorCaseInternalErrorOidcProvidersFiltered": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			filter: &Filter{
				PathPrefix: "/path/",
			},
			getOidcProvidersFilteredMethodErr: &database.Error{
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
				Org:        "123",
			},
			getOidcProvidersFilteredMethodResult: []OidcProvider{
				{
					ID:        "oidcProviderDenied",
					Name:      "oidcProviderDenied",
					Path:      "/path/",
					Urn:       CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "oidcProviderDenied"),
					IssuerURL: "https://oidcProviderDenied.com",
					OidcClients: []OidcClient{
						{
							Name: "client",
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

		testRepo.ArgsOut[GetOidcProvidersFilteredMethod][0] = testcase.getOidcProvidersFilteredMethodResult
		testRepo.ArgsOut[GetOidcProvidersFilteredMethod][1] = testcase.totalResult
		testRepo.ArgsOut[GetOidcProvidersFilteredMethod][2] = testcase.getOidcProvidersFilteredMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		oidcProviders, total, err := testAPI.ListOidcProviders(testcase.requestInfo, testcase.filter)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.expectedOidcProviders, oidcProviders)
		assert.Equal(t, testcase.totalResult, total, "Error in test case %v", x)
	}
}

func TestWorkerAPI_UpdateOidcProvider(t *testing.T) {
	testcases := map[string]struct {
		// API Method args
		requestInfo         RequestInfo
		oidcProviderName    string
		newOidcProviderName string
		newPath             string
		newIssuerUrl        string
		newClients          []string
		// Expected result
		expectedOidcProvider *OidcProvider
		wantError            error
		// Manager Results
		getOidcProviderByNameResult            *OidcProvider
		getGroupMembersResult                  []User
		getGroupsByUserIDResult                []TestUserGroupRelation
		getAttachedPoliciesResult              []TestPolicyGroupRelation
		getUserByExternalIDResult              *User
		updateOidcProviderResult               *OidcProvider
		getOidcProviderByNameMethodSpecialFunc func(string) (*OidcProvider, error)
		// API Errors
		getOidcProviderByNameErr     error
		getUserByExternalIDMethodErr error
		updateOidcProviderMethodErr  error
	}{
		"OKCaseAdmin": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			oidcProviderName:    "oidcProvider1",
			newOidcProviderName: "oidcProviderNewName",
			newPath:             "/new/",
			newIssuerUrl:        "http://oidcProvider1.com",
			newClients:          []string{"newClient1", "newClient2"},
			expectedOidcProvider: &OidcProvider{
				ID:        "12345",
				Name:      "newName",
				Path:      "/new/",
				Urn:       CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/new/", "oidcProviderNewName"),
				IssuerURL: "http://oidcProvider1.com",
				OidcClients: []OidcClient{
					{
						Name: "newClient1",
					},
					{
						Name: "newClient2",
					},
				},
			},
			getOidcProviderByNameResult: &OidcProvider{
				ID:   "12345",
				Name: "oidcProvider1",
				Path: "/path/",
				Urn:  CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "oidcProvider1"),
			},
			updateOidcProviderResult: &OidcProvider{
				ID:        "12345",
				Name:      "newName",
				Path:      "/new/",
				Urn:       CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/new/", "oidcProviderNewName"),
				IssuerURL: "http://oidcProvider1.com",
				OidcClients: []OidcClient{
					{
						Name: "newClient1",
					},
					{
						Name: "newClient2",
					},
				},
			},
		},
		"OKCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			oidcProviderName:    "oidcProvider1",
			newOidcProviderName: "oidcProviderNewName",
			newPath:             "/new/",
			newIssuerUrl:        "http://oidcProvider1.com",
			newClients:          []string{"newClient1", "newClient2"},
			expectedOidcProvider: &OidcProvider{
				ID:        "12345",
				Name:      "newName",
				Path:      "/new/",
				Urn:       CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/new/", "oidcProviderNewName"),
				IssuerURL: "http://oidcProvider1.com",
				OidcClients: []OidcClient{
					{
						Name: "newClient1",
					},
					{
						Name: "newClient2",
					},
				},
			},
			getOidcProviderByNameResult: &OidcProvider{
				ID:   "12345",
				Name: "oidcProvider1",
				Path: "/path/",
				Urn:  CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "oidcProvider1"),
			},
			updateOidcProviderResult: &OidcProvider{
				ID:        "12345",
				Name:      "newName",
				Path:      "/new/",
				Urn:       CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/new/", "oidcProviderNewName"),
				IssuerURL: "http://oidcProvider1.com",
				OidcClients: []OidcClient{
					{
						Name: "newClient1",
					},
					{
						Name: "newClient2",
					},
				},
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Org:  "org1",
						Path: "/path/1/",
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Org:  "org1",
						Path: "/path/",
						Urn:  CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									AUTH_OIDC_ACTION_GET_PROVIDER,
									AUTH_OIDC_ACTION_UPDATE_PROVIDER,
								},
								Resources: []string{
									GetUrnPrefix("", RESOURCE_AUTH_OIDC_PROVIDER, ""),
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
			newOidcProviderName: "%$%&&",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: new name %$%&&",
			},
		},
		"ErrorCaseInvalidPath": {
			newOidcProviderName: "oidcProvider1",
			newPath:             "/$",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: new path /$",
			},
		},
		"ErrorCaseInvalidOidcClientNames": {
			oidcProviderName:    "oidcProvider1",
			newOidcProviderName: "newName",
			newPath:             "/new/",
			newIssuerUrl:        "http://test.com",
			newClients:          []string{"~@"},
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: client name ~@",
			},
		},
		"ErrorCaseInvalidIssuerUrl": {
			oidcProviderName:    "oidcProvider1",
			newOidcProviderName: "newName",
			newPath:             "/new/",
			newIssuerUrl:        "$~",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: issuerUrl $~",
			},
		},
		"ErrorCaseOidcProviderNotFound": {
			oidcProviderName:    "oidcProvider1",
			newOidcProviderName: "newName",
			newPath:             "/new/",
			newIssuerUrl:        "http://test.com",
			wantError: &Error{
				Code: AUTH_OIDC_PROVIDER_BY_NAME_NOT_FOUND,
			},
			getOidcProviderByNameErr: &database.Error{
				Code: database.AUTH_OIDC_PROVIDER_NOT_FOUND,
			},
		},
		"ErrorCaseUnauthorizedUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			oidcProviderName:    "oidcProvider1",
			newOidcProviderName: "newName",
			newPath:             "/new/",
			newIssuerUrl:        "http://test.com",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
			getOidcProviderByNameResult: &OidcProvider{
				ID:   "12345",
				Name: "oidcProvider1",
				Path: "/path/",
				Urn:  CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "oidcProvider1"),
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrorCaseOidcProviderAlreadyExist": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			oidcProviderName:    "oidcProvider1",
			newOidcProviderName: "newName",
			newPath:             "/new/",
			newIssuerUrl:        "http://test.com",
			wantError: &Error{
				Code:    AUTH_OIDC_PROVIDER_ALREADY_EXIST,
				Message: "OIDC Provider name: newName already exists",
			},
			getOidcProviderByNameMethodSpecialFunc: func(name string) (*OidcProvider, error) {
				if name == "oidcProvider1" {
					return &OidcProvider{
						ID:   "12345",
						Name: "oidcProvider1",
						Path: "/path/",
						Urn:  CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "oidcProvider1"),
					}, nil
				}
				return &OidcProvider{
					ID:   "anotherId",
					Name: name,
					Path: "/path/",
					Urn:  CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", name),
				}, nil
			},
		},
		"ErrorCaseGetOidcProviderDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			oidcProviderName:    "oidcProvider1",
			newOidcProviderName: "newName",
			newPath:             "/new/",
			newIssuerUrl:        "http://test.com",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getOidcProviderByNameMethodSpecialFunc: func(name string) (*OidcProvider, error) {
				if name == "oidcProvider1" {
					return &OidcProvider{
						ID:   "12345",
						Name: "oidcProvider1",
						Path: "/path/",
						Urn:  CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "oidcProvider1"),
					}, nil
				}

				return nil, &database.Error{
					Code: database.INTERNAL_ERROR,
				}
			},
		},
		"ErrorCaseUnauthorizedUpdateOidcProvider": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			oidcProviderName:    "oidcProvider1",
			newOidcProviderName: "newName",
			newPath:             "/new/",
			newIssuerUrl:        "http://test.com",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:auth::oidc/path/oidcProvider1",
			},
			getOidcProviderByNameResult: &OidcProvider{
				ID:   "12345",
				Name: "oidcProvider1",
				Path: "/path/",
				Urn:  CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "oidcProvider1"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Org:  "org1",
						Path: "/path/1/",
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Org:  "org1",
						Path: "/path/",
						Urn:  CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									AUTH_OIDC_ACTION_GET_PROVIDER,
								},
								Resources: []string{
									GetUrnPrefix("", RESOURCE_AUTH_OIDC_PROVIDER, ""),
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
		"ErrorCaseDenyUpdateGroup": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			oidcProviderName:    "oidcProvider1",
			newOidcProviderName: "newName",
			newPath:             "/new/",
			newIssuerUrl:        "http://test.com",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:auth::oidc/path/oidcProvider1",
			},
			getOidcProviderByNameResult: &OidcProvider{
				ID:   "12345",
				Name: "oidcProvider1",
				Path: "/path/",
				Urn:  CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "oidcProvider1"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Org:  "org1",
						Path: "/path/1/",
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Org:  "org1",
						Path: "/path/",
						Urn:  CreateUrn("org1", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									AUTH_OIDC_ACTION_GET_PROVIDER,
									AUTH_OIDC_ACTION_UPDATE_PROVIDER,
								},
								Resources: []string{
									GetUrnPrefix("", RESOURCE_AUTH_OIDC_PROVIDER, ""),
								},
							},
							{
								Effect: "deny",
								Actions: []string{
									AUTH_OIDC_ACTION_UPDATE_PROVIDER,
								},
								Resources: []string{
									GetUrnPrefix("", RESOURCE_AUTH_OIDC_PROVIDER, "/path"),
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
			oidcProviderName:    "oidcProvider1",
			newOidcProviderName: "newName",
			newPath:             "/new/",
			newIssuerUrl:        "http://test.com",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:auth::oidc/new/newName",
			},
			getOidcProviderByNameMethodSpecialFunc: func(name string) (*OidcProvider, error) {
				if name == "oidcProvider1" {
					return &OidcProvider{
						ID:   "12345",
						Name: "oidcProvider1",
						Path: "/path/",
						Urn:  CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "oidcProvider1"),
					}, nil
				}

				return nil, &database.Error{
					Code: database.AUTH_OIDC_PROVIDER_NOT_FOUND,
				}
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "group1",
						Org:  "123",
						Path: "/new/",
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Org:  "123",
						Path: "/path/",
						Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									AUTH_OIDC_ACTION_GET_PROVIDER,
									AUTH_OIDC_ACTION_UPDATE_PROVIDER,
								},
								Resources: []string{
									GetUrnPrefix("", RESOURCE_AUTH_OIDC_PROVIDER, "/path"),
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
			oidcProviderName:    "oidcProvider1",
			newOidcProviderName: "newName",
			newPath:             "/new/",
			newIssuerUrl:        "http://test.com",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:auth::oidc/new/newName",
			},
			getOidcProviderByNameMethodSpecialFunc: func(name string) (*OidcProvider, error) {
				if name == "oidcProvider1" {
					return &OidcProvider{
						ID:   "12345",
						Name: "oidcProvider1",
						Path: "/path/",
						Urn:  CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "oidcProvider1"),
					}, nil
				}

				return nil, &database.Error{
					Code: database.AUTH_OIDC_PROVIDER_NOT_FOUND,
				}
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "group1",
						Org:  "123",
						Path: "/new/",
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Org:  "123",
						Path: "/path/",
						Urn:  CreateUrn("123", RESOURCE_POLICY, "/path/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									AUTH_OIDC_ACTION_GET_PROVIDER,
									AUTH_OIDC_ACTION_UPDATE_PROVIDER,
								},
								Resources: []string{
									GetUrnPrefix("", RESOURCE_AUTH_OIDC_PROVIDER, "/"),
								},
							},
							{
								Effect: "deny",
								Actions: []string{
									AUTH_OIDC_ACTION_UPDATE_PROVIDER,
								},
								Resources: []string{
									GetUrnPrefix("", RESOURCE_AUTH_OIDC_PROVIDER, "/new"),
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
			oidcProviderName:    "oidcProvider1",
			newOidcProviderName: "newName",
			newPath:             "/new/",
			newIssuerUrl:        "http://test.com",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:auth::oidc/path/oidcProvider1",
			},
			getOidcProviderByNameMethodSpecialFunc: func(name string) (*OidcProvider, error) {
				if name == "oidcProvider1" {
					return &OidcProvider{
						ID:   "12345",
						Name: "oidcProvider1",
						Path: "/path/",
						Urn:  CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "oidcProvider1"),
					}, nil
				}

				return nil, &database.Error{
					Code: database.AUTH_OIDC_PROVIDER_NOT_FOUND,
				}
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "group1",
						Org:  "123",
						Path: "/new/",
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:         "POLICY-USER-ID",
						Name:       "policyUser",
						Org:        "123",
						Path:       "/path/",
						Urn:        CreateUrn("123", RESOURCE_POLICY, "/path/", "policyUser"),
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
		"ErrorCaseUpdateOidcProviderDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			oidcProviderName:    "oidcProvider1",
			newOidcProviderName: "newName",
			newPath:             "/new/",
			newIssuerUrl:        "http://test.com",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getOidcProviderByNameResult: &OidcProvider{
				ID:   "12345",
				Name: "oidcProvider1",
				Path: "/path/",
				Urn:  CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/path/", "oidcProvider1"),
			},
			updateOidcProviderMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
	}

	for x, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[UpdateOidcProviderMethod][0] = testcase.updateOidcProviderResult
		testRepo.ArgsOut[UpdateOidcProviderMethod][1] = testcase.updateOidcProviderMethodErr
		testRepo.ArgsOut[GetOidcProviderByNameMethod][0] = testcase.getOidcProviderByNameResult
		testRepo.ArgsOut[GetOidcProviderByNameMethod][1] = testcase.getOidcProviderByNameErr
		testRepo.SpecialFuncs[GetOidcProviderByNameMethod] = testcase.getOidcProviderByNameMethodSpecialFunc
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult

		oidcProvider, err := testAPI.UpdateOidcProvider(testcase.requestInfo, testcase.oidcProviderName, testcase.newOidcProviderName,
			testcase.newPath, testcase.newIssuerUrl, testcase.newClients)
		checkMethodResponse(t, x, testcase.wantError, err, testcase.expectedOidcProvider, oidcProvider)
	}
}

func TestWorkerAPI_RemoveOidcProvider(t *testing.T) {
	testcases := map[string]struct {
		//API method args
		requestInfo RequestInfo
		name        string
		// Expected result
		wantError error
		// Manager Results
		getUserByExternalIDResult         *User
		getGroupsByUserIDResult           []TestUserGroupRelation
		getAttachedPoliciesResult         []TestPolicyGroupRelation
		getOidcProviderByNameMethodResult *OidcProvider
		// API Errors
		getUserByExternalIDMethodErr   error
		getOidcProviderByNameMethodErr error
		removeOidcProviderMethodErr    error
		getGroupsByUserIDError         error
	}{
		"OKCaseAdminUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			name: "oidcProvider1",
			getOidcProviderByNameMethodResult: &OidcProvider{
				ID:   "543210",
				Name: "oidcProvider1",
				Path: "/example/",
				Urn:  CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/example/", "oidcProvider1"),
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
			},
		},
		"OkCase": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "oidcProvider1",
			getOidcProviderByNameMethodResult: &OidcProvider{
				ID:   "543210",
				Name: "oidcProvider1",
				Path: "/example/",
				Urn:  CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/example/", "oidcProvider1"),
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("org1", RESOURCE_USER, "/example/", "123456"),
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Org:  "org1",
						Path: "/example/",
						Urn:  CreateUrn("org1", RESOURCE_POLICY, "/example/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									AUTH_OIDC_ACTION_GET_PROVIDER,
									AUTH_OIDC_ACTION_DELETE_PROVIDER,
								},
								Resources: []string{
									GetUrnPrefix("", RESOURCE_AUTH_OIDC_PROVIDER, ""),
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
						Org:  "org1",
						Urn:  CreateUrn("org1", RESOURCE_GROUP, "/example/", "group1"),
					},
				},
			},
		},
		"ErrorCaseInvalidName": {
			name: "invalid*",
			wantError: &Error{
				Code:    INVALID_PARAMETER_ERROR,
				Message: "Invalid parameter: name invalid*",
			},
		},
		"ErrorCaseOidcProviderNotFound": {
			name: "oidcProvider1",
			wantError: &Error{
				Code: AUTH_OIDC_PROVIDER_BY_NAME_NOT_FOUND,
			},
			getOidcProviderByNameMethodErr: &database.Error{
				Code: database.AUTH_OIDC_PROVIDER_NOT_FOUND,
			},
		},
		"ErrorCaseUnauthorizedUser": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "oidcProvider1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "Authenticated user with externalId 123456 not found. Unable to retrieve permissions.",
			},
			getOidcProviderByNameMethodResult: &OidcProvider{
				ID:   "543210",
				Name: "oidcProvider1",
				Path: "/example/",
				Urn:  CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/example/", "oidcProvider1"),
			},
			getUserByExternalIDMethodErr: &database.Error{
				Code: database.USER_NOT_FOUND,
			},
		},
		"ErrorCaseImplicitUnauthorizedDeleteOidcProvider": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "oidcProvider1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:auth::oidc/example/oidcProvider1",
			},
			getOidcProviderByNameMethodResult: &OidcProvider{
				ID:   "543210",
				Name: "oidcProvider1",
				Path: "/example/",
				Urn:  CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/example/", "oidcProvider1"),
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("org1", RESOURCE_USER, "/example/", "123456"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/example/",
						Org:  "org1",
						Urn:  CreateUrn("org1", RESOURCE_GROUP, "/example/", "group1"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Org:  "org1",
						Path: "/example/",
						Urn:  CreateUrn("org1", RESOURCE_POLICY, "/example/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									AUTH_OIDC_ACTION_GET_PROVIDER,
								},
								Resources: []string{
									GetUrnPrefix("", RESOURCE_AUTH_OIDC_PROVIDER, ""),
								},
							},
						},
					},
				},
			},
		},
		"ErrorCaseExplicitUnauthorizedDeleteOidcProvider": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      false,
			},
			name: "oidcProvider1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:auth::oidc/example/oidcProvider1",
			},
			getOidcProviderByNameMethodResult: &OidcProvider{
				ID:   "543210",
				Name: "oidcProvider1",
				Path: "/example/",
				Urn:  CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/example/", "oidcProvider1"),
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("org1", RESOURCE_USER, "/example/", "123456"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/example/",
						Org:  "org1",
						Urn:  CreateUrn("org1", RESOURCE_GROUP, "/example/", "group1"),
					},
				},
			},
			getAttachedPoliciesResult: []TestPolicyGroupRelation{
				{
					Policy: &Policy{
						ID:   "POLICY-USER-ID",
						Name: "policyUser",
						Org:  "org1",
						Path: "/example/",
						Urn:  CreateUrn("org1", RESOURCE_POLICY, "/example/", "policyUser"),
						Statements: &[]Statement{
							{
								Effect: "allow",
								Actions: []string{
									AUTH_OIDC_ACTION_GET_PROVIDER,
									AUTH_OIDC_ACTION_DELETE_PROVIDER,
								},
								Resources: []string{
									GetUrnPrefix("", RESOURCE_AUTH_OIDC_PROVIDER, ""),
								},
							},
							{
								Effect: "deny",
								Actions: []string{
									AUTH_OIDC_ACTION_DELETE_PROVIDER,
								},
								Resources: []string{
									GetUrnPrefix("", RESOURCE_AUTH_OIDC_PROVIDER, "/example/oidcProvider1"),
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
			name: "oidcProvider1",
			wantError: &Error{
				Code:    UNAUTHORIZED_RESOURCES_ERROR,
				Message: "User with externalId 123456 is not allowed to access to resource urn:iws:auth::oidc/example/oidcProvider1",
			},
			getOidcProviderByNameMethodResult: &OidcProvider{
				ID:   "543210",
				Name: "oidcProvider1",
				Path: "/example/",
				Urn:  CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/example/", "oidcProvider1"),
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("org1", RESOURCE_USER, "/example/", "123456"),
			},
			getGroupsByUserIDResult: []TestUserGroupRelation{
				{
					Group: &Group{
						ID:   "GROUP-USER-ID",
						Name: "groupUser",
						Path: "/example/",
						Org:  "org1",
						Urn:  CreateUrn("org1", RESOURCE_GROUP, "/example/", "group1"),
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
		"ErrorCaseDeleteOidcProviderDBErr": {
			requestInfo: RequestInfo{
				Identifier: "123456",
				Admin:      true,
			},
			name: "oidcProvider1",
			wantError: &Error{
				Code: UNKNOWN_API_ERROR,
			},
			getOidcProviderByNameMethodResult: &OidcProvider{
				ID:   "543210",
				Name: "oidcProvider1",
				Path: "/example/",
				Urn:  CreateUrn("", RESOURCE_AUTH_OIDC_PROVIDER, "/example/", "oidcProvider1"),
			},
			getUserByExternalIDResult: &User{
				ID:         "123456",
				ExternalID: "123456",
				Path:       "/path/",
				Urn:        CreateUrn("", RESOURCE_USER, "/path/", "123456"),
			},
			removeOidcProviderMethodErr: &database.Error{
				Code: database.INTERNAL_ERROR,
			},
		},
	}

	for x, testcase := range testcases {
		testRepo := makeTestRepo()
		testAPI := makeTestAPI(testRepo)

		testRepo.ArgsOut[GetOidcProviderByNameMethod][0] = testcase.getOidcProviderByNameMethodResult
		testRepo.ArgsOut[GetOidcProviderByNameMethod][1] = testcase.getOidcProviderByNameMethodErr
		testRepo.ArgsOut[GetUserByExternalIDMethod][0] = testcase.getUserByExternalIDResult
		testRepo.ArgsOut[GetUserByExternalIDMethod][1] = testcase.getUserByExternalIDMethodErr
		testRepo.ArgsOut[GetGroupsByUserIDMethod][0] = testcase.getGroupsByUserIDResult
		testRepo.ArgsOut[GetGroupsByUserIDMethod][1] = testcase.getGroupsByUserIDError
		testRepo.ArgsOut[GetAttachedPoliciesMethod][0] = testcase.getAttachedPoliciesResult
		testRepo.ArgsOut[RemoveOidcProviderMethod][0] = testcase.removeOidcProviderMethodErr

		err := testAPI.RemoveOidcProvider(testcase.requestInfo, testcase.name)
		checkMethodResponse(t, x, testcase.wantError, err, nil, nil)
	}
}
