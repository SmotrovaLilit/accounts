package account

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"accounts/internal/domain/owner"
)

func TestCreateNewAccount(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		id       ID
		ownerID  owner.ID
		currency Currency
		want     Account
		now      time.Time
		wantErr  error
		showFail bool
	}{
		{
			name:     "should return success if account is valid",
			id:       MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
			ownerID:  owner.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
			currency: CurrencyUSD,
			now:      now,
			want: Account{
				id:      MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				ownerID: owner.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
				balance: Money{
					amount:   0,
					currency: CurrencyUSD,
				},
				timestamps: MustNewTimestamps(now, now),
			},
		},
		{
			name:     "should return error if id is invalid",
			id:       ID{},
			ownerID:  owner.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
			currency: CurrencyUSD,
			now:      now,
			showFail: true,
			wantErr:  ErrZeroObjectInputParameters,
		},
		{
			name:     "should return error if ownerID is invalid",
			id:       MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
			ownerID:  owner.ID{},
			currency: CurrencyUSD,
			now:      now,
			showFail: true,
			wantErr:  ErrZeroObjectInputParameters,
		},
		{
			name:     "should return error if currency is invalid",
			id:       MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
			ownerID:  owner.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
			currency: Currency{},
			now:      now,
			showFail: true,
			wantErr:  ErrInvalidCurrency,
		},
		{
			name:     "should return error time now is invalid",
			id:       MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
			ownerID:  owner.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
			currency: CurrencyUSD,
			now:      time.Time{},
			showFail: true,
			wantErr:  ErrEmptyCreatedAt,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateNewAccount(tt.id, tt.ownerID, tt.currency, tt.now)
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
