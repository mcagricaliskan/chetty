package auth

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthenticationHandler struct {
	service authService
}

func RegisterRouters(app *fiber.App, service authService) {
	a := AuthenticationHandler{
		service: service,
	}

	app.Post("/register", a.Register)
	app.Post("/login", a.Login)
}

func (a *AuthenticationHandler) Register(c *fiber.Ctx) error {

	user := RegisterReq{}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(422).JSON(fiber.Map{"message": "Unprocessable Entity"})
	}

	err := a.service.register(c.Context(), &user)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidPassword):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid Password"})
		case errors.Is(err, ErrInternalServer):
			return c.SendStatus(fiber.StatusInternalServerError)
		case errors.Is(err, ErrUserExists):
			return c.SendStatus(fiber.StatusConflict)
		case errors.Is(err, ErrBadRequest):
			return c.SendStatus(fiber.StatusBadRequest)
		default:
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (a AuthenticationHandler) Login(c *fiber.Ctx) error {

	login := LoginReq{}
	if err := c.BodyParser(&login); err != nil {
		return c.Status(422).JSON(err)
	}

	userId, err := a.service.login(c.Context(), &login)
	if err != nil {
		switch err {
		case ErrInternalServer:
			return c.SendStatus(fiber.StatusInternalServerError)
		case ErrUnauthorized:
			return c.SendStatus(fiber.StatusUnauthorized)
		default:
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"id":  userId,
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
