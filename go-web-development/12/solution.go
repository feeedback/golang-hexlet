// https://ru.hexlet.io/courses/go-web-development/lessons/template/exercise_unit
package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/sirupsen/logrus"
)

type (
	CreateItemRequest struct {
		Name  string `json:"name"`
		Price uint   `json:"price"`
	}

	Item struct {
		Name  string `json:"name"`
		Price uint   `json:"price"`
	}
)

var (
	items []Item
)

func main() {
	viewsEngine := html.New("./templates", ".tmpl")

	webApp := fiber.New(fiber.Config{
		Views: viewsEngine,
	})

	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	// BEGIN (write your solution here)
	webApp.Post("/items", func(c *fiber.Ctx) error {
		var item CreateItemRequest

		if err := c.BodyParser(&item); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		items = append(items, Item(item))
		return c.SendStatus(fiber.StatusOK)
	})

	webApp.Get("/items/view", func(c *fiber.Ctx) error {
		return c.Render("items", items)
	})
	// END

	logrus.Fatal(webApp.Listen(":8080"))
}
