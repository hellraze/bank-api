package di

import (
	"bank-api/internal/domain"
	"bank-api/internal/handler/accounts"
	"bank-api/internal/handler/middleware"
	"bank-api/internal/handler/tokens"
	"bank-api/internal/handler/users"
	"bank-api/internal/repository/postgres"
	accounts2 "bank-api/internal/usecase/accounts"
	users2 "bank-api/internal/usecase/users"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	router http.Handler
	pool   *pgxpool.Pool

	usersRepository  *postgres.UserRepository
	createUsers      *users2.CreateUserUseCase
	postUsersHandler *users.POSTUserHandler

	readUsers       *users2.ReadUserUseCase
	getUsersHandler *tokens.POSTTokenHandler

	createAccount      *accounts2.CreateAccountUseCase
	accountHandler     *accounts.POSTAccountHandler
	accountsRepository *postgres.AccountRepository
}

func NewContainer(ctx context.Context) *Container {
	pool, err := CreateConnection(ctx)
	if err != nil {
		fmt.Printf("error: %w", err)
	}
	return &Container{
		pool: pool,
	}
}

func (c *Container) Close() {
	c.pool.Close()
}

func (c *Container) POSTUserHandler() *users.POSTUserHandler {
	if c.postUsersHandler == nil {
		c.postUsersHandler = users.NewPOSTUserHandler(c.CreateUsers(), c.ReadUsers())
	}
	return c.postUsersHandler
}

func (c *Container) CreateUsers() *users2.CreateUserUseCase {
	if c.createUsers == nil {
		c.createUsers = users2.NewCreateUserUseCase(c.UsersRepository())
	}
	return c.createUsers
}

func (c *Container) POSTTokenHandler() *tokens.POSTTokenHandler {
	if c.getUsersHandler == nil {
		c.getUsersHandler = tokens.NewPOSTTokenHandler(c.ReadUsers())
	}
	return c.getUsersHandler
}

func (c *Container) ReadUsers() *users2.ReadUserUseCase {
	if c.readUsers == nil {
		c.readUsers = users2.NewReadUserUseCase(c.UsersRepository())
	}
	return c.readUsers
}

func (c *Container) UsersRepository() domain.UserRepository {
	if c.usersRepository == nil {
		c.usersRepository = postgres.NewUserRepository(c.pool)
	}
	return c.usersRepository
}

func (c *Container) POSTAccountHandler() *accounts.POSTAccountHandler {
	if c.accountHandler == nil {
		c.accountHandler = accounts.NewPOSTAccountHandler(c.CreateAccount())
	}
	return c.accountHandler
}

func (c *Container) CreateAccount() *accounts2.CreateAccountUseCase {
	if c.createAccount == nil {
		c.createAccount = accounts2.NewCreateAccountUseCase(c.AccountsRepository())
	}
	return c.createAccount
}

func (c *Container) AccountsRepository() domain.AccountRepository {
	if c.accountsRepository == nil {
		c.accountsRepository = postgres.NewAccountRepository(c.pool)
	}
	return c.accountsRepository
}

func (c *Container) HTTPRouter() http.Handler {
	if c.router != nil {
		return c.router
	}
	router := mux.NewRouter()
	router.Use(middleware.Recover)
	router.Handle("/api/users", c.POSTUserHandler()).Methods(http.MethodPost)
	router.Handle("/api/tokens", c.POSTTokenHandler()).Methods(http.MethodPost)
	router.Handle("/api/accounts", middleware.AuthMiddleware(c.POSTAccountHandler())).Methods(http.MethodPost)
	c.router = router
	return c.router

}

func CreateConnection(ctx context.Context) (*pgxpool.Pool, error) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(".env file not found")
	}
	dns := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(ctx, dns)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return pool, err
}
