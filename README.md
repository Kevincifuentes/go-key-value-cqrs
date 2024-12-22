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

When finished, you can access to ```root/assets/coverage.html``` to see a user-friendly coverage report for 
all the modules.