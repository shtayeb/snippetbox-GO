## Create local certificates
For Linux and Mac
```shell
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

## Stdout log to a file
```shell
go run ./cmd/web >>/tmp/info.log 2>>/tmp/error.log
```
