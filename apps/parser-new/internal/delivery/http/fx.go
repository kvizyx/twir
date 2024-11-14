package http

import (
	"context"
	"log/slog"

	"go.uber.org/fx"
)

func NewServerFx(params Params, lc fx.Lifecycle) Server {
	server := NewServer(params)

	lc.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				go func() {
					if err := server.Start(); err != nil {
						params.Logger.Error(
							"failed to start http server",
							slog.Any("error", err),
						)
					}
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				return server.Stop(ctx)
			},
		},
	)

	return server
}
