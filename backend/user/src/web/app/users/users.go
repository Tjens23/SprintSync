package users

import (
	"encoding/json"
	"time"

	"SprintSync/src/platform/messaging"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Picture   string    `json:"picture,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var users = make(map[string]*User)

func ListHandler(rabbitMQ *messaging.RabbitMQ) fiber.Handler {
	return func(c fiber.Ctx) error {
		userList := make([]*User, 0, len(users))
		for _, user := range users {
			userList = append(userList, user)
		}
		return c.JSON(userList)
	}
}

func GetHandler(rabbitMQ *messaging.RabbitMQ) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		user, exists := users[id]
		if !exists {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.JSON(user)
	}
}

func CreateHandler(rabbitMQ *messaging.RabbitMQ) fiber.Handler {
	return func(c fiber.Ctx) error {
		var input User
		if err := json.Unmarshal(c.Body(), &input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if input.Email == "" || input.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email and name are required",
			})
		}

		user := &User{
			ID:        uuid.New().String(),
			Email:     input.Email,
			Name:      input.Name,
			Picture:   input.Picture,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		users[user.ID] = user

		rabbitMQ.Publish(messaging.EventUserCreated, map[string]interface{}{
			"user_id": user.ID,
			"email":   user.Email,
			"name":    user.Name,
		})

		return c.Status(fiber.StatusCreated).JSON(user)
	}
}

func UpdateHandler(rabbitMQ *messaging.RabbitMQ) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		user, exists := users[id]
		if !exists {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		var input User
		if err := json.Unmarshal(c.Body(), &input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if input.Name != "" {
			user.Name = input.Name
		}
		if input.Email != "" {
			user.Email = input.Email
		}
		if input.Picture != "" {
			user.Picture = input.Picture
		}
		user.UpdatedAt = time.Now()

		rabbitMQ.Publish(messaging.EventUserUpdated, map[string]interface{}{
			"user_id": user.ID,
			"email":   user.Email,
			"name":    user.Name,
		})

		return c.JSON(user)
	}
}

func DeleteHandler(rabbitMQ *messaging.RabbitMQ) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		user, exists := users[id]
		if !exists {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		delete(users, id)

		rabbitMQ.Publish(messaging.EventUserDeleted, map[string]interface{}{
			"user_id": user.ID,
			"email":   user.Email,
		})

		return c.SendStatus(fiber.StatusNoContent)
	}
}
