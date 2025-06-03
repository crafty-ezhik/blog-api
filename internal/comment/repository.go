package comment

import (
	"github.com/crafty-ezhik/blog-api/internal/models"
	"gorm.io/gorm"
)

type CommentRepository interface {
	FindCommentsByPostID(comment *models.Comment) ([]models.Comment, error)
	CreateCommentByPostID(comment *models.Comment) error
	UpdateCommentByCommentAndPostID(comment *models.Comment) error
	DeleteCommentByCommentAndPostID(comment *models.Comment) error
}

type CommentRepositoryImpl struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepositoryImpl {
	return &CommentRepositoryImpl{
		db: db,
	}
}

func (r *CommentRepositoryImpl) FindCommentsByPostID(comment *models.Comment) ([]models.Comment, error) {
	var comments []models.Comment
	result := r.db.Model(&models.Comment{}).Where(comment).Joins("Author").Joins("Post").Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}

func (r *CommentRepositoryImpl) CreateCommentByPostID(comment *models.Comment) error {
	result := r.db.Create(comment)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *CommentRepositoryImpl) UpdateCommentByCommentAndPostID(comment *models.Comment) error {
	result := r.db.Model(&comment).Where("post_id = ?", comment.PostID).Updates(comment)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (r *CommentRepositoryImpl) DeleteCommentByCommentAndPostID(comment *models.Comment) error {
	result := r.db.Model(&comment).Where("post_id = ?", comment.PostID).Delete(&comment)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
