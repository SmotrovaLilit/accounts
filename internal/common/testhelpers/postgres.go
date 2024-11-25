package testhelpers

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// No-lint.
	_ "github.com/golang-migrate/migrate/v4/source/file"
	// No-lint.
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/testcontainers/testcontainers-go"
	testpostgress "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	// No-lint.
	_ "github.com/lib/pq"
)

type PostgresContainer struct {
	*testpostgress.PostgresContainer
	ConnectionString string
}

func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	pgContainer, err := testpostgress.Run(ctx,
		"postgres:15.3-alpine",
		testpostgress.WithDatabase("test-db"),
		testpostgress.WithUsername("postgres"),
		testpostgress.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	err = runMigrations(connStr)
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connStr,
	}, nil
}

func runMigrations(connStr string) error {
	absPath, err := filepath.Abs("../../migrations")
	if err != nil {
		return err
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errors.New("migrations path is not a directory")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+absPath,
		"postgres", driver)
	if err != nil {
		return err
	}
	return m.Up()
}
