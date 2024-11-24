package app

import (
	"accounts/internal/common/logging"
	"context"
	"log/slog"
)

type Application struct {
	Commands Commands
	Queries  Queries
	Info     Info
}

type Info struct {
	Name      string
	Version   string
	BuildPath string
	GoVersion string
	PID       int
}

type Commands struct {
}

type Queries struct {
}

func (ap Application) AppendCtxWithApplicationLoggingFields(ctx context.Context) context.Context {
	return logging.AppendCtx(
		ctx,
		slog.Any("app", ap.Info.Name),
		slog.Any("app_version", ap.Info.Version),
		slog.Any("build_path", ap.Info.BuildPath),
		slog.Any("go_version", ap.Info.GoVersion),
	)
}
