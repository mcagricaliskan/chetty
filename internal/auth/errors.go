package auth

import "errors"

var (
	ErrUnauthorized    = errors.New("unauthorized")
	ErrInternalServer  = errors.New("internal server error")
	ErrUserExists      = errors.New("user already exists")
	ErrBadRequest      = errors.New("bad request")
	ErrInvalidPassword = errors.New("invalid password")
)
