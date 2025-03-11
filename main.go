package main

import (
	"os"

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
		// CookieHTTPOnly: true,
	})
	// management database connect
	dbs = make(map[string]*gorm.DB)
	// set frontendFolder
	frontendFolder string = "./public"
)

func main() {
	// logrus settings
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)

	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf(".env file load error: %v", err)
	}
	if os.Getenv("ENV") == "development" {
		// two database
		for i := 0; i <= 2; i++ {
			dbName, db := model.InitDsn(i)
			dbs[dbName] = db
			model.Migration(dbName, dbs[dbName])
		}
		frontendFolder = "./dist"
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{AllowOrigins: "http://localhost:8080",
		AllowMethods: "POST, PATCH"}))
	// static folder
	app.Static("/", frontendFolder)
	// middleware
	app.Use(middleware.SessionHandler(store))
	// route
	app.Post("/api/registerHandler", func(c *fiber.Ctx) error {
		return api.RegisterHandler(c, dbs[os.Getenv("DB_NAME")])
	})
	app.Post("/api/loginHandler", func(c *fiber.Ctx) error {
		return loginHandler(c, dbs[os.Getenv("DB_NAME")])
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
	/* --------------------------------- */

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
	/* --------------------------------- */

	app.Post("/api/verify", verifyHandler)

	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile(frontendFolder + "/index.html")
	})

	logrus.Fatal(app.Listen(":3000"))
}

func loginHandler(c *fiber.Ctx, db *gorm.DB) error {
	var data api.LoginData

	if err := c.BodyParser(&data); err != nil {
		logrus.Fatalf("%v", err)
	}
	logrus.Debugf("login page form data: %v", data)

	err := api.Login(c, store, &data, db)
	if err != nil {
		logrus.Infof("%v", err)
		// 401
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": err.Error()})
	}

	return c.JSON(fiber.Map{
		"msg": "登入成功"})
}

func verifyHandler(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
