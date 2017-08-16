# Foulkon demo

This demo shows how Foulkon works, and how to manage it.

## Previous requirements

To run this demo, you have to set some properties and system packages.

### Go configuration

We have to set next environment vars:

 - GOROOT
 - GOPATH
 - GOBIN

On Ubuntu (Directory examples, you can choose your own directories):

```bash

export GOROOT=$HOME/dev/golang/go
export GOPATH=$HOME/dev/sources/golang
export GOBIN=$HOME/dev/sources/golang/bin

```

This directories will be created before. Also, the GOBIN environment variable 
will be in your execution path.

### System packages

This demo works with Docker, so you have to install Docker and Docker Compose.

 - [Docker installation doc](https://docs.docker.com/engine/installation/)
 - [Docker Compose installation doc](https://docs.docker.com/compose/install/)
 
## Start Demo

First, you have to download Foulkon project:

```bash

go get github.com/Tecsisa/foulkon

```

Second, go to Foulkon directory:

```bash

cd $GOPATH/src/github.com/Tecsisa/foulkon

```

Third, execute next command to get all dependencies:

```bash

make bootstrap

```

User login needs a Google client to make UI able to get a user.
To do this, follow the [Google guide](https://developers.google.com/identity/protocols/OpenIDConnect) to get a Google client
set http://localhost:8101/callback in your Authorized redirect URIs, and change next properties in [Docker-compose file](docker/docker-compose.yml):

 - In foulkonworkercompose: 
    - FOULKON_AUTH_CLIENTID for your client id.
 - In foulkondemowebcompose:
    - OIDC_CLIENT_ID for your client id.
    - OIDC_CLIENT_SECRET for your secret.
    
Finally, execute demo command to start demo:

```bash

make bootstrap

```

The applications started are next:

 - Worker: Started on http://localhost:8000
 - Proxy: Started on http://localhost:8001
 - API demo: Started on http://localhost:8100
 - UI demo: Started on http://localhost:8101

Now, you have all suite to try Foulkon, go to [Tour](tour.md) to see an example.

