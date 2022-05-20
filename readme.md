
# laser

A demo project.


## design

This project can be broken down into the following components:

- [directory structures following a standard](https://github.com/golang-standards/project-layout)
	- project root has shared structures and logic
	- `cmd/` for binaries
		- `go.mod` & `package main` in each
		- `integration/` to demo integration testing
- Multi-stage Docker Build for lightweight distribution
	- The final build includes binary & tests (eg. self-proving)
- Embedding frontend with go-bindata
	- Routing demonstrates hardcoded json response, template based page, and file-based delivery


## instructions

Go is a compiled language, so it produces a binary.

This means you cannot simply run a container and change files, you need to restart the container, or rebuild the production scratch container for every change.

**However, because go is compatible on many platforms and can be cross compiled you can just run it locally without the need for dockerization.**

Quick test and build from `cmd/beam`:

	go test -v -race
	go build
	./beam

_Easily check status from [`/health`](http://localhost:3000/health)._

The [Dockerfile](Dockerfile) produces a production scratch container with the beam binary, test binary, and integration binary.

This means the docker image is self-contained and can prove that its own code works by running its own unit tests, or running integration tests pointed at itself.

Here are some docker commands useful for your CICD system:

	TODO

In a CICD system:

- the Dockerfile runs the tests before it builds, so if they fail no build is produced.
- Post-deployment you can run the integration command against the deployed address to validate it


## todo

Working on adding functionality:

- Add go-bindata and deliver files when no registered routes are found
	- _refer to CDN hosted reactjs to reduce dependency load_
- REST endpoints at `/api/service`
	- _Add UI code to support REST endpoints_
- fill out `cmd/integration` with test code that validates REST behavior
	- _self-cleaning, eg. get (should not exist), add, get, update, get, delete, get (should not exist)_
- Introduce authentication
- Introduce authorization
