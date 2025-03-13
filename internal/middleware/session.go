package middleware

import (
	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// check user identity
func SessionHandler(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Method()
		confirmRoutes := map[string]string{
			"/api/verify":         "POST",
			"/api/create-article": "POST",
			"/api/galgame":        "POST",
			"/api/galgame-brand":  "POST",
		}
		if !isPathIn(c.Path(), c.Method(), confirmRoutes) {
			return c.Next()
		}
		sess, err := store.Get(c)
		if err != nil {
			logrus.Fatal(err)
		}
		username := sess.Get("username")
		if username == nil {
			logrus.Infof("visitors is not logged in")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"msg": "visitors is not logged in"})
		}
		logrus.Infof("%s is access %s", username, c.Path())
		return c.Next()
	}
}

func isPathIn(path string, method string, confirmRoutes map[string]string) bool {
	for key, value := range confirmRoutes {
		if key == path && value == method {
			return true
		}
	}
	return false
}
