package accounts

import (
	"bank-api/internal/handler/middleware"
	"bank-api/internal/usecase/accounts"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type POSTDepositAccountHandler struct {
	useCase *accounts.DepositAccountUseCase
}

type POSTDepositAccountRequest struct {
	Deposit int `json:"deposit"`
}

func NewPOSTDepositAccountHandler(useCase *accounts.DepositAccountUseCase) *POSTDepositAccountHandler {
	return &POSTDepositAccountHandler{
		useCase: useCase,
	}
}

func (handler *POSTDepositAccountHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	var body POSTDepositAccountRequest
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	token := middleware.TokenFromContext(request.Context())

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if id, exists := claims["ID"].(string); exists {
			id, _ := uuid.FromString(id)
			depositCommand := &accounts.DepositAccountCommand{
				UserID:  id,
				Deposit: body.Deposit,
			}
			err := handler.useCase.DepositAccountHandler(request.Context(), depositCommand)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
			} else {
				fmt.Println("Идентификатор не найден или не является строкой")
			}
		} else {
			fmt.Println("Неверный токен или отсутствуют данные (claims)")
		}

	}
}
