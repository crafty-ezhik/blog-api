package post

type PostService interface {
	GetAllPosts() error
	GetPostById() error
	CreatePost() error
	UpdatePost() error
	DeletePost() error
}

type PostServiceImpl struct {
	PostRepo PostRepository
}

func NewPostService(postRepo PostRepository) *PostServiceImpl {
	return &PostServiceImpl{
		PostRepo: postRepo,
	}
}
