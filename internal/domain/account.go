package domain

import (
	"context"
	"github.com/gofrs/uuid"
)

type Account struct {
	id      uuid.UUID
	name    string
	userID  uuid.UUID
	balance int
}

func (a *Account) ID() uuid.UUID       { return a.id }
func (a *Account) Name() string        { return a.name }
func (a *Account) UserID() uuid.UUID   { return a.userID }
func (a *Account) Balance() int        { return a.balance }
func (a *Account) Deposit(deposit int) { a.balance += deposit }

func NewAccount(id uuid.UUID, name string, userID uuid.UUID) *Account {
	return &Account{
		id:      id,
		name:    name,
		userID:  userID,
		balance: 0,
	}
}

type AccountRepository interface {
	Save(ctx context.Context, account *Account) error
	FindAccountByName(ctx context.Context, name string, userID uuid.UUID) (*Account, error)
	UpdateAccountBalance(ctx context.Context, id uuid.UUID, deposit int) error
}
