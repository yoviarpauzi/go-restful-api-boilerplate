package usecase

import (
	"context"

	"go-restful-api/internal/domain/entity"
)

type UserUseCase interface {
	GetProfile(ctx context.Context, id string) (*entity.User, error)
	GetAllUsers(ctx context.Context, page, pageSize int) ([]*entity.User, int64, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id string) error
}
