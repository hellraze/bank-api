package domain

import (
	"context"
	"github.com/gofrs/uuid"
)

type User struct {
	id           uuid.UUID
	login        string
	hashPassword []byte
}

func (u *User) ID() uuid.UUID    { return u.id }
func (u *User) Login() string    { return u.login }
func (u *User) Password() []byte { return u.hashPassword }

func NewUser(id uuid.UUID, login string, password []byte) *User {
	return &User{
		id:           id,
		login:        login,
		hashPassword: password,
	}
}

type UserRepository interface {
	Save(ctx context.Context, u *User) error
	FindByLogin(ctx context.Context, login string) (*User, error)
}
