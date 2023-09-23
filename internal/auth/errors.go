package auth

import "errors"

var (
	ErrUnauthoerized  = errors.New("Unauthorized")
	ErrInternalServer = errors.New("Internal Server Error")
)
