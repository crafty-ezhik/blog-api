package post

import (
	"github.com/crafty-ezhik/blog-api/internal/models"
	"gorm.io/gorm"
)

type PostRepository interface {
	FindALL() []models.Post
	FindByID(postID uint) (models.Post, error)
	Create(post models.Post) error
	Update(post models.Post) error
	Delete(post models.Post) error
}

type PostRepositoryImpl struct {
	db *gorm.DB
}
