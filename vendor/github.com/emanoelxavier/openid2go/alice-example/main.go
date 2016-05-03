package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/emanoelxavier/openid2go/openid"
	"github.com/gorilla/context"
	"github.com/justinas/alice"
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
	fmt.Fprintf(w, "Hello %v! this is all I know about you: %+v.", u.ID, u)
}

func timeoutMiddleware(h http.Handler) http.Handler {
	return http.TimeoutHandler(h, 1*time.Second, "timed out")
}

func myMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		fmt.Fprint(w, "Executed my middleware!")
	})
}

var provider *openid.Provider
var configuration *openid.Configuration

type userHandlerAdapter struct {
	h http.Handler
}

func (a *userHandlerAdapter) ServeHTTPWithUser(u *openid.User, rw http.ResponseWriter, req *http.Request) {
	context.Set(req, UserKey, u)
	a.h.ServeHTTP(rw, req)
}

func (a *userHandlerAdapter) myAuthenticateUser(h http.Handler) http.Handler {
	a.h = h
	return openid.AuthenticateUser(configuration, a)
}

func myAuthenticate(h http.Handler) http.Handler {
	return openid.Authenticate(configuration, h)
}

func main() {
	configuration, _ = openid.NewConfiguration(openid.ProvidersGetter(getProviders_googlePlayground))

	adapter := new(userHandlerAdapter)

	http.Handle("/me", alice.New(timeoutMiddleware, myMiddleware, adapter.myAuthenticateUser).ThenFunc(meHandler))
	http.Handle("/authn", alice.New(timeoutMiddleware, myMiddleware, myAuthenticate).ThenFunc(authenticatedHandler))
	http.HandleFunc("/", unauthenticatedHandler)

	http.ListenAndServe(":5103", nil)
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
