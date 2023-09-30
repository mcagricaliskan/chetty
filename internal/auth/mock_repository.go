package auth

type User struct {
	Id         int    `json:"user_id"`
	UserName   string `json:"user_name"`
	DiplayName string `json:"display_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

// MockRepository is a mock implementation of the Repository interface.
type MockRepository struct {
	// Users is a map of user IDs to users.
	Users  map[int]*User
	LastID int
}

// NewMockRepository returns a new instance of MockRepository.
func NewMockRepository() *MockRepository {
	return &MockRepository{
		Users: make(map[int]*User),
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

// CreateUser creates a new user.
func (r *MockRepository) CreateUser(userName string, displayName string, email string, hashedPassword string) error {
	r.LastID++
	user := &User{
		Id:         r.LastID,
		UserName:   userName,
		DiplayName: displayName,
		Email:      email,
		Password:   hashedPassword,
	}
	r.Users[user.Id] = user
	return nil
}
