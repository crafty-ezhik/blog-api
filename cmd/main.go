package main

import (
	"github.com/crafty-ezhik/blog-api/internal/config"
	"github.com/crafty-ezhik/blog-api/internal/routes"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{
		// TODO: Добавить логирование через ZAP
		// TODO: Подключить Swagger
		// TODO: Поменять JSON Decoder/Encoder с encoding/json на bytedance/sonic
	})

	routes.SetupRoutes(app)

	err = app.Listen(":3000")
	if err != nil {
		panic("Error: " + err.Error())
	}
}
