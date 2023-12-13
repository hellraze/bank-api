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
	useCase *accounts.CreateAccountUseCase
}

type POSTAccountRequest struct {
	Name string `json:"name"`
}

func NewPOSTAccountHandler(useCase *accounts.CreateAccountUseCase) *POSTAccountHandler {
	return &POSTAccountHandler{
		useCase: useCase,
	}
}

type POSTAccountResponse struct {
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
			command := &accounts.CreateAccountCommand{
				UserID: userID,
				Name:   body.Name,
			}
			_, err = handler.useCase.CreateAccountHandler(request.Context(), command)
		} else {
			fmt.Println("Идентификатор не найден или не является строкой")
		}
	} else {
		fmt.Println("Неверный токен или отсутствуют данные (claims)")
	}
}
