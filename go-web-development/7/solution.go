// https://ru.hexlet.io/courses/go-web-development/lessons/local-persistence/exercise_unit
package main

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type (
	CreateLinkRequest struct {
		External string `json:"external"`
		Internal string `json:"internal"`
	}

	GetLinkResponse struct {
		Internal string `json:"internal"`
	}
)

var links = make(map[string]string)

func handlerCreateLink(c *fiber.Ctx) error {
	link := CreateLinkRequest{}

	if err := c.BodyParser(&link); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON")
	}

	links[link.External] = link.Internal

	return c.SendStatus(fiber.StatusOK)
}

func handlerGetLink(c *fiber.Ctx) error {
	external, err := url.QueryUnescape(c.Params("external"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid URL")
	}

	internal, ok := links[external]
	if !ok {
		return c.Status(fiber.StatusNotFound).SendString("Link not found")
	}

	response := GetLinkResponse{Internal: internal}
	return c.JSON(response)
}

func main() {
	webApp := fiber.New()
	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	// BEGIN (write your solution here)
	webApp.Post("/links", handlerCreateLink)

	webApp.Get("/links/:external", handlerGetLink)
	// END

	logrus.Fatal(webApp.Listen(":8080"))
}
