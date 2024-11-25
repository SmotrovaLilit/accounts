package account

import (
	"errors"
	"time"
)

var (
	// ErrEmptyCreatedAt is returned when created at is empty.
	ErrEmptyCreatedAt = errors.New("empty created at")

	// ErrEmptyUpdatedAt is returned when updated at is empty.
	ErrEmptyUpdatedAt = errors.New("empty updated at")

	// ErrUpdatedAtBeforeCreatedAt is returned when updated at is before created at.
	ErrUpdatedAtBeforeCreatedAt = errors.New("updated at before created at")
)

// Timestamps represents the creation and update timestamps.
type Timestamps struct {
	createdAt time.Time
	updatedAt time.Time
}

// NewTimestamps creates a new timestamps.
func NewTimestamps(createdAt, updatedAt time.Time) (Timestamps, error) {
	if createdAt.IsZero() {
		return Timestamps{}, ErrEmptyCreatedAt
	}
	if updatedAt.IsZero() {
		return Timestamps{}, ErrEmptyUpdatedAt
	}
	if updatedAt.Before(createdAt) {
		return Timestamps{}, ErrUpdatedAtBeforeCreatedAt
	}
	return Timestamps{
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

// MustNewTimestamps creates a new timestamps and panics if it fails.
func MustNewTimestamps(createdAt, updatedAt time.Time) Timestamps {
	t, err := NewTimestamps(createdAt, updatedAt)
	if err != nil {
		panic(err)
	}
	return t
}

// CreatedAt returns the created at timestamp.
func (t Timestamps) CreatedAt() time.Time { return t.createdAt }

// UpdatedAt returns the updated at timestamp.
func (t Timestamps) UpdatedAt() time.Time { return t.updatedAt }

// IsZero returns true if the timestamps are the zero value.
// It can happen when Timestamps was created with Timestamps{}.
// Timestamps should be created with NewTimestamps to avoid zero value.
func (t Timestamps) IsZero() bool {
	return t.createdAt.IsZero() && t.updatedAt.IsZero()
}

// Update updates the updated at timestamp.
func (t Timestamps) Update(now time.Time) Timestamps {
	t.updatedAt = now
	return t
}
