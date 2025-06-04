package main

import (
	"fmt"
	"github.com/bytedance/sonic"
	db2 "github.com/crafty-ezhik/blog-api/db"
	"github.com/crafty-ezhik/blog-api/internal/auth"
	"github.com/crafty-ezhik/blog-api/internal/comment"
	"github.com/crafty-ezhik/blog-api/internal/config"
	"github.com/crafty-ezhik/blog-api/internal/post"
	"github.com/crafty-ezhik/blog-api/internal/routes"
	"github.com/crafty-ezhik/blog-api/internal/user"
	"github.com/crafty-ezhik/blog-api/pkg/jwt"
	"github.com/crafty-ezhik/blog-api/pkg/logger"
	"github.com/crafty-ezhik/blog-api/pkg/middleware"
	"github.com/crafty-ezhik/blog-api/pkg/validate"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/redis/go-redis/v9"
	"log"
)

/*
2. Добавить логирование
3. Подключить Swagger
5. Добавить CORS middleware
TODO:
	1. Добавить тесты(unit, integration, e2e)
	2. Добавить логирование в сервисный слой и в места, где  логика приложения
	3. Подключить Swagger
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

	db := db2.GetConnection(cfg)

	// Init redis
	logger.Log.Debug("Init Redis Client")
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Db,
	})

	// Init JWT
	jwtService := jwt.NewJWTService(jwt.NewRedisStorage(rdb))
	jwtAuth := jwt.NewJWT(jwtService, cfg.Auth.AccessTTL, cfg.Auth.RefreshTTL, cfg.Auth.SigningKey)

	// Init Validator
	logger.Log.Debug("Init validator")
	myValidator := validator.New()
	v := &validate.XValidator{
		Validator: myValidator,
	}

	// Repositories
	userRepo := user.NewUserRepository(db)
	postRepo := post.NewPostRepository(db)
	commentRepo := comment.NewCommentRepository(db)

	// Services
	userService := user.NewUserService(userRepo)
	authService := auth.NewAuthService(cfg, userRepo, jwtAuth)
	postService := post.NewPostService(postRepo)
	commentService := comment.NewCommentService(commentRepo, postRepo)

	// Handlers
	authHandler := auth.NewAuthHandler(userService, authService, v)
	userHandler := user.NewUserHandler(userService, postService, v)
	postHandler := post.NewPostHandler(postService, v)
	commentHandler := comment.NewCommentHandler(commentService, v)

	// Init Fiber App
	logger.Log.Debug("Init fiber")
	app := fiber.New(fiber.Config{
		JSONDecoder: sonic.Unmarshal,
		JSONEncoder: sonic.Marshal,
		// TODO: Подключить Swagger
	})

	// Middleware для логирование запросов
	app.Use(middleware.LogMiddleware())

	// CORS Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://gofiber.io, https://gofiber.net, http://localhost",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE",
	}))
	//
	routeDeps := routes.RouteDeps{
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		PostHandler:    postHandler,
		CommentHandler: commentHandler,
		JWT:            jwtAuth,
	}

	routes.SetupRoutes(app, routeDeps)

	// Start app
	logger.Log.Debug("Start app...")
	err = app.Listen(fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		panic("Error: " + err.Error())
	}
}
