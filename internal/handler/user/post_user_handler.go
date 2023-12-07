package user

import (
	"bank-api/internal/usecase"
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
)

type POSTUserHandler struct {
	useCase *usecase.CreateUserUseCase
}

type POSTUserRequest struct {
	Login    string
	Password string
}

type POSTUserResponse struct {
	ID uuid.UUID
}

func NewPOSTUserHandler(useCase *usecase.CreateUserUseCase) *POSTUserHandler {
	return &POSTUserHandler{
		useCase: useCase,
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

	user, err := handler.useCase.CreateUserHandler(ctx, command)
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
}
