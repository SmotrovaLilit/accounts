package grpcclient

import (
	"fmt"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	// "time".
)

func NewClient(
	addr string,
	isInsecure bool,
	retryConfig RetryConfig,
) (*grpc.ClientConn, error) {
	limiter := newRateLimiter(40, 40)

	creds := credentials.NewClientTLSFromCert(nil, "")
	transportSecurity := grpc.WithTransportCredentials(creds)
	if isInsecure {
		transportSecurity = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	conn, err := grpc.NewClient(
		addr,
		transportSecurity,

		grpc.WithChainUnaryInterceptor(
			LoggingInterceptor(),
			RetryInterceptor(retryConfig),
			grpc_prometheus.UnaryClientInterceptor,
			limiter.unaryRateLimiterInterceptor,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err)
	}
	return conn, nil
}
