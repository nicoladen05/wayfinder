# wayfinder

A recursive DNS resolver written from scratch as a learning project. I am building it to learn Go and to learn how DNS servers work.

## Usage

```sh
go run ./cmd/wayfinder
go test ./...
```

## Layout

- `cmd/wayfinder`: entry point
- `internal/dnsmsg`: DNS message encoding and decoding (RFC 1035)
- `internal/resolver`: recursive resolution starting at the root servers
- `internal/server`: server that passes incoming questions to the resolver
