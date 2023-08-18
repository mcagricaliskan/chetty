package auth

import (
	"log"
	"the-game-backend/internal/models/auth"

	"github.com/gofiber/fiber/v2"
)

func hashedPassword()

type Authentication struct {
	Repository Repository
}

func (a Authentication) Router(app *fiber.App) {
	app.Post("/register", a.Regsiter)
	app.Post("/login", a.Login)
}

func (a Authentication) Regsiter(c *fiber.Ctx) error {

	user := auth.RegisterModel{}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(422).JSON(err)
	}

	// i can move here to redis if user nubmer grows
	_, err := a.Repository.IsUserExists(c.Context(), user.Username)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Internal Server Error"})
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		// returns
		// {
		// 	"code": 422,
		// 	"message": "Unprocessable Entity"
		// }
		return err
	}

	log.Println(user, hashedPassword)

	return c.SendString("Register")
}

func (a Authentication) Login(c *fiber.Ctx) error {
	user := c.FormValue("user")
	pass := c.FormValue("pass")

	// Throws Unauthorized error
	if user != "john" || pass != "doe" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Create the Claims
	// claims := jwt.MapClaims{
	// 	"name":  "John Doe",
	// 	"admin": true,
	// 	"exp":   time.Now().Add(time.Hour * 72).Unix(),
	// }

	// // Create token
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// // Generate encoded token and send it as response.
	// t, err := token.SignedString([]byte("secret"))
	// if err != nil {
	// 	return c.SendStatus(fiber.StatusInternalServerError)
	// }

	return c.JSON(fiber.Map{"token": "ok"})
}
