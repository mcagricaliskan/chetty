package auth

type RegisterReq struct {
	UserName    string `json:"user_name"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
	EMail       string `json:"email"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Id         string `json:"user_id"`
	UserName   string `json:"user_name"`
	DiplayName string `json:"display_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}
