package middleware

import (
	"fmt"
	"github.com/crafty-ezhik/blog-api/pkg/jwt"
	mock_jwt "github.com/crafty-ezhik/blog-api/pkg/jwt/mock"
	"github.com/crafty-ezhik/blog-api/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Mocks struct {
	JWT *mock_jwt.MockJWTInterface
}

func setup(t *testing.T) *Mocks {
	// 1. Создаем моки
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJWT := mock_jwt.NewMockJWTInterface(ctrl)

	mocks := &Mocks{
		JWT: mockJWT,
	}
	return mocks
}

func TestAuthMiddleware(t *testing.T) {
	logger.Log, _ = zap.NewDevelopment()
	defer logger.Log.Sync()

	mocks := setup(t)
	path := "/api/users/me"
	token := "token"

	tests := []struct {
		Name          string
		headerName    string
		headerValue   string
		token         string
		mockSetup     func(mocks *Mocks)
		secondHandler func(c *fiber.Ctx) error
		expectedCode  int
		expectedBody  string
	}{
		{
			Name:        "Ok",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       token,
			mockSetup: func(mocks *Mocks) {
				mocks.JWT.EXPECT().VerifyToken(token).Return(&jwt.JWTData{UserId: 1}, nil)
			},
			secondHandler: func(c *fiber.Ctx) error {
				userID := c.Locals(UserIDKey)
				return c.Status(200).SendString(fmt.Sprintf("%d", userID))
			},
			expectedCode: 200,
			expectedBody: "1",
		},
		{
			Name:         "Empty Header",
			headerName:   "",
			expectedCode: 401,
			expectedBody: `{"details":"Token is empty","err":"Unauthorized"}`,
		},
		{
			Name:         "Invalid Bearer prefix",
			headerName:   "Authorization",
			headerValue:  "Bearere token",
			expectedCode: 401,
			expectedBody: `{"details":"Token does not start with Bearer","err":"Unauthorized"}`,
		},
		{
			Name:        "Error parsing token",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockSetup: func(mocks *Mocks) {
				mocks.JWT.EXPECT().VerifyToken(token).Return(nil, jwt.ErrInvalidToken)
			},
			expectedCode: 401,
			expectedBody: `{"details":"invalid token","err":"Unauthorized"}`,
		},
		{
			Name:        "Token is invalid",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockSetup: func(mocks *Mocks) {
				mocks.JWT.EXPECT().VerifyToken(token).Return(nil, jwt.ErrInvalidToken)
			},
			expectedCode: 401,
			expectedBody: `{"details":"invalid token","err":"Unauthorized"}`,
		},
		{
			Name:        "Error when getting token version from storage",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockSetup: func(mocks *Mocks) {
				mocks.JWT.EXPECT().VerifyToken(token).Return(nil, jwt.ErrInternalServer)
			},
			expectedCode: 500,
			expectedBody: `{"details":"internal server error","err":"Unauthorized"}`,
		},
		{
			Name:        "The token version is outdated",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockSetup: func(mocks *Mocks) {
				mocks.JWT.EXPECT().VerifyToken(token).Return(nil, jwt.ErrRefreshExpired)
			},
			expectedCode: 401,
			expectedBody: `{"details":"refresh token expired due to logout / password change","err":"Unauthorized"}`,
		},
		{
			Name:        "Token expired",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			mockSetup: func(mocks *Mocks) {
				mocks.JWT.EXPECT().VerifyToken(token).Return(nil, jwt.ErrSessionExpired)
			},
			expectedCode: 401,
			expectedBody: `{"details":"session expired","err":"Unauthorized"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			app := fiber.New()
			// Тут мы сделали 2 обработчика. Первый middleware, а второй просто возвращает userID для проверки
			// корректности работы, что действительно в контекст установлено значение и его можно использовать
			// в других хендлерах
			app.Post(path, AuthMiddleware(mocks.JWT), tt.secondHandler)

			req := httptest.NewRequest(http.MethodPost, path, nil)
			req.Header.Set(tt.headerName, tt.headerValue)

			if tt.mockSetup != nil {
				tt.mockSetup(mocks)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, tt.expectedCode, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Contains(t, string(respBody), tt.expectedBody)
		})
	}
}
