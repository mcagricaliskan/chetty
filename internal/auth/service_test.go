package auth

import (
	"context"
	"log"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/mcagricaliskan/chetty/storage/postgres"
)

func Test_register(t *testing.T) {
	// define embedded-postgres

	embeddedPostgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(5555))
	if err := embeddedPostgres.Start(); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := embeddedPostgres.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	database, err := postgres.Connect("postgres://postgres:postgres@localhost:5555/postgres?sslmode=disable")
	if err != nil {
		log.Fatal("Database Connection Can't Estabilished, Error:", err)
	} else {
		log.Println("Database Connection Estabilished")
	}

	_, err = database.Connection.Exec(context.Background(), `create schema chetty;`)
	if err != nil {
		log.Fatal("Database Schema Can't Created, Error:", err)
	}
	_, err = database.Connection.Exec(context.Background(), `
	CREATE TABLE chetty.users
	(
	  user_id      bigint    NOT NULL GENERATED ALWAYS AS IDENTITY,
	  user_name    text      NOT NULL UNIQUE,
	  display_name text      NOT NULL,
	  password     text      NOT NULL,
	  email        text      NOT NULL UNIQUE,
	  created_at   timestamp NOT NULL,
	  PRIMARY KEY (user_id)
	);`)
	if err != nil {
		log.Fatal("Database Table Can't Created, Error:", err)
	}

	testTable := []struct {
		name string
		req  *RegisterReq
		resp error
	}{
		{
			name: "Success",
			req: &RegisterReq{
				UserName:    "mcagricaliskan",
				DisplayName: "Mehmet Çağrı Çalışkan",
				Password:    "mpyF{0q4THwk{p9dK!ff",
				EMail:       "mcagricaliskan@gmail.com",
			},
			resp: nil,
		},
		{
			name: "DuplicateUserName",
			req: &RegisterReq{
				UserName:    "mcagricaliskan",
				DisplayName: "Mehmet Çağrı Çalışkan",
				Password:    "mpyF{0q4THwk{p9dK!ff",
				EMail:       "mcagricaliskan_test@gmail.com",
			},
			resp: ErrUserExists,
		},
		{
			name: "DuplicateEmail",
			req: &RegisterReq{
				UserName:    "mcagricaliskan_test",
				DisplayName: "Mehmet Çağrı Çalışkan",
				Password:    "mpyF{0q4THwk{p9dK!ff",
				EMail:       "mcagricaliskan@gmail.com",
			},
			resp: ErrUserExists,
		},
		{
			name: "InvalidPasswordNotNumber",
			req: &RegisterReq{
				UserName:    "mcagricaliskan_test",
				DisplayName: "Mehmet Çağrı Çalışkan",
				Password:    "mpyF{qTHwk{pK!ff",
				EMail:       "mcagricaliskan_test@gmail.com",
			},
			resp: ErrInvalidPassword,
		},
		{
			name: "InvalidPasswordNotSpecialChar",
			req: &RegisterReq{
				UserName:    "mcagricaliskan_test",
				DisplayName: "Mehmet Çağrı Çalışkan",
				Password:    "mpyF0q4THwkp9dKff",
				EMail:       "mcagricaliskan_test@gmail.com",
			},
			resp: ErrInvalidPassword,
		},
		{
			name: "InvalidPasswordNotCapitalLetter",
			req: &RegisterReq{
				UserName:    "mcagricaliskan_test",
				DisplayName: "Mehmet Çağrı Çalışkan",
				Password:    "mpyf{0q4thwk{p9dk!ff",
				EMail:       "mcagricaliskan_test@gmail.com",
			},
			resp: ErrInvalidPassword,
		},
		{
			name: "InvalidPasswordNotSmallLetter",
			req: &RegisterReq{
				UserName:    "mcagricaliskan_test",
				DisplayName: "Mehmet Çağrı Çalışkan",
				Password:    "MPYF{0Q4THWK{P9DK!FF",
				EMail:       "mcagricaliskan_test@gmail.com",
			},
			resp: ErrInvalidPassword,
		},
		{
			name: "InvalidPasswordNotLength",
			req: &RegisterReq{
				UserName:    "mcagricaliskan_test",
				DisplayName: "Mehmet Çağrı Çalışkan",
				Password:    "hTM!+n1",
				EMail:       "mcagricaliskan_test@gmail.com",
			},
			resp: ErrInvalidPassword,
		},
	}

	authDatabaseRepository := NewAuthDatabaseRepository(database)
	authService := NewAuthService(authDatabaseRepository)
	ctx := context.Background()

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			err := authService.register(ctx, tc.req)
			if err != tc.resp {
				t.Errorf("Test Failed, Name: %s, Expected: %v, Got: %v", tc.name, tc.resp, err)
			}
		})
	}
}
