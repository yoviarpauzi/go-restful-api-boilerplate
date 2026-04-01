package repository

import "context"

type TransactionManager interface {
	Execute(ctx context.Context, fn func(txCtx context.Context) error) error
}
