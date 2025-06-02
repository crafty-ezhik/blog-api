package post

import (
	"github.com/crafty-ezhik/blog-api/internal/models"
	"github.com/crafty-ezhik/blog-api/pkg/middleware"
	"github.com/crafty-ezhik/blog-api/pkg/req"
	"github.com/crafty-ezhik/blog-api/pkg/validate"
	"github.com/gofiber/fiber/v2"
	"strconv"
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
	v           *validate.XValidator
}

func NewPostHandler(postService PostService, validator *validate.XValidator) *PostHandlerImpl {
	return &PostHandlerImpl{
		PostService: postService,
		v:           validator,
	}
}

func (h *PostHandlerImpl) GetAllPosts(c *fiber.Ctx) error {
	data, err := h.PostService.GetAllPosts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func (h *PostHandlerImpl) GetPostById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Post Id is invalid",
		})
	}
	data, err := h.PostService.GetPostById(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func (h *PostHandlerImpl) CreatePost(c *fiber.Ctx) error {
	body, err := req.HandleBody[CreateRequest](c, h.v)
	if err != nil {
		return nil
	}

	authorID := c.Locals(middleware.UserIDKey).(uint)

	newPost := &models.Post{
		Title:    body.Title,
		Text:     body.Text,
		AuthorID: authorID,
	}
	err = h.PostService.CreatePost(newPost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    newPost,
	})
}

func (h *PostHandlerImpl) UpdatePost(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (h *PostHandlerImpl) DeletePost(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}
