Go OpenId - Gorilla Example
===========
[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/emanoelxavier/openid2go/openid)
[![license](http://img.shields.io/badge/license-MIT-yellowgreen.svg?style=flat)](https://raw.githubusercontent.com/emanoelxavier/openid2go/master/gorilla-example/LICENSE)

This fully working example implements an HTTP server using openid Authentication middlewares and [Gorilla Context](http://www.gorillatoolkit.org/pkg/context) to preserve the user information accross the service application stack.


The AuthenticateUser middleware exported by the package openid2go/openid forwards the user information to a handler that implements the interface openid.UserHandler:


```go
func AuthenticateUser(conf *Configuration, h UserHandler) http.Handler
```

```go
type UserHandler interface {
	ServeHTTPWithUser(*User, http.ResponseWriter, *http.Request)
}
```

This example demonstrates how to create an adapter that implements that interface and use it to store the openid.User into a Gorilla Context. The user information can then be retrieved from the context in another point of the application stack.

## Test

Download and build:
```sh
go get github.com/emanoelxavier/openid2go/gorilla-example
```
```sh
go build github.com/emanoelxavier/openid2go/gorilla-example
```

Run:
```sh
github.com\emanoelxavier\openid2go\alice-example\gorilla-example.exe
```

Once running you can send requests like the ones below:
```sh
GET http://localhost:5100
```
```sh
GET http://localhost:5100/me
Authorization: Bearer eyJhbGciOiJS...
````
```sh
GET http://localhost:5100/authn
Authorization: Bearer eyJhbGciOiJS...
```
The abbreviated token above must be replaced with the IDToken acquired from the [Google OAuth PlayGround](https://developers.google.com/oauthplayground) entering "openid" (without quotes) within the scope field.
