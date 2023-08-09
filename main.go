package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users = make(map[string]string) // username -> hashed password

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Register(c *fiber.Ctx) error {
	user := new(User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
		})
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	users[user.Username] = hashedPassword

	return c.JSON(fiber.Map{
		"message": "Registered successfully",
	})
}

func Login(c *fiber.Ctx) error {
	user := new(User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
		})
	}

	hashedPassword, exists := users[user.Username]
	if !exists || !CheckPasswordHash(user.Password, hashedPassword) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid login credentials",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Logged in successfully",
	})
}

func main() {
	app := fiber.New()

	app.Post("/register", Register)
	app.Post("/login", Login)

	log.Fatal(app.Listen(":3000"))
}
