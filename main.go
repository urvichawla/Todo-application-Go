package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Todo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var todos []Todo

func main() {
	router := mux.NewRouter()

	// API Routes
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/api/todos", getTodos).Methods("GET")
	router.HandleFunc("/api/todos", createTodo).Methods("POST")
	router.HandleFunc("/api/todos/{id}", updateTodo).Methods("PUT")
	router.HandleFunc("/api/todos/{id}", deleteTodo).Methods("DELETE")

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Origin"},
		AllowCredentials: true,
	})

	// Create server
	srv := &http.Server{
		Handler:      c.Handler(router),
		Addr:         ":8081",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Server is running on http://localhost:8081")
	log.Fatal(srv.ListenAndServe())
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Welcome to Todo API",
		"endpoints": []string{
			"GET    /api/todos",
			"POST   /api/todos",
			"PUT    /api/todos/:id",
			"DELETE /api/todos/:id",
		},
	}
	json.NewEncoder(w).Encode(response)
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newTodo Todo
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate todo data
	if newTodo.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Generate ID if not provided
	if newTodo.ID == "" {
		newTodo.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	todos = append(todos, newTodo)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get ID from URL
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Todo ID is required", http.StatusBadRequest)
		return
	}

	var updatedTodo Todo
	if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate todo data
	if updatedTodo.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Ensure the ID in the URL matches the ID in the body
	updatedTodo.ID = id

	// Find and update the todo
	found := false
	for i, todo := range todos {
		if todo.ID == id {
			todos[i] = updatedTodo
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(updatedTodo)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get ID from URL
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Todo ID is required", http.StatusBadRequest)
		return
	}

	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			json.NewEncoder(w).Encode(map[string]string{"message": "Todo deleted successfully"})
			return
		}
	}

	http.Error(w, "Todo not found", http.StatusNotFound)
} 