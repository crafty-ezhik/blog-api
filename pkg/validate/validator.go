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
		Tag          string
		Value        any
	}

	XValidator struct {
		validator *validator.Validate
	}

	GlobalErrorHandlerResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
)

func (v XValidator) validateData(data any) []ErrorResponse {
	var validationErrors []ErrorResponse

	errs := v.validator.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse

			elem.Tag = err.Tag()
			elem.Value = err.Value()
			elem.FailedFields = err.Field()
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}
	return validationErrors
}

func (v XValidator) Validate(data any) *fiber.Error {

	if errs := v.validateData(data); len(errs) > 0 && errs[0].Error {
		errMsg := make([]string, 0)

		for _, err := range errs {
			errMsg = append(errMsg, fmt.Sprintf(
				"[%s]: '%v' | Needs to implement '%s'",
				err.FailedFields,
				err.Value,
				err.Tag,
			))
		}
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: strings.Join(errMsg, " and "),
		}
	}
	return nil
}
