package server_test

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/h0ugetsu/todo-api/internal/server"
	"github.com/h0ugetsu/todo-api/internal/todo"
)

func newTestMux() http.Handler {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo := todo.NewRepository()
	service := todo.NewService(repo)
	handler := todo.NewHandler(service, logger)

	return server.NewRouter(handler)
}

func TestRouter_List(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		mux := newTestMux()

		req := httptest.NewRequest(http.MethodGet, "/todos", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}

		var got []todo.Todo
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("len(got) = %d, want 0", len(got))
		}
	})
}

func TestRouter_Create(t *testing.T) {
	t.Run("create_successfully", func(t *testing.T) {
		mux := newTestMux()

		body := strings.NewReader(`{"description": "牛乳を買う"}`)
		req := httptest.NewRequest(http.MethodPost, "/todos", body)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
		}

		var got todo.Todo
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}
		if got.Description != "牛乳を買う" {
			t.Errorf("Description = %s, want %s", got.Description, "牛乳を買う")
		}
	})

	t.Run("invalid_description", func(t *testing.T) {
		mux := newTestMux()

		body := strings.NewReader(`{"description": ""}`)
		req := httptest.NewRequest(http.MethodPost, "/todos", body)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})
}

func TestRouter_Get(t *testing.T) {
	t.Run("found_one", func(t *testing.T) {
		mux := newTestMux()

		createBody := strings.NewReader(`{"description": "牛乳を買う"}`)
		createReq := httptest.NewRequest(http.MethodPost, "/todos", createBody)
		createRec := httptest.NewRecorder()

		mux.ServeHTTP(createRec, createReq)

		var created todo.Todo
		if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/todos/"+strconv.Itoa(created.ID), nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}

		var got todo.Todo
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}
		if got.ID != created.ID {
			t.Errorf("ID = %d, want %d", got.ID, created.ID)
		}
		if got.Description != created.Description {
			t.Errorf("Description = %q, want %q", got.Description, created.Description)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		mux := newTestMux()

		req := httptest.NewRequest(http.MethodGet, "/todos/"+strconv.Itoa(999), nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
	})
}

func TestRouter_Update(t *testing.T) {
	t.Run("update_successfully", func(t *testing.T) {
		mux := newTestMux()

		createBody := strings.NewReader(`{"description": "更新前"}`)
		createReq := httptest.NewRequest(http.MethodPost, "/todos", createBody)
		createRec := httptest.NewRecorder()

		mux.ServeHTTP(createRec, createReq)

		var created todo.Todo
		if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}

		body := strings.NewReader(`{"description": "更新後"}`)
		req := httptest.NewRequest(http.MethodPatch, "/todos/"+strconv.Itoa(created.ID), body)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}

		var got todo.Todo
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}
		if got.Description != "更新後" {
			t.Errorf("Description = %s, want %s", got.Description, "更新後")
		}
		if got.ID != created.ID {
			t.Errorf("ID = %d, want %d", got.ID, created.ID)
		}
	})

	t.Run("invalid_description", func(t *testing.T) {
		mux := newTestMux()

		createBody := strings.NewReader(`{"description": "更新前"}`)
		createReq := httptest.NewRequest(http.MethodPost, "/todos", createBody)
		createRec := httptest.NewRecorder()

		mux.ServeHTTP(createRec, createReq)

		var created todo.Todo
		if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}

		body := strings.NewReader(`{"description": ""}`)
		req := httptest.NewRequest(http.MethodPatch, "/todos/"+strconv.Itoa(created.ID), body)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		mux := newTestMux()

		body := strings.NewReader(`{"description": "変更後"}`)
		req := httptest.NewRequest(http.MethodPatch, "/todos/"+strconv.Itoa(999), body)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
	})
}

func TestRouter_Delete(t *testing.T) {
	t.Run("delete_successfully", func(t *testing.T) {
		mux := newTestMux()

		createBody := strings.NewReader(`{"description": "牛乳を買う"}`)
		createReq := httptest.NewRequest(http.MethodPost, "/todos", createBody)
		createRec := httptest.NewRecorder()

		mux.ServeHTTP(createRec, createReq)

		var created todo.Todo
		if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}

		deleteReq := httptest.NewRequest(http.MethodDelete, "/todos/"+strconv.Itoa(created.ID), nil)
		deleteRec := httptest.NewRecorder()

		mux.ServeHTTP(deleteRec, deleteReq)

		if deleteRec.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d", deleteRec.Code, http.StatusNoContent)
		}

		getReq := httptest.NewRequest(http.MethodGet, "/todos/"+strconv.Itoa(created.ID), nil)
		getRec := httptest.NewRecorder()

		mux.ServeHTTP(getRec, getReq)

		if getRec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", getRec.Code, http.StatusNotFound)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		mux := newTestMux()

		req := httptest.NewRequest(http.MethodDelete, "/todos/"+strconv.Itoa(999), nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
	})
}
