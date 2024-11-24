package prometheus

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	MetricsAddr string
	server      *http.Server
}

func (p *Server) Run() error {
	p.server = &http.Server{
		Addr:    p.MetricsAddr,
		Handler: promhttp.Handler(),
	}

	slog.Info("starting prometheus metrics server", slog.String("address", p.MetricsAddr))
	err := p.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return err
	}
	if err != nil {
		slog.Error("failed to start prometheus server", slog.Any("error", err))
		return err
	}
	return nil
}

func (p *Server) Shutdown(ctx context.Context) error {
	slog.Info("shutting down prometheus server gracefully")
	err := p.server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to gracefully shut down prometheus server", slog.Any("error", err))
		return err
	}
	slog.Info("prometheus server shut down gracefully")
	return nil
}

func RunWithGroup(g *run.Group, metricsAddr string) {
	prometheusSrv := Server{MetricsAddr: metricsAddr}
	g.Add(
		func() error { return prometheusSrv.Run() },
		func(error) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = prometheusSrv.Shutdown(ctx)
		},
	)
}
