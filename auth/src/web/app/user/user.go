package user

import (
	"SprintSync/src/platform/messaging"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

// Handler for our logged-in user page.
func Handler(store *session.Store, rabbitMQ *messaging.RabbitMQ) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Get session from store
		sess, err := store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		profile := sess.Get("profile")

		if profile == nil {
			return c.Status(fiber.StatusUnauthorized).SendString("User not authenticated")
		}

		// Publish user profile viewed event
		if profileMap, ok := profile.(map[string]interface{}); ok {
			rabbitMQ.Publish(messaging.EventUserProfileViewed, map[string]interface{}{
				"user_id": profileMap["sub"],
				"email":   profileMap["email"],
			})
		}

		return c.JSON(profile)
	}
}
