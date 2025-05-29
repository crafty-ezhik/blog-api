package auth

import (
	"github.com/crafty-ezhik/blog-api/internal/user"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
}

type AuthHandlerImpl struct {
	UserService user.UserService
}

func NewAuthHandler(userService user.UserService) *AuthHandlerImpl {
	return &AuthHandlerImpl{UserService: userService}
}

func (h *AuthHandlerImpl) Login(c *fiber.Ctx) error {

}
