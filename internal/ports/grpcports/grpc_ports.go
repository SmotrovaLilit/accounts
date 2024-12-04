package grpcports

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"accounts/internal/app"
	"accounts/internal/app/commands"
	"accounts/internal/app/query"
	"accounts/internal/domain/account"
	"accounts/internal/domain/owner"
	pb "accounts/internal/gen-proto/accountservice/v1beta"
)

// Server is a gRPC server.
type Server struct {
	app app.Application
	pb.UnimplementedAccountServiceServer
}

// NewGRPCServer creates a new gRPC server.
func NewGRPCServer(application app.Application) Server {
	return Server{
		app: application,
	}
}

// CreateAccount creates an account.
func (s Server) CreateAccount(ctx context.Context, r *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	id, err := account.NewID(r.Id)
	if err != nil {
		return nil, toGrpcError(err)
	}
	ownerID, err := owner.NewID(r.OwnerId)
	if err != nil {
		return nil, toGrpcError(err)
	}
	currency, err := toCurrency(r.Currency)
	if err != nil {
		return nil, toGrpcError(err)
	}
	err = s.app.Commands.CreateAccount.Handle(ctx, commands.CreateAccount{
		ID:       id,
		OwnerID:  ownerID,
		Currency: currency,
	})
	if err != nil {
		return nil, toGrpcError(err)
	}
	return &pb.CreateAccountResponse{}, nil
}

// GetAccount gets an account.
func (s Server) GetAccount(ctx context.Context, r *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	id, err := account.NewID(r.Id)
	if err != nil {
		return nil, toGrpcError(err)
	}
	acc, err := s.app.Queries.GetAccount.Handle(ctx, query.GetAccount{
		ID: id,
	})
	if err != nil {
		return nil, toGrpcError(err)
	}
	return &pb.GetAccountResponse{
		Id:        acc.ID().String(),
		OwnerId:   acc.OwnerID().String(),
		Currency:  toResponseCurrency(acc.BalanceCurrency()),
		Balance:   acc.BalanceAmount(),
		CreatedAt: formatDate(acc.CreatedAt()),
		UpdatedAt: formatDate(acc.UpdatedAt()),
	}, nil
}

// TopUpAccount tops up an account.
func (s Server) TopUpAccount(ctx context.Context, r *pb.TopUpAccountRequest) (*pb.TopUpAccountResponse, error) {
	id, err := account.NewID(r.Id)
	if err != nil {
		return nil, toGrpcError(err)
	}
	err = s.app.Commands.TopUpAccount.Handle(ctx, commands.TopUpAccount{
		ID:     id,
		Amount: r.Amount,
	})
	if err != nil {
		return nil, toGrpcError(err)
	}
	return &pb.TopUpAccountResponse{}, nil
}

func (s Server) TransferFunds(ctx context.Context, r *pb.TransferFundsRequest) (*pb.TransferFundsResponse, error) {
	fromID, err := account.NewID(r.FromAccountId)
	if err != nil {
		return nil, toGrpcError(err)
	}
	toID, err := account.NewID(r.ToAccountId)
	if err != nil {
		return nil, toGrpcError(err)
	}
	err = s.app.Commands.TransferFunds.Handle(ctx, commands.TransferFunds{
		From:   fromID,
		To:     toID,
		Amount: r.Amount,
	})
	if err != nil {
		return nil, toGrpcError(err)
	}
	return &pb.TransferFundsResponse{}, nil
}

// formatDate formats a time.Time to a string ISO 8601.
func formatDate(t time.Time) string {
	return t.Format(time.RFC3339)
}

func newInvalidArgumentError(code pb.ErrorCode, message string) error {
	st := status.New(codes.InvalidArgument, "Invalid argument provided")
	errorDetails := &pb.ErrorDetails{
		Code:    code,
		Message: message,
	}
	stWithDetails, err := st.WithDetails(errorDetails)
	if err != nil {
		return status.Error(codes.Internal, "Internal error")
	}
	return stWithDetails.Err()
}

// toGrpcError converts domain errors to gRPC errors.
// TODO think about how better cover by tests.
func toGrpcError(err error) error {
	switch {
	// Examples of how to handle specific errors
	case errors.Is(err, account.ErrInvalidID):
		return newInvalidArgumentError(pb.ErrorCode_ERROR_CODE_INVALID_ACCOUNT_ID, err.Error())
	case errors.Is(err, owner.ErrInvalidID):
		return newInvalidArgumentError(pb.ErrorCode_ERROR_CODE_INVALID_OWNER_ID, err.Error())
	case errors.As(err, &ErrInvalidCurrency{}):
		return newInvalidArgumentError(pb.ErrorCode_ERROR_CODE_INVALID_CURRENCY, err.Error())
	case errors.As(err, &account.NotFoundError{}):
		return status.Error(codes.NotFound, err.Error())
	// Examples of how to handle generic errors
	case errors.Is(err, account.ErrCannotAddRemoveNegativeAmount):
	case errors.Is(err, account.ErrAmountOverflow):
	case errors.As(err, &account.ErrCurrencyMismatch{}):
	case errors.Is(err, account.ErrNotEnoughBalance):
	case errors.Is(err, account.ErrInvalidAmount):
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return status.Error(codes.Internal, err.Error()) // TODO should hide internal error message
}

func toCurrency(currency pb.Currency) (account.Currency, error) {
	switch currency {
	case pb.Currency_CURRENCY_USD:
		return account.CurrencyUSD, nil
	case pb.Currency_CURRENCY_EUR:
		return account.CurrencyEUR, nil
	case pb.Currency_CURRENCY_UNSPECIFIED:
		return account.Currency{}, ErrInvalidCurrency{Currency: currency.String()}
	default:
		return account.Currency{}, ErrInvalidCurrency{Currency: currency.String()}
	}
}

func toResponseCurrency(currency account.Currency) pb.Currency {
	switch currency {
	case account.CurrencyUSD:
		return pb.Currency_CURRENCY_USD
	case account.CurrencyEUR:
		return pb.Currency_CURRENCY_EUR
	default:
		return pb.Currency_CURRENCY_UNSPECIFIED
	}
}

// ErrInvalidCurrency is returned when an invalid currency is provided.
type ErrInvalidCurrency struct {
	Currency string
}

// Error returns the error message.
func (e ErrInvalidCurrency) Error() string {
	return fmt.Sprintf("invalid currency: %s", e.Currency)
}
