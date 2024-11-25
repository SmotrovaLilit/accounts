package domain

import (
	"fmt"
	"time"

	"accounts/internal/domain/account"
)

// Transfer transfers money from one account to another.
func Transfer(from, to *account.Account, money account.Money, now time.Time) error {
	if from == nil || to == nil {
		panic("from and to accounts cannot be nil")
	}

	// Clone the account to be able to rollback the changes in case of failure in next steps.
	fromCopy := from.Clone()
	err := fromCopy.ReduceBalance(money, now)
	if err != nil {
		return fmt.Errorf("failed reduce balance from account: %w", err)
	}
	err = to.TopUp(money, now)
	if err != nil {
		return fmt.Errorf("failed top up account: %w", err)
	}
	*from = fromCopy

	return nil
}
