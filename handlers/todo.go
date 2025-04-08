package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yourusername/todo-app/models"
	"github.com/yourusername/todo-app/storage"
)

// TodoHandler holds a reference to our storage layer.
type TodoHandler struct {
	Store *storage.MemoryStore
}

// NewTodoHandler creates a new TodoHandler.
func NewTodoHandler(store *storage.MemoryStore) *TodoHandler {
	return &TodoHandler{Store: store}
}

// GetTodos handles GET /todos requests.
func (h *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	todos := h.Store.GetTodos()
	json.NewEncoder(w).Encode(todos)
}

// AddTodo handles POST /todos requests.
func (h *TodoHandler) AddTodo(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Title == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	todo := h.Store.AddTodo(input.Title)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

// UpdateTodo handles PUT /todos/{id} to update completion status.
func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Completed bool `json:"completed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.Store.UpdateTodo(id, input.Completed); err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"result": "updated"})
}

// DeleteTodo handles DELETE /todos/{id} requests.
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	if err := h.Store.DeleteTodo(id); err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"result": "deleted"})
}
