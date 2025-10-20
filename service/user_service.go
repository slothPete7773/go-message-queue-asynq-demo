package service

import (
	messagequeue "asynq-demo/message-queue"
	asynqque "asynq-demo/message-queue/asynq-queue"
	"asynq-demo/repository"
	"errors"
	"fmt"
	"regexp"
)

// UserService defines the interface for user business logic
type UserService interface {
	CreateUser(name, email string, age int) (*repository.User, error)
	GetUserByID(id int) (*repository.User, error)
}

// userService implements UserService
type userService struct {
	repo      repository.UserRepository
	publisher messagequeue.Publisher
}

// NewUserService creates a new instance of UserService
func NewUserService(repo repository.UserRepository, publisher messagequeue.Publisher) UserService {
	return &userService{
		repo:      repo,
		publisher: publisher,
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

	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("error failed to find user: %w", err)
	}

	task := asynqque.NewWelcomeEmailTask(user.ID)
	err = s.publisher.Publish(task)

	if err != nil {
		return nil, fmt.Errorf("error failed to enqueue job: %w", err)
	}

	return user, nil
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
