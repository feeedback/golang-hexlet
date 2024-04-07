// https://ru.hexlet.io/courses/go-web-development/lessons/fiber/exercise_unit
package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

var exchangeRate = map[string]float64{
	"USD/EUR": 0.8,
	"EUR/USD": 1.25,
	"USD/GBP": 0.7,
	"GBP/USD": 1.43,
	"USD/JPY": 110,
	"JPY/USD": 0.0091,
}

func convertHandler(ctx *fiber.Ctx) error {
	from := ctx.Query("from")
	to := ctx.Query("to")

	pair := fmt.Sprintf("%s/%s", from, to)

	rate, ok := exchangeRate[pair]
	if !ok {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.SendString(fmt.Sprintf("%.2f", rate))
}

func main() {
	// BEGIN (write your solution here)
	webApp := fiber.New()

	webApp.Get("/convert", convertHandler)
	// END

	logrus.Fatal(webApp.Listen(":8080"))
}
