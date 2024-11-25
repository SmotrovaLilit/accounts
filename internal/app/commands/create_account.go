package commands

import (
	"context"
	"time"

	"accounts/internal/common/decorator"
	"accounts/internal/domain/account"
	"accounts/internal/domain/owner"
)

// CreateAccount is a command to create an account.
type CreateAccount struct {
	ID       account.ID
	OwnerID  owner.ID
	Currency account.Currency
}

type accountRepository interface {
	Create(ctx context.Context, account account.Account) error
}

type createAccountHandler struct {
	accountRepository accountRepository
}

// CreateAccountHandler is a handler for CreateAccount command.
type CreateAccountHandler decorator.CommandHandler[CreateAccount]

// NewCreateAccountHandler creates a new CreateAccountHandler.
func NewCreateAccountHandler(repo accountRepository) CreateAccountHandler {
	if repo == nil {
		panic("accountRepository cannot be nil")
	}

	return decorator.ApplyCommandDecorators(createAccountHandler{
		accountRepository: repo,
	})
}

// Handle handles the CreateAccount command.
func (h createAccountHandler) Handle(ctx context.Context, cmd CreateAccount) error {
	acc, err := account.CreateNewAccount(cmd.ID, cmd.OwnerID, cmd.Currency, time.Now())
	if err != nil {
		return err
	}
	err = h.accountRepository.Create(ctx, acc)
	if err != nil {
		return err
	}
	return nil
}
