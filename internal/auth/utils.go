package auth

import (
	"errors"
	"fmt"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func validatePassword(password string) error {
	const minPasswordLength = 8
	if len(password) < minPasswordLength {
		return fmt.Errorf("password must be at least %d characters long", minPasswordLength)
	}

	var (
		hasUppercase, hasLowercase, hasDigit, hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUppercase = true
		case unicode.IsLower(char):
			hasLowercase = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUppercase {
		return errors.New("password must have at least one uppercase letter")
	}
	if !hasLowercase {
		return errors.New("password must have at least one lowercase letter")
	}
	if !hasDigit {
		return errors.New("password must have at least one digit")
	}
	if !hasSpecial {
		return errors.New("password must have at least one special character")
	}

	return nil
}
