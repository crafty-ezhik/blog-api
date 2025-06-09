package comment

import (
	"fmt"
	mock_comment "github.com/crafty-ezhik/blog-api/internal/comment/mock"
	"github.com/crafty-ezhik/blog-api/internal/models"
	"github.com/crafty-ezhik/blog-api/pkg/validate"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Mocks struct {
	CommentService *mock_comment.MockCommentService
}

func setup(t *testing.T) (*CommentHandlerImpl, *Mocks) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommentService := mock_comment.NewMockCommentService(ctrl)
	mockValidator := &validate.XValidator{
		Validator: validator.New(),
	}

	commentHandlerImpl := &CommentHandlerImpl{
		CommentService: mockCommentService,
		v:              mockValidator,
	}
	mocks := &Mocks{CommentService: mockCommentService}

	return commentHandlerImpl, mocks
}

func TestCommentHandlerImpl_GetMyComment(t *testing.T) {
	commentHandlerImpl, mocks := setup(t)

	app := fiber.New()
	path := "/my/posts/:postId/comments"
	app.Get(path, commentHandlerImpl.GetMyComment)

	tests := []struct {
		name               string
		postId             uint
		userId             uint
		mockSetup          func(mock *Mocks)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:   "Success",
			userId: 1,
			postId: 1,
			mockSetup: func(mock *Mocks) {
				mock.CommentService.EXPECT().GetCommentsByPostID(uint(1), uint(1)).Return(&models.Comment{}, nil)
			},
			expectedStatusCode: 200,
			expectedBody:       "\"success\": true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/my/posts/1/comments", nil)
			req.AddCookie(&http.Cookie{Name: "user_id", Value: fmt.Sprint(tt.userId)})

			if tt.mockSetup != nil {
				tt.mockSetup(mocks)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
			assert.Contains(t, resp.Body, tt.expectedBody)
		})
	}
}

func TestCommentHandlerImpl_GetUserComments(t *testing.T)    {}
func TestCommentHandlerImpl_GetAllCommentsPost(t *testing.T) {}
func TestCommentHandlerImpl_CreateComments(t *testing.T)     {}
func TestCommentHandlerImpl_UpdateComment(t *testing.T)      {}
func TestCommentHandlerImpl_DeleteComment(t *testing.T)      {}
