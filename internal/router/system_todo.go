package router

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/api"
)

func SystemTodoRouter(routerGroup fiber.Router, store *session.Store, dbs map[string]*gorm.DB) {
	systemTodoGroup := routerGroup.Group("/system-todos")
	systemTodoGroup.Get("/", func(c *fiber.Ctx) error {
		return api.QuerySystemTodo(c, dbs[os.Getenv("DB_NAME3")])
	})
}
