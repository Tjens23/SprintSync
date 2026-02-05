package router

import (
	"encoding/gob"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"

	"SprintSync/src/platform/authenticator"
	"SprintSync/src/web/app/callback"
	"SprintSync/src/web/app/login"
	"SprintSync/src/web/app/logout"
	"SprintSync/src/web/app/user"
)

// New creates and configures a new Fiber v3 router
func New(auth *authenticator.Authenticator) *fiber.App {
	// Create Fiber app
	app := fiber.New()

	// Register custom types for session storage
	gob.Register(map[string]interface{}{})

	// Setup session store with cookie configuration
	store := session.NewStore(session.Config{
		CookieSecure:   false, // Set to true in production with HTTPS
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
	})

	// API routes
	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Auth API",
		})
	})

	// Auth routes
	app.Get("/login", login.Handler(auth, store))

	app.Get("/callback", callback.Handler(auth, store))

	app.Get("/user", user.Handler(store))

	app.Get("/logout", logout.Handler(store))

	return app
}
