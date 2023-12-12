package accounts

import (
	"bank-api/internal/domain"
	"bank-api/pkg"
	"context"
	"github.com/gofrs/uuid"
)

type CreateAccountUseCase struct {
	accountRepository domain.AccountRepository
}

func NewCreateAccountUseCase(accountRepository domain.AccountRepository) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		accountRepository: accountRepository,
	}
}

type CreateAccountCommand struct {
	UserID uuid.UUID
	Name   string
}

func (useCase *CreateAccountUseCase) CreateAccountHandler(ctx context.Context, command *CreateAccountCommand) (*domain.Account, error) {
	account := domain.NewAccount(pkg.GenerateID(), command.Name, command.UserID)
	err := useCase.accountRepository.Save(ctx, account)
	return account, err
}
