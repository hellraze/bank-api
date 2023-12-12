package accounts

import (
	"bank-api/internal/domain"
	"bank-api/internal/handler/middleware"
	"bank-api/internal/usecase/accounts"
	"fmt"
	"net/http"
)

type POSTAccountHandler struct {
	useCase *accounts.CreateAccountUseCase
}

type POSTAccountRequest struct {
	Token string
	Name  string
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
	name := middleware.AccountNameFromContext(request.Context())
	fmt.Println(name)
}
