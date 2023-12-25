package accounts

import (
	"bank-api/internal/domain"
	"context"
	"github.com/gofrs/uuid"
)

type ReadUserAccountsUseCase struct {
	accountRepository domain.AccountRepository
}

func NewReadUserAccountsUseCase(accountRepository domain.AccountRepository) *ReadUserAccountsUseCase {
	return &ReadUserAccountsUseCase{
		accountRepository: accountRepository,
	}
}

type ReadUserAccountsCommand struct {
	Name   string
	UserID uuid.UUID
	Offset int
	Limit  int
}

func (useCase ReadUserAccountsUseCase) ReadUserAccountsHandler(ctx context.Context, command *ReadUserAccountsCommand) ([]domain.Account, error) {
	accounts, err := useCase.accountRepository.FindUserAccountsILike(ctx, command.Name, command.Offset, command.Limit, command.UserID)
	return accounts, err
}
