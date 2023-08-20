package main

import (
	"log"
	"the-game-backend/controllers/auth"
	"the-game-backend/services/postgres"

	"github.com/gofiber/fiber/v2"
)

func main() {

	db, err := postgres.Connect()
	if err != nil {
		log.Fatal("Database Connection Can't Estabilished, Error:", err)
	}

	app := fiber.New()

	auth := auth.Authentication{DB: db}
	auth.Router(app)

	log.Fatal(app.Listen(":33333"))
}
