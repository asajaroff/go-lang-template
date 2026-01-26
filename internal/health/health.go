package health

import (
	"encoding/json"
	"net/http"
)

// Status represents the health status of a component
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
)

// Response represents the health check response
type Response struct {
	Status Status            `json:"status"`
	Checks map[string]Status `json:"checks,omitempty"`
}

// Checker defines the interface for health check components
type Checker interface {
	Check() error
	Name() string
}

// Handler manages health check endpoints
type Handler struct {
	readinessCheckers []Checker
}

// NewHandler creates a new health check handler
func NewHandler(readinessCheckers []Checker) *Handler {
	return &Handler{
		readinessCheckers: readinessCheckers,
	}
}

// Liveness handles the liveness probe endpoint
// This endpoint always returns healthy if the server is running
func (h *Handler) Liveness(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status: StatusHealthy,
	}
	h.writeJSON(w, http.StatusOK, response)
}

// Readiness handles the readiness probe endpoint
// This endpoint checks if all backend dependencies are available
func (h *Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	checks := make(map[string]Status)
	overallStatus := StatusHealthy
	statusCode := http.StatusOK

	for _, checker := range h.readinessCheckers {
		if err := checker.Check(); err != nil {
			checks[checker.Name()] = StatusUnhealthy
			overallStatus = StatusUnhealthy
			statusCode = http.StatusServiceUnavailable
		} else {
			checks[checker.Name()] = StatusHealthy
		}
	}

	response := Response{
		Status: overallStatus,
		Checks: checks,
	}

	h.writeJSON(w, statusCode, response)
}

// writeJSON writes a JSON response
func (h *Handler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
