package usecase

import (
	"context"

	"go-restful-api/internal/domain/entity"
	"go-restful-api/internal/domain/repository"
	domainUseCase "go-restful-api/internal/domain/usecase"
	"go.uber.org/zap"
)

type UserUseCaseImpl struct {
	UserRepo repository.UserRepository
	Log      *zap.Logger
}

func NewUserUseCase(userRepo repository.UserRepository, log *zap.Logger) domainUseCase.UserUseCase {
	return &UserUseCaseImpl{
		UserRepo: userRepo,
		Log:      log,
	}
}

func (u *UserUseCaseImpl) GetProfile(ctx context.Context, id string) (*entity.User, error) {
	return u.UserRepo.FindByID(ctx, id)
}

func (u *UserUseCaseImpl) GetAllUsers(ctx context.Context, page, pageSize int) ([]*entity.User, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return u.UserRepo.FindAll(ctx, offset, pageSize)
}

func (u *UserUseCaseImpl) UpdateUser(ctx context.Context, user *entity.User) error {
	return u.UserRepo.Update(ctx, user)
}

func (u *UserUseCaseImpl) DeleteUser(ctx context.Context, id string) error {
	return u.UserRepo.Delete(ctx, id)
}
