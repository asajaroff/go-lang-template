// Package health implements liveness/readiness endpoints for container and
// Kubernetes probes.
//
// Factor IX (Disposability): Liveness never touches the database — a slow
// or unreachable backing service should not get the process killed.
// Readiness does check the database (Factor IV), so the Helm chart's
// readiness probe keeps a pod out of the Service's endpoints until it can
// actually serve requests.
package health

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

// Pinger is the one method Readiness needs from a backing service — small
// enough that any real dependency (a *pgxpool.Pool, a cache client, ...)
// satisfies it without an adapter, and easy to fake in tests.
type Pinger interface {
	Ping(ctx context.Context) error
}

type status struct {
	Status string `json:"status"`
}

// Liveness always reports ok — see the package doc for why it never
// touches the database.
func Liveness(w http.ResponseWriter, _ *http.Request) {
	writeStatus(w, http.StatusOK, "ok")
}

// Readiness reports ok only if every dependency's Ping succeeds.
func Readiness(logger *slog.Logger, deps ...Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, dep := range deps {
			if err := dep.Ping(r.Context()); err != nil {
				logger.Error("readiness check failed", slog.Any("error", err))
				writeStatus(w, http.StatusServiceUnavailable, "unavailable")
				return
			}
		}
		writeStatus(w, http.StatusOK, "ok")
	}
}

func writeStatus(w http.ResponseWriter, code int, s string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(status{Status: s})
}
