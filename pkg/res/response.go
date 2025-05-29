package res

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

func SuccessResponse(c *fiber.Ctx, data interface{}) {
	c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Data:    data,
	})
}

func ErrorResponse(c *fiber.Ctx, code int, message string, details ...string) {
	errorDetails := ""
	if len(details) > 0 {
		errorDetails = details[0]
	}
	c.Status(code).JSON(Response{
		Success: false,
		Error: &ErrorInfo{
			Message: message,
			Code:    code,
			Details: errorDetails,
		},
	})
}
