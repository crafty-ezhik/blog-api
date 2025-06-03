package routes

import (
	"github.com/crafty-ezhik/blog-api/internal/auth"
	"github.com/crafty-ezhik/blog-api/internal/comment"
	"github.com/crafty-ezhik/blog-api/internal/post"
	"github.com/crafty-ezhik/blog-api/internal/user"
	"github.com/crafty-ezhik/blog-api/pkg/jwt"
	"github.com/crafty-ezhik/blog-api/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

type RouteDeps struct {
	AuthHandler    auth.AuthHandler
	UserHandler    user.UserHandler
	PostHandler    post.PostHandler
	CommentHandler comment.CommentHandler
	JWT            *jwt.JWT
}

func SetupRoutes(app *fiber.App, deps RouteDeps) {
	// Auth
	app.Route("/auth", func(router fiber.Router) {
		router.Post("/register", deps.AuthHandler.Register)
		router.Post("/login", deps.AuthHandler.Login)
		router.Post("/logout", middleware.AuthMiddleware(deps.JWT), deps.AuthHandler.Logout)
		router.Post("/refresh", middleware.AuthMiddleware(deps.JWT), deps.AuthHandler.Refresh)
	})

	api := app.Group("/api", middleware.AuthMiddleware(deps.JWT))

	// Users
	api.Route("users", func(router fiber.Router) {
		router.Get("/me", deps.UserHandler.GetMe)
		router.Get("/:id", deps.UserHandler.GetByID)
		router.Patch("/me", deps.UserHandler.Update)
		router.Delete("/:id", deps.UserHandler.Delete)
		router.Get("/my/posts", deps.UserHandler.GetMyPosts)                           // Получение постов пользователя
		router.Get("/my/posts/:postId/comments", deps.CommentHandler.GetMyComment)     // Получение всех своих комментариев к статье
		router.Get("/:id/posts", deps.UserHandler.GetUserPostsByID)                    // Получение постов по id пользователя
		router.Get("/:id/posts/:postId/comments", deps.CommentHandler.GetUserComments) // Получение всех комментариев к статье по id пользователя
	})

	// Posts
	api.Route("posts", func(router fiber.Router) {
		router.Post("/", deps.PostHandler.CreatePost)      // Создание статьи
		router.Get("/", deps.PostHandler.GetAllPosts)      // Получение всех статей
		router.Get("/:id", deps.PostHandler.GetPostById)   // Получение конкретной статьи
		router.Patch("/:id", deps.PostHandler.UpdatePost)  // Обновление статьи
		router.Delete("/:id", deps.PostHandler.DeletePost) // Удаление статьи

		router.Get("/:id/comments", deps.CommentHandler.GetAllCommentsPost)          // Получение всех комментариев к статье
		router.Post("/:id/comments", deps.CommentHandler.CreateComments)             // Создание комментария к посту
		router.Patch("/:id/comments/:commentId", deps.CommentHandler.UpdateComment)  // Обновление комментария
		router.Delete("/:id/comments/:commentId", deps.CommentHandler.DeleteComment) // Удаление комментария
	})
}
