package auth

import (
	"log"
	"the-game-backend/services/postgres"

	"github.com/gofiber/fiber/v2"
)

type Authentication struct {
	DB *postgres.Postgres
}

func (a Authentication) Router(app *fiber.App) {
	app.Post("/register", a.Regsiter)
	app.Post("/login", a.Login)
}

func (a Authentication) Regsiter(c *fiber.Ctx) error {

	user := RegisterReq{}
	if err := c.BodyParser(&user); err != nil {
		// returns
		// {
		// 	"code": 422,
		// 	"message": "Unprocessable Entity"
		// }
		return c.Status(422).JSON(err)
	}

	// i can move here to redis if user nubmer grows
	isUserExists, err := IsUserExists(a.DB, c.Context(), user.Username)
	if err != nil {
		log.Println("auth -> Register -> IsUserExists -> Error while checking user exists, ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
	}
	if isUserExists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "User Already Exists"})
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		log.Println("auth -> Register -> HashPassword -> Error while hashing password, ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
	}

	genderId, err := getGenderId(user.Gender)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Bad Request"})
	}
	characterGenderId, err := getCharacterGenderId(user.CharacterGender)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Bad Request"})
	}
	log.Println(genderId, characterGenderId)

	err = CreateUser(a.DB, c.Context(), &user, hashedPassword, genderId, characterGenderId)
	if err != nil {
		log.Println("auth -> Register -> CreateUser -> Error while creating user, ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (a Authentication) Login(c *fiber.Ctx) error {

	user := LoginReq{}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(422).JSON(err)
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
