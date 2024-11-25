package grpcclient

import (
	"context"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

// rateLimiter struct to hold the rate limiter.
type rateLimiter struct {
	limiter *rate.Limiter
}

// newRateLimiter creates a new rateLimiter
// limit is the rate limit in requests per second
// burst is the burst limit, the maximum number of requests that can be processed in a short period of time
// Set limit to the desired requests per second (e.g., 100 rps).
// Set burst to control temporary spikes in traffic:
// - Set burst equal to limit for strict rate limiting, ensuring consistent request pacing.
// - Set burst greater than limit to allow short traffic spikes, smoothing out bursty workloads.
func newRateLimiter(limit, burst int) *rateLimiter {
	return &rateLimiter{
		limiter: rate.NewLimiter(rate.Limit(limit), burst),
	}
}

func (rl *rateLimiter) unaryRateLimiterInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := rl.limiter.Wait(ctx)
	if err != nil {
		// Possible errors: context.Canceled, context.DeadlineExceeded, unknown
		return err
	}
	err = invoker(ctx, method, req, reply, cc, opts...)
	return err
}
