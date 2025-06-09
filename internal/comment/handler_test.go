package comment_test

import (
	"fmt"
	"github.com/crafty-ezhik/blog-api/internal/comment"
	mock_comment "github.com/crafty-ezhik/blog-api/mocks/comment"
	"github.com/crafty-ezhik/blog-api/pkg/logger"
	"github.com/crafty-ezhik/blog-api/pkg/middleware"
	"github.com/crafty-ezhik/blog-api/pkg/validate"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Mocks struct {
	CommentService *mock_comment.MockCommentService
}

func setup(t *testing.T) (*comment.CommentHandlerImpl, *Mocks) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommentService := mock_comment.NewMockCommentService(ctrl)
	mockValidator := &validate.XValidator{
		Validator: validator.New(),
	}

	commentHandlerImpl := comment.NewCommentHandler(mockCommentService, mockValidator)
	mocks := &Mocks{CommentService: mockCommentService}

	return commentHandlerImpl, mocks
}

func TestCommentHandlerImpl_GetMyComment(t *testing.T) {
	logger.Log, _ = zap.NewDevelopment()
	defer logger.Log.Sync()

	commentHandlerImpl, mocks := setup(t)

	path := "/my/posts/:postId/comments"

	tests := []struct {
		name               string
		postId             any
		userId             uint
		mockSetup          func(mock *Mocks)
		handlerFunc        func(c *fiber.Ctx) error
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:   "Success",
			userId: 1,
			postId: 1,
			mockSetup: func(mock *Mocks) {
				mock.CommentService.EXPECT().GetCommentsByPostID(uint(1), uint(1)).Return(
					&comment.GetCommentsResponse{
						Comments: []comment.GetCommentResponseBody{
							{
								ID:         0,
								Title:      "TestTitle",
								Content:    "TestContent",
								AuthorName: "TestAuthorName",
								PostTitle:  "TestPostTitle",
							},
						},
					}, nil)
			},
			handlerFunc: func(c *fiber.Ctx) error {
				c.Locals(middleware.UserIDKey, uint(1))
				return c.Next()
			},
			expectedStatusCode: 200,
			expectedBody:       "\"success\":true",
		},
		{
			name:      "Invalid PostId",
			userId:    1,
			postId:    "one",
			mockSetup: nil,
			handlerFunc: func(c *fiber.Ctx) error {
				c.Locals(middleware.UserIDKey, uint(1))
				return c.Next()
			},
			expectedStatusCode: 400,
			expectedBody:       "\"error\":\"Post ID must be an integer",
		},
		{
			name:   "Comments Not Found",
			userId: 1,
			postId: 1,
			mockSetup: func(mock *Mocks) {
				mock.CommentService.EXPECT().GetCommentsByPostID(uint(1), uint(1)).Return(
					nil, gorm.ErrRecordNotFound)
			},
			handlerFunc: func(c *fiber.Ctx) error {
				c.Locals(middleware.UserIDKey, uint(1))
				return c.Next()
			},
			expectedStatusCode: 404,
			expectedBody:       "\"error\":\"Comment not found\"",
		},
		{
			name:   "Server internal error",
			userId: 1,
			postId: 1,
			mockSetup: func(mock *Mocks) {
				mock.CommentService.EXPECT().GetCommentsByPostID(uint(1), uint(1)).Return(
					nil, gorm.ErrInvalidDB)
			},
			handlerFunc: func(c *fiber.Ctx) error {
				c.Locals(middleware.UserIDKey, uint(1))
				return c.Next()
			},
			expectedStatusCode: 500,
			expectedBody:       "\"error\":\"Internal server error\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/my/posts/%v/comments", tt.postId), nil)

			app := fiber.New()
			app.Get(path, tt.handlerFunc, commentHandlerImpl.GetMyComment)

			if tt.mockSetup != nil {
				tt.mockSetup(mocks)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Contains(t, string(respBody), tt.expectedBody)
		})
	}
}

func TestCommentHandlerImpl_GetUserComments(t *testing.T)    {}
func TestCommentHandlerImpl_GetAllCommentsPost(t *testing.T) {}
func TestCommentHandlerImpl_CreateComments(t *testing.T)     {}
func TestCommentHandlerImpl_UpdateComment(t *testing.T)      {}
func TestCommentHandlerImpl_DeleteComment(t *testing.T)      {}
