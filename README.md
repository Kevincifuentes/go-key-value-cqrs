# go-key-value-cqrs

# Description

go-key-value-cqrs is an example project to implement CQRS, DDD and Hexagonal Architecture using Golang programming language.

It will support a Key/Value storage service with REST API capabilities to:
* Get KeyValue by Key
* Post a new KeyValue
* Delete a KeyValue


# Run tests

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