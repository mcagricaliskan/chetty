package auth

import (
	"fmt"

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

func getGenderId(gender string) (int, error) {
	genderId, ok := SEX[gender]
	if !ok {
		return 0, fmt.Errorf("There is no gender like %s", gender)
	}

	return genderId, nil
}

func getCharacterGenderId(characterGender string) (int, error) {
	characterGenderId, ok := SEX[characterGender]
	if !ok {
		return 0, fmt.Errorf("There is no gender like %s", characterGender)
	}

	if characterGenderId == 3 {
		return 0, fmt.Errorf("Character gender can't be other")
	}

	return characterGenderId, nil
}
