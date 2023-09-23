package auth

import (
	"log"
	"os"
	"the-game-backend/storage/postgres"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Authentication struct {
	service AuthService
}

func RegisterRouters(app *fiber.App, database *postgres.Postgres) {
	a := Authentication{}

	app.Post("/register", a.Regsiter)
	app.Post("/login", a.Login)
}

func (a *Authentication) Regsiter(c *fiber.Ctx) error {

	user := RegisterReq{}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(422).JSON(fiber.Map{"message": "Unprocessable Entity"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (a Authentication) Login(c *fiber.Ctx) error {

	login := LoginReq{}
	if err := c.BodyParser(&login); err != nil {
		return c.Status(422).JSON(err)
	}

	isUserExists, user, err := a.repository.GetUser(c.Context(), login.Username)
	if err != nil {
		log.Println("auth -> Login -> GetUser -> Error while getting user, ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
	}
	if !isUserExists {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	if !CheckPasswordHash(login.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"id":  user.Id,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Println("auth -> Login -> SignedString -> Error while signing token, ", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}
