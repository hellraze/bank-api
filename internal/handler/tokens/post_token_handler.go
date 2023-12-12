package tokens

import (
	"bank-api/internal/usecase/tokens"
	"bank-api/internal/usecase/users"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
)

type POSTTokenHandler struct {
	useCase *users.ReadUserUseCase
}

type POSTTokenResponse struct {
	Token string `json:"token"`
}

type POSTTokenRequest struct {
	Login    string
	Password string
}

func NewPOSTTokenHandler(useCase *users.ReadUserUseCase) *POSTTokenHandler {
	return &POSTTokenHandler{
		useCase: useCase,
	}
}

func (handler *POSTTokenHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	secretKey := os.Getenv("SECRET_KEY")
	ctx := request.Context()
	var body POSTTokenRequest
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	password, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	command := &users.ReadUserCommand{
		Login: body.Login,
	}

	user, err := handler.useCase.ReadUserHandler(ctx, command)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	err = bcrypt.CompareHashAndPassword(password, user.Password())

	if err != nil {
		signedToken, err := tokens.NewSignedToken(user.ID(), []byte(secretKey))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		response := &POSTTokenResponse{
			Token: signedToken,
		}
		err = json.NewEncoder(writer).Encode(response)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	} else {
		err = errors.New("wrong password")
		http.Error(writer, err.Error(), http.StatusConflict)
	}
}
