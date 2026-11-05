package metrics_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/asajaroff/go-lang-template/internal/metrics"
)

func TestMetricsEndpointExposesPrometheusFormat(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.Handle("GET /metrics", metrics.Handler())
	instrumented := metrics.Middleware(mux)

	srv := httptest.NewServer(instrumented)
	defer srv.Close()

	if _, err := srv.Client().Get(srv.URL + "/ping"); err != nil {
		t.Fatalf("GET /ping: %v", err)
	}

	resp, err := srv.Client().Get(srv.URL + "/metrics")
	if err != nil {
		t.Fatalf("GET /metrics: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("got status %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/plain") {
		t.Fatalf("got content-type %q, want text/plain prefix", resp.Header.Get("Content-Type"))
	}
}
