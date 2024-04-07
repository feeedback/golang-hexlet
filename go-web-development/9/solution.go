// https://ru.hexlet.io/courses/go-web-development/lessons/validation/exercise_unit
package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type User struct {
	ID      int64
	Email   string
	Age     int
	Country string
}

var users = map[int64]User{}

type (
	CreateUserRequest struct {
		// BEGIN (write your solution here)
		ID      int64  `json:"id" validate:"required,gt=0"`
		Email   string `json:"email" validate:"required,email"`
		Age     int    `json:"age" validate:"required,gte=18,lte=130"`
		Country string `json:"country" validate:"required,oneof=USA Germany France"`
		// END
	}
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func createUser(c *fiber.Ctx) error {
	var req CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
	}

	users[req.ID] = User(req)

	return c.SendStatus(fiber.StatusOK)
}

func main() {
	webApp := fiber.New()
	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	// BEGIN (write your solution here) (write your solution here)
	webApp.Post("/users", createUser)
	// END

	logrus.Fatal(webApp.Listen(":8080"))
}
