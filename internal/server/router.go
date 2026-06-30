package server

import (
	"net/http"

	"github.com/h0ugetsu/todo-api/internal/respond"
	"github.com/h0ugetsu/todo-api/internal/todo"
)

func NewRouter(todoHandler *todo.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /todos", todoHandler.List)
	mux.HandleFunc("POST /todos", todoHandler.Create)
	mux.HandleFunc("GET /todos/{id}", todoHandler.Get)
	mux.HandleFunc("PATCH /todos/{id}", todoHandler.Update)
	mux.HandleFunc("DELETE /todos/{id}", todoHandler.Delete)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		respond.WriteJSON(w, http.StatusOK, map[string]string{
			"code":    http.StatusText(http.StatusOK),
			"message": "ok",
		})
	})

	return mux
}
