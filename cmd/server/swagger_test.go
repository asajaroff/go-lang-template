package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asajaroff/go-lang-template/internal/health"
	httpSwagger "github.com/swaggo/http-swagger"
)

func TestSwaggerEndpoint(t *testing.T) {
	// Set up HTTP routes
	mux := http.NewServeMux()
	healthHandler := health.NewHandler([]health.Checker{})
	mux.HandleFunc("/healthz/liveness", healthHandler.Liveness)
	mux.HandleFunc("/healthz/readiness", healthHandler.Readiness)
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		checkContent   bool
	}{
		{
			name:           "Swagger UI index page",
			path:           "/swagger/index.html",
			expectedStatus: http.StatusOK,
			checkContent:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkContent {
				body := w.Body.String()
				if body == "" {
					t.Errorf("Expected non-empty response body")
				}
			}
		})
	}
}
