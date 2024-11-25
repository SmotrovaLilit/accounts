package account

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// ErrInvalidID is returned when an invalid ID is provided.
var ErrInvalidID = errors.New("invalid id")

// ID represents an account identifier.
type ID struct {
	value uuid.UUID
}

// NewID creates a new ID from a string.
func NewID(stringID string) (ID, error) {
	id, err := uuid.Parse(stringID)
	if err != nil {
		return ID{}, fmt.Errorf("%w: %s", ErrInvalidID, err)
	}
	return ID{id}, nil
}

// MustNewID creates a new ID from a string and panics if it fails.
func MustNewID(s string) ID {
	id, err := NewID(s)
	if err != nil {
		panic(err)
	}
	return id
}

// Value returns the UUID value of the ID.
func (i ID) Value() uuid.UUID {
	return i.value
}

// IsZero returns true if the ID is the zero value.
// It can happen when ID was created with ID{}.
// ID should be created with NewID to avoid zero value.
func (i ID) IsZero() bool {
	return i.value == uuid.Nil
}

// String returns the string representation of the ID.
func (i ID) String() string {
	return i.value.String()
}
