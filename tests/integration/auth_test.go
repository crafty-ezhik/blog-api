package integration

import (
	"bytes"
	"encoding/json"
	"github.com/crafty-ezhik/blog-api/internal/auth"
	"github.com/crafty-ezhik/blog-api/internal/comment"
	"github.com/crafty-ezhik/blog-api/internal/models"
	"github.com/crafty-ezhik/blog-api/internal/post"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AuthIntegrationSuite struct {
	suite.Suite
	app *fiber.App
	db  *gorm.DB
}

// SetupSuite - вызываются настройки для всех тестов
func (s *AuthIntegrationSuite) SetupSuite() {
	s.db = SetupTestDB()
	s.app = SetupApp(s.db)
}

// SetupTest - производятся действия перед каждым тестом
func (s *AuthIntegrationSuite) SetupTest() {
	MigrateTables(s.db)
}

// TearDownTest - производятся действия после каждого теста
func (s *AuthIntegrationSuite) TearDownTest() {
	CleanupTables(s.db)
}

// TearDownSuite - производятся действия после всех тестов
func (s *AuthIntegrationSuite) TearDownSuite() {
	TeardownTestDB(s.db)
}

// TestAuthIntegrationSuite - Запуск всех тестов для AuthIntegrationSuite
func TestAuthIntegrationSuite(t *testing.T) {
	suite.Run(t, new(AuthIntegrationSuite))
}

func (s *AuthIntegrationSuite) Test_Register_Login_Flow() {
	payload := auth.RegisterRequest{
		Email:    "test@test.com",
		Password: "12345678",
		Name:     "Test",
		Age:      30,
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	req = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err = s.app.Test(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var loginResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, &loginResp)
	s.Require().NoError(err)
	s.NotEmpty(loginResp.AccessToken)
	s.NotEmpty(loginResp.RefreshToken)
}

func (s *AuthIntegrationSuite) Test_Create_Post_And_Comment() {
	tokens := registerAndLogin(s.T(), s.app)

	postData := post.CreateRequest{
		Title: "Test",
		Text:  "Test",
	}

	body, _ := json.Marshal(postData)
	req := httptest.NewRequest(http.MethodPost, "/api/posts", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	var postResp struct {
		Success bool `json:"success"`
		Data    models.Post
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, &postResp)
	s.Require().NoError(err)

	commentData := comment.CreateCommentRequest{
		Title:   "Text",
		Content: "Text",
	}

	body, _ = json.Marshal(commentData)
	req = httptest.NewRequest(http.MethodPost, "/api/posts/1/comments", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = s.app.Test(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)
}

func (s *AuthIntegrationSuite) Test_Refresh_Token() {
	tokens := registerAndLogin(s.T(), s.app)

	// Устанавливаем refresh_token в куки
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
	req.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: tokens.RefreshToken,
	})

	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var tokensResp struct {
		AccessToken string `json:"access_token"`
	}
	refreshToken := resp.Cookies()[0].Value
	bodyBytes, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, &tokens)
	s.Require().NoError(err)
	s.NotEqual(tokens.AccessToken, tokensResp.AccessToken)
	s.NotEmpty(refreshToken)
}

func (s *AuthIntegrationSuite) Test_Logout_BlacklistsToken() {
	tokens := registerAndLogin(s.T(), s.app)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: tokens.RefreshToken,
	})
	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	// Проверяем, что токен больше не работает
	req = httptest.NewRequest(http.MethodGet, "/api/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

	resp, err = s.app.Test(req)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusUnauthorized, resp.StatusCode)
}
