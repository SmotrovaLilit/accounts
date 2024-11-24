package grpcclient

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"accounts/internal/common/logging"
)

func LoggingInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		md, _ := metadata.FromOutgoingContext(ctx)

		ctx = logging.AppendCtx(
			ctx,
			slog.String("grpc_client_method", method),
			slog.Any("grpc_client_metadata", md),
			slog.Any("grpc_client_request", req),
		)

		slog.DebugContext(ctx, "Starting client request")

		err := invoker(ctx, method, req, reply, cc, opts...)

		if err != nil {
			slog.ErrorContext(ctx, "Client request failed", slog.Any("error", err))
		} else {
			slog.InfoContext(ctx, "Client request executed successfully", slog.Any("client_grpc_response", reply))
		}

		return err
	}
}
