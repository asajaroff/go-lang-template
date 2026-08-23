# go-lang-template rebuild — status

Rebuilding this repo from scratch as a 12-factor Go service template, mirroring
`fastapi-sample-app`'s structure/tooling but in idiomatic Go. The old ad-hoc scaffold
(Cobra/Viper config, health checks, golang-migrate, slog logging, Swagger docs, optional
Postgres — output of an earlier, unrelated AI session with a mismatched PRD) was wiped per
go-ahead. Nothing has been committed yet — this is all uncommitted working-tree state on the
`master` worktree at `~/Code/github.com/asajaroff/go-lang-template.git/master`.

## Confirmed decisions (already implemented)

- **Router**: stdlib `net/http` (Go 1.22+ pattern routing, e.g. `"GET /items/{id}"`) — no
  external router dependency.
- **DB access**: `sqlc` (type-safe codegen from plain SQL) + `pgx/v5`. Generated code is
  committed under `internal/items/db/`.
- Config via Viper (env-only, no file/CLI-flag precedence layered on) — `internal/config`.
- Logging via stdlib `log/slog`, JSON to stdout — `internal/logging`.
- Metrics via `prometheus/client_golang` + `promhttp`, `/metrics` endpoint, request
  count/latency middleware excluding health/metrics paths from counted traffic —
  `internal/metrics`.
- Migrations via `golang-migrate`, plain SQL files in `migrations/`.
- Single compiled binary with two Cobra subcommands: `server serve` (HTTP server, graceful
  SIGINT/SIGTERM shutdown) and `server migrate up|down` (one-off admin process, Factor XII) —
  `cmd/server/`.
- Sample `items` CRUD resource (create/list/get/update/delete) — `internal/items/`
  (`queries.sql`, generated `db/`, `dto.go`, `handlers.go`), with unit tests against an
  in-memory fake `Querier` (no real Postgres needed for `go test ./...`).
- Multi-stage `Dockerfile`: `golang:1.25-bookworm` builder → static binary (`CGO_ENABLED=0`)
  → `gcr.io/distroless/static-debian12:nonroot` runtime. No shell in the runtime image, so no
  Docker `HEALTHCHECK` (documented tradeoff — Kubernetes probes are the real mechanism).
- `docker-compose.yml` (app + Postgres, `env_file: .env` + a `DATABASE_URL` override for the
  in-network hostname — applying the lesson learned from `fastapi-sample-app`'s code review
  rather than repeating that bug).
- Helm chart at `charts/go-lang-template/` mirroring `fastapi-sample-app`'s chart shape
  (Deployment/Service/ConfigMap/Secret/ServiceAccount/Ingress/HPA/migration-Job/NOTES.txt),
  **with the Python chart's known bugs already fixed here from the start**: `service.yaml`
  uses a named `targetPort: http` (not a separate drift-prone `values.yaml` number), and the
  migration Job carries the same `securityContext`/`podSecurityContext` as the Deployment.
- `.golangci.yml` (v2 config), `Makefile` (`install-tools`, `build`, `run`, `test`, `lint`,
  `fmt`, `sqlc`, `migrate-up/down`, `docker-*`, `compose-*`, `helm-*`).
- Baseline repo tooling copied from `git-repository-template.git` and adapted: `LICENSE`,
  `cliff.toml` (unchanged, generic), `.github/PULL_REQUEST_TEMPLATE.md` (unchanged),
  `.pre-commit-config.yaml` (generic hooks + cspell + conventional-commits kept, Python's
  ruff/uv hooks swapped for `dnephin/pre-commit-golang`'s `go-fmt`/`go-vet-mod`/`go-mod-tidy`
  plus a local `golangci-lint` hook), `cspell.json` (Go/sqlc/pgx/Helm vocabulary).
- `.github/workflows/ci.yml` (go vet, golangci-lint-action, go test, docker build check).
- `README.md` (Quickstart, Tests, sqlc workflow, Docker, Helm, 12-factor mapping table,
  Links) and `docs/CONTRIBUTING.md`, both mirroring the Python template's structure.

**Verified so far**: `go build ./...`, `go vet ./...`, `go test ./...` (all packages pass),
`golangci-lint run ./...` (0 issues), `helm lint`/`helm template` in default + ingress/
autoscaling/secret-enabled states (all render correctly).

## Not done yet

1. **`CLAUDE.md`** — was about to write this (mirroring `fastapi-sample-app`'s CLAUDE.md:
   repo purpose, architecture walkthrough per package, working-in-this-repo conventions) when
   interrupted. No file exists yet at the repo root.
2. **`.gitignore` cspell entry** — not re-checked against the new `.golangci.yml`/`sqlc.yaml`
   additions; run `pre-commit run --all-files` once `.git add -A`'d to see what cspell/other
   hooks flag (the Python template needed a couple of iterations here — expect the same).
3. **`pre-commit run --all-files` has not been run yet** on this repo at all — the
   `.pre-commit-config.yaml`/`cspell.json` are freshly written and unverified. This is the
   next concrete step once resuming.
4. **Docker build/run not verified** — same sandbox blocker as `fastapi-sample-app` (no
   `docker` group membership / dockerd was failing to start on this host; check if that's
   been resolved — see that repo's conversation history for the root cause: a stale running
   kernel with no matching `/usr/lib/modules/<version>` directory, fixed by a reboot).
5. **No commits yet** — `git status --short` currently shows the wipe (many `D` deletions of
   the old scaffold) plus all the new files as modified/untracked, all on the `master`
   worktree/branch. Haven't discussed branch naming for this repo (the Python template used
   `feat/v0.0.1` off `master`) — worth asking before committing, and whether to follow the
   same 6-commit-ish logical breakdown used for `fastapi-sample-app` (scaffold / docker /
   helm / ci / tooling / docs).
6. **`sqlc generate` and `golangci-lint`/`sqlc` binaries were installed into this sandbox's
   `$HOME/go/bin`** (via `go install`) to build and verify this — a fresh clone of this repo
   won't have them until someone runs `make install-tools`. Not a repo problem, just a note
   for continuing in a different environment.
7. Haven't double-checked whether `internal/items/db/` (sqlc-generated code) needs a header
   comment in `queries.sql`/`sqlc.yaml` documented in README (partially done) but not yet
   cross-referenced in CLAUDE.md's architecture section (blocked on item 1).
8. No equivalent of the Python template's `TODO-review.md` yet since no `/code-review` pass
   has been run against this new Go code — worth doing once the above is committed, mirroring
   what surfaced real bugs (Dockerfile/README coupling, Helm targetPort drift, migration Job
   securityContext) in the Python sibling.

## Suggested next steps, in order

1. Write `CLAUDE.md`.
2. `git add -A` and run `pre-commit run --all-files`; iterate on cspell words / hook config
   until clean (same pattern as the Python template's session).
3. `go test ./...` once more after any fixes from step 2 touch Go files.
4. Ask about branch name + commit breakdown, then commit (likely via the `committer` agent,
   same as before).
5. Once docker is usable again: `docker build`, `docker compose up -d db` + `migrate up` +
   `docker run`, `curl /healthz/liveness` and `/healthz/readiness`.
6. Consider a `/code-review` pass before or shortly after the first commit.
