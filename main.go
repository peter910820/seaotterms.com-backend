package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"seaotterms.com-backend/internal/crud"
	"seaotterms.com-backend/internal/model"
)

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
		log.Printf("Received data: %+v\n", data)

		// database handler
		err = crud.Register(&data)
		if err != nil {
			log.Printf("%v\n", err)
			return c.JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "User registered",
		})
	})
	app.Post("/loginHandler", func(c *fiber.Ctx) error {
		return c.SendString("test")
	})

	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile("./public/index.html")
	})

	log.Fatal(app.Listen(":3000"))
}
