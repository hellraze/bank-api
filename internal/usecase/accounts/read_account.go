package accounts

import (
	"bank-api/internal/domain"
	"context"
	"github.com/gofrs/uuid"
)

type ReadAccountUseCase struct {
	accountRepository domain.AccountRepository
}

func NewReadAccountUseCase(accountRepository domain.AccountRepository) *ReadAccountUseCase {
	return &ReadAccountUseCase{
		accountRepository: accountRepository,
	}
}

type ReadAccountCommand struct {
	Name   string
	UserID uuid.UUID
}

func (useCase ReadAccountUseCase) ReadAccountHandler(ctx context.Context, command *ReadAccountCommand) (*domain.Account, error) {
	account, err := useCase.accountRepository.FindAccountByName(ctx, command.Name, command.UserID)
	return account, err
}
