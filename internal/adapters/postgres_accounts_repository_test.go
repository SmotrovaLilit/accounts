package adapters

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"accounts/internal/common/testhelpers"
	"accounts/internal/domain/account"
	"accounts/internal/domain/owner"
)

var container *testhelpers.PostgresContainer

func TestCreateAndAccount(t *testing.T) {
	t.Parallel()
	repo, err := NewPostgresAccountRepository(container.ConnectionString)
	require.NoError(t, err)

	now := time.Now()
	acc := account.MustNewAccount(
		account.MustNewID("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
		owner.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
		account.MustNewMoney(10, account.CurrencyUSD),
		account.MustNewTimestamps(now, now),
	)

	err = repo.Create(context.Background(), acc)
	require.NoError(t, err)

	saved, err := repo.Get(context.Background(), acc.ID())
	require.NoError(t, err)
	require.Equal(t, acc.ID(), saved.ID())
	require.Equal(t, acc.OwnerID(), saved.OwnerID())
	require.Equal(t, acc.Balance(), saved.Balance())
	require.WithinDuration(t, acc.CreatedAt(), saved.CreatedAt(), time.Millisecond)
	require.WithinDuration(t, acc.UpdatedAt(), saved.UpdatedAt(), time.Millisecond)
}

func TestGetAccountNotFound(t *testing.T) {
	t.Parallel()
	// PREPARE
	repo, err := NewPostgresAccountRepository(container.ConnectionString)
	require.NoError(t, err)

	// ACT
	_, err = repo.Get(context.Background(), account.MustNewID("c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12"))

	// ASSERT
	require.Error(t, err)
	require.ErrorAs(t, err, &account.NotFoundError{})
}

func TestUpdateAccount(t *testing.T) {
	// PREPARE
	t.Parallel()
	repo, err := NewPostgresAccountRepository(container.ConnectionString)
	require.NoError(t, err)
	acc := prepareAccount(t, repo, account.MustNewID("bce164bd-0297-473c-8d82-b3012e0821ad"), 10)

	// ACT
	err = repo.UpdateByID(context.Background(), acc.ID(), func(a *account.Account) (bool, error) {
		err := a.TopUp(account.MustNewMoney(103, account.CurrencyUSD), time.Now())
		if err != nil {
			return false, err
		}
		return true, nil
	})

	// ASSERT
	require.NoError(t, err)
	saved, err := repo.Get(context.Background(), acc.ID())
	require.NoError(t, err)
	require.Equal(t, int64(113), saved.BalanceAmount())
	require.WithinDuration(t, time.Now(), saved.UpdatedAt(), time.Second)
}

func TestUpdateAccountNotFound(t *testing.T) {
	// PREPARE
	t.Parallel()
	repo, err := NewPostgresAccountRepository(container.ConnectionString)
	require.NoError(t, err)

	// ACT
	err = repo.UpdateByID(context.Background(), account.MustNewID("c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12"), func(a *account.Account) (bool, error) {
		panic("should not be called")
	})

	// ASSERT
	require.Error(t, err)
	require.ErrorAs(t, err, &account.NotFoundError{})
}

func prepareAccount(t *testing.T, repo *PostgresAccountsRepository, id account.ID, amount int64) account.Account {
	t.Helper()
	acc := account.MustNewAccount(
		id,
		owner.MustNewID("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
		account.MustNewMoney(amount, account.CurrencyUSD),
		account.MustNewTimestamps(time.Now().Add(-10*time.Second), time.Now().Add(-10*time.Second)),
	)
	err := repo.Create(context.Background(), acc)
	require.NoError(t, err)
	return acc
}

func startDB() (*testhelpers.PostgresContainer, bool) {
	cont, err := testhelpers.CreatePostgresContainer(context.Background())
	if err != nil {
		log.Printf("failed to start database: %v", err)
		return nil, false
	}
	container = cont
	return cont, true
}

func TestMain(m *testing.M) {
	var ok bool
	container, ok = startDB()
	if !ok {
		os.Exit(1)
	}
	os.Exit(m.Run())
}
