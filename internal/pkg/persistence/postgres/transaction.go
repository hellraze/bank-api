package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PoolTransactionManager struct {
	Connection *pgxpool.Pool
}

func (p *PoolTransactionManager) Do(ctx context.Context, f func(ctx context.Context) error) error {
	pool := NewPoolConnection(p.Connection)
	tx, err := pool.pool.Begin(ctx)
	if err != nil {
		return err
	}
	txContext := context.WithValue(ctx, "transaction", tx)
	if err = f(txContext); err != nil {
		tx.Rollback(ctx)
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	return err
}

type PoolConnection struct {
	pool *pgxpool.Pool
}

func NewPoolConnection(pool *pgxpool.Pool) *PoolConnection {
	return &PoolConnection{pool: pool}
}
func (c *PoolConnection) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	if tx, ok := ctx.Value("transaction").(pgx.Tx); ok {
		return tx.Query(ctx, sql, args...)
	}

	// No transaction in the context, use the pool
	return c.pool.Query(ctx, sql, args...)
}
func (c *PoolConnection) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	if tx, ok := ctx.Value("transaction").(pgx.Tx); ok {
		return tx.Exec(ctx, sql, arguments...)
	}
	return c.pool.Exec(ctx, sql, arguments...)
}
func (c *PoolConnection) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if tx, ok := ctx.Value("transaction").(pgx.Tx); ok {
		return tx.QueryRow(ctx, sql, args...)
	}
	return c.pool.QueryRow(ctx, sql, args...)
}
