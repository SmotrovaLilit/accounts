package commands

import (
	"context"
	"errors"
	"testing"
	"time"

	"accounts/internal/domain/account"
	"accounts/internal/domain/owner"

	"github.com/stretchr/testify/require"
)

func Test_topUpAccountHandler_Handle(t *testing.T) {
	someErr := errors.New("some error")
	testCases := []struct {
		name        string
		cmd         TopUpAccount
		shouldFail  bool
		expectError error
		repo        *topUpAccountRepositoryMock
	}{
		{
			name: "Should top up account",
			cmd: TopUpAccount{
				ID:     account.MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				Amount: 100,
			},
			shouldFail: false,
			repo: newTopUpAccountRepositoryMock(nil, account.MustNewAccount(
				account.MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				owner.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				account.MustNewMoney(100, account.CurrencyUSD),
				account.MustNewTimestamps(time.Now(), time.Now()),
			)),
		},
		{
			name: "Should fail when error from domain",
			cmd: TopUpAccount{
				ID:     account.MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				Amount: -100,
			},
			repo: newTopUpAccountRepositoryMock(nil, account.MustNewAccount(
				account.MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				owner.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				account.MustNewMoney(100, account.CurrencyUSD),
				account.MustNewTimestamps(time.Now(), time.Now()),
			)),
			shouldFail:  true,
			expectError: account.ErrCannotAddRemoveNegativeAmount,
		},
		{
			name: "Should fail when error from repository",
			cmd: TopUpAccount{
				ID:     account.MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				Amount: 100,
			},
			repo: newTopUpAccountRepositoryMock(someErr, account.MustNewAccount(
				account.MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				owner.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				account.MustNewMoney(100, account.CurrencyUSD),
				account.MustNewTimestamps(time.Now(), time.Now()),
			)),
			shouldFail:  true,
			expectError: someErr,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := tt.repo
			handler := NewTopUpAccountHandler(repo)
			err := handler.Handle(context.Background(), tt.cmd)
			if tt.shouldFail {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectError)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.cmd.Amount+tt.repo.existingAccount.BalanceAmount(), tt.repo.savedAccount.Balance().Amount())
		})
	}
}

type topUpAccountRepositoryMock struct {
	returnErr       error
	savedAccount    account.Account
	existingAccount account.Account
}

func newTopUpAccountRepositoryMock(returnErr error, existing account.Account) *topUpAccountRepositoryMock {
	return &topUpAccountRepositoryMock{returnErr: returnErr, existingAccount: existing}
}

func (m *topUpAccountRepositoryMock) UpdateByID(ctx context.Context, id account.ID, updateFn func(acc *account.Account) (bool, error)) error {
	if m.returnErr != nil {
		return m.returnErr
	}
	acc := m.existingAccount.Clone()
	save, err := updateFn(&acc)
	if save {
		m.savedAccount = acc
	}
	return err
}
