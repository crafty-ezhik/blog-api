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
	"github.com/crafty-ezhik/blog-api/pkg/logger"
	"github.com/crafty-ezhik/blog-api/pkg/validate"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"log"
)

/*
2. Добавить логирование
3. Подключить Swagger
5. Добавить CORS middleware
*/

func main() {
	// Init configuration
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatal(err)
	}

	// Init logger
	err = logger.InitLogger(cfg)
	if err != nil {
		panic(err)
	}
	// Get database connection
	logger.Log.Debug("Получение подключения к базе данных")
	db := db2.GetConnection(cfg)

	// Init redis
	logger.Log.Debug("Инициализация Redis Client")
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Db,
	})

	// Init JWT
	logger.Log.Debug("Инициализация модуля для JWT")
	jwtService := jwt.NewJWTService(jwt.NewRedisStorage(rdb))
	jwtAuth := jwt.NewJWT(jwtService, cfg.Auth.AccessTTL, cfg.Auth.RefreshTTL, cfg.Auth.SigningKey)

	// Init Validator
	logger.Log.Debug("Инициализация валидатора")
	myValidator := validator.New()
	v := &validate.XValidator{
		Validator: myValidator,
	}

	// Repositories
	logger.Log.Debug("Инициализация репозиториев")
	userRepo := user.NewUserRepository(db)

	// Services
	logger.Log.Debug("Инициализация сервисов")
	userService := user.NewUserService(userRepo)
	authService := auth.NewAuthService(cfg, userRepo, jwtAuth)

	// Handlers
	logger.Log.Debug("Инициализация хендлеров")
	authHandler := auth.NewAuthHandler(userService, authService, v)

	// Init Fiber App
	logger.Log.Debug("Инициализация fiber")
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
	logger.Log.Debug("Старт сервера")
	err = app.Listen(fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		panic("Error: " + err.Error())
	}
}
