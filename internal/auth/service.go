package auth

import (
	"context"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// TODO:  bir kere kullanıyoruz zaten burayı service içerisine koyabiliriz
// HashPassword hashes the provided string.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

type authService struct {
	repository AuthDatabaseRepository
}

func NewAuthService(repository AuthDatabaseRepository) *authService {
	return &authService{
		repository: repository,
	}
}

func (a *authService) register(ctx context.Context, RegisterReq *RegisterReq) error {

	err := validatePassword(RegisterReq.Password)
	if err != nil {
		return ErrInvalidPassword
	}
	// if !validateUserName(RegisterReq.UserName) {
	// 	return errInvalidUserName
	// }
	// if !validateDisplayName(RegisterReq.DisplayName) {
	// 	return errInvalidDisplayName
	// }

	// i can move here to redis if user nubmer grows
	isUserExists, err := a.repository.IsUserExists(ctx, RegisterReq.UserName, RegisterReq.EMail)
	if err != nil {
		log.Println("auth -> service -> register -> IsUserExists -> Error while checking user exists, ", err)
		// return fmt.Errorf("%T.CreateUser(): %w", ErrInternalServer)
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
		log.Println("auth -> service -> register -> CreateUser -> Error while creating user: ", err)
		return ErrInternalServer
	}

	return nil
}

func (a *authService) login(ctx context.Context, LoginReq *LoginReq) (userId int, err error) {

	isUserExists, userId, userPassword, err := a.repository.GetUser(ctx, LoginReq.Username)
	if err != nil {
		log.Println("auth -> service.go -> GetUser -> Error while getting user, ", err)
		return 0, ErrInternalServer
	}
	if !isUserExists {
		return 0, ErrUnauthorized
	}

	err = bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(LoginReq.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrUnauthorized
		}
		return 0, ErrInternalServer
	}

	return userId, nil
}

//
