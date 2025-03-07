package middleware

import (
	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"slices"
)

// closure
func SessionHandler(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		confirmRoutes := []string{
			"/api/verify",
			"/api/create-article",
			"/api/galgame",
			"/api/galgame-brand",
		}
		if !isPathIn(c.Path(), confirmRoutes) {
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
		logrus.Debugf("store: %v", store)
		logrus.Infof("%s is access %s", username, c.Path())
		return c.Next()
	}

}

func isPathIn(path string, confirmRoutes []string) bool {
	return slices.Contains(confirmRoutes, path)
}
