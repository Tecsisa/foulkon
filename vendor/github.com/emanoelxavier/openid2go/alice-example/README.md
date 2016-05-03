Go OpenId - Alice Example
===========
[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/emanoelxavier/openid2go/openid)
[![license](http://img.shields.io/badge/license-MIT-yellowgreen.svg?style=flat)](https://raw.githubusercontent.com/emanoelxavier/openid2go/master/alice-example/LICENSE)

This fully working example implements an HTTP server chaining openid Authentication middlewares with various other middlewares using [Alice](https://github.com/justinas/alice).

Alice allows easily chaining middlewares in the form:

```go
func (http.Handler) http.Handler
```

However the Authentication middlewares exported by the package openid2go/openid have slightly different constructors:

```go
func Authenticate(conf *Configuration, h http.Handler) http.Handler
```
```go
func AuthenticateUser(conf *Configuration, h UserHandler) http.Handler
```

This example demonstrates that those middlewares can still be easily chained using Alice with little additional code.

## Test

Download and build:
```sh
go get github.com/emanoelxavier/openid2go/alice-example
```
```sh
go build github.com/emanoelxavier/openid2go/alice-example
```

Run:
```sh
github.com\emanoelxavier\openid2go\alice-example\alice-example.exe
```

Once running you can send requests like the ones below:
```sh
GET http://localhost:5103
```
```sh
GET http://localhost:5103/me
Authorization: Bearer eyJhbGciOiJS...
````
```sh
GET http://localhost:5103/authn
Authorization: Bearer eyJhbGciOiJS...
```
The abbreviated token above must be replaced with the IDToken acquired from the [Google OAuth PlayGround](https://developers.google.com/oauthplayground) entering "openid" (without quotes) within the scope field.