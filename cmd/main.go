package main

import (
	"fmt"
	db2 "github.com/crafty-ezhik/blog-api/db"
	"github.com/crafty-ezhik/blog-api/internal/auth"
	"github.com/crafty-ezhik/blog-api/internal/config"
	"github.com/crafty-ezhik/blog-api/internal/routes"
	"github.com/crafty-ezhik/blog-api/internal/user"
	"github.com/crafty-ezhik/blog-api/pkg/validate"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	// Init configuration
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatal(err)
	}

	// Get database connection
	db := db2.GetConnection(cfg)

	// Init Validator
	myValidator := validator.New()
	v := &validate.XValidator{
		Validator: myValidator,
	}

	// Init Fiber App
	app := fiber.New(fiber.Config{
		// TODO: Добавить логирование через ZAP
		// TODO: Подключить Swagger
		// TODO: Поменять JSON Decoder/Encoder с encoding/json на bytedance/sonic
	})

	// Repositories
	userRepo := user.NewUserRepository(db)

	// Services
	userService := user.NewUserService(userRepo)
	authService := auth.NewAuthService(cfg, userRepo)

	// Handlers
	authHandler := auth.NewAuthHandler(userService, authService, v)

	//
	routeDeps := routes.RouteDeps{
		AuthHandler: authHandler,
	}
	routes.SetupRoutes(app, routeDeps)

	// Start app
	err = app.Listen(fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		panic("Error: " + err.Error())
	}
}
