package postgres

import (
	"bank-api/internal/domain"
	"bank-api/internal/pkg/persistence/postgres"
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository struct {
	pool *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool, transactionManager *postgres.PoolTransactionManager) *AccountRepository {
	return &AccountRepository{
		pool: pool,
	}
}

func (accountRepository *AccountRepository) Save(ctx context.Context, account *domain.Account) error {
	args := pgx.NamedArgs{
		"id":      account.ID(),
		"name":    account.Name(),
		"userID":  account.UserID(),
		"balance": account.Balance(),
	}
	_, err := accountRepository.pool.Exec(ctx, "INSERT INTO bank.account(account_id, name, balance, user_id) VALUES(@id, @name, @balance, @userID)", args)
	return err
}

func (accountRepository *AccountRepository) FindAccountByName(ctx context.Context, name string, userID uuid.UUID) (*domain.Account, error) {
	var (
		id      uuid.UUID
		balance int
	)
	err := accountRepository.pool.QueryRow(ctx, "SELECT * FROM bank.account WHERE name=($1)AND user_id = ($2)", name, userID).Scan(&id, &name, &userID, &balance)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	account := domain.NewAccount(id, name, balance, userID)
	fmt.Println(account)
	return account, err
}

func (accountRepository *AccountRepository) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	var (
		userID  uuid.UUID
		name    string
		balance int
	)
	err := accountRepository.pool.QueryRow(ctx, "SELECT account_id, name, user_id, balance FROM bank.account WHERE user_id = ($1) FOR UPDATE", id).Scan(&id, &name, &userID, &balance)
	if err != nil {
		return nil, err
	}
	account := domain.NewAccount(id, name, balance, userID)
	return account, nil
}

func (accountRepository *AccountRepository) UpdateAccountBalance(ctx context.Context, id uuid.UUID, deposit int) error {
	_, err := accountRepository.pool.Exec(ctx, "UPDATE bank.account SET balance = ($2) WHERE account_id=($1)", id, deposit)
	if err != nil {
		return err
	}
	return err
}

func (accountRepository *AccountRepository) FindUserAccountsILike(ctx context.Context, name string, offset int, limit int, userID uuid.UUID) ([]domain.Account, error) {
	query := `SELECT account_id, name, balance, user_id FROM bank.account WHERE user_id = $1 AND name ILIKE $2 LIMIT $3 OFFSET $4;` //squirell
	rows, err := accountRepository.pool.Query(ctx, query, userID, "%"+name+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var accounts []domain.Account

	for rows.Next() {
		var ( //отдельная структура
			accountID uuid.UUID
			name      string
			balance   int
			userID    uuid.UUID
		)
		err = rows.Scan(
			&accountID,
			&name,
			&balance,
			&userID,
		)
		if err != nil {
			return nil, err
		}
		account := domain.NewAccount(accountID, name, balance, userID)
		accounts = append(accounts, *account)
	}
	return accounts, err
}
