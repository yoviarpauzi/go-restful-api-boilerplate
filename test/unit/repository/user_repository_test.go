package unit

import (
	"context"
	"regexp"
	"testing"
	"time"

	"go-restful-api/internal/domain/entity"
	"go-restful-api/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm: %v", err)
	}

	return gormDB, mock
}

func TestUserRepository_Create(t *testing.T) {
	db, mock := setupTestDB(t)
	log := zap.NewNop()
	repo := repository.NewUserRepository(db, log)

	userID := uuid.New()
	user := &entity.User{
		ID:        userID,
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), user)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByID(t *testing.T) {
	db, mock := setupTestDB(t)
	log := zap.NewNop()
	repo := repository.NewUserRepository(db, log)

	userID := uuid.New()
	expectedUser := &entity.User{
		ID:    userID,
		Name:  "Test User",
		Email: "test@example.com",
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email)

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT \$2`).
		WithArgs(userID.String(), 1).
		WillReturnRows(rows)

	user, err := repo.FindByID(context.Background(), userID.String())

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}
