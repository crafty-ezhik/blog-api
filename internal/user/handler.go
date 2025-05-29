package user

import "github.com/gofiber/fiber/v2"

type AuthHandler interface {
}

type AuthHandlerImpl struct {
	UserService UserService
}

func NewAuthHandler(userService UserService) *AuthHandlerImpl {
	return &AuthHandlerImpl{UserService: userService}
}

func (h *AuthHandlerImpl) Login(c *fiber.Ctx) error {
	
}
