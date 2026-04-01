package unit

import (
	"context"
	"testing"

	"go-restful-api/internal/domain/entity"
	"go-restful-api/internal/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockUserRepository is a mock for UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindAll(ctx context.Context, offset, limit int) ([]*entity.User, int64, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]*entity.User), int64(args.Int(1)), args.Error(2)
}

func TestUserUseCase_GetProfile(t *testing.T) {
	mockRepo := new(MockUserRepository)
	log := zap.NewNop()
	usecase := usecase.NewUserUseCase(mockRepo, log)

	userID := uuid.New()
	expectedUser := &entity.User{
		ID:    userID,
		Name:  "Test User",
		Email: "test@example.com",
	}

	mockRepo.On("FindByID", mock.Anything, userID.String()).Return(expectedUser, nil)

	user, err := usecase.GetProfile(context.Background(), userID.String())

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_GetAllUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	log := zap.NewNop()
	usecase := usecase.NewUserUseCase(mockRepo, log)

	expectedUsers := []*entity.User{
		{Name: "User 1"},
		{Name: "User 2"},
	}

	mockRepo.On("FindAll", mock.Anything, 0, 10).Return(expectedUsers, 2, nil)

	users, total, err := usecase.GetAllUsers(context.Background(), 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, expectedUsers, users)
	mockRepo.AssertExpectations(t)
}
