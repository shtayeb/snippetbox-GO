## Create local certificates
For Linux and Mac
```shell
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

## Stdout log to a file
Add log errs to a file.
```shell
go run ./cmd/web >>/tmp/info.log 2>>/tmp/error.log
```

## Tests
Run all tests at once.
```shell
go test ./...
```

Run a specific test.
```shell
go test -v -run="^TestPing$" ./cmd/web/
```

Skip a specific test
```shell
go test -v -skip="^TestHumanDate$" ./cmd/web/
```

Clear test cache
```shell
go clean -testcache
```

Terminate at first fail
```shell
go test -failfast ./cmd/web
```
