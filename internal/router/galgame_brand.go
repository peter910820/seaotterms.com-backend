package router

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/api"
	"seaotterms.com-backend/internal/middleware"
)

func GalgameBrandRouter(routerGroup fiber.Router, store *session.Store, dbs map[string]*gorm.DB) {
	galgameBrandGroup := routerGroup.Group("/galgame-brand")
	galgameBrandGroup.Get("/", func(c *fiber.Ctx) error {
		return api.QueryAllGalgameBrand(c, dbs[os.Getenv("DB_NAME2")])
	})
	galgameBrandGroup.Get("/:brand", func(c *fiber.Ctx) error {
		return api.QueryGalgameBrand(c, dbs[os.Getenv("DB_NAME2")])
	})
	galgameBrandGroup.Post("/", middleware.CheckLogin(store, dbs[os.Getenv("DB_NAME3")]), func(c *fiber.Ctx) error {
		return api.InsertGalgameBrand(c, dbs[os.Getenv("DB_NAME2")])
	})
	galgameBrandGroup.Patch("/:brand", middleware.CheckLogin(store, dbs[os.Getenv("DB_NAME3")]), func(c *fiber.Ctx) error {
		return api.UpdateGalgameBrand(c, dbs[os.Getenv("DB_NAME2")])
	})
}
