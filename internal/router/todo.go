package router

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/api"
	"seaotterms.com-backend/internal/middleware"
)

func TodoRouter(routerGroup fiber.Router, store *session.Store, dbs map[string]*gorm.DB) {
	todoGroup := routerGroup.Group("/todos")
	dbName := os.Getenv("DB_NAME")

	todoGroup.Get("/:owner", func(c *fiber.Ctx) error {
		return api.QueryTodoByOwner(c, dbs[dbName])
	})
	todoGroup.Post("/", middleware.CheckLogin(store, dbs[dbName]), func(c *fiber.Ctx) error {
		return api.CreateTodo(c, dbs[dbName])
	})
	todoGroup.Patch("/:id", middleware.CheckLogin(store, dbs[dbName]), func(c *fiber.Ctx) error {
		return api.UpdateTodoStatus(c, dbs[dbName])
	})
	todoGroup.Delete("/:id", middleware.CheckLogin(store, dbs[dbName]), func(c *fiber.Ctx) error {
		return api.DeleteTodo(c, dbs[dbName])
	})
}
