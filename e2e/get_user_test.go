package e2e

import (
	"asynq-demo/controller"
	"asynq-demo/repository"
	"asynq-demo/service"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetUser_E2E tests the complete flow of getting a user
// This test covers the full stack: Controller -> Service -> Repository
func TestGetUser_E2E(t *testing.T) {
	// Setup: Initialize the full application stack
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	tests := []struct {
		name           string
		requestPath    string
		expectedStatus int
		expectedUser   *repository.User
		expectError    bool
	}{
		{
			name:           "Valid user ID returns user successfully",
			requestPath:    "/users/1",
			expectedStatus: http.StatusOK,
			expectedUser: &repository.User{
				ID:    1,
				Name:  "test",
				Email: "test@test.com",
				Age:   12,
			},
			expectError: false,
		},
		{
			name:           "Valid user ID with different number",
			requestPath:    "/users/123",
			expectedStatus: http.StatusOK,
			expectedUser: &repository.User{
				ID:    1,
				Name:  "test",
				Email: "test@test.com",
				Age:   12,
			},
			expectError: false,
		},
		{
			name:           "Invalid user ID format returns bad request",
			requestPath:    "/users/abc",
			expectedStatus: http.StatusBadRequest,
			expectedUser:   nil,
			expectError:    true,
		},
		{
			name:           "Missing user ID returns not found (ID=0 validation)",
			requestPath:    "/users/",
			expectedStatus: http.StatusNotFound,
			expectedUser:   nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new HTTP request
			req := httptest.NewRequest(http.MethodGet, tt.requestPath, nil)

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Execute the handler
			userController.GetUser(rr, req)

			// Check the status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			// Check Content-Type header
			if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf("handler returned wrong content type: got %v want %v",
					contentType, "application/json")
			}

			if tt.expectError {
				// For error cases, verify error response structure
				var errorResp controller.ErrorResponse
				if err := json.NewDecoder(rr.Body).Decode(&errorResp); err != nil {
					t.Errorf("failed to decode error response: %v", err)
				}
				if errorResp.Error == "" {
					t.Error("expected error message in response, got empty string")
				}
			} else {
				// For success cases, verify user data
				var user repository.User
				if err := json.NewDecoder(rr.Body).Decode(&user); err != nil {
					t.Errorf("failed to decode user response: %v", err)
				}

				// Verify user data matches expected values
				if user.ID != tt.expectedUser.ID {
					t.Errorf("handler returned unexpected user ID: got %v want %v",
						user.ID, tt.expectedUser.ID)
				}
				if user.Name != tt.expectedUser.Name {
					t.Errorf("handler returned unexpected user name: got %v want %v",
						user.Name, tt.expectedUser.Name)
				}
				if user.Email != tt.expectedUser.Email {
					t.Errorf("handler returned unexpected user email: got %v want %v",
						user.Email, tt.expectedUser.Email)
				}
				if user.Age != tt.expectedUser.Age {
					t.Errorf("handler returned unexpected user age: got %v want %v",
						user.Age, tt.expectedUser.Age)
				}
			}
		})
	}
}

// TestGetUser_E2E_WrongHTTPMethod tests that wrong HTTP methods are rejected
func TestGetUser_E2E_WrongHTTPMethod(t *testing.T) {
	// Setup
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	wrongMethods := []string{
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	}

	for _, method := range wrongMethods {
		t.Run("Method_"+method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/users/1", nil)
			rr := httptest.NewRecorder()

			userController.GetUser(rr, req)

			if status := rr.Code; status != http.StatusMethodNotAllowed {
				t.Errorf("handler returned wrong status code for method %s: got %v want %v",
					method, status, http.StatusMethodNotAllowed)
			}

			var errorResp controller.ErrorResponse
			if err := json.NewDecoder(rr.Body).Decode(&errorResp); err != nil {
				t.Errorf("failed to decode error response: %v", err)
			}
			if errorResp.Error != "method not allowed" {
				t.Errorf("handler returned unexpected error message: got %v want %v",
					errorResp.Error, "method not allowed")
			}
		})
	}
}

// TestGetUser_E2E_ConcurrentRequests tests that the handler can handle concurrent requests
func TestGetUser_E2E_ConcurrentRequests(t *testing.T) {
	// Setup
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	// Number of concurrent requests to send
	concurrentRequests := 10
	done := make(chan bool, concurrentRequests)

	for i := 1; i <= concurrentRequests; i++ {
		go func(userID int) {
			req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
			rr := httptest.NewRecorder()

			userController.GetUser(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("concurrent request failed with status: %v", status)
			}

			var user repository.User
			if err := json.NewDecoder(rr.Body).Decode(&user); err != nil {
				t.Errorf("failed to decode response: %v", err)
			}

			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < concurrentRequests; i++ {
		<-done
	}
}
