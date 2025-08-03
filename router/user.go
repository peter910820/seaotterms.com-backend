package router

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"seaotterms.com-backend/api"
	"seaotterms.com-backend/middleware"
)

func UserRouter(routerGroup fiber.Router, store *session.Store, dbs map[string]*gorm.DB) {
	userGroup := routerGroup.Group("/users")
	dbName := os.Getenv("DB_NAME")

	userGroup.Post("/", func(c *fiber.Ctx) error {
		return api.CreateUser(c, dbs[dbName])
	})
	userGroup.Patch("/:id", middleware.CheckLogin(store, dbs[dbName]), func(c *fiber.Ctx) error {
		return api.UpdateUser(c, dbs[dbName])
	})
}
