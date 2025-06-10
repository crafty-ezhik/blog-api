package user_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/crafty-ezhik/blog-api/internal/models"
	"github.com/crafty-ezhik/blog-api/internal/user"
	mock_post "github.com/crafty-ezhik/blog-api/mocks/post"
	mock_user "github.com/crafty-ezhik/blog-api/mocks/user"
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
	"time"
)

type Mocks struct {
	UserService *mock_user.MockUserService
	PostService *mock_post.MockPostService
}

func setup(t *testing.T) (user.UserHandler, *Mocks) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_user.NewMockUserService(ctrl)
	mockPostService := mock_post.NewMockPostService(ctrl)

	mockValidator := &validate.XValidator{
		Validator: validator.New(),
	}

	userHandler := user.NewUserHandler(mockUserService, mockPostService, mockValidator)
	mocks := &Mocks{
		UserService: mockUserService,
		PostService: mockPostService,
	}
	return userHandler, mocks

}

func TestUserHandlerImpl_GetByID(t *testing.T) {
	logger.Log, _ = zap.NewDevelopment()
	defer logger.Log.Sync()

	userHandler, mocks := setup(t)

	path := "/api/users/:id"

	tests := []struct {
		name               string
		userID             any
		mockSetup          func(mock *Mocks)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:   "Success",
			userID: 1,
			mockSetup: func(mock *Mocks) {
				mock.UserService.EXPECT().GetByID(uint(1)).Return(&models.User{
					Name:      "Test Name",
					Email:     "some@email.com",
					Age:       42,
					CreatedAt: time.Time{},
				}, nil)
			},
			expectedStatusCode: 200,
			expectedBody:       "true",
		},
		{
			name:               "Invalid User ID",
			userID:             "one",
			expectedStatusCode: 400,
			expectedBody:       "id must be an integer",
		},
		{
			name:   "User not found",
			userID: 5123,
			mockSetup: func(mock *Mocks) {
				mock.UserService.EXPECT().GetByID(uint(5123)).Return(nil, gorm.ErrRecordNotFound)
			},
			expectedStatusCode: 404,
			expectedBody:       "user not found",
		},
		{
			name:   "Server internal error",
			userID: 5123,
			mockSetup: func(mock *Mocks) {
				mock.UserService.EXPECT().GetByID(uint(5123)).Return(nil, errors.New("server error"))
			},
			expectedStatusCode: 500,
			expectedBody:       "Something went wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/users/%v", tt.userID), nil)

			if tt.mockSetup != nil {
				tt.mockSetup(mocks)
			}

			app := fiber.New()
			app.Get(path, userHandler.GetByID)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Contains(t, string(respBody), tt.expectedBody)
		})
	}
}
func TestUserHandlerImpl_Update(t *testing.T) {
	logger.Log, _ = zap.NewDevelopment()
	defer logger.Log.Sync()

	userHandler, mocks := setup(t)

	path := "/api/users/me"

	tests := []struct {
		name               string
		userID             any
		payload            user.UpdateUserRequest
		mockSetup          func(mock *Mocks)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:   "Success",
			userID: 1,
			payload: user.UpdateUserRequest{
				Name: "Test Name",
				Age:  53,
			},
			mockSetup: func(mock *Mocks) {
				mock.UserService.EXPECT().Update(uint(1), &models.User{Name: "Test Name", Age: 53}).Return(nil)
			},
			expectedStatusCode: 200,
			expectedBody:       "user updated",
		},
		{
			name:   "Len Name in body < 1",
			userID: 1,
			payload: user.UpdateUserRequest{
				Name: "",
				Age:  42,
			},
			expectedStatusCode: 400,
			expectedBody:       "Invalid field or its absence: [Name]",
		},
		{
			name:   "Len Name in body > 255",
			userID: 1,
			payload: user.UpdateUserRequest{
				Name: "jnoxixcyzkcazkktbsxhzkgpohpjiiahxfpxgllvsdnlksjkzssogohslhwktyutdfrgwxncocsrvwtupexvapgaptsjbwqeffnsttwxsfbpqbhjvruquupworqdamraumqpveprilfijqjquxqmohhrpzitqxlzaunoeywfyfrkukorehficrmeearfrrnjxeszkulihuqzjpoterbenhuhlbedaijuqjhiqoycshusjaekoajmhyzxpbkrjmyp",
				Age:  42,
			},
			expectedStatusCode: 400,
			expectedBody:       "Invalid field or its absence: [Name]",
		},
		{
			name:   "Len Age in body < 1",
			userID: 1,
			payload: user.UpdateUserRequest{
				Name: "Test",
				Age:  0,
			},
			expectedStatusCode: 400,
			expectedBody:       "Invalid field or its absence: [Age]",
		},
		{
			name:   "Len Age in body > 120",
			userID: 1,
			payload: user.UpdateUserRequest{
				Name: "Test",
				Age:  130,
			},
			expectedStatusCode: 400,
			expectedBody:       "Invalid field or its absence: [Age]",
		},
		{
			name:   "Server internal error",
			userID: 1,
			payload: user.UpdateUserRequest{
				Name: "Test Name",
				Age:  53,
			},
			mockSetup: func(mock *Mocks) {
				mock.UserService.EXPECT().Update(uint(1), &models.User{Name: "Test Name", Age: 53}).Return(errors.New("server error"))
			},
			expectedStatusCode: 500,
			expectedBody:       "Something went wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, path, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			if tt.mockSetup != nil {
				tt.mockSetup(mocks)
			}

			app := fiber.New()
			app.Patch(path, func(c *fiber.Ctx) error {
				c.Locals(middleware.UserIDKey, uint(1))
				return c.Next()
			}, userHandler.Update)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Contains(t, string(respBody), tt.expectedBody)
		})
	}
}
func TestUserHandlerImpl_Delete(t *testing.T) {
	logger.Log, _ = zap.NewDevelopment()
	defer logger.Log.Sync()

	userHandler, mocks := setup(t)

	path := "/api/users/:id"

	tests := []struct {
		name               string
		userID             any
		mockSetup          func(mock *Mocks)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:   "Success",
			userID: 1,
			mockSetup: func(mock *Mocks) {
				mock.UserService.EXPECT().Delete(uint(1)).Return(nil)
			},
			expectedStatusCode: 204,
			expectedBody:       "",
		},
		{
			name:               "Invalid User ID",
			userID:             "one",
			expectedStatusCode: 400,
			expectedBody:       "id must be an integer",
		},
		{
			name:   "Server internal error",
			userID: 1,
			mockSetup: func(mock *Mocks) {
				mock.UserService.EXPECT().Delete(uint(1)).Return(errors.New("server error"))
			},
			expectedStatusCode: 500,
			expectedBody:       "Something went wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/users/%v", tt.userID), nil)

			if tt.mockSetup != nil {
				tt.mockSetup(mocks)
			}

			app := fiber.New()
			app.Delete(path, userHandler.Delete)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Contains(t, string(respBody), tt.expectedBody)
		})
	}
}
func TestUserHandlerImpl_GetUserPostsByID(t *testing.T) {}
func TestUserHandlerImpl_GetUserComments(t *testing.T)  {}
