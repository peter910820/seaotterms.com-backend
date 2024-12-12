package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"seaotterms.com-backend/internal/api"
	"seaotterms.com-backend/internal/middleware"
	"seaotterms.com-backend/internal/model"
)

var (
	// init store(session)
	store = session.New(session.Config{
		// CookieHTTPOnly: true,
	})
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
		model.Migration()
		frontendFolder = "./dist"
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{AllowOrigins: "http://localhost:8080",
		AllowMethods: "POST"}))
	// static folder
	app.Static("/", frontendFolder)
	// middleware
	app.Use(middleware.SessionHandler(store))
	// route
	app.Post("/api/registerHandler", api.RegisterHandler)
	app.Post("/api/loginHandler", loginHandler)

	app.Post("/api/create-article", api.CreateArticle)
	app.Post("/api/articles", api.GetArticle)
	app.Post("/api/articles/:articleID", api.GetSingleArticle)
	app.Post("/api/tags", api.GetTags)
	app.Post("/api/tags/:tagName", api.GetTag)

	app.Post("/api/verify", verifyHandler)

	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile(frontendFolder + "/index.html")
	})

	logrus.Fatal(app.Listen(":3000"))
}

func loginHandler(c *fiber.Ctx) error {
	var data api.LoginData

	if err := c.BodyParser(&data); err != nil {
		logrus.Fatalf("%v", err)
	}
	logrus.Debugf("login page form data: %v", data)

	err := api.Login(c, store, &data)
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
