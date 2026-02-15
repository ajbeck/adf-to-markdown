# Contributing

## Prerequisites

- Go `1.25+`
- `GOEXPERIMENT=jsonv2` for normal development/test flows

## Local Development

Use the Makefile targets (they mirror CI):

```bash
make fmt-check
make vet
make build
make test
make test-nojsonv2
```

Run all CI-equivalent checks locally:

```bash
make ci
```

## Benchmarks and Fuzzing

Benchmarks:

```bash
make bench
```

Fuzzing:

```bash
make fuzz
```

Notes:
- `go test` uses built-in caching in package-list mode.
- Use `-count=1` when you need uncached test runs.
- Go manages build/test/fuzz caches automatically; use `go clean -cache`, `go clean -testcache`, or `go clean -fuzzcache` only when needed.

## Versioning

This repo uses semantic versioning with a checked-in `VERSION` file.

- Update `VERSION` for user-visible changes.
- Keep `CHANGELOG.md` updated for release notes.
- Conventional commit impact expected by CI:
  - `fix:` => patch bump
  - `feat:` => minor bump
  - `!` or `BREAKING CHANGE:` => major bump

## Pull Requests

PR workflow validates:

- formatting
- vet
- build
- tests (jsonv2)
- non-jsonv2 stub build test
- `VERSION` bump consistency with commits

Open PRs against `main`.

## Release

Release is performed through GitHub Actions `Release` workflow (`workflow_dispatch`):

1. Ensure `VERSION` and `CHANGELOG.md` are updated.
2. Ensure `make ci` passes locally.
3. Trigger the release workflow.

The release workflow:

- verifies `make vet` and `make test`
- reads `VERSION`
- creates/pushes tags `vX.Y.Z`, `vX.Y`, and `vX`
- creates a GitHub release
