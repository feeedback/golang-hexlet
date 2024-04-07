// https://ru.hexlet.io/courses/go-web-development/lessons/error-handling/exercise_unit
package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
)

type (
	SendPushNotificationRequest struct {
		Message string `json:"message"`
		UserID  int64  `json:"user_id"`
	}

	PushNotification struct {
		Message string `json:"message"`
		UserID  int64  `json:"user_id"`
	}
)

var pushNotificationsQueue []PushNotification

func createAppRoutes() *fiber.App {
	// BEGIN (write your solution here)
	webApp := fiber.New(fiber.Config{
		ReadTimeout:  300 * time.Millisecond,
		WriteTimeout: 300 * time.Millisecond,
	})

	webApp.Use(recover.New())
	// END
	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	webApp.Post("/push/send", func(c *fiber.Ctx) error {
		var req SendPushNotificationRequest
		if err := c.BodyParser(&req); err != nil {
			// BEGIN (write your solution here)
			return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON")
			// END
		}

		pushNotificationsQueue = append(pushNotificationsQueue, PushNotification(req))
		if len(pushNotificationsQueue) > 3 {
			panic("Queue is full")
		}

		return c.SendStatus(fiber.StatusOK)
	})

	return webApp
}

func main() {
	webApp := createAppRoutes()

	logrus.Fatal(webApp.Listen(":8080"))
}
