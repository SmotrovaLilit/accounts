package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/oklog/run"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type Server struct {
	Addr           string
	RegisterServer func(server *grpc.Server)
	Opts           []Option

	server *grpc.Server
}

func (s *Server) Run() error {
	resultOptions := options{}
	for _, opt := range s.Opts {
		opt.apply(&resultOptions)
	}
	// TODO set rate limit from config
	rate := newRateLimiter(50, 50)
	s.server = grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandlerContext(panicRecoveryHandler)),
			newContextFactoryUnaryServerInterceptor(resultOptions.contextFactory),
			tracingUnaryServerInterceptor,
			loggingUnaryServerInterceptor,
			grpc_prometheus.UnaryServerInterceptor,
			rate.unaryRateLimiterInterceptor,
		),
	)

	s.RegisterServer(s.server)
	grpc_prometheus.Register(s.server)
	grpc_prometheus.EnableHandlingTimeHistogram()

	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen %s: %w", s.Addr, err)
	}

	slog.Info("starting grpc server", slog.String("address", s.Addr))
	err = s.server.Serve(listener)
	if errors.Is(err, grpc.ErrServerStopped) {
		return nil
	}
	if err != nil {
		slog.Error("failed to start grpc server", slog.Any("error", err))
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("shutting down grpc server gracefully")
	done := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		slog.Info("grpc server shut down gracefully")
		return nil
	case <-ctx.Done():
		slog.Info("grpc server graceful shutdown timed out, forcing shutdown")
		s.server.Stop()
		return nil
	}
}

func RunWithGroup(
	g *run.Group,
	addr string,
	registerServer func(server *grpc.Server),
	opts ...Option,
) {
	grpcSrv := Server{
		Addr:           addr,
		RegisterServer: registerServer,
		Opts:           opts,
	}
	g.Add(
		func() error { return grpcSrv.Run() },
		func(error) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = grpcSrv.Shutdown(ctx)
		},
	)
}
