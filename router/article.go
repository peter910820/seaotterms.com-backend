package router

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"seaotterms.com-backend/api"
	"seaotterms.com-backend/middleware"
)

func ArticleRouter(routerGroup fiber.Router, store *session.Store, dbs map[string]*gorm.DB) {
	articleGroup := routerGroup.Group("/articles")
	dbName := os.Getenv("DB_NAME")

	articleGroup.Get("/", func(c *fiber.Ctx) error {
		return api.QueryArticle(c, dbs[dbName])
	})
	articleGroup.Post("/", middleware.CheckLogin(store, dbs[os.Getenv("DB_NAME")]), func(c *fiber.Ctx) error {
		return api.CreateArticle(c, dbs[dbName])
	})
}
