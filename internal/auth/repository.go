package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mcagricaliskan/chetty/storage/postgres"
)

type AuthDatabaseRepository interface {
	IsUserExists(ctx context.Context, username string) (isUserExists bool, err error)
	GetUserPassword(ctx context.Context, username string) (isUserExists bool, password string, err error)
	CreateUser(ctx context.Context, RegisterReq *RegisterReq, hashedPassword string) error
	GetUser(ctx context.Context, username string) (isUserExists bool, user User, err error)
}

type authDatabaseRepository struct {
	database *postgres.Postgres
}

func NewAuthDatabaseRepository(database *postgres.Postgres) AuthDatabaseRepository {
	return &authDatabaseRepository{
		database: database,
	}
}

func (r *authDatabaseRepository) IsUserExists(ctx context.Context, username string) (isUserExists bool, err error) {
	err = r.database.Connection.QueryRow(ctx, `select 
			case 
				when count(*) = 1 then true
				else false
			end
		from chetty.users u where u.username = $1`,
		username).Scan(&isUserExists)
	return isUserExists, err
}

func (r *authDatabaseRepository) GetUserPassword(ctx context.Context, username string) (isUserExists bool, password string, err error) {
	err = r.database.Connection.QueryRow(ctx,
		`select when case count(*) > 0 than false else true from chetty.users where username = $1`, username).Scan(&isUserExists, &password)
	return isUserExists, password, err

}

func (r *authDatabaseRepository) CreateUser(ctx context.Context, RegisterReq *RegisterReq, hashedPassword string) error {
	id := uuid.New()
	_, err := r.database.Connection.Exec(ctx, `
		insert into chetty.users 
		(username, password, email, created_at) 
		values ($1, $2, $3, $4, $5, $6, now())`,
		id.String(), RegisterReq.Username, hashedPassword, RegisterReq.EMail)
	return err
}

func (r *authDatabaseRepository) GetUser(ctx context.Context, username string) (isUserExists bool, user User, err error) {
	err = r.database.Connection.QueryRow(ctx, `
		select user_id, password from chetty.users where username = $1`, username).Scan(&user.Id, &user.Password)
	if err == pgx.ErrNoRows {
		return false, user, nil
	}
	isUserExists = true
	return isUserExists, user, err
}
