// https://ru.hexlet.io/courses/go-web-development/lessons/data-serialization/exercise_unit
package main

import (
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type (
	BinarySearchRequest struct {
		Numbers []int `json:"numbers"`
		Target  int   `json:"target"`
	}

	BinarySearchResponse struct {
		TargetIndex int    `json:"target_index"`
		Error       string `json:"error,omitempty"`
	}
)

const targetNotFound = -1

func handlerSearch(c *fiber.Ctx) error {
	reqBody := new(BinarySearchRequest)

	err := c.BodyParser(reqBody)

	if err != nil {
		res := BinarySearchResponse{
			TargetIndex: targetNotFound,
			Error:       "Invalid JSON",
		}
		return c.Status(400).JSON(res)
	}

	targetIndex := sort.SearchInts(reqBody.Numbers, reqBody.Target)

	if targetIndex == len(reqBody.Numbers) || reqBody.Numbers[targetIndex] != reqBody.Target {
		res := BinarySearchResponse{
			TargetIndex: targetNotFound,
			Error:       "Target was not found",
		}
		return c.Status(404).JSON(res)
	}

	res := BinarySearchResponse{
		TargetIndex: targetIndex,
	}
	return c.JSON(res)
}

func main() {
	webApp := fiber.New()
	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	// BEGIN (write your solution here)
	webApp.Post("/search", handlerSearch)
	// END

	logrus.Fatal(webApp.Listen(":8080"))
}
