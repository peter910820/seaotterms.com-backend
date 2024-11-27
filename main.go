package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
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
