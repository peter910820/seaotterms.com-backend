package router

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"seaotterms.com-backend/api"
	"seaotterms.com-backend/middleware"
)

func GalgameBrandRouter(routerGroup fiber.Router, store *session.Store, dbs map[string]*gorm.DB) {
	galgameBrandGroup := routerGroup.Group("/galgame-brand")
	dbName := os.Getenv("DB_NAME2")

	galgameBrandGroup.Get("/", func(c *fiber.Ctx) error {
		return api.QueryAllGalgameBrand(c, dbs[dbName])
	})
	galgameBrandGroup.Get("/:brand", func(c *fiber.Ctx) error {
		return api.QueryGalgameBrand(c, dbs[dbName])
	})
	galgameBrandGroup.Post("/", middleware.CheckOwner(store, dbs[os.Getenv("DB_NAME")]), func(c *fiber.Ctx) error {
		return api.CreateGalgameBrand(c, dbs[dbName])
	})
	galgameBrandGroup.Patch("/:brand", middleware.CheckOwner(store, dbs[os.Getenv("DB_NAME")]), func(c *fiber.Ctx) error {
		return api.UpdateGalgameBrand(c, dbs[dbName])
	})
}
