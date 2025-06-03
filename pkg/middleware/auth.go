package middleware

import (
	"github.com/crafty-ezhik/blog-api/pkg/jwt"
	"github.com/crafty-ezhik/blog-api/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type KeyType string

var UserIDKey KeyType = "user_id"

func AuthMiddleware(jwt *jwt.JWT) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger.Log.Info("Check access token")
		logger.Log.Debug("Check Authorization header")
		rawToken := c.Get("Authorization")
		if rawToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"err":     "Unauthorized",
				"details": "Token is empty",
			})
		}

		logger.Log.Debug("Check Bearer prefix")
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
		c.Locals(UserIDKey, tokenData.UserId)
		logger.Log.Info("Token verification completed successfully")
		return c.Next()
	}
}
