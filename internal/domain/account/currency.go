package account

import (
	"fmt"
)

var (
	// CurrencyUSD represents the USD currency.
	CurrencyUSD = Currency{"USD"}

	// CurrencyEUR represents the EUR currency.
	CurrencyEUR = Currency{"EUR"}

	// Currencies represents a list of available currencies.
	Currencies = []Currency{CurrencyUSD, CurrencyEUR}

	// ErrInvalidCurrency represents an error when the currency is invalid.
	ErrInvalidCurrency = fmt.Errorf("invalid currency, supported currencies are %v", Currencies)
)

// Currency represents a currency.
type Currency struct {
	code string
}

// NewCurrency creates a new currency.
func NewCurrency(code string) (Currency, error) {
	for _, c := range Currencies {
		if c.code == code {
			return c, nil
		}
	}
	return Currency{}, ErrInvalidCurrency
}

// MustNewCurrency creates a new currency and panics if it fails.
func MustNewCurrency(code string) Currency {
	c, err := NewCurrency(code)
	if err != nil {
		panic(err)
	}
	return c
}

// Code returns the currency code.
func (c Currency) Code() string {
	return c.code
}

// IsZero returns true if the currency is the zero value
// It can happen when Currency was created with Currency{}.
// Currency should be created with NewCurrency to avoid zero value.
func (c Currency) IsZero() bool {
	return c.code == ""
}
