package user

import (
	"bank-api/internal/usecase"
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
)

type GETUserHandler struct {
	useCase *usecase.ReadUserUseCase
}

type GETUserRequest struct {
	Login string
}

type GETUserResponse struct {
	ID    uuid.UUID
	Login string
}

func NewGETUserHandler(useCase *usecase.ReadUserUseCase) *GETUserHandler {
	return &GETUserHandler{
		useCase: useCase,
	}
}

func (handler *GETUserHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query().Get("login")
	login := query
	ctx := request.Context()
	command := &usecase.ReadUserCommand{
		Login: login,
	}

	user, err := handler.useCase.ReadUserHandler(ctx, command)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	response := &GETUserResponse{
		ID:    user.ID(),
		Login: user.Login(),
	}
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
