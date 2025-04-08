package storage

import (
	"errors"
	"sync"
	"time"

	"github.com/yourusername/todo-app/models"
)

// MemoryStore holds our tasks in memory
type MemoryStore struct {
	todos []models.Todo
	mu    sync.Mutex
}

// NewMemoryStore creates and returns a new MemoryStore
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		todos: make([]models.Todo, 0),
	}
}

// AddTodo adds a new task to the store.
func (ms *MemoryStore) AddTodo(title string) models.Todo {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	// generate a new ID (for simplicity, using length + 1)
	newID := len(ms.todos) + 1
	todo := models.Todo{
		ID:        newID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}
	ms.todos = append(ms.todos, todo)
	return todo
}

// GetTodos returns all tasks.
func (ms *MemoryStore) GetTodos() []models.Todo {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	return ms.todos
}

// UpdateTodo updates a todo item; for instance, marking it as completed.
func (ms *MemoryStore) UpdateTodo(id int, completed bool) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	for i, t := range ms.todos {
		if t.ID == id {
			ms.todos[i].Completed = completed
			return nil
		}
	}
	return errors.New("todo not found")
}

// DeleteTodo deletes a todo by ID.
func (ms *MemoryStore) DeleteTodo(id int) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	for i, t := range ms.todos {
		if t.ID == id {
			ms.todos = append(ms.todos[:i], ms.todos[i+1:]...)
			return nil
		}
	}
	return errors.New("todo not found")
}
