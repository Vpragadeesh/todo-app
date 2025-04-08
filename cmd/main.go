package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yourusername/todo-app/handlers"
	"github.com/yourusername/todo-app/storage"
)

func main() {
	// Initialize the in-memory store.
	store := storage.NewMemoryStore()

	// Create the todo handler.
	todoHandler := handlers.NewTodoHandler(store)

	// Create a new router.
	router := mux.NewRouter()

	// Define routes.
	router.HandleFunc("/todos", todoHandler.GetTodos).Methods("GET")
	router.HandleFunc("/todos", todoHandler.AddTodo).Methods("POST")
	router.HandleFunc("/todos/{id}", todoHandler.UpdateTodo).Methods("PUT")
	router.HandleFunc("/todos/{id}", todoHandler.DeleteTodo).Methods("DELETE")

	// Start the server.
	log.Println("Server is listening on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
