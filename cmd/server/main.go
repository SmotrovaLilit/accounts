package main

import (
	"accounts/internal/common/logging"
	"accounts/internal/common/prometheus"
	"accounts/internal/common/server/grpcserver"
	"accounts/internal/config"
	pb "accounts/internal/gen-proto/accountservice/v1beta"
	"accounts/internal/ports"
	"accounts/internal/service"
	"context"
	"fmt"
	"github.com/oklog/run"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	logger := logging.NewLogger(os.Stdout, true, slog.LevelDebug)
	slog.SetDefault(logger)

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	application, callback, err := service.NewApplication(ctx, conf)
	if err != nil {
		log.Fatalf("failed to create application: %v", err)
	}
	defer callback()
	ctx = application.AppendCtxWithApplicationLoggingFields(ctx)
	g := &run.Group{}

	// TODO move part of this logic to a common package because implementation
	// of running and gracefully shutting down a server is common for different microservices
	prometheus.RunWithGroup(g, fmt.Sprintf(":%d", conf.Prometheus.Port))
	grpcserver.RunWithGroup(g,
		fmt.Sprintf("%s:%d", conf.ChallengeService.Host, conf.ChallengeService.Port),
		func(server *grpc.Server) {
			svc := ports.NewGRPCServer(application)
			reflection.Register(server)
			pb.RegisterAccountServiceServer(server, svc)
		},
	)

	sigs := make(chan os.Signal)
	g.Add(func() error {
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
		sig := <-sigs
		slog.Info("received signal", slog.Any("signal", sig))
		return nil
	}, func(error) {
		signal.Stop(sigs)
		close(sigs)
	})

	err = g.Run()
	if err != nil {
		log.Fatal(err)
	}
}
