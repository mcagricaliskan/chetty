package auth

type RegisterReq struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	EMail           string `json:"email"`
	Gender          string `json:"gender"`
	CharacterGender string `json:"character_gender"`
	BirthDate       string `json:"birth_date"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Id       string `json:"user_id"`
	Password string `json:"password"`
}
