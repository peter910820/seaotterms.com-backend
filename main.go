package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/session/v2"
	"github.com/joho/godotenv"

	"seaotterms.com-backend/internal/login"
	"seaotterms.com-backend/internal/model"
	"seaotterms.com-backend/internal/register"
)

var store = session.New()

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(".env file error: %v", err)
	}

	if os.Getenv("ENV") == "development" {
		model.Migration()
	}
	app := fiber.New()

	app.Use(cors.New(cors.Config{AllowOrigins: "http://localhost:8080",
		AllowMethods: "POST"}))

	app.Static("/", "./public")

	app.Post("/api/registerHandler", registerHandler)

	app.Post("/api/loginHandler", loginHandler)

	app.Get("*", func(c *fiber.Ctx) error {
		sess := store.Get(c)
		username := sess.Get("username")
		if username == nil {
			fmt.Println("error")
		}
		return c.SendFile("./public/index.html")
	})

	log.Fatal(app.Listen(":3000"))
}

func registerHandler(c *fiber.Ctx) error {
	var data register.RegisterData

	if err := c.BodyParser(&data); err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("Received data: %+v\n", data)

	err := register.Register(&data)
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
	var data login.LoginData

	if err := c.BodyParser(&data); err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("%v\n", data)

	err := login.Login(c, store, &data)
	if err != nil {
		log.Printf("%v\n", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": err.Error()})
	}

	return c.JSON(fiber.Map{
		"msg": "登入成功"})
}
