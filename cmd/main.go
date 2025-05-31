package main

import (
	"fmt"
	"github.com/bytedance/sonic"
	db2 "github.com/crafty-ezhik/blog-api/db"
	"github.com/crafty-ezhik/blog-api/internal/auth"
	"github.com/crafty-ezhik/blog-api/internal/config"
	"github.com/crafty-ezhik/blog-api/internal/routes"
	"github.com/crafty-ezhik/blog-api/internal/user"
	"github.com/crafty-ezhik/blog-api/pkg/jwt"
	"github.com/crafty-ezhik/blog-api/pkg/validate"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"log"
)

/*
1. Сделать генерацию токенов, обновление/отзыв/версионирование refresh токена
2. Добавить логирование
3. Подключить Swagger
4. Поменять стандартный encoding/json на bytedance/sonic
5. Добавить CORS middleware
*/

func main() {
	// Init configuration
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatal(err)
	}

	// Get database connection
	db := db2.GetConnection(cfg)

	// Init redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Db,
	})

	// Init JWT
	jwtService := jwt.NewJWTService(jwt.NewRedisStorage(rdb))
	jwtAuth := jwt.NewJWT(jwtService, cfg.Auth.AccessTTL, cfg.Auth.RefreshTTL, cfg.Auth.SigningKey)

	// Init Validator
	myValidator := validator.New()
	v := &validate.XValidator{
		Validator: myValidator,
	}

	// Repositories
	userRepo := user.NewUserRepository(db)

	// Services
	userService := user.NewUserService(userRepo)
	authService := auth.NewAuthService(cfg, userRepo, jwtAuth)

	// Handlers
	authHandler := auth.NewAuthHandler(userService, authService, v)

	// Init Fiber App
	app := fiber.New(fiber.Config{
		JSONDecoder: sonic.Unmarshal,
		JSONEncoder: sonic.Marshal,
		// TODO: Добавить логирование через ZAP
		// TODO: Подключить Swagger
	})

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
