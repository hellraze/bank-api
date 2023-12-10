package di

import (
	"bank-api/internal/domain"
	"bank-api/internal/handler/user"
	"bank-api/internal/repository/postgres"
	"bank-api/internal/usecase"
	"bank-api/internal/usecase/middleware"
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
	Pool   *pgxpool.Pool

	usersRepository  *postgres.UserRepository
	createUsers      *usecase.CreateUserUseCase
	postUsersHandler *user.POSTUserHandler

	readUsers       *usecase.ReadUserUseCase
	getUsersHandler *user.POSTTokenHandler
}

func NewContainer(ctx context.Context) *Container {
	pool, err := CreateConnection(ctx)
	if err != nil {
		fmt.Printf("error: %w", err)
	}
	return &Container{
		Pool: pool,
	}
}

func (c *Container) POSTUserHandler() *user.POSTUserHandler {
	if c.postUsersHandler == nil {
		c.postUsersHandler = user.NewPOSTUserHandler(c.CreateUsers(), c.ReadUsers())
	}
	return c.postUsersHandler
}

func (c *Container) CreateUsers() *usecase.CreateUserUseCase {
	if c.createUsers == nil {
		c.createUsers = usecase.NewCreateUserUseCase(c.UsersRepository())
	}
	return c.createUsers
}

func (c *Container) POSTTokenHandler() *user.POSTTokenHandler {
	if c.getUsersHandler == nil {
		c.getUsersHandler = user.NewPOSTTokenHandler(c.ReadUsers())
	}
	return c.getUsersHandler
}

func (c *Container) ReadUsers() *usecase.ReadUserUseCase {
	if c.readUsers == nil {
		c.readUsers = usecase.NewReadUserUseCase(c.UsersRepository())
	}
	return c.readUsers
}

func (c *Container) UsersRepository() domain.UserRepository {
	if c.usersRepository == nil {
		c.usersRepository = postgres.NewUserRepository(c.Pool)
	}
	return c.usersRepository
}

func (c *Container) HTTPRouter() http.Handler {
	if c.router != nil {
		return c.router
	}
	router := mux.NewRouter()
	router.Use(middleware.Recover)
	router.Handle("/users", c.POSTUserHandler()).Methods(http.MethodPost)
	router.Handle("/token", c.POSTTokenHandler()).Methods(http.MethodPost)
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
