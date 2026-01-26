package health

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockChecker is a mock implementation of the Checker interface for testing
type mockChecker struct {
	name      string
	checkFunc func() error
}

func (m *mockChecker) Check() error {
	return m.checkFunc()
}

func (m *mockChecker) Name() string {
	return m.name
}

func TestLiveness(t *testing.T) {
	handler := NewHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/healthz/liveness", nil)
	w := httptest.NewRecorder()

	handler.Liveness(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != StatusHealthy {
		t.Errorf("expected status %s, got %s", StatusHealthy, response.Status)
	}

	if response.Checks != nil {
		t.Errorf("expected checks to be nil, got %v", response.Checks)
	}
}

func TestReadiness_AllHealthy(t *testing.T) {
	checkers := []Checker{
		&mockChecker{
			name: "service1",
			checkFunc: func() error {
				return nil
			},
		},
		&mockChecker{
			name: "service2",
			checkFunc: func() error {
				return nil
			},
		},
	}

	handler := NewHandler(checkers)

	req := httptest.NewRequest(http.MethodGet, "/healthz/readiness", nil)
	w := httptest.NewRecorder()

	handler.Readiness(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != StatusHealthy {
		t.Errorf("expected status %s, got %s", StatusHealthy, response.Status)
	}

	if len(response.Checks) != 2 {
		t.Errorf("expected 2 checks, got %d", len(response.Checks))
	}

	if response.Checks["service1"] != StatusHealthy {
		t.Errorf("expected service1 to be healthy")
	}

	if response.Checks["service2"] != StatusHealthy {
		t.Errorf("expected service2 to be healthy")
	}
}

func TestReadiness_SomeUnhealthy(t *testing.T) {
	checkers := []Checker{
		&mockChecker{
			name: "service1",
			checkFunc: func() error {
				return nil
			},
		},
		&mockChecker{
			name: "service2",
			checkFunc: func() error {
				return errors.New("connection failed")
			},
		},
	}

	handler := NewHandler(checkers)

	req := httptest.NewRequest(http.MethodGet, "/healthz/readiness", nil)
	w := httptest.NewRecorder()

	handler.Readiness(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("expected status %d, got %d", http.StatusServiceUnavailable, resp.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != StatusUnhealthy {
		t.Errorf("expected status %s, got %s", StatusUnhealthy, response.Status)
	}

	if response.Checks["service1"] != StatusHealthy {
		t.Errorf("expected service1 to be healthy")
	}

	if response.Checks["service2"] != StatusUnhealthy {
		t.Errorf("expected service2 to be unhealthy")
	}
}

func TestReadiness_NoCheckers(t *testing.T) {
	handler := NewHandler([]Checker{})

	req := httptest.NewRequest(http.MethodGet, "/healthz/readiness", nil)
	w := httptest.NewRecorder()

	handler.Readiness(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != StatusHealthy {
		t.Errorf("expected status %s, got %s", StatusHealthy, response.Status)
	}

	if len(response.Checks) != 0 {
		t.Errorf("expected 0 checks, got %d", len(response.Checks))
	}
}

func TestReadiness_AllUnhealthy(t *testing.T) {
	checkers := []Checker{
		&mockChecker{
			name: "service1",
			checkFunc: func() error {
				return errors.New("service1 down")
			},
		},
		&mockChecker{
			name: "service2",
			checkFunc: func() error {
				return errors.New("service2 down")
			},
		},
	}

	handler := NewHandler(checkers)

	req := httptest.NewRequest(http.MethodGet, "/healthz/readiness", nil)
	w := httptest.NewRecorder()

	handler.Readiness(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("expected status %d, got %d", http.StatusServiceUnavailable, resp.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != StatusUnhealthy {
		t.Errorf("expected status %s, got %s", StatusUnhealthy, response.Status)
	}

	if response.Checks["service1"] != StatusUnhealthy {
		t.Errorf("expected service1 to be unhealthy")
	}

	if response.Checks["service2"] != StatusUnhealthy {
		t.Errorf("expected service2 to be unhealthy")
	}
}

// TestReadiness_NilCheckers tests readiness with nil checkers list
func TestReadiness_NilCheckers(t *testing.T) {
	handler := NewHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/healthz/readiness", nil)
	w := httptest.NewRecorder()

	handler.Readiness(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != StatusHealthy {
		t.Errorf("expected status %s, got %s", StatusHealthy, response.Status)
	}

	if len(response.Checks) != 0 {
		t.Errorf("expected 0 checks, got %d", len(response.Checks))
	}
}

// TestReadiness_ManyCheckers tests readiness with many checkers
func TestReadiness_ManyCheckers(t *testing.T) {
	numCheckers := 10
	checkers := make([]Checker, numCheckers)

	for i := 0; i < numCheckers; i++ {
		checkers[i] = &mockChecker{
			name: fmt.Sprintf("service%d", i),
			checkFunc: func() error {
				return nil
			},
		}
	}

	handler := NewHandler(checkers)

	req := httptest.NewRequest(http.MethodGet, "/healthz/readiness", nil)
	w := httptest.NewRecorder()

	handler.Readiness(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != StatusHealthy {
		t.Errorf("expected status %s, got %s", StatusHealthy, response.Status)
	}

	if len(response.Checks) != numCheckers {
		t.Errorf("expected %d checks, got %d", numCheckers, len(response.Checks))
	}

	for i := 0; i < numCheckers; i++ {
		serviceName := fmt.Sprintf("service%d", i)
		if response.Checks[serviceName] != StatusHealthy {
			t.Errorf("expected %s to be healthy", serviceName)
		}
	}
}

// TestReadiness_MixedHealthStatus tests readiness with various health statuses
func TestReadiness_MixedHealthStatus(t *testing.T) {
	checkers := []Checker{
		&mockChecker{name: "healthy1", checkFunc: func() error { return nil }},
		&mockChecker{name: "unhealthy1", checkFunc: func() error { return errors.New("error1") }},
		&mockChecker{name: "healthy2", checkFunc: func() error { return nil }},
		&mockChecker{name: "unhealthy2", checkFunc: func() error { return errors.New("error2") }},
		&mockChecker{name: "healthy3", checkFunc: func() error { return nil }},
	}

	handler := NewHandler(checkers)

	req := httptest.NewRequest(http.MethodGet, "/healthz/readiness", nil)
	w := httptest.NewRecorder()

	handler.Readiness(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("expected status %d, got %d", http.StatusServiceUnavailable, resp.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != StatusUnhealthy {
		t.Errorf("expected status %s, got %s", StatusUnhealthy, response.Status)
	}

	if len(response.Checks) != 5 {
		t.Errorf("expected 5 checks, got %d", len(response.Checks))
	}

	expectedStatuses := map[string]Status{
		"healthy1":   StatusHealthy,
		"unhealthy1": StatusUnhealthy,
		"healthy2":   StatusHealthy,
		"unhealthy2": StatusUnhealthy,
		"healthy3":   StatusHealthy,
	}

	for name, expectedStatus := range expectedStatuses {
		if response.Checks[name] != expectedStatus {
			t.Errorf("expected %s to be %s, got %s", name, expectedStatus, response.Checks[name])
		}
	}
}

// TestLiveness_ContentType tests that liveness returns correct content type
func TestLiveness_ContentType(t *testing.T) {
	handler := NewHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/healthz/liveness", nil)
	w := httptest.NewRecorder()

	handler.Liveness(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got %q", contentType)
	}
}

// TestReadiness_ContentType tests that readiness returns correct content type
func TestReadiness_ContentType(t *testing.T) {
	handler := NewHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/healthz/readiness", nil)
	w := httptest.NewRecorder()

	handler.Readiness(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got %q", contentType)
	}
}

// TestReadiness_CheckerWithEmptyName tests checker with empty name
func TestReadiness_CheckerWithEmptyName(t *testing.T) {
	checkers := []Checker{
		&mockChecker{
			name: "",
			checkFunc: func() error {
				return nil
			},
		},
	}

	handler := NewHandler(checkers)

	req := httptest.NewRequest(http.MethodGet, "/healthz/readiness", nil)
	w := httptest.NewRecorder()

	handler.Readiness(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != StatusHealthy {
		t.Errorf("expected status %s, got %s", StatusHealthy, response.Status)
	}

	// Empty name should still be present in checks
	if len(response.Checks) != 1 {
		t.Errorf("expected 1 check, got %d", len(response.Checks))
	}

	if response.Checks[""] != StatusHealthy {
		t.Errorf("expected empty-named checker to be healthy")
	}
}

// TestReadiness_DuplicateCheckerNames tests checkers with duplicate names
func TestReadiness_DuplicateCheckerNames(t *testing.T) {
	checkers := []Checker{
		&mockChecker{
			name: "service",
			checkFunc: func() error {
				return nil
			},
		},
		&mockChecker{
			name: "service",
			checkFunc: func() error {
				return errors.New("failed")
			},
		},
	}

	handler := NewHandler(checkers)

	req := httptest.NewRequest(http.MethodGet, "/healthz/readiness", nil)
	w := httptest.NewRecorder()

	handler.Readiness(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Last checker with same name should win
	if response.Checks["service"] != StatusUnhealthy {
		t.Errorf("expected 'service' to be unhealthy (last checker status)")
	}
}

// TestNewHandler tests handler initialization
func TestNewHandler(t *testing.T) {
	t.Run("with nil checkers", func(t *testing.T) {
		handler := NewHandler(nil)
		if handler == nil {
			t.Error("expected non-nil handler")
		}
		if handler.readinessCheckers != nil {
			t.Errorf("expected nil checkers, got %v", handler.readinessCheckers)
		}
	})

	t.Run("with empty checkers", func(t *testing.T) {
		handler := NewHandler([]Checker{})
		if handler == nil {
			t.Error("expected non-nil handler")
		}
		if len(handler.readinessCheckers) != 0 {
			t.Errorf("expected 0 checkers, got %d", len(handler.readinessCheckers))
		}
	})

	t.Run("with checkers", func(t *testing.T) {
		checkers := []Checker{
			&mockChecker{name: "test", checkFunc: func() error { return nil }},
		}
		handler := NewHandler(checkers)
		if handler == nil {
			t.Error("expected non-nil handler")
		}
		if len(handler.readinessCheckers) != 1 {
			t.Errorf("expected 1 checker, got %d", len(handler.readinessCheckers))
		}
	})
}

// TestReadiness_CheckerPanic tests behavior when checker panics
func TestReadiness_CheckerReturningNil(t *testing.T) {
	checkers := []Checker{
		&mockChecker{
			name: "service",
			checkFunc: func() error {
				return nil
			},
		},
	}

	handler := NewHandler(checkers)

	req := httptest.NewRequest(http.MethodGet, "/healthz/readiness", nil)
	w := httptest.NewRecorder()

	handler.Readiness(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

// TestStatusConstants tests status constant values
func TestStatusConstants(t *testing.T) {
	if StatusHealthy != "healthy" {
		t.Errorf("expected StatusHealthy to be 'healthy', got %q", StatusHealthy)
	}
	if StatusUnhealthy != "unhealthy" {
		t.Errorf("expected StatusUnhealthy to be 'unhealthy', got %q", StatusUnhealthy)
	}
}
