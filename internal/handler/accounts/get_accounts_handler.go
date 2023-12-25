package accounts

import (
	"bank-api/internal/domain"
	"bank-api/internal/handler/middleware"
	"bank-api/internal/usecase/accounts"
	"encoding/json"
	"github.com/gofrs/uuid"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type GETUserAccountsHandler struct {
	useCase *accounts.ReadUserAccountsUseCase
}

type GETUserAccountsRequest struct {
	Name   string `json:"name"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

func NewGETUserAccountsHandler(useCase *accounts.ReadUserAccountsUseCase) *GETUserAccountsHandler {
	return &GETUserAccountsHandler{
		useCase: useCase,
	}
}

type GETUserAccountsResponse struct { //добавить сериалайзер
	Accounts []domain.Account
}

func (handler *GETUserAccountsHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	var body GETUserAccountsRequest
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	token := middleware.TokenFromContext(request.Context())

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if id, exists := claims["ID"].(string); exists {
			userID, _ := uuid.FromString(id)
			readCommand := &accounts.ReadUserAccountsCommand{
				Name:   body.Name,
				UserID: userID,
				Offset: body.Offset,
				Limit:  body.Limit,
			}
			accountList, err := handler.useCase.ReadUserAccountsHandler(request.Context(), readCommand)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}

			response := &GETUserAccountsResponse{
				Accounts: accountList,
			}
			err = json.NewEncoder(writer).Encode(response)
		}
	} else {
		http.Error(writer, err.Error(), http.StatusForbidden)
	}
}
