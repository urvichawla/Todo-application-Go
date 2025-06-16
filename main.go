package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var todos []Todo

func main() {
	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Root route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Todo API",
			"endpoints": []string{
				"GET    /api/todos",
				"POST   /api/todos",
				"PUT    /api/todos/:id",
				"DELETE /api/todos/:id",
			},
		})
	})

	// Routes
	router.GET("/api/todos", getTodos)
	router.POST("/api/todos", createTodo)
	router.PUT("/api/todos/:id", updateTodo)
	router.DELETE("/api/todos/:id", deleteTodo)

	router.Run(":8080")
}

func getTodos(c *gin.Context) {
	c.JSON(http.StatusOK, todos)
}

func createTodo(c *gin.Context) {
	var newTodo Todo
	if err := c.BindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo data"})
		return
	}

	// Validate todo data
	if newTodo.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	// Generate ID if not provided
	if newTodo.ID == "" {
		newTodo.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	todos = append(todos, newTodo)
	c.JSON(http.StatusCreated, newTodo)
}

func updateTodo(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Todo ID is required"})
		return
	}

	var updatedTodo Todo
	if err := c.BindJSON(&updatedTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo data"})
		return
	}

	// Validate todo data
	if updatedTodo.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, updatedTodo)
}

func deleteTodo(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Todo ID is required"})
		return
	}

	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
} 