package accounts

import (
	"bank-api/internal/domain"
	"bank-api/internal/handler/middleware"
	"bank-api/internal/usecase/accounts"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type POSTAccountHandler struct {
	useCase     *accounts.CreateAccountUseCase
	readAccount *accounts.ReadAccountUseCase
}

type POSTAccountRequest struct {
	Name string `json:"name"`
}

func NewPOSTAccountHandler(useCase *accounts.CreateAccountUseCase, readAccount *accounts.ReadAccountUseCase) *POSTAccountHandler {
	return &POSTAccountHandler{
		useCase:     useCase,
		readAccount: readAccount,
	}
}

type POSTAccountResponse struct { // добавить сериалайзер
	Account *domain.Account
}

func (handler *POSTAccountHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	var body POSTAccountRequest
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	token := middleware.TokenFromContext(request.Context())

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if id, exists := claims["ID"].(string); exists {
			userID, _ := uuid.FromString(id)
			readCommand := &accounts.ReadAccountCommand{
				UserID: userID,
				Name:   body.Name,
			}
			account, err := handler.readAccount.ReadAccountHandler(request.Context(), readCommand)
			fmt.Println(account)
			if err != nil {
				command := &accounts.CreateAccountCommand{
					UserID: userID,
					Name:   body.Name,
				}
				account, err := handler.useCase.CreateAccountHandler(request.Context(), command)
				if err != nil {
					http.Error(writer, err.Error(), http.StatusInternalServerError)
				}
				response := &POSTAccountResponse{
					Account: account,
				}
				err = json.NewEncoder(writer).Encode(response)
			} else {
				fmt.Println("Идентификатор не найден или не является строкой")
			}
		} else {
			fmt.Println("Неверный токен или отсутствуют данные (claims)")
		}

	}
}
