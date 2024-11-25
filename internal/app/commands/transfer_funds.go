package commands

import (
	"accounts/internal/common/decorator"
	"accounts/internal/domain"
	"accounts/internal/domain/account"
	"context"
	"time"
)

// TransferFunds is a command to transfer funds from one account to another.
type TransferFunds struct {
	From   account.ID
	To     account.ID
	Amount int64
}

// TransferFundsHandler is a handler for TransferFunds command.
type TransferFundsHandler decorator.CommandHandler[TransferFunds]

type transferFundsRepository interface {
	UpdateTwoAccounts(ctx context.Context, id1, id2 account.ID, updateFn func(acc1, acc2 *account.Account) (bool, error)) error
}

type transferFundsHandler struct {
	transferFundsRepository transferFundsRepository
}

// NewTransferFundsHandler creates a new TransferFundsHandler.
func NewTransferFundsHandler(repo transferFundsRepository) TransferFundsHandler {
	if repo == nil {
		panic("transferFundsRepository cannot be nil")
	}

	return decorator.ApplyCommandDecorators(transferFundsHandler{
		transferFundsRepository: repo,
	})
}

// Handle handles the TransferFunds command.
func (h transferFundsHandler) Handle(ctx context.Context, cmd TransferFunds) error {
	return h.transferFundsRepository.UpdateTwoAccounts(ctx, cmd.From, cmd.To, func(acc1, acc2 *account.Account) (bool, error) {
		money, err := account.NewMoney(cmd.Amount, acc1.BalanceCurrency())
		if err != nil {
			return false, err
		}
		err = domain.Transfer(acc1, acc2, money, time.Now())
		if err != nil {
			return false, err
		}
		return true, nil
	})
}
