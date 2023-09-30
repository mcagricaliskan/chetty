package auth

import (
	"context"
	"errors"
)

// MockRepository is a mock implementation of the Repository interface.
type MockRepository struct {
	// Users is a map of user IDs to users.
	Users map[string]*User
}

// NewMockRepository returns a new instance of MockRepository.
func NewMockRepository() *MockRepository {
	return &MockRepository{
		Users: make(map[string]*User),
	}
}

// GetUserByID returns a user by ID.
func (r *MockRepository) GetUserByID(ctx context.Context, id string) (*User, error) {
	if user, ok := r.Users[id]; ok {
		return user, nil
	}

	return nil, errors.New("user not found")
}

//
