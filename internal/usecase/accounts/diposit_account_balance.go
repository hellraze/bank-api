package accounts

import (
	"bank-api/internal/domain"
	"bank-api/internal/pkg/persistence"
	"context"
	"github.com/gofrs/uuid"
)

type DepositAccountUseCase struct {
	accountRepository  domain.AccountRepository
	transactionManager persistence.TransactionManager
}

func NewDepositAccountUseCase(accountRepository domain.AccountRepository, transactionManager persistence.TransactionManager) *DepositAccountUseCase {
	return &DepositAccountUseCase{
		accountRepository:  accountRepository,
		transactionManager: transactionManager,
	}
}

type DepositAccountCommand struct {
	UserID  uuid.UUID
	Deposit int
}

func (useCase DepositAccountUseCase) DepositAccountHandler(ctx context.Context, command *DepositAccountCommand) error {
	return useCase.transactionManager.Do(ctx, func(ctx context.Context) error {
		account, err := useCase.accountRepository.FindByIDForUpdate(ctx, command.UserID)
		if err != nil {
			return err
		}
		account.Deposit(command.Deposit)
		err = useCase.accountRepository.UpdateAccountBalance(ctx, account.UserID(), account.Balance())
		if err != nil {
			return err
		}

		err = useCase.accountRepository.UpdateAccountBalance(ctx, command.UserID, command.Deposit)
		return err
	})

}
