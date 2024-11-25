package service

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/oklog/run"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"accounts/internal/common/client/grpcclient"
	"accounts/internal/common/server/grpcserver"
	"accounts/internal/common/testhelpers"
	"accounts/internal/config"
	pb "accounts/internal/gen-proto/accountservice/v1beta"
	"accounts/internal/ports"
)

const accountServiceAddr = "localhost:8091"

func TestCreateAndGetAccount(t *testing.T) {
	client := newAccountServiceClient(t)
	testCases := []struct {
		name              string
		shouldFail        bool
		expectErrorString string
		request           *pb.CreateAccountRequest
	}{
		{
			name:       "should create account",
			shouldFail: false,
			request: &pb.CreateAccountRequest{
				OwnerId:  "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				Id:       "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				Currency: pb.Currency_CURRENCY_EUR,
			},
		},
		{
			name:       "should fail when owner id is invalid",
			shouldFail: true,
			request: &pb.CreateAccountRequest{
				OwnerId:  "b0eebc99-9c0b-4ef8",
				Id:       "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				Currency: pb.Currency_CURRENCY_EUR,
			},
			expectErrorString: "Invalid argument provided",
		},
		{
			name:       "should fail when account id is invalid",
			shouldFail: true,
			request: &pb.CreateAccountRequest{
				OwnerId:  "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				Id:       "a0eebc99-9c0b-4ef8",
				Currency: pb.Currency_CURRENCY_USD,
			},
			expectErrorString: "Invalid argument provided",
		},
		{
			name:       "should fail when currency is invalid",
			shouldFail: true,
			request: &pb.CreateAccountRequest{
				OwnerId:  "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				Id:       "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				Currency: pb.Currency_CURRENCY_UNSPECIFIED,
			},
			expectErrorString: "Invalid argument provided",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := client.CreateAccount(context.Background(), tt.request)
			if tt.shouldFail {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErrorString)
				return
			}
			require.NoError(t, err)
			resp, err := client.GetAccount(context.Background(), &pb.GetAccountRequest{Id: tt.request.Id})
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, tt.request.Id, resp.Id)
			require.Equal(t, tt.request.OwnerId, resp.OwnerId)
			require.Equal(t, int64(0), resp.Balance)
			require.Equal(t, tt.request.Currency, resp.Currency)
			require.NotEmpty(t, resp.CreatedAt)
			require.NotEmpty(t, resp.UpdatedAt)
			createdTime, err := time.Parse(time.RFC3339, resp.CreatedAt)
			require.NoError(t, err)
			updatedTime, err := time.Parse(time.RFC3339, resp.UpdatedAt)
			require.NoError(t, err)
			require.WithinDuration(t, time.Now(), createdTime, time.Second*2)
			require.WithinDuration(t, time.Now(), updatedTime, time.Second*2)
		})
	}
}

func TestTopUpAccount(t *testing.T) {
	client := newAccountServiceClient(t)
	// PREPARE
	_, err := client.CreateAccount(context.Background(), &pb.CreateAccountRequest{
		OwnerId:  "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a88",
		Id:       "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a88",
		Currency: pb.Currency_CURRENCY_EUR,
	})
	require.NoError(t, err)
	testCases := []struct {
		name              string
		shouldFail        bool
		expectErrorString string
		request           *pb.TopUpAccountRequest
	}{
		{
			name:       "should top up account",
			shouldFail: false,
			request: &pb.TopUpAccountRequest{
				Id:     "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a88",
				Amount: 100,
			},
		},
		{
			name:       "should fail when account id is invalid",
			shouldFail: true,
			request: &pb.TopUpAccountRequest{
				Id:     "a0eebc99-9c0b-4ef8",
				Amount: 100,
			},
			expectErrorString: "Invalid argument provided",
		},
		{
			name:       "should fail when not found",
			shouldFail: true,
			request: &pb.TopUpAccountRequest{
				Id:     "a0eebc99-9c0b-4ef8-bb6d-6bb9bd111111",
				Amount: 100,
			},
			expectErrorString: "not found",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// ACT
			_, err := client.TopUpAccount(context.Background(), tt.request)

			// ASSERT
			if tt.shouldFail {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErrorString)
				return
			}
			require.NoError(t, err)
			resp, err := client.GetAccount(context.Background(), &pb.GetAccountRequest{Id: tt.request.Id})
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, tt.request.Id, resp.Id)
			require.Equal(t, tt.request.Amount, resp.Balance)
			timeUpdated, err := time.Parse(time.RFC3339, resp.UpdatedAt)
			require.WithinDuration(t, time.Now(), timeUpdated, time.Second)
		})
	}
}

func TestTransfer(t *testing.T) {
	client := newAccountServiceClient(t)
	id1 := "a011bc99-9c0b-4ef8-bb6d-6bb9bd380a88"
	id2 := "b011bc99-9c0b-4ef8-bb6d-6bb9bd380a89"
	balance1 := int64(100)
	// PREPARE
	_, err := client.CreateAccount(context.Background(), &pb.CreateAccountRequest{
		OwnerId:  "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a88",
		Id:       id1,
		Currency: pb.Currency_CURRENCY_EUR,
	})
	require.NoError(t, err)
	_, err = client.CreateAccount(context.Background(), &pb.CreateAccountRequest{
		OwnerId:  "c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a88",
		Id:       id2,
		Currency: pb.Currency_CURRENCY_EUR,
	})
	require.NoError(t, err)
	_, err = client.TopUpAccount(context.Background(), &pb.TopUpAccountRequest{
		Id:     id1,
		Amount: balance1,
	})
	require.NoError(t, err)
	testCases := []struct {
		name              string
		shouldFail        bool
		expectErrorString string
		request           *pb.TransferFundsRequest
	}{
		{
			name:       "should transfer funds",
			shouldFail: false,
			request: &pb.TransferFundsRequest{
				FromAccountId: id1,
				ToAccountId:   id2,
				Amount:        43,
			},
		},
		{
			name: "should fail when to account id is invalid",
			request: &pb.TransferFundsRequest{
				FromAccountId: id1,
				ToAccountId:   "b0eebc99-9c0b-4ef8",
				Amount:        43,
			},
			expectErrorString: "Invalid argument provided",
			shouldFail:        true,
		},
		{
			name: "should fail when from  account id is invalid",
			request: &pb.TransferFundsRequest{
				FromAccountId: "a0eebc99-9c0b-4ef8",
				ToAccountId:   id2,
				Amount:        43,
			},
			expectErrorString: "Invalid argument provided",
			shouldFail:        true,
		},
		{
			name: "should fail when not found",
			request: &pb.TransferFundsRequest{
				FromAccountId: id1,
				ToAccountId:   "11111199-9c0b-4ef8-bb6d-6bb9bd111111",
			},
			expectErrorString: "not found",
			shouldFail:        true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// ACT
			_, err := client.TransferFunds(context.Background(), tt.request)

			// ASSERT
			if tt.shouldFail {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErrorString)
				return
			}
			require.NoError(t, err)
			// Check balance of account 1
			resp, err := client.GetAccount(context.Background(), &pb.GetAccountRequest{Id: tt.request.FromAccountId})
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, balance1-tt.request.Amount, resp.Balance)
			timeUpdated, err := time.Parse(time.RFC3339, resp.UpdatedAt)
			require.WithinDuration(t, time.Now(), timeUpdated, time.Second)

			// Check balance of account 2
			resp, err = client.GetAccount(context.Background(), &pb.GetAccountRequest{Id: tt.request.ToAccountId})
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, tt.request.Amount, resp.Balance)
			timeUpdated, err = time.Parse(time.RFC3339, resp.UpdatedAt)
			require.WithinDuration(t, time.Now(), timeUpdated, time.Second)
		})
	}
}

func TestMain(m *testing.M) {
	callback, ok := startService()
	if !ok {
		os.Exit(1)
	}
	code := m.Run()
	callback()
	os.Exit(code)
}

func startService() (func(), bool) {
	pscont, err := testhelpers.CreatePostgresContainer(context.Background())
	if err != nil {
		log.Fatalf("failed to create postgres container: %v", err)
		return nil, false
	}
	app, callback, err := NewApplication(context.Background(), &config.Config{
		Service:     config.Service{},
		Prometheus:  config.Prometheus{},
		DatabaseDSN: pscont.ConnectionString,
	})
	g := run.Group{}
	if err != nil {
		log.Fatalf("failed to create application: %v", err)
		return nil, false
	}
	grpcserver.RunWithGroup(&g,
		accountServiceAddr,
		func(server *grpc.Server) {
			svc := ports.NewGRPCServer(app)
			reflection.Register(server)
			pb.RegisterAccountServiceServer(server, svc)
		},
	)
	go func() {
		err = g.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()

	return callback, true
}

func newAccountServiceClient(t *testing.T) pb.AccountServiceClient {
	t.Helper()

	conn, err := grpcclient.NewClient(accountServiceAddr, true, grpcclient.RetryConfig{
		MaxRetries:     10,
		InitialBackoff: time.Millisecond,
		MaxBackoff:     time.Minute * 2,
	})
	require.NoError(t, err)
	return pb.NewAccountServiceClient(conn)
}
