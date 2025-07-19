package router

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/api"
)

func TagRouter(routerGroup fiber.Router, store *session.Store, dbs map[string]*gorm.DB) {
	tagGroup := routerGroup.Group("/tags")
	dbName := os.Getenv("DB_NAME")

	tagGroup.Get("/", func(c *fiber.Ctx) error {
		return api.QueryTag(c, dbs[dbName])
	})
	tagGroup.Get("/:tagName", func(c *fiber.Ctx) error {
		return api.QueryArticleForTag(c, dbs[dbName])
	})
}
