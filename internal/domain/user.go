package domain

import (
	"context"
	"github.com/gofrs/uuid"
)

type User struct {
	id           uuid.UUID
	login        string
	hashPassword string
}

func (u *User) ID() uuid.UUID    { return u.id }
func (u *User) Login() string    { return u.login }
func (u *User) Password() string { return u.hashPassword }

func NewUser(login string, password string) *User {
	return &User{
		id:           uuid.Must(uuid.NewV7()),
		login:        login,
		hashPassword: password,
	}
}

type UserRepository interface {
	Save(ctx context.Context, u *User) error
}
