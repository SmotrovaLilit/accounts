package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"accounts/internal/domain/account"
	"accounts/internal/domain/owner"
)

func TestTransfer(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name              string
		from              account.Account
		to                account.Account
		money             account.Money
		expectBalanceFrom account.Money
		expectBalanceTo   account.Money

		shouldFail  bool
		expectError error
	}{
		{
			name: "should transfer money from one account to another",
			from: account.MustNewAccount(
				account.MustNewID("10eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				owner.MustNewID("40eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				account.MustNewMoney(100, account.CurrencyUSD),
				account.MustNewTimestamps(now, now),
			),
			to: account.MustNewAccount(
				account.MustNewID("20eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				owner.MustNewID("30eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				account.MustNewMoney(100, account.CurrencyUSD),
				account.MustNewTimestamps(now, now),
			),
			money:             account.MustNewMoney(50, account.CurrencyUSD),
			expectBalanceFrom: account.MustNewMoney(50, account.CurrencyUSD),
			expectBalanceTo:   account.MustNewMoney(150, account.CurrencyUSD),
		},
		{
			name: "should fail when transfer reduce failed",
			from: account.MustNewAccount(
				account.MustNewID("10eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				owner.MustNewID("40eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				account.MustNewMoney(100, account.CurrencyUSD),
				account.MustNewTimestamps(now, now),
			),
			to: account.MustNewAccount(
				account.MustNewID("20eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				owner.MustNewID("30eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				account.MustNewMoney(100, account.CurrencyUSD),
				account.MustNewTimestamps(now, now),
			),
			money:             account.MustNewMoney(150, account.CurrencyUSD),
			expectBalanceFrom: account.MustNewMoney(100, account.CurrencyUSD),
			expectBalanceTo:   account.MustNewMoney(100, account.CurrencyUSD),
			shouldFail:        true,
			expectError:       account.ErrNotEnoughBalance,
		},
		{
			name: "should fail when transfer increase failed",
			from: account.MustNewAccount(
				account.MustNewID("10eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				owner.MustNewID("40eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				account.MustNewMoney(100, account.CurrencyUSD),
				account.MustNewTimestamps(now, now),
			),
			to: account.MustNewAccount(
				account.MustNewID("20eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				owner.MustNewID("30eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				account.MustNewMoney(100, account.CurrencyEUR),
				account.MustNewTimestamps(now, now),
			),
			money:             account.MustNewMoney(50, account.CurrencyUSD),
			expectBalanceFrom: account.MustNewMoney(100, account.CurrencyUSD),
			expectBalanceTo:   account.MustNewMoney(100, account.CurrencyEUR),
			shouldFail:        true,
			expectError: account.ErrCurrencyMismatch{
				Expected: account.CurrencyEUR.Code(),
				Got:      account.CurrencyUSD.Code(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Transfer(&tt.from, &tt.to, tt.money, now)
			if tt.shouldFail {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectError)
				// Check if the balance is not changed
				require.Equal(t, tt.expectBalanceFrom, tt.from.Balance())
				require.Equal(t, tt.expectBalanceTo, tt.to.Balance())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectBalanceFrom, tt.from.Balance())
			require.Equal(t, tt.expectBalanceTo, tt.to.Balance())
		})
	}
}
