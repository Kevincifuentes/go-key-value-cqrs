[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
# go-key-value-cqrs

# Index
* [Description](#description)
* [Run server](#run-server)
  * [Run server Locally](#run-server-locally)
  * [Run server using docker](#run-server-using-docker)
* [Run tests](#run-tests)
  * [Run tests locally](#run-tests-locally)
  * [Run tests using docker](#run-tests-using-docker)
* [Debug Mode](#debug-mode)

# Description

go-key-value-cqrs is an example project to implement CQRS, DDD and Hexagonal Architecture using Golang programming language.

It will support a Key/Value storage service with REST API capabilities to:
* Get KeyValue by Key
* Post a new KeyValue
* Delete a KeyValue

# Run server

There are two ways of running the server: simply [locally](#run-server-locally) using the code or using [docker](#run-server-using-docker):

## Run server Locally
To run the server, you don't need any configuration because it uses default values but if you want to configure it, 
you will need a `.env` file with at least the following properties:
```txt
SERVER_HOST=localhost
SERVER_PORT=8080
OPENAPI_RELATIVE_PATH=./api/keyvalue/api.yml
DEBUG_SERVER_HOST=localhost
DEBUG_SERVER_PORT=8081
```
* SERVER_HOST: usually localhost or 127.0.0.1
* SERVER_PORT: in which port will the server run
* OPENAPI_RELATIVE_PATH: this represents the path where the OPENAPI yml file is stored. It is used for Request validation
middleware.
* DEBUG_SERVER_HOST (ONLY needed if `debug` tag is used): similar to SERVER_HOST but for profiling
* DEBUG_SERVER_PORT (ONLY needed if `debug` tag is used): similar to SERVER_PORT but for profiling

Then, you can simply run: 
```bash
go run server.go
```
or on [DEBUG MODE](#debug-mode):
```bash
go run -tags debug server.go
```

## Run server using docker

If you have docker installed, this is the most straightforward way of starting the KeyValue server, as well as running
all tests (getting the report/coverage for them). The repository includes a `Dockerfile` that describes all the stages
that can be run:
* Build project image
* Run all the tests
* Get the tests report

The file also describes all the steps needed to perform all those stages without the necessity to have any dependencies
nor libraries installed beforehand.

We can **build the project image** with the following command:

```bash
docker build -t keyvalueserver-app .
```

If you want to enable the [DEBUG_MODE](#debug-mode), we can run the following instead:

```bash
docker build -t keyvalueserver-app-debug . --build-arg DEBUG_MODE=true 
```

> [!TIP]
> You can still add a `.env` file as stated on [Run server locally chapter](#run-server-locally) with all those variables.
> The Dockerfile will consider that file if it's present at the same level.

Then, we can simply run our application like so (using `keyvalueserver-app` or `keyvalueserver-app-debug`):

``` bash
docker run -p 8080:8080 keyvalueserver-app
```

or

```bash
docker run -p 8080:8080 -p 8081:8081 keyvalueserver-app-debug 
```

You can still overwrite the environment variables if you didn't provide any `.env` passing them as `-env` option at the 
run command. For example, if we want to change the server's port:

```bash
docker run -p 9090:9090 -e SERVER_PORT=9090 keyvalueserver-app 
```

# Run tests

There are two ways of running the server: simply [locally](#locally) using the code or using [docker](#using-docker):

## Run tests locally

If you have your own environment ready with Go (1.23.4 version), you can directly run the following bash script to run 
all tests and generate coverage report on ``assets`` directory:
``` bash
bash ./runAllTests.sh
```

This will display a console report like the following extract:

```bash
Processing module /Users/mongatanga/Documents/Personal/Repositories/go-key-value-cqrs
?       main    [no test files]
Processing module /Users/mongatanga/Documents/Personal/Repositories/go-key-value-cqrs/internal/infrastructure
        go-key-value-cqrs/infrastructure/api/model              coverage: 0.0% of statements
        go-key-value-cqrs/infrastructure/api            coverage: 0.0% of statements
ok      go-key-value-cqrs/infrastructure/persistence    0.433s  coverage: 6.7% of statements in go-key-value-cqrs/...
Processing module /Users/mongatanga/Documents/Personal/Repositories/go-key-value-cqrs/internal/e2e
        go-key-value-cqrs/e2e/client            coverage: 0.0% of statements
ok      go-key-value-cqrs/e2e   0.202s  coverage: 37.9% of statements in go-key-value-cqrs/...
Processing module /Users/mongatanga/Documents/Personal/Repositories/go-key-value-cqrs/internal/domain
ok      go-key-value-cqrs/domain        0.198s  coverage: 87.0% of statements in go-key-value-cqrs/...
Processing module /Users/mongatanga/Documents/Personal/Repositories/go-key-value-cqrs/internal/application
ok      go-key-value-cqrs/application/queries/cqrs/querybus     0.183s  coverage: 88.2% of statements in go-key-value-cqrs/...
ok      go-key-value-cqrs/application/queries/keyvalue/getvalue 0.342s  coverage: 27.5% of statements in go-key-value-cqrs/...
```

When finished, you can access to ```go-key-value-cqrs/assets/coverage.html``` to see a user-friendly coverage report for 
all the modules.

## Run tests using docker

If we just want to get the test report, we can simply run the ``test-out`` stage like so:

``` bash
docker build --output="type=local,dest=./assets" --target=test-out .
```

The command above will generate a directory `assets` with the test execution, which will be the execution of the
`runAllTests.sh` bash script output: `coverage.html` and `final.out`. Both files describe the coverage.

# Debug Mode
The server can be started with the debug mode on, both for metrics or profiling. We just need to follow two steps:
* Set up DEBUG_SERVER configuration to the properties wanted:
```
SERVER_HOST=localhost
SERVER_PORT=8080
OPENAPI_RELATIVE_PATH=./api/keyvalue/api.yml
DEBUG_SERVER_HOST=localhost <<<<<
DEBUG_SERVER_PORT=8081 <<<<<
```
* When building the server, we should add a tag of `debug` to add the necessary dependencies (this process is done to 
avoid importing `pprof` package by default). Consequently, we should build it like so:
```bash
go build -o assets/keyvalueserver -tags debug
```

This will generate an executable with the capabilities described on 
[pprof documentation](https://pkg.go.dev/net/http/pprof). You can access all that information on the URL (f.e. http://localhost:8081/debug/pprof/).