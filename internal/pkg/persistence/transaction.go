package persistence

import "context"

type TransactionManager interface {
	Do(ctx context.Context, f func(ctx context.Context) error) error
}

type MockTransactionManager struct {
}

func (m MockTransactionManager) Do(ctx context.Context, f func(ctx context.Context) error) error {
	return f(ctx)
}
