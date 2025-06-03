package comment

import "time"

type CreateCommentRequest struct {
	Title   string `json:"title" validate:"required,max=255"`
	Content string `json:"content" validate:"required,max=255"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required,max=255"`
}

type GetCommentsResponse struct {
	Comments []GetCommentResponseBody `json:"comments"`
}

type GetCommentResponseBody struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	AuthorName string    `json:"author_name"`
	PostTitle  string    `json:"post_title"`
	CreatedAt  time.Time `json:"created_at"`
}
