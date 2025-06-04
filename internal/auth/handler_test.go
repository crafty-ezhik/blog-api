package auth

import (
	"bytes"
	"encoding/json"
	"github.com/crafty-ezhik/blog-api/internal/config"
	"github.com/crafty-ezhik/blog-api/internal/models"
	"github.com/crafty-ezhik/blog-api/internal/user"
	mock_user "github.com/crafty-ezhik/blog-api/internal/user/mock"
	"github.com/crafty-ezhik/blog-api/pkg/jwt"
	mock_jwt "github.com/crafty-ezhik/blog-api/pkg/jwt/mock"
	"github.com/crafty-ezhik/blog-api/pkg/logger"
	"github.com/crafty-ezhik/blog-api/pkg/validate"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type Mocks struct {
	UserRepo     *mock_user.MockUserRepository
	BlackList    *mock_jwt.MockBlackListStorage
	TokenVersion *mock_jwt.MockTokenVersionStorage
}

func setup(t *testing.T) (*AuthHandlerImpl, *Mocks) {
	// 1. Создаем моки
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_user.NewMockUserRepository(ctrl)
	mockBlackList := mock_jwt.NewMockBlackListStorage(ctrl)
	mockTokenVersion := mock_jwt.NewMockTokenVersionStorage(ctrl)

	// 2. Создаем экземпляр конфига
	cfg := &config.Config{
		Auth: config.AuthConfig{
			SigningKey: "FKI/0XYt3YksmneW8QxCRWdlYbINIzPdp4fpiTqXXqs=",
			SecretKey:  "",
			AccessTTL:  time.Duration(30) * time.Minute,
			RefreshTTL: time.Duration(30) * time.Hour,
		},
	}

	// 3. Создаем экземпляр jwtService
	jwtService := jwt.NewJWTService(mockBlackList, mockTokenVersion)
	jwtAuth := jwt.NewJWT(jwtService, cfg.Auth.AccessTTL, cfg.Auth.RefreshTTL, cfg.Auth.SigningKey)

	// 4. Создаем экземпляр AuthService и UserService
	authService := &AuthServiceimpl{
		jwtAuth:  jwtAuth,
		cfg:      cfg,
		UserRepo: mockUserRepo,
	}
	userService := &user.UserServiceImpl{
		UserRepo: mockUserRepo,
	}

	// 5. Создаем валидатор
	v := &validate.XValidator{
		Validator: validator.New(),
	}

	// 6. Создаем AuthHandler
	authHandler := &AuthHandlerImpl{
		UserService: userService,
		AuthService: authService,
		v:           v,
	}

	mocks := &Mocks{
		UserRepo:     mockUserRepo,
		BlackList:    mockBlackList,
		TokenVersion: mockTokenVersion,
	}
	return authHandler, mocks
}

func TestAuthHandlerImpl_Login(t *testing.T) {
	logger.Log, _ = zap.NewDevelopment()
	defer logger.Log.Sync()

	authHandler, mocks := setup(t)

	app := fiber.New()

	app.Post("/auth/login", authHandler.Login)

	t.Run("login success", func(t *testing.T) {
		payload := LoginRequest{
			Email:    "test@test.com",
			Password: "123456",
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		require.NoError(t, err)

		mocks.UserRepo.EXPECT().FindByEmail(payload.Email).Return(
			&models.User{
				ID:       1,
				Email:    payload.Email,
				Password: string(hashedPassword)}, nil)

		mocks.TokenVersion.EXPECT().GetVersion(uint(1)).Return(uint(1), nil).Times(2)

		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var loginResponse LoginResponse
		err = json.Unmarshal(respBody, &loginResponse)
		require.NoError(t, err)

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NotEmpty(t, loginResponse.AccessToken)
		assert.NotEmpty(t, loginResponse.RefreshToken)
	})
}

func TestAuthHandlerImpl_Register(t *testing.T) {
	logger.Log, _ = zap.NewDevelopment()
	defer logger.Log.Sync()

	authHandler, mocks := setup(t)

	// 1. Создаем fiber app
	app := fiber.New()

	// 2. Регистрируем проверяемый маршрут
	app.Post("/auth/register", authHandler.Register)

	t.Run("Successful registration", func(t *testing.T) {
		payload := RegisterRequest{
			Email:    "test@test.com",
			Password: "123456",
			Name:     "Test User",
			Age:      20,
		}
		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		// Ожидаем, что пользователя нет в базе
		mocks.UserRepo.EXPECT().FindByEmail(payload.Email).Return(nil, nil)

		// Ожидаем вызов Create
		mocks.UserRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(user *models.User) error {
			assert.Equal(t, payload.Name, user.Name)
			assert.Equal(t, payload.Age, user.Age)
			assert.Equal(t, payload.Email, user.Email)
			assert.NotEmpty(t, user.Password)
			return nil
		})

		resp, err := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Contains(t, string(bodyBytes), "successfully")
	})

	t.Run("Missing email", func(t *testing.T) {
		payload := RegisterRequest{
			Password: "123456",
			Name:     "Test User",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		bodyBytes, _ := io.ReadAll(resp.Body)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Contains(t, string(bodyBytes), "false")
	})
}
func TestAuthHandlerImpl_Logout(t *testing.T) {}

func TestAuthHandlerImpl_Refresh(t *testing.T) {

}
