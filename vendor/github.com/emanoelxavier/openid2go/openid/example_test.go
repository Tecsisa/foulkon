package openid_test

import (
	"fmt"
	"net/http"

	"github.com/emanoelxavier/openid2go/openid"
)

func AuthenticatedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "The user was authenticated!")
}

func AuthenticatedHandlerWithUser(u *openid.User, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "The user was authenticated! The token was issued by %v and the user is %+v.", u.Issuer, u)
}

func UnauthenticatedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Reached without authentication!")
}

// This example demonstrates how to use of the openid middlewares to validate incoming
// ID Tokens in the HTTP Authorization header with the format 'Bearer id_token'.
// It initializes the Configuration with the desired providers (OPs) and registers two
// middlewares: openid.Authenticate and openid.AuthenticateUser.
// The former will validate the ID Token and fail the call if the token is not valid.
// The latter will do the same but forward the user's information extracted from the token to the next handler.
func Example() {
	configuration, err := openid.NewConfiguration(openid.ProvidersGetter(getProviders_googlePlayground))

	if err != nil {
		panic(err)
	}

	http.Handle("/user", openid.AuthenticateUser(configuration, openid.UserHandlerFunc(AuthenticatedHandlerWithUser)))
	http.Handle("/authn", openid.Authenticate(configuration, http.HandlerFunc(AuthenticatedHandler)))
	http.HandleFunc("/unauth", UnauthenticatedHandler)

	http.ListenAndServe(":5100", nil)
}

// getProviders returns the identity providers that will authenticate the users of the underlying service.
// A Provider is composed by its unique issuer and the collection of client IDs registered with the provider that
// are allowed to call this service.
// On this example Google OP is the provider of choice and the client ID used corresponds
// to the Google OAUTH Playground https://developers.google.com/oauthplayground
func getProviders_googlePlayground() ([]openid.Provider, error) {
	provider, err := openid.NewProvider("https://accounts.google.com", []string{"407408718192.apps.googleusercontent.com"})

	if err != nil {
		return nil, err
	}

	return []openid.Provider{provider}, nil
}
