# Contributing

Thanks for contributing! This project follows [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/#summary) and [Git Flow](https://about.gitlab.com/blog/what-is-gitflow/#what-is-gitflow) (see Branching and releases below).

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
- Go formatting (`gofmt`) and `go mod tidy` staying in sync, via [pre-commit-golang](https://github.com/dnephin/pre-commit-golang)
- Static analysis via [golangci-lint](https://golangci-lint.run/) (requires `make install-tools` once)
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

Allowed types are: `feat, fix, refactor, test, build, ci, chore, docs, perf, revert, style`. The scope is free-form, e.g. `chore(deps): bump pgx`.

## Branching and releases

Git Flow: `develop` is the integration branch and `master` holds released code.

1. Branch from `develop` (`feat/...`, `fix/...`) and open a PR back to `develop`.
2. To release, branch `release/vX.Y.Z` from `develop`, run `make changelog TAG=vX.Y.Z`, and commit the updated `CHANGELOG.md`.
3. Merge the release branch into `master`, tag it `vX.Y.Z` and push the tag, then merge `master` back into `develop`.

Pushing a `v*` tag runs `.github/workflows/release.yml`: it builds and pushes the image to `ghcr.io/<owner>/<repo>` and creates a GitHub Release with notes from git-cliff.

## Changelog

`CHANGELOG.md` is generated, not hand-written. Regenerate it from commit history with [git-cliff](https://git-cliff.org/) (uses Docker):

```bash
make changelog TAG=v0.1.0
```

See `cliff.toml` for how commits are grouped into sections.

## Dependency updates

Dependabot (`.github/dependabot.yml`) opens weekly PRs for Go modules, GitHub Actions and the Dockerfile base images, titled `build(deps): ...` / `ci(deps): ...` so they pass the commit hook.

## Pull requests

Please fill out the PR template, including the PR checklist and PR type. Make sure your commits pass the hooks above, and that CI (`.github/workflows/ci.yml`) is green, before opening a PR.
