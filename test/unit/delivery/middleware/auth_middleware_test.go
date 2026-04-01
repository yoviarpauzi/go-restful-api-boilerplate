package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-restful-api/internal/delivery/http/middleware"
	"go-restful-api/internal/infrastructure/token"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	config := viper.New()
	config.Set("PASETO_ACCESS_TOKEN_SECRET", "01234567890123456789012345678912") // 32 chars
	config.Set("PASETO_REFRESH_TOKEN_SECRET", "01234567890123456789012345678912")
	config.Set("PASETO_ACCESS_TOKEN_DURATION", time.Minute)
	config.Set("PASETO_REFRESH_TOKEN_DURATION", time.Hour)

	tokenGenerator, _ := token.NewTokenGenerator(config)
	authMiddleware := middleware.NewAuth(tokenGenerator)

	app := fiber.New()
	app.Get("/protected", authMiddleware, func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	t.Run("No Authorization Header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var body map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&body)
		assert.False(t, body["success"].(bool))
	})

	t.Run("Valid Token", func(t *testing.T) {
		userID := "user-123"
		signedToken, _ := tokenGenerator.GenerateAccessToken(userID)

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+signedToken)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Invalid Token Format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "InvalidFormat token")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
