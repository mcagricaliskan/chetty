package main

import (
	"log"
	"the-game-backend/app/controllers/auth"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users = make(map[string]string) // username -> hashed password

// func Register(c *fiber.Ctx) error {
// 	user := new(User)

// 	if err := c.BodyParser(user); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Failed to parse request",
// 		})
// 	}

// 	hashedPassword, err := HashPassword(user.Password)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Failed to hash password",
// 		})
// 	}

// 	users[user.Username] = hashedPassword

// 	return c.JSON(fiber.Map{
// 		"message": "Registered successfully",
// 	})
// }

// func Login(c *fiber.Ctx) error {
// 	user := new(User)

// 	if err := c.BodyParser(user); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Failed to parse request",
// 		})
// 	}

// 	hashedPassword, exists := users[user.Username]
// 	if !exists || !CheckPasswordHash(user.Password, hashedPassword) {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": "Invalid login credentials",
// 		})
// 	}

// 	return c.JSON(fiber.Map{
// 		"message": "Logged in successfully",
// 	})
// }

func main() {
	app := fiber.New()

	app.Post("/register", auth.Register)
	// app.Post("/login", auth.Login)

	log.Fatal(app.Listen(":33333"))
}
