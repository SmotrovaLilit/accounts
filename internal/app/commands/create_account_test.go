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

func Test_createAccountHandler_Handle(t *testing.T) {
	someErr := errors.New("some error")
	testCases := []struct {
		name        string
		cmd         CreateAccount
		shouldFail  bool
		expectError error
		repo        *createAccountRepositoryMock
	}{
		{
			name: "Should create account",
			cmd: CreateAccount{
				ID:       account.MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				OwnerID:  owner.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				Currency: account.CurrencyUSD,
			},
			shouldFail: false,
			repo:       newCreateAccountRepositoryMock(nil),
		},
		{
			name: "Should fail when error from domain",
			cmd: CreateAccount{
				ID:       account.ID{},
				OwnerID:  owner.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				Currency: account.CurrencyUSD,
			},
			shouldFail:  true,
			expectError: account.ErrZeroObjectInputParameters,
			repo:        newCreateAccountRepositoryMock(nil),
		},
		{
			name: "Should fail when error from repository",
			cmd: CreateAccount{
				ID:       account.MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				OwnerID:  owner.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				Currency: account.CurrencyUSD,
			},
			shouldFail:  true,
			expectError: someErr,
			repo:        newCreateAccountRepositoryMock(someErr),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := tt.repo
			h := NewCreateAccountHandler(repo)
			err := h.Handle(context.Background(), tt.cmd)
			if tt.shouldFail {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectError)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.cmd.ID, repo.createdAccount.ID())
			require.Equal(t, tt.cmd.OwnerID, repo.createdAccount.OwnerID())
			require.Equal(t, account.MustNewMoney(int64(0), tt.cmd.Currency), repo.createdAccount.Balance())
			require.WithinDuration(t, time.Now(), repo.createdAccount.CreatedAt(), time.Millisecond*100)
			require.WithinDuration(t, time.Now(), repo.createdAccount.UpdatedAt(), time.Millisecond*100)
		})
	}
}

type createAccountRepositoryMock struct {
	returnErr      error
	createdAccount account.Account
}

func newCreateAccountRepositoryMock(err error) *createAccountRepositoryMock {
	return &createAccountRepositoryMock{returnErr: err}
}

func (m *createAccountRepositoryMock) Create(ctx context.Context, account account.Account) error {
	m.createdAccount = account
	return m.returnErr
}
