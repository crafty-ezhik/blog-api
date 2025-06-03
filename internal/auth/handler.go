package auth

import (
	"github.com/crafty-ezhik/blog-api/internal/user"
	"github.com/crafty-ezhik/blog-api/pkg/jwt"
	"github.com/crafty-ezhik/blog-api/pkg/logger"
	"github.com/crafty-ezhik/blog-api/pkg/req"
	"github.com/crafty-ezhik/blog-api/pkg/validate"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
}

type AuthHandlerImpl struct {
	UserService user.UserService
	AuthService AuthService
	v           *validate.XValidator
}

func NewAuthHandler(userService user.UserService, authService AuthService, validator *validate.XValidator) *AuthHandlerImpl {
	logger.Log.Debug("Init auth handler")
	return &AuthHandlerImpl{
		UserService: userService,
		AuthService: authService,
		v:           validator,
	}
}

func (h *AuthHandlerImpl) Login(c *fiber.Ctx) error {
	body, err := req.HandleBody[LoginRequest](c, h.v)
	if err != nil {
		return nil
	}
	responseData, cookie, err := h.AuthService.Login(body)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Установка refresh в куки
	c.Cookie(cookie)
	return c.Status(fiber.StatusOK).JSON(responseData)
}

func (h *AuthHandlerImpl) Logout(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"err":     jwt.ErrInBlackList,
		})
	}

	cookie, err := h.AuthService.Logout(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"err":     err.Error(),
		})
	}

	// Удаляем refresh из cookie
	c.Cookie(cookie)

	// Отдаем пользователю ответ, что он успешно вышел из системы
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "logged out",
	})
}

func (h *AuthHandlerImpl) Register(c *fiber.Ctx) error {
	body, err := req.HandleBody[RegisterRequest](c, h.v)
	if err != nil {
		return nil
	}
	ok, err := h.AuthService.Register(body)
	if err != nil || !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "User already exists",
		})
	}

	response := RegisterResponse{
		Success: true,
		Message: "You have successfully registered",
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *AuthHandlerImpl) Refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"err":     jwt.ErrInBlackList,
		})
	}
	tokens, cookie, err := h.AuthService.Refresh(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"err":     err.Error(),
		})
	}
	// Передаем refresh токен в куки
	c.Cookie(cookie)

	// Отдаем новый access токен пользователю
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": tokens.AccessToken,
	})
}
