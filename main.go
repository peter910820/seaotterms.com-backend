package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"

	"seaotterms.com-backend/internal/api"
	"seaotterms.com-backend/internal/model"
)

var store = session.New()

func main() {
	frontendFolder := "./public"

	err := godotenv.Load()
	if err != nil {
		log.Fatalf(".env file error: %v", err)
	}

	if os.Getenv("ENV") == "development" {
		model.Migration()
		frontendFolder = "./dist"
	}
	app := fiber.New()

	app.Use(cors.New(cors.Config{AllowOrigins: "http://localhost:8080",
		AllowMethods: "POST"}))

	app.Static("/", frontendFolder)

	app.Post("/api/registerHandler", registerHandler)
	app.Post("/api/loginHandler", loginHandler)
	app.Post("/api/create-article", createArticle)
	app.Post("/api/check-session", sessionHandler)

	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile(frontendFolder + "/index.html")
	})

	log.Fatal(app.Listen(":3000"))
}

func registerHandler(c *fiber.Ctx) error {
	var data api.RegisterData

	if err := c.BodyParser(&data); err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("Received data: %+v\n", data)

	err := api.Register(&data)
	if err != nil {
		log.Printf("%v\n", err)
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
		log.Fatalf("%v", err)
	}
	fmt.Printf("%v\n", data)

	err := api.Login(c, store, &data)
	if err != nil {
		log.Printf("%v\n", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": err.Error()})
	}

	return c.JSON(fiber.Map{
		"msg": "登入成功"})
}

func createArticle(c *fiber.Ctx) error {
	err := sessionHandler(c)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return c.SendStatus(fiber.StatusOK)
}

func sessionHandler(c *fiber.Ctx) error {
	log.Println("check session")
	err := api.CheckSession(c, store)
	if err != nil {
		log.Printf("%v\n", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": err.Error()})
	}
	return c.SendStatus(fiber.StatusOK)
}
