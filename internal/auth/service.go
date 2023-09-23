package auth

import (
	"context"
	"fmt"
	"log"
)

type AuthService interface {
	register(ctx context.Context, RegisterReq *RegisterReq) error
	login(ctx context.Context, LoginReq *LoginReq) error
}

type authService struct {
	repository AuthDatabaseRepository
}

func NewAuthService(repository AuthDatabaseRepository) AuthService {
	return &authService{
		repository: repository,
	}
}

func (a authService) register(ctx context.Context, RegisterReq *RegisterReq) error {
	// i can move here to redis if user nubmer grows
	isUserExists, err := a.repository.IsUserExists(ctx, RegisterReq.Username)
	if err != nil {
		log.Println("auth -> Register -> IsUserExists -> Error while checking user exists, ", err)
		return fmt.Errorf("Internal Server Error")
	}
	if isUserExists {
		return fmt.Errorf("User Already Exists")
	}

	hashedPassword, err := HashPassword(RegisterReq.Password)
	if err != nil {
		log.Println("auth -> Register -> HashPassword -> Error while hashing password, ", err)
		return fmt.Errorf("Internal Server Error")
	}

	genderId, err := getGenderId(RegisterReq.Gender)
	if err != nil {
		log.Println("auth -> Register -> getGenderId -> Error while getting gender id, ", err)
		return fmt.Errorf("Bad Request")
	}

	characterGenderId, err := getCharacterGenderId(RegisterReq.CharacterGender)
	if err != nil {
		return fmt.Errorf("Bad Request")
	}

	err = a.repository.CreateUser(ctx, RegisterReq, hashedPassword, genderId, characterGenderId)
	if err != nil {
		log.Println("auth -> Register -> CreateUser -> Error while creating user, ", err)
		return fmt.Errorf("Internal Server Error")
	}

	return nil
}

func (a authService) login(ctx context.Context, LoginReq *LoginReq) error {
	// TODO
	return nil
}
