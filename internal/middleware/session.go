package middleware

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// check user identity
func SessionHandler(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		confirmRoutes := map[string]string{
			"/api/verify":         "POST", // handle front-end router verify
			"/api/create-article": "POST",
			"/api/galgame":        "POST",
			"/api/galgame-brand":  "POST",
		}
		confirmRoutesPrefix := map[string]string{
			"/api/galgame/":       "PATCH",
			"/api/galgame-brand/": "PATCH",
		}
		if isPathIn(c.Path(), c.Method(), confirmRoutes) {
			err := checkLogin(c, store)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"msg":      "visitors is not logged in",
					"username": "",
				})
			}
			return c.Next()
		}
		if isPathPrefix(c.Path(), c.Method(), confirmRoutesPrefix) {
			err := checkLogin(c, store)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"msg":      "visitors is not logged in",
					"username": "",
				})
			}
			return c.Next()
		}
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

func isPathPrefix(path string, method string, confirmRoutesPrefix map[string]string) bool {
	for key, value := range confirmRoutesPrefix {
		if strings.HasPrefix(path, key) && value == method {
			return true
		}
	}
	return false
}

func checkLogin(c *fiber.Ctx, store *session.Store) error {
	sess, err := store.Get(c)
	if err != nil {
		logrus.Fatal(err)
	}
	username := sess.Get("username")
	if username == nil {
		logrus.Error("visitors is not logged in")
		return errors.New("visitors is not logged in")
	}
	c.Locals("username", username)
	logrus.Infof("%s is access %s", username, c.Path())
	return nil
}
