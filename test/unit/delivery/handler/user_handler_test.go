package unit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-restful-api/internal/delivery/http/handler"
	"go-restful-api/internal/delivery/http/response"
	"go-restful-api/internal/domain/entity"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockUserUseCase is a mock for UserUseCase interface
type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) GetProfile(ctx context.Context, id string) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUseCase) GetAllUsers(ctx context.Context, page, pageSize int) ([]*entity.User, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*entity.User), int64(args.Int(1)), args.Error(2)
}

func (m *MockUserUseCase) UpdateUser(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserUseCase) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUserHandler_GetByID(t *testing.T) {
	app := fiber.New()
	mockUseCase := new(MockUserUseCase)
	log := zap.NewNop()
	validate := validator.New()
	h := handler.NewUserHandler(mockUseCase, log, validate)

	userID := uuid.New().String()
	user := &entity.User{
		ID:    uuid.MustParse(userID),
		Name:  "Test User",
		Email: "test@example.com",
	}

	app.Get("/users/:id", h.GetByID)

	mockUseCase.On("GetProfile", mock.Anything, userID).Return(user, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/"+userID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result response.SuccessResponse
	_ = json.NewDecoder(resp.Body).Decode(&result)

	assert.True(t, result.Success)
	assert.Equal(t, "fetch profile successfully", result.Message)
	
	// Check data field
	data := result.Data.(map[string]interface{})
	assert.Equal(t, "Test User", data["name"])
	assert.Equal(t, "test@example.com", data["email"])
}
