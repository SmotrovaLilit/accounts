package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"accounts/internal/domain/account"
)

// PostgresAccountsRepository is a repository for accounts that uses PostgreSQL.
type PostgresAccountsRepository struct {
	db *sql.DB
}

// NewPostgresAccountRepository creates a new PostgresAccountsRepository.
func NewPostgresAccountRepository(connStr string) (*PostgresAccountsRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresAccountsRepository{db: db}, nil
}

// Create inserts a new account into the database.
func (p PostgresAccountsRepository) Create(ctx context.Context, acc account.Account) error {
	query := `
		INSERT INTO accounts (id, owner_id, balance, currency, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := p.db.ExecContext(ctx, query, acc.ID().Value(), acc.OwnerID().Value(), acc.BalanceAmount(), acc.BalanceCurrencyCode(), acc.CreatedAt(), acc.UpdatedAt())
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}
	return nil
}

// Close closes the database connection.
func (p PostgresAccountsRepository) Close() error {
	return p.db.Close()
}

type accountEntity struct {
	ID        uuid.UUID
	OwnerID   uuid.UUID
	Balance   int64
	Currency  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Get retrieves an account from the database.
// It returns a NotFoundError if the account is not found.
func (p PostgresAccountsRepository) Get(ctx context.Context, id account.ID) (account.Account, error) {
	query := `
		SELECT id, owner_id, balance, currency, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`
	row := p.db.QueryRowContext(ctx, query, id.Value())
	entity := accountEntity{}
	err := row.Scan(&entity.ID, &entity.OwnerID, &entity.Balance, &entity.Currency, &entity.CreatedAt, &entity.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return account.Account{}, account.NotFoundError{ID: id}
	}
	if err != nil {
		return account.Account{}, fmt.Errorf("failed to get account: %w", err)
	}
	return convertToAccount(entity), nil
}

// UpdateByID updates an account in the database.
// It returns a NotFoundError if the account is not found.
// It reads the account from the database (FOR UPDATE), calls the update function, and then updates the account in the database.
// See examples of such approach here https://threedots.tech/post/database-transactions-in-go/
func (p PostgresAccountsRepository) UpdateByID(ctx context.Context, id account.ID, updateFn func(acc *account.Account) (bool, error)) error {
	return runInTx(p.db, func(tx *sql.Tx) error {
		query := `
			SELECT id, owner_id, balance, currency, created_at, updated_at
			FROM accounts
			WHERE id = $1 FOR UPDATE
		`
		row := tx.QueryRowContext(ctx, query, id.Value())

		var entity accountEntity
		err := row.Scan(
			&entity.ID,
			&entity.OwnerID,
			&entity.Balance,
			&entity.Currency,
			&entity.CreatedAt,
			&entity.UpdatedAt,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return account.NotFoundError{ID: id}
			}
			return fmt.Errorf("failed to scan account: %w", err)
		}

		acc := convertToAccount(entity)

		// Apply the update function
		updated, err := updateFn(&acc)
		if err != nil {
			return fmt.Errorf("update function failed: %w", err)
		}

		// If no update is necessary, return early
		if !updated {
			return nil
		}

		// Update the account in the database
		updateQuery := `
			UPDATE accounts
			SET balance = $1, updated_at = $2
			WHERE id = $3
		`
		_, err = tx.ExecContext(ctx, updateQuery, acc.BalanceAmount(), acc.UpdatedAt(), acc.ID().Value())
		if err != nil {
			return fmt.Errorf("failed to update account: %w", err)
		}

		return nil
	})
}

func (p PostgresAccountsRepository) UpdateTwoAccounts(ctx context.Context, id1, id2 account.ID, updateFn func(acc1, acc2 *account.Account) (bool, error)) error {
	return runInTx(p.db, func(tx *sql.Tx) error {
		query := `
			SELECT id, owner_id, balance, currency, created_at, updated_at
			FROM accounts
			WHERE id = $1 OR id = $2
			FOR UPDATE
		`
		rows, err := tx.QueryContext(ctx, query, id1.Value(), id2.Value())
		if err != nil {
			return fmt.Errorf("failed to query accounts: %w", err)
		}
		if err := rows.Err(); err != nil {
			return fmt.Errorf("failed to query accounts: %w", err)
		}
		defer rows.Close()

		var entities []accountEntity
		for rows.Next() {
			var entity accountEntity
			err := rows.Scan(
				&entity.ID,
				&entity.OwnerID,
				&entity.Balance,
				&entity.Currency,
				&entity.CreatedAt,
				&entity.UpdatedAt,
			)
			if err != nil {
				return fmt.Errorf("failed to scan account: %w", err)
			}
			entities = append(entities, entity)
		}

		if len(entities) != 2 {
			if len(entities) == 0 {
				return errors.Join(account.NotFoundError{ID: id1}, account.NotFoundError{ID: id2})
			}
			if entities[0].ID == id1.Value() {
				return account.NotFoundError{ID: id2}
			}
			return account.NotFoundError{ID: id1}
		}

		acc1 := convertToAccount(entities[0])
		acc2 := convertToAccount(entities[1])

		// Apply the update function
		updated, err := updateFn(&acc1, &acc2)
		if err != nil {
			return fmt.Errorf("update function failed: %w", err)
		}

		// If no update is necessary, return early
		if !updated {
			return nil
		}

		// Update the accounts in the database
		updateQuery := `
			UPDATE accounts
			SET balance = $1, updated_at = $2
			WHERE id = $3
		`
		_, err = tx.ExecContext(ctx, updateQuery, acc1.BalanceAmount(), acc1.UpdatedAt(), acc1.ID().Value())
		if err != nil {
			return fmt.Errorf("failed to update account 1: %w", err)
		}
		_, err = tx.ExecContext(ctx, updateQuery, acc2.BalanceAmount(), acc2.UpdatedAt(), acc2.ID().Value())
		if err != nil {
			return fmt.Errorf("failed to update account 2: %w", err)
		}

		return nil
	})
}

func convertToAccount(entity accountEntity) account.Account {
	return account.UnmarshalAccountFromStorage(
		entity.ID.String(),
		entity.OwnerID.String(),
		entity.Balance,
		entity.Currency,
		entity.CreatedAt,
		entity.UpdatedAt,
	)
}
