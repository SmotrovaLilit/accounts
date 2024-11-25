package account

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"accounts/internal/domain/owner"
)

func TestNewAccount(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		id          ID
		ownerID     owner.ID
		balance     Money
		timestamps  Timestamps
		want        Account
		showFail    bool
		wantErrText string
	}{
		{
			name:       "should return success if account is valid",
			id:         MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
			ownerID:    owner.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
			balance:    MustNewMoney(1, CurrencyUSD),
			timestamps: MustNewTimestamps(now, now),
			want: Account{
				id:         MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				ownerID:    owner.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				balance:    MustNewMoney(1, CurrencyUSD),
				timestamps: MustNewTimestamps(now, now),
			},
		},
		{
			name:        "should return error if id is invalid",
			id:          ID{},
			ownerID:     owner.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
			balance:     MustNewMoney(1, CurrencyUSD),
			timestamps:  MustNewTimestamps(now, now),
			showFail:    true,
			wantErrText: "invalid input parameters, use constructors to create objects",
		},
		{
			name:        "should return error if ownerID is invalid",
			id:          MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
			ownerID:     owner.ID{},
			balance:     MustNewMoney(1, CurrencyUSD),
			timestamps:  MustNewTimestamps(now, now),
			showFail:    true,
			wantErrText: "invalid input parameters, use constructors to create objects",
		},
		{
			name:        "should return error if balance is invalid",
			id:          MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
			ownerID:     owner.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
			balance:     Money{},
			timestamps:  MustNewTimestamps(now, now),
			showFail:    true,
			wantErrText: "invalid input parameters, use constructors to create objects",
		},
		{
			name:        "should return error if timestamps is invalid",
			id:          MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
			ownerID:     owner.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
			balance:     MustNewMoney(1, CurrencyUSD),
			timestamps:  Timestamps{},
			showFail:    true,
			wantErrText: "invalid input parameters, use constructors to create objects",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccount(tt.id, tt.ownerID, tt.balance, tt.timestamps)
			if tt.showFail {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.wantErrText)
				return
			}
			require.NoError(t, err)
			require.True(t, reflect.DeepEqual(tt.want, got))
		})
	}
}

func TestAccount_TopUp(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name       string
		account    Account
		money      Money
		want       Money
		shouldFail bool
		wantErr    error
	}{
		{
			name:    "should return success if money is valid",
			account: MustNewAccount(MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"), owner.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"), MustNewMoney(1, CurrencyUSD), MustNewTimestamps(time.Now(), time.Now())),
			money:   MustNewMoney(1, CurrencyUSD),
			want:    MustNewMoney(2, CurrencyUSD),
		},
		{
			name:       "should return error if money is invalid",
			account:    MustNewAccount(MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"), owner.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"), MustNewMoney(1, CurrencyUSD), MustNewTimestamps(time.Now(), time.Now())),
			money:      Money{},
			shouldFail: true,
			wantErr:    ErrZeroObjectInputParameters,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.account.TopUp(tt.money, now)
			if tt.shouldFail {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, tt.account.Balance())
			require.WithinDuration(t, now, tt.account.UpdatedAt(), time.Millisecond)
		})
	}
}
