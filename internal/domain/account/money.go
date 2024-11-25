package account

import "fmt"

var (
	OverdraftLimit                   = int64(-10000) // 100.00
	ErrInvalidAmount                 = fmt.Errorf("invalid amount, must be greater than %d", OverdraftLimit)
	ErrCannotAddRemoveNegativeAmount = fmt.Errorf("cannot add/remove negative amount")
	ErrAmountOverflow                = fmt.Errorf("amount overflow")
	ErrNotEnoughBalance              = fmt.Errorf("not enough balance")
)

// Money represents a money value with a currency.
type Money struct {
	amount   int64
	currency Currency
}

// NewMoney creates a new money value.
func NewMoney(amount int64, currency Currency) (Money, error) {
	if amount < OverdraftLimit {
		return Money{}, ErrInvalidAmount
	}

	if currency.IsZero() {
		return Money{}, ErrInvalidCurrency
	}

	return Money{
		amount:   amount,
		currency: currency,
	}, nil
}

// MustNewMoney creates a new money value and panics if it fails.
func MustNewMoney(amount int64, currency Currency) Money {
	m, err := NewMoney(amount, currency)
	if err != nil {
		panic(err)
	}
	return m
}

// IsZeroObject returns true if the money was created in wrong way using Money{}.
// Money should be created with NewMoney to avoid zero object value.
func (m Money) IsZeroObject() bool {
	return m.currency.IsZero()
}

// Amount returns the money amount.
func (m Money) Amount() int64 {
	return m.amount
}

// CurrencyCode returns the currency code.
func (m Money) CurrencyCode() string {
	return m.currency.Code()
}

// Add adds money to the current money value.
func (m Money) Add(money Money) (Money, error) {
	if money.IsZeroObject() {
		return Money{}, ErrZeroObjectInputParameters
	}
	if m.currency != money.currency {
		return Money{}, ErrCurrencyMismatch{Expected: m.currency.Code(), Got: money.currency.Code()}
	}
	if money.amount < 0 {
		return Money{}, ErrCannotAddRemoveNegativeAmount
	}
	amount, overflow := safeAddInt64(m.amount, money.amount)
	if overflow {
		return Money{}, ErrAmountOverflow
	}
	return Money{
		amount:   amount,
		currency: m.currency,
	}, nil
}

func (m Money) Reduce(money Money) (Money, error) {
	if money.IsZeroObject() {
		return Money{}, ErrZeroObjectInputParameters
	}
	if m.currency != money.currency {
		return Money{}, ErrCurrencyMismatch{Expected: m.currency.Code(), Got: money.currency.Code()}
	}
	if money.amount < 0 {
		return Money{}, ErrCannotAddRemoveNegativeAmount
	}
	amount := m.amount - money.amount
	if amount < 0 {
		return Money{}, ErrNotEnoughBalance
	}
	return Money{
		amount:   amount,
		currency: m.currency,
	}, nil
}

// ErrCurrencyMismatch is returned when the currency is different.
type ErrCurrencyMismatch struct {
	Expected string
	Got      string
}

// Error returns the error message.
func (e ErrCurrencyMismatch) Error() string {
	return fmt.Sprintf("currency mismatch: expected %s, got %s", e.Expected, e.Got)
}

// safeAddInt64 safely adds two int64 numbers and returns whether an overflow occurred.
// It returns the sum and a boolean indicating whether an overflow occurred.
func safeAddInt64(a, b int64) (int64, bool) {
	sum := a + b
	if (a > 0 && b > 0 && sum < 0) || (a < 0 && b < 0 && sum > 0) {
		return 0, true
	}
	return sum, false
}
