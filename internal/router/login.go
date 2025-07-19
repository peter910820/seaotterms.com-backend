package router

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/api"
)

func LoginRouter(routerGroup fiber.Router, store *session.Store, dbs map[string]*gorm.DB) {
	loginGroup := routerGroup.Group("/login")
	dbName := os.Getenv("DB_NAME")

	loginGroup.Post("/", func(c *fiber.Ctx) error {
		return api.Login(c, store, dbs[dbName])
	})
}
