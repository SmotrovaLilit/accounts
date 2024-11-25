package account

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCurrency(t *testing.T) {

	testCases := []struct {
		name       string
		code       string
		want       Currency
		shouldFail bool
		wantErr    error
	}{
		{
			name: "should return USD currency",
			code: "USD",
			want: CurrencyUSD,
		},
		{
			name: "should return EUR currency",
			code: "EUR",
			want: CurrencyEUR,
		},
		{
			name:       "should return error if currency is invalid",
			code:       "INVALID",
			shouldFail: true,
			wantErr:    ErrInvalidCurrency,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCurrency(tt.code)
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
