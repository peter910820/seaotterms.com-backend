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
	"seaotterms.com-backend/internal/router"
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

	// route group
	apiGroup := app.Group("/api") // main api route group

	router.ArticleRouter(apiGroup, store, dbs)
	router.GalgameRouter(apiGroup, store, dbs)
	router.GalgameBrandRouter(apiGroup, store, dbs)
	router.UserRouter(apiGroup, store, dbs)
	router.TodoRouter(apiGroup, store, dbs)
	router.TodoTopicRouter(apiGroup, store, dbs)
	router.TagRouter(apiGroup, store, dbs)

	// route
	/* --------------------------------- */
	// old route
	apiGroup.Post("/registerHandler", func(c *fiber.Ctx) error {
		return api.RegisterHandler(c, dbs[os.Getenv("DB_NAME3")])
	})
	apiGroup.Post("/loginHandler", func(c *fiber.Ctx) error {
		return api.Login(c, store, dbs[os.Getenv("DB_NAME3")])
	})
	/* --------------------------------- */
	// verify identity
	apiGroup.Post("/verify", middleware.CheckLogin(store, dbs[os.Getenv("DB_NAME3")]), func(c *fiber.Ctx) error {
		return api.Verify(c, store)
	})

	/* --------------------------------- */
	// match all routes
	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile(frontendFolder + "/index.html")
	})

	logrus.Fatal(app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
