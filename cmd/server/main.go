package main

import (
	"accounts/internal/common/server/httpserver"
	"accounts/internal/ports/httpports"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"

	"accounts/internal/common/logging"
	"accounts/internal/config"
	"accounts/internal/service"

	_ "github.com/lib/pq"
	"github.com/oklog/run"
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
	//prometheus.RunWithGroup(g, fmt.Sprintf(":%d", conf.Prometheus.Port))
	// grpcserver.RunWithGroup(g,
	//	fmt.Sprintf("%s:%d", conf.Service.Host, conf.Service.Port),
	//	func(server *grpc.Server) {
	//		svc := grpc2.NewGRPCServer(application)
	//		reflection.Register(server)
	//		pb.RegisterAccountServiceServer(server, svc)
	//	},
	//)
	httpserver.RunWithGroup(g, fmt.Sprintf("%s:%d", conf.HTTPService.Host, conf.HTTPService.Port), func(router chi.Router) http.Handler {
		svc := httpports.NewHTTPServer(application)
		router.Post("/accounts", svc.CreateAccount)
		router.Get("/accounts/{id}", svc.GetAccount)
		return router
	})

	sigs := make(chan os.Signal, 1)
	g.Add(func() error {
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
		sig := <-sigs
		slog.InfoContext(ctx, "received signal", slog.Any("signal", sig))
		return nil
	}, func(error) {
		signal.Stop(sigs)
		close(sigs)
	})

	err = g.Run()
	if err != nil {
		slog.ErrorContext(ctx, "application encountered an error", slog.Any("error", err))
	}
}
