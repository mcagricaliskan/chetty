package auth

import (
	"context"
	"the-game-backend/services/postgres"

	"github.com/google/uuid"
)

func IsUserExists(db *postgres.Postgres, ctx context.Context, username string) (isUserExists bool, err error) {
	err = db.Connection.QueryRow(ctx, `select 
			case 
				when count(*) = 1 then true
				else false
			end
		from thegame.users u where u.username = $1`,
		username).Scan(&isUserExists)
	return isUserExists, err
}

func GetUserPassword(db *postgres.Postgres, ctx context.Context, username string) (isUserExists bool, password string, err error) {
	err = db.Connection.QueryRow(ctx,
		`select when case count(*) > 0 than false else true from thegame.users where username = $1`, username).Scan(&isUserExists, &password)
	return isUserExists, password, err

}

func CreateUser(db *postgres.Postgres, ctx context.Context, user *RegisterReq, hashedPassword string, genderId int, characterGenderId int) error {
	id := uuid.New()
	_, err := db.Connection.Exec(ctx, `
		insert into thegame.users 
		(user_uuid, username, password, email, birth_date, gender_id, character_gender_id, created_at) 
		values ($1, $2, $3, $4, $5, $6, $7, now())`,
		id.String(), user.Username, hashedPassword, user.EMail, user.BirthDate, genderId, characterGenderId)
	return err
}
