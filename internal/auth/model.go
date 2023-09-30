package auth

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	EMail    string `json:"email"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Id       string `json:"user_id"`
	Password string `json:"password"`
}
