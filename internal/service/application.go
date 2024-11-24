package service

import (
	"accounts/internal/app"
	"accounts/internal/config"
	"context"
	"os"
	"runtime/debug"
)

func NewApplication(_ context.Context, conf *config.Config) (a app.Application, callback func(), err error) {
	buildInfo, _ := debug.ReadBuildInfo()
	pid := os.Getpid()

	//accountRepository, err := adapters.NewPostgresSolutionRepository(conf.DatabaseDSN)
	//if err != nil {
	//	return app.Application{}, nil, err
	//}

	return app.Application{
			Commands: app.Commands{
				// TODO  add commands
			},
			Queries: app.Queries{
				// TODO add queries
			},
			Info: app.Info{
				Name:      "accounts",
				Version:   "v1beta1",
				BuildPath: buildInfo.Main.Path,
				GoVersion: buildInfo.GoVersion,
				PID:       pid,
			},
		}, func() {
			// TODO close resources
		}, nil
}
