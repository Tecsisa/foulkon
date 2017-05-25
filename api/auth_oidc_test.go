package api

import (
	"testing"

	"github.com/Tecsisa/foulkon/database"
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
