package middleware

import (
	"github.com/crafty-ezhik/blog-api/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func AuthMiddleware(jwt *jwt.JWT) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rawToken := c.Get("Authorization")
		if rawToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"err":     "Unauthorized",
				"details": "Token is empty",
			})
		}

		if !strings.HasPrefix(rawToken, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"err":     "Unauthorized",
				"details": "Token does not start with Bearer",
			})
		}
		tokenString := strings.TrimPrefix(rawToken, "Bearer ")
		tokenData, err := jwt.VerifyToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"err":     "Unauthorized",
				"details": err.Error(),
			})
		}
		c.Locals("user_id", tokenData.UserId)
		return c.Next()
	}
}
