package middleware

import (
	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// closure
func SessionHandler(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		confirmRoutes := []string{"/create-article", "/api/verify"}
		if !isPathIn(c.Path(), confirmRoutes) {
			return c.Next()
		}
		logrus.Infof("path %s execute middleware", c.Path())
		logrus.Debugf("store: %v", store)
		// check session
		logrus.Info("check session")
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

func isPathIn(path string, confirmRoutes []string) bool {
	for _, r := range confirmRoutes {
		if path == r {
			return true
		}
	}
	return false
}
