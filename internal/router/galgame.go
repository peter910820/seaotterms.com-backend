package router

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/api"
	"seaotterms.com-backend/internal/middleware"
)

func GalgameRouter(routerGroup fiber.Router, store *session.Store, dbs map[string]*gorm.DB) {
	galgameGroup := routerGroup.Group("/galgame")
	galgameGroup.Get("/s/:name", middleware.CheckLogin(store, dbs[os.Getenv("DB_NAME3")]), func(c *fiber.Ctx) error {
		return api.QueryGalgame(c, dbs[os.Getenv("DB_NAME2")])
	})
	galgameGroup.Get("/:brand", func(c *fiber.Ctx) error {
		return api.QueryGalgameByBrand(c, dbs[os.Getenv("DB_NAME2")])
	})
	galgameGroup.Patch("/develop/:name", middleware.CheckLogin(store, dbs[os.Getenv("DB_NAME3")]), func(c *fiber.Ctx) error {
		return api.UpdateGalgameDevelop(c, dbs[os.Getenv("DB_NAME2")])
	})
	galgameGroup.Post("/", middleware.CheckLogin(store, dbs[os.Getenv("DB_NAME3")]), func(c *fiber.Ctx) error {
		return api.InsertGalgame(c, dbs[os.Getenv("DB_NAME2")])
	})
}
