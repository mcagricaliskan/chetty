package auth

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthenticationFiber struct {
	service AuthService
}

func RegisterRouters(app *fiber.App, service AuthService) {
	a := AuthenticationFiber{
		service: service,
	}

	app.Post("/register", a.Register)
	app.Post("/login", a.Login)
}

func (a *AuthenticationFiber) Register(c *fiber.Ctx) error {

	user := RegisterReq{}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(422).JSON(fiber.Map{"message": "Unprocessable Entity"})
	}

	err := a.service.register(c.Context(), &user)
	if err != nil {
		switch err {
		case ErrInternalServer:
			return c.SendStatus(fiber.StatusInternalServerError)
		case ErrUserExists:
			return c.SendStatus(fiber.StatusConflict)
		case ErrBadRequest:
			return c.SendStatus(fiber.StatusBadRequest)
		default:
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

func (a AuthenticationFiber) Login(c *fiber.Ctx) error {

	login := LoginReq{}
	if err := c.BodyParser(&login); err != nil {
		return c.Status(422).JSON(err)
	}

	user, err := a.service.login(c.Context(), &login)
	if err != nil {
		switch err {
		case ErrInternalServer:
			return c.SendStatus(fiber.StatusInternalServerError)
		case ErrUnauthoerized:
			return c.SendStatus(fiber.StatusUnauthorized)
		default:
			return c.SendStatus(fiber.StatusInternalServerError)
		}
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
