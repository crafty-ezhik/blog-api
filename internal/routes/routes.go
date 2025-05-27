package routes

import "github.com/gofiber/fiber/v2"

type Routes struct {
}

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Auth
	api.Route("/auth", func(router fiber.Router) {
		router.Post("/register", pass)
		router.Post("/login", pass)
		router.Post("/logout", pass)     // TODO: Нужна проверка access токена
		router.Post("/logout-all", pass) // TODO: Нужна проверка access токена
		router.Post("/refresh", pass)    // TODO: Нужна проверка refresh токена
	})

	// Users
	api.Route("users", func(router fiber.Router) {
		router.Get("/me", pass)                         // TODO: Нужна проверка access токена
		router.Get("/my/posts", pass)                   // Получение постов пользователя
		router.Get("/:id/posts", pass)                  // Получение постов по id пользователя
		router.Get("/my/posts/:postId/comments", pass)  // Получение всех своих комментариев к статье
		router.Get("/:id/posts/:postId/comments", pass) // Получение всех комментариев к статье по id пользователя
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
