package account

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func TestNewTimestamps(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name       string
		createdAt  time.Time
		updatedAt  time.Time
		want       Timestamps
		shouldFail bool
		wantErr    error
	}{
		{
			name:      "should return success if timestamps are valid",
			createdAt: now,
			updatedAt: now.Add(time.Second),
			want: Timestamps{
				createdAt: now,
				updatedAt: now.Add(time.Second),
			},
		},
		{
			name:       "should return error if created at is empty",
			createdAt:  time.Time{},
			updatedAt:  now,
			shouldFail: true,
			wantErr:    ErrEmptyCreatedAt,
		},
		{
			name:       "should return error if updated at is empty",
			createdAt:  now,
			updatedAt:  time.Time{},
			shouldFail: true,
			wantErr:    ErrEmptyUpdatedAt,
		},
		{
			name:       "should return error if updated at is before created at",
			createdAt:  now,
			updatedAt:  now.Add(-time.Second),
			shouldFail: true,
			wantErr:    ErrUpdatedAtBeforeCreatedAt,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTimestamps(tt.createdAt, tt.updatedAt)
			if tt.shouldFail {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.True(t, reflect.DeepEqual(tt.want, got))
		})
	}
}
