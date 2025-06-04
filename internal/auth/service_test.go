package auth

import (
	"github.com/crafty-ezhik/blog-api/internal/config"
	"github.com/crafty-ezhik/blog-api/internal/models"
	mock_user "github.com/crafty-ezhik/blog-api/internal/user/mock"
	"github.com/crafty-ezhik/blog-api/pkg/jwt"
	mock_jwt "github.com/crafty-ezhik/blog-api/pkg/jwt/mock"
	"github.com/crafty-ezhik/blog-api/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func TestAuthServiceImpl_Login(t *testing.T) {
	logger.Log, _ = zap.NewDevelopment()
	defer logger.Log.Sync()

	// 1. Объявляем контроллер и финишируем его
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 2. Создаем мок userRepository, BlackList, TokenVersion передав туда контроллер
	mockUserRepo := mock_user.NewMockUserRepository(ctrl)
	mockBlackList := mock_jwt.NewMockBlackListStorage(ctrl)
	mockTokenVersion := mock_jwt.NewMockTokenVersionStorage(ctrl)

	// 3. Создаем экземпляр конфига
	cfg := &config.Config{
		Auth: config.AuthConfig{
			SigningKey: "FKI/0XYt3YksmneW8QxCRWdlYbINIzPdp4fpiTqXXqs=",
			SecretKey:  "",
			AccessTTL:  time.Duration(30) * time.Minute,
			RefreshTTL: time.Duration(30) * time.Hour,
		},
	}

	// 4. Создаем экземпляр jwtService
	jwtService := jwt.NewJWTService(mockBlackList, mockTokenVersion)
	jwtAuth := jwt.NewJWT(jwtService, cfg.Auth.AccessTTL, cfg.Auth.RefreshTTL, cfg.Auth.SigningKey)

	// 5. Создаем экземпляр AuthService
	authService := &AuthServiceimpl{
		cfg:      cfg,
		jwtAuth:  jwtAuth,
		UserRepo: mockUserRepo,
	}

	// 6. Начало тестов

	// 7. Тест успешной авторизации
	t.Run("Successful login", func(t *testing.T) {
		email := "test@example.com"
		password := "test"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		// 7.1 Проверка отсутствия ошибок
		require.NoError(t, err)

		// 7.2 Создаем модель искомого пользователя
		user := &models.User{
			ID:       1,
			Name:     "TestUser",
			Email:    email,
			Password: string(hashedPassword),
		}

		// 7.3 Выполняем запрос на получения пользователя с переданным email
		mockUserRepo.EXPECT().FindByEmail(email).Return(user, nil)
		mockTokenVersion.EXPECT().GetVersion(user.ID).Return(uint(1), nil).Times(2)

		// 7.4 Создаем имитация запроса
		resp, _, err := authService.Login(&LoginRequest{
			Email:    email,
			Password: password,
		})

		// 7.5 Производит проверку на наличие ошибкок и полей в ответе
		assert.NoError(t, err)
		assert.NotNil(t, resp.AccessToken)
		assert.NotNil(t, resp.RefreshToken)
	})
}

func TestAuthServiceImpl_Register(t *testing.T) {

}
