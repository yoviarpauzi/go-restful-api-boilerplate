package repository

import (
	"context"

	"go-restful-api/internal/domain/entity"
	domainRepository "go-restful-api/internal/domain/repository"
	"go-restful-api/internal/infrastructure/database"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func NewUserRepository(db *gorm.DB, log *zap.Logger) domainRepository.UserRepository {
	return &UserRepositoryImpl{
		DB:  db,
		Log: log,
	}
}

func (r *UserRepositoryImpl) getDB(ctx context.Context) *gorm.DB {
	return database.GetDBFromContext(ctx, r.DB)
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	return r.getDB(ctx).Create(user).Error
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	return r.getDB(ctx).Save(user).Error
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.getDB(ctx).Delete(&entity.User{}, "id = ?", id).Error
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	err := r.getDB(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.getDB(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindAll(ctx context.Context, offset, limit int) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	// Get total count
	err := r.getDB(ctx).Model(&entity.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated data
	err = r.getDB(ctx).Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, total, err
	}
	return users, total, nil
}
