package user

import (
	"bank-api/internal/usecase"
	"encoding/json"
	"errors"
	"net/http"
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
	SignedToken string `json:"signedToken"`
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
		signedToken, err := usecase.NewSignedToken(user.ID(), []byte("asd"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		response := &POSTUserResponse{
			SignedToken: signedToken,
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
