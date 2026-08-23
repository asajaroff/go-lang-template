# Contributing

Thanks for contributing! This project follows [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/#summary) and [Git Flow](https://about.gitlab.com/blog/what-is-gitflow/#what-is-gitflow).

## Setup

```bash
go mod download        # installs dependencies
make install-tools     # installs sqlc + golangci-lint
pre-commit install     # wires up the pre-commit and commit-msg hook stages
cp .env.example .env   # local-only config — never commit this file
```

## Before you commit

`pre-commit install` (above) wires up hooks that run automatically on every commit and check:

- JSON/YAML validity, trailing whitespace, end-of-file newlines, oversized files, missing shebangs, and accidental AWS credentials
- Go formatting (`gofmt`), `go vet`, and `go mod tidy` staying in sync, via [pre-commit-golang](https://github.com/dnephin/pre-commit-golang)
- Static analysis via [golangci-lint](https://golangci-lint.run/) (requires `make install-tools` once)
- Spelling, via [cspell](https://cspell.org/), against both changed files and your commit message
- Commit message format, via [conventional-pre-commit](https://github.com/compilerla/conventional-pre-commit), enforced with `--strict`

Run all hooks against the whole repo at any time with:

```bash
pre-commit run --all-files
```

## Running the app and tests

```bash
docker compose up -d db             # local Postgres
go run ./cmd/server migrate up      # apply migrations
go run ./cmd/server serve
go test ./...                       # unit tests — no external services needed
```

See `Makefile` for shortcuts (`make run`, `make test`, `make lint`, `make fmt`, `make migrate-up`).

## Adding or changing the DB schema

1. Add a migration: create `migrations/00000X_description.up.sql` / `.down.sql` by hand
   (there's no autogenerate step here — `golang-migrate` migrations are plain SQL you write).
2. Add or edit a query in `internal/items/queries.sql` (or a new resource's equivalent file).
3. Regenerate the type-safe Go code: `make sqlc` (wraps `sqlc generate`). Commit the
   regenerated files under `internal/items/db/` along with your migration and query changes.

## Commit messages

Use the Conventional Commits format: `type(scope): message`.

Allowed scopes are: `feat, fix, refactor, tests, build, ci, changelog, docs, perf`.

If cspell flags a legitimate word (a tool name, proper noun, or project-specific term), add it to the `words` list in `cspell.json` rather than rewording around it.

## Changelog

`CHANGELOG` is generated, not hand-written. Regenerate it from commit history with [git-cliff](https://git-cliff.org/):

```bash
docker run -t -v "$(pwd)":/app/ "orhunp/git-cliff:${TAG:-latest}"
```

See `cliff.toml` for how commits are grouped into sections.

## Pull requests

Please fill out the PR template, including the PR checklist and PR type. Make sure your commits pass the hooks above, and that CI (`.github/workflows/ci.yml`) is green, before opening a PR.
