package callback

import (
	"SprintSync/src/platform/authenticator"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

// Handler for our callback.
func Handler(auth *authenticator.Authenticator, store *session.Store) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Get session from store
		sess, err := store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// Verify state parameter
		if c.Query("state") != sess.Get("state") {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid state parameter.")
		}

		// Exchange an authorization code for a token.
		token, err := auth.Exchange(c.Context(), c.Query("code"))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Failed to exchange an authorization code for a token.")
		}

		idToken, err := auth.VerifyIDToken(c.Context(), token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to verify ID Token.")
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		sess.Set("access_token", token.AccessToken)
		sess.Set("profile", profile)
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// Redirect to logged in page.
		return c.Redirect().To("/user")
	}
}
