package auth

import (
	"context"
	"the-game-backend/storage/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AuthRepository interface {
	IsUserExists(ctx context.Context, username string) (isUserExists bool, err error)
	GetUserPassword(ctx context.Context, username string) (isUserExists bool, password string, err error)
	CreateUser(ctx context.Context, user *RegisterReq, hashedPassword string, genderId int, characterGenderId int) error
	GetUser(ctx context.Context, username string) (isUserExists bool, user User, err error)
}

type AuthRepo struct {
	database *postgres.Postgres
}

func (r *AuthRepo) IsUserExists(ctx context.Context, username string) (isUserExists bool, err error) {
	err = r.database.Connection.QueryRow(ctx, `select 
			case 
				when count(*) = 1 then true
				else false
			end
		from thegame.users u where u.username = $1`,
		username).Scan(&isUserExists)
	return isUserExists, err
}

func (r *AuthRepo) GetUserPassword(ctx context.Context, username string) (isUserExists bool, password string, err error) {
	err = r.database.Connection.QueryRow(ctx,
		`select when case count(*) > 0 than false else true from thegame.users where username = $1`, username).Scan(&isUserExists, &password)
	return isUserExists, password, err

}

func (r *AuthRepo) CreateUser(ctx context.Context, user *RegisterReq, hashedPassword string, genderId int, characterGenderId int) error {
	id := uuid.New()
	_, err := r.database.Connection.Exec(ctx, `
		insert into thegame.users 
		(user_uuid, username, password, email, birth_date, gender_id, character_gender_id, created_at) 
		values ($1, $2, $3, $4, $5, $6, $7, now())`,
		id.String(), user.Username, hashedPassword, user.EMail, user.BirthDate, genderId, characterGenderId)
	return err
}

func (r *AuthRepo) GetUser(ctx context.Context, username string) (isUserExists bool, user User, err error) {
	err = r.database.Connection.QueryRow(ctx, `
		select user_uuid, password from thegame.users where username = $1`, username).Scan(&user.Id, &user.Password)
	if err == pgx.ErrNoRows {
		return false, user, nil
	}
	isUserExists = true
	return isUserExists, user, err
}
