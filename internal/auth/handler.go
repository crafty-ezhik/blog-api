package auth

import (
	"github.com/crafty-ezhik/blog-api/internal/user"
	"github.com/crafty-ezhik/blog-api/pkg/jwt"
	"github.com/crafty-ezhik/blog-api/pkg/req"
	"github.com/crafty-ezhik/blog-api/pkg/validate"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	LogoutAll(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
}

type AuthHandlerImpl struct {
	UserService user.UserService
	AuthService AuthService
	v           *validate.XValidator
}

func NewAuthHandler(userService user.UserService, authService AuthService, validator *validate.XValidator) *AuthHandlerImpl {
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
	responseData, err := h.AuthService.Login(body)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(responseData)
}

func (h *AuthHandlerImpl) Logout(c *fiber.Ctx) error {

	return nil
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
			"error":   err,
		})
	}

	response := RegisterResponse{
		Success: true,
		Message: "You have successfully registered",
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}
func (h *AuthHandlerImpl) LogoutAll(c *fiber.Ctx) error {
	return nil
}

func (h *AuthHandlerImpl) Refresh(c *fiber.Ctx) error {
	resfreshToken := c.Cookies("refresh_token")
	if resfreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"err":     jwt.ErrTokenNotFound,
		})
	}
}
