package commands

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"accounts/internal/domain/account"
	"accounts/internal/domain/owner"
)

func Test_transferFundsHandler_Handle(t *testing.T) {
	someErr := errors.New("some error")
	testCases := []struct {
		name        string
		cmd         TransferFunds
		shouldFail  bool
		expectError error
		repo        *transferFundsRepositoryMock
	}{
		{
			name: "Should transfer funds",
			cmd: TransferFunds{
				From:   account.MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				To:     account.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				Amount: 47,
			},
			repo: &transferFundsRepositoryMock{
				acc1: account.MustNewAccount(
					account.MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
					owner.MustNewID("b01ebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
					account.MustNewMoney(100, account.CurrencyUSD),
					account.MustNewTimestamps(time.Now(), time.Now()),
				),
				acc2: account.MustNewAccount(
					account.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
					owner.MustNewID("b02ebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
					account.MustNewMoney(100, account.CurrencyUSD),
					account.MustNewTimestamps(time.Now(), time.Now()),
				),
			},
		},
		{
			name: "Should fail when error from repository",
			cmd: TransferFunds{
				From:   account.MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				To:     account.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				Amount: 100,
			},
			repo: &transferFundsRepositoryMock{
				acc1: account.MustNewAccount(
					account.MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
					owner.MustNewID("b01ebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
					account.MustNewMoney(100, account.CurrencyUSD),
					account.MustNewTimestamps(time.Now(), time.Now()),
				),
				acc2: account.MustNewAccount(
					account.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
					owner.MustNewID("b02ebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
					account.MustNewMoney(100, account.CurrencyUSD),
					account.MustNewTimestamps(time.Now(), time.Now()),
				),
				returnError: someErr,
			},
			shouldFail:  true,
			expectError: someErr,
		},
		{
			name: "Should fail when error from domain",
			cmd: TransferFunds{
				From:   account.MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				To:     account.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				Amount: -100,
			},
			repo: &transferFundsRepositoryMock{
				acc1: account.MustNewAccount(
					account.MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
					owner.MustNewID("b01ebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
					account.MustNewMoney(100, account.CurrencyUSD),
					account.MustNewTimestamps(time.Now(), time.Now()),
				),
				acc2: account.MustNewAccount(
					account.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
					owner.MustNewID("b02ebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
					account.MustNewMoney(100, account.CurrencyUSD),
					account.MustNewTimestamps(time.Now(), time.Now()),
				),
			},
			shouldFail:  true,
			expectError: account.ErrCannotAddRemoveNegativeAmount,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := tt.repo
			handler := NewTransferFundsHandler(repo)
			err := handler.Handle(context.Background(), tt.cmd)
			if tt.shouldFail {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectError)
				return
			}
			require.NoError(t, err)
			require.Equal(t, account.MustNewMoney(repo.acc1.BalanceAmount()-tt.cmd.Amount, repo.acc1.BalanceCurrency()), repo.savedAcc1.Balance())
			require.Equal(t, account.MustNewMoney(repo.acc2.BalanceAmount()+tt.cmd.Amount, repo.acc2.BalanceCurrency()), repo.savedAcc2.Balance())
		})
	}
}

type transferFundsRepositoryMock struct {
	returnError error
	acc1        account.Account
	acc2        account.Account
	savedAcc1   account.Account
	savedAcc2   account.Account
}

func (t *transferFundsRepositoryMock) UpdateTwoAccounts(ctx context.Context, id1, id2 account.ID, updateFn func(acc1 *account.Account, acc2 *account.Account) (bool, error)) error {
	acc1 := t.acc1.Clone()
	acc2 := t.acc2.Clone()

	save, err := updateFn(&acc1, &acc2)
	if err != nil {
		return err
	}
	if save {
		t.savedAcc1 = acc1
		t.savedAcc2 = acc2
	}

	return t.returnError
}
