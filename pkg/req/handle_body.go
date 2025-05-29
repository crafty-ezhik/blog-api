package req

import (
	"github.com/crafty-ezhik/blog-api/pkg/res"
	"github.com/crafty-ezhik/blog-api/pkg/validate"
	"github.com/gofiber/fiber/v2"
)

func HandleBody[T any](c *fiber.Ctx, validator *validate.XValidator) (*T, error) {
	var body T
	if err := c.BodyParser(&body); err != nil {
		res.ErrorResponse(c, fiber.StatusBadRequest, "Failed to parse request body", err.Error())
		return nil, err
	}

	err := validator.Validate(body)
	if err != nil {
		res.ErrorResponse(c, fiber.StatusBadRequest, "Validation error", err.Error())
		return nil, err
	}
	return &body, nil
}
