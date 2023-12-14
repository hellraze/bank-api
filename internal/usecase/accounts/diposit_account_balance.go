package accounts

import (
	"bank-api/internal/domain"
	"context"
	"github.com/gofrs/uuid"
)

type DepositAccountUseCase struct {
	accountRepository domain.AccountRepository
}

func NewDepositAccountUseCase(accountRepository domain.AccountRepository) *DepositAccountUseCase {
	return &DepositAccountUseCase{
		accountRepository: accountRepository,
	}
}

type DepositAccountCommand struct {
	ID      uuid.UUID
	Deposit int
}

func (useCase DepositAccountUseCase) DepositAccountHandler(ctx context.Context, command *DepositAccountCommand) error {
	err := useCase.accountRepository.UpdateAccountBalance(ctx, command.ID, command.Deposit)
	return err
}
