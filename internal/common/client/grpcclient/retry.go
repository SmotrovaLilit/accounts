package grpcclient

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RetryConfig struct {
	MaxRetries     int           `json:"max_retries"`
	InitialBackoff time.Duration `json:"initial_backoff"`
	MaxBackoff     time.Duration `json:"max_backoff"`
}

var (
	retriableCodes = map[codes.Code]struct{}{
		codes.Internal:          {},
		codes.Unavailable:       {},
		codes.DeadlineExceeded:  {},
		codes.ResourceExhausted: {},
	}
)

func RetryInterceptor(config RetryConfig) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		var err error
		backoff := config.InitialBackoff

		for i := 0; i <= config.MaxRetries; i++ {
			err = invoker(ctx, method, req, reply, cc, opts...)
			if err == nil {
				return nil
			}
			if _, ok := retriableCodes[status.Code(err)]; !ok {
				// It is not a retriable error
				return err
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}

			slog.WarnContext(ctx, "Retry",
				"method", method,
				"attempt", i+1,
				"max_attempt", config.MaxRetries,
				"error", err,
			)
			backoff = time.Duration(float64(backoff) * 2)
			if backoff > config.MaxBackoff {
				backoff = config.MaxBackoff
			}
		}
		return err
	}
}
