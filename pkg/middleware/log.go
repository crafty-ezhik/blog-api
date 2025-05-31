package middleware

import (
	"github.com/crafty-ezhik/blog-api/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"time"
)

func LogMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		c.Set("X-Request-Time", start.Format(time.RFC3339))

		logger.Log.Info("Incoming request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
		)

		err := c.Next()

		latency := time.Since(start)
		status := c.Response().StatusCode()

		fields := []zap.Field{
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", status),
			zap.Duration("latency", latency),
		}

		switch {
		case status >= 500:
			logger.Log.Error("Server error", fields...)
		case status >= 400:
			logger.Log.Warn("Client error", fields...)
		default:
			logger.Log.Info("Request processed", fields...)
		}

		return err
	}
}
