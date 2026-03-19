package runners

import (
	"context"

	"github.com/not-for-prod/clay/server"
	"go.uber.org/fx"

	"github.com/dehwyy/x-balance/internal/config"
	balancehandler "github.com/dehwyy/x-balance/internal/delivery/api/balance"
	transactionhandler "github.com/dehwyy/x-balance/internal/delivery/api/transaction"
	userhandler "github.com/dehwyy/x-balance/internal/delivery/api/user"
	balancev1 "github.com/dehwyy/x-balance/internal/generated/pb/balance/v1"
	transactionsv1 "github.com/dehwyy/x-balance/internal/generated/pb/transactions/v1"
	usersv1 "github.com/dehwyy/x-balance/internal/generated/pb/users/v1"
)

type RunAPIOpts struct {
	fx.In
	LC           fx.Lifecycle
	Shutdowner   fx.Shutdowner
	Config       *config.Config
	UserH        *userhandler.Handler
	BalanceH     *balancehandler.Handler
	TransactionH *transactionhandler.Handler
}

func RunAPI(opts RunAPIOpts) error {
	srv := server.NewServer(opts.Config.Port)

	opts.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.Run(
					usersv1.NewUserServiceServiceDesc(opts.UserH),
					balancev1.NewBalanceServiceServiceDesc(opts.BalanceH),
					transactionsv1.NewTransactionServiceServiceDesc(opts.TransactionH),
				); err != nil {
					_ = opts.Shutdowner.Shutdown(fx.ExitCode(1))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Stop(ctx)
		},
	})
	return nil
}
