// https://ru.hexlet.io/courses/go-web-development/lessons/routing/exercise_unit
package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

var postLikes = map[string]int64{}

func handleGetLike(c *fiber.Ctx) error {
	postID := c.Params("post_id")
	likes, ok := postLikes[postID]
	if !ok {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.SendString(strconv.FormatInt(likes, 10))
}

func handleCreateLike(c *fiber.Ctx) error {
	postID := c.Params("post_id")

	postLikes[postID]++

	likes := postLikes[postID]

	status := fiber.StatusOK
	if likes == 1 {
		status = fiber.StatusCreated
	}

	return c.Status(status).SendString(strconv.FormatInt(postLikes[postID], 10))
}

func main() {
	webApp := fiber.New(fiber.Config{Immutable: true})

	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Go to /likes/12345")
	})

	// BEGIN
	webApp.Get("/likes/:post_id", handleGetLike)

	webApp.Post("/likes/:post_id", handleCreateLike)
	// END

	logrus.Fatal(webApp.Listen(":8080"))
}
