package auth

type RegisterModel struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Mail      string `json:"mail"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birthDate"`
}

type LoginModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
