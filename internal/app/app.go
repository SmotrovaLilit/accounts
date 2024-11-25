package app

import (
	"context"
	"log/slog"

	"accounts/internal/app/commands"
	"accounts/internal/app/query"
	"accounts/internal/common/logging"
)

// Application is an application.
type Application struct {
	Commands Commands
	Queries  Queries
	Info     Info
}

// Info is an application info.
type Info struct {
	Name      string
	Version   string
	BuildPath string
	GoVersion string
	PID       int
}

// Commands is a set of commands.
type Commands struct {
	CreateAccount commands.CreateAccountHandler
	TopUpAccount  commands.TopUpAccountHandler
	TransferFunds commands.TransferFundsHandler
}

// Queries is a set of queries.
type Queries struct {
	GetAccount query.GetAccountQueryHandler
}

// AppendCtxWithApplicationLoggingFields appends application logging fields to the context.
func (ap Application) AppendCtxWithApplicationLoggingFields(ctx context.Context) context.Context {
	return logging.AppendCtx(
		ctx,
		slog.Any("app", ap.Info.Name),
		slog.Any("app_version", ap.Info.Version),
		slog.Any("build_path", ap.Info.BuildPath),
		slog.Any("go_version", ap.Info.GoVersion),
	)
}
