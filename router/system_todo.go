package router

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"seaotterms.com-backend/api"
	"seaotterms.com-backend/middleware"
)

func SystemTodoRouter(routerGroup fiber.Router, store *session.Store, dbs map[string]*gorm.DB) {
	systemTodoGroup := routerGroup.Group("/system-todos")
	dbName := os.Getenv("DB_NAME")

	systemTodoGroup.Get("/", func(c *fiber.Ctx) error {
		return api.QuerySystemTodo(c, dbs[dbName])
	})

	systemTodoGroup.Post("/", middleware.CheckLogin(store, dbs[dbName]), func(c *fiber.Ctx) error {
		return api.CreateSystemTodo(c, dbs[dbName])
	})

	systemTodoGroup.Patch("/:id", middleware.CheckOwner(store, dbs[dbName]), func(c *fiber.Ctx) error {
		return api.UpdateSystemTodo(c, dbs[dbName])
	})

	// quick update
	systemTodoGroup.Patch("/quick/:id", middleware.CheckOwner(store, dbs[dbName]), func(c *fiber.Ctx) error {
		return api.QuickUpdateSystemTodo(c, dbs[dbName])
	})

	systemTodoGroup.Delete("/:id", middleware.CheckOwner(store, dbs[dbName]), func(c *fiber.Ctx) error {
		return api.DeleteSystemTodo(c, dbs[dbName])
	})
}
