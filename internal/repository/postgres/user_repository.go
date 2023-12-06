package postgres

import (
	"bank-api/internal/domain"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	Pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		Pool: pool,
	}
}

func (userRepository *UserRepository) Save(ctx context.Context, user *domain.User) error {
	args := pgx.NamedArgs{
		"id":       user.ID(),
		"login":    user.Login(),
		"password": user.Password(),
	}
	_, err := userRepository.Pool.Query(ctx, "INSERT INTO bank.user(id, login, hash_password) VALUES(@id, @login, @password)", args)
	return err
}
