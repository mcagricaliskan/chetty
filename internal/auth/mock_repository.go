package auth

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

// IsUserExists returns true if the user exists.
func (r *MockRepository) IsUserExists(userName string, userEmail string) (isUserExists bool, err error) {
	for _, user := range r.Users {
		if user.UserName == userName || user.Email == userEmail {
			return true, nil
		}
	}
	return false, nil
}
