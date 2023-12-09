package user

import (
	"bank-api/internal/usecase"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gofrs/uuid"
)

type POSTUserHandler struct {
	useCase     *usecase.CreateUserUseCase
	readUseCase *usecase.ReadUserUseCase
}

type POSTUserRequest struct {
	Login    string
	Password []byte
}

type POSTUserResponse struct {
	ID uuid.UUID
}

func NewPOSTUserHandler(useCase *usecase.CreateUserUseCase, readUseCase *usecase.ReadUserUseCase) *POSTUserHandler {
	return &POSTUserHandler{
		useCase:     useCase,
		readUseCase: readUseCase,
	}
}

func (handler *POSTUserHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var body POSTUserRequest
	ctx := request.Context()
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	command := &usecase.CreateUserCommand{
		Login:    body.Login,
		Password: body.Password,
	}
	readCommand := &usecase.ReadUserCommand{
		Login: body.Login,
	}

	user, err := handler.readUseCase.ReadUserHandler(ctx, readCommand)
	if err != nil {
		user, err = handler.useCase.CreateUserHandler(ctx, command)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		response := &POSTUserResponse{
			ID: user.ID(),
		}
		err = json.NewEncoder(writer).Encode(response)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	} else {
		err = errors.New("user already exists")
		http.Error(writer, err.Error(), http.StatusConflict)
	}

}
