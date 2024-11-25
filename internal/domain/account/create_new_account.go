package account

import (
	"time"

	"accounts/internal/domain/owner"
)

// CreateNewAccount creates a new account with a zero balance.
func CreateNewAccount(id ID, ownerID owner.ID, currency Currency, now time.Time) (Account, error) {
	balance, err := NewMoney(int64(0), currency)
	if err != nil {
		return Account{}, err
	}
	t, err := NewTimestamps(now, now)
	if err != nil {
		return Account{}, err
	}
	return NewAccount(id, ownerID, balance, t)
}
