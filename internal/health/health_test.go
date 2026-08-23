package health_test

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asajaroff/go-lang-template/internal/health"
)

type fakePinger struct{ err error }

func (f fakePinger) Ping(context.Context) error { return f.err }

func TestLiveness(t *testing.T) {
	rec := httptest.NewRecorder()
	health.Liveness(rec, httptest.NewRequest(http.MethodGet, "/healthz/liveness", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestReadinessOK(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	h := health.Readiness(logger, fakePinger{})

	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/healthz/readiness", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestReadinessFailsOnUnreachableDependency(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	h := health.Readiness(logger, fakePinger{err: errors.New("connection refused")})

	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, "/healthz/readiness", nil))

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
}
