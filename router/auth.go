package router

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"seaotterms.com-backend/api"
	"seaotterms.com-backend/middleware"
)

// this router is use to check identity for front-end routes
func AuthRouter(routerGroup fiber.Router, store *session.Store, dbs map[string]*gorm.DB) {
	authGroup := routerGroup.Group("/auth")
	dbName := os.Getenv("DB_NAME")

	authGroup.Get("/", middleware.CheckLogin(store, dbs[dbName]), func(c *fiber.Ctx) error {
		return api.AuthLogin(c, store)
	})
	// check if you are the website owner
	authGroup.Get("/root", middleware.CheckOwner(store, dbs[dbName]), func(c *fiber.Ctx) error {
		return api.AuthLogin(c, store)
	})
	// authGroup.Get("/specific", middleware.CheckLogin(store, dbs[dbName]), func(c *fiber.Ctx) error {
	// 	return api.AuthLogin(c, store)
	// })
}
