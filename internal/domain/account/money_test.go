package account

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestNewMoney(t *testing.T) {
	tests := []struct {
		name     string
		amount   int64
		currency Currency
		want     Money
		showFail bool
		wantErr  error
	}{
		{
			name:     "should return success if money is valid",
			amount:   -10000,
			currency: CurrencyUSD,
			want: Money{
				amount:   -10000,
				currency: CurrencyUSD,
			},
		},
		{
			name:     "should return error if amount is invalid",
			amount:   -10001,
			currency: CurrencyUSD,
			showFail: true,
			wantErr:  ErrInvalidAmount,
		},
		{
			name:     "should return error if currency is invalid",
			amount:   -10000,
			currency: Currency{},
			showFail: true,
			wantErr:  ErrInvalidCurrency,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMoney(tt.amount, tt.currency)
			if tt.showFail {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.True(t, reflect.DeepEqual(tt.want, got))
		})
	}
}

func TestMoney_Add(t *testing.T) {
	tests := []struct {
		name       string
		existing   Money
		moneyToAdd Money
		want       Money
		shouldFail bool
		wantErr    error
	}{
		{
			name: "should return success if money is added",
			existing: Money{
				amount:   10,
				currency: CurrencyUSD,
			},
			moneyToAdd: Money{
				amount:   10,
				currency: CurrencyUSD,
			},
			want: Money{
				amount:   20,
				currency: CurrencyUSD,
			},
		},
		{
			name: "should return error if currency is not the same",
			existing: Money{
				amount:   10,
				currency: CurrencyUSD,
			},
			moneyToAdd: Money{
				amount:   10,
				currency: CurrencyEUR,
			},
			shouldFail: true,
			wantErr:    ErrCurrencyMismatch{Expected: CurrencyUSD.Code(), Got: CurrencyEUR.Code()},
		},
		{
			name: "should return error if amount is negative",
			existing: Money{
				amount:   10,
				currency: CurrencyUSD,
			},
			moneyToAdd: Money{
				amount:   -10,
				currency: CurrencyUSD,
			},
			shouldFail: true,
			wantErr:    ErrCannotAddRemoveNegativeAmount,
		},
		{
			name: "should return error if amount is overflow",
			existing: Money{
				amount:   9223372036854775807,
				currency: CurrencyUSD,
			},
			moneyToAdd: Money{
				amount:   1,
				currency: CurrencyUSD,
			},
			shouldFail: true,
			wantErr:    ErrAmountOverflow,
		},
		{
			name: "should return error if money is zero object",
			existing: Money{
				amount:   10,
				currency: CurrencyUSD,
			},
			moneyToAdd: Money{},
			shouldFail: true,
			wantErr:    ErrZeroObjectInputParameters,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.existing.Add(tt.moneyToAdd)
			if tt.shouldFail {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestMoney_Reduce(t *testing.T) {
	tests := []struct {
		name          string
		existing      Money
		moneyToReduce Money
		want          Money
		shouldFail    bool
		wantErr       error
	}{
		{
			name: "should return success if money is reduced",
			existing: Money{
				amount:   10,
				currency: CurrencyUSD,
			},
			moneyToReduce: Money{
				amount:   10,
				currency: CurrencyUSD,
			},
			want: Money{
				amount:   0,
				currency: CurrencyUSD,
			},
		},
		{
			name: "should return error if currency is not the same",
			existing: Money{
				amount:   10,
				currency: CurrencyUSD,
			},
			moneyToReduce: Money{
				amount:   10,
				currency: CurrencyEUR,
			},
			shouldFail: true,
			wantErr:    ErrCurrencyMismatch{Expected: CurrencyUSD.Code(), Got: CurrencyEUR.Code()},
		},
		{
			name: "should return error if amount is negative",
			existing: Money{
				amount:   10,
				currency: CurrencyUSD,
			},
			moneyToReduce: Money{
				amount:   -10,
				currency: CurrencyUSD,
			},
			shouldFail: true,
			wantErr:    ErrCannotAddRemoveNegativeAmount,
		},
		{
			name: "should return error if not enough balance",
			existing: Money{
				amount:   10,
				currency: CurrencyUSD,
			},
			moneyToReduce: Money{
				amount:   11,
				currency: CurrencyUSD,
			},
			shouldFail: true,
			wantErr:    ErrNotEnoughBalance,
		},
		{
			name: "should return error if money is zero object",
			existing: Money{
				amount:   10,
				currency: CurrencyUSD,
			},
			moneyToReduce: Money{},
			shouldFail:    true,
			wantErr:       ErrZeroObjectInputParameters,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.existing.Reduce(tt.moneyToReduce)
			if tt.shouldFail {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
