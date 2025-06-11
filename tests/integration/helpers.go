package integration

import (
	"bytes"
	"encoding/json"
	"github.com/crafty-ezhik/blog-api/internal/auth"
	"github.com/crafty-ezhik/blog-api/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func registerAndLogin(t testing.TB, app *fiber.App) jwt.Tokens {
	payload := auth.RegisterRequest{
		Email:    "test@test.com",
		Password: "12345678",
		Name:     "Test",
		Age:      30,
	}
	body, _ := json.Marshal(payload)

	// Register
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	app.Test(req)

	// Login
	req = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	var loginResp auth.LoginResponse
	bodyBytes, _ := io.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &loginResp)

	return jwt.Tokens{
		AccessToken:  loginResp.AccessToken,
		RefreshToken: loginResp.RefreshToken,
	}
}
