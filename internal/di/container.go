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

	depositAccount            *accounts2.DepositAccountUseCase
	postDepositAccountHandler *accounts.POSTDepositAccountHandler
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

func (c *Container) POSTDepositAccountHandler() *accounts.POSTDepositAccountHandler {
	if c.postDepositAccountHandler == nil {
		c.postDepositAccountHandler = accounts.NewPOSTDepositAccountHandler(c.DepositAccount())
	}
	return c.postDepositAccountHandler
}

func (c *Container) DepositAccount() *accounts2.DepositAccountUseCase {
	if c.depositAccount == nil {
		c.depositAccount = accounts2.NewDepositAccountUseCase(c.AccountsRepository())
	}
	return c.depositAccount
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

	publicRouter := router.PathPrefix("/api").Subrouter()
	publicRouter.Handle("/users", c.POSTUserHandler()).Methods(http.MethodPost)
	publicRouter.Handle("/tokens", c.POSTTokenHandler()).Methods(http.MethodPost)

	securedRouter := router.PathPrefix("/api").Subrouter()
	securedRouter.Use(middleware.AuthMiddleware)
	securedRouter.Handle("/accounts", c.POSTAccountHandler()).Methods(http.MethodPost)
	securedRouter.Handle("/deposit", c.POSTDepositAccountHandler()).Methods(http.MethodPost)
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
