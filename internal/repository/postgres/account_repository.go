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
		"id":      account.ID(),
		"name":    account.Name(),
		"userID":  account.UserID(),
		"balance": account.Balance(),
	}
	_, err := accountRepository.Pool.Exec(ctx, "INSERT INTO bank.account(account_id, name, balance, user_id) VALUES(@id, @name, @balance, @userID)", args)
	return err
}

func (accountRepository *AccountRepository) FindAccountByName(ctx context.Context, name string, userID uuid.UUID) (*domain.Account, error) {
	var (
		id      uuid.UUID
		balance int
	)
	err := accountRepository.Pool.QueryRow(ctx, "SELECT * FROM bank.account WHERE name=($1)AND user_id = ($2)", name, userID).Scan(&id, &name, &userID, &balance)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	account := domain.NewAccount(id, name, userID)
	fmt.Println(account)
	return account, err
}

func (accountRepository *AccountRepository) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	var (
		userID uuid.UUID
		name   string
	)
	err := accountRepository.Pool.QueryRow(ctx, "SELECT account_id, name, user_id FROM bank.account WHERE user_id = ($1) FOR UPDATE", id).Scan(&id, &name, &userID)
	if err != nil {
		return nil, err
	}
	account := domain.NewAccount(id, name, userID)
	return account, nil
}

func (accountRepository *AccountRepository) UpdateAccountBalance(ctx context.Context, id uuid.UUID, deposit int) error {
	account, err := accountRepository.FindByIDForUpdate(ctx, id)
	if err != nil {
		return err
	}
	account.Deposit(deposit)
	_, err = accountRepository.Pool.Exec(ctx, "UPDATE bank.account SET balance = $2 WHERE account_id=$1 FOR UPDATE", account.UserID(), account.Balance())
	if err != nil {
		return err
	}
	return err
}
