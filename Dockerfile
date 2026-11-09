# syntax=docker/dockerfile:1

# --- builder -------------------------------------------------------------
# Factor II (Dependencies): resolved and built from a locked, committed
# manifest (go.mod/go.sum) — nothing pulled from ambient system packages.
FROM golang:1.27-bookworm AS builder

WORKDIR /src

# Download deps first so this layer is cached independently of source
# changes (Factor V: reproducible builds should not re-resolve deps on
# every code edit).
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/go/pkg/mod go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY migrations/ ./migrations/

# CGO disabled + static build: a single self-contained binary that runs on
# a scratch-like distroless base with no libc dependency.
RUN --mount=type=cache,target=/root/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# --- runtime ---------------------------------------------------------------
# Factor X (Dev/prod parity): this exact image is what docker-compose runs
# locally and what the Helm chart deploys to every environment.
# distroless: no shell, no package manager, minimal attack surface; the
# nonroot variant already runs as an unprivileged user (uid 65532).
FROM gcr.io/distroless/static-debian12:nonroot AS runtime

WORKDIR /app
COPY --from=builder /out/server ./server
COPY --from=builder /src/migrations ./migrations

USER nonroot:nonroot

# Factor VII (Port binding): the binary reads $PORT itself (see
# internal/config) and binds it directly — no shell-form CMD needed to
# expand an env var, since distroless has no shell to run one in anyway.
EXPOSE 8000

# No Docker HEALTHCHECK here — distroless has no shell/curl to run one in,
# and Kubernetes' own liveness/readiness probes (see the Helm chart) are
# the actual mechanism this template relies on for that. docker-compose
# uses `depends_on: condition: service_started` for the same reason.
ENTRYPOINT ["/app/server"]
CMD ["serve"]
