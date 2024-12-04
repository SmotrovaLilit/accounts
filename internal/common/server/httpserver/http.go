package httpserver

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/oklog/run"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-chi/chi/v5"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Server struct {
	addr           string
	server         *http.Server
	registerServer func(router chi.Router) http.Handler
}

func (s *Server) Run() error {
	apiRouter := chi.NewRouter()
	setMiddlewares(apiRouter)

	rootRouter := chi.NewRouter()
	// we are mounting all APIs under /api path
	rootRouter.Mount("/api", s.registerServer(apiRouter))
	s.server = &http.Server{
		Addr:              s.addr,
		Handler:           promhttp.Handler(),
		ReadHeaderTimeout: 5 * time.Second, // TODO make configurable
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	slog.InfoContext(ctx, "shutting down http server gracefully")
	err := s.server.Shutdown(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to gracefully shut down http server", slog.Any("error", err))
		return err
	}
	slog.InfoContext(ctx, "http server shut down gracefully")
	return nil
}

// RunWithGroup runs the HTTP server with the provided run.Group.
func RunWithGroup(g *run.Group, addr string, registerServer func(router chi.Router) http.Handler) {
	srv := Server{addr: addr, registerServer: registerServer}
	g.Add(
		srv.Run,
		func(error) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = srv.Shutdown(ctx)
		},
	)
}

func setMiddlewares(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	// TODO setup logging middleware
	// router.Use()
	router.Use(middleware.Recoverer)

	addCorsMiddleware(router)

	router.Use(
		middleware.SetHeader("X-Content-Type-Options", "nosniff"),
		middleware.SetHeader("X-Frame-Options", "deny"),
	)
	router.Use(middleware.NoCache)
}

func addCorsMiddleware(router *chi.Mux) {
	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ";")
	if len(allowedOrigins) == 0 {
		return
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(corsMiddleware.Handler)
}
