package main

import (
	"asynq-demo/controller"
	"asynq-demo/repository"
	"asynq-demo/service"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Initialize repository layer
	userRepo := repository.NewUserRepository()

	// Initialize service layer with repository dependency
	userService := service.NewUserService(userRepo)

	// Initialize controller layer with service dependency
	userController := controller.NewUserController(userService)

	// Setup HTTP routes
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			userController.CreateUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userController.GetUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"healthy"}`)
	})

	// Start server
	port := ":8088"
	fmt.Printf("User Management Service starting on port %s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET    /health          - Health check")
	fmt.Println("  GET    /users           - Get all users")
	fmt.Println("  POST   /users           - Create a new user")
	fmt.Println("  GET    /users/{id}      - Get user by ID")
	fmt.Println("  PUT    /users/{id}      - Update user by ID")
	fmt.Println("  DELETE /users/{id}      - Delete user by ID")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
