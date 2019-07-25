# restarter
Restart applications.

It provides an HTTP API to trigger the restart of the pods of a given
application (based on the host).

**IMPORTANT**: This is an internal service, don't expose this on the internet.

## Usage

A Makefile is provided to enable easily building, testing and running
of the service.

### `make help`
Show list of available "commands" (targets)

### `make run`
Compiles and runs the service on `http://localhost:8080` (or the
specified `$PORT`)

**NOTE**: This will mount your `~/.kube/config` file so the running container
will use your current k8s context/credentials.

### `make test`
Compiles the test code and runs it

### `make static` (default)
Statically compile the service binary.

### `make docker-image`
Builds a docker image as defined in Dockerfile

### `make docker-run`
Builds and runs the service in a docker container


## Configuration
The application doesn't require any configuration to work.

| Env variable         | Default |  Details |
| -------------------- | ------- | -------- |
| `PORT`               | `8000`  | port on which the server listen to |

**NOTE**: The server will try to load the kubernetes configuration from
in-cluster first (this is the case when running the server within a k8s
cluster) and fallback to load it from `$HOME/.kube/config` when this fails.

If that fails as well the server will not start.


### Endpoints

#### `POST /restart`

Restarts the application with the given host.
Restart requests have to have the following JSON format:

```json
{
    "host": "example.com",          // Required.
    "reason": "reason for restart"  // Optional. Default: "data-updated"
}
```

The restart is triggerred by adding the following annotations to the `Deployment`
and to its pods's template (with the current time and the given reason):

```json
{
    "annotations": {
        "restartedAt": "2019-07-25T14:17:33Z",
        "restartReason": "My dog ate my homework"
    }
}
```

##### Example
```bash
$ curl -v -XPOST --data '{"host": "example.com", "reason": "My dog ate my homework"}' 127.0.0.1:8000/restart
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 8000 (#0)
> POST /restart HTTP/1.1
> Host: 127.0.0.1:8000
> User-Agent: curl/7.54.0
> Accept: */*
> Content-Length: 83
> Content-Type: application/x-www-form-urlencoded
>
* upload completely sent off: 83 out of 83 bytes
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Thu, 25 Jul 2019 14:17:33 GMT
< Content-Length: 25
<
{"message":"Restarted."}
```

### `GET /healthz` (healthcheck)
This will respond with a `200 OK` and a small body.
It's used by kubernetes (or whatever) to check that the server is still
responding.

#### Example
```bash
$ curl -v -XGET 127.0.0.1:8000/healthz
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 8000 (#0)
> GET /healthz HTTP/1.1
> Host: 127.0.0.1:8000
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Thu, 25 Jul 2019 14:16:23 GMT
< Content-Length: 26
<
{"message":"OK 👍🏼"}
```


## Dependencies

Dependencies are managed using [Go Modules](https://github.com/golang/go/wiki/Modules).

Dependences are vendored in the `/vendor` which is checked in Git.


### Add a new dependency

1. `$ go get foo/bar`
2. Edit your code to import foo/bar

### Upgrade a dependency

As per instructions [here](https://github.com/golang/go/wiki/Modules#how-to-upgrade-and-downgrade-dependencies)

1. `$ go get foo/bar`

This will upgrade to the latest version of `foo/bar` with a semver tag.
Alternatively, `go get foo/bar@v1.2.3` will get a specific version.

## Docker image
The [`Dockerfile`](/) uses 2 stages one for building and the final image.

### builder stage

### final stage
The actual image running the service is just scratch with the binary compiled
statically (`-ldflags '-extldflags "-static"'`) to keep the docker image to the
bare minimum.

See this article on containerising Go application: https://www.cloudreach.com/blog/containerize-this-golang-dockerfiles/
