# Contributing

## Guidelines for pull requests

- Write tests for any changes.
- Separate unrelated changes into multiple pull requests.
- For bigger changes, make sure you start a discussion first by creating an issue and explaining the intended change.
- Use [conventional changelog conventions](https://github.com/bcoe/conventional-changelog-standard/blob/master/convention.md) in your commit messages.

## Development dependencies

- go 1.13.6 or higher
- docker 18.06.1-ce or higher (multistage build is required)
- [artillery](https://artillery.io/) for performance tests

## Setting up a development machine

1. Install project dependencies

```
make install
```

If your IDE does not support go modules or the project is inside GOPATH, run `go mod vendor` to put dependencies into a vendor folder

2. Verify your environment is working (format, lint code, run all tests)
```
make verify
```

## Performance tests

Performance tests are performed using [artillery](https://artillery.io/). Follow instructions under https://artillery.io/docs/getting-started/ to install artillery.

Configuration for performance tests is under `test/performance/config.yaml`, which consist of two phases:

- warm up phase ramping up from 10 to 30 concurrent users over 60 seconds
- sustained load of 30 concurrent users over 5 minutes

To run performance tests execute:

```
ETHEREUM_API_BASE_URL=http://localhost:8080 make test-performance
```

This command will create a json report under `test/performance/<timestamp>` folder, which can be used to generate an html report through `artillery report <path-to-file>`

## During development

Commits to this codebase should follow the [conventional changelog conventions](https://github.com/bcoe/conventional-changelog-standard/blob/master/convention.md).

- `make verify` - A check to be run before pushing any changes
- `make build` - builds a server executable under `build` directory
- `make docker` - builds the server and creates a docker image for release
- `make docker-publish` - publishes docker image
- `make watch` - runs unit tests on every change
