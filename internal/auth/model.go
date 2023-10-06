package auth

type RegisterReq struct {
	UserName    string `json:"user_name" validate:"required"`
	DisplayName string `json:"display_name" validate:"required"`
	Password    string `json:"password" validate:"required"`
	EMail       string `json:"email" validate:"required,email"`
}

// func (r *RegisterReq) Validate() error {
// 	// minimum 8 lenght, at least one uppercase, one lowercase, one number and one special character, without regex
// 	if len(r.Password) < 8 ||  {
// 		return ErrBadRequest
// 	}
// 	if r.UserName == "" || r.DisplayName == "" || r.EMail == "" {
// 		return ErrBadRequest
// 	}

// 	return nil
// }

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
