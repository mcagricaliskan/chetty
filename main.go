package main

import (
	"log"
	"the-game-backend/controllers/auth"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	//connection := pgx.ConnConfig{}

	auth := auth.Authentication{Repository: auth.AuthenticationRepository{Connection: connection}}
	auth.Router(app)

	log.Fatal(app.Listen(":33333"))
}
