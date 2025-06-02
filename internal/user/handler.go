package user

import (
	"github.com/crafty-ezhik/blog-api/internal/models"
	"github.com/crafty-ezhik/blog-api/pkg/middleware"
	"github.com/crafty-ezhik/blog-api/pkg/req"
	"github.com/crafty-ezhik/blog-api/pkg/validate"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type UserHandler interface {
	GetByID(c *fiber.Ctx) error
	GetMe(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error

	GetMyPostsByID(c *fiber.Ctx) error
	GetUserPostsByID(c *fiber.Ctx) error
	GetMyCommentsByPostID(c *fiber.Ctx) error
	GetUserCommentsByPostID(c *fiber.Ctx) error
}

type UserHandlerImpl struct {
	UserService UserService
	v           *validate.XValidator
}

func NewUserHandler(userService UserService, validator *validate.XValidator) *UserHandlerImpl {
	return &UserHandlerImpl{
		UserService: userService,
		v:           validator,
	}
}

// region:Операции с пользователем
func (h *UserHandlerImpl) GetByID(c *fiber.Ctx) error {
	// Если нет параметров пути, то выводится данные под которыми выполнен вход
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "id must be an integer",
		})
	}
	result, err := h.UserService.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "user not found",
		})
	}

	data := GetByIDResponse{
		Name:      result.Name,
		Age:       result.Age,
		Email:     result.Email,
		CreatedAt: result.CreatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func (h *UserHandlerImpl) GetMe(c *fiber.Ctx) error {
	ctxUserID, ok := c.Locals(middleware.UserIDKey).(uint)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "user id must be a uint",
		})
	}
	result, err := h.UserService.GetByID(ctxUserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "user not found",
		})
	}

	data := GetByIDResponse{
		Name:      result.Name,
		Age:       result.Age,
		Email:     result.Email,
		CreatedAt: result.CreatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func (h *UserHandlerImpl) Update(c *fiber.Ctx) error {
	body, err := req.HandleBody[UpdateUserRequest](c, h.v)
	if err != nil {
		return nil
	}

	ctxUserID, ok := c.Locals(middleware.UserIDKey).(uint)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "user id must be a uint",
		})
	}

	updated := &models.User{
		Name: body.Name,
		Age:  body.Age,
	}
	err = h.UserService.Update(ctxUserID, updated)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "user updated",
	})
}

func (h *UserHandlerImpl) Delete(c *fiber.Ctx) error {
	userId := c.Params("id")
	userID, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "id must be an integer",
		})
	}

	err = h.UserService.Delete(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"success": true,
		"message": "User deleted",
	})
}

// endregion

// region: Операции с пользователем, постами и комментариями
func (h *UserHandlerImpl) GetMyPostsByID(c *fiber.Ctx) error {
	userID := c.Locals(middleware.UserIDKey).(uint)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    userID,
	})
}

func (h *UserHandlerImpl) GetUserPostsByID(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "id must be an integer",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    userId,
	})
}

func (h *UserHandlerImpl) GetMyCommentsByPostID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (h *UserHandlerImpl) GetUserCommentsByPostID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

// endregion
