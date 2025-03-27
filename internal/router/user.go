package router

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/api"
	"seaotterms.com-backend/internal/middleware"
)

func UserRouter(routerGroup fiber.Router, store *session.Store, dbs map[string]*gorm.DB) {
	userGroup := routerGroup.Group("/users")
	userGroup.Post("/", func(c *fiber.Ctx) error {
		return api.CreateUser(c, dbs[os.Getenv("DB_NAME3")])
	})
	userGroup.Patch("/:id", middleware.CheckLogin(store, dbs[os.Getenv("DB_NAME3")]), func(c *fiber.Ctx) error {
		return api.UpdateUser(c, dbs[os.Getenv("DB_NAME3")])
	})
}
