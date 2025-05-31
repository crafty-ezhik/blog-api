package routes

import (
	"github.com/crafty-ezhik/blog-api/internal/auth"
	"github.com/crafty-ezhik/blog-api/pkg/jwt"
	"github.com/crafty-ezhik/blog-api/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

type RouteDeps struct {
	AuthHandler auth.AuthHandler
	JWT         *jwt.JWT
}

func SetupRoutes(app *fiber.App, deps RouteDeps) {
	// Auth
	app.Route("/auth", func(router fiber.Router) {
		router.Post("/register", deps.AuthHandler.Register)
		router.Post("/login", deps.AuthHandler.Login)
		router.Post("/logout", middleware.AuthMiddleware(deps.JWT), deps.AuthHandler.Logout)   // TODO: Нужна проверка access токена
		router.Post("/refresh", middleware.AuthMiddleware(deps.JWT), deps.AuthHandler.Refresh) // TODO: Нужна проверка refresh токена
	})

	api := app.Group("/api")

	// Users
	api.Route("users", func(router fiber.Router) {
		router.Get("/me", middleware.AuthMiddleware(deps.JWT), pass)       // TODO: Нужна проверка access токена
		router.Get("/my/posts", middleware.AuthMiddleware(deps.JWT), pass) // Получение постов пользователя
		router.Get("/:id/posts", pass)                                     // Получение постов по id пользователя
		router.Get("/my/posts/:postId/comments", pass)                     // Получение всех своих комментариев к статье
		router.Get("/:id/posts/:postId/comments", pass)                    // Получение всех комментариев к статье по id пользователя
	})

	// Posts
	api.Route("posts", func(router fiber.Router) {
		router.Post("/", pass)      // Создание статьи
		router.Get("/", pass)       // Получение всех статей
		router.Get("/:id", pass)    // Получение конкретной статьи
		router.Put("/:id", pass)    // Обновление статьи
		router.Delete("/:id", pass) // Удаление статьи

		router.Get("/:id/comments", pass)               // Получение всех комментариев к статье
		router.Post("/:id/comments", pass)              // Создание комментария к посту
		router.Put("/:id/comments/:commentId", pass)    // Обновление комментария
		router.Delete("/:id/comments/:commentId", pass) // Удаление комментария
	})
}

func pass(c *fiber.Ctx) error {
	return c.Status(fiber.StatusTeapot).JSON(fiber.Map{
		"success": true,
	})
}
