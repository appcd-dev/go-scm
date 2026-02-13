# AGENTS.md

## Repo purpose
[![Go Doc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](http://godoc.org/github.com/drone/go-scm/scm)

## Tech stack
- Go

## Build / test / lint / run (best-effort)

### Bootstrap
- (not detected; check README/Makefile/package scripts)

### Build
- `go build ./...`

### Test
- `go test ./...`

### Lint / format
- `golangci-lint run  # if configured`

### Run
- (not detected; check README/Makefile/package scripts)

## Project layout (where to look)
- Start with `README.md` and `.github/workflows/`.
- Common conventions (if present): `src/`, `lib/`, `cmd/`, `internal/`, `packages/`.
- Look for build entrypoints: `Makefile`, `package.json`, `go.mod`, `pyproject.toml`.

## CI / workflows
- (none detected)

## Common gotchas
- Some commands above are **best-effort guesses** based on repository signals; prefer README/Makefile/package scripts when they disagree.
- If CI fails, check workflow logs for exact tool versions and required secrets.
