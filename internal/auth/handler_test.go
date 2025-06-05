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

	tests := []struct {
		name               string
		payload            LoginRequest
		mockSetup          func(mocks *Mocks)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name: "Login success",
			payload: LoginRequest{
				Email:    "test@test.com",
				Password: "123456",
			},
			mockSetup: func(mocks *Mocks) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
				mocks.UserRepo.EXPECT().FindByEmail("test@test.com").Return(
					&models.User{
						ID:       1,
						Email:    "test@test.com",
						Password: string(hashedPassword)}, nil)

				mocks.TokenVersion.EXPECT().GetVersion(uint(1)).Return(uint(1), nil).Times(2)
			},
			expectedStatusCode: 200,
			expectedBody:       `access_token`,
		},
		{
			name: "Invalid email format",
			payload: LoginRequest{
				Email:    "test@test....com",
				Password: "123456",
			},
			mockSetup:          nil,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `Invalid field or its absence: [Email]`,
		},
		{
			name: "Invalid password",
			payload: LoginRequest{
				Email:    "test@test.com",
				Password: "1234567",
			},
			mockSetup: func(mocks *Mocks) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
				mocks.UserRepo.EXPECT().FindByEmail("test@test.com").Return(
					&models.User{
						ID:       1,
						Email:    "test@test.com",
						Password: string(hashedPassword)}, nil)
				mocks.TokenVersion.EXPECT().GetVersion(uint(1)).Return(uint(1), nil).Times(2)
			},
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       `invalid credentials`,
		},
		{
			name: "Empty request body",
			payload: LoginRequest{
				Email:    "",
				Password: "",
			},
			mockSetup:          func(mocks *Mocks) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `Invalid field or its absence: [Email] and [Password]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

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

func TestAuthHandlerImpl_Register(t *testing.T) {
	logger.Log, _ = zap.NewDevelopment()
	defer logger.Log.Sync()

	authHandler, mocks := setup(t)

	// 1. Создаем fiber app
	app := fiber.New()

	// 2. Регистрируем проверяемый маршрут
	app.Post("/auth/register", authHandler.Register)

	tests := []struct {
		name               string
		payload            RegisterRequest
		mockSetup          func(mocks *Mocks)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name: "Successful registration",
			payload: RegisterRequest{
				Email:    "test@test.com",
				Password: "123456",
				Name:     "Test User",
				Age:      20,
			},
			mockSetup: func(mocks *Mocks) {
				// Ожидаем, что пользователя нет в базе
				mocks.UserRepo.EXPECT().FindByEmail("test@test.com").Return(nil, nil)

				// Ожидаем вызов Create
				mocks.UserRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(user *models.User) error {
					assert.Equal(t, "Test User", user.Name)
					assert.Equal(t, 20, user.Age)
					assert.Equal(t, "test@test.com", user.Email)
					assert.NotEmpty(t, user.Password)
					return nil
				})
			},
			expectedStatusCode: http.StatusCreated,
			expectedBody:       `You have successfully registered`,
		},
		{
			name: "User exists",
			payload: RegisterRequest{
				Email:    "test@test.com",
				Password: "123456",
				Name:     "Test User",
				Age:      20,
			},
			mockSetup: func(mocks *Mocks) {
				mocks.UserRepo.EXPECT().FindByEmail("test@test.com").Return(&models.User{}, nil)
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `User already exists`,
		},
		{
			name: "Missing email",
			payload: RegisterRequest{
				Password: "123456",
				Name:     "Test User",
				Age:      20,
			},
			mockSetup:          func(mocks *Mocks) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `false`,
		},
		{
			name: "Missing password",
			payload: RegisterRequest{
				Email: "test@test.com",
				Name:  "Test User",
				Age:   20,
			},
			mockSetup:          func(mocks *Mocks) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `false`,
		},
		{
			name: "Missing name",
			payload: RegisterRequest{
				Email:    "test@test.com",
				Password: "123456",
				Age:      20,
			},
			mockSetup:          func(mocks *Mocks) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `false`,
		},
		{
			name: "Missing age",
			payload: RegisterRequest{
				Email:    "test@test.com",
				Password: "123456",
				Name:     "Test User",
			},
			mockSetup:          func(mocks *Mocks) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `false`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

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
func TestAuthHandlerImpl_Logout(t *testing.T) {}

func TestAuthHandlerImpl_Refresh(t *testing.T) {

}
