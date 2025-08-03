package router

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"seaotterms.com-backend/api"
	"seaotterms.com-backend/middleware"
)

func TodoTopicRouter(routerGroup fiber.Router, store *session.Store, dbs map[string]*gorm.DB) {
	todoTopicGroup := routerGroup.Group("/todo-topics")
	dbName := os.Getenv("DB_NAME")

	todoTopicGroup.Get("/:owner", func(c *fiber.Ctx) error {
		return api.QueryTodoTopic(c, dbs[dbName])
	})
	todoTopicGroup.Post("/", middleware.CheckLogin(store, dbs[dbName]), func(c *fiber.Ctx) error {
		return api.CreateTodoTopic(c, dbs[dbName])
	})
}
