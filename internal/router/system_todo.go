package router

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/api"
	"seaotterms.com-backend/internal/middleware"
)

func SystemTodoRouter(routerGroup fiber.Router, store *session.Store, dbs map[string]*gorm.DB) {
	systemTodoGroup := routerGroup.Group("/system-todos")

	systemTodoGroup.Get("/", func(c *fiber.Ctx) error {
		return api.QuerySystemTodo(c, dbs[os.Getenv("DB_NAME3")])
	})

	systemTodoGroup.Post("/", middleware.CheckLogin(store, dbs[os.Getenv("DB_NAME3")]), func(c *fiber.Ctx) error {
		return api.CreateSystemTodo(c, dbs[os.Getenv("DB_NAME3")])
	})

	systemTodoGroup.Patch("/:id", middleware.CheckOwner(store, dbs[os.Getenv("DB_NAME3")]), func(c *fiber.Ctx) error {
		return api.UpdateSystemTodo(c, dbs[os.Getenv("DB_NAME3")])
	})

	// quick update
	systemTodoGroup.Patch("/quick/:id", middleware.CheckOwner(store, dbs[os.Getenv("DB_NAME3")]), func(c *fiber.Ctx) error {
		return api.QuickUpdateSystemTodo(c, dbs[os.Getenv("DB_NAME3")])
	})

	systemTodoGroup.Delete("/:id", middleware.CheckOwner(store, dbs[os.Getenv("DB_NAME3")]), func(c *fiber.Ctx) error {
		return api.DeleteSystemTodo(c, dbs[os.Getenv("DB_NAME3")])
	})
}
