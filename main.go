package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"seaotterms.com-backend/crud"
	"seaotterms.com-backend/model"
)

// type FormData struct {
// 	username      string
// 	email         string
// 	password      string
// 	checkPassword string
// }

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

	app.Post("/registerHandler", func(c *fiber.Ctx) error {
		var data map[string]interface{}

		if err := c.BodyParser(&data); err != nil {
			log.Fatalf("%v", err)
		}
		fmt.Printf("Received data: %+v\n", data)

		// database handler
		crud.Register()

		return c.SendString("test")
	})
	app.Post("/loginHandler", func(c *fiber.Ctx) error {
		return c.SendString("test")
	})

	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile("./public/index.html")
	})

	log.Fatal(app.Listen(":3000"))
}
