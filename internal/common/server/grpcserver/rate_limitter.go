package grpcserver

import (
	"context"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

// rateLimiter struct to hold the rate limiter
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

// unaryRateLimiterInterceptor intercepts unary requests for rate limiting
func (rl *rateLimiter) unaryRateLimiterInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	// Rate limiting
	if !rl.limiter.Allow() {
		return nil, grpcstatus.Errorf(
			grpccodes.ResourceExhausted,
			"rate limit exceeded: limit: %f, burst: %d",
			rl.limiter.Limit(), rl.limiter.Burst(),
		)
	}
	return handler(ctx, req)
}
