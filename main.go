package main

import (
	"github.com/gemm123/vkrf-ticket/config"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v\n", err)
	}

	db := config.InitConnPool()
	defer db.Close()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	//api := app.Group("/api")
	//
	//v1 := api.Group("/v1")

	app.Listen(":3001")
}
