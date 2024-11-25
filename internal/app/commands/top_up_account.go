package commands

import (
	"context"
	"time"

	"accounts/internal/common/decorator"
	"accounts/internal/domain/account"
)

// TopUpAccount is a command to top up an account.
type TopUpAccount struct {
	ID     account.ID
	Amount int64
}

// TopUpAccountHandler is a handler for TopUpAccount command.
type TopUpAccountHandler decorator.CommandHandler[TopUpAccount]

type updateAccountRepository interface {
	UpdateByID(ctx context.Context, id account.ID, updateFn func(acc *account.Account) (bool, error)) error
}
type topUpAccountHandler struct {
	accountRepository updateAccountRepository
}

// NewTopUpAccountHandler creates a new TopUpAccountHandler.
func NewTopUpAccountHandler(repo updateAccountRepository) TopUpAccountHandler {
	if repo == nil {
		panic("accountRepository cannot be nil")
	}

	return decorator.ApplyCommandDecorators(topUpAccountHandler{
		accountRepository: repo,
	})
}

// Handle handles the TopUpAccount command.
func (h topUpAccountHandler) Handle(ctx context.Context, cmd TopUpAccount) error {
	return h.accountRepository.UpdateByID(ctx, cmd.ID, func(acc *account.Account) (bool, error) {
		money, err := account.NewMoney(cmd.Amount, acc.BalanceCurrency())
		if err != nil {
			return false, err
		}
		err = acc.TopUp(money, time.Now())
		if err != nil {
			return false, err
		}
		return true, nil
	})
}
