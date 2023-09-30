package auth

import (
	"context"
	"log"
)

type AuthService interface {
	register(ctx context.Context, RegisterReq *RegisterReq) error
	login(ctx context.Context, LoginReq *LoginReq) (userId int, err error)
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
	isUserExists, err := a.repository.IsUserExists(ctx, RegisterReq.UserName, RegisterReq.EMail)
	if err != nil {
		log.Println("auth -> service -> register -> IsUserExists -> Error while checking user exists, ", err)
		return ErrInternalServer
	}
	if isUserExists {
		return ErrUserExists
	}

	hashedPassword, err := HashPassword(RegisterReq.Password)
	if err != nil {
		log.Println("auth -> service -> register -> HashPassword -> Error while hashing password, ", err)
		return ErrInternalServer
	}

	err = a.repository.CreateUser(ctx, RegisterReq.UserName, RegisterReq.DisplayName, RegisterReq.EMail, hashedPassword)
	if err != nil {
		log.Println("auth -> service -> register -> CreateUser -> Error while creating user, ", err)
		return ErrInternalServer
	}

	return nil
}

func (a authService) login(ctx context.Context, LoginReq *LoginReq) (userId int, err error) {

	isUserExists, userId, userPassword, err := a.repository.GetUser(ctx, LoginReq.Username)
	if err != nil {
		log.Println("auth -> service.go -> GetUser -> Error while getting user, ", err)
		return 0, ErrInternalServer
	}
	if !isUserExists {
		return 0, ErrUnauthoerized
	}

	if !CheckPasswordHash(userPassword, userPassword) {
		return 0, ErrUnauthoerized
	}

	return userId, nil
}
