package domain

import (
	"context"
	"github.com/gofrs/uuid"
)

type Account struct {
	id     uuid.UUID
	name   string
	userID uuid.UUID
}

func (a *Account) ID() uuid.UUID     { return a.id }
func (a *Account) Name() string      { return a.name }
func (a *Account) UserID() uuid.UUID { return a.userID }

func NewAccount(id uuid.UUID, name string, userID uuid.UUID) *Account {
	return &Account{
		id:     id,
		name:   name,
		userID: userID,
	}
}

type AccountRepository interface {
	Save(ctx context.Context, account *Account) error
	FindAccount(ctx context.Context, name string, userID uuid.UUID) (*Account, error)
}
