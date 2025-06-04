package auth

import (
	"errors"
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

		// 7.5 Производит проверку на наличие ошибок и полей в ответе
		assert.NoError(t, err)
		assert.NotNil(t, resp.AccessToken)
		assert.NotNil(t, resp.RefreshToken)
	})

	// 8. Случай, когда пользователь не найден
	t.Run("User not found", func(t *testing.T) {
		email := "notfound@example.com"
		mockUserRepo.EXPECT().FindByEmail(email).Return(nil, errors.New(ErrInvalidCredentials))

		resp, _, err := authService.Login(&LoginRequest{
			Email:    email,
			Password: "any",
		})
		assert.Error(t, err)
		assert.EqualError(t, err, ErrInvalidCredentials)
		assert.Nil(t, resp)
	})

	t.Run("Wrong password", func(t *testing.T) {
		email := "test@example.com"
		correctPassword := "test"
		wrongPassword := "wrong"

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)
		require.NoError(t, err)

		user := &models.User{
			ID:       1,
			Name:     "TestUser",
			Email:    email,
			Password: string(hashedPassword),
		}

		mockUserRepo.EXPECT().FindByEmail(email).Return(user, nil)

		resp, _, err := authService.Login(&LoginRequest{
			Email:    email,
			Password: wrongPassword,
		})
		assert.Error(t, err)
		assert.EqualError(t, err, ErrInvalidCredentials)
		assert.Nil(t, resp)
	})
}

func TestAuthServiceImpl_Register(t *testing.T) {
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
	t.Run("Successful register", func(t *testing.T) {
		email := "test@example.com"
		request := &RegisterRequest{
			Name:     "Name",
			Email:    email,
			Password: "test",
			Age:      30,
		}

		// Ожидаем, что пользователя нет в базе
		mockUserRepo.EXPECT().FindByEmail(email).Return(nil, nil)

		// Ожидаем вызов Create
		mockUserRepo.EXPECT().Create(gomock.All()).DoAndReturn(func(user *models.User) error {
			assert.Equal(t, request.Name, user.Name)
			assert.Equal(t, request.Age, user.Age)
			assert.Equal(t, request.Email, user.Email)
			assert.NotEmpty(t, user.Password)
			return nil
		})

		ok, err := authService.Register(request)
		assert.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("User already exists", func(t *testing.T) {
		email := "test@example.com"
		request := &RegisterRequest{
			Email:    email,
			Password: "test",
		}

		mockUserRepo.EXPECT().FindByEmail(email).Return(&models.User{Email: email}, nil)

		ok, err := authService.Register(request)
		assert.Error(t, err)
		assert.False(t, ok)
		assert.Contains(t, err.Error(), ErrUserExisted)
	})
}
