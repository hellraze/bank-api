package postgres

import (
	"bank-api/internal/domain"
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository struct {
	Pool *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{
		Pool: pool,
	}
}

func (accountRepository *AccountRepository) Save(ctx context.Context, account *domain.Account) error {
	args := pgx.NamedArgs{
		"id":     account.ID(),
		"name":   account.Name(),
		"userID": account.UserID(),
	}
	_, err := accountRepository.Pool.Exec(ctx, "INSERT INTO bank.account(account_id, name, user_id) VALUES(@id, @name, @userID)", args)
	return err
}

func (accountRepository *AccountRepository) FindAccount(ctx context.Context, name string, userID uuid.UUID) (*domain.Account, error) {
	var (
		id uuid.UUID
	)
	err := accountRepository.Pool.QueryRow(ctx, "SELECT * FROM bank.account WHERE name=($1)", name).Scan(&id, &name, &userID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	account := domain.NewAccount(id, name, userID)
	return account, err
}
