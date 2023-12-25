package accounts

import (
	"bank-api/internal/domain"
	postgres2 "bank-api/internal/pkg/persistence/postgres"
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
	UserID  uuid.UUID
	Deposit int
}

func (useCase DepositAccountUseCase) DepositAccountHandler(ctx context.Context, command *DepositAccountCommand) error {
	err := postgres2.PoolTransactionManager.Do(ctx, useCase.accountRepository.UpdateAccountBalance(ctx, command.UserID, command.Deposit))
	return err
}
