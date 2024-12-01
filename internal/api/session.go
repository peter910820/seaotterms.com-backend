package api

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func CheckSession(c *fiber.Ctx, store *session.Store) error {
	sess, err := store.Get(c)
	if err != nil {
		log.Fatal(err.Error())
	}
	username := sess.Get("username")
	if username == nil {
		return errors.New("user is signout")
	}
	return nil

}
