package database

import (
	"context"

	"gorm.io/gorm"
)

type contextKey string

const dbKey contextKey = "db_transaction"

// GetDBFromContext returns the DB from context, or the default DB if not found
func GetDBFromContext(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(dbKey).(*gorm.DB); ok {
		return tx
	}
	return defaultDB.WithContext(ctx)
}

// TransactionManager implementation using GORM
type TransactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

func (tm *TransactionManager) Execute(ctx context.Context, fn func(txCtx context.Context) error) error {
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Inject the transaction into the context
		txCtx := context.WithValue(ctx, dbKey, tx)
		return fn(txCtx)
	})
}
