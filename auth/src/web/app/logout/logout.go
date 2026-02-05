package logout

import (
	"net/url"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

// Handler for our logout.
func Handler(store *session.Store) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Get session from store
		sess, err := store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// Destroy the session
		if err := sess.Destroy(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		logoutUrl, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		scheme := "http"
		if c.Protocol() == "https" {
			scheme = "https"
		}

		// Use BaseURL() to get the full base URL including scheme, host, and port
		returnTo, err := url.Parse(scheme + "://" + string(c.Request().Host()))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		parameters := url.Values{}
		parameters.Add("returnTo", returnTo.String())
		parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
		logoutUrl.RawQuery = parameters.Encode()

		return c.Redirect().To(logoutUrl.String())
	}
}
