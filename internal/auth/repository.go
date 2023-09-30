package auth

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/mcagricaliskan/chetty/storage/postgres"
)

type AuthDatabaseRepository interface {
	IsUserExists(ctx context.Context, userName string, userEmail string) (isUserExists bool, err error)
	CreateUser(ctx context.Context, userName string, displayName string, email string, hashedPassword string) error
	GetUser(ctx context.Context, username string) (isUserExists bool, userId int, userPassword string, err error)
}

type authDatabaseRepository struct {
	database *postgres.Postgres
}

func NewAuthDatabaseRepository(database *postgres.Postgres) AuthDatabaseRepository {
	return &authDatabaseRepository{
		database: database,
	}
}

func (r *authDatabaseRepository) IsUserExists(ctx context.Context, userName string, userEmail string) (isUserExists bool, err error) {
	err = r.database.Connection.QueryRow(ctx, `select 
			case 
				when count(*) = 1 then true
				else false
			end
		from chetty.users u where u.user_name = $1 or u.email = $2`,
		userName, userEmail).Scan(&isUserExists)
	return isUserExists, err
}

func (r *authDatabaseRepository) CreateUser(ctx context.Context, userName string, displayName string, email string, hashedPassword string) error {
	_, err := r.database.Connection.Exec(ctx, `
		insert into chetty.users
		(user_name, display_name, password, email, created_at)
		values ($1, $2, $3, $4, now())`,
		userName, displayName, hashedPassword, email)
	return err
}

func (r *authDatabaseRepository) GetUser(ctx context.Context, username string) (isUserExists bool, userId int, userPassword string, err error) {
	err = r.database.Connection.QueryRow(ctx, `
		select user_id, display_name, password from chetty.users where username = $1`, username).Scan(&userId, &userPassword)
	if err == pgx.ErrNoRows {
		return false, userId, userPassword, nil
	}
	isUserExists = true
	return isUserExists, userId, userPassword, err
}
