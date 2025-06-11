package integration

import (
	"fmt"
	"github.com/bytedance/sonic"
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
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupApp(testDB *gorm.DB) *fiber.App {
	cfg, err := config.LoadConfig("/configs")
	if err != nil {
		log.Errorf("Error loading configs: %v", err)
	}

	err = logger.InitLogger(cfg)
	if err != nil {
		panic(err)
	}

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
	userRepo := user.NewUserRepository(testDB)
	postRepo := post.NewPostRepository(testDB)
	commentRepo := comment.NewCommentRepository(testDB)

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
	app := fiber.New(fiber.Config{
		JSONDecoder: sonic.Unmarshal,
		JSONEncoder: sonic.Marshal,
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
	return app
}
