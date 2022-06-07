
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

Working on adding functionality:

- Move `cmd/integration` into `cmd/beam/integration`
	- _integration tests are dependent on the binary produced, so putting it into `cmd/` directly was a mistake_
- Move handlers and storage to `pkg/{api,postgres}`
- Document the fact that `pkg/postgres` initialization helps control order dependent operations
	- _for example, relational tables must be created after the tables they relate to, which is harder to control if kept loosely attached to each models storage component._
- Rename to match model (_no longer need long-names_)
- Abstract api from postgres dependencies via interfaces
	- Add note to api login operation noting the pre/post handling could be added to secure
		- pre-login rate limit and counter check, post-failure expiring counter increment
			- something like redis would be best for this; using postgres gets messy/expensive
	- _need to find best method to abstract Token from User with regards to Login_
- Identify best pattern to combine token and user authentication operations
	- _perhaps a composite model and interface expecting two independent operations?_
- Update `cmd/beam` assembly to use the new `pkg` imports and assemble components
	- _Utilize inversion-of-control from `pkg/api` to accept `httprouter` and attach routes_
- new routes (replacing some old ones)
	- `/oauth/token` used instead of `/login`; accepts Basic & Bearer authentication
		- Bearer checks refresh token (eg. a hash /w base64 user prefix), _not a jwt_
		- add another note that the base64 user prefix allows us to add access token rate-limiting
	- `/oauth/revoke`; an alias to delete a refresh token, but which must check ownership?
	- add `/services` and `/users` routes to demonstrate `html/template` stdlib package
		- _render a basic html table of all results; pagination could be added in future iterations_
		- _this is primarily to demonstrate server-side rendering as an alternative to react & api calls._

- revisit auth behavior by adding shared secret based jwt validation wrapper
	- _If no shared secret is provided use `crypto/rand` to generate one at launch_
	- _we can add notes that keypair jwt is also possible, which would work for authentication with separated systems by exposing a `/public.key` route allowing third parties to validate our systems tokens._


- define tests in `cmd/beam/integration`
	- _accept `ADDRESS` and `CREDENTIALS` env vars_
		- _allows it to be run against any environment_
		- _can use predefined credentials for testing_
			- we may want to support different tiers of credentials (eg. to validate users with varying permission levels)
	- Cleanup cycle should be executed both pre and post execution
		- _pre-cleanup avoids canceled tests with leftover/orphaned resources_

- begin writing reactjs into main.js to produce SPA UI
	- _we can replace server-side rendered pages_


## developer reference commands

Curl to create a user:

	curl -X POST -d "username=username&email=user@gmail.com&password=password" http://localhost:3000/register

Can collect a refresh token:

	REFRESH_TOKEN=$(curl -X POST -u username:password http://localhost:3000/login 2> /dev/null)

Can send that refresh token to get an access token:

	ACCESS_TOKEN=$(curl -h "Authorization: Refresh $REFRESH_TOKEN" http://localhost:3000/oauth/token)

Can create a service using that access token:

	curl -X POST -h "Authorization: Bearer $ACCESS_TOKEN" http://localhost:3000/api/service -d '{"name":"","technology":"go","poc":"user@email.com"}'

_Substitute with username and password stored in the database or it will fail with 401._


# refrences

- [Refresh Tokens](https://www.oauth.com/oauth2-servers/making-authenticated-requests/refreshing-an-access-token/)

Sample json response from Refresh Tokens reference:

	{
	  "access_token": "BWjcyMzY3ZDhiNmJkNTY",
	  "refresh_token": "Srq2NjM5NzA2OWJjuE7c",
	  "token_type": "Bearer",
	  "expires": 3600
	}

_Not sure if passing expires matters, since the service has to detect expired tokens from 401 response anyways._

- [Storing Passwords Securely With PostgreSQL and Pgcrypto](https://x-team.com/blog/storing-secure-passwords-with-postgresql/)


---

Just adding a note here that I did try to write a modified Users table and query that could conditionally update a timestamp and counter to deal with failed logins, but there was no sane way to pull tokens using it, and the query itself was pretty ugly to look at:

	UPDATE users SET
		failed_logins = sq.new_failed_logins,
		failed_login_at = sq.new_failed_login_at
	FROM (SELECT
	CASE WHEN failed_logins > 3
		OR transaction_timestamp() < failed_login_at + (users.failed_logins * interval '1 minute')
		OR password != crypt('abc123', password)
		THEN failed_logins + 1
		ELSE 0
	END new_failed_logins,
	CASE WHEN failed_logins > 3
		OR transaction_timestamp() < failed_login_at + (users.failed_logins * interval '1 minute')
		OR password != crypt('abc123', password)
		THEN transaction_timestamp()
		ELSE failed_login_at
	END new_failed_login_at,
	CASE WHEN failed_logins > 3
		OR transaction_timestamp() < failed_login_at + (users.failed_logins * interval '1 minute')
		OR password != crypt('abc123', password)
		THEN false
		ELSE true
	END success
	FROM users
	WHERE email = 'user@email.com') AS sq
	WHERE email = 'user@email.com'
	RETURNING sq.success;

_To keep the SQL sane it would make far more sense to properly secure the login route with a redis or cache dbms system for performance and separate checks to rate limit or block on a number failed logins with an expiring counter, which is something redis has built-in support for._
