<img align="right" src="https://github.com/Tecsisa/foulkon/blob/master/dist/foulkon.png" width="150">
# Foulkon
[![Build Status](https://travis-ci.org/Tecsisa/foulkon.svg?branch=master)](https://travis-ci.org/Tecsisa/foulkon)
[![Join the chat at https://gitter.im/Tecsisa/foulkon](https://badges.gitter.im/Tecsisa/foulkon.svg)](https://gitter.im/Tecsisa/foulkon?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/tecsisa/foulkon)](https://goreportcard.com/report/github.com/tecsisa/foulkon)
[![codecov](https://codecov.io/gh/Tecsisa/foulkon/branch/master/graph/badge.svg)](https://codecov.io/gh/Tecsisa/foulkon)
[![Release](https://img.shields.io/badge/release-v0.2.0-blue.svg)](https://github.com/Tecsisa/foulkon/blob/master/CHANGELOG.md#v020-2016-10-06)
 
__Foulkon__ is an authorization server that allows or denies the access to web resources.

## Installation / usage

This project generates 2 apps:

- Worker: This is the authorization server itself.
- Proxy: This proxies the requests to the authorization server (worker).

Installation/deployment docs using Go binaries or Docker:<br />
- [Worker](doc/deploy/worker.md)
- [Proxy](doc/deploy/proxy.md)

## Documentation

Specification docs:
- [Specification](doc/spec/README.md)
- [Use case](doc/spec/usecase.md)
- [Internal IAM Actions](doc/spec/action.md)
- [Authorization flow](doc/spec/authorization.md)

API docs:
- [User](doc/api/user.md)
- [Group](doc/api/group.md)
- [Policy](doc/api/policy.md)
- [Resource](doc/api/resource.md)

You can also import this [Postman collection](schema/postman.json) file with all API methods.

## Limitations

Since validation is different in each identity provider, Foulkon needs __ID Token__ instead of __Access Token__ in order to check user permissions
in Authorization header with type bearer.
E.g.:

```
GET /example/resource HTTP/1.1
  Host: server.example.com
  Authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6IjFlOWdkazcifQ.ewogImlzcyI6ICJodHRwOi8vc2VydmVyLmV4YW1wbGUuY29tIiwKICJzdWIiOiAiMjQ4Mjg5NzYxMDAx
  IiwKICJhdWQiOiAiczZCaGRSa3F0MyIsCiAibm9uY2UiOiAibi0wUzZfV3pBMk1qIiwKICJleHAiOiAxMzExMjgxOTcwLAogImlhdCI6IDEzMTEyODA5NzAKfQ.ggW8hZ1EuVLuxNuuIJKX_V8
  a_OMXzR0EHR9R6jgdqrOOF4daGU96Sr_P6qJp6IcmD3HP99Obi1PRs-cwh3LO-p146waJ8IhehcwL7F09JdijmBqkvPeB2T9CJNqeGpe-gccMg4vfKjkM8FcGvnzZUN4_KSP0aAp1tOJ1zZwgj
  xqGByKHiOtX7TpdQyHE5lcMiKPXfEIQILVq0pc_E2DzL7emopWoaoZTF_m0_N0YzFC6g6EJbOEoRoSK5hoDalrcvRYLSrQAZZKflyuVCyixEoV9GfNQC3_osjzw2PAithfubEEBLuVVk4XUVrWO
  LrLl0nx7RkKU8NXNHq-rvKMzqg
```

## Development

For local development, first make sure Go is properly installed according to [Go install doc](https://golang.org/doc/install) (Also, include $GOBIN environment var in your $PATH). Then run next command in project root path:

```
$ make bootstrap
```

Afterwards type make test. This will run the tests. If this exits with exit status 0, then everything is working!

```
$ make test
```

To compile a development version of Foulkon, run make dev. This will put the Foulkon binaries in the $GOPATH/bin folder:

```
$ make dev
...
$ $GOPATH/bin/worker
$ $GOPATH/bin/proxy
...
```

## Contribution policy

Contributions via GitHub pull requests are gladly accepted from their original author. Along with any pull requests, please state that the contribution is your original work and that you license the work to the project under the project's open source license. Whether or not you state this explicitly, by submitting any copyrighted material via pull request, email, or other means you agree to license the material under the project's open source license and warrant that you have the legal authority to do so.

Please make sure to follow these conventions:
- For each contribution there must be a ticket (GitHub issue) with a short descriptive name, e.g. "run go imports in Makefile"
- Work should happen in a branch named "ISSUE-DESCRIPTION", e.g. "32-go-imports-in-Makefile"
- Before a PR can be merged, all commits must be squashed into one with its message made up from the ticket name and the ticket id, e.g. "better go files formatting: run go imports in Makefile (closes #32)"

#### Questions

If you have a question, preferably use Gitter chat. As an alternative, prepend your issue with `[question]`.

## License

This code is open source software licensed under the [Apache 2.0 License](http://www.apache.org/licenses/LICENSE-2.0.html).
