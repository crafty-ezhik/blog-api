package post

type CreateRequest struct {
	Title string `json:"title" validate:"required,max=255"`
	Text  string `json:"text" validate:"required"`
}
