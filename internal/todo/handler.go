package todo

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/h0ugetsu/todo-api/internal/utils"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.ListTodos()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.WriteJSON(w, http.StatusOK, todos)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	todo, err := h.service.CreateTodo(req.Description)
	if errors.Is(err, ErrInvalidDescription) {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
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

	todo, err := h.service.GetTodo(id)
	if errors.Is(err, ErrNotFound) {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
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

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "id must be an integer")
		return
	}

	todo, err := h.service.UpdateTodo(id, req.Description, req.Completed)
	if errors.Is(err, ErrNotFound) {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	if errors.Is(err, ErrInvalidDescription) {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
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

	err = h.service.DeleteTodo(id)
	if errors.Is(err, ErrNotFound) {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
