package adapters

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresAccountsRepository struct {
	DB *sql.DB
}

func NewPostgresSolutionRepository(connStr string) (*PostgresAccountsRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresAccountsRepository{DB: db}, nil
}
