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
	repository AuthRepository
}

func RegisterRouters(app *fiber.App, database *postgres.Postgres) {
	a := Authentication{
		repository: &AuthRepo{database: database}}

	app.Post("/register", a.Regsiter)
	app.Post("/login", a.Login)
}

func (a *Authentication) Regsiter(c *fiber.Ctx) error {

	user := RegisterReq{}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(422).JSON(fiber.Map{"message": "Unprocessable Entity"})
	}

	// i can move here to redis if user nubmer grows
	isUserExists, err := a.repository.IsUserExists(c.Context(), user.Username)
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
		log.Println("auth -> Register -> getGenderId -> Error while getting gender id, ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Bad Request"})
	}
	characterGenderId, err := getCharacterGenderId(user.CharacterGender)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Bad Request"})
	}
	log.Println(genderId, characterGenderId)

	err = a.repository.CreateUser(c.Context(), &user, hashedPassword, genderId, characterGenderId)
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
