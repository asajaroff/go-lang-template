# go-lang-template

A batteries-included template for a [12-factor](https://12factor.net/) Go service: stdlib
`net/http` (Go 1.22+ pattern routing), `sqlc` + `pgx` for a Postgres-backed sample resource,
`golang-migrate` migrations, a distroless multi-stage `Dockerfile`, `docker-compose.yml` for
local dev, and a sample Helm chart for deploying to Kubernetes.

It also carries this author's usual repo tooling: [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/#summary)
+ [Git Flow](https://about.gitlab.com/blog/what-is-gitflow/#what-is-gitflow), enforced by
`pre-commit` hooks, plus a `git-cliff`-generated `CHANGELOG`.

## Using this template

Click **Use this template** on GitHub, clone your new repo, then rename the project. Set the
three variables, then run the rest from the repo root:

```bash
NEW_MODULE=github.com/acme/my-service   # your Go module path
NEW_NAME=my-service                      # used for the chart, image, app name and DB name

git rm -q TODO.md                        # template status notes, not yours
git grep -lz -e go-lang-template -e go_lang_template -- . ':!go.sum' | xargs -0 sed -i \
  -e "s#github.com/asajaroff/go-lang-template#$NEW_MODULE#g" \
  -e "s#go-lang-template#$NEW_NAME#g" \
  -e "s#go_lang_template#${NEW_NAME//-/_}#g"
git mv charts/go-lang-template charts/$NEW_NAME
go mod tidy && make test
```

Then replace the sample `items` resource with your own (see `CLAUDE.md`), update the
`LICENSE`, and run `pre-commit install`.

## Quickstart

Requires Go 1.26 or newer (see `go.mod`).

```bash
go mod download                 # install dependencies (Factor II — Go modules are the
                                 # native equivalent of a package manager here, no extra tool needed)
cp .env.example .env            # local-only config (Factor III) — never commit this file
docker compose up -d db         # a throwaway local Postgres (Factor IV)
go run ./cmd/server migrate up  # run migrations as a one-off process (Factor XII)
go run ./cmd/server serve
```

Or run the whole stack (app + Postgres) in containers:

```bash
docker compose up --build
```

`GET /healthz/liveness` is the liveness probe, `GET /healthz/readiness` checks the DB
connection, and `GET /metrics` serves Prometheus/VictoriaMetrics-format metrics.

A `Makefile` wraps the common commands: `make build`, `make run`, `make test`, `make lint`,
`make fmt`, `make migrate-up`, `make compose-up`, `make helm-lint`, etc. Run
`make install-tools` once to get `sqlc` and `golangci-lint` on your `PATH`.

## Tests

```bash
go test ./...
```

The `items` handler tests run against an in-memory fake implementing the same `Querier`
interface `sqlc` generates — no external services required, the same dependency-substitution
idea the Python sibling template uses for its DB session.

## Database access: sqlc

`internal/items/queries.sql` holds plain SQL; `sqlc generate` (or `make sqlc`) turns it into
type-safe Go in `internal/items/db/` (committed — contributors don't need `sqlc` installed to
build, only to change a query). Add a new query by editing `queries.sql` and regenerating;
add a new table via a `golang-migrate` migration in `migrations/`.

## Docker

```bash
docker build -t go-lang-template .
docker run --rm -p 8000:8000 --env-file .env go-lang-template
```

The runtime image is `distroless/static-debian12:nonroot` — no shell, no package manager,
runs as an unprivileged user by default. The compiled binary reads `$PORT` itself (Factor
VII), so there's no shell-form `CMD` needed to expand it the way the Python sibling template
needs.

## Helm

```bash
helm lint charts/go-lang-template
helm template go-lang-template charts/go-lang-template
helm install go-lang-template charts/go-lang-template \
  --set image.repository=<your-registry>/go-lang-template \
  --set image.tag=<tag>
```

The chart expects a Secret carrying `DATABASE_URL` to already exist in the target namespace
(see `values.yaml`'s `env.existingSecret`) — it does not template secret values or bundle a
database itself, since backing services are attached resources, not chart dependencies. See
`charts/go-lang-template/templates/NOTES.txt` (also printed after `helm install`) for how to
create that Secret. A pre-install/pre-upgrade hook Job runs `migrate up` (the same binary,
same image, different subcommand) before each release rolls out.

## How this maps to the 12 factors

| # | Factor | Where |
|---|---|---|
| I | Codebase | This repo |
| II | Dependencies | `go.mod`/`go.sum`, Go's own module system — no extra package manager; Dockerfile builds from the same locked manifest, no system packages |
| III | Config | `internal/config` (Viper, env-only) reads everything from env vars; `.env.example` documents them; `.env` is gitignored |
| IV | Backing services | Postgres reached only via `DATABASE_URL`; swappable without code changes — `docker-compose.yml` supplies one locally, the Helm chart expects one injected via Secret |
| V | Build, release, run | `Dockerfile` builds one static binary/image; CI (`.github/workflows/ci.yml`) builds/tests it; Helm/K8s runs it |
| VI | Processes | The app is stateless — the DB pool is the only long-lived resource, held outside any single request |
| VII | Port binding | The binary reads `$PORT` itself and binds it directly — see `internal/config` and `cmd/server/serve.go` |
| VIII | Concurrency | Scale via replica count / the chart's `HorizontalPodAutoscaler`, not in-process goroutine tuning |
| IX | Disposability | `cmd/server/serve.go` opens the DB pool fast on startup and shuts down gracefully on SIGINT/SIGTERM; `/healthz/readiness` gates traffic until the DB is reachable |
| X | Dev/prod parity | `docker-compose.yml` and the Helm chart run the exact same image |
| XI | Logs | `internal/logging` writes structured JSON to stdout via `log/slog` — no files, no rotation |
| XII | Admin processes | `go run ./cmd/server migrate up` locally; a Helm hook Job in the cluster (`templates/migration-job.yaml`), same binary |

## Links & Resources

- [12-factor app](https://12factor.net/)
- [Go](https://go.dev/)
- [sqlc](https://sqlc.dev/)
- [pgx](https://github.com/jackc/pgx)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [Helm](https://helm.sh/)
- [git hooks](https://git-scm.com/docs/githooks)
- [pre-commit](https://pre-commit.com/)
- [Keep a CHANGELOG](https://keepachangelog.com/en/1.0.0/)
- [Git Flow Cheatsheet](https://danielkummer.github.io/git-flow-cheatsheet/)
- [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/#summary)
