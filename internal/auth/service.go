package auth

import (
	"context"
	"log"
)

type AuthService interface {
	register(ctx context.Context, RegisterReq *RegisterReq) error
	login(ctx context.Context, LoginReq *LoginReq) (user User, err error)
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
		return ErrInternalServer
	}
	if isUserExists {
		return ErrUserExists
	}

	hashedPassword, err := HashPassword(RegisterReq.Password)
	if err != nil {
		log.Println("auth -> Register -> HashPassword -> Error while hashing password, ", err)
		return ErrInternalServer
	}

	genderId, err := getGenderId(RegisterReq.Gender)
	if err != nil {
		log.Println("auth -> Register -> getGenderId -> Error while getting gender id, ", err)
		return ErrBadRequest
	}

	characterGenderId, err := getCharacterGenderId(RegisterReq.CharacterGender)
	if err != nil {
		return ErrBadRequest
	}

	err = a.repository.CreateUser(ctx, RegisterReq, hashedPassword, genderId, characterGenderId)
	if err != nil {
		log.Println("auth -> Register -> CreateUser -> Error while creating user, ", err)
		return ErrInternalServer
	}

	return nil
}

func (a authService) login(ctx context.Context, LoginReq *LoginReq) (user User, err error) {

	isUserExists, user, err := a.repository.GetUser(ctx, LoginReq.Username)
	if err != nil {
		log.Println("auth -> service.go -> GetUser -> Error while getting user, ", err)
		return user, ErrInternalServer
	}
	if !isUserExists {
		return user, ErrUnauthoerized
	}

	if !CheckPasswordHash(LoginReq.Password, user.Password) {
		return user, ErrUnauthoerized
	}

	return user, nil
}
