package users

import (
	"bank-api/internal/usecase/tokens"
	"bank-api/internal/usecase/users"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type POSTUserHandler struct {
	useCase     *users.CreateUserUseCase
	readUseCase *users.ReadUserUseCase
}

type POSTUserRequest struct {
	Login    string
	Password []byte
}

type POSTUserResponse struct {
	SignedToken string `json:"signedToken"`
}

func NewPOSTUserHandler(useCase *users.CreateUserUseCase, readUseCase *users.ReadUserUseCase) *POSTUserHandler {
	return &POSTUserHandler{
		useCase:     useCase,
		readUseCase: readUseCase,
	}
}

func (handler *POSTUserHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var body POSTUserRequest
	secretKey := os.Getenv("SECRET_KEY")
	ctx := request.Context()
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	command := &users.CreateUserCommand{
		Login:    body.Login,
		Password: body.Password,
	}
	readCommand := &users.ReadUserCommand{
		Login: body.Login,
	}

	user, err := handler.readUseCase.ReadUserHandler(ctx, readCommand)
	if err != nil {
		user, err = handler.useCase.CreateUserHandler(ctx, command)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		signedToken, err := tokens.NewSignedToken(user.ID(), []byte(secretKey))
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
		err = errors.New("users already exists")
		http.Error(writer, err.Error(), http.StatusConflict)
	}

}
