package post

import (
	"github.com/crafty-ezhik/blog-api/internal/models"
	"gorm.io/gorm"
)

type PostRepository interface {
	FindALL() ([]models.Post, error)
	FindByID(postID uint) (*models.Post, error)
	Create(post *models.Post) error
	Update(postID uint, updatedFields *models.Post) error
	Delete(postID uint) error
}

type PostRepositoryImpl struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepositoryImpl {
	return &PostRepositoryImpl{
		db: db,
	}
}

func (repo *PostRepositoryImpl) FindALL() ([]models.Post, error) {
	var posts []models.Post
	result := repo.db.Find(&posts)
	return posts, result.Error
}

func (repo *PostRepositoryImpl) FindByID(postID uint) (*models.Post, error) {
	var post models.Post
	result := repo.db.First(&post, postID)
	return &post, result.Error
}

func (repo *PostRepositoryImpl) Create(post *models.Post) error {
	return repo.db.Create(post).Error
}

func (repo *PostRepositoryImpl) Update(postID uint, updatedFields *models.Post) error {
	return repo.db.Model(models.Post{ID: postID}).Updates(updatedFields).Error
}

func (repo *PostRepositoryImpl) Delete(postID uint) error {
	return repo.db.Delete(&models.Post{ID: postID}).Error
}
