package auth

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetUserPassword(c *fiber.Ctx, username string) (bool, string, error)
	IsUserExists(c *fiber.Ctx, username string) (bool, error)
}

type AuthenticationRepository struct {
	Connection *pgxpool.Pool
}

func (a AuthenticationRepository) Register() {
	return
}

func (a AuthenticationRepository) IsUserExists(ctx context.Context, username string) (isUserExists bool, err error) {
	err = a.Connection.QueryRow(ctx, `select when case count(*) > 0 than false else true from thegame.users where username = $1`, username).Scan(&isUserExists)
	return isUserExists, err
}

func (a AuthenticationRepository) GetUserPassword(ctx context.Context, username string) (isUserExists bool, password string, err error) {
	err = a.Connection.QueryRow(ctx,
		`select when case count(*) > 0 than false else true from thegame.users where username = $1`, username).Scan(&isUserExists, &password)
	return isUserExists, password, err

}
