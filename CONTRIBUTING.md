# Contributing

Thanks for contributing to timestamp.

## Setup

1. Fork and clone the repository.
2. `go mod download`
3. `go test ./...`

## Pull requests

- Run `go fmt ./...` and `go vet ./...`.
- Preserve monotonic guarantees in all code paths.
- Update `docs/` when initialization or configuration changes.
- Avoid adding network calls on the hot path after init.

## Changelog

Note changes under `## [Unreleased]` in [CHANGELOG.md](CHANGELOG.md).
