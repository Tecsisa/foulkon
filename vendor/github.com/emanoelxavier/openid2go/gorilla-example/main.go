package main

import (
	"fmt"
	"net/http"

	"github.com/emanoelxavier/openid2go/openid"
	"github.com/gorilla/context"
)

const UserKey = 0

func authenticatedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "The user was authenticated successfully!")
}

func unauthenticatedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Reached without authentication!")
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, UserKey).(*openid.User)
	fmt.Fprintf(w, "Hello %v! This is all I know about you: %+v", u.ID, u)
}

type userHandlerAdapter struct {
	h http.HandlerFunc
}

func (uh userHandlerAdapter) ServeHTTPWithUser(u *openid.User, rw http.ResponseWriter, req *http.Request) {
	context.Set(req, UserKey, u)
	uh.h.ServeHTTP(rw, req)
}

func main() {
	configuration, _ := openid.NewConfiguration(openid.ProvidersGetter(getProviders_googlePlayground))

	handlerAdapter := userHandlerAdapter{h: meHandler}

	http.Handle("/me", openid.AuthenticateUser(configuration, handlerAdapter))
	http.Handle("/authn", openid.Authenticate(configuration, http.HandlerFunc(authenticatedHandler)))
	http.HandleFunc("/", unauthenticatedHandler)

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
