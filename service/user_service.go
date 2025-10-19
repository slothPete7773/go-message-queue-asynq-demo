package service

import (
	"asynq-demo/repository"
	"errors"
	"regexp"
)

// UserService defines the interface for user business logic
type UserService interface {
	CreateUser(name, email string, age int) (*repository.User, error)
	GetUserByID(id int) (*repository.User, error)
	GetAllUsers() ([]*repository.User, error)
	UpdateUser(id int, name, email string, age int) (*repository.User, error)
	DeleteUser(id int) error
}

// userService implements UserService
type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

// CreateUser creates a new user with validation
func (s *userService) CreateUser(name, email string, age int) (*repository.User, error) {
	// Business logic validation
	if err := s.validateUserInput(name, email, age); err != nil {
		return nil, err
	}

	user := &repository.User{
		Name:  name,
		Email: email,
		Age:   age,
	}

	return s.repo.Create(user)
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(id int) (*repository.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	return s.repo.GetByID(id)
}

// GetAllUsers retrieves all users
func (s *userService) GetAllUsers() ([]*repository.User, error) {
	return s.repo.GetAll()
}

// UpdateUser updates an existing user with validation
func (s *userService) UpdateUser(id int, name, email string, age int) (*repository.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	// Check if user exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Business logic validation
	if err := s.validateUserInput(name, email, age); err != nil {
		return nil, err
	}

	user := &repository.User{
		ID:    id,
		Name:  name,
		Email: email,
		Age:   age,
	}

	return s.repo.Update(user)
}

// DeleteUser deletes a user by ID
func (s *userService) DeleteUser(id int) error {
	if id <= 0 {
		return errors.New("invalid user ID")
	}

	return s.repo.Delete(id)
}

// validateUserInput validates user input according to business rules
func (s *userService) validateUserInput(name, email string, age int) error {
	if name == "" {
		return errors.New("name is required")
	}

	if len(name) < 2 {
		return errors.New("name must be at least 2 characters long")
	}

	if len(name) > 100 {
		return errors.New("name must be less than 100 characters")
	}

	if email == "" {
		return errors.New("email is required")
	}

	if !isValidEmail(email) {
		return errors.New("invalid email format")
	}

	if age < 0 {
		return errors.New("age cannot be negative")
	}

	if age > 150 {
		return errors.New("age must be less than 150")
	}

	return nil
}

// isValidEmail validates email format
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
