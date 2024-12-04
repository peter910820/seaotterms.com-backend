package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"

	"seaotterms.com-backend/internal/api"
	"seaotterms.com-backend/internal/middleware"
	"seaotterms.com-backend/internal/model"
)

// init store(session)
var store = session.New(session.Config{
	// CookieHTTPOnly: true,
})

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
	// set frontendFolder
	frontendFolder := "./public"
	if os.Getenv("ENV") == "development" {
		model.Migration()
		frontendFolder = "./dist"
	}
	app := fiber.New()

	app.Use(cors.New(cors.Config{AllowOrigins: "http://localhost:8080",
		AllowMethods: "POST"}))
	// middleware
	app.Use(middleware.SessionHandler(store))

	app.Static("/", frontendFolder)

	app.Post("/api/registerHandler", registerHandler)
	app.Post("/api/loginHandler", loginHandler)

	app.Post("/api/create-article", createArticle)
	app.Post("/api/articles", api.GetArticle)

	app.Post("/api/verify", verifyHandler)

	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile(frontendFolder + "/index.html")
	})

	logrus.Fatal(app.Listen(":3000"))
}

func registerHandler(c *fiber.Ctx) error {
	var data api.RegisterData

	if err := c.BodyParser(&data); err != nil {
		logrus.Fatalf("%v", err)
	}
	logrus.Debugf("Received data: %+v", data)

	err := api.Register(&data)
	if err != nil {
		logrus.Infof("%v", err)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"msg": "註冊成功",
	})
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

func createArticle(c *fiber.Ctx) error {
	var data api.ArticleData

	if err := c.BodyParser(&data); err != nil {
		logrus.Fatalf("%v", err)
	}
	err := api.CreateArticle(&data)
	if err != nil {
		// 500
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func verifyHandler(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
