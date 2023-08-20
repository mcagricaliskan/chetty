package auth

import (
	"context"
	"the-game-backend/services/postgres"
)

func IsUserExists(db *postgres.Postgres, ctx context.Context, username string) (isUserExists bool, err error) {
	err = db.Connection.QueryRow(ctx, `select when case count(*) > 0 than false else true from thegame.users where username = $1`, username).Scan(&isUserExists)
	return isUserExists, err
}

func GetUserPassword(db *postgres.Postgres, ctx context.Context, username string) (isUserExists bool, password string, err error) {
	err = db.Connection.QueryRow(ctx,
		`select when case count(*) > 0 than false else true from thegame.users where username = $1`, username).Scan(&isUserExists, &password)
	return isUserExists, password, err

}
