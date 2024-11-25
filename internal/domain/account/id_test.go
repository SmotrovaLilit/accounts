package account

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewID(t *testing.T) {
	t.Run("should return success if id is valid", func(t *testing.T) {
		v, err := NewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
		require.NoError(t, err)
		require.Equal(t, "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", v.String())
	})
	t.Run("should return error if id is invalid", func(t *testing.T) {
		_, err := NewID("invalid")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrInvalidID)
	})
}
