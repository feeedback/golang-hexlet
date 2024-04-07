// https://ru.hexlet.io/courses/go-web-development/lessons/auth/exercise_unit
package main

import (
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type (
	SignUpRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	SignInRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	SignInResponse struct {
		JWTToken string `json:"jwt_token"`
	}

	ProfileResponse struct {
		Email string `json:"email"`
	}

	User struct {
		Email    string
		password string
	}
)

var (
	webApiPort = ":8080"

	users = map[string]User{}

	secretKey = []byte("qwerty123456")

	contextKeyUser = "user"
)

func main() {
	webApp := fiber.New()
	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	// BEGIN (write your solution here) (write your solution here)
	publicGroup := webApp.Group("")

	publicGroup.Post("/signup", func(c *fiber.Ctx) error {
		var request SignUpRequest

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		_, exists := users[request.Email]
		if exists {
			return c.Status(fiber.StatusConflict).SendString("User with this email already exists")
		}

		users[request.Email] = User{
			Email:    request.Email,
			password: request.Password,
		}

		return c.SendStatus(fiber.StatusOK)
	})

	publicGroup.Post("/signin", func(c *fiber.Ctx) error {
		var request SignInRequest

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		user, exists := users[request.Email]
		if !exists || user.password != request.Password {
			return c.Status(fiber.StatusUnprocessableEntity).SendString("Invalid credentials")
		}

		payload := jwt.MapClaims{
			"sub": user.Email,
			"exp": time.Now().Add(time.Hour * 72).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

		t, err := token.SignedString(secretKey)
		if err != nil {
			logrus.WithError(err).Error("JWT token signing")
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(SignInResponse{JWTToken: t})
	})

	authorizedGroup := webApp.Group("")

	authorizedGroup.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: secretKey,
		},
		ContextKey: contextKeyUser,
	}))

	authorizedGroup.Get("/profile", func(c *fiber.Ctx) error {
		userToken := c.Context().Value(contextKeyUser).(*jwt.Token)
		claims := userToken.Claims.(jwt.MapClaims)
		email := claims["sub"].(string)

		user, exists := users[email]
		if !exists {
			return c.Status(fiber.StatusNotFound).SendString("User not found")
		}

		return c.JSON(ProfileResponse{Email: user.Email})
	})
	// END

	logrus.Fatal(webApp.Listen(webApiPort))
}
