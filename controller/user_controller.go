package controller

import (
	"asynq-demo/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// UserController handles HTTP requests for user operations
type UserController struct {
	service service.UserService
}

// NewUserController creates a new instance of UserController
func NewUserController(service service.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string `json:"message"`
}

// CreateUser handles POST /users
func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		c.respondError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.respondError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := c.service.CreateUser(req.Name, req.Email, req.Age)
	if err != nil {
		c.respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.respondJSON(w, user, http.StatusCreated)
}

// GetUser handles GET /users/{id}
func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		c.respondError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := c.extractIDFromPath(r.URL.Path)
	if err != nil {
		c.respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := c.service.GetUserByID(id)
	if err != nil {
		c.respondError(w, err.Error(), http.StatusNotFound)
		return
	}

	c.respondJSON(w, user, http.StatusOK)
}

// extractIDFromPath extracts the ID from the URL path
func (c *UserController) extractIDFromPath(path string) (int, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return 0, nil
	}

	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// respondJSON sends a JSON response
func (c *UserController) respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// respondError sends an error response
func (c *UserController) respondError(w http.ResponseWriter, message string, statusCode int) {
	c.respondJSON(w, ErrorResponse{Error: message}, statusCode)
}
