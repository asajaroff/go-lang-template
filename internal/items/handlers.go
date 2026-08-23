// Package items is the sample CRUD resource — stands in for whatever your
// service actually manages. Swap this out; keep the pattern (migration →
// sqlc query → handler → test).
package items

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	itemsdb "github.com/asajaroff/go-lang-template/internal/items/db"
)

// RegisterRoutes wires the CRUD handlers into mux. q is itemsdb.Querier
// (not the concrete *itemsdb.Queries) so tests can substitute a fake.
func RegisterRoutes(mux *http.ServeMux, logger *slog.Logger, q itemsdb.Querier) {
	h := &handler{logger: logger, q: q}
	mux.HandleFunc("POST /items", h.create)
	mux.HandleFunc("GET /items", h.list)
	mux.HandleFunc("GET /items/{id}", h.get)
	mux.HandleFunc("PATCH /items/{id}", h.update)
	mux.HandleFunc("DELETE /items/{id}", h.delete)
}

type handler struct {
	logger *slog.Logger
	q      itemsdb.Querier
}

func (h *handler) create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	if err := req.validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	item, err := h.q.CreateItem(r.Context(), itemsdb.CreateItemParams{
		Name:        req.Name,
		Description: toPgText(req.Description),
	})
	if err != nil {
		h.logger.Error("create item failed", slog.Any("error", err))
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, toResponse(item))
}

func (h *handler) list(w http.ResponseWriter, r *http.Request) {
	dbItems, err := h.q.ListItems(r.Context())
	if err != nil {
		h.logger.Error("list items failed", slog.Any("error", err))
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	out := make([]response, 0, len(dbItems))
	for _, item := range dbItems {
		out = append(out, toResponse(item))
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *handler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	item, err := h.q.GetItem(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "item not found")
		return
	}
	if err != nil {
		h.logger.Error("get item failed", slog.Any("error", err))
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, toResponse(item))
}

func (h *handler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var req updateRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	if err := req.validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	item, err := h.q.UpdateItem(r.Context(), itemsdb.UpdateItemParams{
		ID:          id,
		Name:        toPgText(req.Name),
		Description: toPgText(req.Description),
	})
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "item not found")
		return
	}
	if err != nil {
		h.logger.Error("update item failed", slog.Any("error", err))
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, toResponse(item))
}

func (h *handler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	n, err := h.q.DeleteItem(r.Context(), id)
	if err != nil {
		h.logger.Error("delete item failed", slog.Any("error", err))
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if n == 0 {
		writeError(w, http.StatusNotFound, "item not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return uuid.UUID{}, false
	}
	return id, true
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

type errorBody struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorBody{Error: message})
}
