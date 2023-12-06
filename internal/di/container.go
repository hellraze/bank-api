package di

import (
	"bank-api/internal/domain"
	"bank-api/internal/handler"
	"bank-api/internal/repository/postgres"
	"bank-api/internal/usecase"
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	router           http.Handler
	Pool             *pgxpool.Pool
	usersRepository  *postgres.UserRepository
	createUsers      *usecase.CreateUserUseCase
	postUsersHandler *handler.POSTUserHandler
}

func NewContainer(ctx context.Context) *Container {
	pool, err := postgres.CreateConnection(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return &Container{
		Pool: pool,
	}
}

func (c *Container) POSTUserHandler() *handler.POSTUserHandler {
	if c.postUsersHandler == nil {
		c.postUsersHandler = handler.NewPOSTUserHandler(c.CreateUsers())
	}
	return c.postUsersHandler
}

func (c *Container) CreateUsers() *usecase.CreateUserUseCase {
	if c.createUsers == nil {
		c.createUsers = usecase.NewCreateUserUseCase(c.UsersRepository())
	}
	return c.createUsers
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
	c.router = router
	return c.router

}
