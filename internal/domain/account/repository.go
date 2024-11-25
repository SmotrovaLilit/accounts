package account

import (
	"context"
	"fmt"
)

// NotFoundError is an error for a not found account.
type NotFoundError struct {
	ID ID
}

// Error returns the error message.
func (e NotFoundError) Error() string {
	return fmt.Sprintf("account '%s' not found", e.ID.String())
}

// Repository is a repository for accounts.
type Repository interface {
	Create(ctx context.Context, acc Account) error
	Get(ctx context.Context, id ID) (Account, error)
}
