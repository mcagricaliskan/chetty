package main

import (
	"log"
	"os"
	"the-game-backend/controllers/auth"
	"the-game-backend/services/postgres"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	if os.Getenv("ENVIROMENT") != "PRODUCTION" {
		godotenv.Load()
	}

	db, err := postgres.Connect()
	if err != nil {
		log.Fatal("Database Connection Can't Estabilished, Error:", err)
	} else {
		log.Println("Database Connection Estabilished")
	}

	app := fiber.New()

	auth := auth.Authentication{DB: db}
	auth.Router(app)

	log.Fatal(app.Listen(":33333"))
}
