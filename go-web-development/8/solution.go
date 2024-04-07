// https://ru.hexlet.io/courses/go-web-development/lessons/local-persistence/exercise_unit
package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type (
	GetTaskResponse struct {
		ID       int64  `json:"id"`
		Desc     string `json:"description"`
		Deadline int64  `json:"deadline"`
	}

	CreateTaskRequest struct {
		Desc     string `json:"description"`
		Deadline int64  `json:"deadline"`
	}

	CreateTaskResponse struct {
		ID int64 `json:"id"`
	}

	UpdateTaskRequest struct {
		Desc     string `json:"description"`
		Deadline int64  `json:"deadline"`
	}

	Task struct {
		ID       int64
		Desc     string
		Deadline int64
	}
)

var (
	taskIDCounter int64 = 1
	tasks               = make(map[int64]Task)
)

func updateTaskHandler(c *fiber.Ctx) error {
	taskID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	existingTask, ok := tasks[taskID]
	if !ok {
		return c.Status(fiber.StatusNotFound).SendString("Not Found")
	}

	var updateTaskReq UpdateTaskRequest
	if err := c.BodyParser(&updateTaskReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	existingTask.Desc = updateTaskReq.Desc
	existingTask.Deadline = updateTaskReq.Deadline
	tasks[taskID] = existingTask

	return c.Status(fiber.StatusOK).SendString("OK")
}

func getTaskHandler(c *fiber.Ctx) error {
	taskID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	task, ok := tasks[taskID]
	if !ok {
		return c.Status(fiber.StatusNotFound).SendString("Not Found")
	}

	return c.Status(fiber.StatusOK).JSON(GetTaskResponse(task))
}

func deleteTaskHandler(c *fiber.Ctx) error {
	taskID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	_, ok := tasks[taskID]
	if !ok {
		return c.Status(fiber.StatusNotFound).SendString("Not Found")
	}

	delete(tasks, taskID)

	return c.Status(fiber.StatusOK).SendString("OK")
}

func createTaskHandler(c *fiber.Ctx) error {
	var createTaskReq CreateTaskRequest
	if err := c.BodyParser(&createTaskReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	taskID := taskIDCounter
	taskIDCounter++
	newTask := Task{
		ID:       taskID,
		Desc:     createTaskReq.Desc,
		Deadline: createTaskReq.Deadline,
	}

	tasks[taskID] = newTask

	return c.Status(fiber.StatusOK).JSON(CreateTaskResponse{ID: taskID})
}

func main() {
	webApp := fiber.New()
	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	// BEGIN
	webApp.Post("/tasks", createTaskHandler)
	webApp.Patch("/tasks/:id", updateTaskHandler)
	webApp.Get("/tasks/:id", getTaskHandler)
	webApp.Delete("/tasks/:id", deleteTaskHandler)
	// END

	logrus.Fatal(webApp.Listen(":8080"))
}
