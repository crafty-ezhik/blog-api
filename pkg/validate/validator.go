package validate

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type (
	ErrorResponse struct {
		Error        bool
		FailedFields string
	}

	XValidator struct {
		Validator *validator.Validate
	}

	GlobalErrorHandlerResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
)

func (v XValidator) Validate(data any) *fiber.Error {

	if errs := v.validateData(data); len(errs) > 0 && errs[0].Error {
		errMsg := make([]string, 0)
		for _, err := range errs {
			errMsg = append(errMsg, fmt.Sprintf(
				"[%s]",
				err.FailedFields,
			))
		}
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid field or its absence: " + strings.Join(errMsg, " and "),
		}
	}
	return nil
}

func (v XValidator) validateData(data any) []ErrorResponse {
	var validationErrors []ErrorResponse

	errs := v.Validator.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse
			elem.FailedFields = err.Field()
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}
	return validationErrors
}
