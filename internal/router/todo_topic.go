package router

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/api"
	"seaotterms.com-backend/internal/middleware"
)

func TodoTopicRouter(routerGroup fiber.Router, store *session.Store, dbs map[string]*gorm.DB) {
	todoTopicGroup := routerGroup.Group("/todo-topics")
	todoTopicGroup.Get("/:owner", func(c *fiber.Ctx) error {
		return api.QueryTodoTopic(c, dbs[os.Getenv("DB_NAME3")])
	})
	todoTopicGroup.Post("/", middleware.CheckLogin(store, dbs[os.Getenv("DB_NAME3")]), func(c *fiber.Ctx) error {
		return api.InsertTodoTopic(c, dbs[os.Getenv("DB_NAME3")])
	})
}
