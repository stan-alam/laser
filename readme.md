
# laser

A demo project.


## design

This project can be broken down into the following components:

- [directory structures following a standard](https://github.com/golang-standards/project-layout)
- Multi-stage Docker Build for lightweight distribution
- Embedding a frontend with the builtin embed package

The `embed` package does not allow parent directory relative paths, so the `embed.go` file at the root is where the import occurs to match the `web/` directory structure standard.

For `go.mod` you can run `go get ...` at the root to pull in dependencies, allowing you to build all sub-directory components.  You can also update `go/mod` by running `go get -u ... && go mod tidy`.

The separated `cmd/` structure allows you to create multiple binaries off the same library components, which could be used to separate functionality into `X/` following the directory standards making for very modular coding that works for monolithic and independent binaries.

We've included `cmd/integration` to demonstrate an independent integration testing package, versioned alongside the project.  _The only potential concern would be if multiple binaries with separate functionality were created, in which case you might create an `cmd/*/integration/` directory structure, and change the Dockerfile build instructions as needed._

The root `Dockerfile` produces a final self-contained image with the `cmd/beam` binary and both unit and integration tests, allowing you to run the tests on itself, giving you excellent synchronized versioning capabilities.

Routing demonstrates hard-coded JSON response at `/health`, dynamic responses behind `/api`, template responses at `/users` and `/services`, and finally an embedded file-based front-end when no explicitly registered routes match, including the index, _which could eventually be turned into a REACT SPA._


## instructions

**Go is a cross platform compiled language, which means it produces a binary, so for local development you can just run the code without containerization, otherwise every change will require you to recompile the binary.**

Start at the root and install dependencies leveraging go mod:

	go get ...

Quick test and build from `cmd/beam`:

	go test -v
	go build
	./beam

_Easily check status from [`/health`](http://localhost:3000/health)._


### postgresql

For local development you may need to [enable ssl](https://www.postgresql.org/docs/current/ssl-tcp.html) or change the connection string to explicitly disable it (eg. `postgres://username:password@localhost/db_name?sslmode=disable`).


### docker

The [Dockerfile](Dockerfile) produces a production scratch container with the beam binary, test binary, and integration binary.

This means the docker image is self-contained and can prove that its own code works by running its own unit tests, or running integration tests pointed at itself.

To build the image and tag by git hash:

	docker build . -t beam:$(git rev-parse --short HEAD)

_We'll assume the `test` tag for shorter subsequent commands._

Run the unit tests　from the test binary within the container:

	docker run --rm beam:test /go/bin/beam.test -test.v

Run the benchmarks from the test binary within the container:

	docker run --rm beam:test /go/bin/beam.test -test.v -test.run=X -test.benchtime=10s -test.benchmem -test.bench=.

Run the default command and access it from [`http://localhost:3000/`](http://localhost:3000/):

	docker run --rm beam:test

Run the integration tests from the integration binary within the container, and point at the local running copy:

	docker run --rm -e "ADDRESS=http://localhost:3000/" beam:test go/bin/integration

_The majority of these commands can be fully automated from a CICD system, and the outputs can be used to make decisions such as failing unit tests stops before deployment, failing integration tests in one environment prevents promotion to the next environment, and so on..._

Further, anyone should be able to run the integration command against any environment they have access to, allowing you to rapidly validate whether an environment is functioning correctly.


## todo

- add `jwt` generation and abstraction to `pkg/api`
	- add shared secret to `cmd/beam/config.go`
		- if no secret exists, generate one using `crypto/rand`
- define custom claims in `pkg/api` with permissions
- add permission check wrapper, combined with cors wrapper to `pkg/api/auth.go`

- Debatable cleanup tasks:
	- do we move `cmd/beam/health.go` into `pkg/api/`?
		- _probably not, too much cruft to build health struct, init for time tracking, and passing main.Version._
	- do we abstract httprouter in `pkg/api`?
		- _probably not, the behavior differs too much_
	- do we move or replicate relevant documentation into `cmd/beam/readme.md`?
	- do we swap shared secret for key pair jwt?
		- _can generate keypair when not supplied_
		- _can add `/public.key` route to allow external signature validation from separated systems._

- define tests in `cmd/beam/integration`
	- _accept `ADDRESS` and `CREDENTIALS` env vars_
		- _allows it to be run against any environment_
		- _can use predefined credentials for testing_
			- we may want to support different tiers of credentials (eg. to validate users with varying permission levels)
	- Cleanup cycle should be executed both pre and post execution
		- _pre-cleanup avoids canceled tests with leftover/orphaned resources_

- begin writing reactjs into main.js to produce SPA UI
	- _we can replace server-side rendered pages_


# refrences

- [Refresh Tokens](https://www.oauth.com/oauth2-servers/making-authenticated-requests/refreshing-an-access-token/)
- [Storing Passwords Securely With PostgreSQL and Pgcrypto](https://x-team.com/blog/storing-secure-passwords-with-postgresql/)
