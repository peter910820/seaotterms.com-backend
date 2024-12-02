package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// closure
func SessionHandler(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Printf("store: %v", store)
		confirmRoutes := []string{"/create-article", "/api/check-session"}
		if !isPathIn(c.Path(), confirmRoutes) {
			return c.Next()
		}
		// check session
		log.Println("check session")
		sess, err := store.Get(c)
		if err != nil {
			log.Fatal(err.Error())
		}
		username := sess.Get("username")
		if username == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"msg": "user is signout"})
		}
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
