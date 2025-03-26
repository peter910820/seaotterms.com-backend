package router

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/api"
	"seaotterms.com-backend/internal/middleware"
)

func ArticleRouter(routerGroup fiber.Router, store *session.Store, dbs map[string]*gorm.DB) {
	articleGroup := routerGroup.Group("/articles")
	articleGroup.Get("/", func(c *fiber.Ctx) error {
		return api.QueryArticle(c, dbs[os.Getenv("DB_NAME")])
	})
	articleGroup.Post("/", middleware.CheckLogin(store, dbs[os.Getenv("DB_NAME3")]), func(c *fiber.Ctx) error {
		return api.CreateArticle(c, dbs[os.Getenv("DB_NAME")])
	})
	articleGroup.Get("/:articleID", func(c *fiber.Ctx) error {
		return api.QuerySingleArticle(c, dbs[os.Getenv("DB_NAME")])
	})
}
