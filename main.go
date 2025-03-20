package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/api"
	"seaotterms.com-backend/internal/middleware"
	"seaotterms.com-backend/internal/model"
)

var (
	// init store(session)
	store = session.New(session.Config{
		Expiration: 24 * time.Hour,
		// CookieHTTPOnly: true,
	})
	// management database connect
	dbs = make(map[string]*gorm.DB)
	// set frontendFolder
	frontendFolder string = "./dist"
)

func init() {
	// init logrus settings
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)
	// init env file
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf(".env file load error: %v", err)
	}
}

func main() {
	// init migration
	for i := 0; i <= 2; i++ {
		dbName, db := model.InitDsn(i)
		dbs[dbName] = db
		model.Migration(dbName, dbs[dbName])
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{AllowOrigins: "http://localhost:8080",
		AllowMethods: "POST, PATCH"}))
	// static folder
	app.Static("/", frontendFolder)
	// middleware
	app.Use(middleware.SessionHandler(store, dbs[os.Getenv("DB_NAME3")]))
	app.Use(middleware.AuthenticationManagementHandler(store, dbs[os.Getenv("DB_NAME3")]))

	// route
	/* --------------------------------- */
	// old route
	app.Post("/api/registerHandler", func(c *fiber.Ctx) error {
		return api.RegisterHandler(c, dbs[os.Getenv("DB_NAME3")])
	})
	app.Post("/api/loginHandler", func(c *fiber.Ctx) error {
		return api.Login(c, store, dbs[os.Getenv("DB_NAME3")])
	})
	app.Post("/api/create-article", func(c *fiber.Ctx) error {
		return api.ArticleHandler(c, dbs[os.Getenv("DB_NAME")])
	})
	app.Post("/api/articles", func(c *fiber.Ctx) error {
		return api.GetArticle(c, dbs[os.Getenv("DB_NAME")])
	})
	app.Post("/api/articles/:articleID", func(c *fiber.Ctx) error {
		return api.GetSingleArticle(c, dbs[os.Getenv("DB_NAME")])
	})
	app.Post("/api/tags", func(c *fiber.Ctx) error {
		return api.GetTags(c, dbs[os.Getenv("DB_NAME")])
	})
	app.Post("/api/tags/:tagName", func(c *fiber.Ctx) error {
		return api.GetTag(c, dbs[os.Getenv("DB_NAME")])
	})
	/* --------------------------------- */
	// new route
	app.Get("/api/galgame/s/:name", func(c *fiber.Ctx) error {
		return api.QueryGalgame(c, dbs[os.Getenv("DB_NAME2")])
	})
	app.Get("/api/galgame/:brand", func(c *fiber.Ctx) error {
		return api.QueryGalgameByBrand(c, dbs[os.Getenv("DB_NAME2")])
	})
	app.Patch("/api/galgame/develop/:name", func(c *fiber.Ctx) error {
		return api.UpdateGalgameDevelop(c, dbs[os.Getenv("DB_NAME2")])
	})
	app.Post("/api/galgame", func(c *fiber.Ctx) error {
		return api.InsertGalgame(c, dbs[os.Getenv("DB_NAME2")])
	})
	app.Get("/api/galgame-brand", func(c *fiber.Ctx) error {
		return api.QueryAllGalgameBrand(c, dbs[os.Getenv("DB_NAME2")])
	})
	app.Get("/api/galgame-brand/:brand", func(c *fiber.Ctx) error {
		return api.QueryGalgameBrand(c, dbs[os.Getenv("DB_NAME2")])
	})
	app.Post("/api/galgame-brand", func(c *fiber.Ctx) error {
		return api.InsertGalgameBrand(c, dbs[os.Getenv("DB_NAME2")])
	})
	app.Patch("/api/galgame-brand/:brand", func(c *fiber.Ctx) error {
		return api.UpdateGalgameBrand(c, dbs[os.Getenv("DB_NAME2")])
	})
	/* --------------------------------- */
	app.Patch("/api/users/:id", func(c *fiber.Ctx) error {
		return api.UpdateUser(c, dbs[os.Getenv("DB_NAME3")])
	})
	/* --------------------------------- */
	app.Get("/api/todos/:owner", func(c *fiber.Ctx) error {
		return api.QueryTodoByOwner(c, dbs[os.Getenv("DB_NAME3")])
	})
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		return api.InsertTodo(c, dbs[os.Getenv("DB_NAME3")])
	})
	app.Post("/api/switch_todo/:id", func(c *fiber.Ctx) error {
		return api.SwitchTodo(c, dbs[os.Getenv("DB_NAME3")])
	})
	/* --------------------------------- */
	app.Get("/api/todo_topics", func(c *fiber.Ctx) error {
		return api.QueryTodoTopic(c, dbs[os.Getenv("DB_NAME3")])
	})
	app.Post("/api/todo_topics", func(c *fiber.Ctx) error {
		return api.InsertTodoTopic(c, dbs[os.Getenv("DB_NAME3")])
	})
	/* --------------------------------- */
	// verify identity
	app.Post("/api/verify", func(c *fiber.Ctx) error {
		return api.Verify(c, store)
	})

	/* --------------------------------- */
	// match all routes
	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile(frontendFolder + "/index.html")
	})

	logrus.Fatal(app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
