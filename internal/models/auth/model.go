package auth

type RegisterRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Mail      string `json:"mail"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birthDate"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
