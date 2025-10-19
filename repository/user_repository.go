package repository

// User represents the user entity
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *User) (*User, error)
	GetByID(id int) (*User, error)
	GetAll() ([]*User, error)
	Update(user *User) (*User, error)
	Delete(id int) error
}

type userRepository struct {
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

var (
	mockUser = User{
		ID:    1,
		Name:  "test",
		Email: "test@test.com",
		Age:   12,
	}
)

func (r *userRepository) Create(user *User) (*User, error) {
	return &mockUser, nil
}

func (r *userRepository) GetByID(id int) (*User, error) {
	return &mockUser, nil
}

func (r *userRepository) GetAll() ([]*User, error) {
	return []*User{
		&mockUser,
	}, nil
}

func (r *userRepository) Update(user *User) (*User, error) {
	return &mockUser, nil
}

func (r *userRepository) Delete(id int) error {
	return nil
}
