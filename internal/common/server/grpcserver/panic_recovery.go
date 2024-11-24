package grpcserver

import (
	"accounts/internal/common/safe"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func panicRecoveryHandler(ctx context.Context, p interface{}) (err error) {
	safe.DefaultRecover(ctx, p)
	return status.Errorf(codes.Internal, "internal server error")
}
