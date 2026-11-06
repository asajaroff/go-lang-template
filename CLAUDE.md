# CLAUDE.md

Guidance for working in this repo.

## What this is

A 12-factor Go service template: stdlib `net/http` router, Postgres via `sqlc` + `pgx/v5`,
Cobra/Viper CLI, `slog` JSON logging, Prometheus metrics, a Dockerfile, docker-compose and a
Helm chart. The `items` package is a sample CRUD resource meant to be swapped out.
Requires Go 1.26+ (see `go.mod`). `README.md` has the 12-factor mapping table.

## Commands

```bash
make install-tools   # sqlc + golangci-lint (needed once per machine)
make test            # go test ./... (no Postgres needed)
make lint            # go vet + golangci-lint
make sqlc            # regenerate internal/items/db after editing queries or migrations
make run             # go run ./cmd/server serve
make migrate-up      # go run ./cmd/server migrate up (also: migrate-down)
make helm-lint       # lint the chart in charts/go-lang-template
```

Configuration is env-only; copy `.env.example` to `.env` for local runs.

## Layout

- `cmd/server/`: one binary, two Cobra subcommands. `serve` runs the HTTP server with
  graceful SIGINT/SIGTERM shutdown. `migrate up|down` is a one-off admin process (Factor XII).
  `serve` never runs migrations.
- `internal/config/`: Viper, env vars only.
- `internal/logging/`: `slog` JSON to stdout.
- `internal/database/`: pgx pool setup.
- `internal/health/`: `/healthz/liveness` (never touches the DB) and `/healthz/readiness`
  (pings dependencies through the small `Pinger` interface).
- `internal/metrics/`: `/metrics` plus request-count/latency middleware. Health and metrics
  paths are excluded from the counts.
- `internal/items/`: sample resource: `handlers.go`, `dto.go`, `queries.sql`, and the
  generated `db/`.
- `migrations/`: plain SQL for golang-migrate. These are also sqlc's schema source.
- `charts/go-lang-template/`: Helm chart, including a migration Job.

## Conventions

- Don't edit `internal/items/db/` by hand. Change `queries.sql` or `migrations/`, then
  `make sqlc`, and commit the generated output.
- New resources follow the `items` pattern: migration, sqlc query, handler, test.
- Handlers depend on `itemsdb.Querier` (an interface), so tests use the in-memory fake in
  `internal/items/fake_querier_test.go` rather than a real database.
- Routes use Go 1.22+ patterns (`"GET /items/{id}"`). Don't add an external router.
- Keep the Helm Service `targetPort` as the named port `http`, and keep the migration Job's
  `securityContext` the same as the Deployment's.
- The runtime image is distroless with no shell, so there is no Docker `HEALTHCHECK`;
  Kubernetes probes are the health mechanism.
- Git Flow: branch off `develop`, PR back to `develop`; releases are tagged on `master`.
- Commits: conventional commits with a scope (`chore(deps): ...`); the changelog comes from
  git-cliff (`make changelog TAG=vX.Y.Z`). `pre-commit install` also installs the strict `commit-msg` hook; any scope is
  allowed, the types are listed in `.pre-commit-config.yaml`.
- Run `make lint test` and `pre-commit run --all-files` before pushing.
