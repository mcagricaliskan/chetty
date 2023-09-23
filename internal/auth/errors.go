package auth

import "errors"

var (
	ErrUnauthoerized  = errors.New("Unauthorized")
	ErrInternalServer = errors.New("Internal Server Error")
	ErrUserExists     = errors.New("User Already Exists")
	ErrBadRequest     = errors.New("Bad Request")
)
