package usecase

import (
	"context"

	"go-restful-api/internal/domain/entity"
)

type AuthUseCase interface {
	Register(ctx context.Context, user *entity.User) error
	Login(ctx context.Context, email, password string) (string, string, string, error) // Returns AccessToken, RefreshToken, UserID, error
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
	ResetPassword(ctx context.Context, email, newPassword string) error
}
