package main

import (
	"context"
	"fmt"
	"net"

	"github.com/dehwyy/txmanagerfx/pkg/txmanager"
	"github.com/dehwyy/x-balance/internal/config"
	balancehandler "github.com/dehwyy/x-balance/internal/delivery/api/balance"
	transactionhandler "github.com/dehwyy/x-balance/internal/delivery/api/transaction"
	userhandler "github.com/dehwyy/x-balance/internal/delivery/api/user"
	balancepb "github.com/dehwyy/x-balance/internal/generated/pb"
	"github.com/dehwyy/x-balance/internal/runners"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func main() {
	app := fx.New(
		fx.Provide(
			config.Load,
			runners.NewDB,
			runners.NewRedis,
			func() zerolog.Logger { return log.Logger },
		),
		txmanager.NewGorm(),
		runners.ReposModule,
		runners.GatewayModule,
		runners.ServicesModule,
		runners.WorkersModule,
		fx.Provide(
			userhandler.New,
			balancehandler.New,
			transactionhandler.New,
		),
		fx.Invoke(func(
			lc fx.Lifecycle,
			userH *userhandler.Handler,
			balanceH *balancehandler.Handler,
			txH *transactionhandler.Handler,
			cfg *config.Config,
			logger zerolog.Logger,
		) {
			srv := grpc.NewServer()
			balancepb.RegisterUserServiceServer(srv, userH)
			balancepb.RegisterBalanceServiceServer(srv, balanceH)
			balancepb.RegisterTransactionServiceServer(srv, txH)

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
					if err != nil {
						return err
					}
					go func() {
						if err := srv.Serve(lis); err != nil {
							logger.Error().Err(err).Msg("grpc server error")
						}
					}()
					logger.Info().Int("port", cfg.Port).Msg("gRPC server started")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					srv.GracefulStop()
					return nil
				},
			})
		}),
	)

	app.Run()
}
