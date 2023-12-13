package postgres

import (
	"bank-api/internal/domain"
	"context"
	"fmt"
	"github.com/gofrs/uuid"
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
	_, err := userRepository.Pool.Exec(ctx, "INSERT INTO bank.user(user_id, login, hash_password) VALUES(@id, @login, @password)", args)
	return err
}

func (userRepository *UserRepository) FindByLogin(ctx context.Context, login string) (*domain.User, error) {
	var (
		id       uuid.UUID
		password []byte
	)
	err := userRepository.Pool.QueryRow(ctx, "SELECT * FROM bank.user WHERE login=($1)", login).Scan(&id, &login, &password)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	user := domain.NewUser(id, login, password)
	return user, nil
}
