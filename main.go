package main

import (
	"log"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/mcagricaliskan/chetty/internal/auth"
	"github.com/mcagricaliskan/chetty/storage/postgres"
)

func main() {

	if os.Getenv("ENVIROMENT") != "PRODUCTION" {
		godotenv.Load()
	}

	databaseUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatal("Database URL Not Found")
	}

	database, err := postgres.Connect(databaseUrl)
	if err != nil {
		log.Fatal("Database Connection Can't Estabilished, Error:", err)
	} else {
		log.Println("Database Connection Estabilished")
	}

	app := fiber.New()

	auth.RegisterRouters(app,
		*auth.NewAuthService(
			auth.NewAuthDatabaseRepository(database),
		),
	)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
	}))

	log.Fatal(app.Listen(":33333"))
}
