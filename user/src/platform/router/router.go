package router

import (
	"SprintSync/src/platform/messaging"
	"SprintSync/src/web/app/users"

	"github.com/gofiber/fiber/v3"
)

func New(rabbitMQ *messaging.RabbitMQ) *fiber.App {
	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "User Service API",
		})
	})

	api := app.Group("/api/v1")

	api.Get("/users", users.ListHandler(rabbitMQ))
	api.Get("/users/:id", users.GetHandler(rabbitMQ))
	api.Post("/users", users.CreateHandler(rabbitMQ))
	api.Put("/users/:id", users.UpdateHandler(rabbitMQ))
	api.Delete("/users/:id", users.DeleteHandler(rabbitMQ))

	return app
}
