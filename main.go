package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"seaotterms.com-backend/model"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(".env file error: %v", err)
	}

	if os.Getenv("ENV") == "development" {
		model.Migration()
		fmt.Println("Successfully")
	}

	app := fiber.New()
	app.Static("/", "./public")

	app.Post("/loginHandler", func(c *fiber.Ctx) error {
		return c.SendString("test")
	})

	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile("./public/index.html")
	})

	log.Fatal(app.Listen(":3000"))
}
