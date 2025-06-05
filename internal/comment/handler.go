package comment

import (
	"errors"
	"github.com/crafty-ezhik/blog-api/pkg/logger"
	"github.com/crafty-ezhik/blog-api/pkg/middleware"
	"github.com/crafty-ezhik/blog-api/pkg/req"
	"github.com/crafty-ezhik/blog-api/pkg/validate"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
)

type CommentHandler interface {
	GetMyComment(c *fiber.Ctx) error
	GetUserComments(c *fiber.Ctx) error
	GetAllCommentsPost(c *fiber.Ctx) error
	CreateComments(c *fiber.Ctx) error
	UpdateComment(c *fiber.Ctx) error
	DeleteComment(c *fiber.Ctx) error
}

type CommentHandlerImpl struct {
	CommentService CommentService
	v              *validate.XValidator
}

func NewCommentHandler(commentService CommentService, validator *validate.XValidator) *CommentHandlerImpl {
	logger.Log.Debug("Init comment handler")
	return &CommentHandlerImpl{
		CommentService: commentService,
		v:              validator,
	}
}

func (h *CommentHandlerImpl) GetAllCommentsPost(c *fiber.Ctx) error {
	postID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Post ID must be an integer",
		})
	}
	return h.getComments(c, uint(postID), 0)
}

func (h *CommentHandlerImpl) GetMyComment(c *fiber.Ctx) error {
	userID := c.Locals(middleware.UserIDKey).(uint)

	postID, err := strconv.Atoi(c.Params("postId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Post ID must be an integer",
		})
	}
	return h.getComments(c, uint(postID), userID)
}

func (h *CommentHandlerImpl) GetUserComments(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "User ID must be an integer",
		})
	}
	postID, err := strconv.Atoi(c.Params("postId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Post ID must be an integer",
		})
	}
	return h.getComments(c, uint(postID), uint(userID))
}

func (h *CommentHandlerImpl) CreateComments(c *fiber.Ctx) error {
	userID := c.Locals(middleware.UserIDKey).(uint)

	postID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Post ID must be an integer",
		})
	}
	body, err := req.HandleBody[CreateCommentRequest](c, h.v)
	if err != nil {
		return nil
	}
	err = h.CommentService.CreateCommentByPostID(uint(postID), userID, body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Something went wrong",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Comment created successfully",
	})
}

func (h *CommentHandlerImpl) UpdateComment(c *fiber.Ctx) error {
	userID := c.Locals(middleware.UserIDKey).(uint)

	commentID, err := strconv.Atoi(c.Params("commentId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Comment ID must be an integer",
		})
	}
	postID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Post ID must be an integer",
		})
	}

	body, err := req.HandleBody[UpdateCommentRequest](c, h.v)
	if err != nil {
		return nil
	}

	err = h.CommentService.UpdateComment(uint(commentID), uint(postID), userID, body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Comment updated successfully",
	})
}

func (h *CommentHandlerImpl) DeleteComment(c *fiber.Ctx) error {
	userID := c.Locals(middleware.UserIDKey).(uint)

	commentID, err := strconv.Atoi(c.Params("commentId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Comment ID must be an integer",
		})
	}
	postID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Post ID must be an integer",
		})
	}

	err = h.CommentService.DeleteComment(uint(commentID), uint(postID), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"success": true,
		"message": "Comment deleted successfully",
	})
}

func (h *CommentHandlerImpl) getComments(c *fiber.Ctx, postID, userID uint) error {
	data, err := h.CommentService.GetCommentsByPostID(postID, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Comment not found",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}
