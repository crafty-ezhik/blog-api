package comment

import (
	"errors"
	"github.com/crafty-ezhik/blog-api/internal/models"
	"github.com/crafty-ezhik/blog-api/internal/post"
	"github.com/crafty-ezhik/blog-api/pkg/logger"
	"gorm.io/gorm"
)

//go:generate mockgen -source=service.go -destination=mocks/comment_service_mock.go

type CommentService interface {
	GetCommentsByPostID(postID, userID uint) (*GetCommentsResponse, error)
	CreateCommentByPostID(postID, authorID uint, comment *CreateCommentRequest) error
	UpdateComment(commentID, PostID, userID uint, updatedFields *UpdateCommentRequest) error
	DeleteComment(commentID, PostID, userID uint) error
}

type CommentServiceImpl struct {
	CommentRepo CommentRepository
	PostRepo    post.PostRepository
}

func NewCommentService(commentRepo CommentRepository, postRepo post.PostRepository) *CommentServiceImpl {
	logger.Log.Debug("Init comment service")
	return &CommentServiceImpl{
		CommentRepo: commentRepo,
		PostRepo:    postRepo,
	}
}

func (s *CommentServiceImpl) GetCommentsByPostID(postID, userID uint) (*GetCommentsResponse, error) {
	findComment := &models.Comment{}

	switch userID {
	case 0:
		findComment.PostID = postID
	default:
		findComment.PostID = postID
		findComment.AuthorID = userID
	}

	commentList, err := s.CommentRepo.FindCommentsByPostID(findComment)
	if err != nil {
		return nil, err
	}
	if len(commentList) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	result := &GetCommentsResponse{}
	for _, comment := range commentList {
		item := GetCommentResponseBody{
			ID:         comment.ID,
			Title:      comment.Title,
			Content:    comment.Content,
			AuthorName: comment.Author.Name,
			PostTitle:  comment.Post.Title,
			CreatedAt:  comment.CreatedAt,
		}
		result.Comments = append(result.Comments, item)
	}
	return result, nil
}

func (s *CommentServiceImpl) CreateCommentByPostID(postID, authorID uint, comment *CreateCommentRequest) error {
	newComment := &models.Comment{
		PostID:   postID,
		AuthorID: authorID,
		Title:    comment.Title,
		Content:  comment.Content,
	}

	err := s.CommentRepo.CreateCommentByPostID(newComment)
	if err != nil {
		return err
	}
	return nil
}

func (s *CommentServiceImpl) UpdateComment(commentID, postID, userID uint, fields *UpdateCommentRequest) error {
	comment := &models.Comment{
		ID:      commentID,
		PostID:  postID,
		Content: fields.Content,
	}
	if ok, err := s.checkPermission(postID, userID); err != nil || !ok {
		return errors.New("permission denied")
	}
	err := s.CommentRepo.UpdateCommentByCommentAndPostID(comment)
	if err != nil {
		return err
	}
	return nil
}

func (s *CommentServiceImpl) DeleteComment(commentID, PostID, userID uint) error {
	comment := &models.Comment{
		ID:     commentID,
		PostID: PostID,
	}
	if ok, err := s.checkPermission(PostID, userID); err != nil || !ok {
		return errors.New("permission denied")
	}

	err := s.CommentRepo.DeleteCommentByCommentAndPostID(comment)
	if err != nil {
		return err
	}
	return nil
}

func (s *CommentServiceImpl) checkPermission(postID, userID uint) (bool, error) {
	postCheck, err := s.PostRepo.FindByID(postID)
	if err != nil {
		return false, err
	}
	if postCheck.AuthorID != userID {
		return false, nil
	}
	return true, nil
}
