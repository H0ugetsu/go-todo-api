package todo

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/h0ugetsu/todo-api/internal/utils"
)

const maxRequestBodySize = 1 << 20 // 1MB

type Handler struct {
	service Service
	logger  *slog.Logger
}

func NewHandler(service Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.ListTodos(r.Context())
	if err != nil {
		h.logger.Error("failed to list todos", "error", err)
		utils.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.WriteJSON(w, http.StatusOK, todos)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Description string `json:"description"`
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	todo, err := h.service.CreateTodo(r.Context(), req.Description)
	if errors.Is(err, ErrInvalidDescription) {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		h.logger.Error("failed to create todo", "error", err)
		utils.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.WriteJSON(w, http.StatusCreated, todo)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "id must be an integer")
		return
	}

	todo, err := h.service.GetTodo(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		h.logger.Error("failed to get todo", "error", err)
		utils.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.WriteJSON(w, http.StatusOK, todo)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Description *string `json:"description"`
		Completed   *bool   `json:"completed"`
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "id must be an integer")
		return
	}

	todo, err := h.service.UpdateTodo(r.Context(), id, req.Description, req.Completed)
	if errors.Is(err, ErrNotFound) {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	if errors.Is(err, ErrInvalidDescription) {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		h.logger.Error("failed to update todo", "error", err)
		utils.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.WriteJSON(w, http.StatusOK, todo)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "id must be an integer")
		return
	}

	err = h.service.DeleteTodo(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		h.logger.Error("failed to delete todo", "error", err)
		utils.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
