package auth

import (
	"log"
	"testing"
	"the-game-backend/storage/postgres"

	"github.com/gofiber/fiber/v2"
)

func TestRegister(t *testing.T) {

	app := fiber.New()

	database, err := postgres.Connect()
	if err != nil {
		log.Fatal("Database Connection Can't Estabilished, Error:", err)
	} else {
		log.Println("Database Connection Estabilished")
	}

	RegisterRouters(app, database)
}
