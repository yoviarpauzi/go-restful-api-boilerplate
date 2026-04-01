package unit

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-restful-api/internal/delivery/http/handler"
	"go-restful-api/internal/delivery/http/request"
	"go-restful-api/internal/delivery/http/response"
	"go-restful-api/internal/domain/entity"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockAuthUseCase is a mock for AuthUseCase interface
type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) Register(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockAuthUseCase) Login(ctx context.Context, email, password string) (string, string, string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.String(1), args.String(2), args.Error(3)
}

func (m *MockAuthUseCase) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	args := m.Called(ctx, userID, oldPassword, newPassword)
	return args.Error(0)
}

func (m *MockAuthUseCase) ResetPassword(ctx context.Context, email, newPassword string) error {
	args := m.Called(ctx, email, newPassword)
	return args.Error(0)
}

func TestAuthHandler_Register(t *testing.T) {
	app := fiber.New()
	mockUseCase := new(MockAuthUseCase)
	log := zap.NewNop()
	validate := validator.New()
	h := handler.NewAuthHandler(mockUseCase, log, validate)

	app.Post("/auth/register", h.Register)

	t.Run("Success", func(t *testing.T) {
		reqBody := request.RegisterRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.On("Register", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
			return u.Name == reqBody.Name && u.Email == reqBody.Email
		})).Return(nil).Run(func(args mock.Arguments) {
			u := args.Get(1).(*entity.User)
			u.ID = uuid.New()
			u.CreatedAt = time.Now()
			u.UpdatedAt = time.Now()
		})

		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var result response.SuccessResponse
		_ = json.NewDecoder(resp.Body).Decode(&result)
		assert.True(t, result.Success)
		assert.Equal(t, "register successfully", result.Message)
	})

	t.Run("Validation Error", func(t *testing.T) {
		reqBody := request.RegisterRequest{
			Name: "", // Name is required
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	app := fiber.New()
	mockUseCase := new(MockAuthUseCase)
	log := zap.NewNop()
	validate := validator.New()
	h := handler.NewAuthHandler(mockUseCase, log, validate)

	app.Post("/auth/login", h.Login)

	t.Run("Success", func(t *testing.T) {
		reqBody := request.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		userID := uuid.New().String()
		mockUseCase.On("Login", mock.Anything, "test@example.com", "password123").Return("mock-token", "mock-refresh-token", userID, nil)

		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		
		// Check Cookie
		cookies := resp.Cookies()
		var refreshTokenCookie *http.Cookie
		for _, v := range cookies {
			if v.Name == "refresh_token" {
				refreshTokenCookie = v
			}
		}
		assert.NotNil(t, refreshTokenCookie)
		assert.Equal(t, "mock-refresh-token", refreshTokenCookie.Value)

		var result response.SuccessResponse
		_ = json.NewDecoder(resp.Body).Decode(&result)
		assert.True(t, result.Success)
		assert.Contains(t, result.Message, "login successfully")
		
		data := result.Data.(map[string]interface{})
		assert.Equal(t, "mock-token", data["access_token"])
		assert.Equal(t, userID, data["user_id"])
	})

	t.Run("Unauthorized", func(t *testing.T) {
		reqBody := request.LoginRequest{
			Email:    "wrong@example.com",
			Password: "wrong",
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.On("Login", mock.Anything, "wrong@example.com", "wrong").Return("", "", "", assert.AnError)

		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestAuthHandler_ChangePassword(t *testing.T) {
	app := fiber.New()
	mockUseCase := new(MockAuthUseCase)
	log := zap.NewNop()
	validate := validator.New()
	h := handler.NewAuthHandler(mockUseCase, log, validate)

	userID := "user-123"
	app.Post("/auth/change-password", func(c *fiber.Ctx) error {
		c.Locals("currentUser", userID)
		return h.ChangePassword(c)
	})

	t.Run("Success", func(t *testing.T) {
		reqBody := request.ChangePasswordRequest{
			OldPassword: "old-password",
			NewPassword: "new-password",
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.On("ChangePassword", mock.Anything, userID, "old-password", "new-password").Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/auth/change-password", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result response.SuccessResponse
		_ = json.NewDecoder(resp.Body).Decode(&result)
		assert.True(t, result.Success)
		assert.Equal(t, "password changed successfully", result.Message)
	})
}

func TestAuthHandler_ResetPassword(t *testing.T) {
	app := fiber.New()
	mockUseCase := new(MockAuthUseCase)
	log := zap.NewNop()
	validate := validator.New()
	h := handler.NewAuthHandler(mockUseCase, log, validate)

	app.Post("/auth/reset-password", h.ResetPassword)

	t.Run("Success", func(t *testing.T) {
		reqBody := request.ResetPasswordRequest{
			Email:       "test@example.com",
			NewPassword: "new-password",
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.On("ResetPassword", mock.Anything, "test@example.com", "new-password").Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/auth/reset-password", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result response.SuccessResponse
		_ = json.NewDecoder(resp.Body).Decode(&result)
		assert.True(t, result.Success)
		assert.Equal(t, "password reset successfully", result.Message)
	})
}
