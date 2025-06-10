package post_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/crafty-ezhik/blog-api/internal/models"
	"github.com/crafty-ezhik/blog-api/internal/post"
	mock_post "github.com/crafty-ezhik/blog-api/mocks/post"
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
	PostService *mock_post.MockPostService
}

func setup(t *testing.T) (post.PostHandler, *Mocks) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostService := mock_post.NewMockPostService(ctrl)
	mockValidator := &validate.XValidator{
		Validator: validator.New(),
	}

	postHandler := post.NewPostHandler(mockPostService, mockValidator)
	mocks := &Mocks{
		PostService: mockPostService,
	}
	return postHandler, mocks
}

func TestPostHandlerImpl_GetPostById(t *testing.T) {
	logger.Log, _ = zap.NewDevelopment()
	defer logger.Log.Sync()

	postHandler, mocks := setup(t)

	path := "/api/posts/:id"

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
				mock.PostService.EXPECT().GetPostById(uint(1)).Return(&models.Post{
					ID:       1,
					Title:    "TestTitle",
					Text:     "TestText",
					AuthorID: 1,
				}, nil)
			},
			expectedStatusCode: 200,
			expectedBody:       "true",
		},
		{
			name:               "Invalid Post Id",
			postId:             "one",
			expectedStatusCode: 400,
			expectedBody:       "Post Id is invalid",
		},
		{
			name:   "Server Internal Error",
			postId: 1,
			mockSetup: func(mock *Mocks) {
				mock.PostService.EXPECT().GetPostById(uint(1)).Return(nil, errors.New("server Error"))
			},
			expectedStatusCode: 500,
			expectedBody:       "Something went wrong",
		},
		{
			name:   "Post Not Found",
			postId: 99,
			mockSetup: func(mock *Mocks) {
				mock.PostService.EXPECT().GetPostById(uint(99)).Return(nil, gorm.ErrRecordNotFound)
			},
			expectedStatusCode: 404,
			expectedBody:       "Post not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/posts/%v", tt.postId), nil)

			app := fiber.New()
			app.Get(path, postHandler.GetPostById)

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

func TestPostHandlerImpl_GetAllPosts(t *testing.T) {
	logger.Log, _ = zap.NewDevelopment()
	defer logger.Log.Sync()

	postHandler, mocks := setup(t)

	path := "/api/posts/"

	tests := []struct {
		name               string
		mockSetup          func(mock *Mocks)
		expectedStatusCode int
		expectedBody       string
		expectedLengthData int
	}{
		{
			name: "Success",
			mockSetup: func(mock *Mocks) {
				mock.PostService.EXPECT().GetAllPosts().Return([]models.Post{
					models.Post{}, models.Post{}}, nil)
			},
			expectedStatusCode: 200,
			expectedBody:       "true",
			expectedLengthData: 2,
		},
		{
			name: "Server internal Error",
			mockSetup: func(mock *Mocks) {
				mock.PostService.EXPECT().GetAllPosts().Return([]models.Post{}, errors.New("server Error"))
			},
			expectedStatusCode: 500,
			expectedBody:       "Something went wrong",
		},
		{
			name: "Posts Not Found",
			mockSetup: func(mock *Mocks) {
				mock.PostService.EXPECT().GetAllPosts().Return([]models.Post{}, gorm.ErrRecordNotFound)
			},
			expectedStatusCode: 404,
			expectedBody:       "Posts not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)

			app := fiber.New()
			app.Get(path, postHandler.GetAllPosts)

			if tt.mockSetup != nil {
				tt.mockSetup(mocks)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Contains(t, string(respBody), tt.expectedBody)

			if tt.expectedLengthData > 0 {
				var jsonData map[string]interface{}
				err = json.Unmarshal(respBody, &jsonData)
				require.NoError(t, err)
				assert.Len(t, jsonData, tt.expectedLengthData)
			}

		})
	}
}

func TestPostHandlerImpl_CreatePost(t *testing.T) {
	logger.Log, _ = zap.NewDevelopment()
	defer logger.Log.Sync()

	postHandler, mocks := setup(t)

	path := "/api/posts/"

	tests := []struct {
		name               string
		payload            post.CreateRequest
		mockSetup          func(mock *Mocks)
		expectedStatusCode int
		expectedBody       any
	}{
		{
			name: "Success",
			payload: post.CreateRequest{
				Title: "TestTitle",
				Text:  "TestText",
			},
			mockSetup: func(mock *Mocks) {
				mock.PostService.EXPECT().CreatePost(gomock.Any()).Return(nil)
			},
			expectedStatusCode: 201,
			expectedBody:       "true",
		},
		{
			name: "Empty Title in body",
			payload: post.CreateRequest{
				Text: "TestText",
			},
			expectedStatusCode: 400,
			expectedBody:       "Invalid field or its absence: [Title]",
		},
		{
			name: "Empty Text in body",
			payload: post.CreateRequest{
				Title: "TestText",
			},
			expectedStatusCode: 400,
			expectedBody:       "Invalid field or its absence: [Text]",
		},
		{
			name: "Server Internal Error",
			payload: post.CreateRequest{
				Title: "TestTitle",
				Text:  "TestText",
			},
			mockSetup: func(mock *Mocks) {
				mock.PostService.EXPECT().CreatePost(gomock.Any()).Return(errors.New("server Error"))
			},
			expectedStatusCode: 500,
			expectedBody:       "Something went wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, path, bytes.NewBuffer(body))
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
				postHandler.CreatePost)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Contains(t, string(respBody), tt.expectedBody)

		})
	}
}

func TestPostHandlerImpl_DeletePost(t *testing.T) {
	logger.Log, _ = zap.NewDevelopment()
	defer logger.Log.Sync()

	postHandler, mocks := setup(t)

	path := "api/posts/:id"

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
				mock.PostService.EXPECT().DeletePost(uint(1)).Return(nil)
			},
			expectedStatusCode: 204,
			expectedBody:       "",
		},
		{
			name:               "Invalid Post Id",
			postId:             "one",
			expectedStatusCode: 400,
			expectedBody:       "Post Id is invalid",
		},
		{
			name:   "Server internal Error",
			postId: 1,
			mockSetup: func(mock *Mocks) {
				mock.PostService.EXPECT().DeletePost(uint(1)).Return(errors.New("server Error"))
			},
			expectedStatusCode: 500,
			expectedBody:       "Something went wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/posts/%v", tt.postId), nil)

			app := fiber.New()
			app.Delete(path, postHandler.DeletePost)

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

func TestPostHandlerImpl_UpdatePost(t *testing.T) {
	logger.Log, _ = zap.NewDevelopment()
	defer logger.Log.Sync()

	postHandler, mocks := setup(t)

	path := "/api/posts/:id"

	tests := []struct {
		name               string
		postId             any
		payload            post.UpdateRequest
		mockSetup          func(mock *Mocks)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:   "Success",
			postId: 1,
			payload: post.UpdateRequest{
				Title: "TestTitle",
				Text:  "TestText",
			},
			mockSetup: func(mock *Mocks) {
				mock.PostService.EXPECT().UpdatePost(uint(1), gomock.Any()).Return(nil)
			},
			expectedStatusCode: 200,
			expectedBody:       "post updated",
		},
		{
			name:               "Invalid Post Id",
			postId:             "one",
			expectedStatusCode: 400,
			expectedBody:       "Post Id is invalid",
		},
		{
			name:   "Server internal Error",
			postId: 1,
			payload: post.UpdateRequest{
				Title: "TestTitle",
				Text:  "TestText",
			},
			mockSetup: func(mock *Mocks) {
				mock.PostService.EXPECT().UpdatePost(uint(1), gomock.Any()).Return(errors.New("server Error"))
			},
			expectedStatusCode: 500,
			expectedBody:       "Something went wrong",
		},
		{
			name:   "Empty Text in body",
			postId: 1,
			payload: post.UpdateRequest{
				Title: "TestText",
			},
			expectedStatusCode: 400,
			expectedBody:       "Invalid field or its absence: [Text]",
		},
		{
			name:   "Empty Title in body",
			postId: 1,
			payload: post.UpdateRequest{
				Text: "TestText",
			},
			expectedStatusCode: 400,
			expectedBody:       "Invalid field or its absence: [Title]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/posts/%v", tt.postId), bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			if tt.mockSetup != nil {
				tt.mockSetup(mocks)
			}

			app := fiber.New()
			app.Put(path, postHandler.UpdatePost)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Contains(t, string(respBody), tt.expectedBody)
		})
	}
}
