package service

import (
	"context"
	"os"
	"runtime/debug"

	"accounts/internal/adapters"
	"accounts/internal/app"
	"accounts/internal/app/commands"
	"accounts/internal/app/query"
	"accounts/internal/config"
)

func NewApplication(_ context.Context, conf *config.Config) (a app.Application, callback func(), err error) {
	buildInfo, _ := debug.ReadBuildInfo()
	pid := os.Getpid()

	accountRepository, err := adapters.NewPostgresAccountRepository(conf.DatabaseDSN)
	if err != nil {
		return app.Application{}, nil, err
	}

	return app.Application{
			Commands: app.Commands{
				CreateAccount: commands.NewCreateAccountHandler(accountRepository),
				TopUpAccount:  commands.NewTopUpAccountHandler(accountRepository),
				TransferFunds: commands.NewTransferFundsHandler(accountRepository),
			},
			Queries: app.Queries{
				GetAccount: query.NewGetAccountQueryHandler(accountRepository),
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
