package post

import (
	"github.com/gofiber/fiber/v2"
)

type PostHandler interface {
	GetAllPosts(c *fiber.Ctx) error
	GetPostById(c *fiber.Ctx) error
	CreatePost(c *fiber.Ctx) error
	UpdatePost(c *fiber.Ctx) error
	DeletePost(c *fiber.Ctx) error
}

type PostHandlerImpl struct {
	PostService PostService
}

func NewPostHandler(postService PostService) *PostHandlerImpl {
	return &PostHandlerImpl{
		PostService: postService,
	}
}
