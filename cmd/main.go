package main

import (
	"github.com/crafty-ezhik/blog-api/internal/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		// TODO: Добавить логирование через ZAP
		// TODO: Подключить Swagger
		// TODO: Поменять JSON Decoder/Encoder с encoding/json на bytedance/sonic
	})

	app.Get("/", hello)
	routes.SetupRoutes(app)

	err := app.Listen(":3000")
	if err != nil {
		panic("Error: " + err.Error())
	}
}

func hello(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello, World!",
	})
}
