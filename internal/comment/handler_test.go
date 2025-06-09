package comment_test

import (
	"bytes"
	"encoding/json"
	"errors"
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

func TestCommentHandlerImpl_GetUserComments(t *testing.T) {
	logger.Log, _ = zap.NewDevelopment()
	defer logger.Log.Sync()

	commentHandlerImpl, mocks := setup(t)

	path := "/:id/posts/:postId/comments"

	tests := []struct {
		name               string
		postId             any
		userId             any
		mockSetup          func(mock *Mocks)
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
			expectedStatusCode: 200,
			expectedBody:       "\"success\":true",
		},
		{
			name:               "Invalid PostId",
			userId:             1,
			postId:             "one",
			mockSetup:          nil,
			expectedStatusCode: 400,
			expectedBody:       "\"error\":\"Post ID must be an integer",
		},
		{
			name:               "Invalid UserId",
			userId:             "one",
			postId:             1,
			mockSetup:          nil,
			expectedStatusCode: 400,
			expectedBody:       "\"error\":\"User ID must be an integer",
		},
		{
			name:   "Comments Not Found",
			userId: 1,
			postId: 1,
			mockSetup: func(mock *Mocks) {
				mock.CommentService.EXPECT().GetCommentsByPostID(uint(1), uint(1)).Return(
					nil, gorm.ErrRecordNotFound)
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
			expectedStatusCode: 500,
			expectedBody:       "\"error\":\"Internal server error\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%v/posts/%v/comments", tt.userId, tt.postId), nil)

			app := fiber.New()
			app.Get(path, commentHandlerImpl.GetUserComments)

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
func TestCommentHandlerImpl_GetAllCommentsPost(t *testing.T) {
	logger.Log, _ = zap.NewDevelopment()
	defer logger.Log.Sync()

	commentHandlerImpl, mocks := setup(t)

	path := "/:id/comments"

	tests := []struct {
		name               string
		postId             any
		mockSetup          func(mock *Mocks)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:   "Success",
			postId: 1,
			mockSetup: func(mock *Mocks) {
				mock.CommentService.EXPECT().GetCommentsByPostID(uint(1), uint(0)).Return(
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
			expectedStatusCode: 200,
			expectedBody:       "\"success\":true",
		},
		{
			name:               "Invalid PostId",
			postId:             "one",
			mockSetup:          nil,
			expectedStatusCode: 400,
			expectedBody:       "\"error\":\"Post ID must be an integer",
		},
		{
			name:   "Comments Not Found",
			postId: 1,
			mockSetup: func(mock *Mocks) {
				mock.CommentService.EXPECT().GetCommentsByPostID(uint(1), uint(0)).Return(
					nil, gorm.ErrRecordNotFound)
			},
			expectedStatusCode: 404,
			expectedBody:       "\"error\":\"Comment not found\"",
		},
		{
			name:   "Server internal error",
			postId: 1,
			mockSetup: func(mock *Mocks) {
				mock.CommentService.EXPECT().GetCommentsByPostID(uint(1), uint(0)).Return(
					nil, gorm.ErrInvalidDB)
			},
			expectedStatusCode: 500,
			expectedBody:       "\"error\":\"Internal server error\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%v/comments", tt.postId), nil)

			app := fiber.New()
			app.Get(path, commentHandlerImpl.GetAllCommentsPost)
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
func TestCommentHandlerImpl_CreateComments(t *testing.T) {
	logger.Log, _ = zap.NewDevelopment()
	defer logger.Log.Sync()

	commentHandlerImpl, mocks := setup(t)

	path := "/posts/:id/comments"

	tests := []struct {
		name               string
		postId             any
		userId             any
		payload            comment.CreateCommentRequest
		mockSetup          func(mock *Mocks)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:   "Success",
			userId: 1,
			postId: 1,
			payload: comment.CreateCommentRequest{
				Title:   "TestTitle",
				Content: "TestContent",
			},
			mockSetup: func(mock *Mocks) {
				mock.CommentService.EXPECT().CreateCommentByPostID(uint(1), uint(1), gomock.Any()).Return(nil)
			},
			expectedStatusCode: 201,
			expectedBody:       "Comment created successfully",
		},
		{
			name:   "Empty Title",
			userId: 1,
			postId: 1,
			payload: comment.CreateCommentRequest{
				Title:   "",
				Content: "TestContent",
			},
			expectedStatusCode: 400,
			expectedBody:       "Invalid field or its absence: [Title]",
		},
		{
			name:   "Empty Content",
			userId: 1,
			postId: 1,
			payload: comment.CreateCommentRequest{
				Title:   "TestTitle",
				Content: "",
			},
			expectedStatusCode: 400,
			expectedBody:       "Invalid field or its absence: [Content]",
		},
		{
			name:   "Server internal error",
			userId: 1,
			postId: 1,
			payload: comment.CreateCommentRequest{
				Title:   "TestTitle",
				Content: "TestContent",
			},
			mockSetup: func(mock *Mocks) {
				mock.CommentService.EXPECT().CreateCommentByPostID(uint(1), uint(1), gomock.Any()).Return(errors.New("error"))
			},
			expectedStatusCode: 500,
			expectedBody:       "Something went wrong",
		},
		{
			name:   "Invalid PostId",
			userId: 1,
			postId: "one",
			payload: comment.CreateCommentRequest{
				Title:   "TestTitle",
				Content: "TestContent",
			},
			expectedStatusCode: 400,
			expectedBody:       "Post ID must be an integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/posts/%v/comments", tt.postId), bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			if tt.mockSetup != nil {
				tt.mockSetup(mocks)
			}

			app := fiber.New()
			app.Post(path,
				func(c *fiber.Ctx) error {
					c.Locals(middleware.UserIDKey, uint(1))
					return c.Next()
				},
				commentHandlerImpl.CreateComments)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Contains(t, string(respBody), tt.expectedBody)
		})
	}
}
func TestCommentHandlerImpl_UpdateComment(t *testing.T) {
	logger.Log, _ = zap.NewDevelopment()
	defer logger.Log.Sync()

	commentHandlerImpl, mocks := setup(t)
	path := "/posts/:id/comments/:commentId"

	tests := []struct {
		name               string
		commentId          any
		postId             any
		payload            comment.UpdateCommentRequest
		mockSetup          func(mock *Mocks)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:      "Success",
			postId:    1,
			commentId: 1,
			payload: comment.UpdateCommentRequest{
				Content: "NewContent",
			},
			mockSetup: func(mock *Mocks) {
				mock.CommentService.EXPECT().UpdateComment(uint(1), uint(1), uint(1),
					&comment.UpdateCommentRequest{Content: "NewContent"}).Return(nil)
			},
			expectedStatusCode: 200,
			expectedBody:       "Comment updated successfully",
		},
		{
			name:               "Invalid CommentId",
			postId:             1,
			commentId:          "one",
			expectedStatusCode: 400,
			expectedBody:       "Comment ID must be an integer",
		},
		{
			name:               "Invalid PostId",
			postId:             "one",
			commentId:          1,
			expectedStatusCode: 400,
			expectedBody:       "Post ID must be an integer",
		},
		{
			name:      "Server internal error",
			commentId: 1,
			postId:    1,
			payload: comment.UpdateCommentRequest{
				Content: "NewContent",
			},
			mockSetup: func(mock *Mocks) {
				mock.CommentService.EXPECT().UpdateComment(uint(1), uint(1), uint(1),
					&comment.UpdateCommentRequest{Content: "NewContent"}).Return(errors.New("error"))
			},
			expectedStatusCode: 500,
			expectedBody:       "Something went wrong",
		},
		{
			name:      "Permission denied",
			commentId: 1,
			postId:    2,
			payload: comment.UpdateCommentRequest{
				Content: "NewContent",
			},
			mockSetup: func(mock *Mocks) {
				mock.CommentService.EXPECT().UpdateComment(uint(1), uint(2), uint(1),
					&comment.UpdateCommentRequest{Content: "NewContent"}).Return(comment.ErrPermissionDenied)
			},
			expectedStatusCode: 403,
			expectedBody:       "Permission denied",
		},
		{
			name:               "Invalid Body",
			commentId:          1,
			postId:             2,
			payload:            comment.UpdateCommentRequest{},
			expectedStatusCode: 400,
			expectedBody:       "Invalid field or its absence: [Content]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch,
				fmt.Sprintf("/posts/%v/comments/%v", tt.postId, tt.commentId),
				bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			if tt.mockSetup != nil {
				tt.mockSetup(mocks)
			}

			app := fiber.New()
			app.Patch(path,
				func(c *fiber.Ctx) error {
					c.Locals(middleware.UserIDKey, uint(1))
					return c.Next()
				},
				commentHandlerImpl.UpdateComment)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Contains(t, string(respBody), tt.expectedBody)

		})
	}
}
func TestCommentHandlerImpl_DeleteComment(t *testing.T) {
	logger.Log, _ = zap.NewDevelopment()
	defer logger.Log.Sync()

	commentHandlerImpl, mocks := setup(t)

	path := "/posts/:id/comments/:commentId"

	tests := []struct {
		name               string
		postId             any
		commentId          any
		payload            comment.CreateCommentRequest
		mockSetup          func(mock *Mocks)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:      "Success",
			postId:    1,
			commentId: 1,
			mockSetup: func(mock *Mocks) {
				mock.CommentService.EXPECT().DeleteComment(uint(1), uint(1), uint(1)).Return(nil)
			},
			expectedStatusCode: 204,
			expectedBody:       "",
		},
		{
			name:               "Invalid CommentId",
			postId:             1,
			commentId:          "one",
			expectedStatusCode: 400,
			expectedBody:       "Comment ID must be an integer",
		},
		{
			name:               "Invalid PostId",
			postId:             "one",
			commentId:          1,
			expectedStatusCode: 400,
			expectedBody:       "Post ID must be an integer",
		},
		{
			name:      "Server internal error",
			postId:    1,
			commentId: 1,
			mockSetup: func(mock *Mocks) {
				mock.CommentService.EXPECT().DeleteComment(uint(1), uint(1), uint(1)).Return(errors.New("error"))
			},
			expectedStatusCode: 500,
			expectedBody:       "error",
		},
		{
			name:      "Permission denied",
			postId:    1,
			commentId: 1,
			mockSetup: func(mock *Mocks) {
				mock.CommentService.EXPECT().DeleteComment(uint(1), uint(1), uint(1)).Return(comment.ErrPermissionDenied)
			},
			expectedStatusCode: 403,
			expectedBody:       "Permission denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/posts/%v/comments/%v", tt.postId, tt.commentId), nil)

			if tt.mockSetup != nil {
				tt.mockSetup(mocks)
			}
			app := fiber.New()
			app.Delete(path,
				func(c *fiber.Ctx) error {
					c.Locals(middleware.UserIDKey, uint(1))
					return c.Next()
				},
				commentHandlerImpl.DeleteComment)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Contains(t, string(respBody), tt.expectedBody)
		})
	}
}
