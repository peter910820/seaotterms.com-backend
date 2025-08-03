package router

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"seaotterms.com-backend/api"
	"seaotterms.com-backend/middleware"
)

func GalgameRouter(routerGroup fiber.Router, store *session.Store, dbs map[string]*gorm.DB) {
	galgameGroup := routerGroup.Group("/galgame")
	dbName := os.Getenv("DB_NAME2")

	galgameGroup.Get("/s/:name", middleware.CheckLogin(store, dbs[dbName]), func(c *fiber.Ctx) error {
		return api.QueryGalgame(c, dbs[os.Getenv("DB_NAME2")])
	})
	galgameGroup.Get("/:brand", func(c *fiber.Ctx) error {
		return api.QueryGalgameByBrand(c, dbs[dbName])
	})
	galgameGroup.Patch("/develop/:name", middleware.CheckOwner(store, dbs[os.Getenv("DB_NAME")]), func(c *fiber.Ctx) error {
		return api.UpdateGalgameDevelop(c, dbs[dbName])
	})
	galgameGroup.Post("/", middleware.CheckOwner(store, dbs[os.Getenv("DB_NAME")]), func(c *fiber.Ctx) error {
		return api.CreateGalgame(c, dbs[dbName])
	})
}
