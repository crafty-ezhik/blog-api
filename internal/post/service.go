package post

import (
	"github.com/crafty-ezhik/blog-api/internal/models"
	"github.com/crafty-ezhik/blog-api/pkg/logger"
)
//go:generate mockgen -source=service.go -destination=mock/post_service_mock.go


type PostService interface {
	GetAllPosts() ([]models.Post, error)
	GetPostById(postID uint) (*models.Post, error)
	GetPostsByAuthorID(authorID uint) ([]models.Post, error)
	CreatePost(post *models.Post) error
	UpdatePost(postID uint, updatedFields *models.Post) error
	DeletePost(postID uint) error
}

type PostServiceImpl struct {
	PostRepo PostRepository
}

func NewPostService(postRepo PostRepository) *PostServiceImpl {
	logger.Log.Debug("Init post service")
	return &PostServiceImpl{
		PostRepo: postRepo,
	}
}

func (s *PostServiceImpl) GetAllPosts() ([]models.Post, error) {
	return s.PostRepo.FindALL()
}

func (s *PostServiceImpl) GetPostById(postID uint) (*models.Post, error) {
	return s.PostRepo.FindByID(postID)
}

func (s *PostServiceImpl) GetPostsByAuthorID(authorID uint) ([]models.Post, error) {
	return s.PostRepo.FindByUserID(authorID)
}

func (s *PostServiceImpl) CreatePost(post *models.Post) error {
	return s.PostRepo.Create(post)
}

func (s *PostServiceImpl) UpdatePost(postID uint, updatedFields *models.Post) error {
	return s.PostRepo.Update(postID, updatedFields)
}

func (s *PostServiceImpl) DeletePost(postID uint) error {
	return s.PostRepo.Delete(postID)
}
