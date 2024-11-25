package account

import (
	"errors"
	"time"

	"accounts/internal/domain/owner"
)

// ErrZeroObjectInputParameters is an error for zero object input parameters.
var ErrZeroObjectInputParameters = errors.New("invalid input parameters, use constructors to create objects")

// Account is a bank account.
type Account struct {
	id         ID
	ownerID    owner.ID
	balance    Money
	timestamps Timestamps
}

// NewAccount is a constructor for Account.
func NewAccount(id ID, ownerID owner.ID, balance Money, timestamps Timestamps) (Account, error) {
	if id.IsZero() || ownerID.IsZero() || balance.IsZeroObject() || timestamps.IsZero() {
		return Account{}, ErrZeroObjectInputParameters
	}
	return Account{
		id:         id,
		ownerID:    ownerID,
		balance:    balance,
		timestamps: timestamps,
	}, nil
}

// MustNewAccount is a constructor for Account that panics on error.
func MustNewAccount(id ID, ownerID owner.ID, balance Money, timestamps Timestamps) Account {
	acc, err := NewAccount(id, ownerID, balance, timestamps)
	if err != nil {
		panic(err)
	}
	return acc
}

// UnmarshalAccountFromStorage creates an account from storage.
// !!! This function is used just for converting data from storage to domain model.
// !!! It should not be used in the application or domain layers.
// We trust that the data in the storage is correct, and we don't need to validate it.
// It is risky to use this function in the application or domain layers because it doesn't validate the input data.
func UnmarshalAccountFromStorage(idS, ownerIDs string, balanceAmount int64, balanceCurrencyCode string, createdAt, updatedAt time.Time) Account {
	return MustNewAccount(MustNewID(idS), owner.MustNewID(ownerIDs), Money{balanceAmount, Currency{code: balanceCurrencyCode}}, Timestamps{createdAt, updatedAt})
}

// ID returns the account ID.
func (a *Account) ID() ID {
	return a.id
}

// OwnerID returns the account owner ID.
func (a *Account) OwnerID() owner.ID {
	return a.ownerID
}

// Balance returns the account balance.
func (a *Account) Balance() Money {
	return a.balance
}

// CreatedAt returns the created at timestamp.
func (a *Account) CreatedAt() time.Time {
	return a.timestamps.CreatedAt()
}

// UpdatedAt returns the updated at timestamp.
func (a *Account) UpdatedAt() time.Time {
	return a.timestamps.UpdatedAt()
}

// BalanceAmount returns the account balance amount.
func (a *Account) BalanceAmount() int64 {
	return a.balance.Amount()
}

// BalanceCurrencyCode returns the account balance currency code.
func (a *Account) BalanceCurrencyCode() string {
	return a.balance.CurrencyCode()
}

// BalanceCurrency returns the account balance currency.
func (a *Account) BalanceCurrency() Currency {
	return a.balance.currency
}

// TopUp adds money to the account balance.
func (a *Account) TopUp(money Money, now time.Time) error {
	newBalance, err := a.balance.Add(money)
	if err != nil {
		return err
	}
	a.balance = newBalance
	a.timestamps = a.timestamps.Update(now)

	return nil
}

// Clone creates a copy of the account.
func (a *Account) Clone() Account {
	return MustNewAccount(a.id, a.ownerID, a.balance, a.timestamps)
}

func (a *Account) ReduceBalance(money Money, now time.Time) error {
	newBalance, err := a.balance.Reduce(money)
	if err != nil {
		return err
	}
	a.balance = newBalance
	a.timestamps = a.timestamps.Update(now)
	return nil
}
