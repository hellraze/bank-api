package di

import (
	"bank-api/internal/domain"
	"bank-api/internal/handler/user"
	"bank-api/internal/repository/postgres"
	"bank-api/internal/usecase"
	"context"
	"github.com/joho/godotenv"
	"log"
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
	getUsersHandler *user.GETUserHandler
}

func NewContainer(ctx context.Context) *Container {
	pool, err := CreateConnection(ctx)
	if err != nil {
		log.Fatal(err)
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

func (c *Container) GETUserHandler() *user.GETUserHandler {
	if c.getUsersHandler == nil {
		c.getUsersHandler = user.NewGETUserHandler(c.ReadUsers())
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
	router.Handle("/users", c.POSTUserHandler()).Methods(http.MethodPost)
	router.Handle("/users", c.GETUserHandler()).Methods(http.MethodGet)
	c.router = router
	return c.router

}

func CreateConnection(ctx context.Context) (*pgxpool.Pool, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env file not found")
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
