# go-lang-template rebuild — status

Rebuilding this repo from scratch as a 12-factor Go service template, mirroring
`fastapi-sample-app`'s structure/tooling but in idiomatic Go. The old ad-hoc scaffold
(Cobra/Viper config, health checks, golang-migrate, slog logging, Swagger docs, optional
Postgres — output of an earlier, unrelated AI session with a mismatched PRD) was wiped per
go-ahead. The work is committed on `feat/refactor` and open as PR #1 against `master`
(includes the 2026-11-05 dependency bump to Go 1.26).

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
- Multi-stage `Dockerfile`: `golang:1.26-bookworm` builder → static binary (`CGO_ENABLED=0`)
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
  `.pre-commit-config.yaml` (generic hooks + conventional-commits kept, Python's
  ruff/uv hooks swapped for `dnephin/pre-commit-golang`'s `go-fmt`/`go-mod-tidy`
  plus a local `golangci-lint` hook).
- `.github/workflows/ci.yml` (go vet, golangci-lint-action, go test, docker build check).
- `README.md` (Quickstart, Tests, sqlc workflow, Docker, Helm, 12-factor mapping table,
  Links) and `docs/CONTRIBUTING.md`, both mirroring the Python template's structure.

**Verified (2026-11-05)**: `go build`, `go vet`, `go test ./...` pass; `golangci-lint run ./...`
reports 0 issues; `sqlc generate` leaves the tree clean. Helm lint/template passed earlier.

## Not done yet

1. **`CLAUDE.md`** — not written yet (repo purpose, per-package architecture, conventions).
2. **`pre-commit run --all-files`** — passes (2026-11-05); cspell removed, `go-vet` is now a local hook.
3. **Docker build/run not verified** — `docker info` fails on this host (stale kernel, needs a
   reboot). Then: `docker build`, compose db + `migrate up`, curl `/healthz/*`.
4. **Code-review findings** — a `/code-review` pass ran on 2026-11-05; see PR #1 / follow-ups.
5. **`pre-commit-golang`** stays at v0.5.1 (newer releases dropped hooks we use).
6. **Distroless runtime** is still `static-debian12`; check for a Debian 13 image.
7. `make install-tools` is needed in a fresh environment for `sqlc` and `golangci-lint`.

## Suggested next steps

1. Triage code-review findings.
2. Write `CLAUDE.md`.
3. Verify Docker once usable.
