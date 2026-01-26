package health

import (
	"encoding/json"
	"errors"
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
