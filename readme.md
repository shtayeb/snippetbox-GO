## Clone the repo 
```shell
git clone https://github.com/shtayeb/snippetbox-GO.git
```

## Setup the DB
Create a postgres database named `snippetbox`
```sql
CREATE DATABASE snippetbox OWNER go_user;
```

## Create local certificates
For Linux and Mac
```shell
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

## Start the server
Start in debug mode
```shell
go run ./cmd/web/ -debug
```

```shell
$ go run ./cmd/web --help

  -addr string
        HTTP network address (default ":4000")
  -debug
        Enable debug mode
  -dsn string
        Database source (default "postgres://go_user:go_1234@localhost/snippetbox")
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

## Resources
- [Let's GO](https://lets-go.alexedwards.net/) Book by Alex Edwards
