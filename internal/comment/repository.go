package comment

import (
	"github.com/crafty-ezhik/blog-api/internal/models"
	"gorm.io/gorm"
)

type CommentRepository interface {
	FindCommentsByPostID(postID uint) ([]models.Comment, error)
	CreateCommentByPostID(comment models.Comment) error
	UpdateCommentByCommentAndPostID(comment models.Comment) error
	DeleteCommentByCommentAndPostID(comment models.Comment) error
}

type CommentRepositoryImpl struct {
	db *gorm.DB
}
