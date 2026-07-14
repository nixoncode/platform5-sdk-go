# Contributing

## Prerequisites

- Go 1.25+

## Setup

```sh
go mod download
go build ./...
```

## Development

```sh
go vet ./...
go test ./... -v
```

## Pull Request Process

1. Run `go vet ./... && go build ./... && go test ./... -count=1`
2. Update documentation if needed
3. Open a PR against `main`
