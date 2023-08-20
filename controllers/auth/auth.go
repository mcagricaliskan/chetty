package auth

import (
	"log"
	"os"
	"the-game-backend/services/postgres"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Authentication struct {
	DB *postgres.Postgres
}

func (a Authentication) RegisterRouters(app *fiber.App) {
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

	login := LoginReq{}
	if err := c.BodyParser(&login); err != nil {
		return c.Status(422).JSON(err)
	}

	isUserExists, user, err := GetUser(a.DB, c.Context(), login.Username)
	if err != nil || !isUserExists {
		log.Println("auth -> Login -> GetUser -> Error while getting user, ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
	}
	if !isUserExists {
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
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}
