package items_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asajaroff/go-lang-template/internal/items"
)

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()
	logger := slog.New(slog.DiscardHandler)
	items.RegisterRoutes(mux, logger, newFakeQuerier())
	return httptest.NewServer(mux)
}

func postJSON(t *testing.T, srv *httptest.Server, path string, body any) *http.Response {
	t.Helper()
	buf, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}
	resp, err := srv.Client().Post(srv.URL+path, "application/json", bytes.NewReader(buf))
	if err != nil {
		t.Fatalf("POST %s: %v", path, err)
	}
	return resp
}

func decode(t *testing.T, resp *http.Response, dst any) {
	t.Helper()
	defer func() { _ = resp.Body.Close() }()
	if err := json.NewDecoder(resp.Body).Decode(dst); err != nil {
		t.Fatalf("decode response: %v", err)
	}
}

func TestCreateAndGetItem(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp := postJSON(t, srv, "/items", map[string]string{"name": "widget"})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create: got status %d", resp.StatusCode)
	}
	var created map[string]any
	decode(t, resp, &created)

	getResp, err := srv.Client().Get(srv.URL + "/items/" + created["id"].(string))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("get: got status %d", getResp.StatusCode)
	}
}

func TestCreateRejectsBlankName(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp := postJSON(t, srv, "/items", map[string]string{"name": ""})
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("got status %d, want %d", resp.StatusCode, http.StatusUnprocessableEntity)
	}
}

func TestListItems(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	postJSON(t, srv, "/items", map[string]string{"name": "a"})
	postJSON(t, srv, "/items", map[string]string{"name": "b"})

	resp, err := srv.Client().Get(srv.URL + "/items")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	var out []map[string]any
	decode(t, resp, &out)
	if len(out) != 2 {
		t.Fatalf("got %d items, want 2", len(out))
	}
}

func TestUpdateItem(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	created := postJSON(t, srv, "/items", map[string]string{"name": "old"})
	var item map[string]any
	decode(t, created, &item)

	req, _ := http.NewRequest(
		http.MethodPatch,
		srv.URL+"/items/"+item["id"].(string),
		bytes.NewReader([]byte(`{"name":"new"}`)),
	)
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	var updated map[string]any
	decode(t, resp, &updated)
	if updated["name"] != "new" {
		t.Fatalf("got name %v, want new", updated["name"])
	}
}

func TestDeleteItem(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	created := postJSON(t, srv, "/items", map[string]string{"name": "temp"})
	var item map[string]any
	decode(t, created, &item)
	id := item["id"].(string)

	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/items/"+id, nil)
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("delete: got status %d", resp.StatusCode)
	}

	getResp, _ := srv.Client().Get(srv.URL + "/items/" + id)
	if getResp.StatusCode != http.StatusNotFound {
		t.Fatalf("get after delete: got status %d, want 404", getResp.StatusCode)
	}
}

func TestGetMissingItem404(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp, _ := srv.Client().Get(srv.URL + "/items/00000000-0000-0000-0000-000000000000")
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("got status %d, want 404", resp.StatusCode)
	}
}
