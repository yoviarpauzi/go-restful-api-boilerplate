package usecase

import (
	"context"
	"errors"

	"go-restful-api/internal/domain/entity"
	"go-restful-api/internal/domain/repository"
	domainUseCase "go-restful-api/internal/domain/usecase"
	"go-restful-api/internal/infrastructure/token"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCaseImpl struct {
	UserRepo       repository.UserRepository
	TokenGenerator *token.TokenGenerator
	Log            *zap.Logger
	TxManager      repository.TransactionManager
}

func NewAuthUseCase(userRepo repository.UserRepository, generator *token.TokenGenerator, log *zap.Logger, txManager repository.TransactionManager) domainUseCase.AuthUseCase {
	return &AuthUseCaseImpl{
		UserRepo:       userRepo,
		TokenGenerator: generator,
		Log:            log,
		TxManager:      txManager,
	}
}

func (a *AuthUseCaseImpl) Register(ctx context.Context, user *entity.User) error {
	return a.TxManager.Execute(ctx, func(txCtx context.Context) error {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
		return a.UserRepo.Create(txCtx, user)
	})
}

func (a *AuthUseCaseImpl) Login(ctx context.Context, email, password string) (string, string, string, error) {
	user, err := a.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", "", errors.New("invalid email or password")
	}

	accessToken, err := a.TokenGenerator.GenerateAccessToken(user.ID.String())
	if err != nil {
		return "", "", "", err
	}

	refreshToken, err := a.TokenGenerator.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return "", "", "", err
	}

	return accessToken, refreshToken, user.ID.String(), nil
}

func (a *AuthUseCaseImpl) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	return a.TxManager.Execute(ctx, func(txCtx context.Context) error {
		user, err := a.UserRepo.FindByID(txCtx, userID)
		if err != nil {
			return err
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
		if err != nil {
			return errors.New("old password does not match")
		}

		newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(newHashedPassword)
		return a.UserRepo.Update(txCtx, user)
	})
}

func (a *AuthUseCaseImpl) ResetPassword(ctx context.Context, email, newPassword string) error {
	return a.TxManager.Execute(ctx, func(txCtx context.Context) error {
		user, err := a.UserRepo.FindByEmail(txCtx, email)
		if err != nil {
			return err
		}

		newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(newHashedPassword)
		return a.UserRepo.Update(txCtx, user)
	})
}
