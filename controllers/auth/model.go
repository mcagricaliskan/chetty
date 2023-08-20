package auth

type RegisterReq struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Mail      string `json:"mail"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birthDate"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
